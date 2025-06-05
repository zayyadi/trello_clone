import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { apiClient } from '../../lib/apiClient'; // Import apiClient (named export)

// --- Thunks for Board List ---
export const fetchUserBoards = createAsyncThunk(
  'boards/fetchUserBoards',
  async (_, { rejectWithValue }) => {
    try {
      const response = await apiClient.get('/boards');
      return response.data.data;
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not fetch boards');
    }
  }
);

export const createNewBoard = createAsyncThunk(
  'boards/createNewBoard',
  async (boardData, { rejectWithValue }) => {
    try {
      const response = await apiClient.post('/boards', boardData);
      return response.data.data;
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not create board');
    }
  }
);

export const updateBoard = createAsyncThunk(
  'boards/updateBoard',
  async ({ boardId, boardData }, { rejectWithValue }) => {
    try {
      const response = await apiClient.put(`/boards/${boardId}`, boardData);
      return response.data.data;
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not update board');
    }
  }
);

export const deleteBoard = createAsyncThunk(
  'boards/deleteBoard',
  async (boardId, { rejectWithValue }) => {
    try {
      await apiClient.delete(`/boards/${boardId}`);
      return boardId; // Return the ID of the deleted board
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not delete board');
    }
  }
);

// --- Thunks for Specific Board, its Lists, and Cards ---
export const fetchBoardDetails = createAsyncThunk(
  'boards/fetchBoardDetails',
  async (boardId, { rejectWithValue }) => {
    try {
      // Backend's GET /boards/:boardID should return board with its lists and cards preloaded
      // If not, you'd make separate calls for lists and then cards for each list.
      // Assuming backend returns: { id, name, description, owner, lists: [{id, name, cards: [{id, title, ...}]}] }
      const response = await apiClient.get(`/boards/${boardId}`);
      return response.data.data; // This should be the full board object
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not fetch board details');
    }
  }
);

export const updateList = createAsyncThunk(
  'boards/updateList',
  async ({ listId, listData }, { rejectWithValue }) => {
    try {
      const response = await apiClient.put(`/lists/${listId}`, listData);
      return response.data.data; // Should be the updated list
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not update list');
    }
  }
);

export const deleteList = createAsyncThunk(
  'boards/deleteList',
  async (listId, { rejectWithValue }) => {
    try {
      await apiClient.delete(`/lists/${listId}`);
      return listId; // Return the ID of the deleted list
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not delete list');
    }
  }
);

export const addListToBoard = createAsyncThunk(
  'boards/addListToBoard',
  async ({ boardId, name }, { rejectWithValue }) => {
    try {
      const response = await apiClient.post(`/boards/${boardId}/lists`, { name });
      return response.data.data; // Should be the newly created list
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not add list');
    }
  }
);

export const addCardToList = createAsyncThunk(
  'boards/addCardToList',
  async ({ listId, title, description, dueDate, assignedUserID, supervisorID }, { rejectWithValue }) => { // Added new fields
    try {
      const response = await apiClient.post(`/lists/${listId}/cards`, { 
          title, 
          description, // Optional
          dueDate, // Optional
          assignedUserID, // Optional
          supervisorID // Optional
        });
      return response.data.data;
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not add card');
    }
  }
);

export const updateCard = createAsyncThunk(
  'boards/updateCard',
  async ({ cardId, cardData }, { rejectWithValue }) => {
    try {
      const response = await apiClient.put(`/cards/${cardId}`, cardData);
      return response.data.data; // Should be the updated card
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not update card');
    }
  }
);

export const updateCardDetails = createAsyncThunk(
  'boards/updateCardDetails',
  async (cardData, { rejectWithValue }) => { // cardData: { cardId, title, description, dueDate, status, assignedUserID, supervisorID }
    const { cardId, ...updatePayload } = cardData;
    try {
      const response = await apiClient.put(`/cards/${cardId}`, updatePayload);
      return response.data.data; // Updated card
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not update card');
    }
  }
);

export const deleteCard = createAsyncThunk(
  'boards/deleteCard',
  async (cardId, { rejectWithValue }) => {
    try {
      await apiClient.delete(`/cards/${cardId}`);
      return cardId; // Return the ID of the deleted card
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not delete card');
    }
  }
);

export const moveCard = createAsyncThunk(
  'boards/moveCard',
  async ({ cardId, targetListId, newPosition }, { rejectWithValue }) => {
    try {
      // Note: Backend expects 1-based indexing for position
      const response = await apiClient.patch(`/cards/${cardId}/move`, { targetListId, newPosition });
      return response.data.data; // Updated card
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not move card');
    }
  }
);

// Thunk for updating card order within the same list (if backend supports direct position update)
export const updateCardOrderInList = createAsyncThunk(
  'boards/updateCardOrderInList',
  async ({ cardId, newPosition, listId /* for optimistic update reference */ }, { rejectWithValue }) => {
    try {
      // Assuming PATCH /cards/:cardID can update position if targetListID is not changed
      // Or use the existing move endpoint if it handles same-list moves gracefully.
      // For now, using the general move endpoint but ensuring targetListId is the current list.
      const response = await apiClient.put(`/cards/${cardId}`, { position: newPosition });
      // The backend's PUT /cards/:cardID should handle reordering within its list.
      // It might need more info, or the backend needs to be smart.
      // A better backend endpoint might be PATCH /lists/:listId/cards/reorder { cardOrders: [{cardId, newPosition}]}
      // For now, this is a simplified call. The backend needs to handle this specific PUT for position.
      return { ...response.data.data, originalListId: listId }; // Return updated card
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not reorder card');
    }
  }
);

export const addCommentToCard = createAsyncThunk(
  'boards/addCommentToCard',
  async ({ cardId, text }, { rejectWithValue }) => {
    try {
      const response = await apiClient.post(`/cards/${cardId}/comments`, { text });
      return { cardId, comment: response.data.data }; // Return cardId to find the card in state
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not add comment');
    }
  }
);

export const deleteCommentFromCard = createAsyncThunk(
  'boards/deleteCommentFromCard',
  async ({ cardId, commentId }, { rejectWithValue }) => {
    try {
      await apiClient.delete(`/cards/${cardId}/comments/${commentId}`);
      return { cardId, commentId }; // Return IDs to remove comment from state
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not delete comment');
    }
  }
);


const initialState = {
  userBoards: [], // List of all boards for the dashboard
  userBoardsStatus: 'idle',
  userBoardsError: null,

  currentBoard: null, // The board currently being viewed { id, name, lists: [...] }
  currentBoardStatus: 'idle',
  currentBoardError: null,

  // Status for list/card operations
  listCardOpStatus: 'idle',
  listCardOpError: null,
};

const boardsSlice = createSlice({
  name: 'boards',
  initialState,
  reducers: {
    // Optimistic updates for card movements
    optimisticallyUpdateCardOrder: (state, action) => {
      const { listId, orderedCards } = action.payload;
      if (state.currentBoard && state.currentBoard.lists) {
        const listIndex = state.currentBoard.lists.findIndex(l => l.id === listId);
        if (listIndex !== -1) {
          state.currentBoard.lists[listIndex].cards = orderedCards;
        }
      }
    },
    optimisticallyMoveCardBetweenLists: (state, action) => {
      const { sourceListId, destinationListId, sourceCards, destinationCards } = action.payload;
      if (state.currentBoard && state.currentBoard.lists) {
        const sourceListIndex = state.currentBoard.lists.findIndex(l => l.id === sourceListId);
        const destListIndex = state.currentBoard.lists.findIndex(l => l.id === destinationListId);
        if (sourceListIndex !== -1) {
          state.currentBoard.lists[sourceListIndex].cards = sourceCards;
        }
        if (destListIndex !== -1) {
          state.currentBoard.lists[destListIndex].cards = destinationCards;
        }
      }
    },
    clearCurrentBoard: (state) => {
        state.currentBoard = null;
        state.currentBoardStatus = 'idle';
        state.currentBoardError = null;
    }
  },
  extraReducers: (builder) => {
    builder
      // User Boards List
      .addCase(fetchUserBoards.pending, (state) => {
        state.userBoardsStatus = 'loading';
      })
      .addCase(fetchUserBoards.fulfilled, (state, action) => {
        state.userBoardsStatus = 'succeeded';
        state.userBoards = action.payload;
      })
      .addCase(fetchUserBoards.rejected, (state, action) => {
        state.userBoardsStatus = 'failed';
        state.userBoardsError = action.payload;
      })
      .addCase(createNewBoard.fulfilled, (state, action) => {
        state.userBoards.push(action.payload); // Add to the list
      })
      .addCase(updateBoard.fulfilled, (state, action) => {
        const updatedBoard = action.payload;
        // Update in userBoards list
        const boardIndex = state.userBoards.findIndex(board => board.id === updatedBoard.id);
        if (boardIndex !== -1) {
          state.userBoards[boardIndex] = updatedBoard;
        }
        // Update if it's the current board
        if (state.currentBoard && state.currentBoard.id === updatedBoard.id) {
          state.currentBoard = { ...state.currentBoard, ...updatedBoard };
        }
      })
      .addCase(deleteBoard.fulfilled, (state, action) => {
        const deletedBoardId = action.payload;
        state.userBoards = state.userBoards.filter(board => board.id !== deletedBoardId);
        if (state.currentBoard && state.currentBoard.id === deletedBoardId) {
          state.currentBoard = null; // Clear current board if deleted
        }
      })
      // Current Board Details
      .addCase(fetchBoardDetails.pending, (state) => {
        state.currentBoardStatus = 'loading';
      })
      .addCase(fetchBoardDetails.fulfilled, (state, action) => {
        state.currentBoardStatus = 'succeeded';
        // Ensure lists and cards are sorted by position from backend
        // Create a deep copy of the payload to avoid direct mutation of action.payload
        const boardData = JSON.parse(JSON.stringify(action.payload));
        if (boardData.lists) {
            boardData.lists.forEach(list => {
                if (list.cards) {
                    list.cards.sort((a, b) => a.position - b.position);
                } else {
                    list.cards = []; // Ensure cards array exists
                }
            });
            boardData.lists.sort((a,b) => a.position - b.position);
        } else {
            boardData.lists = []; // Ensure lists array exists
        }
        state.currentBoard = boardData;
      })
      .addCase(fetchBoardDetails.rejected, (state, action) => {
        state.currentBoardStatus = 'failed';
        state.currentBoardError = action.payload;
      })
      // Add List
      .addCase(addListToBoard.pending, (state) => {
        state.listCardOpStatus = 'loading';
      })
      .addCase(addListToBoard.fulfilled, (state, action) => {
        state.listCardOpStatus = 'succeeded';
        if (state.currentBoard && state.currentBoard.lists) {
          state.currentBoard.lists.push({...action.payload, cards: []}); // Add new list with empty cards array
        } else if (state.currentBoard) {
            state.currentBoard.lists = [{...action.payload, cards: []}];
        }
      })
      .addCase(addListToBoard.rejected, (state, action) => {
        state.listCardOpStatus = 'failed';
        state.listCardOpError = action.payload;
      })
      .addCase(updateList.fulfilled, (state, action) => {
        state.listCardOpStatus = 'succeeded';
        const updatedList = action.payload;
        if (state.currentBoard && state.currentBoard.lists) {
          const listIndex = state.currentBoard.lists.findIndex(l => l.id === updatedList.id);
          if (listIndex !== -1) {
            // Preserve cards if not returned by update payload
            const existingCards = state.currentBoard.lists[listIndex].cards;
            state.currentBoard.lists[listIndex] = { ...updatedList, cards: existingCards || [] };
            // If position changed, re-sort lists
            state.currentBoard.lists.sort((a, b) => a.position - b.position);
          }
        }
      })
      .addCase(updateList.rejected, (state, action) => {
        state.listCardOpStatus = 'failed';
        state.listCardOpError = action.payload;
      })
      .addCase(deleteList.fulfilled, (state, action) => {
        state.listCardOpStatus = 'succeeded';
        const deletedListId = action.payload;
        if (state.currentBoard && state.currentBoard.lists) {
          state.currentBoard.lists = state.currentBoard.lists.filter(l => l.id !== deletedListId);
        }
      })
      .addCase(deleteList.rejected, (state, action) => {
        state.listCardOpStatus = 'failed';
        state.listCardOpError = action.payload;
      })
      // Add Card
      .addCase(addCardToList.pending, (state) => {
        state.listCardOpStatus = 'loading';
      })
      .addCase(addCardToList.fulfilled, (state, action) => {
        state.listCardOpStatus = 'succeeded';
        if (state.currentBoard && state.currentBoard.lists) {
          const listIndex = state.currentBoard.lists.findIndex(l => l.id === action.payload.listID);
          if (listIndex !== -1) {
            if (!state.currentBoard.lists[listIndex].cards) {
                state.currentBoard.lists[listIndex].cards = [];
            }
            state.currentBoard.lists[listIndex].cards.push(action.payload);
          }
        }
      })
      .addCase(addCardToList.rejected, (state, action) => {
        state.listCardOpStatus = 'failed';
        state.listCardOpError = action.payload;
      })
      .addCase(updateCardDetails.fulfilled, (state, action) => {
        state.listCardOpStatus = 'succeeded';
        const updatedCard = action.payload;
        if (state.currentBoard && state.currentBoard.lists) {
          const listIndex = state.currentBoard.lists.findIndex(l => l.id === updatedCard.listID);
          if (listIndex !== -1) {
            const cardIndex = state.currentBoard.lists[listIndex].cards.findIndex(c => c.id === updatedCard.id);
            if (cardIndex !== -1) {
              state.currentBoard.lists[listIndex].cards[cardIndex] = updatedCard;
            } else { // Should not happen if card exists
              state.currentBoard.lists[listIndex].cards.push(updatedCard);
            }
             // Optional: re-sort if position could have changed, though this thunk isn't for position
            // state.currentBoard.lists[listIndex].cards.sort((a, b) => a.position - b.position);
          }
        }
      })
      .addCase(updateCardDetails.pending, (state) => { state.listCardOpStatus = 'loading'; })
      .addCase(updateCardDetails.rejected, (state, action) => {
        state.listCardOpStatus = 'failed';
        state.listCardOpError = action.payload;
      })
      .addCase(deleteCard.fulfilled, (state, action) => {
        state.listCardOpStatus = 'succeeded';
        const deletedCardId = action.payload;
        if (state.currentBoard && state.currentBoard.lists) {
          state.currentBoard.lists.forEach(list => {
            list.cards = list.cards.filter(card => card.id !== deletedCardId);
          });
        }
      })
      .addCase(deleteCard.rejected, (state, action) => {
        state.listCardOpStatus = 'failed';
        state.listCardOpError = action.payload;
      })
      // Move Card (covers both reorder within list and move between lists via backend)
      // After backend confirms, we might refetch or rely on optimistic updates being correct
      // For simplicity, if backend returns the updated board structure or list, we update it.
      // Or, we can refetch the board details, but this is less efficient.
      // The moveCard thunk now returns just the card. The backend should have handled reordering.
      // We need to manually adjust the state or refetch.
      // Let's assume after a move, we might want to refetch the board for consistency or handle complex state updates.
      // For now, optimistic updates are primary, and this `fulfilled` case can be a no-op if optimistic is good,
      // or it can trigger a refetch of the board.
      .addCase(moveCard.fulfilled, (state, action) => {
        state.listCardOpStatus = 'succeeded';
        // The optimistic updates should have handled the UI.
        // If strict consistency is needed, refetch the board:
        // state.currentBoardStatus = 'idle'; // to trigger refetch on component
        // Or, more granularly update the specific card and lists if payload is rich enough
        const movedCard = action.payload;
        if (state.currentBoard && state.currentBoard.lists) {
            // Remove card from any old position (if it was already there due to optimistic)
            state.currentBoard.lists.forEach(list => {
                list.cards = list.cards.filter(c => c.id !== movedCard.id);
            });
            // Add card to its new list and position
            const targetList = state.currentBoard.lists.find(l => l.id === movedCard.listID);
            if (targetList) {
                targetList.cards.push(movedCard);
                targetList.cards.sort((a, b) => a.position - b.position);
            }
        }

      })
      .addCase(moveCard.rejected, (state, action) => {
          state.listCardOpStatus = 'failed';
          state.listCardOpError = action.payload;
          // Here you might want to revert optimistic updates if they were applied
          // This requires storing pre-optimistic state, which is complex.
          // Simplest: notify user and they might need to refresh or try again.
      })
      .addCase(updateCardOrderInList.fulfilled, (state, action) => {
          state.listCardOpStatus = 'succeeded';
          // Similar to moveCard, optimistic updates are primary.
          // This ensures the card's position is updated from backend response if different.
          const updatedCard = action.payload;
          const listId = action.payload.originalListId; // from thunk meta or modified payload
           if (state.currentBoard && state.currentBoard.lists) {
            const listIndex = state.currentBoard.lists.findIndex(l => l.id === listId);
            if (listIndex !== -1) {
                const cardIndex = state.currentBoard.lists[listIndex].cards.findIndex(c => c.id === updatedCard.id);
                if (cardIndex !== -1) {
                    state.currentBoard.lists[listIndex].cards[cardIndex] = updatedCard;
                } else {
                     state.currentBoard.lists[listIndex].cards.push(updatedCard); // If not found, add
                }
                state.currentBoard.lists[listIndex].cards.sort((a, b) => a.position - b.position);
            }
        }
      })
      .addCase(addCommentToCard.fulfilled, (state, action) => {
        state.listCardOpStatus = 'succeeded';
        const { cardId, comment } = action.payload;
        if (state.currentBoard && state.currentBoard.lists) {
          state.currentBoard.lists.forEach(list => {
            const card = list.cards.find(c => c.id === cardId);
            if (card) {
              if (!card.comments) {
                card.comments = [];
              }
              card.comments.push(comment);
            }
          });
        }
      })
      .addCase(addCommentToCard.rejected, (state, action) => {
        state.listCardOpStatus = 'failed';
        state.listCardOpError = action.payload;
      })
      .addCase(deleteCommentFromCard.fulfilled, (state, action) => {
        state.listCardOpStatus = 'succeeded';
        const { cardId, commentId } = action.payload;
        if (state.currentBoard && state.currentBoard.lists) {
          state.currentBoard.lists.forEach(list => {
            const card = list.cards.find(c => c.id === cardId);
            if (card && card.comments) {
              card.comments = card.comments.filter(c => c.id !== commentId);
            }
          });
        }
      })
      .addCase(deleteCommentFromCard.rejected, (state, action) => {
        state.listCardOpStatus = 'failed';
        state.listCardOpError = action.payload;
      });
  },
});

export const { optimisticallyUpdateCardOrder, optimisticallyMoveCardBetweenLists, clearCurrentBoard } = boardsSlice.actions;

export const selectUserBoards = (state) => state.boards.userBoards;
export const selectUserBoardsStatus = (state) => state.boards.userBoardsStatus;
export const selectCurrentBoard = (state) => state.boards.currentBoard;
export const selectCurrentBoardStatus = (state) => state.boards.currentBoardStatus;
export const selectCurrentBoardError = (state) => state.boards.currentBoardError;


export default boardsSlice.reducer;
