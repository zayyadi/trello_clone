package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zayyadi/trello/dto" // Import new dto package
	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/services"
)

type BoardHandler struct {
	boardService services.BoardServiceInterface // Use interface type
}

func NewBoardHandler(boardService services.BoardServiceInterface) *BoardHandler { // Use interface type
	return &BoardHandler{boardService: boardService}
}

func (h *BoardHandler) CreateBoard(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}

	var req dto.CreateBoardRequest // Use dto type
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	board, err := h.boardService.CreateBoard(req.Name, req.Description, userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusCreated, "Board created successfully", dto.MapBoardToResponse(board, true, false)) // Use dto mapper
}

func (h *BoardHandler) GetBoardByID(c *gin.Context) {
	userID, _ := c.Get("userID")
	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid board ID")
		return
	}

	// Optionally load lists and members here or have service do it based on flags
	// For now, load everything for GetByID
	fullBoard, err := h.boardService.GetBoardByID(uint(boardID), userID.(uint)) // Service already preloads owner, members
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	// Manually fetch lists if not preloaded by GetBoardByID (depends on service implementation)
	// Assuming lists are not preloaded by default in BoardService.GetBoardByID for this specific path
	// but GetListsByBoardID service method would be used.
	// For simplicity, we will use the BoardResponse mapping which can conditionally include them.
	RespondWithSuccess(c, http.StatusOK, "Board retrieved successfully", dto.MapBoardToResponse(fullBoard, true, true)) // Use dto mapper
}

func (h *BoardHandler) GetBoardsForUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	boards, err := h.boardService.GetBoardsForUser(userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	var boardResponses []dto.BoardResponse // Use dto type
	for _, b := range boards {
		boardResponses = append(boardResponses, dto.MapBoardToResponse(&b, true, false)) // Use dto mapper
	}
	RespondWithSuccess(c, http.StatusOK, "Boards retrieved successfully", boardResponses)
}

func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	userID, _ := c.Get("userID")
	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid board ID")
		return
	}

	var req dto.UpdateBoardRequest // Use dto type
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	board, err := h.boardService.UpdateBoard(uint(boardID), req.Name, req.Description, userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusOK, "Board updated successfully", dto.MapBoardToResponse(board, true, false)) // Use dto mapper
}

func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	userID, _ := c.Get("userID")
	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid board ID")
		return
	}

	err = h.boardService.DeleteBoard(uint(boardID), userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusOK, "Board deleted successfully", nil)
}

func (h *BoardHandler) AddMemberToBoard(c *gin.Context) {
	userID, _ := c.Get("userID")
	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid board ID")
		return
	}

	var req dto.AddMemberRequest // Use dto type
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	if req.Email == nil && req.UserID == nil {
		RespondWithError(c, http.StatusBadRequest, "Either email or userID must be provided to add a member")
		return
	}
	if req.Email != nil && req.UserID != nil {
		RespondWithError(c, http.StatusBadRequest, "Provide either email or userID, not both")
		return
	}

	boardMember, err := h.boardService.AddMemberToBoard(uint(boardID), req.Email, req.UserID, userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusCreated, "Member added to board successfully", dto.MapBoardMemberToResponse(boardMember)) // Use dto mapper
}

func (h *BoardHandler) RemoveMemberFromBoard(c *gin.Context) {
	userID, _ := c.Get("userID")
	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid board ID")
		return
	}

	memberUserIDStr := c.Param("memberUserID")
	memberUserID, err := strconv.ParseUint(memberUserIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid member user ID")
		return
	}

	err = h.boardService.RemoveMemberFromBoard(uint(boardID), uint(memberUserID), userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusOK, "Member removed from board successfully", nil)
}

func (h *BoardHandler) GetBoardMembers(c *gin.Context) {
	userID, _ := c.Get("userID")
	boardIDStr := c.Param("boardID")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid board ID")
		return
	}

	members, err := h.boardService.GetBoardMembers(uint(boardID), userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}

	var memberResponses []dto.BoardMemberResponse // Use dto type
	for _, m := range members {
		memberResponses = append(memberResponses, dto.MapBoardMemberToResponse(&m)) // Use dto mapper
	}
	RespondWithSuccess(c, http.StatusOK, "Board members retrieved successfully", memberResponses)
}

// Mapping functions MapBoardToResponse and MapBoardMemberToResponse are now in dto/board_dto.go
// MapUserToResponse is in dto/auth_dto.go (which will be imported as part of dto package)
// MapListToResponse will be in dto/list_dto.go
