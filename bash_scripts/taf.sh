#!/bin/bash

# Path to the config.yaml file
config_yaml_file="$PWD/../python/utils/config.yaml"

# Extracting environment variables from the config.yaml file
migration_script=$(grep 'migration_script' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
TAF_PATH=$(grep 'TAF_PATH' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
RTL_ATL_Main=$(grep 'RTL_ATL_Main' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')
log_path=$(grep 'log_path' $config_yaml_file | awk -F': ' '{print $2}' | sed 's/^"//' | sed 's/"$//')

# Log file path
logfiletimestamp=$(date +"%Y-%m-%d_%H:%M:%S")
log_file="taf_$logfiletimestamp.log"

# Function to check whether the log directory exists and create it if it does not
check_log_directory() {
    if [ ! -d "$log_path" ]; then
        mkdir -p "$log_path"
        chmod 755 "$log_path"
    fi
}

# Function to check if TAF directory exists and create it if it does not
check_taf_directory() {
    if [ ! -d "$TAF_PATH" ]; then
        mkdir -p "$TAF_PATH"
        chmod 755 "$TAF_PATH"
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

# Function to check if the file exists
check_file_exists() {
    file_path="$1"
    if [[ ! -f "$file_path" ]]; then
        echo "The file $file_path does not exist."
        echo "Creating the file $file_path."
        touch "$file_path"
        chmod 755 "$file_path"
        if [[ ! -f "$file_path" ]]; then
            echo "The file $file_path could not be created."
            echo "Exiting Script Execution."
            return 1
        fi
    fi
}

# Main Execution Function
taf_main() {
    check_log_directory
    check_taf_directory
    echo "*****TAF Modification Script Execution Begins*****"
    # Choose the machine for which the RTL and ATL value will be modified
    echo "Choose the machine for which you want to modify the RTL value."
    echo "1: Vehicle Computer 1"
    echo "2: Vehicle Computer 2"
    echo "Enter your choice:"
    read -r machine_choice
    # Check the machine user's choice
    if [[ "$machine_choice" == "1" ]]; then
        machine_name="vehicle_computer_1"
    elif [[ "$machine_choice" == "2" ]]; then
        machine_name="vehicle_computer_2"
    else
        echo "Invalid choice. Please enter 1 or 2."
        echo "Exiting Script Execution."
        exit 1
    fi
    
    # Choose the RTL value for the chosen machine
    echo "Use default RTL values for the $machine_name? (Yes/No)"
    read -r rtl_choice
    rtl_choice=$(echo "$rtl_choice" | tr '[:upper:]' '[:lower:]')
    # Check the RTL user's choice
    if [[ "$rtl_choice" == "no" || "$rtl_choice" == "n" ]]; then
        echo "Enter the RTL values for $machine_name in range 0-1 with 2 decimal places. Example: 0.25"
        echo "And make sure that Belief + Disbelief + Uncertainty = 1."
        echo "Enter RTL Belief Threshold Value (0-1):"
        read -r rtl_belief
        # Check if the input is a valid number
        if [[ ! "$rtl_belief" =~ ^[0-9]+(\.[0-9]{1,2})?$ ]]; then
            echo "Invalid input. Please enter a valid number in range 0-1 with 2 decimal places."
            echo "Exiting Script Execution."
            exit 1
        fi
        # Check if the input is in the range 0-1
        if (( $(echo "$rtl_belief > 1" | bc -l) )) || (( $(echo "$rtl_belief < 0" | bc -l) )); then
            echo "Invalid input. Please enter a valid number in range 0-1 with 2 decimal places."
            echo "Exiting Script Execution."
            exit 1
        fi
        echo "Enter RTL Disbelief Threshold Value (0-1):"
        read -r rtl_disbelief
        # Check if the input is a valid number
        if [[ ! "$rtl_disbelief" =~ ^[0-9]+(\.[0-9]{1,2})?$ ]]; then
            echo "Invalid input. Please enter a valid number in range 0-1 with 2 decimal places."
            echo "Exiting Script Execution."
            exit 1
        fi
        # Check if the input is in the range 0-1
        if (( $(echo "$rtl_disbelief > 1" | bc -l) )) || (( $(echo "$rtl_disbelief < 0" | bc -l) )); then
            echo "Invalid input. Please enter a valid number in range 0-1 with 2 decimal places."
            echo "Exiting Script Execution."
            exit 1
        fi
        echo "Enter RTL Uncertainty Threshold Value (0-1):"
        read -r rtl_uncertainty
        # Check if the input is a valid number
        if [[ ! "$rtl_uncertainty" =~ ^[0-9]+(\.[0-9]{1,2})?$ ]]; then
            echo "Invalid input. Please enter a valid number in range 0-1 with 2 decimal places."
            echo "Exiting Script Execution."
            exit 1
        fi
        # Check if the input is in the range 0-1
        if (( $(echo "$rtl_uncertainty > 1" | bc -l) )) || (( $(echo "$rtl_uncertainty < 0" | bc -l) )); then
            echo "Invalid input. Please enter a valid number in range 0-1 with 2 decimal places."
            echo "Exiting Script Execution."
            exit 1
        fi
        # Check if the sum of Belief, Disbelief and Uncertainty Thresholds is equal to 1
        rtl=$(echo "$rtl_belief + $rtl_disbelief + $rtl_uncertainty" | bc -l)
        if (( $(echo "$rtl != 1" | bc -l) )); then
            echo "The sum of Belief, Disbelief and Uncertainty Thresholds is not equal to 1."
            echo "Please enter the values again."
            echo "Exiting Script Execution."
            exit 1
        fi
        echo "Using RTL values: Belief=$rtl_belief, Disbelief=$rtl_disbelief, Uncertainty=$rtl_uncertainty"
    elif [[ "$rtl_choice" == "yes" || "$rtl_choice" == "y" ]]; then
        rtl_belief=0.60
        rtl_disbelief=0.20
        rtl_uncertainty=0.20
        rtl=$(echo "$rtl_belief + $rtl_disbelief + $rtl_uncertainty" | bc -l)
        echo "Using RTL values: Belief=$rtl_belief, Disbelief=$rtl_disbelief, Uncertainty=$rtl_uncertainty"
    else
        echo "Invalid choice. Please enter Yes or No."
        echo "Exiting Script Execution."
        exit 1
    fi

    # Choose the ATL value for the chosen machine
    echo "Use default ATL values for the $machine_name? (Yes/No)"
    read -r atl_choice
    atl_choice=$(echo "$atl_choice" | tr '[:upper:]' '[:lower:]')
    # Check the ATL user's choice
    if [[ "$atl_choice" == "no" || "$atl_choice" == "n" ]]; then
        echo "Enter the ATL values for $machine_name in range 0-1 with 2 decimal places. Example: 0.25"
        echo "And make sure that Belief + Disbelief + Uncertainty = 1."
        echo "Enter ATL Belief Value (0-1):"
        read -r atl_belief
        # Check if the input is a valid number
        if [[ ! "$atl_belief" =~ ^[0-9]+(\.[0-9]{1,2})?$ ]]; then
            echo "Invalid input. Please enter a valid number in range 0-1 with 2 decimal places."
            echo "Exiting Script Execution."
            exit 1
        fi
        # Check if the input is in the range 0-1
        if (( $(echo "$atl_belief > 1" | bc -l) )) || (( $(echo "$atl_belief < 0" | bc -l) )); then
            echo "Invalid input. Please enter a valid number in range 0-1 with 2 decimal places."
            echo "Exiting Script Execution."
            exit 1
        fi
        echo "Enter ATL Disbelief Value (0-1):"
        read -r atl_disbelief
        # Check if the input is a valid number
        if [[ ! "$atl_disbelief" =~ ^[0-9]+(\.[0-9]{1,2})?$ ]]; then
            echo "Invalid input. Please enter a valid number in range 0-1 with 2 decimal places."
            echo "Exiting Script Execution."
            exit 1
        fi
        # Check if the input is in the range 0-1
        if (( $(echo "$atl_disbelief > 1" | bc -l) )) || (( $(echo "$atl_disbelief < 0" | bc -l) )); then
            echo "Invalid input. Please enter a valid number in range 0-1 with 2 decimal places."
            echo "Exiting Script Execution."
            exit 1
        fi
        echo "Enter ATL Uncertainty Value (0-1):"
        read -r atl_uncertainty
        # Check if the input is a valid number
        if [[ ! "$atl_uncertainty" =~ ^[0-9]+(\.[0-9]{1,2})?$ ]]; then
            echo "Invalid input. Please enter a valid number in range 0-1 with 2 decimal places."
            echo "Exiting Script Execution."
            exit 1
        fi
        # Check if the input is in the range 0-1
        if (( $(echo "$atl_uncertainty > 1" | bc -l) )) || (( $(echo "$atl_uncertainty < 0" | bc -l) )); then
            echo "Invalid input. Please enter a valid number in range 0-1 with 2 decimal places."
            echo "Exiting Script Execution."
            exit 1
        fi
        # Check if the sum of Belief, Disbelief and Uncertainty Thresholds is equal to 1
        atl=$(echo "$atl_belief + $atl_disbelief + $atl_uncertainty" | bc -l)
        if (( $(echo "$atl != 1" | bc -l) )); then
            echo "The sum of Belief, Disbelief and Uncertainty Thresholds is not equal to 1."
            echo "Please enter the values again."
            echo "Exiting Script Execution."
            exit 1
        fi
        echo "Using ATL values: Belief=$atl_belief, Disbelief=$atl_disbelief, Uncertainty=$atl_uncertainty"
    elif [[ "$atl_choice" == "yes" || "$atl_choice" == "y" ]]; then
        atl_belief=0.20
        atl_disbelief=0.40
        atl_uncertainty=0.40
        atl=$(echo "$atl_belief + $atl_disbelief + $atl_uncertainty" | bc -l)
        echo "Using ATL values: Belief=$atl_belief, Disbelief=$atl_disbelief, Uncertainty=$atl_uncertainty"
    else
        echo "Invalid choice. Please enter Yes or No."
        echo "Exiting Script Execution."
        exit 1
    fi

    # Add the RTL and ATL values to the RTL-ATL file
    # Check if the RTL-ATL file exists
    check_file_exists "$RTL_ATL_Main"
    systemtime=$(date +"%Y%m%d%H%M%S")
    echo "$systemtime,$machine_name,$rtl_belief,$rtl_disbelief,$rtl_uncertainty,$atl_belief,$atl_disbelief,$atl_uncertainty" >> "$RTL_ATL_Main"
    echo "Added RTL and ATL values to the RTL-ATL Main file."
    if [[ "$machine_name" == "vehicle_computer_1" ]]; then
        if (( $(echo "$rtl_belief > $atl_belief" | bc -l) )) || (( $(echo "$rtl_disbelief < $atl_disbelief" | bc -l) )) || (( $(echo "$rtl_uncertainty < $atl_uncertainty" | bc -l) )); then
            source "$migration_script" "VC2"
            echo "Invoked Migration of applications from Vehicle Computer 1 to Vehicle Computer 2."
        fi
    elif [[ "$machine_name" == "vehicle_computer_2" ]]; then
        if (( $(echo "$rtl_belief > $atl_belief" | bc -l) )) || (( $(echo "$rtl_disbelief < $atl_disbelief" | bc -l) )) || (( $(echo "$rtl_uncertainty < $atl_uncertainty" | bc -l) )); then
            source "$migration_script" "VC1"
            echo "Invoked Migration of applications from Vehicle Computer 2 to Vehicle Computer 1."
        fi
    else
        echo "Invalid machine name. Exiting Script Execution."
        exit 1
    fi

    echo "*****TAF Modification Script Execution Ends*****"
}

# Execute the main function
taf_main
