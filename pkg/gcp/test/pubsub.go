package test

import (
	"context"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PubSubStream struct {
	Topic        *pubsub.Topic
	Subscription *pubsub.Subscription
	client       *pubsub.Client
	conn         *grpc.ClientConn
}

func NewPubSubStreamWithContext(ctx context.Context, srv *pstest.Server, projectID, topicID, subscriptionID string) *PubSubStream {

	conn, err := grpc.Dial(srv.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client, err := pubsub.NewClient(ctx, projectID, option.WithGRPCConn(conn))
	if err != nil {
		panic(err)
	}
	topic, err := client.CreateTopic(ctx, topicID)
	if err != nil {
		panic(err)
	}
	subs, err := client.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		panic(err)
	}

	sub := client.Subscription(subs.ID())

	sub.ReceiveSettings.Synchronous = false
	sub.ReceiveSettings.NumGoroutines = -1
	sub.ReceiveSettings.MaxOutstandingMessages = -1
	sub.ReceiveSettings.MaxOutstandingBytes = -1

	return &PubSubStream{
		Subscription: sub,
		Topic:        topic,
		client:       client,
		conn:         conn,
	}
}

func (s *PubSubStream) Close() error {
	s.Topic.Delete(context.Background())
	s.Subscription.Delete(context.Background())
	s.client.Close()
	s.conn.Close()
	return nil
}
