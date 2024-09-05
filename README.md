# toture-testing-consensus

This repository implements a tool to test the consensus algorithms.
The tool can be used to test the consensus algorithm under different fault conditions.

## Requirements

```sudo setcap cap_net_admin,cap_net_raw+ep $(which tc)```
```sudo setcap cap_net_admin=eip  /usr/sbin/xtables-nft-multi```
```sudo modprobe ip_tables```
```sudo modprobe nfnetlink_queue```
```sudo apt-get install libcap2-bin```

## Precautions

Run this program in a VM, to avoid any problem in your host machine

This project is currently under development, and is ""not"" ready for production use.
