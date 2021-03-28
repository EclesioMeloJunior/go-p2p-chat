package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery"
)

const DiscoveryInterval = time.Hour
const DiscoveryServiceTag = "pubsub-chat-example"

func main() {
	nicknameFlag := flag.String("nickname", "", "nickname to use in chat. Will be generated if empty")
	roomFlag := flag.String("room", "default-chat-room", "name of the chat to join")
	flag.Parse()

	ctx := context.Background()

	h, err := libp2p.New(ctx, libp2p.DisableRelay())
	if err != nil {
		panic(err)
	}

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}

	_, addr, err := newDHTDiscovery(ctx, h)
	if err != nil {
		panic(err)
	}

	err = setupDiscovery(ctx, h)
	if err != nil {
		panic(err)
	}

	nick := *nicknameFlag
	if len(nick) == 0 {
		nick = defaulNick(h.ID())
	}

	room := *roomFlag

	cr, err := JoinChatRoom(ctx, ps, h.ID(), nick, room)
	if err != nil {
		panic(err)
	}

	ui := NewChatUI(cr)

	if err = ui.Run(addr); err != nil {
		printErr("error runing chat ui: %s", err)
	}
}

func defaulNick(p peer.ID) string {
	return fmt.Sprintf("%s-%s", os.Getenv("USER"), shortPeerId(p))
}

func shortPeerId(p peer.ID) string {
	return p.Pretty()[len(p.Pretty())-8:]
}

type discoveryNotifee struct {
	h host.Host
}

func (d *discoveryNotifee) HandlePeerFound(pf peer.AddrInfo) {
	fmt.Printf("Discovered a new peer %s\n", pf.ID.Pretty())
	err := d.h.Connect(context.Background(), pf)

	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pf.ID.Pretty(), err)
	}
}

func setupDiscovery(ctx context.Context, h host.Host) error {
	disc, err := discovery.NewMdnsService(ctx, h, DiscoveryInterval, DiscoveryServiceTag)

	if err != nil {
		return err
	}

	d := discoveryNotifee{
		h: h,
	}

	disc.RegisterNotifee(&d)
	return nil
}

func printErr(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
}
