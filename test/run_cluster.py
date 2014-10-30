#!/usr/bin/python

# start a gossip cluster.
# usage: ./run_cluster.py [number_of_hosts]

import subprocess
import sys
import time

gogpath = "../gog"
filename = "peerlist.txt"
genlist = "./genlist.py"

process = []

def startNode(hostport):
    p = subprocess.Popen([gogpath, "-addr", hostport, "-peer-file", filename], stdin=subprocess.PIPE)
    process.append(p)

def startNodes():
    f = open(filename)
    lines = f.readlines()

    for line in lines:
        startNode(line[:-1])
        time.sleep(2)

def generateList(n):
    subprocess.call([genlist, filename, str(n)])

def main():
    if len(sys.argv) != 3:
        print "usage:", sys.argv[0], "[gogpath]", "[number_of_hosts]"
        return

    gogpath = sys.argv[1]
    n = int(sys.argv[2])

    # generate host list
    generateList(n)

    # run node
    startNodes()

    while True:
        continue


if __name__ == "__main__":
    main()
