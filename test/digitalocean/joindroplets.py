#!/usr/bin/python

import subprocess
import sys
import json

token = ""
specialname = "gogs"
num = 40

def main():
    if len(sys.argv) != 3:
        print "usage: %s [token] [num_per_nodes]" %sys.argv[0]
        return -1

    global token
    global num
    token = sys.argv[1]
    num = int(sys.argv[2])

    nullf = open("/dev/null", "w")

    # List droplets
    p = subprocess.Popen(["curl", "-H", "Content-Type: application/json", "-H", "Authorization: Bearer "+token, "https://api.digitalocean.com/v2/droplets?page=1&per_page=1024"], stdout=subprocess.PIPE, stderr=nullf)
    droplets = p.stdout.read()
    p.wait()
    nullf.close()

    jdroplets =  json.loads(droplets)["droplets"]

    for i in range(0, len(jdroplets)):
        dropid = jdroplets[i]["id"]
        name = jdroplets[i]["name"]
        dropip = jdroplets[i]["networks"]["v4"][0]["ip_address"]
        if name == specialname:
            continue
        for j in range(0, num):
            restaddr = dropip+":"+str(9000+j)
            print "joining %d, addr %s" % (dropid, restaddr)
            subprocess.call(["curl", "http://"+restaddr+"/api/join", "-d", "peer=104.236.9.169:8000"])

if __name__ == '__main__':
    main()
