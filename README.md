# toture-testing-consensus

## Requirements

This tool uses Linux traffic control (```tc```), and assumes that a non-root user can run the ```tc``` command. 

To enable that, run ```sudo setcap cap_net_admin,cap_net_raw+ep $(which tc)```

```sudo setcap cap_net_admin=eip  /usr/sbin/xtables-nft-multi```

sudo modprobe ip_tables
sudo modprobe nfnetlink_queue
sudo apt-get install libcap2-bin

## Precautions

Run this program in a VM, to avoid any problem in your host machine

## Dummy
