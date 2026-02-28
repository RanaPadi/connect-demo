#!/usr/bin/env python

import os
import sys
import optparse
from colorama import Fore, Back, Style
import json
from kafka import KafkaConsumer, KafkaProducer

############################################

# we need to import some python modules from the $SUMO_HOME/tools directory
if 'SUMO_HOME' in os.environ:
    tools = os.path.join(os.environ['SUMO_HOME'], 'tools')
    sys.path.append(tools)
else:
    sys.exit("Please declare environment variable 'SUMO_HOME'")

from sumolib import checkBinary  # Checks for the binary in environ vars
import traci

def get_options():
    opt_parser = optparse.OptionParser()
    opt_parser.add_option("--nogui", action="store_true",
                         default=False, help="run the commandline version of sumo")
    options, args = opt_parser.parse_args()
    return options

############################################

# Kafka configuration
KAFKA_BOOTSTRAP_SERVERS = 'localhost:9092'
VEHICLE_STATE_TOPIC = 'sumo_info'
ACC_CONTROL_COMMANDS_TOPIC = 'acc_control_commands'

# Kafka consumer and producer
consumer = KafkaConsumer(ACC_CONTROL_COMMANDS_TOPIC,
                         bootstrap_servers=[KAFKA_BOOTSTRAP_SERVERS],
                         value_deserializer=lambda m: json.loads(m.decode('utf-8')))

producer = KafkaProducer(bootstrap_servers=[KAFKA_BOOTSTRAP_SERVERS],
                         value_serializer=lambda m: json.dumps(m).encode('utf-8'))

############################################

# contains TraCI control loop
def run(ego):
    sumo_info = {}
    control_command = {}
    
    while traci.simulation.getMinExpectedNumber() > 0:
        traci.simulationStep()

        # Check if the Ego vehicle is present in the simulation
        vehIDList = traci.vehicle.getIDList()
        if ego not in vehIDList:
            continue

        # Get the distance from the Ego to the leader and other information
        leader = traci.vehicle.getLeader(ego, 100)
        egoSpeed = traci.vehicle.getSpeed(ego)
        sumo_info = {
            'ego_id': ego,
            'ego_speed': egoSpeed,
            'leader_gap': leader[1],
            'leader_speed': traci.vehicle.getSpeed(leader[0]),
            'desired_speed': traci.vehicle.getAllowedSpeed(ego)
        }
        producer.send(VEHICLE_STATE_TOPIC, value=sumo_info)

        # Poll messages with a short timeout to avoid blocking the simulation
        polled_messages = consumer.poll(timeout_ms=100)  # Adjust the timeout as needed

        # Process polled messages
        for tp, msgs in polled_messages.items():
            for message in msgs:
                control_command = message.value
                if control_command['ego_id'] == ego:
                    ego_speed_new = control_command['ego_speed_new']
                    mode = control_command['mode']
                    # Process the control command
                    if mode == 'speed_control':
                        print(f"{Fore.GREEN} [Vehicle {control_command['ego_id']}] [{control_command['mode']}] Speed {control_command['ego_speed_new']}") 
                    elif mode == 'gap_control':
                        print(f"{Fore.RED} [Vehicle {control_command['ego_id']}] [{control_command['mode']}] Speed {control_command['ego_speed_new']}")

                    traci.vehicle.setSpeed(ego, ego_speed_new)


    traci.close()
    sys.stdout.flush()

############################################

# main entry point
if __name__ == "__main__":
    
    
    options = get_options()

    # check binary
    if options.nogui:
        sumoBinary = checkBinary('sumo')
    else:
        sumoBinary = checkBinary('sumo-gui')

    print("--- Starting SUMO Demo ---")
    ego = input("Vehicle ID:" )

# traci starts sumo as a subprocess and then this script connects and runs
    traci.start([sumoBinary, "-c", "demo.sumocfg",
                             "--tripinfo-output", "tripinfo.xml"])
    
    
    
    run(ego)
