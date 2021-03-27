package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const ChatRoomBufSize = 128

type ChatRoom struct {
	Messages chan *ChatMessage
	ctx      context.Context
	ps       *pubsub.PubSub
	topic    *pubsub.Topic
	sub      *pubsub.Subscription

	roomName string
	nick     string
	self     peer.ID
}

func (cr *ChatRoom) readLoop() {
	for {
		msg, err := cr.sub.Next(cr.ctx)
		if err != nil {
			close(cr.Messages)
			return
		}

		if msg.ReceivedFrom == cr.self {
			continue
		}

		cm := new(ChatMessage)
		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			continue
		}

		cr.Messages <- cm
	}
}

func (cr *ChatRoom) Publish(message string) error {
	m := &ChatMessage{
		Message:    message,
		SenderId:   cr.self.Pretty(),
		SenderNick: cr.nick,
	}

	msgBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return cr.topic.Publish(cr.ctx, msgBytes)
}

func (cr *ChatRoom) ListPeers() []peer.ID {
	return cr.ps.ListPeers(topicName(cr.roomName))
}

type ChatMessage struct {
	Message    string
	SenderId   string
	SenderNick string
}

func JoinChatRoom(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID, nick string, room string) (*ChatRoom, error) {
	topic, err := ps.Join(topicName(room))
	if err != nil {
		return nil, err
	}

	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	cr := &ChatRoom{
		ctx:      ctx,
		ps:       ps,
		topic:    topic,
		sub:      sub,
		self:     selfID,
		nick:     nick,
		roomName: room,
		Messages: make(chan *ChatMessage, ChatRoomBufSize),
	}

	go cr.readLoop()
	return cr, nil
}

func topicName(room string) string {
	return fmt.Sprintf("chat-room:%s", room)
}
