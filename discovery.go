package main

import (
	"context"
	"fmt"

	"github.com/ipfs/go-ipns"
	"github.com/libp2p/go-libp2p-core/host"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

var (
	bootstrapPeerAddrs []multiaddr.Multiaddr
	routingDiscovery   *discovery.RoutingDiscovery
)

var bootstrapPeers = []string{
	//"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	//"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
	//"/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",
	//"/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",
	//"/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd",
	"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",            // mars.i.ipfs.io
	"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",           // pluto.i.ipfs.io
	"/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",           // saturn.i.ipfs.io
	"/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",             // venus.i.ipfs.io
	"/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd",            // earth.i.ipfs.io
	"/ip6/2604:a880:1:20::203:d001/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",  // pluto.i.ipfs.io
	"/ip6/2400:6180:0:d0::151:6001/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",  // saturn.i.ipfs.io
	"/ip6/2604:a880:800:10::4a:5001/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64", // venus.i.ipfs.io
	"/ip6/2a03:b0c0:0:1010::23:1001/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd", // earth.i.ipfs.io

}

type dhtDiscovery struct {
	h   host.Host
	dht *dht.IpfsDHT
}

func newDHTDiscovery(ctx context.Context, h host.Host) (*dhtDiscovery, string, error) {
	kad, err := dht.New(
		ctx, h, dht.NamespacedValidator("ipns", ipns.Validator{}))

	if err != nil {
		return nil, "", err
	}

	addr := fmt.Sprintf("Host MultiAddress: %s/ipfs/%s (%s)\n", h.Addrs()[0].String(), h.ID().Pretty(), h.ID().String())

	disc := &dhtDiscovery{
		h:   h,
		dht: kad,
	}

	return disc, addr, nil
}

type NullValidator struct{}

func (nv NullValidator) Select(key string, values [][]byte) (int, error) {
	for _, b := range values {
		fmt.Printf("NullValidator select: %s - %s", key, string(b))
	}
	return 0, nil
}

func (nv NullValidator) Validate(key string, value []byte) error {
	fmt.Printf("NullValidator validate: %s - %s\n", key, string(value))
	return nil
}

func stringsToAddr(addrs []string) ([]multiaddr.Multiaddr, error) {
	maddrs := make([]multiaddr.Multiaddr, 0)
	for _, addrString := range addrs {
		addr, err := multiaddr.NewMultiaddr(addrString)
		if err != nil {
			return nil, err
		}

		maddrs = append(maddrs, addr)
	}

	return maddrs, nil
}
