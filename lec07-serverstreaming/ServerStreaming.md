server / main.go

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

client / main.go
    // protoBuffer에 설정된 값의 이름을 그대로 사용하여 보내고자 하는 숫자를 담는다.
	var val int32 = 5
	in := &pb.Number{Value: val}

	stream, err := client.GetServerResponse(context.Background(), in)
	if err != nil {
		log.Fatalf("open stream error %v", err)
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