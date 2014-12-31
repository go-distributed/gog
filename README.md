gog
===

Gossip Over Gophers

[![GoDoc] (https://godoc.org/github.com/go-distributed/gog?status.png)](https://godoc.org/github.com/go-distributed/gog)

[Writeup] (https://drive.google.com/file/d/0BzkE2qUvDHtpTFdZRVFCWDI0dWM/view?usp=sharing)

###How to build:

```shell
$ go build
```

Show usage:

```shell
$ ./gog -h
```


Run a standalone node:

```shell
$ ./gog
```

To join an existing cluster:

1. Form a two-node-cluster

    ```shell
    $ ./gog -addr="localhost:8000" -rest-addr="localhost:8001"
    $ ./gog -addr="localhost:8002" -rest-addr="localhost:8003"
    ```

2. Let the first node join the second node

    ```shell
    $ curl http://localhost:8001/api/join -d peer=localhost:8002
    ```

3. Show the view in the first node

    ```shell
    $ curl http://localhost:8001/api/list
    {"active_view":[{"id":"localhost:8002","address":"localhost:8002"}],"passive_view":[]}
    ```

You can also provide a json file contains a list of nodes:
```shell
$ cat peers.json
["localhost:8000", "localhost:8001", "localhost:8002", "localhost:8003"]

$ curl -d @peers.json -H "Content-Type: application/json" http://localhost:8001/api/join
```

To broadcast a message:

```shell
$ curl http://localhost:8001/api/broadcast -d message=hello
```
