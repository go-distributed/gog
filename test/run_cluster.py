#!/usr/bin/python

# start a gossip cluster.
# usage: ./run_cluster.py [number_of_hosts]

import os
import subprocess
import sys
import time
import random

timer_addr = "localhost:11000"
startport = 8000
n = 10
gogpath = "../gog"
peerfile = "peers.json"
genlist = "./genlist.py"
logdir = "./log"
userscript = "./query_handler.sh"

process = []
addrs = []
rest_addrs = []

def startNode(addr, rest_addr):
    stdoutdir = logdir + "/stdout"
    stderrdir = logdir + "/stderr"

    if not os.path.exists(stdoutdir):
        os.makedirs(stdoutdir)
    if not os.path.exists(stderrdir):
        os.makedirs(stderrdir)

    stdoutpath = stdoutdir + "/" + addr
    stderrpath = stderrdir + "/" + addr

    stdoutf = open(stdoutpath, "w+")
    stderrf = open(stderrpath, "w+")

    p = subprocess.Popen([gogpath, "-addr", addr, "-rest-addr", rest_addr, "-user-message-handler", userscript, "-peer-file", peerfile],
                         stdin=subprocess.PIPE, stdout=stdoutf, stderr=stderrf)
    process.append(p)

def startNodes():
    print "Launching %d nodes..." %n

    for i in range(0, n):
        addr = "localhost:%d" % (startport+i)
        rest_addr = "localhost:%d" % (startport+i+n)

        global addrs
        global rest_addrs

        addrs.append(addr)
        rest_addrs.append(rest_addr)

        print "Launching node %d..." %i
        startNode(addr, rest_addr)

    print "All nodes launched.\n"

def joinNodes():
    print "Joining %d nodes..." %n

    for i in range(0, len(rest_addrs)):
        rest_addr = rest_addrs[i]
        #subprocess.call(["curl", "-d", "@"+peerfile, "-H", "Content-Type: application/json", "http://"+rest_addr+"/api/join"])
        subprocess.call(["curl", "http://"+rest_addr+"/api/join", "-d", "peer=localhost:8000"])

def listViews():
    nullf = open("/dev/null", "w")
    for i in range(0, len(rest_addrs)):
        rest_addr = rest_addrs[i]
        subprocess.call(["curl", "http://"+rest_addr+"/api/list"], stdout=nullf, stderr=nullf)

def main():
    if len(sys.argv) != 4:
        print "usage:", sys.argv[0], "[gogpath]", "[number_of_hosts]", "[number_of_kill]"
        return

    global gopath
    global n

    gogpath = sys.argv[1]
    n = int(sys.argv[2])
    m = int(sys.argv[3])

    subprocess.call([genlist, peerfile, str(n)])

    # run nodes
    startNodes()

    time.sleep(1)

    # connecting nodes
    joinNodes()

    for i in range(0, 5):
        print "Listing views %d..." %i
        time.sleep(1)
        listViews()

    print "Broadcasting message"
    # start timer
    subprocess.call(["curl", "http://"+timer_addr+"/start"])

    # send message
    subprocess.call(["curl", "http://"+rest_addrs[0]+"/api/broadcast", "-d", "message=hello"])

    while True:
        continue

    #print "listing..."
    #i = 0
    #for p in process:
    #    print "list %d" %i
    #    i = i + 1
    #    p.stdin.write("list\n")
    #
    #time.sleep(1)
    #print "sending..."
    #process[0].stdin.write("hello\n")
    #
    #time.sleep(5)
    #
    #print "randomly killing %d nodes..." %m
    #for i in range (0, m):
    #    index = random.randint(0, n-1)
    #    p = process[index]
    #    while p is None:
    #        index = (index + 1) % n
    #        p = process[index]
    #    process[index].kill()
    #    process[index] = None
    #    
    #time.sleep(5)
    #i = 0
    #for p in process:
    #    if not p is None:
    #        print "list %d" %i
    #        p.stdin.write("list\n")
    #    i = i + 1
    #
    #for p in process:
    #    if not p is None:
    #        print "sending..."
    #        p.stdin.write("world\n")
    #        break
    #
    #time.sleep(10)

if __name__ == "__main__":
    main()
