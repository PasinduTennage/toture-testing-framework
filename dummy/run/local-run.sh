timeout=$1

dummy_path="dummy/bin/dummy"

pkill dummy; pkill dummy; pkill dummy; pkill dummy; pkill dummy

echo "Killed previously running instances"

nohup ./${dummy_path} --name 1 --ui &
nohup ./${dummy_path} --name 2 --ui &
nohup ./${dummy_path} --name 3 --ui &
nohup ./${dummy_path} --name 4 --ui &
nohup ./${dummy_path} --name 5 --ui &

echo "Started 5 dummy replicas"

sleep ${timeout}

pkill dummy; pkill dummy; pkill dummy; pkill dummy; pkill dummy
rm nohup.out
echo "Killed all dummy servers"