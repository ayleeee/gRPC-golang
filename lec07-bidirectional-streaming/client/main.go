package main

import (
	"io"
	"log"
	"strconv"

	pb "lec07-bidirectional-streaming/proto"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// 5개의 데이터를 한 번에 처리하기 위함
	msg [5]string
	req pb.Message
)

func main() {
	log.SetFlags(0)

	// 서버와의 연결
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect : %s", err)
	}
	defer conn.Close()

	client := pb.NewBidirectionalClient(conn)
	stream, err := client.GetServerResponse(context.Background())
	if err != nil {
		log.Fatalf("problem : %s", err)
	}

	done := make(chan bool)

	// 1. 데이터를 만들어서 서버단으로 전송한다.
	go func() {
		for i := 0; i < 5; i++ {
			msg[i] = "message # " + strconv.Itoa(i)
			req = pb.Message{Message: msg[i]}
			if err := stream.Send(&req); err != nil {
				log.Fatalf("can not send %v", err)
			}
			log.Println("[client to server] message # " + req.GetMessage())
		}
		if err := stream.CloseSend(); err != nil {
			log.Println(err)
		}

	}()

	// 2. 클라이언트 -> 서버로 향한 정보를 다시 서버가 클라이언트 단으로 전송한다.
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			resMsg := res.Message
			log.Println("[server to client] message # " + resMsg)
		}
	}()
	<-done

}
