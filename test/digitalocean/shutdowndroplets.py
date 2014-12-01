#!/usr/bin/python

import subprocess
import sys
import json
import random

num = 0
token = "2f016dd850e30bb47ba912d1cdb1779d4cc7a7c7a7e1a8ac8cd16ab867714e1f"
specialname = "gogs"

def main():
    if len(sys.argv) != 2:
        print "usage:", sys.argv[0], "[number of droplets]"
        return -1

    global num
    num = int(sys.argv[1])

    nullf = open("/dev/null", "w")

    # List droplets
    p = subprocess.Popen(["curl", "-H", "Content-Type: application/json", "-H", "Authorization: Bearer "+token, "https://api.digitalocean.com/v2/droplets?page=1&per_page=1024"], stdout=subprocess.PIPE, stderr=nullf)
    droplets = p.stdout.read()
    p.wait()
    nullf.close()

    jdroplets =  json.loads(droplets)["droplets"]

    # Shuffle droplets
    random.shuffle(jdroplets)

    if num >= len(jdroplets):
        num = len(jdroplets)

    for i in range(0, num):
        dropid = jdroplets[i]["id"]
        name = jdroplets[i]["name"]
        if name == specialname:
            continue
        print "shutting down %d" % dropid
        subprocess.call(["curl", "-X", "POST", "-H", "Content-Type: application/json", "-H", "Authorization: Bearer "+token, "-d",  "{\"type\":\"power_off\"}", "https://api.digitalocean.com/v2/droplets/"+str(dropid)+"/actions"])

if __name__ == '__main__':
    main()
