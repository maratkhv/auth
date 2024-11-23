package grpc

import (
	"context"
	"errors"
	"time"

	"auth/protobuf/generated/auth/v1"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type serverAPI struct {
	auth.UnimplementedAuthServer
	api Auth
}

type Auth interface {
	Login(ctx context.Context, login string, password string) (token string, keySinceUnix int64, err error)
	Register(ctx context.Context, login string, password string) (id int64, err error)
	GetPublicKey(ctx context.Context) (key []byte, createTime time.Time, err error)
}

func RegisterServer(srv *grpc.Server, authImpl Auth) {
	auth.RegisterAuthServer(srv, &serverAPI{api: authImpl})
}

func (s *serverAPI) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if req.GetLogin() == "" {
		return nil, errors.New("login is required")
	}
	if req.GetPassword() == "" {
		return nil, errors.New("no password")
	}

	token, t, err := s.api.Login(ctx, req.Login, req.Password)
	return &auth.LoginResponse{
		Token: token,
		PubKeySince: &timestamppb.Timestamp{
			Seconds: t,
		},
	}, err
}

func (s *serverAPI) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if req.GetLogin() == "" {
		return nil, errors.New("login is required")
	}
	if req.GetPassword() == "" {
		return nil, errors.New("no password")
	}
	id, err := s.api.Register(ctx, req.Login, req.Password)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{Id: id}, nil
}

func (s *serverAPI) GetPublicKey(ctx context.Context, r *auth.GetPublicKeyRequest) (*auth.GetPublicKeyResponse, error) {
	key, t, err := s.api.GetPublicKey(ctx)
	return &auth.GetPublicKeyResponse{
		PubKey: key,
		CreatedAt: &timestamppb.Timestamp{
			Seconds: t.Unix(),
		},
	}, err
}
