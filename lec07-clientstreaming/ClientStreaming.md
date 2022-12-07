**server / main.go**

```
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
```

**client / main.go**

```
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
	
```
