package test

import (
	"context"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewPubSubStream(srv *pstest.Server, projectID, topicID, subscriptionID string) (*pubsub.Topic, *pubsub.Subscription, context.CancelFunc) {

	ctx, cancel := context.WithCancel(context.Background())
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

	cancelFunc := func() {
		cancel()
		defer conn.Close()
		defer client.Close()
		defer topic.Delete(ctx)
	}
	return topic, sub, cancelFunc
}
