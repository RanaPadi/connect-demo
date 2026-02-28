#!/bin/bash

# Path to the config.yaml file
config_yaml_file="$PWD/../python/utils/config.yaml"

# Extracting environment variables from the config.yaml file
python_venv_path=$(grep 'python_venv_path' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
python_requirements_file=$(grep 'python_requirements_file' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
python_acc_application=$(grep 'python_acc_application' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
vehicle_computer_password=$(grep 'vehicle_computer_password' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
log_path=$(grep 'log_path' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')

# Log file path
logfiletimestamp=$(date +"%Y-%m-%d_%H:%M:%S")
log_file="python_setup.log"

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

# Function to check and install Python, PIP, venv, and Tkinter if not already installed
check_and_install_python() {
    # Define an associative array mapping package managers to their update and install commands
    declare -A pkg_managers=(
        ["apt-get"]="sudo apt-get update -y && sudo apt-get install -y"
        ["yum"]="sudo yum update -y && sudo yum install -y"
        ["dnf"]="sudo dnf update -y && sudo dnf install -y"
        ["zypper"]="sudo zypper refresh && sudo zypper install -y"
        ["pacman"]="sudo pacman -Syu --noconfirm"
        ["brew"]="brew update && brew install"
        ["pkg"]="sudo pkg update && sudo pkg install -y"
        ["port"]="sudo port selfupdate && sudo port install"
        ["emerge"]="sudo emerge --sync && sudo emerge -av"
    )

    # Define a list of required packages
    packages=("python3" "python3-pip" "python3-venv" "python3-tk")

    # Detect the available package manager
    for manager in "${!pkg_managers[@]}"; do
        if command -v $manager &>/dev/null; then
            install_cmd="${pkg_managers[$manager]}"
            break
        fi
    done

    if [ -z "$install_cmd" ]; then
        echo "Package manager not found. Please install the required packages manually."
        exit 1
    fi

    # Install or update the required packages
    for pkg in "${packages[@]}"; do
        if command -v dpkg &>/dev/null && dpkg -s "$pkg" &>/dev/null; then
            echo "$pkg is already installed. Updating $pkg..."
        else
            echo "$pkg is not installed. Installing $pkg..."
        fi
        eval "$install_cmd $pkg"
        echo "$pkg is successfully installed/updated."
    done

    return 0
}

# Function to check and create a Python Virtual environment for Kafka
check_and_create_python_virtualenv() {
  	if [ ! -d "$python_venv_path" ]; then
        echo "Python virtual environment for Kafka does not exist. Creating now..."
     	mkdir -p "$python_venv_path"
        chmod 755 "$python_venv_path"
        python3 -m venv "$python_venv_path"
        echo "Created Python virtual environment at $python_venv_path"
        if ! source "$python_venv_path/bin/activate"; then
            echo "Failed to activate Python virtual environment."
            return 1
        else
            deactivate
        fi
    else
        echo "Python virtual environment for Kafka already exists at: $python_venv_path"
        if ! source "$python_venv_path/bin/activate"; then
            echo "Failed to activate Python virtual environment."
            return 1
        else
            deactivate
        fi
  	fi
    return 0
}

#Function to activate Python Virtual Environment
activate_python_virtualenv() {
    if [ -f "$python_venv_path/bin/activate" ]; then
        if ! source "$python_venv_path/bin/activate"; then
            echo "Failed to activate Python virtual environment."
            return 1
        fi
        echo "Activated Python virtual environment at: $python_venv_path"
        return 0
    else
        echo "Python Virtual Environment is not available."
        return 1
    fi
}

# Function to deactivate Python Virtual Environment
deactivate_python_virtualenv() {
    if [ -f "$python_venv_path/bin/activate" ]; then
        if ! deactivate; then
            echo "Failed to deactivate Python virtual environment."
            return 1
        fi
        echo "Deactivated Python Virtual Environment at: $python_venv_path"
        return 0
    else
        echo "Python Virtual Environment is not available."
        return 1
    fi
}

# Function to install Python Libraries from Requirements file
install_python_libraries() {
    activate_python_virtualenv
    echo "Kafka Python Library has been successfully installed."
    if [ -f "$python_requirements_file" ]; then
        echo "Contents of the Python Requirements file are:"
        cat "$python_requirements_file"
        echo "Installing the above Python Libraries..."
        pip install -r "$python_requirements_file"
        pip install git+https://github.com/dpkp/kafka-python.git
        echo "Python Libraries have been successfully installed."
        deactivate_python_virtualenv
        return 0
    else
        echo "Python Requirements file is not available."
        deactivate_python_virtualenv
        return 1
    fi
}

# Function to Uninstall Python Libraries
uninstall_python_libraries() {
    activate_python_virtualenv
    if [ -f "$python_requirements_file" ]; then
        echo "Contents of the Python Requirements file are:"
        cat "$python_requirements_file"
        echo ""
        echo "Uninstalling Python Libraries..."
        pip uninstall -r "$python_requirements_file" -y
        echo "Python Libraries have been successfully uninstalled."
        deactivate_python_virtualenv
        return 0
    else
        echo "Python Requirements file is not available."
        deactivate_python_virtualenv
        return 1
    fi
}

# Function to check if Python Adaptive Cruise Control Application is available
check_python_acc_application_available() {
    if [ -f "$python_acc_application" ]; then
        echo "Python Adaptive Cruise Control Application is available."
        return 0
    else
        echo "Python Adaptive Cruise Control Application is not available."
        return 1
    fi
}

# Function to check if Python Adaptive Cruise Control Application is running
check_python_acc_application_running() {
    if pgrep -f "python3 $python_acc_application v2" &>/dev/null; then
        echo "Python Adaptive Cruise Control Application is already running."
        return 0
    else
        echo "Python Adaptive Cruise Control Application is not running."
        return 1
    fi
}

# Function to enable Python Adaptive Cruise Control Application
enable_python_acc_application() {
    activate_python_virtualenv
    if [ -f "$python_acc_application" ]; then
        echo "Enabling Python Adaptive Cruise Control Application..."
        python3 "$python_acc_application" v2 &>/dev/null &
        if pgrep -f "python3 $python_acc_application v2" &>/dev/null; then
            echo "Python Adaptive Cruise Control Application is successfully enabled."
            deactivate_python_virtualenv
            return 0
        else
            echo "Failed to enable Python Adaptive Cruise Control Application."
            deactivate_python_virtualenv
            return 1
        fi
    else
        echo "Python Adaptive Cruise Control Application is not available."
        deactivate_python_virtualenv
        return 1
    fi
}

# Function to disable Python Adaptive Cruise Control Application
disable_python_acc_application() {
    if pgrep -f "python3 $python_acc_application v2" &>/dev/null; then
        echo "Disabling Python Adaptive Cruise Control Application..."
        pkill -f "python3 $python_acc_application v2"
        if ! pgrep -f "python3 $python_acc_application v2" &>/dev/null; then
            echo "Python Adaptive Cruise Control Application is successfully disabled."
            return 0
        else
            echo "Failed to disable Python Adaptive Cruise Control Application."
            return 1
        fi
    else
        echo "Python Adaptive Cruise Control Application is not running."
        return 1
    fi
}

# Function to display the user options
display_user_options() {
    echo "Python Setup Options:"
    echo "ipl: Install Python Libraries for Demo."
    echo "upl: Uninstall Python Libraries for Demo."
    echo "eac: Enable Python Adaptive Cruise Control Application."
    echo "dac: Disable Python Adaptive Cruise Control Application."
    echo "Enter your choice:"
}

# Function to setup Python
python_setup() {
    python_setup_user_choice=$1
    echo "Python Setup User Choice: $python_setup_user_choice"
    case $python_setup_user_choice in        
            ipl)
                check_and_install_python && check_and_create_python_virtualenv && install_python_libraries
                ;;
            upl)
                check_python_acc_application_running && disable_python_acc_application
                check_and_install_python && check_and_create_python_virtualenv && uninstall_python_libraries
                ;;
            eac)
                # check_and_install_python && check_and_create_python_virtualenv && install_python_libraries
                # Commented out the above line to speed up the application enabling process
                check_python_acc_application_available && ! check_python_acc_application_running && enable_python_acc_application
                ;;
            dac)
                check_python_acc_application_available && check_python_acc_application_running && disable_python_acc_application
                ;;
            *)
                echo "Invalid choice. Exiting the Python Setup Execution."
                ;;
        esac
}

# Main Execution Function
python_setup_main() {
    check_log_directory
    check_sudo_permission $vehicle_computer_password
    if [ "$#" -ne 1 ]; then
        display_user_options
        read -r python_setup_user_choice
        python_setup "$python_setup_user_choice"
    else    
        python_setup "$1"
    fi
}

# Execute the main function
python_setup_main "$@"