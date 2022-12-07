package main

import (
	"log"
	"strconv"

	pb "lec07-clientstreaming/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.SetFlags(0)

	// 서버와의 연결
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect : %s", err)
	}
	defer conn.Close()

	client := pb.NewClientStreamingClient(conn)

	stream, err := client.GetServerResponse(context.Background())
	if err != nil {
		log.Fatalf("problem : %s", err)
	}

	done := make(chan bool)

	go func() {
		for i := 1; i <= 5; i++ {
			msg := &pb.Message{
				Message: "[client to server] message #" + strconv.Itoa(i),
			}
			// msg를 서버단으로 전송한다.
			stream.Send(msg)
			log.Println(msg.GetMessage())

		}
		// 서버단에서 데이터의 개수에 대한 정보를 보낸다.
		// stream을 닫고 해당 정보를 받는다.
		response, _ := stream.CloseAndRecv()
		log.Println("[server to client] " + strconv.Itoa(int(response.GetValue())))
		close(done)
	}()
	<-done

}
