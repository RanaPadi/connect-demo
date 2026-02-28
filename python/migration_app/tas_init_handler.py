import os
import sys
from typing import Any

# Update the system path to include the project root directory
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))

from python.migration_app.tas_sub_handler import create_tas_sub_request
from python.utils.utils import setup_logger, load_config, TopicBoundProducer
from python.utils.models import TASMessage


# Set up logger for debugging
logger = setup_logger("migration_app")
config = load_config()

def create_tas_init_request(producer: TopicBoundProducer) -> None:
    """
    Creates and sends a TAS_INIT_REQUEST message.

    This function initializes a TAS_INIT_REQUEST message with a predefined trust model template
    and sends it using the provided Kafka producer.

    Args:
        producer (Any): The Kafka producer object used to send the message.

    Returns:
        None
    """
    logger.info("Creating TAS_INIT_REQUEST")

    # trust models: VCM or BRUSSELS
    msg = TASMessage(
        messageType="TAS_INIT_REQUEST",
        message={"trustModelTemplate": "VCM@0.0.1"}
    )

    producer.send_message(msg)


def handle_tas_init_response(data: Any, producer: Any) -> None:
    """
    Handles the TAS_INIT_RESPONSE message and initiates a TAS_SUBSCRIBE_REQUEST.

    This function extracts the session ID from the TAS_INIT_RESPONSE message data
    and creates a TAS_SUBSCRIBE_REQUEST using this session ID.

    Args:
        data (Any): The data received from the TAS_INIT_RESPONSE message.
        producer (Any): The Kafka producer object used to send the subscription request.

    Returns:
        None
    """
    logger.info("Received TAS_INIT_RESPONSE message")

    session_id = data.get("message", {}).get("sessionId", "")

    logger.info(f"Extracted Session ID: {session_id}")

    if session_id == "":
        logger.error("Session ID not found in TAS_INIT_RESPONSE message")
        exit(-1)

    create_tas_sub_request(session_id=session_id, producer=producer)
