timeout=$1

dummy_path="dummy/bin/dummy"

pkill dummy; pkill dummy; pkill dummy; pkill dummy; pkill dummy

echo "Killed previously running instances"

nohup ./${dummy_path} --name 1 --ui &
nohup ./${dummy_path} --name 2 --ui &
nohup ./${dummy_path} --name 3 --ui &
nohup ./${dummy_path} --name 4 --ui &
nohup ./${dummy_path} --name 5 --ui &

sleep 10

xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=10100"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=20100"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=30100"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=40100"
xdg-open "http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=50100"

echo "Started 5 dummy replicas"

sleep ${timeout}

pkill dummy; pkill dummy; pkill dummy; pkill dummy; pkill dummy
rm nohup.out
echo "Killed all dummy servers"