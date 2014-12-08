#!/bin/bash

if [[ $# != 3 ]]; then
    echo "usage $0 [origin_file] [number_of_gog_per_machine] [dir]"
    exit -1
fi

file=$1
num=$2
service_dir=$3

for i in `seq 0 $num`
do
    let addr=8000+$i
    let rest=9000+$i

    sed -e s/gog8000/gog$addr/ -e s/8000/$addr/ -e s/9000/$rest/ $file > $service_dir/agent_${addr}.service
done
