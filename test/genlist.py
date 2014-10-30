#!/usr/bin/python

# Generate peer list.
# usage: ./genlist.py [filename] [number_of_hosts]

import sys

startport = 8000

def main():
    if len(sys.argv) != 3:
        print "usage:", sys.argv[0], "[filename] [number_of_hosts]"
        return

    filename = sys.argv[1]
    n = int(sys.argv[2])
    f = open(filename, "w+")

    for i in range (0, int(n)):
        line = "localhost:%d\n" % (startport+i)
        f.write(line)
    f.close()

if __name__ == "__main__":
    main()
