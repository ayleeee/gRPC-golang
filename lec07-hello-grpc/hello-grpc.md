아래와 같은 함수를 가지고 간단한 제곱 연산을 수행한다.
func MyFunction (num int) int {
    return num * num 
}

server / main.go 에 rpc 함수를 정의한다.
func (s *server) MyFunction(ctx context.Context, in *pb.MyNumber) (*pb.MyNumber, error) {
	return &pb.MyNumber{Value: in.GetValue() * in.GetValue()}, nil
}

client / main.go 에서 응답을 주고 받을 통로와 변수의 값을 전달할 통로를 만든다.
stub := pb.NewMyServiceClient(conn)

ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

res, err := stub.MyFunction(ctx, &pb.MyNumber{
		Value: 4,
})