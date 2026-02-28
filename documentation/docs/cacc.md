# CACC

The CACC (Cooperative Adaptive Cruise Control) simulation is a Python script that simulates a vehicle's behavior in a CACC system.
The script uses Kafka to communicate with other components, such as the SUMO simulation in the system and simulate real-time data exchange.
This is the app that is being migrated by the migration app.
It is a standalone component that can be run independently from other modules in this project. However, it is designed to work in conjunction with the SUMO simulation.
It is necessary to have the Kafka server running to start the script.

It sends the following messages to the Kafka server at topic `control_commands_data`:

- CCData with: ego_id, mode, ego_speed_new, hostname

## How to run

1. start the virtual environment (if not done already): `cd connect-demo` `source ./kafka_python_venv/bin/activate`
2. to start the simulation: `python3 connect-demo/python/cc.py <ego_id>` e.g. `v2`