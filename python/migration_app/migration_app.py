import os
import sys
import time

# Update the system path to include the project root directory
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))

# Import necessary functions and classes from the project modules
from python.migration_app.tas_init_handler import handle_tas_init_response, create_tas_init_request
from python.migration_app.tas_notify_handler import handle_tas_notify_response
from python.migration_app.tas_sub_handler import handle_tas_sub_response
from python.migration_app.tas_tear_handler import handle_tas_tear_response
from python.migration_app.tas_unsub_handler import handle_tas_unsub_response, create_tas_unsub_request
from python.utils.utils import load_config, TopicBoundConsumer, setup_logger, TopicBoundProducer, \
    poll_and_extract_kafka_message

def main(producer: TopicBoundProducer, consumer: TopicBoundConsumer) -> None:
    """
    Main entry point for the migration application. This function manages the lifecycle of the application,
    including initialization, message processing, error handling, and graceful shutdown.

    The function does the following:
    - Initializes the TAS request by calling `create_tas_init_request`.
    - Enters a loop where it continuously polls for incoming Kafka messages.
    - Based on the type of the received message, it delegates the processing to the appropriate handler function.
    - In case of a user interruption (e.g., Ctrl+C), it initiates the unsubscribe and teardown process for cleanup.
    - Ensures that Kafka producer and consumer resources are closed properly when the program exits.

    Args:
        producer (TopicBoundProducer): The Kafka producer instance used to send messages or requests to Kafka topics.
        consumer (TopicBoundConsumer): The Kafka consumer instance used to poll and retrieve messages from Kafka topics.

    Returns:
        None: This function does not return a value. It runs indefinitely until interrupted by the user or an error.

    Side Effects:
        - Sends initial and unsubscribe requests to Kafka via the `producer`.
        - Processes Kafka messages based on their `messageType`, invoking specific handler functions for each message type.
        - Logs errors or events throughout the execution using the `logger`.
        - Initiates a clean shutdown procedure in case of interruption or error, including unsubscribing and tearing down resources.

    Exceptions:
        - KeyboardInterrupt: The function handles a user interrupt (Ctrl+C) to trigger an unsubscribe and teardown sequence.
        - AssertionError: Raised if the session ID or subscription ID is not set before trying to unsubscribe.
        - General Exception: Any error encountered during message processing or resource management is logged, and processing continues to the next message.

    Notes:
        - The `session_id` and `subscription_id` must be set when unsubscribing. If they are not available, the program raises an AssertionError.
        - The program uses multiple message types such as `TAS_INIT_RESPONSE`, `TAS_SUBSCRIBE_RESPONSE`, and `TAS_NOTIFY`. Each message type is processed by a corresponding handler function.
        - The application performs resource cleanup using the `finally` block to ensure Kafka producer and consumer are closed gracefully.

    Example Usage:
        1. The application starts by initializing the `TAS` request using `create_tas_init_request`.
        2. The message loop continuously listens for and processes Kafka messages.
        3. Upon receiving a `TAS_SUBSCRIBE_RESPONSE`, the subscription details are saved and the program continues to handle incoming messages.
        4. If interrupted, the application attempts to unsubscribe and tear down resources gracefully.
    """


    # Create the initial TAS request
    create_tas_init_request(producer=producer)
    session_id, subscription_id = None, None


    try:
        while True:
            # Poll Kafka messages
            data = poll_and_extract_kafka_message(consumer=consumer)
            if data is None:
                continue
            message_type = data.get("messageType", "")
            try:
                # Handle messages based on their type
                if message_type == "TAS_INIT_RESPONSE":
                    handle_tas_init_response(data=data, producer=producer)
                elif message_type == "TAS_SUBSCRIBE_RESPONSE":
                    session_id, subscription_id = handle_tas_sub_response(data=data)
                elif message_type == "TAS_NOTIFY":
                    handle_tas_notify_response(data=data)


            except Exception as e:
                # Log errors that occur during message processing
                logger.error(f"An error occurred during processing: {e}")
                continue  # Skip to the next iteration of the loop to continue processing

    except KeyboardInterrupt:
        # Handle user interruption by unsubscribing and tearing down
        logger.info("Interruption by user. Starting UNSUBSCRIBE and TEARDOWN...")
        try:
            assert session_id and subscription_id, "Session ID and Subscription ID must be set to unsubscribe."
            create_tas_unsub_request(session_id=session_id, subscription_id=subscription_id, producer=producer)

            while True:
                # Poll Kafka messages for unsubscribe and teardown responses
                data = poll_and_extract_kafka_message(consumer=consumer)
                if data is None:
                    continue
                message_type = data.get("messageType", "")
                if message_type == "TAS_UNSUBSCRIBE_RESPONSE":
                    handle_tas_unsub_response(session_id=session_id, producer=producer)
                elif message_type == "TAS_TEARDOWN_RESPONSE":
                    handle_tas_tear_response(data=data, producer=producer, consumer=consumer)
                    break  # Exit the inner loop after handling teardown
        except AssertionError as ae:
            logger.error(f"Assertion error: {ae}. Exiting program...")
        except Exception as e:
            logger.error(f"An error occurred during the unsubscribe and teardown process: {e}")

    finally:
        # Ensure producers and consumers are closed properly
        try:
            producer.close()
            consumer.close()
        except Exception as close_error:
            logger.error(f"An error occurred while closing resources: {close_error}")

        logger.info("User interrupted or error occurred. Exiting gracefully.")


if __name__ == "__main__":
    # Load configuration and set up logging
    config = load_config()
    logger = setup_logger("migration_app")

    # Create instances of Kafka consumer and producer
    consumer = TopicBoundConsumer(server="remote_server", topic=config['kafka']['migrationApp_topic'], logger=logger)
    producer = TopicBoundProducer(server="remote_server", topic=config['kafka']['taf_topic'], logger=logger)

    logger.info("Migration application started.")

    # Start the main function
    main(producer=producer, consumer=consumer)
