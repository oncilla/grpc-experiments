package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
	greetv1 "github.com/oncilla/grpc-experiments/http3/gen/greet/v1"
	"github.com/oncilla/grpc-experiments/http3/gen/greet/v1/greetv1connect"

	"github.com/bufbuild/connect-go"
)

var (
	certFile = flag.String("cert", "./testdata/cert.pem", "certificate file")
	keyFile  = flag.String("key", "./testdata/priv.key", "key file")
)

type GreetServer struct{}

func (s *GreetServer) Greet(
	ctx context.Context,
	req *connect.Request[greetv1.GreetRequest],
) (*connect.Response[greetv1.GreetResponse], error) {

	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func main() {
	greeter := &GreetServer{}
	mux := http.NewServeMux()
	path, handler := greetv1connect.NewGreetServiceHandler(greeter)
	mux.Handle(path, handler)

	err := http3.ListenAndServeQUIC("localhost:8080", *certFile, *keyFile, mux)
	log.Fatal(err)
}
