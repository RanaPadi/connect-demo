import optparse
import os
import socket
import sys
import traci
from sumolib import checkBinary

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))
from python.sumo.sumo_logic import produce_out_message, update_ego_values
from python.utils.utils import load_config, get_absolute_path, setup_logger, filter_ego_message, TopicBoundProducer, TopicBoundConsumer
from python.utils.models import CCData

# Load configuration settings from a config file
config = load_config()
logger = setup_logger("sumo")

def run(
    ego_producer: TopicBoundProducer,
    leader_gap_producer: TopicBoundProducer,
    leader_pos_producer: TopicBoundProducer,
    cc_data_consumer: TopicBoundConsumer
) -> None:
    """
    Main loop for running the SUMO simulation, handling the production and consumption of messages.

    Args:
        ego_producer (TopicBoundProducer): Producer for ego vehicle messages.
        leader_gap_producer (TopicBoundProducer): Producer for leader gap messages.
        leader_pos_producer (TopicBoundProducer): Producer for leader position messages.
        cc_data_consumer (TopicBoundConsumer): Consumer for control command data.

    Returns:
        None

    Example usage:
    python3 demo.py
    """
    veh_id_list = ('v1', 'v2', 'v3')

    try:
        while traci.simulation.getMinExpectedNumber() > 0:
            # Advance the simulation by one step
            traci.simulationStep()

            for ego_id in veh_id_list:
                if ego_id not in traci.vehicle.getIDList():
                    logger.debug(f"Ego vehicle {ego_id} not found. Skipping this cycle.")
                    continue
                # Produce messages for the ego vehicle, leader gap, and leader position
                produce_out_message(
                    ego_id=ego_id,
                    ego_producer=ego_producer,
                    leader_gap_producer=leader_gap_producer,
                    leader_pos_producer=leader_pos_producer
                )

                # Get control command data for the ego vehicle and update its values
                cc_data = filter_ego_message(ego_id=ego_id, consumer=cc_data_consumer)
                cc_data = CCData(**cc_data) if cc_data is not None else None

                if cc_data is not None:
                    update_ego_values(cc_data=cc_data)

        # Close the TraCI connection after the simulation ends
        traci.close()
        sys.stdout.flush()

    except traci.exceptions.FatalTraCIError as e:
        logger.error(f"FatalTraCIError occurred: {e}")
    except KeyboardInterrupt:
        logger.warning("Simulation stopped by user. Shutting down ACC...")

if __name__ == "__main__":
    """
    Main entry point for the SUMO simulation script.    
    
    Example usage: python3 demo.py
    """
    # Ensure SUMO_HOME environment variable is set
    if 'SUMO_HOME' in os.environ:
        tools = os.path.join(os.environ['SUMO_HOME'], 'tools')
        sys.path.append(tools)
        logger.debug(f"Added {tools} to system path")
    else:
        logger.critical("SUMO_HOME environment variable not set.")
        sys.exit("Please set the 'SUMO_HOME' environment variable.")

    # Parse command-line options
    opt_parser = optparse.OptionParser()
    opt_parser.add_option(
        "--nogui",
        action="store_true",
        default=False,
        help="Run the command-line version of SUMO."
    )
    options, args = opt_parser.parse_args()
    logger.debug(f"Options: {options}")

    # Create Kafka producers and consumer
    ego_producer = TopicBoundProducer(topic=config['kafka']['ego_topic'], logger=logger)
    leader_gap_producer = TopicBoundProducer(topic=config['kafka']['leader_gap_topic'], logger=logger)
    leader_pos_producer = TopicBoundProducer(topic=config['kafka']['leader_position_topic'], logger=logger)
    cc_data_consumer = TopicBoundConsumer(topic=config['kafka']['cc_topic'], logger=logger)

    # Determine the SUMO binary to use based on command-line options
    if options.nogui:
        sumoBinary = checkBinary('sumo')
    else:
        sumoBinary = checkBinary('sumo-gui')

    hostname = socket.gethostname()

    logger.info(f"Starting SUMO simulation on hostname: {hostname}")

    # Start the SUMO simulation
    traci.start([
        sumoBinary,
        "-c", get_absolute_path('python/sumo/demo.sumocfg'),
        "--tripinfo-output", get_absolute_path('tripinfo.xml')
    ])

    # Run the main loop
    run(
        ego_producer=ego_producer,
        leader_gap_producer=leader_gap_producer,
        leader_pos_producer=leader_pos_producer,
        cc_data_consumer=cc_data_consumer
    )
