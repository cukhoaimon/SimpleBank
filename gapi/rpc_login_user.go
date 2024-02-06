package gapi

import (
	"context"
	"database/sql"
	"errors"
	db "github.com/cukhoaimon/SimpleBank/db/sqlc"
	"github.com/cukhoaimon/SimpleBank/pb"
	"github.com/cukhoaimon/SimpleBank/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user [%s] is not found in database", req.GetUsername())
		}

		return nil, status.Errorf(codes.Internal, "fail to get user [%s] from database", req.GetUsername())
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "wrong password")

	}

	token, accessTokenPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.TokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "fail to create access token")
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "fail to create refresh access token")
	}

	arg := db.CreateSessionParams{
		ID:           accessTokenPayload.Id,
		Username:     accessTokenPayload.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	}

	session, err := server.store.CreateSession(ctx, arg)
	if err != nil {
		return nil, status.Error(codes.Internal, "fail to create session")
	}

	res := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           token,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
		User:                  convertUser(user),
	}

	return res, nil
}
