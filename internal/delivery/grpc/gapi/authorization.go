package gapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/cukhoaimon/SimpleBank/pkg/token"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	authorizationHeader = "authorization"
	authorizationType   = "bearer"
)

func (handler *Handler) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	value := md.Get(authorizationHeader)
	if len(value) == 0 {
		return nil, errors.New("missing authorization header")
	}

	authHeader := value[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, errors.New("invalid authorization format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationType {
		return nil, fmt.Errorf("unsupported authorization header: %s", authType)
	}

	accessToken := fields[1]
	payload, err := handler.TokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
