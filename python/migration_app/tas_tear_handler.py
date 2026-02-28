import os
import sys
from typing import Any

# Update the system path to include the project root directory
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))

# Import necessary functions and classes
from python.utils.utils import setup_logger, load_config, TopicBoundConsumer
from python.utils.models import TASMessage

# Set up logger for debugging
logger = setup_logger("migration_app")
config = load_config()

def create_tas_tear_request(session_id: str, producer: Any) -> None:
    """
    Creates and sends a TAS_TEARDOWN_REQUEST message.

    This function constructs a TAS_TEARDOWN_REQUEST message with the given session ID
    and sends it using the provided Kafka producer.

    Args:
        session_id (str): The session ID associated with the teardown request.
        producer (Any): The Kafka producer object used to send the message.

    Returns:
        None
    """
    logger.info("Creating TAS_TEARDOWN_REQUEST")

    # Create a TASMessage with the message type "TAS_TEARDOWN_REQUEST"
    msg = TASMessage(
        messageType="TAS_TEARDOWN_REQUEST",
        message={"sessionId": session_id},
    )

    # Send the message using the producer
    producer.send_message(msg)

def handle_tas_tear_response(data: Any, producer: Any, consumer: TopicBoundConsumer) -> None:
    """
    Handles the TAS_TEARDOWN_RESPONSE message. Closes the producer and consumer and exits the program.

    This function is triggered when a TAS_TEARDOWN_RESPONSE message is received. It closes
    the Kafka producer and consumer, then exits the program gracefully.

    Args:
        data (Any): The data received from the TAS_TEARDOWN_RESPONSE message.
        producer (Any): The Kafka producer object to be closed.
        consumer (Any): The Kafka consumer object to be closed.

    Returns:
        None
    """
    logger.info("Received TAS_TEARDOWN_RESPONSE. Closing producer and consumer and exiting program...")

    try:
        # Close the producer and consumer
        producer.close()
        consumer.close()
        logger.info("Producer and consumer closed successfully.")
    except Exception as e:
        # Log any error that occurs while closing resources
        logger.error(f"An error occurred while closing resources: {e}")
    finally:
        # Exit the program with a status code of 0 (indicating success)
        logger.info("Exiting the program now.")
        sys.exit(0)
