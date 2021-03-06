package main

import (
	"crypto/x509"
	"log"
	"context"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func connexionPublisher(ctx context.Context, tracer opentracing.Tracer, address string, filename string, scope ...string) pubsub.PublisherClient {
	var err error

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Println(err)
	}

	creds := credentials.NewClientTLSFromCert(pool, "")

	perRPC, err := oauth.NewServiceAccountFromFile(filename, scope...)
	if err != nil {
		log.Println(err)
	}

	conn, err := grpc.Dial(
		"pubsub.googleapis.com:443",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)),
	)
	if err != nil {
		log.Println(err)
	}

	return pubsub.NewPublisherClient(conn)
}


func (s server) publishmessage(mess *pubsub.PubsubMessage) {

	var request pubsub.PublishRequest
	request.Topic = *publictopic
	request.Messages = append(request.Messages, mess)
	ctx := context.Background()

	if _, err := s.pub.Publish(ctx, &request); err != nil {
		log.Println(err)
	}
}
