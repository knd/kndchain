# kndchain

## Setup

- [Golang](https://golang.org/dl/) (>= 1.12)
- [Redis](https://redis.io/topics/quickstart) (stable)

## Quick start

```
# Run redis
$ redis-server -daemonize yes
$ redis-cli ping
PONG

$ git clone https://github.com/knd/kndchain.git
$ cd cmd/miner
$ go build main.go miner.go

# Run node with mining
$ ./main -mining=true

# Run node w/o mining
$ ./main

# Clean up
rm -rf /tmp/kndchainKeys/*
rm -rf /tmp/kndchainDatadir/*
```

## For help

```
$ ./main -h
(ツ) ./main -h
Usage of ./main:
  -address string
    	provide pubkeyhex/ address used for transactions or mining reward
  -chainDatadir string
    	directory to store blockchain data (default "/tmp/kndchainDatadir")
  -keysDatadir string
    	directory to store keys (default "/tmp/kndchainKeys")
  -mining
    	enable mining option
```

## Simulate 2 miners (with the former acting as beacon node)

### Terminal 1

```
$ cd cmd/miner
$ go build main.go miner.go
$ ./main -mining=true
```

### Terminal 2

```
$ mkdir /tmp/anotherminerDatadir
$ mkdir /tmp/anotherminerKeys
$ cd cmd/anotherminer
$ go build main.go miner.go
$ ./main -chainDatadir=/tmp/anotherminerDatadir -keysDatadir=/tmp/anotherminerKeys -beaconURL=http://localhost:3001 -mining=true
```
