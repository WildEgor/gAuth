package proto

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"net"
)

// GRPCServer struct
type GRPCServer struct {
	// vt *vt.CentrifugeValidateTokenHandler
}

func NewGRPCServer(
// vt *vt.CentrifugeValidateTokenHandler,
) *GRPCServer {
	return &GRPCServer{
		// vt: vt,
	}
}

func (s *GRPCServer) mustEmbedUnimplementedAdderServer() {
}

// Add method for calculate X + Y
func (s *GRPCServer) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {

	//res, err := s.vt.Handle(req)
	//if err != nil {
	//	return nil, err
	//}

	res := 255

	return &AddResponse{
		Result: int32(res),
	}, nil
}

func (s *GRPCServer) Init() (*grpc.Server, error) {
	// Create new gRPC server instance
	instance := grpc.NewServer()
	srv := &GRPCServer{}

	// Register gRPC server
	RegisterAdderServer(instance, srv)

	// Listen on port 8080
	l, err := net.Listen("tcp", ":8088")
	if err != nil {
		return nil, err
	}

	go func() {
		// Start gRPC server
		if err := instance.Serve(l); err != nil {
			log.Fatal(err)
		}
	}()

	return instance, nil
}
