#!/usr/bin/python

# start a gossip cluster.
# usage: ./run_cluster.py [number_of_hosts]

import os
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

    p = subprocess.Popen([gogpath, "-addr", hostport, "-peer-file", filename],
                         stdin=subprocess.PIPE, stdout=stdoutf, stderr=stderrf)
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

    time.sleep(1)

    print "listing..."
    i = 0
    for p in process:
        print "list %d" %i
        i = i + 1
        p.stdin.write("list\n")

    time.sleep(1)
    print "sending..."
    process[0].stdin.write("hello\n")

    time.sleep(10)
    i = 0
    for p in process:
        print "list %d" %i
        i = i + 1
        p.stdin.write("list\n")

    print "sending..."
    process[0].stdin.write("world\n")
    while True:
        continue

if __name__ == "__main__":
    main()
