#!/usr/bin/python

import glob
import sys

logdir = "./log"

def main():
    if len(sys.argv) != 2:
        print "usage:", sys.argv[0], "[log_dir]"
        return

    global logdir

    logdir = sys.argv[1]
    stdoutdir = logdir + "/stdout"
    stderrdir = logdir + "/stderr"

    files = glob.glob(stdoutdir + "/*")

    print "hello"

    # print those isolated nodes.
    for name in files:
        f = open(name)
        line = f.read()
        f.close()
        if not "hello" in line:
            print name

    #print "world"
    #
    #for name in files:
    #    f = open(name)
    #    line = f.read()
    #    f.close()
    #    if not "hello" in line:
    #        print name

if __name__ == "__main__":
    main()
