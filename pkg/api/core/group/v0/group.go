package v0

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/core/group"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/store"
)

type GroupHandler struct {
	group.GroupHandlerInterface

	store store.Store
}

func NewGroupHandler(store store.Store) *GroupHandler {
	return &GroupHandler{
		store: store,
	}
}

func (h *GroupHandler) FindAll(ctx *gin.Context) {
	groupList := []*core.Group{}
	f := func(n int) []interface{} {
		m := []interface{}{}
		for i := 0; i < n; i++ {
			group := &core.Group{}
			groupList = append(groupList, group)
			m = append(m, group)
		}
		return m
	}
	h.store.List("group/", f)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"groups": groupList,
	})
}

func (h *GroupHandler) Find(ctx *gin.Context) {
	groupID := ctx.Param("group_id")
	obj := h.store.Get(getKey(groupID))
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Group `%s` is not found.", groupID), nil)
		return
	}

	group := obj.(core.Group)
	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"group": group,
	})
}

func (h *GroupHandler) Create(ctx *gin.Context) {
	var request core.Group

	err := ctx.Bind(&request)
	if err != nil {
		log.Println(err)
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if request.ID == "" {
		log.Println("id is empty")
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: ID is empty."), nil)
		return
	}

	key := getKey(request.ID)
	obj := h.store.Get(key)
	if obj != nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: group `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"group": request,
	})
}

func (h *GroupHandler) Update(ctx *gin.Context) {
	groupID := ctx.Param("group_id")
	var request core.Group

	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if request.ID != groupID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: can't change id."), nil)
		return
	}

	key := getKey(request.ID)
	obj := h.store.Get(key)
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Error: group `%s` is not found.", request.ID), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"group": request,
	})
}

func (h *GroupHandler) Delete(ctx *gin.Context) {
	groupID := ctx.Param("group_id")

	key := getKey(groupID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)
}

func getKey(id string) string {
	return "group/" + id
}
