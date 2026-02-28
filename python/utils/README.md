# Utility Functions, Models, and Configuration Documentation

## Overview

This file provides a collection of utility functions, classes, and Pydantic models that are used throughout the project. These components are essential for tasks such as loading configuration files, setting up logging, interacting with Kafka, and defining data models. Additionally, there is a `config.yaml` file that stores various configurations needed across different parts of the project. The code in this file is shared across multiple modules, ensuring consistency and reducing redundancy in the project.

## Example

Refer to the `example.py` file for a simple example of sending and receiving messages using the Kafka consumer and producer classes. This example demonstrates how to use the utility functions and classes to interact with Kafka topics and process messages.
In the example file the TASMessage model is used to create a message and send it to the kafka topic, then the consumer is used to receive the message and print it to the console.

```python
import os
import sys

from python.utils.models import TASMessage
from python.utils.utils import load_config, setup_logger, TopicBoundConsumer, TopicBoundProducer, poll_and_extract_kafka_message

# Adjust the system path to include the parent directories
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))

############################################

def main(producer: TopicBoundProducer, consumer:TopicBoundConsumer) -> None:

  # how to receive messages from Kafka
  while True:
    # Poll Kafka messages
    data = poll_and_extract_kafka_message(consumer=consumer)
    if data is None:
      continue
    print(data)

  ######################

  # how to create TAS Message
  message = TASMessage(
    messageType="EXAMPLE TYPE",
    message={"EXAMPLE KEY": "EXAMPLE VALUE"}
  )

  # how to send message
  producer.send_message(message)

############################################

if __name__ == "__main__":
  # Load configuration and set up logging
  config = load_config()
  logger = setup_logger("EXAMPLE TOPIC")

  # Create instances of Kafka consumer and producer
  consumer = TopicBoundConsumer(server="remote_server", topic="EXAMPLE TOPIC", logger=logger)
  producer = TopicBoundProducer(server="remote_server", topic="EXAMPLE TOPIC", logger=logger)
```

## Functions, Classes, and Models

### 1. Path Handling Functions

- **`get_project_root`**:
  - Retrieves the root directory of the project. This is used for constructing absolute paths relative to the project root.
  
- **`get_absolute_path(relative_path: str)`**:
  - Converts a relative path to an absolute path based on the project root. This is useful for accessing files or directories within the project.

### 2. Logger Setup

- **`setup_logger(custom_name: str)`**:
  - Configures and returns a logger with both file and console handlers. The logger is essential for capturing and tracking application events, errors, and other significant information.

### 3. Configuration Management

- **`load_config()`**:
  - Loads and returns the configuration from a `config.yaml` file. The configuration is critical for setting up environment-specific parameters, such as Kafka server details, simulation settings, and environment variables.

### 4. Kafka Consumer and Producer

- **`TopicBoundConsumer(KafkaConsumer)`**:
  - A specialized Kafka consumer class that initializes a consumer bound to a specific topic and server, using the configuration loaded by `load_config`. It includes methods for polling messages from Kafka.
  
- **`TopicBoundProducer(KafkaProducer)`**:
  - A specialized Kafka producer class that initializes a producer bound to a specific topic and server. It includes methods for sending messages to Kafka topics.

### 5. Kafka Message Handling Functions

- **`poll_and_extract_kafka_message(consumer: TopicBoundConsumer)`**:
  - Polls messages from the Kafka consumer and extracts the data from the first message received. This is useful for processing real-time data streams.

- **`filter_ego_message(ego_id: str, consumer: TopicBoundConsumer)`**:
  - Filters Kafka messages to find one that matches a specific `ego_id`. This function is important for identifying and processing messages associated with a particular entity.

- **`process_trigger_message(switch: bool, consumer: TopicBoundConsumer)`**:
  - Processes Kafka messages to check for a trigger condition and toggles a switch accordingly. This can be used to manage state changes in response to real-time events.

### 6. Pydantic Models

The file also defines a series of Pydantic models that represent the structure of various data types used in the project. Pydantic models ensure that data is validated and properly structured before it is used in the application.

- **`KafkaData`**:
  - A base model that includes an optional `ego_id`. It is the parent model for other data models.

- **`CCData`**:
  - Inherits from `KafkaData` and adds fields such as `mode`, `ego_speed_new`, and `hostname`.

- **`SumoEgoData`**:
  - Extends `KafkaData` and includes fields for `ego_speed`, `ego_position`, `desired_speed`, and `leader_id`.

- **`SumoLGapData`**:
  - Inherits from `KafkaData` and includes fields for `leader_id`, `leader_speed`, and `leader_gap`.

- **`SumoLPosData`**:
  - Extends `KafkaData` and includes fields for `leader_id`, `leader_speed`, and `leader_position`.

- **`TASMessage`**:
  - A comprehensive model representing a TAS message, with both fixed and dynamic fields. It includes fields such as `sender`, `serviceType`, `responseTopic`, `requestId`, `messageType`, `message`, and an optional `subscriberTopic`.

- **`create_model_instance(data: Optional[Dict], model_class: Type[BaseModel])`**:
  - A utility function that creates an instance of a Pydantic model based on provided data, ensuring that all fields are populated with `None` if no data is provided.

## Configuration File: `config.yaml`

The `config.yaml` file is a crucial part of the project, providing configuration for various components and settings. It includes:
To create a config.yaml specific to your user, do `connect-demo/bash_scripts/create_config.sh`. This will create a config.yaml file in the python/utils directory with the correct paths and configurations.

- **Kafka Configuration**: Details about Kafka servers, topics, and related services.
  
- **CACC (Cooperative Adaptive Cruise Control) Configuration**: Parameters related to desired speed, control gains, and gap control.

- **Simulation Settings**: Includes timeout values and other simulation-specific configurations.

- **Environment Variables**: Defines paths, usernames, IP addresses, and other environment-specific settings used in the project.

This configuration file ensures that the application can be easily configured and modified without changing the codebase, promoting flexibility and scalability.

## Test Kafka/Consumer/Producer

the `test_consumer.py` `test_producer.py` `test_kafka.py` are used to test the kafka connection and the consumer and producer classes very quickly and for debugging purposes.
They are not to be used in the final production code.
To start and execute them, follow the main readme file regarding the virtual environment and the python setup, then run as normal python scripts.

## Usage

These utilities, classes, models, and configurations are intended to be imported and used across various parts of the project. They help in maintaining a clean and organized codebase by centralizing common functionality, such as configuration management, logging, Kafka interaction, and data validation.

By leveraging these utilities and the `config.yaml` file, developers can ensure that the application is consistently configured and that data is correctly structured and validated throughout the system.
