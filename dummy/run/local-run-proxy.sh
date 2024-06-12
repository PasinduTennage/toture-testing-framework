dummy_path="dummy/bin/dummy"
torture_path="torture/bin/torture"

/bin/bash build.sh

pkill dummy; pkill dummy; pkill dummy; pkill dummy; pkill dummy
pkill torture; pkill torture; pkill torture; pkill torture; pkill torture
rm nohup.out
echo "Killed previously running instances"

mkdir -p logs/dummy
mkdir -p logs/torture

nohup ./${dummy_path} --name 1 --ui --config dummy/configuration/local-proxy/1.cfg --debugOn --debugLevel 12 --interArrivalTime 10> logs/dummy/1.log &
nohup ./${dummy_path} --name 2 --ui --config dummy/configuration/local-proxy/2.cfg --debugOn --debugLevel 12 --interArrivalTime 10 > logs/dummy/2.log &
nohup ./${dummy_path} --name 3 --ui --config dummy/configuration/local-proxy/3.cfg --debugOn --debugLevel 12 --interArrivalTime 10 > logs/dummy/3.log &
nohup ./${dummy_path} --name 4 --ui --config dummy/configuration/local-proxy/4.cfg --debugOn --debugLevel 12 --interArrivalTime 10 > logs/dummy/4.log &
nohup ./${dummy_path} --name 5 --ui --config dummy/configuration/local-proxy/5.cfg --debugOn --debugLevel 12 --interArrivalTime 10 > logs/dummy/5.log &

echo "Started 5 dummy replicas"

sleep 10

xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=10200"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=20200"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=30200"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=40200"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=50200"

echo "Started 5 uis"

nohup ./${torture_path} --name 11 --config torture/configuration/local-config.cfg  --replicaConfig torture/configuration/local_consensus_config/11.cfg --debugOn --debugLevel 2 --attacker lProxy > logs/torture/11.log &
nohup ./${torture_path} --name 12 --config torture/configuration/local-config.cfg  --replicaConfig torture/configuration/local_consensus_config/12.cfg --debugOn --debugLevel 2 --attacker lProxy > logs/torture/12.log &
nohup ./${torture_path} --name 13 --config torture/configuration/local-config.cfg  --replicaConfig torture/configuration/local_consensus_config/13.cfg --debugOn --debugLevel 2 --attacker lProxy > logs/torture/13.log &
nohup ./${torture_path} --name 14 --config torture/configuration/local-config.cfg  --replicaConfig torture/configuration/local_consensus_config/14.cfg --debugOn --debugLevel 2 --attacker lProxy > logs/torture/14.log &
nohup ./${torture_path} --name 15 --config torture/configuration/local-config.cfg  --replicaConfig torture/configuration/local_consensus_config/15.cfg --debugOn --debugLevel 2 --attacker lProxy > logs/torture/15.log &

./${torture_path} --name 16 --config torture/configuration/local-config.cfg  --debugOn --debugLevel 2 --attacker lProxy --isController > logs/torture/16.log

sleep 10

pkill dummy; pkill dummy; pkill dummy; pkill dummy; pkill dummy
pkill torture; pkill torture; pkill torture; pkill torture; pkill torture

rm nohup.out
echo "Finished tests"