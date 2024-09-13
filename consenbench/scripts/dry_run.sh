git pull origin main
/bin/bash build.sh
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 3 --controller_operation_type bootstrap --log_file /home/pasindu/bench/log.log
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 3 --controller_operation_type copy      --log_file /home/pasindu/bench/log.log --consensus_algorithm baxos
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 3 --controller_operation_type run       --log_file /home/pasindu/bench/log.log --consensus_algorithm baxos  --attack_duration 60 --attack basic
