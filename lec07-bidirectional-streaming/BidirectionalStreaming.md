**server / main.go**

```
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
```

**client / main.go**

```
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
```
