#!/bin/bash

if [[ $# != 1 ]]; then
    echo "usage: $0 [number_of_agents]"
    exit -1
fi

num=$1
gopath=/root/gopher
interface=eth0

ipaddr=$(ifconfig $interface | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}')

rootpath=$gopath/src/github.com/go-distributed/gog
logpath=$rootpath/test/log

userscript=$rootpath/test/query_handler.sh

for i in `seq 1 $num`;
do
    let agentport=8000+$i-1
    let restport=9000+$i-1

    agentaddr=$ipaddr:$agentport
    restaddr=$ipaddr:$restport

    mkdir $logpath/$agentaddr

    $rootpath/gog -addr=$agentaddr -rest-addr=$restaddr -user-message-handler=$userscript > $logpath/$agentaddr/stdout 2>$logpath/$agentaddr/stderr &
done
