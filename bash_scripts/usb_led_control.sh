#!/bin/bash

# Path to the config.yaml file
config_yaml_file="$PWD/../python/utils/config.yaml"

# Extracting environment variables from the config.yaml file
blink_tool_path=$(grep 'blink_tool_path' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
log_path=$(grep 'log_path' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')

# Log file path
logfiletimestamp=$(date +"%Y-%m-%d_%H:%M:%S")
log_file="usb_led_control_$logfiletimestamp.log"

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

# Function to check if the blink1-tool is installed
check_blink_tool() {
    if [ ! -f "$blink_tool_path" ]; then
        echo "Error: blink1-tool not found. Please install the blink1-tool."
        exit 1
    fi
}

# Function to control the USB LED
usb_led_control() {
    $blink_tool_path -l 0 --"$1"
}

# Function to display the user options
display_user_options() {
    echo "Choose USB LED color:"
    echo "r/R: Red"
    echo "g/G: Green"
    echo "y/Y: Yellow"
    echo "o/O: Off"
    echo "q/Q: Quit"
    echo "Enter your choice:"
}

# Function to control the USB LED based on user input
usb_led_control_choice() {
    usb_led_control_user_choice=$1
    case $usb_led_control_user_choice in
        r|R)
            echo "Turning USB LED Red..."
            usb_led_control red
            echo "USB LED is now Red."
            ;;
        g|G)
            echo "Turning USB LED Green..."
            usb_led_control green
            echo "USB LED is now Green."
            ;;
        y|Y)
            echo "Turning USB LED Yellow..."
            usb_led_control yellow
            echo "USB LED is now Yellow."
            ;;
        o|O)
            echo "Turning USB LED Off..."
            usb_led_control off
            echo "USB LED is now Off."
            ;;
        q|Q)
            echo "Exiting..."
            exit 0
            ;;
        *)
            echo "Invalid choice. Exiting..."
            ;;
    esac
}

# Main function to control the USB LED based on user input
usb_led_control_main() {
    check_log_directory
    check_blink_tool
    if [ "$#" -ne 1 ]; then
        display_user_options
        read -r usb_led_control_user_choice
        usb_led_control_choice "$usb_led_control_user_choice"
    else
        usb_led_control_choice "$1"
    fi
}

# Call Main function
usb_led_control_main "$@"