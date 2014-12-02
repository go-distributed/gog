#!/usr/bin/python

import subprocess
import sys

num = 0
token = ""

def main():
    if len(sys.argv) != 3:
        print "usage:", sys.argv[0], "[number of droplets] [token]"
        return -1

    global num
    num = int(sys.argv[1])
    global token
    token = sys.argv[2]

    f = open("create_droplets.json")
    json_data = f.read()
    f.close()

    for i in range(0, num):
        subprocess.call(["curl", "-X", "POST", "-H", "Content-Type: application/json", "-H", "Authorization: Bearer "+token, "-d", json_data, "https://api.digitalocean.com/v2/droplets"])

if __name__ == '__main__':
    main()
