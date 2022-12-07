package main

import (
	"io"
	"log"
	"net"

	pb "lec07-bidirectional-streaming/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedBidirectionalServer
}

// 클라이언트가 5개의 정보를 보낼 예정이고 한꺼번에 처리할 것이기 때문에 배열을 만들어준다.
var receivedMsg [5]string
var i int = -1

func (s server) GetServerResponse(stream pb.Bidirectional_GetServerResponseServer) error {
	log.Println("Server processing gRPC bidirectional streaming")

	for {
		// 클라이언트 단으로 부터 정보를 받는다.
		req, err := stream.Recv()
		i++
		// 만약 데이터가 더 이상 오지 않는다면
		if err == io.EOF {
			// receivedMsg에 저장되어 있는 것들을 하나씩 꺼내서 다시 클라이언트 단으로 보낸다.
			for j := range receivedMsg {
				resp := pb.Message{Message: receivedMsg[j]}
				if err := stream.Send(&resp); err != nil {
					log.Printf("send error %v", err)
				}
			}
			return nil
		}
		// 아직 데이터가 오는 중이라면
		if err != io.EOF {
			reqMsg := req.Message
			// 순서대로 receivedMsg에 집어넣는다.
			receivedMsg[i] = reqMsg
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}
	}
}

func main() {
	log.SetFlags(0)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBidirectionalServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	log.Println("Starting server. Listening on port 50051.")
	grpcServer.Serve(lis)
}
