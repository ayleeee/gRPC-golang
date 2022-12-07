package main

import (
	"context"
	pb "lec07-hello-grpc/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedMyServiceServer
}

// 제곱 연산을 수행하는 함수를 rpc 함수로 작성
func (s *server) MyFunction(ctx context.Context, in *pb.MyNumber) (*pb.MyNumber, error) {
	return &pb.MyNumber{Value: in.GetValue() * in.GetValue()}, nil
}

func main() {
	log.SetFlags(0)
	// listener 생성
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// grpc 서버 생성
	grpcServer := grpc.NewServer()
	pb.RegisterMyServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	// grpc 서버와 listener 연결
	log.Println("Starting server. Listening on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
