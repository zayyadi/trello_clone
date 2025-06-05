package services

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrBoardNotFound       = errors.New("board not found")
	ErrListNotFound        = errors.New("list not found")
	ErrCardNotFound        = errors.New("card not found")
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrForbidden           = errors.New("forbidden: insufficient permissions")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrEmailExists         = errors.New("email already exists")
	ErrUsernameExists      = errors.New("username already exists")
	ErrUserAlreadyMember   = errors.New("user is already a member of this board")
	ErrBoardMemberNotFound = errors.New("board member not found")
	ErrCannotRemoveOwner   = errors.New("cannot remove the board owner")
	ErrInvalidInput        = errors.New("invalid input")
	ErrSameListMove        = errors.New("card is already in the target list; use reorder instead")
	ErrPositionOutOfBound  = errors.New("position out of bounds")
	ErrUserNotCollaborator = errors.New("user is not a collaborator on this card")
)
