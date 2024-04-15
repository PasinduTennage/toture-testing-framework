import time
import random
import threading
import os
import psutil

def get_open_ports(server_name):
    # Get all processes for server_name
    all_processes = psutil.process_iter(['pid', 'name'])

    # Filter processes based on server_name
    server_processes = [p.info['pid'] for p in all_processes if p.info['name'] == server_name]

    # For each server process, get the set of open ports
    open_ports = []
    for pid in server_processes:
        process = psutil.Process(pid)
        connections = process.connections()
        ports = [conn.laddr.port for conn in connections if conn.status == 'LISTEN']
        open_ports.append(ports)

    return open_ports

# class Local is used to attack the local machine
# params
# server_name: name of the server process
# threshold: how many nodes are attacked at the same time
# test_time: test duration in seconds
# view_time: view change timeout in milliseconds
# benchmark_type: {none, delay, loss, duplicate, reorder, corrupt, parition}
# epoch_time: epoch timeout in milliseconds -- threshold number of nodes are attacked for epoch time


class Local:
    def __init__(self, server_name, threshold, test_time, view_time=0, benchmark_type="none", epoch_time=0):
        self.ports = get_open_ports(server_name)
        self.threshold = threshold
        self.test_time = test_time
        self.view_time = view_time
        self.benchmark_type = benchmark_type
        self.epoch_time = epoch_time

    # external interface to start the attack
    def attack(self):
        if self.benchmark_type == "none": 
            self.__attack_none()
        elif self.benchmark_type == "delay" or self.benchmark_type == "loss" or self.benchmark_type == "duplicate" or self.benchmark_type == "reorder" or self.benchmark_type == "corrupt" or self.benchmark_type == "partition": 
            self.__attack_qc()   
        else:
            SystemExit("Invalid benchmark type")     
    
    # external interface to stop the attack by removing all qc rules
    def stopAttack(self):
        self.__execute("tc qdisc del dev lo root")
        return                   
        
    # sleep for self.time seconds duration
    def __attack_none(self):
        time.sleep(self.test_time)
        return
        
        
    # attack threshold number of nodes at the same time, and keep it for epoch_time, change the set of nodes every epoch_time
    def __attack_qc(self):
        # run for test_time
        start_time = time.time()
        while time.time() - start_time < self.test_time:
            # randomly select threshold number of processes from the all available replicas and make the attack_nodes array 
            attack_nodes = random.sample(self.ports, self.threshold)
            # for each node in attack nodes, concurrently start attacking it using a thread per port of a replica, and then have a barrier to wait for all threads to finish
            n = len(attack_nodes)  # number of processes to attack
            threads = []
            for i in range(n):
                for j in range(len(attack_nodes[i])):
                    start_str = ""
                    if self.benchmark_type == "delay":
                        start_str = "tc qdisc add dev lo root handle 1: prio; tc qdisc add dev lo parent 1:3 handle 30: netem delay "+str(int(self.view_time*1.5))+"ms; tc filter add dev lo protocol ip parent 1:0 prio 3 u32 match ip sport "+str(attack_nodes[i][j])+" 0xffff flowid 1:3"
                    elif self.benchmark_type == "parition":
                        start_str = "tc qdisc add dev lo root handle 1: prio; tc qdisc add dev lo parent 1:3 handle 30: netem delay" +str(int(self.view_time*5))+"ms; tc filter add dev lo protocol ip parent 1:0 prio 3 u32 match ip sport "+str(attack_nodes[i][j])+" 0xffff flowid 1:3"
                    elif self.benchmark_type == "loss":
                        start_str = "tc qdisc add dev lo root handle 1: prio;  tc qdisc add dev lo parent 1:3 handle 30: netem loss 25%;  tc filter add dev lo protocol ip parent 1:0 prio 3 u32 match ip sport "+str(attack_nodes[i][j])+" 0xffff flowid 1:3"
                    elif self.benchmark_type == "duplicate":
                        start_str = "tc qdisc add dev lo root handle 1: prio;  tc qdisc add dev lo parent 1:3 handle 30: netem duplicate 25%;  tc filter add dev lo protocol ip parent 1:0 prio 3 u32 match ip sport "+str(attack_nodes[i][j])+" 0xffff flowid 1:3"
                    elif self.benchmark_type == "reorder":
                        start_str = "tc qdisc add dev lo root handle 1: prio;  tc qdisc add dev lo parent 1:3 handle 30: netem reorder 25%;  tc filter add dev lo protocol ip parent 1:0 prio 3 u32 match ip sport "+str(attack_nodes[i][j])+" 0xffff flowid 1:3"
                    elif self.benchmark_type == "corrupt":
                        start_str = "tc qdisc add dev lo root handle 1: prio;  tc qdisc add dev lo parent 1:3 handle 30: netem corrupt 25%;  tc filter add dev lo protocol ip parent 1:0 prio 3 u32 match ip sport "+str(attack_nodes[i][j])+" 0xffff flowid 1:3"            
                    else:
                        SystemExit("Invalid benchmark type")

                    t = threading.Thread(target=self.__execute, args=(start_str,))
                    threads.append(t)
                    t.start()

            # Wait for all threads to start the attack
            for t in threads:
                t.join()

            #  sleep for epoch_time number of miliseconds
            time.sleep(self.epoch_time/1000)
            

            # stop the attack on all nodes
            self.stopAttack()

    # wrapper for Linux tc,string s contains the complete tc command     
    def __execute(self, s):
        # execute the command s in ubuntu
        os.system(s)
        pass