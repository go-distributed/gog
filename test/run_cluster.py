#!/usr/bin/python

# start a gossip cluster.
# usage: ./run_cluster.py [number_of_hosts]

import os
import subprocess
import sys
import time
import random

n = 10
gogpath = "../gog"
filename = "peerlist.txt"
genlist = "./genlist.py"
logdir = "./log"
delay = "0"
droprate = "10"


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

    p = subprocess.Popen([gogpath, "-addr", hostport, "-peer-file", filename,
                          "-delay", delay, "-droprate", droprate],
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
    if len(sys.argv) != 7:
        print "usage:", sys.argv[0], "[gogpath]", "[number_of_hosts]", "[number_of_kill]", "[log_dir]", "[delay]", "[droprate]"
        return

    global gopath
    global n
    global logdir
    global delay
    global droprate

    gogpath = sys.argv[1]
    n = int(sys.argv[2])
    m = int(sys.argv[3])
    logdir = sys.argv[4]
    delay = sys.argv[5]
    droprate = sys.argv[6]

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

    time.sleep(5)

    print "randomly killing %d nodes..." %m
    for i in range (0, m):
        index = random.randint(0, n-1)
        p = process[index]
        while p is None:
            index = (index + 1) % n
            p = process[index]
        process[index].kill()
        process[index] = None
        
    time.sleep(5)
    i = 0
    for p in process:
        if not p is None:
            print "list %d" %i
            p.stdin.write("list\n")
        i = i + 1

    for p in process:
        if not p is None:
            print "sending..."
            p.stdin.write("world\n")
            break

    time.sleep(10)

if __name__ == "__main__":
    main()
