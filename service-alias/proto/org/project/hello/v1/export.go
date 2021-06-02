package hello

import grpc "google.golang.org/grpc"

func EchoService_serviceDesc() grpc.ServiceDesc {
	return _EchoService_serviceDesc
}
