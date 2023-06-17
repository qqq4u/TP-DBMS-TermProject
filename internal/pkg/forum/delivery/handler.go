package handler

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/models"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/pkg/forum"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/utils"
	"net/http"
)

type Handler struct {
	uc forum.ForumUsecase
}

func NewForumHandler(forumUsecase forum.ForumUsecase) *Handler {
	return &Handler{
		uc: forumUsecase,
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, ok := vars["nickname"]
	if !ok {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	userOut, err := h.uc.GetUser(r.Context(), nickname)
	if err != nil {
		utils.Response(w, http.StatusNotFound, nickname)
		return
	}
	utils.Response(w, http.StatusOK, userOut)
	return
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, ok := vars["nickname"]
	if !ok {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	user := models.User{}
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	user.Nickname = nickname

	result, err := h.uc.CreateUser(r.Context(), user)
	if errors.Is(err, models.ErrorConflict) {
		utils.Response(w, http.StatusConflict, result)
		return
	} else if errors.Is(err, models.ErrorInternal) {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusCreated, result[0])
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, found := vars["nickname"]
	if !found {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	user := models.User{}
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	user.Nickname = nickname

	finalUser, err := h.uc.UpdateUser(r.Context(), user)
	if errors.Is(err, models.ErrorConflict) {
		utils.Response(w, http.StatusConflict, models.Error{Message: "Can't update data"})
		return
	}
	if errors.Is(err, models.ErrorNotFound) {
		utils.Response(w, http.StatusNotFound, models.Error{Message: "User doesn't exists"})
		return
	}
	utils.Response(w, http.StatusOK, finalUser)
}

func (h *Handler) CreateForum(w http.ResponseWriter, r *http.Request) {
	forumInfo := models.Forum{}
	err := easyjson.UnmarshalFromReader(r.Body, &forumInfo)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	result, err := h.uc.CreateForum(r.Context(), forumInfo)
	if errors.Is(err, models.ErrorConflict) {
		utils.Response(w, http.StatusConflict, result)
		return
	} else if errors.Is(err, models.ErrorNotFound) {
		utils.Response(w, http.StatusNotFound, "Forum owner is not found")
		return
	} else if errors.Is(err, models.ErrorInternal) {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusCreated, result)
}
