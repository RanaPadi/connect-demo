#!/usr/bin/env python

import os
import sys
import optparse

# we need to import some python modules from the $SUMO_HOME/tools directory
if 'SUMO_HOME' in os.environ:
    tools = os.path.join(os.environ['SUMO_HOME'], 'tools')
    sys.path.append(tools)
else:
    sys.exit("please declare environment variable 'SUMO_HOME'")


from sumolib import checkBinary  # Checks for the binary in environ vars
import traci


from colorama import Fore, Back, Style
#print(Fore.RED + 'some red text')
#print(Back.GREEN + 'and with a green background')
#print(Style.DIM + 'and in dim text')
#print(Style.RESET_ALL)
#print('back to normal now')



def get_options():
    opt_parser = optparse.OptionParser()
    opt_parser.add_option("--nogui", action="store_true",
                         default=False, help="run the commandline version of sumo")
    options, args = opt_parser.parse_args()
    return options


# contains TraCI control loop
def run():
    step = 0
    oldSpeed = 0
    ego = "v2"
    
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
        	
        	#get the ego speed and print colored based on acceleration
        	if egoSpeed < oldSpeed:
        		print (Fore.RED +"Current EGO's speed:"+ str(egoSpeed))
        	else:
        		if egoSpeed > oldSpeed:
        			print (Fore.GREEN +"Current EGO's speed:"+ str(egoSpeed))
        		else:
        			print (Fore.YELLOW +"Current EGO's speed:"+ str(egoSpeed))
        	
        	
        	#If no leader or distance greater than 20, go to max speed.
        	if leader == None or leader[1] >=20.0:
        		AllowedSpeed = traci.vehicle.getAllowedSpeed(ego)
        		traci.vehicle.setSpeed(ego,AllowedSpeed)
        	#If too close
        	else:
        		#wait until have some speed
        		if(egoSpeed >= 3):
        			traci.vehicle.slowDown(ego,egoSpeed-1,5)
        oldSpeed = egoSpeed
        step += 1

    traci.close()
    sys.stdout.flush()


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
