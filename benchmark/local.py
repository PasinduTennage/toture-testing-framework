import time
import random
import threading
import os

class Local:
    def __init__(self, ports, threshold, test_time, view_time=0, benchmark_type="none", epoch_time=0):
        self.benchmark_type = benchmark_type # {none, delay, loss, duplicate, reorder, corrupt, ddoS, parition}
        self.ports = ports                   # list of all local TCP ports to attack
        self.threshold = threshold           # how many nodes are attacked at the same time
        self.test_time = test_time           # test duration in seconds
        self.view_time = view_time           # view change timeout in milliseconds
        self.epoch_time = epoch_time         # epoch timeout in milliseconds -- threshold number of nodes are attacked for epoch time

    def attack(self):
        if self.benchmark_type == "none": 
            self.__attack_none()
        elif self.benchmark_type == "delay" or self.benchmark_type == "loss" or self.benchmark_type == "duplicate" or self.benchmark_type == "reorder" or self.benchmark_type == "corrupt": 
            self.__attack_qc()
        elif self.benchmark_type == "ddoS":
            self.attack_DdoS()
        elif self.benchmark_type == "partition":
            self.attack_partition()        
        else:
            SystemExit("Invalid benchmark type")                        
        
    # sleep for self.time seconds duration
    def __attack_none(self):
        time.sleep(self.test_time)
        return
        
        
    # increase the egress latency of threshold number of nodes at the same time, and keep it for epoch_time, change the set of nodes every epoch_time
    def __attack_qc(self):
        # run for test_time
        start_time = time.time()
        while time.time() - start_time < self.test_time:
            # randomly select threshold number of ports from the ports array and make the attack_nodes array 
            attack_nodes = random.sample(self.ports, self.threshold)
            # for each node in attack nodes, concurrently start attacking it using a thread per node, and then have a barrier to wait for all threads to finish
            n = len(attack_nodes)  # Number of threads
            barrier = threading.Barrier(n)

            threads = []
            for i in range(n):
                start_str = ""
                if self.benchmark_type == "delay":
                    start_str = "tc qdisc add dev lo root handle 1: prio; tc qdisc add dev lo parent 1:3 handle 30: netem delay"+str(int(self.view_time*1.5))+"ms; tc filter add dev lo protocol ip parent 1:0 prio 3 u32 match ip sport "+str(attack_nodes[i])+" 0xffff flowid 1:3"
                elif self.benchmark_type == "parition":
                    start_str = "tc qdisc add dev lo root handle 1: prio; tc qdisc add dev lo parent 1:3 handle 30: netem delay"+str(int(self.view_time*5))+"ms; tc filter add dev lo protocol ip parent 1:0 prio 3 u32 match ip sport "+str(attack_nodes[i])+" 0xffff flowid 1:3"
                
                elif self.benchmark_type == "loss":
                    start_str = "tc qdisc add dev lo root handle 1: prio;  tc qdisc add dev lo parent 1:3 handle 30: netem loss 25%;  tc filter add dev lo protocol ip parent 1:0 prio 3 u32 match ip sport "+str(attack_nodes[i])+" 0xffff flowid 1:3"
                elif self.benchmark_type == "duplicate":
                    start_str = "tc qdisc add dev lo root handle 1: prio;  tc qdisc add dev lo parent 1:3 handle 30: netem duplicate 25%;  tc filter add dev lo protocol ip parent 1:0 prio 3 u32 match ip sport "+str(attack_nodes[i])+" 0xffff flowid 1:3"
                elif self.benchmark_type == "reorder":
                    start_str = "tc qdisc add dev lo root handle 1: prio;  tc qdisc add dev lo parent 1:3 handle 30: netem reorder 25%;  tc filter add dev lo protocol ip parent 1:0 prio 3 u32 match ip sport "+str(attack_nodes[i])+" 0xffff flowid 1:3"
                elif self.benchmark_type == "corrupt":
                    start_str = "tc qdisc add dev lo root handle 1: prio;  tc qdisc add dev lo parent 1:3 handle 30: netem corrupt 25%;  tc filter add dev lo protocol ip parent 1:0 prio 3 u32 match ip sport "+str(attack_nodes[i])+" 0xffff flowid 1:3"            
                else:
                    SystemExit("Invalid benchmark type")

                t = threading.Thread(target=self.__execute, args=(start_str,))
                threads.append(t)
                t.start()

            # Wait for all threads to start the attack
            for t in threads:
                t.join()

            # wait for an epoch time
            time.sleep(self.epoch_time)

            # stop the attack on all nodes
            self.__execute(" tc qdisc del dev lo root")

    # wrapper for Linux tc,string s contains the complete tc command     
    def __execute(self, s):
        # execute the command s in ubuntu
        os.system(s)
        pass