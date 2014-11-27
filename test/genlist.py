#!/usr/bin/python

# Generate peer list.
# usage: ./genlist.py [filename] [number_of_hosts]

import sys
import json

startport = 8000

def main():
    if len(sys.argv) != 3:
        print "usage:", sys.argv[0], "[filename] [number_of_hosts]"
        return

    filename = sys.argv[1]
    n = int(sys.argv[2])
    f = open(filename, "w+")

    peers = []
    for i in range (0, int(n)):
        peer = "localhost:%d" % (startport+i)
        peers.append(peer)

    f.write(json.dumps(peers))
    f.close()

if __name__ == "__main__":
    main()
