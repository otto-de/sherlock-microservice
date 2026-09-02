package test

import (
	"context"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"cloud.google.com/go/pubsub/v2/pstest"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PubSubStream struct {
	Publisher  *pubsub.Publisher
	Subscriber *pubsub.Subscriber
	client     *pubsub.Client
	conn       *grpc.ClientConn
}

func NewPubSubStreamWithContext(ctx context.Context, srv *pstest.Server, projectID, topicID, subscriptionID string) *PubSubStream {

	conn, err := grpc.NewClient(srv.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client, err := pubsub.NewClient(ctx, projectID, option.WithGRPCConn(conn))
	if err != nil {
		panic(err)
	}
	topic, err := client.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{
		Name: "projects/" + projectID + "/topics/" + topicID,
	})
	if err != nil {
		panic(err)
	}
	_, err = client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name:  "projects/" + projectID + "/subscriptions/" + subscriptionID,
		Topic: topic.Name,
	})
	if err != nil {
		panic(err)
	}

	sub := client.Subscriber("projects/" + projectID + "/subscriptions/" + subscriptionID)

	sub.ReceiveSettings.NumGoroutines = -1
	sub.ReceiveSettings.MaxOutstandingMessages = -1
	sub.ReceiveSettings.MaxOutstandingBytes = -1

	return &PubSubStream{
		Subscriber: sub,
		Publisher:  client.Publisher("projects/" + projectID + "/topics/" + topicID),
		client:     client,
		conn:       conn,
	}
}

func (s *PubSubStream) Close() error {
	s.Publisher.Stop()
	s.client.SubscriptionAdminClient.DeleteSubscription(context.Background(), &pubsubpb.DeleteSubscriptionRequest{
		Subscription: s.Subscriber.String(),
	})
	s.client.TopicAdminClient.DeleteTopic(context.Background(), &pubsubpb.DeleteTopicRequest{
		Topic: s.Publisher.String(),
	})
	s.client.Close()
	s.conn.Close()
	return nil
}
