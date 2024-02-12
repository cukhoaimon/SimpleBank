package gapi

import (
	"context"
	"database/sql"
	"errors"
	"github.com/cukhoaimon/SimpleBank/internal/delivery/grpc/pb"
	db "github.com/cukhoaimon/SimpleBank/internal/usecase/sqlc"
	"github.com/cukhoaimon/SimpleBank/internal/usecase/val"
	"github.com/cukhoaimon/SimpleBank/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (handler *Handler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := handler.Store.GetUser(ctx, req.GetUsername())
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

	token, accessTokenPayload, err := handler.TokenMaker.CreateToken(req.Username, handler.Config.TokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "fail to create access token")
	}

	refreshToken, refreshTokenPayload, err := handler.TokenMaker.CreateToken(req.Username, handler.Config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "fail to create refresh access token")
	}

	metadata := handler.extractMetadata(ctx)
	arg := db.CreateSessionParams{
		ID:           accessTokenPayload.Id,
		Username:     accessTokenPayload.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	}

	session, err := handler.Store.CreateSession(ctx, arg)
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

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
