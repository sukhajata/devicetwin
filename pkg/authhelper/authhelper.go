package authhelper

import (
	"context"
	"errors"
	pbAuth "github.com/sukhajata/ppauth"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
	"time"
)

// GetContextWithAuth returns a new context with authhelper metadata
func GetContextWithAuth(token string) (context.Context, context.CancelFunc) {
	md := metadata.Pairs("authorization", token)
	outgoing := metadata.NewOutgoingContext(context.Background(), md)
	//set timeout
	return context.WithTimeout(outgoing, 5*time.Second)
}

// GetTokenFromContext gets a token from a context
func GetTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("Could not get metadata")
	}
	token := md["authorization"]
	//WriteToLog(fmt.Sprintf("Using token |%v|\n", token))
	if token == nil || len(token) == 0 {
		return "", errors.New("Missing token")
	}

	return token[0], nil
}

// GetTokenFromHeader - get the token from an http authhelper header
func GetTokenFromHeader(r *http.Request) (string, error) {
	var header = r.Header.Get("authorization")
	header = strings.TrimSpace(header)
	if header == "" {
		return "", errors.New("No token")
	}

	splitHeader := strings.Split(header, "Bearer ")
	if len(splitHeader) != 2 {
		return "", errors.New("Malformed authhelper header")
	}
	token := splitHeader[1]

	return token, nil
}

// CheckToken checks that a token is valid and has one of required roles
// returns error if validation fails, returns username if valid
func CheckToken(grpcAuthClient pbAuth.AuthServiceClient, token string, allowedRoles []string) (string, error) {
	if len(token) == 0 {
		return "", errors.New("Missing token")
	}
	authRequest := &pbAuth.AuthRequest{
		Token:        token,
		AllowedRoles: allowedRoles,
	}
	authResponse, err := grpcAuthClient.CheckAuth(context.Background(), authRequest)
	if err != nil {
		return "", err
	}
	if authResponse.GetResult() == false {
		return "", errors.New(authResponse.GetMessage())
	}

	return authResponse.GetUsername(), nil
}
