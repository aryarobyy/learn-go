package controller

import (
	"auth/internal/helper"
	"auth/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserController struct {
	service service.UserService
}

func NewUserController(s service.UserService) *UserController {
	return &UserController{service: s}
}

func (h *UserController) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, strErr := strconv.Atoi(idStr)
	if strErr != nil {
		helper.RespondError(w, http.StatusBadRequest, strErr)
		return
	}

	res, err := h.service.GetById(r.Context(), id)
	if err != nil {
		helper.RespondError(w, http.StatusBadRequest, err)
		return
	}

	helper.RespondSuccess(w, http.StatusAccepted, res)
}

func (h *UserController) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	res, err := h.service.GetByUsername(r.Context(), username)
	if err != nil {
		helper.RespondError(w, http.StatusBadRequest, err)
		return
	}

	helper.RespondSuccess(w, http.StatusAccepted, res)
}
