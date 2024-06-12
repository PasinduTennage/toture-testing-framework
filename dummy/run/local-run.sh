timeout=$1
interArrivalTime=$2

dummy_path="dummy/bin/dummy"

pkill dummy; pkill dummy; pkill dummy; pkill dummy; pkill dummy

echo "Killed previously running instances"

nohup ./${dummy_path} --name 1 --ui --interArrivalTime 100 &
nohup ./${dummy_path} --name 2 --ui --interArrivalTime 100 &
nohup ./${dummy_path} --name 3 --ui --interArrivalTime 100 &
nohup ./${dummy_path} --name 4 --ui --interArrivalTime 100 &
nohup ./${dummy_path} --name 5 --ui --interArrivalTime 100 &

sleep 10

xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=10200"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=20200"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=30200"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=40200"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=50200"

echo "Started 5 dummy replicas"

sleep ${timeout}

pkill dummy; pkill dummy; pkill dummy; pkill dummy; pkill dummy
rm nohup.out
echo "Killed all dummy servers"