#!/bin/bash

# Path to the template and final config.yaml file
template_config_yaml_file="$PWD/python/utils/template_config.yaml"
config_yaml_file="$PWD/python/utils/config.yaml"

# Ask the user for their password
echo "Please enter your password:"
read -s ACTUAL_USER_PASSWORD

# Get the user's home directory
ACTUAL_USER_HOME_DIRECTORY="$HOME"
ACTUAL_USER_NAME="$USER"

# Replace placeholders in the template file and save to a new file
sed -e "s|\$PWD|$PWD|g" \
    -e "s|\$USER_PASSWORD|$ACTUAL_USER_PASSWORD|g" \
    -e "s|\$USER_NAME|$ACTUAL_USER_NAME|g" \
    "$template_config_yaml_file" > "$config_yaml_file"

echo "Configuration file 'config.yaml' has been created at $config_yaml_file"