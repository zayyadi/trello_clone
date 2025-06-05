package handlers

import (
	"net/http"
	"strconv"

	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/services"

	"github.com/gin-gonic/gin"
)

type ListHandler struct {
	listService *services.ListService
}

func NewListHandler(listService *services.ListService) *ListHandler {
	return &ListHandler{listService: listService}
}

func (h *ListHandler) CreateList(c *gin.Context) {
	userID, _ := c.Get("userID")
	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid board ID")
		return
	}

	var req CreateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	list, err := h.listService.CreateList(req.Name, uint(boardID), userID.(uint), req.Position)
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusCreated, "List created successfully", MapListToResponse(list, false)) // Don't include cards by default
}

func (h *ListHandler) GetListsByBoardID(c *gin.Context) {
	userID, _ := c.Get("userID")
	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid board ID")
		return
	}

	lists, err := h.listService.GetListsByBoardID(uint(boardID), userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	var listResponses []ListResponse
	for _, l := range lists {
		listResponses = append(listResponses, MapListToResponse(&l, true)) // Include cards when fetching lists for a board
	}
	RespondWithSuccess(c, http.StatusOK, "Lists retrieved successfully", listResponses)
}

func (h *ListHandler) UpdateList(c *gin.Context) {
	userID, _ := c.Get("userID")
	listIDStr := c.Param("listID")
	listID, err := strconv.ParseUint(listIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid list ID")
		return
	}

	var req UpdateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	list, err := h.listService.UpdateList(uint(listID), req.Name, req.Position, userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusOK, "List updated successfully", MapListToResponse(list, false))
}

func (h *ListHandler) DeleteList(c *gin.Context) {
	userID, _ := c.Get("userID")
	listIDStr := c.Param("listID")
	listID, err := strconv.ParseUint(listIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid list ID")
		return
	}

	err = h.listService.DeleteList(uint(listID), userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusOK, "List deleted successfully", nil)
}

// Helper to map model.List to ListResponse
func MapListToResponse(list *models.List, includeCards bool) ListResponse {
	if list == nil {
		return ListResponse{}
	}
	resp := ListResponse{
		ID:        list.ID,
		Name:      list.Name,
		BoardID:   list.BoardID,
		Position:  list.Position,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
	if includeCards && len(list.Cards) > 0 {
		resp.Cards = []CardResponse{}
		for _, c := range list.Cards {
			resp.Cards = append(resp.Cards, MapCardToResponse(&c, true)) // Include assigned user for cards
		}
	}
	return resp
}
