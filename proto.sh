protoc -I ./api/proto \
  --go_opt=module=gAuth --go_out=internal \
  --go-grpc_opt=module=gAuth --go-grpc_out=internal \
  api/proto/*.proto

read key