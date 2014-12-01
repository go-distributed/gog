#!/usr/bin/python

import subprocess
import sys

num = 0
token = "2f016dd850e30bb47ba912d1cdb1779d4cc7a7c7a7e1a8ac8cd16ab867714e1f"

def main():
    if len(sys.argv) != 2:
        print "usage:", sys.argv[0], "[number of droplets]"
        return -1

    global num
    num = int(sys.argv[1])

    f = open("create_droplets.json")
    json_data = f.read()
    f.close()

    for i in range(0, num):
        subprocess.call(["curl", "-X", "POST", "-H", "Content-Type: application/json", "-H", "Authorization: Bearer "+token, "-d", json_data, "https://api.digitalocean.com/v2/droplets"])

if __name__ == '__main__':
    main()
