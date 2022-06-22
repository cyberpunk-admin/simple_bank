package gapi

import (
	"context"
	db "github.com/simplebank/db/sqlc"
	"github.com/simplebank/pb"
	"github.com/simplebank/util"

	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "not such user: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "falied to get user: %s", err)
	}

	err = util.CheckPassword(req.GetPassword(), user.HashPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password: %s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(req.GetUsername(), server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "falied to create access token %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(req.GetUsername(), server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "falied to refresh access token %s", err)
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserName:     user.UserName,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "falied to create session %s", err)
	}

	rsp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  convertUser(user),
	}
	return rsp, nil
}
