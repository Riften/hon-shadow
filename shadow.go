package main

import (
	//"crypto/rand"
	"context"
	"fmt"
	host "github.com/Riften/hon-shadow/host"
	"github.com/Riften/hon-shadow/service"
	//"github.com/libp2p/go-libp2p-core/crypto"
)


// The context governs the lifetime of the libp2p node
func main() {
	ctx := context.Background()
	h, err := host.NewHost(ctx)
	service.Test()
	fmt.Printf("Hello World, my hosts ID is %s\n", h.ID())
}