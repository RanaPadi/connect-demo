#!/bin/bash

# Path to the config.yaml file
config_yaml_file="$PWD/../python/utils/config.yaml"

# Extracting environment variables from the config.yaml file
CONNECT_DEMO_PATH=$(grep 'CONNECT_DEMO_PATH' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
vehicle_computer_username=$(grep 'vehicle_computer_username' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
vehicle_computer_1_ip=$(grep 'vehicle_computer_1_ip' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
vehicle_computer_2_ip=$(grep 'vehicle_computer_2_ip' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
vehicle_computer_password=$(grep 'vehicle_computer_password' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
log_path=$(grep 'log_path' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')

# Log file path
logfiletimestamp=$(date +"%Y-%m-%d_%H:%M:%S")
log_file="migration_$logfiletimestamp.log"
benchmark_log_file="$log_path/benchmark.log"

# Function to check whether the log directory exists
check_log_directory() {
    if [ ! -d "$log_path" ]; then
        mkdir -p "$log_path"
        chmod 755 "$log_path"
    fi
}

# Function to log messages with timestamp
log_message() {
    local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
    local message="$1"
    echo "$message"
    echo "[$timestamp] $message" >> "$log_path/$log_file"
}

# Redirect stdout to the log file
exec > >(while read -r line; do log_message "$line"; done)

# Function to check for sudo permission
check_sudo_permission() {
    local password="$1"
    echo "$password" | sudo -Sv
    if ! sudo -v ; then
        echo "Sorry, unable to get sudo permission. Exiting Script Execution."
        exit 1
    fi
}

# Function to check and install sshpass
check_and_install_sshpass() {
if ! command -v sshpass &> /dev/null
then
    echo "sshpass could not be found. Installing sshpass..."
    check_sudo_permission "$vehicle_computer_password"
    sudo apt-get install sshpass
    if [ $? -eq 0 ]; then
        echo "sshpass installed successfully."
    else
        echo "Failed to install sshpass. Exiting the script execution..."
        exit 1
    fi
fi
}

# Function to enable applications on the Vehicle Computer machines
enable_disable_applications() {
    # Check if username and hostname/IP are provided
    if [ $# -ne 4 ]; then
        echo "Usage: enable_applications <username> <hostname_or_ip>"
        return 1
    fi

    # Assign input parameters to variables
    vehicle_computer_username="$1"
    vehicle_computer_ip="$2"
    color="$3"
    acc_python_script_option="$4"

    # Connect Demo path on the Vehicle Computer
    remote_connect_demo_path=$CONNECT_DEMO_PATH

    # Python Setup Script file path
    python_setup_script="$remote_connect_demo_path/bash_scripts/python_setup.sh"

    # USB LED Control Script file path
    usb_led_control_script="$remote_connect_demo_path/bash_scripts/usb_led_control.sh"

    # Execute SSH command to enable applications
    ssh_application_command="$python_setup_script $acc_python_script_option"
    ssh_led_command="$usb_led_control_script $color"
    #ssh "$vehicle_computer_username@$vehicle_computer_ip" "$ssh_application_command"
    sshpass -p "$vehicle_computer_password" ssh "$vehicle_computer_username@$vehicle_computer_ip" "$ssh_application_command"
    #ssh "$vehicle_computer_username@$vehicle_computer_ip" "$ssh_led_command"
    sshpass -p "$vehicle_computer_password" ssh "$vehicle_computer_username@$vehicle_computer_ip" "$ssh_led_command"
}

# Function to log the duration of the migration process
log_migration_duration() {
    local start_time=$1
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    echo "$(date +"%Y-%m-%d %H:%M:%S,%3N") - INFO - Benchmark 8: ${duration}.0ms." | tee -a "$benchmark_log_file"
}

# Function to display the user options
display_user_options() {
    echo "Migration Options:"
    echo "VC1: Migrate applications to Vehicle Computer 1."
    echo "VC2: Migrate applications to Vehicle Computer 2."
    echo "Enter your choice:"
}

# Function to invoke migration of the Kafka Python applications
migration() {
    local start_time=$(date +%s)
    if [ "$1" == "VC1" ]; then
        echo "Migrating applications to Vehicle Computer 1..."
        enable_disable_applications "$vehicle_computer_username" "$vehicle_computer_1_ip" "g" "eac"
        enable_disable_applications "$vehicle_computer_username" "$vehicle_computer_2_ip" "r" "dac"
        echo "Migration of applications to Vehicle Computer 1 completed."
    elif [ "$1" == "VC2" ]; then
        echo "Migrating applications to Vehicle Computer 2..."
        enable_disable_applications "$vehicle_computer_username" "$vehicle_computer_2_ip" "g" "eac"
        enable_disable_applications "$vehicle_computer_username" "$vehicle_computer_1_ip" "r" "dac"
        echo "Migration of applications to Vehicle Computer 2 completed."
    else
        echo "Invalid input. Please provide the correct Vehicle Computer name."
        return
    fi
    log_migration_duration "$start_time"
}

# Main function to migrate the Kafka Python applications
migration_main() {
    check_log_directory
    check_and_install_sshpass
    if [ "$#" -ne 1 ]; then
        display_user_options
        read -r migration_user_choice
        migration "$migration_user_choice"
    else    
        migration "$1"
    fi
}

# Call the main function
migration_main "$@"