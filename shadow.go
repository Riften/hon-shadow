package main

import (
	//"crypto/rand"
	"context"
	"fmt"
	host "github.com/Riften/hon-shadow/host"
	//"github.com/libp2p/go-libp2p-core/crypto"
)

type shadowPeer struct {

}


// The context governs the lifetime of the libp2p node
func main() {
	ctx := context.Background()
	h, err := host.NewHost(ctx);  if err != nil {fmt.Print(err.Error())}
	
	fmt.Printf("Hello World, my hosts ID is %s\n", h.ID())
}