# Torture-testing-consensus

This repository implements a networked tool to test the robustness of consensus algorithms.

## Precautions

This project is currently in its **early development stage** and is **not yet ready for use**.

## Requirements

```sudo setcap cap_net_admin,cap_net_raw+ep $(which tc)```
```sudo setcap cap_net_admin=eip /usr/sbin/xtables-nft-multi```
```sudo modprobe ip_tables```
```sudo modprobe nfnetlink_queue```
```sudo apt-get install libcap2-bin```
```sudo apt install -y protobuf-compiler```



