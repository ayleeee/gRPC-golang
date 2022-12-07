package main

import (
	"context"
	"log"
	"time"

	pb "lec07-hello-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.SetFlags(0)

	// 서버와의 연결
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can not connect to : %v", err)
	}
	defer conn.Close()

	stub := pb.NewMyServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// rpc 함수에 필요한 값 전달하기
	res, err := stub.MyFunction(ctx, &pb.MyNumber{
		Value: 4,
	})

	if err != nil {
		log.Fatalf("can not calculate", err)
	}
	log.Printf("gRPC result: %v", res.GetValue())
}
