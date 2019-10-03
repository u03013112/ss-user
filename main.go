package main

import (
	"log"
	"net"

	pb "github.com/u03013112/ss-pb/user"
	"github.com/u03013112/ss-user/user"
	"google.golang.org/grpc"
)

const (
	port = ":50000"
)

func main() {
	user.InitDB()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listen %s", port)
	s := grpc.NewServer()
	pb.RegisterSSUserServer(s, &user.Srv{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
