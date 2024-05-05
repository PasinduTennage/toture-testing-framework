nohup ./dummy/bin/dummy --ports "10001,10002,10003,10004,10005" &
nohup ./dummy/bin/dummy --ports "11001,11002,11003,11004,11005" &
nohup ./dummy/bin/dummy --ports "12001,12002,12003,12004,12005" &
nohup ./dummy/bin/dummy --ports "13001,13002,13003,13004,13005" &
nohup ./dummy/bin/dummy --ports "14001,14002,14003,14004,14005" &

sleep 200

pkill dummy
pkill dummy
pkill dummy
pkill dummy
pkill dummy