package handlers

import (
	"net/http"
	"strconv"

	"github.com/zayyadi/trello/dto" // Import new dto package
	"github.com/zayyadi/trello/services"

	"github.com/gin-gonic/gin"
)

type ListHandler struct {
	listService *services.ListService // Assuming ListService doesn't need an interface yet, or it's services.ListServiceInterface
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

	var req dto.CreateListRequest // Use dto type
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	list, err := h.listService.CreateList(req.Name, uint(boardID), userID.(uint), req.Position)
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusCreated, "List created successfully", dto.MapListToResponse(list, false)) // Use dto mapper
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
	var listResponses []dto.ListResponse // Use dto type
	for _, l := range lists {
		listResponses = append(listResponses, dto.MapListToResponse(&l, true)) // Use dto mapper
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

	var req dto.UpdateListRequest // Use dto type
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	list, err := h.listService.UpdateList(uint(listID), req.Name, req.Position, userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusOK, "List updated successfully", dto.MapListToResponse(list, false)) // Use dto mapper
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

// MapListToResponse function is now in dto/list_dto.go
// MapCardToResponse will be in dto/card_dto.go
