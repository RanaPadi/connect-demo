# SUMO simulation

The SUMO simulation is a Python script that simulates a vehicle's behavior in a SUMO simulation.
The script uses Kafka to communicate with the CACC simulation.
It is a standalone component that can be run independently from other modules in this project. However, it is designed to work in conjunction with the CACC simulation.
It is necessary to have the Kafka server running to start the script.

It sends the following messages to the Kafka server at topic `ego_data`, `leader_position_data`, `leader_gap_data`:

- SumoData with: ego_id, mode, ego_speed_new, hostname
- SumoEgoData with: ego_id, ego_speed, ego_position, desired_speed, leader_id
- SumoLGapData with: ego_id, leader_id, leader_speed, leader_gap

the latter two depending on which one is needed. To switch topics, you need to send any Kafka message to the `switch_data`.

## How to run

1. start the virtual environment (if not done already): `cd connect-demo` `source ./kafka_python_venv/bin/activate`
2. Start Sumo using:
    1. `cd connect-demo/python/sumo`
    2. `python3 sumo/demo.py`
