#!/usr/bin/env python

import os
import sys
import optparse
import json
from kafka import KafkaConsumer, KafkaProducer
from colorama import Fore, Back, Style
from sumolib import checkBinary
import traci

############################################

# Kafka configuration
KAFKA_BOOTSTRAP_SERVERS = '134.60.77.96:9092'
VEHICLE_STATE_TOPIC = 'sumo_info'
ACC_CONTROL_COMMANDS_TOPIC = 'acc_control_commands'
TIMEOUT = 100  # Adjust the timeout as needed

# Kafka consumer and producer
consumer = KafkaConsumer(ACC_CONTROL_COMMANDS_TOPIC,
                         bootstrap_servers=[KAFKA_BOOTSTRAP_SERVERS],
                         value_deserializer=lambda m: json.loads(m.decode('utf-8')))


producer = KafkaProducer(bootstrap_servers=[KAFKA_BOOTSTRAP_SERVERS],
                         value_serializer=lambda m: json.dumps(m).encode('utf-8'))

############################################

# we need to import some python modules from the $SUMO_HOME/tools directory
if 'SUMO_HOME' in os.environ:
    tools = os.path.join(os.environ['SUMO_HOME'], 'tools')
    sys.path.append(tools)
else:
    sys.exit("Please declare environment variable 'SUMO_HOME'")

def get_options():
    opt_parser = optparse.OptionParser()
    opt_parser.add_option("--nogui", action="store_true",
                         default=False, help="run the commandline version of sumo")
    options, args = opt_parser.parse_args()
    return options

############################################

# contains TraCI control loopiimport optparse
import json
from kafka import KafkaConsumer, KafkaProducer
from colorama import Fore, Back, Style
from sumolib import checkBinary
import traci

############################################

# Kafka configuration
KAFKA_BOOTSTRAP_SERVERS = '134.60.77.96:9092'
VEHICLE_STATE_TOPIC = 'sumo_info'
ACC_CONTROL_COMMANDS_TOPIC = 'acc_control_commands'
TIMEOUT = 100  # Adjust the timeout as needed

# Kafka consumer and producer
consumer = KafkaConsumer(ACC_CONTROL_COMMANDS_TOPIC,
                         bootstrap_servers=[KAFKA_BOOTSTRAP_SERVERS],
                         value_deserializer=lambda m: json.loads(m.decode('utf-8')))


producer = KafkaProducer(bootstrap_servers=[KAFKA_BOOTSTRAP_SERVERS],
                         value_serializer=lambda m: json.dumps(m).encode('utf-8'))

############################################

# we need to import some python modules from the $SUMO_HOME/tools directory
if 'SUMO_HOME' in os.environ:
    tools = os.path.join(os.environ['SUMO_HOME'], 'tools')
    sys.path.append(tools)
else:
    sys.exit("Please declare environment variable 'SUMO_HOME'")

def get_options():
    opt_parser = optparse.OptionParser()
    opt_parser.add_option("--nogui", action="store_true",
                         default=False, help="run the commandline version of sumo")
    options, args = opt_parser.parse_args()
    return options

############################################

# contains TraCI control loop
def run():
    try:
        while traci.simulation.getMinExpectedNumber() > 0:
            traci.simulationStep()
            vehIDList = traci.vehicle.getIDList()

            for ego in vehIDList:
                ego_speed = traci.vehicle.getSpeed(ego)
                
                # Get the distance from the vehicle to the leader and other information
                leader = traci.vehicle.getLeader(ego, 10000)
                if leader is not None:
                    leader_speed = traci.vehicle.getSpeed(leader[0])
                else:
                    leader = (None, 1000)
                    leader_speed = ego_speed
                
                sumo_info = {
                    'ego_id': ego,
                    'ego_speed': ego_speed,
                    'leader_gap': leader[1],
                    'leader_speed': leader_speed
                }
                producer.send(VEHICLE_STATE_TOPIC, value=sumo_info)

                # Poll messages with a short timeout to avoid blocking the simulation
                polled_messages = consumer.poll(timeout_ms=TIMEOUT)  # Adjust the timeout as needed

                # Process polled messages
                for tp, msgs in polled_messages.items():
                    for message in msgs:
                        control_command = message.value
                        if control_command['ego_id'] == ego:
                            ego_speed_new = control_command['ego_speed_new']
                            desired_speed = control_command['desired_speed']
                            leader_gap = control_command['leader_gap']
                            mode = control_command['mode']

                            # Process the control command
                            if mode == 'speed_control':
                                print(f"{Fore.GREEN} [Vehicle {ego}] [{mode}] | Speed: {ego_speed_new} | Desired Speed: {desired_speed} | Leader Gap: {leader_gap}")
                            elif mode == 'gap_control':
                                print(f"{Fore.RED} [Vehicle {ego}] [{mode}] | Speed: {ego_speed_new} | Desired Speed: {desired_speed} | Leader Gap: {leader_gap}")

                            traci.vehicle.setSpeed(ego, ego_speed_new)
                            traci.vehicle.setMaxSpeed(ego, desired_speed)
                            
        traci.close()
        sys.stdout.flush()
    except traci.exceptions.FatalTraCIError as e:
        print(f"FatalTraCIError occurred: {e}")
    except KeyboardInterrupt:
        print("Simulation stopped by user...")

############################################

# main entry point
if __name__ == "__main__":
    
    options = get_options()

    # check binary
    if options.nogui:
        sumoBinary = checkBinary('sumo')
    else:
        sumoBinary = checkBinary('sumo-gui')
        
    print(f"Starting SUMO with {sumoBinary}")

    # traci starts sumo as a subprocess and then this script connects and runs
    traci.start([sumoBinary, "-c", "demo.sumocfg",
                             "--tripinfo-output", "tripinfo.xml"])
    
    run()
