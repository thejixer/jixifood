package handlers

import (
	pb "github.com/thejixer/jixifood/generated/menu"
	"github.com/thejixer/jixifood/services/menu/internal/logic"
)

type MenuLogicInterface interface{}

type MenuHandler struct {
	pb.UnimplementedMenuServiceServer
	MenuLogic MenuLogicInterface
}

func NewMenuHandler(menuLogic *logic.MenuLogic) *MenuHandler {
	return &MenuHandler{
		MenuLogic: menuLogic,
	}
}
