#!/usr/bin/env python

import os
import sys
import optparse
from sumolib import checkBinary  # Checks for the binary in environ vars
import traci
from colorama import Fore, Back, Style
import matplotlib.pyplot as plt
from collections import deque


# we need to import some python modules from the $SUMO_HOME/tools directory
if 'SUMO_HOME' in os.environ:
    tools = os.path.join(os.environ['SUMO_HOME'], 'tools')
    sys.path.append(tools)
else:
    sys.exit("please declare environment variable 'SUMO_HOME'")
    
############################################    

def get_options():
    opt_parser = optparse.OptionParser()
    opt_parser.add_option("--nogui", action="store_true",
                         default=False, help="run the commandline version of sumo")
    options, args = opt_parser.parse_args()
    return options

############################################

SPEED_CONTROL_GAIN = 0.4  # This is a simplification; actual gain should be based on vehicle and conditions
GAP_CONTROL_GAIN_SPEED = 0.07
GAP_CONTROL_GAIN_SPACE = 0.23
DESIRED_TIME_GAP = 1.1

def calculate_acc_control(ego_id, sumo_info):
    # Placeholder for speed and gap control calculation based on vehicle states
    ego_speed = sumo_info['ego_speed']
    # leader_id = sumo_info['leader_id']
    leader_gap = sumo_info['leader_gap']
    leader_speed = sumo_info['leader_speed']
    desired_speed = sumo_info['desired_speed']

    
##### acc logic according to porfyri et al. 2018

    gap_deviation = leader_gap - ego_speed * DESIRED_TIME_GAP
    speed_deviation = leader_speed - ego_speed
    
    if gap_deviation < 0.2 and speed_deviation < 0.1:
        accel = GAP_CONTROL_GAIN_SPACE * gap_deviation + GAP_CONTROL_GAIN_SPEED * speed_deviation
        mode = 'speed_control'
    else:
        accel = SPEED_CONTROL_GAIN * (desired_speed - ego_speed)
        mode = 'gap_control'
    
    ego_speed_new = ego_speed + accel
    
    return {'ego_id': ego_id, 'ego_speed_new': ego_speed_new, 'mode': mode}

############################################

# contains TraCI control loop
def run():
    step = 0
    ego = "v2"
    sumo_info = {}
    # MAX NO. OF POINTS TO STORE for plotting
    que = deque(maxlen = 40)
    
    while traci.simulation.getMinExpectedNumber() > 0:
        traci.simulationStep()
        #print(step)
        
        #get the list of vehicles and check if the Ego is there        
        vehIDList = traci.vehicle.getIDList()
        if ego not in vehIDList: continue
        else:
            
            #get the distance from the Ego to the leader
            leader = traci.vehicle.getLeader(ego,100)
            print(leader)
            
            #get the ego speed
            egoSpeed = traci.vehicle.getSpeed(ego)

            sumo_info['ego_speed'] = egoSpeed
            sumo_info['leader_gap'] = leader[1]
            sumo_info['leader_speed'] = traci.vehicle.getSpeed(leader[0])
            sumo_info['desired_speed'] = traci.vehicle.getAllowedSpeed(ego)
            
            #calculate speed
            control_command = calculate_acc_control(ego, sumo_info)
            
            ego_speed_new = control_command['ego_speed_new']
            mode = control_command['mode']
            
            if mode == 'speed_control':
                print (Fore.GREEN + "[Speed control mode] Current EGO speed: " + str(ego_speed_new))
            if mode == 'gap_control':
                print (Fore.RED + "[Gap control mode] Current EGO speed: " + str(ego_speed_new))
            
            traci.vehicle.setSpeed(ego, ego_speed_new)
            
            ##plot speed
            #que.append(ego_speed_new)
            #plt.plot(que)
            #plt.scatter(range(len(que)),que) #insert points on top of each speed
            #plt.autoscale
            ##draw and clear plot
            #plt.title("Vehicle speed")  # Add a title to the axes.
            #plt.draw()
            #plt.pause(1e-17)
            #plt.clf()            
            
            
            
            
            
            #plot speed
            que.append(ego_speed_new)
            plt.plot(que)
            plt.scatter(range(len(que)),que) #insert points on top of each speed
            plt.autoscale
            #draw and clear plot
            plt.title("Vehicle speed")  # Add a title to the axes.
            plt.draw()
            plt.pause(1e-17)
            plt.clf() 
            
            
        step += 1

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

    # traci starts sumo as a subprocess and then this script connects and runs
    traci.start([sumoBinary, "-c", "demo.sumocfg",
                             "--tripinfo-output", "tripinfo.xml"])
    run()
