#!/usr/bin/python

import subprocess
import sys
import json
import random

num = 0
token = ""
specialname = "gog"
specialnameserf = "serf"

def main():
    if len(sys.argv) != 3:
        print "usage:", sys.argv[0], "[number of droplets] [token]"
        return -1

    global num
    num = int(sys.argv[1])

    global token
    token = sys.argv[2]

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
        if name == specialname or name == specialnameserf:
            continue
        print "destroying %d" % dropid
        subprocess.call(["curl", "-X", "DELETE", "-H", "Content-Type: application/json", "-H", "Authorization: Bearer "+token, "https://api.digitalocean.com/v2/droplets/"+str(dropid)])

if __name__ == '__main__':
    main()
