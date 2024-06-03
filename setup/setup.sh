pwd=$(pwd)
. "${pwd}"/setup/ip.sh

rm dummy/bin/dummy
rm torture/bin/torture
echo "Removed old binaries"

/bin/bash build.sh
echo "built binaries"

mkdir -p temp

cp -r dummy/                     temp
cp -r torture/                   temp

zip -r temp.zip temp/
rm -r temp

kill_instances="pkill dummy ; pkill torture; pkill dummy ; pkill torture; pkill dummy ; pkill torture; pkill dummy ; pkill torture; pkill dummy ; pkill torture; pkill dummy ; pkill torture"

local_zip_path="temp.zip"
remote_home_path="/home/${username}/torture/"
reset_directory="rm -rf /home/${username}/torture; mkdir -p /home/${username}/torture/ "

command="sudo apt-get update; sudo apt-get install unzip;sudo apt-get install zip; cd /home/${username}/torture && unzip temp.zip"

for index in "${!replicas[@]}";
do
    echo "copying files to replica ${index}"
    sshpass ssh "${replicas[${index}]}" -i ${cert} "${reset_directory};${kill_instances}"
    scp -i ${cert} ${local_zip_path} "${replicas[${index}]}":${remote_home_path}
    sshpass ssh "${replicas[${index}]}" -i ${cert} "${command}"
done

rm ${local_zip_path}

echo "setup complete"