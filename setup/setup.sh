# zip the project into a binaries.zip
# copy the zip to remote
# extract the zip file, install the c library, build hotstuff
# create the logs/

pwd=$(pwd)
. "${pwd}"/experiments/remote/ip.sh

rm client/bin/client
rm replica/bin/replica

echo "Removed old binaries"

rm -r logs/ ; mkdir logs/

mkdir -p logs/temp/hotstuff

cp -r client/                     logs/temp/hotstuff
cp -r common/                     logs/temp/hotstuff
cp -r configuration/              logs/temp/hotstuff
cp -r proto/                      logs/temp/hotstuff
cp -r replica/                    logs/temp/hotstuff
cp    go.mod                      logs/temp/hotstuff
cp    go.sum                      logs/temp/hotstuff

cp    replica/execbin/exelayers.zip logs/temp/

zip -r logs/binaries.zip logs/temp/
rm -r logs/temp/

reset_directory="sudo rm -r /home/${username}/hotstuff; mkdir /home/${username}/hotstuff"
kill_instances="pkill replica ; pkill client"

local_zip_path="logs/binaries.zip"
remote_home_path="/home/${username}/hotstuff/"

command1="sudo apt-get update; sudo apt-get install unzip;sudo apt-get install zip; cd /home/${username}/hotstuff && unzip binaries.zip"
command2="mv /home/${username}/hotstuff/logs/temp/exelayers.zip /home/${username}/hotstuff/exelayers.zip"
command3="cd /home/${username}/hotstuff && unzip exelayers.zip; sudo apt-get -y install cmake"
command4="sudo apt install libssl-dev; sudo apt install openssl; cd /home/${username}/hotstuff/exelayers-main && mkdir build && cd build/ && cmake .. && make && ./tests && sudo make install && cd ../example && gcc example.c -o out -lexelayers -lcrypto -lsecp256k1"
command5="sudo rm -rf /usr/local/go ; cd /home/${username}/hotstuff/ && wget https://go.dev/dl/go1.19.5.linux-amd64.tar.gz   && sudo tar -C /usr/local -xzf go1.19.5.linux-amd64.tar.gz && export PATH=\$PATH:/usr/local/go/bin && go version"
command6="export PATH=\$PATH:/usr/local/go/bin; cd /home/${username}/hotstuff/logs/temp/hotstuff ; go get -u github.com/golang/protobuf/protoc-gen-go; go get go.dedis.ch/kyber/v3"
command7="export PATH=\$PATH:/usr/local/go/bin; cd /home/${username}/hotstuff/logs/temp/hotstuff ; go get -u google.golang.org/grpc"
#command8="sudo printf \"\n export LD_LIBRARY_PATH=/usr/local/lib\" >>  ~/.bashrc && source ~/.bashrc"
command9="mv /home/${username}/hotstuff/logs/temp/hotstuff/ /home/${username}/hotstuff/"
command10="export PATH=\$PATH:/usr/local/go/bin ;  cd /home/${username}/hotstuff/hotstuff ; go mod tidy; go build -v -o ./client/bin/client ./client/"
command11="export PATH=\$PATH:/usr/local/go/bin ; cd /home/${username}/hotstuff/hotstuff ; go mod tidy; go build -v -o ./replica/bin/replica ./replica/"
command12="cd /home/${username}/hotstuff/hotstuff && mkdir logs"

for index in "${!replicas[@]}";
do
    echo "copying files to replica ${index}"
    sshpass ssh "${replicas[${index}]}" -i ${cert} "${reset_directory};${kill_instances}"
    scp -i ${cert} ${local_zip_path} "${replicas[${index}]}":${remote_home_path}
    sshpass ssh "${replicas[${index}]}" -i ${cert} "${command1}; ${command2}; ${command3};${command4}; ${command5}; ${command6}; ${command7}; ${command9}; ${command10}; ${command11}; ${command12}"
#   sshpass ssh "${replicas[${index}]}" -i ${cert} "sudo printf \"\n127.0.0.1  ${replica_names[${index}]}\" >> /etc/hosts"
done

rm ${local_zip_path}

echo "setup complete"