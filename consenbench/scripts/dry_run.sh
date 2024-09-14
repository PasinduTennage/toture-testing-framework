git pull origin main
/bin/bash build.sh
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 0 --controller_operation_type bootstrap  --consensus_algorithm baxos  --attack_duration 60 --attack basic --device enp1s0
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 0 --controller_operation_type copy       --consensus_algorithm baxos  --attack_duration 60 --attack basic --device enp1s0
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 0 --controller_operation_type run        --consensus_algorithm baxos  --attack_duration 60 --attack basic --device enp1s0
