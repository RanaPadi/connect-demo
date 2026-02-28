import os
import sys
from typing import Any, Tuple

# Update the system path to include the project root directory
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))

# Import necessary functions and classes
from python.utils.models import TASMessage
from python.utils.utils import setup_logger, load_config, TopicBoundProducer

# Set up logger for debugging
logger = setup_logger("migration_app")
config = load_config()

def create_tas_sub_request(session_id: str, producer: TopicBoundProducer) -> None:
    """
    Creates and sends a TAS_SUBSCRIBE_REQUEST message.

    Constructs a subscription request message using the provided session ID and sends it
    using the given Kafka producer. The message includes a filter and a trigger.

    Args:
        session_id (str): The session ID for the subscription request.
        producer (Any): The Kafka producer object used to send the message.

    Returns:
        None
    """
    logger.info("Creating TAS_SUBSCRIBE_REQUEST")

    # Create a TASMessage with the message type "TAS_SUBSCRIBE_REQUEST"
    # leave filter blank for now
    msg = TASMessage(
        messageType="TAS_SUBSCRIBE_REQUEST",
        message={
            "sessionId": session_id,
            "subscribe": {
                "filter": [
                ]
            },
            "trigger": "ACTUAL_TRUSTWORTHINESS_LEVEL"
        },
    )

    # Send the message using the producer
    producer.send_message(msg)

def handle_tas_sub_response(data: Any) -> Tuple[str, str]:
    """
    Handles the TAS_SUBSCRIBE_RESPONSE message.

    Extracts the session ID and subscription ID from the response message and returns them.

    Args:
        data (Any): The data received from the TAS_SUBSCRIBE_RESPONSE message.

    Returns:
        Tuple[str, str]: A tuple containing the session ID and subscription ID.
    """
    logger.info("Received TAS_SUBSCRIBE_RESPONSE")

    # Extract subscription ID from the response data. Session ID extracted again for synchronous programming purposes.
    session_id = data.get("message", {}).get("sessionId", "")
    subscription_id = data.get("message", {}).get("subscriptionId", "")

    logger.info(f"Session ID: {session_id}, Subscription ID: {subscription_id}")
    if not session_id or not subscription_id:
        logger.error("Session ID or Subscription ID not found in the response data")
        logger.error(data.get("message", {}).get("error", ""))
        exit(-1)


    return session_id, subscription_id
