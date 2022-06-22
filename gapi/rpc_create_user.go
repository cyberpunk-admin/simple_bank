package gapi

import (
	"context"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	db "github.com/simplebank/db/sqlc"
	"github.com/simplebank/pb"
	"github.com/simplebank/util"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	HashPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed hash password: %s", err)
	}
	arg := db.CreateUserParams{
		UserName:     req.GetUsername(),
		HashPassword: HashPassword,
		FullName:     req.GetFullName(),
		Email:        req.GetEmail(),
	}
	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exist: %s", err)
			}
		}
		return nil, status.Errorf(codes.AlreadyExists, "failed to create user: %s", err)
	}

	rsq := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsq, nil
}
