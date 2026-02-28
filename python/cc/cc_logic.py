import math
import os
import sys
from typing import Optional, Any, Tuple

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))
from python.utils.utils import load_config
from python.utils.models import CCData, SumoEgoData, SumoLGapData, SumoLPosData

# Load configuration settings from a config file
config = load_config()

def calculate_gap(ego_data: SumoEgoData, leader_pos_data: SumoLPosData) -> float:
    """
    Calculate the Euclidean distance (gap) between the ego vehicle and the leader vehicle based on their positions.

    Args:
        ego_data (SumoEgoData): Data of the ego vehicle including its position.
        leader_pos_data (SumoLPosData): Data of the leader vehicle including its position.

    Returns:
        float: The Euclidean distance (gap) between the ego vehicle and the leader vehicle.
    """
    ego_position: Tuple[float, float] = ego_data.ego_position
    leader_position: Tuple[float, float] = leader_pos_data.leader_position

    # Calculate the Euclidean distance between the ego and leader positions
    leader_gap = math.sqrt((leader_position[0] - ego_position[0]) ** 2 +
                           (leader_position[1] - ego_position[1]) ** 2)

    return leader_gap

def calculate_cc_control(
    ego_data: SumoEgoData,
    leader_data: Optional[Any],
    hostname: str
) -> CCData:
    """
    Calculate the Control Command (CC) for the ego vehicle based on its data and the leader vehicle's data.

    Args:
        ego_data (SumoEgoData): Data of the ego vehicle including its speed and desired speed.
        leader_data (Optional[Any]): Data of the leader vehicle which can be either SumoLGapData or SumoLPosData.
        hostname (str): Hostname of the machine running the script.

    Returns:
        CCData: The control command data including the new speed for the ego vehicle and the control mode.
    """
    # Retrieve configuration parameters
    speed_control_gain = config['cacc']['speed_control_gain']
    gap_control_gain_speed = config['cacc']['gap_control_gain_speed']
    gap_control_gain_space = config['cacc']['gap_control_gain_space']
    desired_time_gap = config['cacc']['desired_time_gap']
    gap_closing_gain_speed = config['cacc']['gap_closing_gain_speed']
    gap_closing_gain_space = config['cacc']['gap_closing_gain_space']

    ego_speed: float = ego_data.ego_speed
    desired_speed: float = ego_data.desired_speed

    # Default values for leader gap and speed
    leader_gap = 1000.0
    leader_speed = ego_speed

    # Process leader data based on its type
    if isinstance(leader_data, SumoLGapData):
        if leader_data.leader_gap is not None:
            leader_gap = leader_data.leader_gap
            leader_speed = leader_data.leader_speed
    elif isinstance(leader_data, SumoLPosData):
        if leader_data.leader_position is not None:
            leader_gap = calculate_gap(ego_data, leader_data)
            leader_speed = leader_data.leader_speed

    # Calculate deviations
    gap_deviation: float = leader_gap - ego_speed * desired_time_gap
    speed_deviation: float = leader_speed - ego_speed

    # Determine control mode and calculate acceleration
    if gap_deviation < 0.2 and speed_deviation < 0.1:
        accel: float = gap_control_gain_space * gap_deviation + gap_control_gain_speed * speed_deviation
        mode: str = 'gap_control'
    elif leader_gap < 1.5:
        accel: float = gap_closing_gain_space * gap_deviation + gap_closing_gain_speed * speed_deviation
        mode: str = 'gap_closing_control'
    else:
        accel: float = speed_control_gain * (desired_speed - ego_speed)
        mode: str = 'speed_control'

    # Calculate the new speed for the ego vehicle
    ego_speed_new: float = ego_speed + accel

    return CCData(
        ego_id=ego_data.ego_id,
        mode=mode,
        ego_speed_new=ego_speed_new,
        hostname=hostname
    )
