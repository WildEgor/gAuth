package proto

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type ProxyService struct{}

func NewProxyService() *ProxyService {
	return &ProxyService{}
}

func (s *ProxyService) Connect(ctx context.Context, request *ConnectRequest) (*ConnectResponse, error) {
	log.Info("Connect..")
	log.Info(request.Client)
	return &ConnectResponse{}, nil
}

func (s *ProxyService) Refresh(ctx context.Context, request *RefreshRequest) (*RefreshResponse, error) {
	log.Info("Refresh..")
	return &RefreshResponse{}, nil
}

func (s *ProxyService) Subscribe(ctx context.Context, request *SubscribeRequest) (*SubscribeResponse, error) {
	log.Info("Subscribe..")
	return &SubscribeResponse{}, nil
}

func (s *ProxyService) Publish(ctx context.Context, request *PublishRequest) (*PublishResponse, error) {
	log.Info("Publish..")
	log.Info(string(request.Data))
	return &PublishResponse{}, nil
}

func (s *ProxyService) RPC(ctx context.Context, request *RPCRequest) (*RPCResponse, error) {
	log.Info("RPC..")
	return &RPCResponse{}, nil
}

func (s *ProxyService) mustEmbedUnimplementedCentrifugoProxyServer() {
	panic("implement me")
}
