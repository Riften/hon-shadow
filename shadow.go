package main

import (
	"context"
	//"crypto/rand"
	"fmt"

	"github.com/libp2p/go-libp2p"
	//"github.com/libp2p/go-libp2p-core/crypto"
)


// The context governs the lifetime of the libp2p node
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// To construct a simple host with all the default settings, just use `New`
	h, err := libp2p.New(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello World, my hosts ID is %s\n", h.ID())
}