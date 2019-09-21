# kndchain

## Quick start

```
git clone https://github.com/knd/kndchain.git
cd cmd/miner
go build main.go miner.go

# Run node with mining
./main -mining=true

# Run node w/o mining
./main
```

## For help

```
$ ./main -h
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
