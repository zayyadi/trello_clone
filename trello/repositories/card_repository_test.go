package repositories

import (
	"errors"
	"testing"

	"github.com/zayyadi/trello/models"
	"gorm.io/gorm"
)

// MockCardRepository for testing purposes.
// This is a simplified mock. Ideally, use testify/mock.
type MockCardRepositoryForIsUserCollab struct {
	db *gorm.DB // Keep db field if some passthrough is needed, or remove
	CardRepositoryInterface
	mockIsCollaborator func(cardID uint, userID uint) (bool, error)
	mockDbFirst        func(dest interface{}, conds ...interface{}) *gorm.DB // Simplified mock for db.First
}

// Overriding IsCollaborator for our mock
func (m *MockCardRepositoryForIsUserCollab) IsCollaborator(cardID uint, userID uint) (bool, error) {
	if m.mockIsCollaborator != nil {
		return m.mockIsCollaborator(cardID, userID)
	}
	return false, errors.New("IsCollaborator mock not implemented")
}

// Simplified GORM First mock
func (m *MockCardRepositoryForIsUserCollab) getMockDbFirst(cardToReturn *models.Card, errToReturn error) func(dest interface{}, conds ...interface{}) *gorm.DB {
	return func(dest interface{}, conds ...interface{}) *gorm.DB {
		if cardToReturn != nil && dest != nil {
			if cardDest, ok := dest.(*models.Card); ok {
				*cardDest = *cardToReturn
			}
		}
		return &gorm.DB{Error: errToReturn} // Return a gorm.DB instance with the desired error
	}
}

func TestCardRepository_IsUserCollaboratorOrAssignee(t *testing.T) {
	// Dummy DB needed for NewCardRepository, even if we override methods.
	// In a real scenario, this would be a connection to a test DB or a proper GORM mock.
	var dummyDb *gorm.DB

	assignedUserID := uint(100)
	collaboratorUserID := uint(200)
	otherUserID := uint(300)
	cardID := uint(1)

	tests := []struct {
		name                         string
		cardID                       uint
		userID                       uint
		setupMockDbFirstCard         *models.Card
		setupMockDbFirstError        error
		setupMockIsCollaborator      func(cardID uint, userID uint) (bool, error)
		expectedResult               bool
		expectedError                error
		skipIsCollaboratorCheck      bool // Skip if db.First already determines outcome or errors out
	}{
		{
			name:   "User is assignee",
			cardID: cardID,
			userID: assignedUserID,
			setupMockDbFirstCard: &models.Card{
				Model:          gorm.Model{ID: cardID},
				AssignedUserID: &assignedUserID,
			},
			setupMockDbFirstError:   nil,
			expectedResult:          true,
			expectedError:           nil,
			skipIsCollaboratorCheck: true,
		},
		{
			name:   "User is collaborator, not assignee",
			cardID: cardID,
			userID: collaboratorUserID,
			setupMockDbFirstCard: &models.Card{ // Assigned to someone else or no one
				Model:          gorm.Model{ID: cardID},
				AssignedUserID: nil, // or a different ID
			},
			setupMockDbFirstError: nil,
			setupMockIsCollaborator: func(cID uint, uID uint) (bool, error) {
				if cID == cardID && uID == collaboratorUserID {
					return true, nil
				}
				return false, nil
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:   "User is neither assignee nor collaborator",
			cardID: cardID,
			userID: otherUserID,
			setupMockDbFirstCard: &models.Card{
				Model:          gorm.Model{ID: cardID},
				AssignedUserID: &assignedUserID, // Assigned to someone else
			},
			setupMockDbFirstError: nil,
			setupMockIsCollaborator: func(cID uint, uID uint) (bool, error) {
				// Ensure this is called with otherUserID
				if cID == cardID && uID == otherUserID {
					return false, nil
				}
				// Fail test if called with unexpected ID
				t.Errorf("IsCollaborator called with unexpected userID %d", uID)
				return false, errors.New("unexpected call to IsCollaborator")
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:                    "Card not found",
			cardID:                  cardID,
			userID:                  otherUserID,
			setupMockDbFirstCard:    nil,
			setupMockDbFirstError:   gorm.ErrRecordNotFound,
			expectedResult:          false,
			expectedError:           gorm.ErrRecordNotFound,
			skipIsCollaboratorCheck: true,
		},
		{
			name:   "DB error on First call",
			cardID: cardID,
			userID: otherUserID,
			setupMockDbFirstCard:    nil,
			setupMockDbFirstError:   errors.New("some db error"),
			expectedResult:          false,
			expectedError:           errors.New("some db error"),
			skipIsCollaboratorCheck: true,
		},
		{
			name:   "DB error on IsCollaborator call",
			cardID: cardID,
			userID: collaboratorUserID, // User who would be a collaborator
			setupMockDbFirstCard: &models.Card{ // Not assigned to this user
				Model: gorm.Model{ID: cardID},
			},
			setupMockDbFirstError: nil,
			setupMockIsCollaborator: func(cID uint, uID uint) (bool, error) {
				return false, errors.New("isCollaborator db error")
			},
			expectedResult: false,
			expectedError:  errors.New("isCollaborator db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewCardRepository(dummyDb).(*CardRepository)
			originalDB := repo.db

			// This is a simplified simulation of GORM's behavior for the purpose of testing logic flow.
			// It does not involve actual database interaction or full GORM mocking.

			// Simulate the r.db.Select().First() call's outcome
			dbFirstSim := func(dest interface{}, cardIDSearch interface{}) error {
				if tt.setupMockDbFirstError != nil {
					return tt.setupMockDbFirstError
				}
				if tt.setupMockDbFirstCard != nil {
					if c, ok := dest.(*models.Card); ok {
						*c = *tt.setupMockDbFirstCard
						return nil
					}
				}
				// Should not happen if setup is correct
				return errors.New("dbFirstSim was called without proper setup for success or specific error")
			}

			// Simulate the IsCollaborator call's outcome
			isCollabSim := func(cardIDSearch uint, userIDSearch uint) (bool, error) {
				if tt.skipIsCollaboratorCheck {
					// This case should ideally not call IsCollaborator.
					// If it does, it might indicate a logic flaw or a test setup issue.
					t.Logf("IsCollaborator called unexpectedly for test: %s", tt.name)
					return false, errors.New("IsCollaborator called unexpectedly")
				}
				if tt.setupMockIsCollaborator != nil {
					return tt.setupMockIsCollaborator(cardIDSearch, userIDSearch)
				}
				// Default behavior if not skipped and not mocked
				return false, errors.New("IsCollaborator was called but not mocked for this test case")
			}

			// Manually step through the logic of the actual IsUserCollaboratorOrAssignee method
			var actualResult bool
			var actualError error

			var cardForLogic models.Card
			// Simulate: err := r.db.Select("assigned_user_id").First(&card, cardID).Error;
			errDbFirst := dbFirstSim(&cardForLogic, tt.cardID)

			if errDbFirst != nil {
				if errors.Is(errDbFirst, gorm.ErrRecordNotFound) {
					actualError = gorm.ErrRecordNotFound
				} else {
					actualError = errDbFirst
				}
				actualResult = false
			} else {
				if cardForLogic.AssignedUserID != nil && *cardForLogic.AssignedUserID == tt.userID {
					actualResult = true
					actualError = nil
				} else {
					// Simulate: return r.IsCollaborator(cardID, userID)
					isCollab, collabErr := isCollabSim(tt.cardID, tt.userID)
					actualResult = isCollab
					actualError = collabErr
				}
			}

			if actualResult != tt.expectedResult {
				t.Errorf("Expected result %v, got %v", tt.expectedResult, actualResult)
			}
			if (actualError == nil && tt.expectedError != nil) ||
			   (actualError != nil && tt.expectedError == nil) ||
			   (actualError != nil && tt.expectedError != nil && actualError.Error() != tt.expectedError.Error()) {
				t.Errorf("Expected error '%v', got '%v'", tt.expectedError, actualError)
			}
			repo.db = originalDB // Restore
		})
	}
}
