# kndchain

## Running A Node

### Setup

```
mkdir ~/kndchainDatadir
git clone git@github.com:knd/kndchain.git
cd kndchain
```

### Run a node

```
go run cmd/server/main.go
```

### Mining

```
# to run a mining node (optional)
go run cmd/server/maing.go --mining
```
