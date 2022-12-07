package main

import (
	"context"
	"io"
	"log"

	pb "lec07-serverstreaming/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.SetFlags(0)

	// 서버와의 연결
	conn, err := grpc.Dial("localhost:50005", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can not connect to %v", err)
	}
	defer conn.Close()
	client := pb.NewServerStreamingClient(conn)

	// rpc 함수에서 필요로 하는 값 설정
	var val int32 = 5
	in := &pb.Number{Value: val}

	stream, err := client.GetServerResponse(context.Background(), in)
	if err != nil {
		log.Fatalf("stream error %v", err)
	}

	done := make(chan bool)

	go func() {
		for {
			// 서버단에서 전송하는 데이터를 받는다.
			res, err := stream.Recv()
			// 더 이상 데이터가 오지 않으면 끝난 것이므로 종료한다.
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			// 받은 데이터를 프린트한다.
			log.Printf("[server to client] %s", res.Message)
		}
	}()
	<-done
}
