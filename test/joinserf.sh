#!/bin/bash

if [[ $# != 2 ]]; then
    echo "usage: $0 [number_of_agents] [join_addr]"
    exit -1
fi

num=$1
joinaddr=$2
gopath=/root/gopher
interface=eth0

ipaddr=$(ifconfig $interface | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}')
echo $ipaddr

rootpath=$gopath/src/github.com/go-distributed/gog
logpath=$rootpath/test/serflog

for i in `seq 1 $num`;
do
    let agentport=8000+$i-1
    let restport=9000+$i-1

    agentaddr=$ipaddr:$agentport
    restaddr=$ipaddr:$restport

    $gopath/src/github.com/hashicorp/serf/bin/serf join -rpc-addr=$restaddr $joinaddr
done
