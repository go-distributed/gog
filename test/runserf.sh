#!/bin/bash

if [[ $# != 1 ]]; then
    echo "usage: $0 [number_of_agents]"
    exit -1
fi

num=$1
gopath=/home/yifan/gopher
interface=lo

ipaddr=$(ifconfig $interface | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}')
echo $ipaddr

rootpath=$gopath/src/github.com/go-distributed/gog
logpath=$rootpath/test/serflog

userscript=$rootpath/test/query_handler.sh

for i in `seq 1 $num`;
do
    let agentport=8000+$i-1
    let restport=9000+$i-1

    agentaddr=$ipaddr:$agentport
    restaddr=$ipaddr:$restport

    mkdir $logpath/$agentaddr

    $gopath/src/github.com/hashicorp/serf/bin/serf agent -node=$agentaddr -bind=$agentaddr -rpc-addr=$restaddr -event-handler user=$userscript > $logpath/$agentaddr/stdout 2>$logpath/$agentaddr/stderr &
done
