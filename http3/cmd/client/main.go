package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
	greetv1 "github.com/oncilla/grpc-experiments/http3/gen/greet/v1"
	"github.com/oncilla/grpc-experiments/http3/gen/greet/v1/greetv1connect"

	"github.com/bufbuild/connect-go"
)

func main() {
	client := greetv1connect.NewGreetServiceClient(
		&http.Client{
			Transport: &http3.RoundTripper{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		"https://localhost:8080",
		connect.WithGRPC(),
	)
	res, err := client.Greet(
		context.Background(),
		connect.NewRequest(&greetv1.GreetRequest{Name: "Jane"}),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.Greeting)
}
