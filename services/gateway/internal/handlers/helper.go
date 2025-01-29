package handlers

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"
	authPB "github.com/thejixer/jixifood/generated/auth"
	"github.com/thejixer/jixifood/shared/constants"
	apperrors "github.com/thejixer/jixifood/shared/errors"
	"google.golang.org/grpc/metadata"
)

func GetME(c echo.Context, h *HandlerService) (*authPB.User, error) {
	tokenString, err := GetToken(c)

	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("auth", tokenString))

	resp, err := h.gc.AuthClient.Me(ctx, &authPB.Empty{})
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func GetToken(c echo.Context) (string, error) {
	req := c.Request()
	authSlice := req.Header["Auth"]

	if len(authSlice) == 0 {
		return "", apperrors.ErrMissingToken
	}

	s := strings.Split(authSlice[0], " ")

	if len(s) != 2 || s[0] != constants.Bearer {
		return "", apperrors.ErrBadTokenFormat
	}

	return s[1], nil
}

func ContextWithCredentials(c echo.Context) (context.Context, error) {
	tokenString, err := GetToken(c)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs("auth", tokenString)), nil
}
