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

func (h *Handler) GetForumDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, ok := vars["slug"]
	if !ok {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	userOut, err := h.uc.GetForum(r.Context(), nickname)
	if err != nil {
		utils.Response(w, http.StatusNotFound, nickname)
		return
	}
	utils.Response(w, http.StatusOK, userOut)
	return
}

func (h *Handler) CreateThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	thread := models.Thread{}
	err := easyjson.UnmarshalFromReader(r.Body, &thread)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	thread.Forum = slug

	result, err := h.uc.CreateThread(r.Context(), thread)
	if errors.Is(err, models.ErrorConflict) {
		utils.Response(w, http.StatusConflict, result)
		return
	} else if errors.Is(err, models.ErrorInternal) {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	} else if errors.Is(err, models.ErrorNotFound) {
		utils.Response(w, http.StatusNotFound, "Slug owner is not found")
		return
	}

	utils.Response(w, http.StatusCreated, result)
}

func (h *Handler) GetThreads(w http.ResponseWriter, r *http.Request) {
	var limit, since, desc string
	vars := mux.Vars(r)
	slug, found := vars["slug"]
	if !found {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}
	query := r.URL.Query()
	if limits := query["limit"]; len(limits) > 0 {
		limit = limits[0]
	}
	if sinces := query["since"]; len(sinces) > 0 {
		since = sinces[0]
	}
	if descs := query["desc"]; len(descs) > 0 {
		desc = descs[0]
	}

	result, err := h.uc.GetThreads(r.Context(), slug, limit, since, desc)
	if errors.Is(err, models.ErrorInternal) {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	} else if errors.Is(err, models.ErrorNotFound) {
		utils.Response(w, http.StatusNotFound, "Forum not found")
		return
	}

	utils.Response(w, http.StatusOK, result)

}

func (h *Handler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, found := vars["slug_or_id"]
	if !found {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	thread, err := h.uc.CheckThreadByIdOrSlug(r.Context(), slugOrId)
	if errors.Is(err, models.ErrorInternal) {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	} else if errors.Is(err, models.ErrorNotFound) {
		utils.Response(w, http.StatusNotFound, "Thread not found")
		return
	}

	posts := models.PostsList{}
	err = easyjson.UnmarshalFromReader(r.Body, &posts)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if len(posts) == 0 {
		utils.Response(w, http.StatusCreated, posts)
		return
	}

	createdPosts, err := h.uc.CreatePosts(r.Context(), posts, thread)
	if errors.Is(err, models.ErrorInternal) {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	} else if errors.Is(err, models.ErrorConflict) {
		utils.Response(w, http.StatusConflict, "Conflict during inserting")
		return
	}
	utils.Response(w, http.StatusCreated, createdPosts)

}

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, found := vars["slug_or_id"]
	if !found {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	thread, err := h.uc.CheckThreadByIdOrSlug(r.Context(), slugOrId)
	if errors.Is(err, models.ErrorInternal) {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	} else if errors.Is(err, models.ErrorNotFound) {
		utils.Response(w, http.StatusNotFound, "Thread not found")
		return
	}

	vote := models.Vote{}
	err = easyjson.UnmarshalFromReader(r.Body, &vote)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if thread.ID != 0 {
		vote.Thread = thread.ID
	}

	_, _ := h.uc.Vote(r.Context(), vote, thread) //ЗДЕСЬ!

	threadUpdated, _ := h.uc.CheckThreadByIdOrSlug(r.Context(), slugOrId)
	utils.Response(w, http.StatusOK, threadUpdated)
}
