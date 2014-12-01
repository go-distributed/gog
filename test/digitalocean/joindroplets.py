#!/usr/bin/python

import subprocess
import sys
import json

token = ""
specialname = "gogs"

def main():
    if len(sys.argv) != 2:
        print "usage: %s [token]" %sys.argv[0]
        return -1

    global token
    token = sys.argv[1]

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
        print "joining %d, ip %s" % (dropid, dropip)
        subprocess.call(["curl", "http://"+dropip+":8425/api/join", "-d", "peer=104.236.9.169:8424"])

if __name__ == '__main__':
    main()
