interArrivalTime=$1

dummy_path="dummy/bin/dummy"
torture_path="torture/bin/torture"

/bin/bash build.sh

pkill dummy; pkill dummy; pkill dummy; pkill dummy; pkill dummy
pkill torture; pkill torture; pkill torture; pkill torture; pkill torture
rm nohup.out

echo "Killed previously running instances"

mkdir -p logs/dummy
mkdir -p logs/torture

nohup ./${dummy_path} --name 1 --ui --interArrivalTime ${interArrivalTime} > logs/dummy/1.log &
nohup ./${dummy_path} --name 2 --ui --interArrivalTime ${interArrivalTime} > logs/dummy/2.log &
nohup ./${dummy_path} --name 3 --ui --interArrivalTime ${interArrivalTime} > logs/dummy/3.log &
nohup ./${dummy_path} --name 4 --ui --interArrivalTime ${interArrivalTime} > logs/dummy/4.log &
nohup ./${dummy_path} --name 5 --ui --interArrivalTime ${interArrivalTime} > logs/dummy/5.log &

sleep 10

xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=10200"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=20200"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=30200"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=40200"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=50200"

echo "Started 5 dummy replicas"

sudo /bin/bash dummy/run/local-run-torture.sh

sleep 10

pkill dummy; pkill dummy; pkill dummy; pkill dummy; pkill dummy
pkill torture; pkill torture; pkill torture; pkill torture; pkill torture

rm nohup.out
echo "Finished tests"