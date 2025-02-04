package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	menuPB "github.com/thejixer/jixifood/generated/menu"
	"github.com/thejixer/jixifood/services/gateway/internal/adapters"
	apperrors "github.com/thejixer/jixifood/shared/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/thejixer/jixifood/shared/models"
)

func (h *HandlerService) HandleCreateCategory(c echo.Context) error {
	fmt.Println("create category 1")
	ctx, err := ContextWithCredentials(c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
	}

	fmt.Println("create category 2")
	body := models.CreateCategoryDto{}

	if err := c.Bind(&body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
	}

	if err := c.Validate(body); err != nil {
		fmt.Println(err)
		return WriteReponse(c, http.StatusBadRequest, apperrors.ErrInputRequirements.Error())
	}
	fmt.Println("create category 3")

	d := &menuPB.CreateCategoryRequest{
		Name:           body.Name,
		Description:    body.Description,
		IsQuantifiable: body.IsQuantifiable,
	}
	resp, err := h.gc.MenuClient.CreateCategory(ctx, d)
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
		case codes.PermissionDenied:
			return WriteReponse(c, http.StatusForbidden, apperrors.ErrForbidden.Error())
		default:
			return WriteReponse(c, http.StatusInternalServerError, apperrors.ErrUnexpected.Error())
		}
	}

	return c.JSON(http.StatusOK, adapters.MapPBCategoryToCategoryDTO(resp))
}
