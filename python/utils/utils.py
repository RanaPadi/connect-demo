import json
import logging
import os
import sys
import yaml
from datetime import datetime
from typing import TypeVar, Any, Dict, Optional
from kafka import KafkaConsumer, KafkaProducer
from pydantic import BaseModel

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))

# Determine the project root directory (two levels up from this file)
PROJECT_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))


# helper functions for paths
def get_project_root() -> str:
    """
    Retrieves the root directory of the project.

    This function returns the path to the root directory of the project,
    which is typically used to construct absolute paths for other files
    and directories.

    Returns:
        str: The path to the root directory of the project.

    Example:
        root_path = get_project_root()
    """
    return PROJECT_ROOT


def get_absolute_path(relative_path: str) -> str:
    """
    Converts a relative path to an absolute path based on the project root.

    This function combines the project root directory with a relative path
    to produce an absolute path.

    Args:
        relative_path (str): The relative path to be combined with the project root.

    Returns:
        str: The absolute path derived from the project root and the relative path.

    Example:
        abs_path = get_absolute_path('data/file.txt')
    """
    return os.path.join(PROJECT_ROOT, relative_path)


def setup_logger(custom_name: str) -> logging.Logger:
    """
    Sets up and configures a logger with both file and console handlers.

    This function creates a logger that writes logs to a file and outputs them
    to the console. The log file is named using a custom name and the current date
    and time.

    Set log level for File handler and Console handler here.
        - File handler: >INFO
        - Console handler: >INFO / DEBUG

    Args:
        custom_name (str): The name to be used for the logger.

    Returns:
        logging.Logger: The configured logger instance.

    Example:
        logger = setup_logger("my_app")
        logger.info("Logger is set up.")
    """
    # Generate log file name with current date and time
    timestamp = datetime.now().strftime("%Y-%m-%d_%H:%M:%S")
    log_file_name = f"{custom_name}_{timestamp}.log"

    # Get the absolute path for the log file
    log_file_path = get_absolute_path('logs')
    os.makedirs(log_file_path, exist_ok=True)
    log_file = os.path.join(log_file_path, log_file_name)

    # Create or get the existing logger
    logger = logging.getLogger(custom_name)

    # Clear existing handlers if they exist
    if logger.hasHandlers():
        logger.handlers.clear()

    logger.setLevel(logging.DEBUG)  # Set the default logging level

    # Create handlers
    file_handler = logging.FileHandler(log_file)
    console_handler = logging.StreamHandler(sys.stdout)

    # Set the level for each handler
    file_handler.setLevel(logging.DEBUG) # recommended: DEBUG or INFO
    console_handler.setLevel(logging.INFO) # recommended: INFO

    # Create formatters and add them to the handlers
    formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
    file_handler.setFormatter(formatter)
    console_handler.setFormatter(formatter)

    # Add handlers to the logger
    logger.addHandler(file_handler)
    logger.addHandler(console_handler)

    return logger


logger = setup_logger("utils")

def load_config() -> dict:
    """
    Loads the configuration from a YAML file.

    This function reads a YAML configuration file from a predefined path
    and returns the configuration data as a dictionary. It logs the
    configuration data if loaded successfully or exits the program in
    case of an error.

    Returns:
        dict: The configuration data loaded from the YAML file.

    Raises:
        SystemExit: Exits the program if the configuration file is not found
                    or if there is a YAML parsing error.

    Example:
        config = load_config()
        Access config: config['kafka']['local_server']
    """

    config_path = get_absolute_path('python/utils/config.yaml')
    try:
        with open(config_path, 'r') as file:
            logger.debug(f"Configuration loaded from {config_path}")
            config_dict = yaml.safe_load(file)
            logger.debug(f"Configuration: {config_dict}")
    except yaml.YAMLError as e:
        logger.error(f"Error loading configuration: {e}")
        sys.exit(1)
    except FileNotFoundError:
        logger.error(f"Configuration file not found: {config_path}")
        sys.exit(1)
    return config_dict


config = load_config()


# Update the TopicBoundConsumer class to use the config dictionary
class TopicBoundConsumer(KafkaConsumer):
    """
    Initializes a Kafka consumer for a specific topic and server.

    This method sets up a Kafka consumer using the given topic and server
    configuration. It also initializes a logger to log the consumer's activities.

    Args:
        topic (str): The Kafka topic to subscribe to.
        server (str, optional): The server configuration key to use for connecting
            to Kafka. Defaults to "local_server".
        logger (logging.Logger, optional): A logger instance to use for logging
            messages. If not provided, a default logger for the current module
            will be used.

    Example:
        consumer = KafkaConsumer("my_topic", server="production_server")

    """
    def __init__(self, topic: str, server:str="local_server",  logger=None):
        super().__init__(
            topic,
            bootstrap_servers=[config['kafka'][server]],
            value_deserializer=lambda m: json.loads(m.decode('utf-8'))
        )
        self.logger = logger or logging.getLogger(__name__)
        self.logger.info(f"Kafka Consumer created for {topic}@{config['kafka'][server]}.")

    def poll_kafka_messages(self):
        #self.logger.debug(f"Polling Kafka messages from {self.subscription()}.")
        return super().poll(timeout_ms=config['simulation']['timeout'])

    def close_consumer(self):
        self.logger.info("Closing Kafka Consumer.")
        super().close()


# Update the TopicBoundProducer class to use the config dictionary
class TopicBoundProducer(KafkaProducer):
    """
    Initializes a Kafka producer that is bound to a specific topic and server.

    This method sets up a Kafka producer using the provided server configuration.
    It binds the producer to a specific topic, enabling message production directly
    to that topic. A logger is also initialized for logging activities related to
    the producer.

    Args:
        topic (str): The Kafka topic to which this producer will send messages.
        server (str, optional): The server configuration key for Kafka bootstrap
            servers. Defaults to "local_server".
        logger (logging.Logger, optional): A logger instance to use for logging
            messages. If not provided, a default logger for the current module will
            be used.

    Example:
        producer = TopicBoundProducer("my_topic", server="production_server")
    """
    def __init__(self, topic: str, server:str="local_server", logger=None):
        super().__init__(
            bootstrap_servers=[config['kafka'][server]],
            value_serializer=lambda m:json.dumps(m.model_dump()).encode('utf-8')
        )
        self.topic = topic
        self.logger = logger or logging.getLogger(__name__)
        self.logger.info(f"Kafka Producer created for {topic}@{config['kafka'][server]}.")
        self.server = server

    def send_message(self, message: Any):
        self.logger.info(f"Sending {json.dumps(message.model_dump()).encode('utf-8')} to topic {self.topic}@{config['kafka'][self.server]}.")
        super().send(self.topic, value=message)

    def close_producer(self):
        self.logger.info("Closing Kafka Producer.")
        super().close()


T = TypeVar("T", bound=Any)


def poll_and_extract_kafka_message(consumer: TopicBoundConsumer) -> Optional[Dict[str, Any]]:
    """
    Polls messages from the Kafka consumer and extracts the first message's data.

    This function retrieves messages from the Kafka topic via the given consumer,
    extracts the data from the first message, and logs it. It returns the data
    as a Python object.

    Args:
        consumer (TopicBoundConsumer): The Kafka consumer instance to use for polling
            messages.

    Returns:
        Any: The data extracted from the first Kafka message, or `None` if no messages
            are retrieved.

    Example:
        data = poll_and_extract_kafka_message(my_consumer)
    """
    kafka_messages = consumer.poll_kafka_messages()
    logger.debug(f"Received messages: {kafka_messages}")
    for topic_partition, messages in kafka_messages.items():
        for message in messages:
            data = message.value
            #logger.debug(f"Received data: {data}")
            if data is not None:
                logger.info(f"Extracted data: {data}")
            return data


def filter_ego_message(ego_id: str, consumer: TopicBoundConsumer) -> Dict[str, Any]:
    """
    Filters Kafka messages to find one with the specified ego_id.

    This function polls for messages from the Kafka topic and checks if any message
    contains the specified `ego_id`. If such a message is found, it logs the data
    and returns it.

    Args:
        ego_id (str): The ego_id to filter messages by.
        consumer (TopicBoundConsumer): The Kafka consumer instance to use for polling
            messages.

    Returns:
        Any: The data from the message with the matching `ego_id`, or `None` if no such
            message is found.

    Example:
        data = filter_ego_message("my_ego_id", my_consumer)
    """
    data = poll_and_extract_kafka_message(consumer=consumer)
    if data is not None:
        if data.get('ego_id') == ego_id:
            logger.info(f"Extracted data: {data}")
            return data


def process_trigger_message(switch: bool, consumer: TopicBoundConsumer) -> bool:
    """
    Processes Kafka messages to toggle a switch based on the trigger condition.

    This function polls for messages and checks if the message contains a `trigger`
    field set to `True`. If so, it toggles the provided `switch` and logs the new state.
    The function returns the updated switch state.

    Args:
        switch (bool): The current state of the switch to be toggled.
        consumer (TopicBoundConsumer): The Kafka consumer instance to use for polling
            messages.

    Returns:
        bool: The updated state of the switch after processing the Kafka messages.

    Example:
        new_switch_state = process_trigger_message(current_switch, my_consumer)
    """
    data = poll_and_extract_kafka_message(consumer=consumer)
    if data is not None:
        # Check if the trigger is present and True
        if data.get('trigger') is True:
            switch = not switch
            logger.info(f"Switch toggled to: {switch}")
            return switch
        return switch
