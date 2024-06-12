nohup ./${torture_path} --name 11 --config torture/configuration/local-config.cfg  --replicaConfig torture/configuration/local_consensus_config/11.cfg --debugOn --debugLevel 2 --attacker localNetEm > logs/torture/11.log &
nohup ./${torture_path} --name 12 --config torture/configuration/local-config.cfg  --replicaConfig torture/configuration/local_consensus_config/12.cfg --debugOn --debugLevel 2 --attacker localNetEm > logs/torture/12.log &
nohup ./${torture_path} --name 13 --config torture/configuration/local-config.cfg  --replicaConfig torture/configuration/local_consensus_config/13.cfg --debugOn --debugLevel 2 --attacker localNetEm > logs/torture/13.log &
nohup ./${torture_path} --name 14 --config torture/configuration/local-config.cfg  --replicaConfig torture/configuration/local_consensus_config/14.cfg --debugOn --debugLevel 2 --attacker localNetEm > logs/torture/14.log &
nohup ./${torture_path} --name 15 --config torture/configuration/local-config.cfg  --replicaConfig torture/configuration/local_consensus_config/15.cfg --debugOn --debugLevel 2 --attacker localNetEm > logs/torture/15.log &

./${torture_path} --name 16 --config torture/configuration/local-config.cfg  --debugOn --debugLevel 2 --attacker localNetEm --isController > logs/torture/16.log

