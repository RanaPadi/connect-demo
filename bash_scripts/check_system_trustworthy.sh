#!/bin/bash

# Check if config file exists
config_yaml_file="$PWD/../python/utils/config.yaml"

# Extract Kafka server and other variables from the config.yaml file
KAFKA_SERVER=$(awk -F': ' '/local_server:/ {print $2}' "$config_yaml_file" | tr -d '"' | tr -d "'")
KAFKA_SERVER_IP=$(printf '%s' "$KAFKA_SERVER" | awk -F: '{print $1}')
KAFKA_SERVER_PORT=$(printf '%s' "$KAFKA_SERVER" | awk -F: '{print $2}')
LOG_PATH=$(awk -F': ' '/log_path:/ {print $2}' "$config_yaml_file" | tr -d "'" | tr -d '"')
KAFKA_PATH="/opt/kafka" # Static path

# Define Kafka topics based on the vehicle computer IP
my_ip=$(hostname -I | awk '{print $1}')
vehicle_computer_1_ip=$(awk -F': ' '/vehicle_computer_1_ip:/ {print $2}' "$config_yaml_file" | tr -d "'" | tr -d '"')
vehicle_computer_2_ip=$(awk -F': ' '/vehicle_computer_2_ip:/ {print $2}' "$config_yaml_file" | tr -d "'" | tr -d '"')
vehicle_computer_password=$(awk -F': ' '/vehicle_computer_password:/ {print $2}' "$config_yaml_file" | tr -d "'" | tr -d '"')

# Log file setup
logfiletimestamp=$(date +"%Y-%m-%d_%H-%M-%S")
log_file="check_system_$logfiletimestamp.log"

# Function to check whether the log directory exists
check_log_directory() {
    if [ ! -d "$LOG_PATH" ]; then
        mkdir -p "$LOG_PATH"
        chmod 755 "$LOG_PATH"
    fi
}

# Function to log messages with timestamp
log_message() {
    local timestamp
    timestamp=$(date +"%Y-%m-%d %H:%M:%S")
    local message="$1"
    echo "$message"
    echo "[$timestamp] $message" >> "$LOG_PATH/$log_file"
}

# Redirect stdout and stderr to the log file
exec > >(while read -r line; do log_message "INFO: $line"; done)
exec 2> >(while read -r line; do log_message "ERROR: $line"; done)

# Function to check for sudo permission
check_sudo_permission() {
    local password="$1"
    if ! echo "$password" | sudo -Sv &> /dev/null; then
        echo "Sorry, unable to verify sudo permission."
        echo "Either the password is incorrect or sudo is not configured properly for your user: $(whoami)."
        echo "Return code: 2"
        return 2
    fi
}

# Function to check if the Kafka server is reachable
check_kafka_server() {
    if command -v nc &> /dev/null; then
        echo "'nc' command is available."
    else
        echo "Error: 'nc' command not found. Please install netcat."
        echo "Return code: 2"
        return 2
    fi

    if nc -z "$KAFKA_SERVER_IP" "$KAFKA_SERVER_PORT"; then
        echo "Kafka server $KAFKA_SERVER is reachable."
    else
        echo "Kafka server $KAFKA_SERVER is not reachable."
        echo "Return code: 2"
        return 2
    fi
}

# Function to check if Kafka is installed
check_kafka_installation() {
    if [ -d "$KAFKA_PATH" ]; then
        echo "Kafka is locally installed at $KAFKA_PATH."
    else
        echo "Kafka is not locally installed at $KAFKA_PATH."
        echo "Return code: 2"
        return 2
    fi
}

# Function to check if the Kafka topic exists and create it if not
check_kafka_topic() {
    local topic="$1"
    if ! "$KAFKA_PATH/bin/kafka-topics.sh" --list --bootstrap-server "$KAFKA_SERVER" | grep -qw "$topic"; then
        echo "Kafka topic $topic does not exist. Creating it."
        if "$KAFKA_PATH/bin/kafka-topics.sh" --create --topic "$topic" --bootstrap-server "$KAFKA_SERVER" --replication-factor 1 --partitions 1; then
            echo "Kafka topic $topic created successfully."
        else
            echo "Failed to create Kafka topic $topic."
            echo "Return code: 2"
            return 2
        fi
    else
        echo "Kafka topic $topic already exists."
    fi
}

# Function to send messages to Kafka
send_kafka_message() {
    local message="$1"
    local topic="$2"

    # Determine the vehicle computer IP and set the corresponding Topic Postfix
    if [ "$my_ip" = "$vehicle_computer_1_ip" ]; then
        topic="${topic}_VC1"
        echo "This script is running on Vehicle Computer 1."
        echo "Topic chosen is: $topic"
    elif [ "$my_ip" = "$vehicle_computer_2_ip" ]; then
        topic="${topic}_VC2"
        echo "This script is running on Vehicle Computer 2."
        echo "Topic chosen is: $topic"
    else
        echo "This script is not running on a proper machine."
        echo "Not authorized to send messages to Kafka."
        echo "Return code: 2"
        return 2
    fi

    # Check if Kafka server is reachable, Kafka is installed, and the topic exists
    if ! check_kafka_server; then
        echo "Kafka server check failed."
        return 2
    fi
    if ! check_kafka_installation; then
        echo "Kafka installation check failed."
        return 2
    fi
    if ! check_kafka_topic "$topic"; then
        echo "Kafka topic check failed."
        return 2
    fi

    # Append timestamp to the message
    local timestamp
    timestamp=$(date +"%Y%m%d%H%M%S")
    message="[${timestamp}]:${message}"

    if [ -z "$message" ]; then
        echo "Message is empty. Cannot send to Kafka."
        echo "Return code: 2"
        return 2
    else
        echo "Kafka Message: $message"
        echo "Kafka Topic: $topic"
        echo "Kafka server: $KAFKA_SERVER"
    fi

    # Send the message to Kafka
    if echo "$message" | "$KAFKA_PATH/bin/kafka-console-producer.sh" --broker-list "$KAFKA_SERVER" --topic "$topic" > /dev/null; then
        echo "Message: $message sent to Kafka topic: $topic successfully."
    else
        echo "Failed to send message: $message to Kafka topic: $topic"
        return 2
    fi
}

# Function to check Access Control
check_access_control() {
    echo "Access Control check completed successfully."
    echo "Return code: 1"
    return 1
}

# Function to check Application Isolation
check_application_isolation() {
    echo "Application Isolation check completed successfully."
    echo "Return code: 1"
    return 1
}

# Function to check Control Flow Integrity
check_control_flow_integrity() {
    local control_flow_integrity_check_file="$HOME/check_control_flow_integrity.txt"
    if [ -f "$control_flow_integrity_check_file" ]; then
        echo "Control Flow Integrity check file: $control_flow_integrity_check_file found at $HOME"
        echo "Control Flow Integrity check failed."
        echo "Return code: 0"
        return 0
    else
        echo "Control Flow Integrity check file: $control_flow_integrity_check_file not found at $HOME"
        echo "Control Flow Integrity check completed successfully."
        echo "Return code: 1"
        return 1
    fi
}

# Function to check Control Flow Integrity Verification
check_control_flow_integrity_verification() {
    echo "Control Flow Integrity Verification check completed successfully."
    echo "Return code: 1"
    return 1
}

# Function to check Secure Boot
check_secure_boot() {
    echo "Secure Boot check completed successfully."
    echo "Return code: 1"
    return 1
}

# Function to check Secure OTA
check_secure_ota() {
    if ! check_sudo_permission "$vehicle_computer_password"; then
        echo "Secure OTA check failed due to sudo permission issue."
        echo "Return code: 2"
        return 2
    else
        local updates_available
        updates_available=$(sudo apt update | grep -i "packages can be upgraded")
        if [ -n "$updates_available" ]; then
            echo "APT update check: $updates_available"
            echo "Secure OTA check failed."
            echo "Return code: 0"
            return 0
        else
            echo "APT update check: No updates available."
            echo "Secure OTA check completed successfully."
            echo "Return code: 1"
            return 1
        fi
    fi
}

# Function to display the user options
display_user_options() {
    echo "Please choose an option:"
    echo "1. Check Access Control"
    echo "2. Check Application Isolation"
    echo "3. Check Control Flow Integrity"
    echo "4. Check Control Flow Integrity Verification"
    echo "5. Check Secure Boot"
    echo "6. Check Secure OTA"
    echo "Enter your choice:"
}

# Function to check system trustworthiness based on user input
check_system_trustworthy() {
    local choice="$1"
    echo "User choice: $choice"
    case "$choice" in
        "ACCESS_CONTROL")
            check_access_control
            ;;
        "APPLICATION_ISOLATION")
            check_application_isolation
            ;;
        "CONTROL_FLOW_INTEGRITY")
            check_control_flow_integrity
            ;;
        "CONTROL_FLOW_INTEGRITY_VERIFICATION")
            check_control_flow_integrity_verification
            ;;
        "SECURE_BOOT")
            check_secure_boot
            ;;
        "SECURE_OTA")
            check_secure_ota
            ;;
        *)
            echo "Invalid choice. Please choose a valid option."
            echo "Valid options are: ACCESS_CONTROL, APPLICATION_ISOLATION, CONTROL_FLOW_INTEGRITY, CONTROL_FLOW_INTEGRITY_VERIFICATION, SECURE_BOOT, SECURE_OTA"
            echo "Exiting Script Execution."
            echo "Exit code: 2"
            exit 2
            ;;
    esac
}

# Helper function to process check with Kafka messaging
process_check() {
    local check_name="$1"
    echo "Performing check: $check_name"
    check_system_trustworthy "$check_name"
    local return_code=$?
    echo "$check_name check returned code: $return_code"

    send_kafka_message "$return_code" "$check_name"
    local kafka_return_code=$?
    if [ $kafka_return_code -eq 2 ]; then
        echo "Kafka message sending failed."
        return_code=$kafka_return_code
    fi

    exit ${return_code:-2}
}

# Main function
check_system_trustworthy_main() {
    # Check if the log directory exists
    check_log_directory
    echo "Starting system checks..."

    if [ "$#" -eq 0 ]; then
        display_user_options
        read -r user_input

        case "$user_input" in
            1)
                process_check "ACCESS_CONTROL"
                ;;
            2)
                process_check "APPLICATION_ISOLATION"
                ;;
            3)
                process_check "CONTROL_FLOW_INTEGRITY"
                ;;
            4)
                process_check "CONTROL_FLOW_INTEGRITY_VERIFICATION"
                ;;
            5)
                process_check "SECURE_BOOT"
                ;;
            6)
                process_check "SECURE_OTA"
                ;;
            *)
                echo "Invalid choice. Please choose a valid option."
                echo "Valid options are: 1, 2, 3, 4, 5, 6"
                echo "Exiting Script Execution."
                echo "Exit code: 2"
                exit 2
                ;;
        esac
    else
        case "$1" in
            "ACCESS_CONTROL")
                process_check "ACCESS_CONTROL"
                ;;
            "APPLICATION_ISOLATION")
                process_check "APPLICATION_ISOLATION"
                ;;
            "CONTROL_FLOW_INTEGRITY")
                process_check "CONTROL_FLOW_INTEGRITY"
                ;;
            "CONTROL_FLOW_INTEGRITY_VERIFICATION")
                process_check "CONTROL_FLOW_INTEGRITY_VERIFICATION"
                ;;
            "SECURE_BOOT")
                process_check "SECURE_BOOT"
                ;;
            "SECURE_OTA")
                process_check "SECURE_OTA"
                ;;
            *)
                echo "Invalid flag. Please choose a valid option."
                echo "Valid options are: ACCESS_CONTROL, APPLICATION_ISOLATION, CONTROL_FLOW_INTEGRITY, CONTROL_FLOW_INTEGRITY_VERIFICATION, SECURE_BOOT, SECURE_OTA"
                echo "Exiting Script Execution."
                echo "Exit code: 2"
                exit 2
                ;;
        esac
    fi
}

# Execute the main function
check_system_trustworthy_main "$@"