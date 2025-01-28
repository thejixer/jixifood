package handlers

import (
	"context"
	"fmt"
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
		fmt.Println("xxx")
		return echo.NewHTTPError(http.StatusBadRequest, "bad input")
	}

	if err := c.Validate(body); err != nil {
		fmt.Println("yyy")
		fmt.Println(err)
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}
	fmt.Println("lol")
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
