# kndchain

## Running A Node

### Setup

```
mkdir ~/kndchainDatadir
mkdir ~/kndchainKeys
git clone git@github.com:knd/kndchain.git
cd kndchain
```

### Run a node

```
go run cmd/production/main.go
```

### Mining

```
# to run a mining node (optional)
go run cmd/production/main.go mining
```
