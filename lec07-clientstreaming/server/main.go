package main

import (
	"io"
	"log"
	"net"

	pb "lec07-clientstreaming/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedClientStreamingServer
}

// 서버는 클라이언트 단에서 보낸 데이터의 개수를 세어 나온 값을 클라이언트 단으로 전송한다.
func (s *server) GetServerResponse(stream pb.ClientStreaming_GetServerResponseServer) error {
	log.Println("Server processing gRPC client streaming")
	var count int32 = 0

	for {
		// 클라이언트 단에서 받은 데이터의 내용은 중요하지 않기 때문에 _ 처리하였다.
		_, err := stream.Recv()

		// 데이터의 전송이 끝나면
		if err == io.EOF {
			// 다시 클라이언트 단에 카운트된 데이터의 개수를 보내고 끝낸다.
			stream.SendAndClose(&pb.Number{Value: count})
			return nil
		}
		if err != nil {
			return err
		}
		count++
	}
}

func main() {
	log.SetFlags(0)

	// Listener 생성
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	// grpc 서버 생성
	grpcServer := grpc.NewServer()
	pb.RegisterClientStreamingServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	// grpc 서버와 Listener 연결
	log.Println("Starting server. Listening on port 50051.")
	grpcServer.Serve(lis)
}
