package handlers

import (
	"errors"

	apperrors "github.com/thejixer/jixifood/shared/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGetRequesterError(err error) error {

	if errors.Is(err, apperrors.ErrMissingMetaData) ||
		errors.Is(err, apperrors.ErrMissingToken) ||
		errors.Is(err, apperrors.ErrUnauthorized) {
		return status.Error(codes.Unauthenticated, apperrors.ErrUnauthorized.Error())
	}
	return status.Error(codes.Internal, apperrors.ErrUnexpected.Error())

}
