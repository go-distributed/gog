#!/bin/bash

ipaddr=$(ifconfig eth0 | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}')

rootpath=/root/gopher/src/github.com/go-distributed/gog
logpath=$rootpath/test/log

userscript=$rootpath/test/query_handler.sh

$rootpath/gog -addr=$ipaddr:8424  -user-message-handler=$userscript > $logpath/stdout 2>$logpath/stderr
