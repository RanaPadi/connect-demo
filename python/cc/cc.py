import argparse
import os
import socket
import sys
import time

# Update the system path to include the project root directory
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))

# Import necessary functions and classes
from cc_logic import calculate_cc_control
from python.utils.models import SumoEgoData, SumoLPosData, SumoLGapData, create_model_instance
from python.utils.utils import load_config, filter_ego_message, \
    TopicBoundProducer, TopicBoundConsumer, setup_logger, process_trigger_message

# Load configuration and setup logger
config = load_config()
hostname = socket.gethostname()
logger = setup_logger("acc")


def process_messages(
        ego_id: str,
        ego_data_consumer: TopicBoundConsumer,
        leader_gap_consumer: TopicBoundConsumer,
        leader_position_consumer: TopicBoundConsumer,
        cc_data_producer: TopicBoundProducer,
        hostname: str,
        switch_consumer: TopicBoundConsumer,
        switch: bool
) -> None:
    """
    Main loop for processing messages from Kafka topics and calculating control commands.

    Args:
        ego_id (str): The ID of the ego vehicle.
        ego_data_consumer (TopicBoundConsumer): Consumer for ego vehicle data.
        leader_gap_consumer (TopicBoundConsumer): Consumer for leader gap data.
        leader_position_consumer (TopicBoundConsumer): Consumer for leader position data.
        cc_data_producer (TopicBoundProducer): Producer to send calculated control data.
        hostname (str): Hostname of the machine running the script.
        switch_consumer (TopicBoundConsumer): Consumer for switch messages that determine which leader data to use.
        switch (bool): Boolean indicating the current state of the switch.

    Returns:
        None
    """
    while True:
        # Check for trigger message and update the switch state
        switch = process_trigger_message(switch=switch, consumer=switch_consumer)
        logger.debug(
            f"Using {'leader_position_consumer' if switch else 'leader_gap_consumer'} based on switch={switch}"
        )

        # Select the appropriate consumer based on the switch state
        consumer = leader_position_consumer if switch else leader_gap_consumer

        # Process ego vehicle data
        ego_data = filter_ego_message(ego_id=ego_id, consumer=ego_data_consumer)
        if ego_data is None:
            logger.debug("Ego data is None, skipping this cycle.")
            continue

        # Process leader vehicle data
        leader_data = filter_ego_message(ego_id=ego_id, consumer=consumer)
        if leader_data is None:
            logger.debug("Leader data is None, skipping this cycle.")
            continue

        # Create a model instance based on the type of leader data
        leader_data_instance = create_model_instance(leader_data, SumoLPosData if switch else SumoLGapData)

        # Calculate control commands and send the data
        cc_data = calculate_cc_control(
            ego_data=SumoEgoData(**ego_data),
            leader_data=leader_data_instance,
            hostname=hostname
        )
        cc_data_producer.send_message(cc_data)

        # Short sleep to prevent tight loop and excessive CPU usage
        time.sleep(0.1)


def main(
        ego_id: str,
        ego_data_consumer: TopicBoundConsumer,
        leader_gap_consumer: TopicBoundConsumer,
        leader_position_consumer: TopicBoundConsumer,
        cc_data_producer: TopicBoundProducer,
        hostname: str,
        switch_consumer: TopicBoundConsumer,
        switch: bool
) -> None:
    """
    Main function to run the ACC/CACC control system.

    Args:
        ego_id (str): The ID of the ego vehicle.
        ego_data_consumer (TopicBoundConsumer): Consumer for ego vehicle data.
        leader_gap_consumer (TopicBoundConsumer): Consumer for leader gap data.
        leader_position_consumer (TopicBoundConsumer): Consumer for leader position data.
        cc_data_producer (TopicBoundProducer): Producer to send calculated control data.
        hostname (str): Hostname of the machine running the script.
        switch_consumer (TopicBoundConsumer): Consumer for switch messages that determine which leader data to use.
        switch (bool): Boolean indicating the current state of the switch.

    Returns:
        None

    Example usage: python3 cc.py v2
    """
    try:
        while True:
            process_messages(
                ego_id=ego_id,
                ego_data_consumer=ego_data_consumer,
                leader_gap_consumer=leader_gap_consumer,
                leader_position_consumer=leader_position_consumer,
                cc_data_producer=cc_data_producer,
                hostname=hostname,
                switch_consumer=switch_consumer,
                switch=switch
            )
    except KeyboardInterrupt:
        logger.warning("Simulation stopped by user. Shutting down ACC...")
    finally:
        # Ensure consumers and producers are closed properly
        ego_data_consumer.close_consumer()
        leader_gap_consumer.close_consumer()
        leader_position_consumer.close_consumer()
        cc_data_producer.close_producer()


if __name__ == "__main__":
    """
    Main entry point for the ACC/CACC controller script.
    
    Example usage: python3 cc.py v2
    """
    parser = argparse.ArgumentParser(
        description=(
            'ACC/CACC Controller for SUMO Vehicles. ACC is default, use --cacc to enable CACC mode. '
            'Specify ego_id (e.g., v2). Use --gui to enable GUI (optional).'
        )
    )
    parser.add_argument('ego_id', type=str, help='Ego Vehicle ID')
    parser.add_argument('--gui', action='store_true', help='Enable GUI (optional)')
    args = parser.parse_args()

    # Log the initial information
    info_string: str = (
        f"*** Starting ACC: ***\n"
        f"Hostname: {hostname}\n"
        f"Vehicle ID: {args.ego_id}\n"
    )
    logger.info(info_string)

    # Initialize consumers and producer
    ego_data_consumer = TopicBoundConsumer(topic=config['kafka']['ego_topic'], logger=logger)
    leader_gap_consumer = TopicBoundConsumer(topic=config['kafka']['leader_gap_topic'], logger=logger)
    leader_position_consumer = TopicBoundConsumer(topic=config['kafka']['leader_position_topic'], logger=logger)
    cc_data_producer = TopicBoundProducer(topic=config['kafka']['cc_topic'], logger=logger)
    switch_consumer = TopicBoundConsumer(topic=config['kafka']['switch_topic'], logger=logger)
    switch = False

    # Uncomment to enable GUI thread if needed
    # if args.gui:
    #     gui_thread = threading.Thread(target=run_gui, args=(args.ego_id,), daemon=True)
    #     gui_thread.start()
    #     logger.info("GUI enabled.")

    main(
        ego_id=args.ego_id,
        ego_data_consumer=ego_data_consumer,
        leader_gap_consumer=leader_gap_consumer,
        leader_position_consumer=leader_position_consumer,
        cc_data_producer=cc_data_producer,
        hostname=hostname,
        switch_consumer=switch_consumer,
        switch=switch
    )
