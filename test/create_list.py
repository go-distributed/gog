#!/usr/bin/python

import sys

startport = 8000

def main(filename, n):
    f = open(filename, "w+")
    for i in range(0, int(n)):
        line = "localhost:%d\n" % (startport+i)
        f.write(line)
    f.close()

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print "usage:", sys.argv[0], "[filename] [number_of_hosts]"
    else:
        main(sys.argv[1], sys.argv[2])
