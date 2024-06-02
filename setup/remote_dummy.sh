arrivalRate=$1
exeEnable=$2
txModel=$3
sigType=$4
payloadSize=$5
averageInputMax=$6
averageOutputMax=$7
totalUsers=$8
buffferSize=$9
emptyRequestSize=${10}
clientBatchSize=${11}
replicabatchSize=${12}
iteration=${13}

# A remote test that
#     1. Spawns 5 replicas
#     2. Bootstrap servers
#     3. Spawns 5 clients

pwd=$(pwd)
. "${pwd}"/experiments/remote/ip.sh

rm nohup.out

hotstuff_path="/hotstuff/hotstuff/replica/bin/replica"
ctl_path="/hotstuff/hotstuff/client/bin/client"
config="/home/${username}/hotstuff/hotstuff/configuration/remote/configuration.yml"
output_path="/home/${username}/hotstuff/hotstuff/logs/"


local_output_path="logs/${arrivalRate}/${exeEnable}/${txModel}/${sigType}/${payloadSize}/${averageInputMax}/${averageOutputMax}/${totalUsers}/${buffferSize}/${emptyRequestSize}/${clientBatchSize}/${replicabatchSize}/${iteration}/"
rm -r "${local_output_path}"; mkdir -p "${local_output_path}"

for index in "${!replicas[@]}";
do
  sshpass ssh "${replicas[${index}]}"  -i ${cert}  "pkill replica; pkill client;pkill replica; pkill client;pkill replica; pkill client; rm -r ${output_path}; mkdir ${output_path};mkdir ${output_path}blockchain$(($index + 1))"
  sleep 2
done

echo "Killed previously running instances"

echo "starting replicas"

nohup ssh ${replica1}  -i ${cert}   "sudo ldconfig; .${hotstuff_path} --config ${config} --name 1   --logFilePath ${output_path} --batchSize ${replicabatchSize} --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --buffferSize ${buffferSize} " > "${local_output_path}"1.log &
sleep 2
nohup ssh ${replica2}  -i ${cert}   "sudo ldconfig; .${hotstuff_path} --config ${config} --name 2   --logFilePath ${output_path} --batchSize ${replicabatchSize} --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --buffferSize ${buffferSize} " > "${local_output_path}"2.log &
sleep 2
nohup ssh ${replica3}  -i ${cert}   "sudo ldconfig; .${hotstuff_path} --config ${config} --name 3   --logFilePath ${output_path} --batchSize ${replicabatchSize} --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --buffferSize ${buffferSize} " > "${local_output_path}"3.log &
sleep 2
nohup ssh ${replica4}  -i ${cert}   "sudo ldconfig; .${hotstuff_path} --config ${config} --name 4   --logFilePath ${output_path} --batchSize ${replicabatchSize} --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --buffferSize ${buffferSize} " > "${local_output_path}"4.log &
sleep 2
nohup ssh ${replica5}  -i ${cert}   "sudo ldconfig; .${hotstuff_path} --config ${config} --name 5   --logFilePath ${output_path} --batchSize ${replicabatchSize} --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --buffferSize ${buffferSize} " > "${local_output_path}"5.log &
sleep 2
nohup ssh ${replica6}  -i ${cert}   "sudo ldconfig; .${hotstuff_path} --config ${config} --name 6   --logFilePath ${output_path} --batchSize ${replicabatchSize} --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --buffferSize ${buffferSize} " > "${local_output_path}"6.log &
sleep 2
nohup ssh ${replica7}  -i ${cert}   "sudo ldconfig; .${hotstuff_path} --config ${config} --name 7   --logFilePath ${output_path} --batchSize ${replicabatchSize} --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --buffferSize ${buffferSize} " > "${local_output_path}"7.log &
sleep 2
nohup ssh ${replica8}  -i ${cert}   "sudo ldconfig; .${hotstuff_path} --config ${config} --name 8   --logFilePath ${output_path} --batchSize ${replicabatchSize} --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --buffferSize ${buffferSize} " > "${local_output_path}"8.log &
sleep 2
nohup ssh ${replica9}  -i ${cert}   "sudo ldconfig; .${hotstuff_path} --config ${config} --name 9   --logFilePath ${output_path} --batchSize ${replicabatchSize} --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --buffferSize ${buffferSize} " > "${local_output_path}"9.log &
sleep 2
nohup ssh ${replica10} -i ${cert}   "sudo ldconfig; .${hotstuff_path} --config ${config} --name 10  --logFilePath ${output_path} --batchSize ${replicabatchSize} --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --buffferSize ${buffferSize} " > "${local_output_path}"10.log &
sleep 2
nohup ssh ${replica11} -i ${cert}   "sudo ldconfig; .${hotstuff_path} --config ${config} --name 11  --logFilePath ${output_path} --batchSize ${replicabatchSize} --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --buffferSize ${buffferSize} " > "${local_output_path}"11.log &
sleep 2
nohup ssh ${replica12} -i ${cert}   "sudo ldconfig; .${hotstuff_path} --config ${config} --name 12  --logFilePath ${output_path} --batchSize ${replicabatchSize} --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --buffferSize ${buffferSize} " > "${local_output_path}"12.log &

echo "Started replicas"

sleep 15

nohup ssh ${replica13} -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 21  --logFilePath ${output_path} --requestType status --operationType 1 --exeEnable 0 " > "${local_output_path}"status1.log &

echo "Sent initial status to bootstrap"

sleep 30

nohup ssh ${replica14} -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 22  --logFilePath ${output_path} --requestType status --operationType 3 --exeEnable 0 " > "${local_output_path}"status3.log &

echo "Sent hotstuff start indication"

sleep 30

echo "Starting client[s]"

nohup ssh ${replica13}  -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 21  --logFilePath ${output_path} --batchSize ${clientBatchSize} --arrivalRate ${arrivalRate} --requestType request --designatedReplica 1  --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --emptyRequestSize ${emptyRequestSize} " > "${local_output_path}"21.log &
sleep 2
nohup ssh ${replica14}  -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 22  --logFilePath ${output_path} --batchSize ${clientBatchSize} --arrivalRate ${arrivalRate} --requestType request --designatedReplica 2  --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --emptyRequestSize ${emptyRequestSize} " > "${local_output_path}"22.log &
sleep 2
nohup ssh ${replica15}  -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 23  --logFilePath ${output_path} --batchSize ${clientBatchSize} --arrivalRate ${arrivalRate} --requestType request --designatedReplica 3  --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --emptyRequestSize ${emptyRequestSize} " > "${local_output_path}"23.log &
sleep 2
nohup ssh ${replica16}  -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 24  --logFilePath ${output_path} --batchSize ${clientBatchSize} --arrivalRate ${arrivalRate} --requestType request --designatedReplica 4  --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --emptyRequestSize ${emptyRequestSize} " > "${local_output_path}"24.log &
sleep 2
nohup ssh ${replica17}  -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 25  --logFilePath ${output_path} --batchSize ${clientBatchSize} --arrivalRate ${arrivalRate} --requestType request --designatedReplica 5  --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --emptyRequestSize ${emptyRequestSize} " > "${local_output_path}"25.log &
sleep 2
nohup ssh ${replica18}  -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 26  --logFilePath ${output_path} --batchSize ${clientBatchSize} --arrivalRate ${arrivalRate} --requestType request --designatedReplica 6  --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --emptyRequestSize ${emptyRequestSize} " > "${local_output_path}"26.log &
sleep 2
nohup ssh ${replica19}  -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 27  --logFilePath ${output_path} --batchSize ${clientBatchSize} --arrivalRate ${arrivalRate} --requestType request --designatedReplica 7  --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --emptyRequestSize ${emptyRequestSize} " > "${local_output_path}"27.log &
sleep 2
nohup ssh ${replica20}  -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 28  --logFilePath ${output_path} --batchSize ${clientBatchSize} --arrivalRate ${arrivalRate} --requestType request --designatedReplica 8  --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --emptyRequestSize ${emptyRequestSize} " > "${local_output_path}"28.log &
sleep 2
nohup ssh ${replica21}  -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 29  --logFilePath ${output_path} --batchSize ${clientBatchSize} --arrivalRate ${arrivalRate} --requestType request --designatedReplica 9  --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --emptyRequestSize ${emptyRequestSize} " > "${local_output_path}"29.log &
sleep 2
nohup ssh ${replica22}  -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 30  --logFilePath ${output_path} --batchSize ${clientBatchSize} --arrivalRate ${arrivalRate} --requestType request --designatedReplica 10 --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --emptyRequestSize ${emptyRequestSize} " > "${local_output_path}"30.log &
sleep 2
nohup ssh ${replica23}  -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 31  --logFilePath ${output_path} --batchSize ${clientBatchSize} --arrivalRate ${arrivalRate} --requestType request --designatedReplica 11 --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --emptyRequestSize ${emptyRequestSize} " > "${local_output_path}"31.log &
sleep 2
nohup ssh ${replica24}  -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 32  --logFilePath ${output_path} --batchSize ${clientBatchSize} --arrivalRate ${arrivalRate} --requestType request --designatedReplica 12 --exeEnable ${exeEnable} --txModel ${txModel} --sigType ${sigType} --payloadSize ${payloadSize} --averageInputMax ${averageInputMax} --averageOutputMax ${averageOutputMax} --totalUsers ${totalUsers} --emptyRequestSize ${emptyRequestSize} " > "${local_output_path}"32.log &

sleep 300

echo "Completed Client[s]"

nohup ssh ${replica13} -i ${cert}   "sudo ldconfig; .${ctl_path} --config ${config} --name 21  --logFilePath ${output_path} --requestType status --operationType 4 --exeEnable 0 " > "${local_output_path}"status4.log &

sleep 20

echo "Finish test"