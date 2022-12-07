package main

import (
	"log"
	"net"
	"strconv"

	pb "lec07-serverstreaming/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedServerStreamingServer
}

func (s server) GetServerResponse(in *pb.Number, srv pb.ServerStreaming_GetServerResponseServer) error {
	msg := "Server processing gRPC server-streaming"
	log.Println(msg)
	// 1부터 클라이언트 단에서 입력한 숫자만큼 루프를 돌면서 message # 숫자 형태로 저장한다.
	// 그리고 저장한 것을 protobuffer에서 설정한 Message에 넣어 클라이언트 단으로 보낸다.
	for i := 1; i <= int(in.Value); i++ {
		resp := pb.Message{
			Message: "message #" + strconv.Itoa(i),
		}
		srv.Send(&resp)
	}
	return nil
}

func main() {
	log.SetFlags(0)

	// listener 셋팅
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// grpc 서버 만들기
	s := grpc.NewServer()
	pb.RegisterServerStreamingServer(s, &server{})

	log.Println("start server. Listening on port 50051")

	// grpc 서버와 listener 연결
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
