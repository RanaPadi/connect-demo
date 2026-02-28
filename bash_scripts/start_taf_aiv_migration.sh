#!/bin/bash

# Path to the config.yaml file
config_yaml_file="$PWD/../python/utils/config.yaml"

# Extracting environment variables from the config.yaml file
python_migration_script=$(grep 'python_migration_script' "$config_yaml_file" | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
python_venv_path=$(grep 'python_venv_path' "$config_yaml_file" | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
log_path=$(grep 'log_path' "$config_yaml_file" | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')

# Log file path
logfiletimestamp=$(date +"%Y-%m-%d_%H:%M:%S")
log_file="$log_path/start_taf_aiv_migration_$logfiletimestamp.log"

# Function to check whether the log directory exists
check_log_directory() {
    if [ ! -d "$log_path" ]; then
        mkdir -p "$log_path"
        chmod 755 "$log_path"
    fi
}

# Function to log messages with timestamp
log_message() {
    local timestamp
    timestamp=$(date +"%Y-%m-%d %H:%M:%S")
    local message="$1"
    echo "$message"
    echo "[$timestamp] $message" >> "$log_file"
}

# Redirect stdout to the log file
exec > >(while read -r line; do log_message "$line"; done) 2>&1

# Array to hold process IDs for cleanup
declare -a pids

# Cleanup function to terminate specific processes
cleanup() {
    echo "Cleaning up the processes..."
    if [ ${#pids[@]} -eq 0 ]; then
        echo "No processes to clean."
        return
    else
        for pid in "${pids[@]}"; do
            if ps -p "$pid" > /dev/null; then
                kill "$pid"
                sleep 2
                if ps -p "$pid" > /dev/null; then
                    kill -9 "$pid"
                    echo "Process $pid killed forcefully."
                else
                    echo "Process $pid killed."
                fi
            fi
        done
        echo "Cleanup completed."
    fi
}

# Function to start TAF process
start_taf() {
    local taf_dir="$HOME/ACT/taf"
    local taf_exec="$taf_dir/go-taf"
    if [[ -d $taf_dir && -f $taf_exec ]]; then
        echo "Starting TAF process..."
        cd "$taf_dir" || exit
        export TAF_CONFIG="$taf_dir/res/taf.json"
        ./go-taf &>/dev/null &
        local taf_pid=$!
        sleep 5
        if ps -p "$taf_pid" > /dev/null; then
            echo "TAF process started successfully."
            pids+=("$taf_pid")
        else
            echo "Failed to start TAF process."
            cleanup
            exit 1
        fi
    else
        echo "TAF setup error: Directory or executable not found."
        cleanup
        exit 1
    fi
}

# Function to start AIV process
start_aiv() {
    local aiv_dir="$HOME/ACT/aiv"
    local aiv_exec="$aiv_dir/aiv.py"
    if [[ -d $aiv_dir && -f $aiv_exec ]]; then
        echo "Starting AIV process..."
        cd "$aiv_dir" || exit
        python3 aiv.py --broker_ip=connect-kafka.euprojects.net:3020 --mode=mutable &>/dev/null &
        local aiv_pid=$!
        sleep 5
        if ps -p "$aiv_pid" > /dev/null; then
            echo "AIV process started successfully."
            pids+=("$aiv_pid")
        else
            echo "Failed to start AIV process."
            cleanup
            exit 1
        fi
    else
        echo "AIV setup error: Directory or executable not found."
        cleanup
        exit 1
    fi
}

# Function to copy keys for TAF/AIV communication
copy_keys() {
    local taf_dir="$HOME/ACT/taf/res/cert"
    local aiv_dir="$HOME/ACT/aiv"
    local taf_key_1="$taf_dir/ecdsa_public_key.pem"
    local taf_key_2="$taf_dir/attestationCertificate.pem"
    local aiv_key="$aiv_dir/aiv_public_key.pem"
    if [[ -f $taf_key_1 && -f $taf_key_2 && -f $aiv_key ]]; then
        echo "Copying keys for TAF/AIV communication..."
        cp "$taf_key_1" "$aiv_dir"
        cp "$taf_key_2" "$aiv_dir"
        cp "$aiv_key" "$taf_dir"
        echo "Keys copied successfully."
    else
        echo "Key files not found."
        cleanup
        exit 1
    fi
}

# Function to start the Migration App
start_migration_app() {
    local migration_app_exec="$python_migration_script"
    if [[ -f $migration_app_exec && -d $python_venv_path ]]; then
        echo "Starting Migration App..."
        source "$python_venv_path/bin/activate"
        python3 "$migration_app_exec" &>/dev/null &
        local migration_pid=$!
        sleep 5
        if ps -p "$migration_pid" > /dev/null; then
            echo "Migration App started successfully."
            pids+=("$migration_pid")
        else
            echo "Failed to start Migration App."
            cleanup
            exit 1
        fi
    else
        echo "Migration App setup error: Directory or executable or virtual environment not found."
        cleanup
        exit 1
    fi
}

# Kill TAF/AIV/Migration App processes
kill_processes() {
    local pids=($(pgrep -f "go-taf|aiv.py|migration_app.py"))
    if [ ${#pids[@]} -eq 0 ]; then
        echo "Nothing to kill."
        return
    fi

    echo "Killing TAF/AIV/Migration App processes..."
    for pid in "${pids[@]}"; do
        kill "$pid"
        sleep 2  # Allow some time for the process to terminate gracefully
        if ps -p "$pid" > /dev/null; then
            kill -9 "$pid"
            echo "Process $pid killed forcefully."
        else
            echo "Process $pid killed."
        fi
    done
    echo "All processes terminated."
}

main() {
    check_log_directory
    echo "Enter 1 to start the TAF/AIV/Migration App processes."
    echo "Enter 2 to kill the TAF/AIV/Migration App processes."
    read -p "Enter your choice: " choice
    case $choice in
        1)
            start_taf
            start_aiv
            copy_keys
            start_migration_app
            ;;
        2)
            kill_processes
            ;;
        *)
            echo "Invalid choice. Exiting..."
            exit 1
            ;;
    esac
}

# Execute the main function
main
