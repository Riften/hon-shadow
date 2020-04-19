package shadow
import p2phost "github.com/libp2p/go-libp2p-core/host"


// Define shadow peer
type shadowPeer interface{
	Host() p2phost.Host
}
