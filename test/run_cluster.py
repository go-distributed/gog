#!/usr/bin/python

# start a gossip cluster.
# usage: ./run_cluster.py [number_of_hosts]

import subprocess
import sys
import time

n = 10
gogpath = "../gog"
filename = "peerlist.txt"
genlist = "./genlist.py"
logdir = "./log"

process = []

def startNode(hostport):
    logpath = logdir + "/" + hostport
    logf = open(logpath, "w+")
    p = subprocess.Popen([gogpath, "-addr", hostport, "-peer-file", filename],
                         stdin=subprocess.PIPE, stdout=logf, stderr=logf)
    process.append(p)

def startNodes():
    print "Launching %d nodes..." %n
    f = open(filename)
    lines = f.readlines()

    i = 0
    for line in lines:
        print "Launching node %d..." %i
        startNode(line[:-1])
        i = i + 1

    print "All nodes launched.\n"

def generateList(n):
    print "Generating peer list file..."
    subprocess.call([genlist, filename, str(n)])
    print "Ok.\n"

def main():
    if len(sys.argv) != 4:
        print "usage:", sys.argv[0], "[gogpath]", "[number_of_hosts]", "[log_dir]"
        return

    global gopath
    global n
    global logdir

    gogpath = sys.argv[1]
    n = int(sys.argv[2])
    logdir = sys.argv[3]

    # generate host list
    generateList(n)

    # run node
    startNodes()

    for p in process:
        p.stdin.write("list\n")

    time.sleep(1)
    process[0].stdin.write("hello world\n")

    while True:
        continue

if __name__ == "__main__":
    main()
