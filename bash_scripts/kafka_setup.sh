#!/bin/bash

# Path to the config.yaml file
config_yaml_file="$PWD/../python/utils/config.yaml"

# Extracting environment variables from the config.yaml file
KAFKA_PATH=$(awk -F': ' '/KAFKA_PATH:/ {print $2}' "$config_yaml_file" | tr -d '"' | tr -d "'")
BOOTSTRAP_SERVER=$(awk -F': ' '/local_server:/ {print $2}' "$config_yaml_file" | tr -d "'" | tr -d '"')
log_path=$(awk -F': ' '/log_path:/ {print $2}' "$config_yaml_file" | tr -d "'" | tr -d '"')
my_ip=$(hostname -I | awk '{print $1}')
main_computer_ip=$(awk -F': ' '/main_computer_ip:/ {print $2}' "$config_yaml_file" | tr -d "'" | tr -d '"')
vehicle_computer_1_ip=$(awk -F': ' '/vehicle_computer_1_ip:/ {print $2}' "$config_yaml_file" | tr -d "'" | tr -d '"')
vehicle_computer_2_ip=$(awk -F': ' '/vehicle_computer_2_ip:/ {print $2}' "$config_yaml_file" | tr -d "'" | tr -d '"')
vehicle_computer_username=$(awk -F': ' '/vehicle_computer_username:/ {print $2}' "$config_yaml_file" | tr -d "'" | tr -d '"')
vehicle_computer_password=$(awk -F': ' '/vehicle_computer_password:/ {print $2}' "$config_yaml_file" | tr -d "'" | tr -d '"')
readarray -t kafka_topics < <(awk -F': ' '/topic:|TOPIC:/ {print $2}' "$config_yaml_file" | tr -d '"' | tr -d "'")

# Log file path
logfiletimestamp=$(date +%Y-%m-%d_%H-%M-%S)
log_file="kafka_setup_$logfiletimestamp.log"

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

# Redirect stdout and stderr to the log file
exec > >(while read -r line; do log_message "INFO: $line"; done)
exec 2> >(while read -r line; do log_message "ERROR: $line"; done)

# Function to check for sudo permission
check_sudo_permission() {
    local password="$1"
    if ! echo "$password" | sudo -Sv &> /dev/null; then
        echo "Sorry, unable to verify sudo permission."
        echo "Either the password is incorrect or sudo is not configured properly for your user: $(whoami)."
        echo "Exiting the script execution."
        exit 1
    fi
}

# Function to check if wget is installed and install it if necessary
check_and_install_wget() {
    if command -v wget &>/dev/null; then
        echo "wget is already available."
        return 0
    fi

    echo "wget is not installed. Installing now..."

    # Define an associative array mapping package managers to their install commands
    declare -A pkg_managers=(
        ["apt-get"]="sudo apt-get update -y && sudo apt-get install wget -y"
        ["yum"]="sudo yum update -y && sudo yum install wget -y"
        ["dnf"]="sudo dnf update -y && sudo dnf install wget -y"
        ["zypper"]="sudo zypper refresh && sudo zypper install wget -y"
        ["pacman"]="sudo pacman -Syu --noconfirm wget"
        ["brew"]="brew update && brew install wget"
        ["pkg"]="sudo pkg update && sudo pkg install wget -y"
        ["port"]="sudo port selfupdate && sudo port install wget"
        ["emerge"]="sudo emerge --sync && sudo emerge -av wget"
    )

    # Loop through the package managers and execute the install command if found
    for manager in "${!pkg_managers[@]}"; do
        if command -v "$manager" &>/dev/null; then
            eval "${pkg_managers[$manager]}"
            echo "wget installed successfully using ${manager}."
            return 0
        fi
    done

    echo "Package manager not found. Please install wget manually."
    exit 1
}

# Function to check if Java is installed and install it if necessary
check_and_install_java() {
    if command -v java &>/dev/null; then
        echo "Java is already available."
        return 0
    fi

    echo "Java is not installed. Installing now..."

    # Define an associative array mapping package managers to their install commands for Java
    declare -A pkg_managers=(
        ["apt-get"]="sudo apt-get update -y && sudo apt-get install -y openjdk-8-jdk"
        ["yum"]="sudo yum update -y && sudo yum install -y java-1.8.0-openjdk"
        ["dnf"]="sudo dnf update -y && sudo dnf install -y java-1.8.0-openjdk"
        ["zypper"]="sudo zypper refresh && sudo zypper install -y java-1_8_0-openjdk"
        ["pacman"]="sudo pacman -Syu --noconfirm jdk8-openjdk"
        ["brew"]="brew update && brew install openjdk@8"
        ["pkg"]="sudo pkg update && sudo pkg install -y openjdk8"
        ["port"]="sudo port selfupdate && sudo port install openjdk8"
        ["emerge"]="sudo emerge --sync && sudo emerge -av openjdk-bin"
    )

    # Loop through the package managers and execute the install command if found
    for manager in "${!pkg_managers[@]}"; do
        if command -v "$manager" &>/dev/null; then
            eval "${pkg_managers[$manager]}"
            echo "Java installed successfully using $manager."
            return 0
        fi
    done

    echo "Package manager not found. Please install Java manually."
    exit 1
}

# Function to check if Kafka Server and Zookeeper are installed
check_kafka_server_zookeeper_installed() {
    echo "Checking if Kafka Server and Zookeeper are installed..."
    # Check if Kafka installation directory exists
    if [ -d "/opt/kafka" ]; then
        echo "Kafka is already installed in /opt/kafka"
        return 0
    fi
    if [ -d "/usr/local/kafka" ]; then
        echo "Kafka is already installed in /usr/local/kafka"
        return 0
    fi
    # Check if Kafka binaries are in the PATH
    kafka_command_path=$(command -v kafka-server-start.sh)
    if [ -x "$kafka_command_path" ]; then
        echo "Kafka binaries are found at: $kafka_command_path"
        return 0
    fi
    # Check if Kafka processes are running
    if pgrep -f "kafka-server" >/dev/null; then
        echo "Kafka Server and Zookeeper are already running."
        return 0
    fi
    # Check if Kafka service is enabled
    if systemctl is-active --quiet kafka; then
        echo "Kafka service is enabled"
        return 0
    fi
    # Check if Kafka configuration files exist
    kafka_config_path=$(find / -name server.properties 2>/dev/null)
    if [ -f "$kafka_config_path" ]; then
        echo "Kafka configuration files are found at: ${kafka_config_path}"
        return 0
    fi
    # Check Kafka version
    if [ -x "$(command -v kafka-server-start.sh)" ]; then
        kafka_version=$(kafka-server-start.sh -version 2>&1)
        echo "Kafka is already installed with version: ${kafka_version}"
        return 0
    fi
    echo "Kafka Server and zookeeper are not installed. Please install them first."
    return 1
}

# Function to download and install Kafka Server and Zookeeper
install_kafka_server_zookeeper() {
    # Download Kafka binary archive
    echo "Downloading Kafka Server and zookeeper..."
    ##if ! wget -q https://downloads.apache.org/kafka/3.7.0/kafka_2.13-3.7.0.tgz -O /tmp/kafka.tgz; then
    if ! wget -q https://dlcdn.apache.org/kafka/3.8.0/kafka_2.13-3.8.0.tgz -O /tmp/kafka.tgz; then
        echo "Failed to download Kafka."
        echo "Please check the kafka download url command at line: 87 in kafka_setup.sh script or internet connection and try again."
        echo "Exiting the Kafka Server and zookeeper installation."
        exit 1
    fi
    # Extract Kafka archive to the specified path
    echo "Extracting Kafka Server and zookeeper to ${KAFKA_PATH}..."
    sudo mkdir -p -m 775 "$KAFKA_PATH"
    sudo tar -xzf /tmp/kafka.tgz -C "$KAFKA_PATH" --strip-components=1
    # Cleanup temporary files
    rm /tmp/kafka.tgz
    # Give $USER ownership of the Kafka installation directory
    sudo chown -R "$vehicle_computer_username":"$vehicle_computer_username" "$KAFKA_PATH"
    sudo chmod -R 775 "$KAFKA_PATH"
    echo "Kafka Server and zookeeper has been installed at the path ${KAFKA_PATH}"
    return 0
}

# Function to configure server.properties file
configure_server_properties() {
    if [ -f "$KAFKA_PATH/config/server.properties" ]; then
        echo "server.properties file already exists."
        # Create a copy of the default server.properties file
        cp "$KAFKA_PATH/config/server.properties" "$KAFKA_PATH/config/server.properties.backup"
        echo "Backup of the server.properties file is created at ${KAFKA_PATH}/config/server.properties.backup"
        # Add the following line to the end of the server.properties file
        echo "listeners=PLAINTEXT://0.0.0.0:9092" >> "$KAFKA_PATH/config/server.properties"
        echo "advertised.listeners=PLAINTEXT://192.168.10.101:9092" >> "$KAFKA_PATH/config/server.properties"
        echo "Server properties are configured successfully."
        return 0
    else
        echo "server.properties file not found. Please check the Kafka installation."
        return 1
    fi
}

# Uninstall Kafka Server and Zookeeper
uninstall_kafka_server_zookeeper() {
    if check_kafka_server_zookeeper_installed; then
        if check_sudo_permission "$vehicle_computer_password" ; then
            echo "Uninstalling Apache Kafka."
            sudo rm -rf "$KAFKA_PATH"
            echo "Kafka Server and zookeeper are uninstalled."
            return 0
        else
            echo "Unable to uninstall Kafka Server and zookeeper."
            return 1
        fi
    else
        echo "Kafka Server and zookeeper are not installed."
        return 0
    fi
}

# Function to check if Kafka Server and Zookeeper are running
check_kafka_server_zookeeper_running() {
    if pgrep -f "kafka-server" >/dev/null; then
        echo "Kafka Server and Zookeeper are already running."
        return 0
    else
        echo "Kafka Server and Zookeeper are not running."
        return 1
    fi
}

# Function to Start Kafka Server and Zookeeper
start_kafka_server_zookeeper() {
    echo "Starting Kafka Server and Zookeeper..."
    $KAFKA_PATH/bin/zookeeper-server-start.sh $KAFKA_PATH/config/zookeeper.properties >> "$log_path/$log_file" &
    sleep 5
    $KAFKA_PATH/bin/kafka-server-start.sh $KAFKA_PATH/config/server.properties >> "$log_path/$log_file" &
    sleep 5
    echo "Successfully started Kafka Server and Zookeeper."
    return 0
}

#Function to Stop Kafka Server and Zookeeper
stop_kafka_server_zookeeper() {
    echo "Stopping Kafka Server and Zookeeper..."
    sudo $KAFKA_PATH/bin/kafka-server-stop.sh $KAFKA_PATH/config/server.properties >> "$log_path/$log_file" &
    sleep 5
    sudo $KAFKA_PATH/bin/zookeeper-server-stop.sh $KAFKA_PATH/config/zookeeper.properties >> "$log_path/$log_file" &
    sleep 5
    echo "Successfully stopped Kafka Server and Zookeeper."
    return 0
}

# Function to check if Kafka is set up for auto-start
check_kafka_autostart() {
    if sudo systemctl is-enabled kafka &> /dev/null; then
        echo "Kafka Server and zookeeper are already set up for auto-start."
        return 0
    else
        echo "Kafka Server and zookeeper are not set up for auto-start. "
        return 1
    fi
}

# Function to create kafka autostart commands for systemd unit file
create_kafka_autostart_commands_script() {
	if [ ! -f "$KAFKA_PATH/kafka-autostart-commands.sh" ]; then
       cat <<EOF | tee "$KAFKA_PATH/kafka-autostart-commands.sh" > /dev/null
#!/bin/bash

kafka_zookeeper_start_command() { 
	$KAFKA_PATH/bin/zookeeper-server-start.sh $KAFKA_PATH/config/zookeeper.properties 
	sleep 5
}

kafka_server_start_command() { 
	$KAFKA_PATH/bin/kafka-server-start.sh $KAFKA_PATH/config/server.properties
	sleep 5
}

kafka_zookeeper_start_command &
kafka_server_start_command &

wait
EOF
      chmod 775 "$KAFKA_PATH/kafka-autostart-commands.sh"
   fi
}

# Function to create kafka autostop commands for systemd unit file
create_kafka_autostop_commands_script() {
	if [ ! -f "$KAFKA_PATH/kafka-autostop-commands.sh" ]; then
       cat <<EOF | tee "$KAFKA_PATH/kafka-autostop-commands.sh" > /dev/null
#!/bin/bash

kafka_zookeeper_stop_command() { 
	$KAFKA_PATH/bin/zookeeper-server-stop.sh $KAFKA_PATH/config/zookeeper.properties 
	sleep 5
}

kafka_server_stop_command() { 
	$KAFKA_PATH/bin/kafka-server-stop.sh $KAFKA_PATH/config/server.properties
	sleep 5
}

kafka_zookeeper_stop_command &
kafka_server_stop_command &

wait
EOF
      chmod 775 "$KAFKA_PATH/kafka-autostop-commands.sh"
   fi
}

# Function to create systemd unit file for Kafka auto-start
create_kafka_systemd_unit_file() {
    cat <<EOF | sudo tee "/etc/systemd/system/kafka.service" > /dev/null
[Unit]
Description=Apache Kafka
After=network.target

[Service]
Type=simple
ExecStart=$KAFKA_PATH/kafka-autostart-commands.sh
ExecStop=$KAFKA_PATH/kafka-autostop-commands.sh
Restart=on-abnormal

[Install]
WantedBy=default.target
EOF
}

# Function to enable auto-start for Kafka
enable_kafka_autostart() {
    echo "Enabling auto-start for Kafka Server and zookeeper."
    create_kafka_autostart_commands_script
	create_kafka_autostop_commands_script
    create_kafka_systemd_unit_file
    sudo systemctl daemon-reload
    sudo systemctl enable kafka
    sudo systemctl start kafka
    sudo systemctl status kafka
    echo "Auto-start for Kafka Server and zookeeper are enabled."
    return 0
}

# Function to disable auto-start for Kafka
disable_kafka_autostart() {
     echo "Disabling auto-start for Kafka Server and zookeeper."
     sudo systemctl stop kafka
     sudo systemctl disable kafka
     sudo rm -f "/etc/systemd/system/kafka.service"
     rm -f "$KAFKA_PATH/kafka-autostart-commands.sh" "$KAFKA_PATH/kafka-autostop-commands.sh"
     echo "Auto-start for Kafka Server and zookeeper are disabled."
     return 0
}

# Function to extract Kafka topics from config.yaml and create them
create_kafka_topics() {
    echo "Creating Kafka topics..."
    local failed_flag=0
    for topic in "${kafka_topics[@]}"; do
        # check if topic already exists (quietly)
        if "$KAFKA_PATH/bin/kafka-topics.sh" --list --bootstrap-server "$BOOTSTRAP_SERVER" \
             | grep -qw "$topic"; then
            echo "Kafka topic $topic already exists."
            continue
        fi
        # create the topic
        echo "Creating Kafka topic $topic..."
        if ! "$KAFKA_PATH/bin/kafka-topics.sh" --create --topic "$topic" --bootstrap-server "$BOOTSTRAP_SERVER"; then
            echo "Failed to create Kafka topic $topic."
            failed_flag=$((failed_flag + 1))
        else
            echo "Kafka topic $topic is created."
        fi
        echo "-----------------------------------"
    done
    if [ "$failed_flag" -ne 0 ]; then
        echo "Failed to create some Kafka topics. Please check the logs."
        return 1
    fi
    echo "All Kafka topics are created successfully."
    return 0
}

# Function to list Kafka topics
list_kafka_topics() {
    echo "Listing Kafka topics..."
    if ! "$KAFKA_PATH/bin/kafka-topics.sh" --list --bootstrap-server "$BOOTSTRAP_SERVER"; then
        echo "Failed to list Kafka topics."
        return 1
    fi
    return 0
}

# Function to display Kafka commands
display_kafka_commands() {
    echo "-----------------------------------"
    echo "To create a Kafka topic, use the following command:"
    echo "${KAFKA_PATH}/bin/kafka-topics.sh --create --topic <topic-name> --bootstrap-server ${BOOTSTRAP_SERVER}"
    echo "Example: ${KAFKA_PATH}/bin/kafka-topics.sh --create --topic test-topic --bootstrap-server ${BOOTSTRAP_SERVER}"
    echo "-----------------------------------"
    echo "To delete a Kafka topic, use the following command:"
    echo "${KAFKA_PATH}/bin/kafka-topics.sh --delete --topic <topic-name> --bootstrap-server ${BOOTSTRAP_SERVER}"
    echo "Example: ${KAFKA_PATH}/bin/kafka-topics.sh --delete --topic test-topic --bootstrap-server ${BOOTSTRAP_SERVER}"
    echo "-----------------------------------"
    echo "To write messages to a Kafka topic, use the following command:"
    echo "${KAFKA_PATH}/bin/kafka-console-producer.sh --topic <topic-name> --bootstrap-server ${BOOTSTRAP_SERVER}"
    echo "Example: ${KAFKA_PATH}/bin/kafka-console-producer.sh --topic test-topic --bootstrap-server ${BOOTSTRAP_SERVER}"
    echo "-----------------------------------"
    echo "To read messages from a Kafka topic, use the following command:"
    echo "${KAFKA_PATH}/bin/kafka-console-consumer.sh --topic <topic-name> --from-beginning --bootstrap-server ${BOOTSTRAP_SERVER}"
    echo "Example: ${KAFKA_PATH}/bin/kafka-console-consumer.sh --topic test-topic --from-beginning --bootstrap-server ${BOOTSTRAP_SERVER}"
    echo "-----------------------------------"
    echo "To list already exisitng kafka topics, use the following command:"
    echo "${KAFKA_PATH}/bin/kafka-topics.sh --list --bootstrap-server ${BOOTSTRAP_SERVER}"
    echo "Example: ${KAFKA_PATH}/bin/kafka-topics.sh --list --bootstrap-server ${BOOTSTRAP_SERVER}"
    echo "-----------------------------------"
    echo "To check the status of Kafka Server and Zookeeper, use the following command:"
    echo "sudo systemctl status kafka"
    echo "-----------------------------------"
    echo "To check the logs of Kafka Server and Zookeeper, use the following command:"
    echo "sudo journalctl -u kafka"
}

# Main Execution Function for Main Computer
kafka_setup_main_computer() {
    check_log_directory
    echo "*****Apache Kafka Setup Script Execution Begins******"
    echo "Detected main computer with IP address: ${my_ip}"
    check_sudo_permission "$vehicle_computer_password"
    echo "1: Install Kafka binaries for Main Computer"
    echo "2: Uninstall Kafka binaries for Main Computer"
    echo "3: Start Kafka Server and zookeeper for Main Computer"
    echo "4: Stop Kafka Server and zookeeper for Main Computer"
    echo "5: Enable Kafka Server and zookeeper Auto-start for Main Computer"
    echo "6: Disable Kafka Server and zookeeper Auto-start for Main Computer"
    echo "7: Create Kafka Topics"
    echo "8: List Already Active Kafka Topics"
    echo "9: Display Kafka Commands"
    echo "Enter your choice:"
    read -r setup_user_choice
    case $setup_user_choice in
        1)
            ! check_kafka_server_zookeeper_installed && check_and_install_wget && check_and_install_java && install_kafka_server_zookeeper && configure_server_properties
            ;;
        2)  
            if check_kafka_server_zookeeper_installed; then
                if check_kafka_autostart; then
                    echo "Kafka Server and zookeeper are set up for auto-start. Disabling it now..."
                    disable_kafka_autostart
                fi
                if check_kafka_server_zookeeper_running; then
                    echo "Kafka Server and zookeeper are running. Stopping them now..."
                    stop_kafka_server_zookeeper
                fi
                # Uninstall Kafka Server and zookeeper
                echo "Uninstalling Kafka Server and zookeeper..."
                uninstall_kafka_server_zookeeper
            fi
            ;;
        3)
            check_kafka_server_zookeeper_installed && ! check_kafka_server_zookeeper_running && start_kafka_server_zookeeper
            ;;
        4)
            check_kafka_server_zookeeper_installed && check_kafka_server_zookeeper_running && stop_kafka_server_zookeeper
            ;;            
        5)
            if check_kafka_server_zookeeper_installed; then
                if check_kafka_autostart; then
                    echo "Kafka Server and zookeeper are already set up for auto-start."
                else
                    if ! check_kafka_server_zookeeper_running; then
                        echo "Kafka Server and zookeeper are not running."
                        echo "Enabling auto-start for Kafka Server and zookeeper."
                        enable_kafka_autostart
                    else
                        echo "Kafka Server and zookeeper are running. Stopping them now..."
                        stop_kafka_server_zookeeper
                        echo "Enabling auto-start for Kafka Server and zookeeper."
                        enable_kafka_autostart
                    fi
                fi
            fi
            ;;
        6)
            check_kafka_server_zookeeper_installed && check_kafka_autostart && disable_kafka_autostart
            ;;
        7)
            check_kafka_server_zookeeper_installed && check_kafka_server_zookeeper_running && create_kafka_topics
            ;;
        8)
            check_kafka_server_zookeeper_installed && check_kafka_server_zookeeper_running && list_kafka_topics
            ;;
        9)
            display_kafka_commands
            ;;
        *)
            echo "Invalid choice. Exiting the Kafka Setup Execution."
            ;;
    esac
    echo "*****Apache Kafka Setup Script Execution Ends*****"
    exit 0
}

# Main Execution Function for VC1/VC2
kafka_setup_vc() {
    check_log_directory
    echo "*****Apache Kafka Setup Script Execution Begins******"
    echo "Detected VC1/VC2 computer with IP address: ${my_ip}"
    check_sudo_permission "$vehicle_computer_password"
    echo "1: Install Kafka binaries for VC1/VC2"
    echo "2: Uninstall Kafka binaries for VC1/VC2"
    echo "3: List Already Active Kafka Topics"
    echo "4: Display Kafka Commands"
    echo "Enter your choice:"
    read -r setup_user_choice
    case $setup_user_choice in
        1)
            ! check_kafka_server_zookeeper_installed && check_and_install_wget && check_and_install_java && install_kafka_server_zookeeper
            ;;
        2)
            check_kafka_server_zookeeper_installed && uninstall_kafka_server_zookeeper
            ;;
        3)
            check_kafka_server_zookeeper_installed && list_kafka_topics
            ;;
        4)
            display_kafka_commands
            ;;
        *)
            echo "Invalid choice. Exiting the Kafka Setup Execution."
            exit 1
            ;;
    esac
    echo "*****Apache Kafka Setup Script Execution Ends*****"
    exit 0
}

# Function to check computer IP address and limit the user options based on it
kafka_setup_main() {
    # If the script is running on Main Computer then run the kafka_setup_main_computer function
    if [ "$my_ip" == "$main_computer_ip" ]; then
        kafka_setup_main_computer
    # If the script is running on VC1 or VC2 then run the kafka_setup_vc function
    elif [ "$my_ip" == "$vehicle_computer_1_ip" ] || [ "$my_ip" == "$vehicle_computer_2_ip" ]; then
        kafka_setup_vc
    else
        echo "This script is not intended to run on this computer. Exiting the script."
        exit 1
    fi
}

# Execute the main function
kafka_setup_main