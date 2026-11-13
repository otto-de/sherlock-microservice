package test

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"cloud.google.com/go/pubsub/v2/pstest"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PubSubStream struct {
	Topic        *pubsub.Publisher
	Subscription *pubsub.Subscriber
	topicName    string
	subName      string
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
	topicName := fmt.Sprintf("projects/%s/topics/%s", projectID, topicID)
	subName := fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subscriptionID)

	_, err = client.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{Name: topicName})
	if err != nil {
		panic(err)
	}
	_, err = client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name:  subName,
		Topic: topicName,
	})
	if err != nil {
		panic(err)
	}

	sub := client.Subscriber(subName)

	sub.ReceiveSettings.NumGoroutines = -1
	sub.ReceiveSettings.MaxOutstandingMessages = -1
	sub.ReceiveSettings.MaxOutstandingBytes = -1

	return &PubSubStream{
		Subscription: sub,
		Topic:        client.Publisher(topicName),
		topicName:    topicName,
		subName:      subName,
		client:       client,
		conn:         conn,
	}
}

func (s *PubSubStream) Close() error {
	s.Topic.Stop()
	s.client.TopicAdminClient.DeleteTopic(context.Background(), &pubsubpb.DeleteTopicRequest{Topic: s.topicName})
	s.client.SubscriptionAdminClient.DeleteSubscription(context.Background(), &pubsubpb.DeleteSubscriptionRequest{Subscription: s.subName})
	s.client.Close()
	s.conn.Close()
	return nil
}
