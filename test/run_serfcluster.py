#!/usr/bin/python

# start a gossip cluster.
# usage: ./run_cluster.py [number_of_hosts]

import os
import subprocess
import sys
import time

n = 10
filename = "peerlist.txt"
genlist = "./genlist.py"
logdir = "./log"

process = []

def joinNodes():
    print "Joining %d nodes..." %n

    f = open(filename)
    lines = f.readlines()
    f.close()

    i = 0

    first_rpcaddr = lines[1][:-1]
    while i < len(lines):
        print "Joining node %d..." %(i/2)
        hostport = lines[i][:-1]
        subprocess.call(["serf", "join", "-rpc-addr", first_rpcaddr, hostport], stdout=None)
        i = i + 2

    print "All nodes joined.\n"
    

def startNode(hostport, rpcaddr):
    stdoutdir = logdir + "/stdout"
    stderrdir = logdir + "/stderr"

    if not os.path.exists(stdoutdir):
        os.makedirs(stdoutdir)
    if not os.path.exists(stderrdir):
        os.makedirs(stderrdir)

    stdoutpath = stdoutdir + "/" + hostport
    stderrpath = stderrdir + "/" + hostport

    stdoutf = open(stdoutpath, "w+")
    stderrf = open(stderrpath, "w+")

    p = subprocess.Popen(["serf", "agent", "-node", hostport, "-bind", hostport, "-rpc-addr", rpcaddr],
                         stdin=subprocess.PIPE, stdout=stdoutf, stderr=stderrf)
    process.append(p)


def startNodes():
    print "Launching %d nodes..." %n
    f = open(filename)
    lines = f.readlines()
    f.close()

    i = 0
    while i < len(lines):
        print "Launching node %d..." %(i/2)
        startNode(lines[i][:-1], lines[i+1][:-1])
        i = i + 2

    print "All nodes launched.\n"

def generateList(n):
    print "Generating peer list file..."
    subprocess.call([genlist, filename, str(2*n)])
    print "Ok.\n"

def main():
    if len(sys.argv) != 3:
        print "usage:", sys.argv[0], "[number_of_hosts]", "[log_dir]"
        return

    global gopath
    global n
    global logdir

    n = int(sys.argv[1])
    logdir = sys.argv[2]

    # generate host list
    generateList(n)

    # run node
    startNodes()

    time.sleep(1)

    joinNodes()

    # list member
    subprocess.call(["serf", "members", "-rpc-addr", "localhost:8001"])

    # send a message
    subprocess.call(["serf", "event", "-rpc-addr", "localhost:8001", "hello", "hello"])
    
    while True:
        continue

if __name__ == "__main__":
    main()
