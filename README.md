# toture-testing-consensus



## Requirements

This tool uses Linux traffic control (```tc```), and assumes that a non-root user can run the ```tc``` command. 

To enable that, run ```sudo setcap cap_net_admin,cap_net_raw+ep $(which tc)```
