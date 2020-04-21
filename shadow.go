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
	shadowCtx, _ := context.WithCancel(ctx)
	h, err := host.NewHost(ctx)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	shadowService := service.NewShadowService(ctx, h)
	shadowService.Start()
	fmt.Printf("Hello World, my hosts ID is %s\n", h.ID())
	select{
		case <-shadowCtx.Done():
			fmt.Printf("Routine end.")
			break
	}

}