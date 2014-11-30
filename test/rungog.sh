#!/bin/bash

ipaddr=$(ifconfig eth0 | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}')

rootpath=$GOPATH/src/github.com/go-distributed/gog
logpath=$rootpath/test/log

$rootpath/gog -addr=$ipaddr:8424 > $logpath/stdout 2>$logpath/stderr
