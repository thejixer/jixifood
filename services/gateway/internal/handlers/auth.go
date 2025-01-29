package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	authPB "github.com/thejixer/jixifood/generated/auth"
	"github.com/thejixer/jixifood/shared/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *HandlerService) HandleRequestOTP(c echo.Context) error {
	body := models.RequestOTPDto{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad input")
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	d := &authPB.RequestOtpRequest{
		PhoneNumber: body.PhoneNumber,
	}
	resp, err := h.gc.AuthClient.RequestOtp(context.Background(), d)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return WriteReponse(c, http.StatusInternalServerError, "Internal server error")
		}
		switch st.Code() {
		case codes.InvalidArgument:
			return WriteReponse(c, http.StatusBadRequest, st.Message())
		case codes.Internal:
			return WriteReponse(c, http.StatusInternalServerError, "Internal server error")
		default:
			// Handle other error codes
			return WriteReponse(c, http.StatusInternalServerError, "An unexpected error occurred")
		}

	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HandlerService) HandleVerifyOTP(c echo.Context) error {
	body := models.VerifyOTPDto{}

	if err := c.Bind(&body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	d := &authPB.VerifyOtpRequest{
		PhoneNumber: body.PhoneNumber,
		Otp:         body.OTP,
	}
	res, err := h.gc.AuthClient.VerifyOtp(context.Background(), d)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return WriteReponse(c, http.StatusInternalServerError, "Internal server error")
		}
		switch st.Code() {
		case codes.InvalidArgument:
			return WriteReponse(c, http.StatusBadRequest, "provided values did not meet the requirements")
		case codes.Internal:
			return WriteReponse(c, http.StatusInternalServerError, "Internal server error")
		case codes.Unauthenticated:
			return WriteReponse(c, http.StatusUnauthorized, "code doesn't match")
		default:
			// Handle other error codes
			return WriteReponse(c, http.StatusInternalServerError, "An unexpected error occurred")
		}
	}

	return c.JSON(http.StatusOK, res)

}
