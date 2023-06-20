package handler

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/models"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/pkg/forum"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/utils"
	"net/http"
	"strconv"
	"strings"
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
	if limitTmp := query["limit"]; len(limitTmp) > 0 {
		limit = limitTmp[0]
	}
	if sinceTmp := query["since"]; len(sinceTmp) > 0 {
		since = sinceTmp[0]
	}
	if descTmp := query["desc"]; len(descTmp) > 0 {
		desc = descTmp[0]
	}

	result, err := h.uc.GetThreads(r.Context(), slug, limit, since, desc)
	if errors.Is(err, models.ErrorNotFound) {
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
	if errors.Is(err, models.ErrorNotFound) {
		utils.Response(w, http.StatusNotFound, "Can't find post author")
		return
	} else if errors.Is(err, models.ErrorConflict) {
		utils.Response(w, http.StatusConflict, models.Error{Message: "Wrong post parent"})
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

	if err = h.uc.Vote(r.Context(), vote); errors.Is(err, models.ErrorNotFound) {
		utils.Response(w, http.StatusNotFound, "User not found")
		return
	}

	threadUpdated, _ := h.uc.CheckThreadByIdOrSlug(r.Context(), slugOrId)
	utils.Response(w, http.StatusOK, threadUpdated)
}

func (h *Handler) GetThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, found := vars["slug_or_id"]
	if !found {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	result, err := h.uc.CheckThreadByIdOrSlug(r.Context(), slugOrId)
	if errors.Is(err, models.ErrorNotFound) {
		utils.Response(w, http.StatusNotFound, "Thread not found")
		return
	}

	utils.Response(w, http.StatusOK, result)
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, found := vars["id"]
	if !found {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	query := r.URL.Query()
	var related []string
	if relatedTmp := query["related"]; len(relatedTmp) > 0 {
		related = strings.Split(relatedTmp[0], ",")
	}

	result, err := h.uc.GetPost(r.Context(), id, related)
	if errors.Is(err, models.ErrorNotFound) {
		utils.Response(w, http.StatusNotFound, "Post not found")
		return
	}

	utils.Response(w, http.StatusOK, result)
}

func (h *Handler) GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	var limit, since, desc, sort string

	vars := mux.Vars(r)
	slugOrId, _ := vars["slug_or_id"]

	query := r.URL.Query()
	if limitTmp := query["limit"]; len(limitTmp) > 0 {
		limit = limitTmp[0]
	}
	if sinceTmp := query["since"]; len(sinceTmp) > 0 {
		since = sinceTmp[0]
	}
	if descTmp := query["desc"]; len(descTmp) > 0 {
		desc = descTmp[0]
	}
	if sortTmp := query["sort"]; len(sortTmp) > 0 {
		sort = sortTmp[0]
	}

	thread, err := h.uc.CheckThreadByIdOrSlug(r.Context(), slugOrId)
	if errors.Is(err, models.ErrorNotFound) {
		utils.Response(w, http.StatusNotFound, "Thread not found")
		return
	}

	result, _ := h.uc.GetThreadPosts(r.Context(), limit, since, desc, sort, thread.ID)

	utils.Response(w, http.StatusOK, result)
}

func (h *Handler) UpdateThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, _ := vars["slug_or_id"]
	thread := models.Thread{}
	easyjson.UnmarshalFromReader(r.Body, &thread)

	idInt, err := strconv.Atoi(slugOrId)
	if err != nil {
		thread.Slug = slugOrId
	} else {
		thread.ID = idInt
	}

	result, err := h.uc.UpdateThread(r.Context(), thread)

	if errors.Is(err, models.ErrorNotFound) {
		utils.Response(w, http.StatusNotFound, "Thread not found")
		return
	}

	utils.Response(w, http.StatusOK, result)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	var limit, since, desc string
	vars := mux.Vars(r)
	slug, _ := vars["slug"]

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

	result, err := h.uc.GetUsers(r.Context(), slug, limit, since, desc)
	if errors.Is(err, models.ErrorInternal) {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	} else if errors.Is(err, models.ErrorNotFound) {
		utils.Response(w, http.StatusNotFound, "Forum not found")
		return
	}

	utils.Response(w, http.StatusOK, result)

}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, _ := vars["id"]
	id, _ := strconv.Atoi(idStr)
	postUpdateInfo := models.PostUpdate{ID: id}
	easyjson.UnmarshalFromReader(r.Body, &postUpdateInfo)

	result, err := h.uc.UpdatePost(r.Context(), postUpdateInfo)
	if errors.Is(err, models.ErrorNotFound) {
		utils.Response(w, http.StatusNotFound, "Post not found")
		return
	}

	utils.Response(w, http.StatusOK, result)
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	utils.Response(w, http.StatusOK, h.uc.GetStatus())
}

func (h *Handler) Clear(w http.ResponseWriter, r *http.Request) {
	h.uc.Clear()
	utils.Response(w, http.StatusOK, nil)
}
