duration=$1

dummy_path="dummy/bin/dummy"

pkill dummy; pkill dummy; pkill dummy; pkill dummy; pkill dummy

echo "Killed previously running instances"

nohup ./${dummy_path} --name 1 &
nohup ./${dummy_path} --name 2 &
nohup ./${dummy_path} --name 3 &
nohup ./${dummy_path} --name 4 &
nohup ./${dummy_path} --name 5 &

echo "Started 5 dummy replicas"

sleep ${duration}

pkill dummy; pkill dummy; pkill dummy; pkill dummy; pkill dummy

echo "Killed previously running instances"
