#!/bin/bash

if [[ $# != 1 ]]; then
    echo "usage $0 [dir]"
    exit -1
fi

dir=$1

for file in $dir/*
do
    echo $file
    fleetctl start $file
done
