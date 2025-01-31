package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	authPB "github.com/thejixer/jixifood/generated/auth"
	"github.com/thejixer/jixifood/services/gateway/internal/adapters"
	apperrors "github.com/thejixer/jixifood/shared/errors"
	"github.com/thejixer/jixifood/shared/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *HandlerService) HandleRequestOTP(c echo.Context) error {
	body := models.RequestOTPDto{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
	}

	d := &authPB.RequestOtpRequest{
		PhoneNumber: body.PhoneNumber,
	}
	resp, err := h.gc.AuthClient.RequestOtp(context.Background(), d)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrInternal.Error())
		}
		switch st.Code() {
		case codes.InvalidArgument:
			return WriteReponse(c, http.StatusBadRequest, st.Message())
		case codes.Internal:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrInternal.Error())
		default:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrUnexpected.Error())
		}

	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HandlerService) HandleVerifyOTP(c echo.Context) error {
	body := models.VerifyOTPDto{}

	if err := c.Bind(&body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
	}

	d := &authPB.VerifyOtpRequest{
		PhoneNumber: body.PhoneNumber,
		Otp:         body.OTP,
	}
	res, err := h.gc.AuthClient.VerifyOtp(context.Background(), d)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrInternal.Error())
		}
		switch st.Code() {
		case codes.InvalidArgument:
			return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
		case codes.Internal:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrInternal.Error())
		case codes.Unauthenticated:
			return WriteReponse(c, http.StatusUnauthorized, apperrors.ErrCodeMismatch.Error())
		default:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrUnexpected.Error())
		}
	}

	return c.JSON(http.StatusOK, res)

}

func (h *HandlerService) HandleME(c echo.Context) error {

	ctx, err := ContextWithCredentials(c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
	}

	resp, err := h.gc.AuthClient.Me(ctx, &authPB.Empty{})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrInternal.Error())
		}
		switch st.Code() {
		case codes.Unauthenticated:
			return WriteReponse(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
		case codes.Internal:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrInternal.Error())
		default:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrUnexpected.Error())
		}
	}

	return c.JSON(http.StatusOK, adapters.MapPBUserToUserDTO(resp))
}

func (h *HandlerService) HandleCreateUser(c echo.Context) error {
	ctx, err := ContextWithCredentials(c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
	}

	body := models.CreateUserDto{}

	if err := c.Bind(&body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
	}

	d := &authPB.CreateUserRequest{
		PhoneNumber: body.PhoneNumber,
		Name:        body.Name,
		RoleId:      body.RoleID,
	}
	resp, err := h.gc.AuthClient.CreateUser(ctx, d)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrInternal.Error())
		}
		switch st.Code() {
		case codes.InvalidArgument:
			if strings.Contains(err.Error(), apperrors.ErrDuplicatePhone.Error()) {
				return WriteReponse(c, http.StatusBadRequest, apperrors.ErrDuplicatePhone.Error())
			}
			return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
		case codes.Unauthenticated:
			return WriteReponse(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
		case codes.Internal:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrInternal.Error())
		case codes.PermissionDenied:
			return WriteReponse(c, http.StatusForbidden, apperrors.ErrForbidden.Error())
		default:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrUnexpected.Error())
		}
	}

	return c.JSON(http.StatusOK, adapters.MapPBUserToUserDTO(resp))
}

func (h *HandlerService) HandleChangeUserRole(c echo.Context) error {
	ctx, err := ContextWithCredentials(c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
	}

	body := models.ChangeUserRoleDto{}

	if err := c.Bind(&body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
	}

	i := c.Param("id")
	intuserID, err := strconv.Atoi(i)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
	}
	userID := uint64(intuserID)
	d := &authPB.ChangeUserRoleRequest{
		UserId: userID,
		RoleId: body.RoleID,
	}
	resp, err := h.gc.AuthClient.ChangeUserRole(ctx, d)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrInternal.Error())
		}
		switch st.Code() {
		case codes.InvalidArgument:
			return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
		case codes.Unauthenticated:
			return WriteReponse(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
		case codes.Internal:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrInternal.Error())
		case codes.PermissionDenied:
			return WriteReponse(c, http.StatusForbidden, apperrors.ErrForbidden.Error())
		default:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrUnexpected.Error())
		}
	}

	return c.JSON(http.StatusOK, adapters.MapPBUserToUserDTO(resp))

}

func (h *HandlerService) HandleEditProfile(c echo.Context) error {
	ctx, err := ContextWithCredentials(c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
	}

	body := models.EditProfileDto{}

	if err := c.Bind(&body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
	}

	d := &authPB.EditProfileRequest{
		Name: body.Name,
	}
	resp, err := h.gc.AuthClient.EditProfile(ctx, d)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrInternal.Error())
		}
		switch st.Code() {
		case codes.InvalidArgument:
			return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
		case codes.Unauthenticated:
			return WriteReponse(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
		case codes.Internal:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrInternal.Error())
		default:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrUnexpected.Error())
		}
	}
	return c.JSON(http.StatusOK, adapters.MapPBUserToUserDTO(resp))

}

func (h *HandlerService) HandleQueryUsers(c echo.Context) error {
	ctx, err := ContextWithCredentials(c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
	}

	t := c.QueryParam("text")
	p := c.QueryParam("page")
	l := c.QueryParam("limit")

	// var page uint64
	// var limit uint64
	// var text string

	page, err := strconv.Atoi(p)
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 10
	}

	d := &authPB.QueryUsersRequest{
		Text:  t,
		Page:  uint64(page),
		Limit: uint64(limit),
	}
	resp, err := h.gc.AuthClient.QueryUsers(ctx, d)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrInternal.Error())
		}
		switch st.Code() {
		case codes.InvalidArgument:
			return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
		case codes.Unauthenticated:
			return WriteReponse(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
		case codes.Internal:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrInternal.Error())
		case codes.PermissionDenied:
			return WriteReponse(c, http.StatusForbidden, apperrors.ErrForbidden.Error())
		default:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrUnexpected.Error())
		}
	}

	var users []*models.UserDto
	for _, user := range resp.Data {
		u := adapters.MapPBUserToUserDTO(user)
		users = append(users, u)
	}
	result := models.QueryUsersDto{
		Data:        users,
		HasNextPage: resp.HasNextPage,
		Total:       resp.Total,
	}
	return c.JSON(http.StatusOK, result)
}
