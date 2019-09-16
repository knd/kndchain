package main

import (
	"fmt"
	"os"

	"github.com/knd/kndchain/pkg/crypto"
	"github.com/knd/kndchain/pkg/wallet"
)

func main() {
	secp265k1 := crypto.NewSecp256k1Generator()
	pathToKeyDir := os.Args[1]
	wal := wallet.NewWallet(secp265k1, nil, 0, &pathToKeyDir)

	fmt.Printf("Your pubKey: %s\n", wal.PubKeyHex())
}
