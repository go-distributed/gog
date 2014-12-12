#!/usr/bin/python

import subprocess
import sys
import json

token = ""
specialname = "gog"
serfname = "serf"
num = 50
nodeaddr = ""

def main():
    if len(sys.argv) != 4:
        print "usage: %s [token] [num_per_nodes] [node_addr]" %sys.argv[0]
        return -1

    global token
    global num
    global nodeaddr

    token = sys.argv[1]
    num = int(sys.argv[2])
    nodeaddr = sys.argv[3]

    nullf = open("/dev/null", "w")

    # List droplets
    p = subprocess.Popen(["curl", "-H", "Content-Type: application/json", "-H", "Authorization: Bearer "+token, "https://api.digitalocean.com/v2/droplets?page=1&per_page=1024"], stdout=subprocess.PIPE, stderr=nullf)
    droplets = p.stdout.read()
    p.wait()
    nullf.close()

    jdroplets =  json.loads(droplets)["droplets"]

    total = 0
    for i in range(0, len(jdroplets)):
        dropid = jdroplets[i]["id"]
        name = jdroplets[i]["name"]
        dropip = jdroplets[i]["networks"]["v4"][0]["ip_address"]
        if name == specialname or name == serfname:
            continue
        for j in range(0, num):
            restaddr = dropip+":"+str(9000+j)
            print "joining %d, addr %s" % (total, restaddr)
            total = total + 1
            subprocess.call(["curl", "http://"+restaddr+"/api/join", "-d", "peer="+nodeaddr+":8000"])

if __name__ == '__main__':
    main()
