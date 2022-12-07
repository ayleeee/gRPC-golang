사전에 설치해야하는 것

(1) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
(2) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

.proto 파일을 사용하기 위해선 

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    [.proto가 담겨 있는 폴더 이름]/[name].proto