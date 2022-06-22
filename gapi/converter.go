package gapi

import (
	db "github.com/simplebank/db/sqlc"
	"github.com/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username: user.UserName,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangeAt),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
