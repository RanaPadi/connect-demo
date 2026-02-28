import os
import sys
from typing import Any

# Update the system path to include the project root directory
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))

# Import necessary functions and classes
from python.migration_app.tas_tear_handler import create_tas_tear_request
from python.utils.utils import setup_logger, load_config, TopicBoundProducer
from python.utils.models import TASMessage

# Set up logger for debugging
logger = setup_logger("migration_app")
config = load_config()

def create_tas_unsub_request(session_id: str, subscription_id: str, producer: TopicBoundProducer) -> None:
    """
    Creates and sends a TAS_UNSUBSCRIBE_REQUEST message.

    This function constructs a TAS_UNSUBSCRIBE_REQUEST message with the given session ID
    and subscription ID and sends it using the provided Kafka producer.

    Args:
        session_id (str): The session ID associated with the unsubscribe request.
        subscription_id (str): The subscription ID associated with the unsubscribe request.
        producer (Any): The Kafka producer object used to send the message.

    Returns:
        None
    """
    logger.info("Creating TAS_UNSUBSCRIBE_REQUEST")

    # Create a TASMessage with the message type "TAS_UNSUBSCRIBE_REQUEST"
    msg = TASMessage(
        messageType="TAS_UNSUBSCRIBE_REQUEST",
        message={"sessionId": session_id, "subscriptionId": subscription_id},
    )

    # Send the message using the producer
    producer.send_message(msg)

def handle_tas_unsub_response(session_id: str, producer: TopicBoundProducer) -> None:
    """
    Handles the TAS_UNSUBSCRIBE_RESPONSE message. Initiates a TAS_TEARDOWN_REQUEST.

    This function is triggered when a TAS_UNSUBSCRIBE_RESPONSE message is received. It
    initiates a TAS_TEARDOWN_REQUEST using the provided session ID and Kafka producer.

    Args:
        session_id (str): The session ID extracted from the unsubscribe response.
        producer (Any): The Kafka producer object used to send the teardown request.

    Returns:
        None
    """
    logger.info("Received TAS_UNSUBSCRIBE_RESPONSE. Sending TEARDOWN request...")

    # Initiate a TAS_TEARDOWN_REQUEST with the provided session ID
    create_tas_tear_request(session_id=session_id, producer=producer)
