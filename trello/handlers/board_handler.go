package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/services"
)

type BoardHandler struct {
	boardService *services.BoardService
}

func NewBoardHandler(boardService *services.BoardService) *BoardHandler {
	return &BoardHandler{boardService: boardService}
}

func (h *BoardHandler) CreateBoard(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}

	var req CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	board, err := h.boardService.CreateBoard(req.Name, req.Description, userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusCreated, "Board created successfully", MapBoardToResponse(board, true, false)) // Include owner, not lists/members by default
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
	RespondWithSuccess(c, http.StatusOK, "Board retrieved successfully", MapBoardToResponse(fullBoard, true, true))
}

func (h *BoardHandler) GetBoardsForUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	boards, err := h.boardService.GetBoardsForUser(userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	var boardResponses []BoardResponse
	for _, b := range boards {
		boardResponses = append(boardResponses, MapBoardToResponse(&b, true, false)) // Include owner, not lists/members
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

	var req UpdateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	board, err := h.boardService.UpdateBoard(uint(boardID), req.Name, req.Description, userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusOK, "Board updated successfully", MapBoardToResponse(board, true, false))
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

	var req AddMemberRequest
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
	RespondWithSuccess(c, http.StatusCreated, "Member added to board successfully", MapBoardMemberToResponse(boardMember))
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

	var memberResponses []BoardMemberResponse
	for _, m := range members {
		memberResponses = append(memberResponses, MapBoardMemberToResponse(&m))
	}
	RespondWithSuccess(c, http.StatusOK, "Board members retrieved successfully", memberResponses)
}

// Helper to map model.Board to BoardResponse
func MapBoardToResponse(board *models.Board, includeOwner, includeMembersAndLists bool) BoardResponse {
	if board == nil {
		return BoardResponse{}
	}
	resp := BoardResponse{
		ID:          board.Model.ID,
		Name:        board.Name,
		Description: board.Description,
		OwnerID:     board.OwnerID,
		CreatedAt:   board.Model.CreatedAt,
		UpdatedAt:   board.Model.UpdatedAt,
	}
	if includeOwner && board.Owner.ID != 0 {
		resp.Owner = MapUserToResponse(&board.Owner)
	}
	if includeMembersAndLists {
		if len(board.Members) > 0 {
			resp.Members = []BoardMemberResponse{}
			for _, m := range board.Members {
				resp.Members = append(resp.Members, MapBoardMemberToResponse(&m))
			}
		}
		if len(board.Lists) > 0 {
			resp.Lists = []ListResponse{}
			for _, l := range board.Lists {
				// For lists within board response, don't include cards by default to avoid huge payloads
				resp.Lists = append(resp.Lists, MapListToResponse(&l, false))
			}
		}
	}
	return resp
}

// Helper to map model.BoardMember to BoardMemberResponse
func MapBoardMemberToResponse(member *models.BoardMember) BoardMemberResponse {
	if member == nil {
		return BoardMemberResponse{}
	}
	return BoardMemberResponse{
		BoardID:   member.BoardID,
		UserID:    member.UserID,
		User:      MapUserToResponse(&member.User),
		CreatedAt: member.CreatedAt,
	}
}
