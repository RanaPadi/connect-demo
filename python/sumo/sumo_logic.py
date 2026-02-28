import os
import sys
import traci
from typing import Tuple, Optional

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))
from python.utils.utils import TopicBoundProducer, setup_logger
from python.utils.models import SumoEgoData, CCData, SumoLGapData, SumoLPosData

# Set up logger for debugging
logger = setup_logger("sumo")


def produce_out_message(
        ego_id: str,
        ego_producer: TopicBoundProducer,
        leader_gap_producer: TopicBoundProducer,
        leader_pos_producer: TopicBoundProducer
) -> None:
    """
    Produce and send out messages for ego vehicle and leader data.

    Args:
        ego_id (str): The ID of the ego vehicle.
        ego_producer (TopicBoundProducer): Producer for sending ego vehicle data.
        leader_gap_producer (TopicBoundProducer): Producer for sending leader gap data.
        leader_pos_producer (TopicBoundProducer): Producer for sending leader position data.

    Returns:
        None
    """
    try:
        # Retrieve ego vehicle data
        ego_speed: float = traci.vehicle.getSpeed(ego_id)
        ego_position: Tuple[float, float] = traci.vehicle.getPosition(ego_id)
        desired_speed: float = traci.vehicle.getMaxSpeed(ego_id)

        # Retrieve leader data
        leader = traci.vehicle.getLeader(ego_id, 10000)
        if leader:
            leader_id: str = leader[0]
            leader_gap: float = leader[1]
            leader_speed: float = traci.vehicle.getSpeed(leader_id)
            leader_position: Tuple[float, float] = traci.vehicle.getPosition(leader_id)
        else:
            leader_id, leader_gap, leader_speed, leader_position = (None, None, None, None)

        # Create and send messages
        sumo_ego_data = SumoEgoData(
            ego_id=ego_id,
            ego_speed=ego_speed,
            ego_position=ego_position,
            desired_speed=desired_speed,
            leader_id=leader_id
        )
        ego_producer.send_message(sumo_ego_data)

        sumo_lgap_data = SumoLGapData(
            ego_id=ego_id,
            leader_id=leader_id,
            leader_speed=leader_speed,
            leader_gap=leader_gap
        )
        leader_gap_producer.send_message(sumo_lgap_data)

        sumo_lpos_data = SumoLPosData(
            ego_id=ego_id,
            leader_id=leader_id,
            leader_speed=leader_speed,
            leader_position=leader_position
        )
        leader_pos_producer.send_message(sumo_lpos_data)

        logger.debug(f"Produced messages for ego_id={ego_id} with leader_id={leader_id}")

    except Exception as e:
        logger.error(f"Error producing messages for ego_id={ego_id}: {e}")


def update_ego_values(cc_data: CCData) -> None:
    """
    Update the speed of the ego vehicle based on the control command data.

    Args:
        cc_data (CCData): The control command data.

    Returns:
        None
    """
    try:
        # Update ego vehicle speed
        traci.vehicle.setSpeed(cc_data.ego_id, cc_data.ego_speed_new)
        # Uncomment the following line if desired speed update is needed
        # traci.vehicle.setMaxSpeed(cc_data.ego_id, cc_data.desired_speed)

        logger.debug(f"Updated speed for ego_id={cc_data.ego_id} to {cc_data.ego_speed_new}")

    except Exception as e:
        logger.error(f"Error updating values for ego_id={cc_data.ego_id}: {e}")
