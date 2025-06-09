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
      const response = await apiClient.get(`/boards/${boardId}`);
      return response.data.data;
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
      return response.data.data;
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
      return listId;
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
      return response.data.data;
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not add list');
    }
  }
);

export const addCardToList = createAsyncThunk(
  'boards/addCardToList',
  async ({ listId, title, description, dueDate, assignedUserID, supervisorID }, { rejectWithValue }) => {
    const payload = { title, description, dueDate, assignedUserID, supervisorID };
    console.log('Adding card with payload:', payload);
    try {
      const response = await apiClient.post(`/lists/${listId}/cards`, payload);
      console.log('Add card response:', response.data);
      return response.data.data;
    } catch (error) {
      console.error('Add card error:', error.response?.data || error.message);
      return rejectWithValue(error.response?.data?.message || 'Could not add card');
    }
  }
);

export const updateCard = createAsyncThunk(
  'boards/updateCard',
  async ({ cardId, cardData }, { rejectWithValue }) => {
    try {
      const response = await apiClient.put(`/cards/${cardId}`, cardData);
      return response.data.data;
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not update card');
    }
  }
);

export const updateCardDetails = createAsyncThunk(
  'boards/updateCardDetails',
  async (cardData, { rejectWithValue }) => {
    const { cardId, ...updatePayload } = cardData;
    console.log('Updating card with payload:', updatePayload);
    try {
      const response = await apiClient.put(`/cards/${cardId}`, updatePayload);
      console.log('Update card response:', response.data);
      return response.data.data;
    } catch (error) {
      console.error('Update card error:', error.response?.data || error.message);
      return rejectWithValue(error.response?.data?.message || 'Could not update card');
    }
  }
);

export const deleteCard = createAsyncThunk(
  'boards/deleteCard',
  async (cardId, { rejectWithValue }) => {
    try {
      await apiClient.delete(`/cards/${cardId}`);
      return cardId;
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not delete card');
    }
  }
);

export const moveCard = createAsyncThunk(
  'boards/moveCard',
  async ({ cardId, targetListId, newPosition }, { rejectWithValue }) => {
    try {
      const response = await apiClient.patch(`/cards/${cardId}/move`, { targetListId, newPosition });
      return response.data.data;
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not move card');
    }
  }
);

export const updateCardOrderInList = createAsyncThunk(
  'boards/updateCardOrderInList',
  async ({ cardId, newPosition, listId }, { rejectWithValue }) => {
    try {
      const response = await apiClient.put(`/cards/${cardId}`, { position: newPosition });
      return { ...response.data.data, originalListId: listId };
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
      return { cardId, comment: response.data.data };
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
      return { cardId, commentId };
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not delete comment');
    }
  }
);

// --- Card Collaborator Thunks ---
export const fetchCardCollaborators = createAsyncThunk(
  'boards/fetchCardCollaborators',
  async (cardId, { rejectWithValue }) => {
    try {
      const response = await apiClient.get(`/cards/${cardId}/collaborators`);
      return { cardId, collaborators: response.data.data };
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not fetch collaborators');
    }
  }
);

export const addCardCollaborator = createAsyncThunk(
  'boards/addCardCollaborator',
  async ({ cardId, userIdentifier }, { rejectWithValue }) => {
    try {
      const response = await apiClient.post(`/cards/${cardId}/collaborators`, userIdentifier);
      return { cardId, collaborator: response.data.data };
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not add collaborator');
    }
  }
);

export const removeCardCollaborator = createAsyncThunk(
  'boards/removeCardCollaborator',
  async ({ cardId, userIdToRemove }, { rejectWithValue }) => {
    try {
      await apiClient.delete(`/cards/${cardId}/collaborators/${userIdToRemove}`);
      return { cardId, userIdRemoved: userIdToRemove };
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Could not remove collaborator');
    }
  }
);


const initialState = {
  userBoards: [],
  userBoardsStatus: 'idle',
  userBoardsError: null,

  currentBoard: null,
  currentBoardStatus: 'idle',
  currentBoardError: null,

  listCardOpStatus: 'idle',
  listCardOpError: null,
};

const boardsSlice = createSlice({
  name: 'boards',
  initialState,
  reducers: {
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
    },
    // Add the WebSocket event reducers here
    boardUpdatedByWS: (state, action) => {
      const updatedBoardData = action.payload;
      if (state.currentBoard && state.currentBoard.id === updatedBoardData.id) {
        state.currentBoard.name = updatedBoardData.name;
        state.currentBoard.description = updatedBoardData.description;
      }
      const boardIndex = state.userBoards.findIndex(b => b.id === updatedBoardData.id);
      if (boardIndex !== -1) {
        state.userBoards[boardIndex] = { ...state.userBoards[boardIndex], ...updatedBoardData };
      }
    },
    boardDeletedByWS: (state, action) => {
      const { id: deletedBoardId } = action.payload;
      state.userBoards = state.userBoards.filter(board => board.id !== deletedBoardId);
      if (state.currentBoard && state.currentBoard.id === deletedBoardId) {
        state.currentBoard = null;
        state.currentBoardStatus = 'idle';
      }
    },
    listCreatedByWS: (state, action) => {
      const newList = action.payload;
      if (state.currentBoard && state.currentBoard.id === newList.boardID) {
        if (!state.currentBoard.lists) state.currentBoard.lists = [];
        if (!state.currentBoard.lists.find(l => l.id === newList.id)) {
          state.currentBoard.lists.push({ ...newList, cards: newList.cards || [] });
          state.currentBoard.lists.sort((a, b) => a.position - b.position);
        }
      }
    },
    listUpdatedByWS: (state, action) => {
      const updatedList = action.payload;
      if (state.currentBoard && state.currentBoard.id === updatedList.boardID) {
        const listIndex = state.currentBoard.lists.findIndex(l => l.id === updatedList.id);
        if (listIndex !== -1) {
          const existingCards = state.currentBoard.lists[listIndex].cards;
          state.currentBoard.lists[listIndex] = { ...state.currentBoard.lists[listIndex], ...updatedList, cards: existingCards || [] };
          state.currentBoard.lists.sort((a, b) => a.position - b.position);
        }
      }
    },
    listDeletedByWS: (state, action) => {
      const { id: deletedListId, boardID } = action.payload;
      if (state.currentBoard && state.currentBoard.id === boardID) {
        state.currentBoard.lists = state.currentBoard.lists.filter(l => l.id !== deletedListId);
      }
    },
    cardCreatedByWS: (state, action) => {
      const newCard = action.payload;
      if (state.currentBoard) {
        const list = state.currentBoard.lists.find(l => l.id === newCard.listID);
        if (list) {
          if (!list.cards) list.cards = [];
          if (!list.cards.find(c => c.id === newCard.id)) {
            list.cards.push(newCard);
            list.cards.sort((a, b) => a.position - b.position);
          }
        }
      }
    },
    cardUpdatedByWS: (state, action) => {
      const updatedCard = action.payload;
      if (state.currentBoard) {
        const list = state.currentBoard.lists.find(l => l.id === updatedCard.listID);
        if (list) {
          const cardIndex = list.cards.findIndex(c => c.id === updatedCard.id);
          if (cardIndex !== -1) {
            list.cards[cardIndex] = { ...list.cards[cardIndex], ...updatedCard };
          } else {
            list.cards.push(updatedCard); // If card moved list and updated
            list.cards.sort((a, b) => a.position - b.position);
          }
        }
      }
    },
    cardDeletedByWS: (state, action) => {
      const { id: deletedCardId, listID, boardID } = action.payload;
      if (state.currentBoard && state.currentBoard.id === boardID) {
        const list = state.currentBoard.lists.find(l => l.id === listID);
        if (list) {
          list.cards = list.cards.filter(c => c.id !== deletedCardId);
        }
      }
    },
    cardMovedByWS: (state, action) => {
      const { cardId, oldListId, newListId, oldPosition, newPosition, boardId, updatedCard } = action.payload;
      if (state.currentBoard && state.currentBoard.id === boardId) {
        const oldList = state.currentBoard.lists.find(l => l.id === oldListId);
        if (oldList) {
          oldList.cards = oldList.cards.filter(c => c.id !== cardId);
        }
        const newList = state.currentBoard.lists.find(l => l.id === newListId);
        if (newList) {
          if (!newList.cards.find(c => c.id === cardId)) {
            newList.cards.push(updatedCard);
          } else {
            const cardIndex = newList.cards.findIndex(c => c.id === cardId);
            newList.cards[cardIndex] = updatedCard;
          }
          newList.cards.sort((a, b) => a.position - b.position);
        }
      }
    },
    // TODO: BOARD_MEMBER_ADDED, BOARD_MEMBER_REMOVED
    // TODO: CARD_COLLABORATOR_ADDED, CARD_COLLABORATOR_REMOVED
  },
  extraReducers: (builder) => {
    builder
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
        state.userBoards.push(action.payload);
      })
      .addCase(updateBoard.fulfilled, (state, action) => {
        const updatedBoard = action.payload;
        const boardIndex = state.userBoards.findIndex(board => board.id === updatedBoard.id);
        if (boardIndex !== -1) {
          state.userBoards[boardIndex] = updatedBoard;
        }
        if (state.currentBoard && state.currentBoard.id === updatedBoard.id) {
          state.currentBoard = { ...state.currentBoard, ...updatedBoard };
        }
      })
      .addCase(deleteBoard.fulfilled, (state, action) => {
        const deletedBoardId = action.payload;
        state.userBoards = state.userBoards.filter(board => board.id !== deletedBoardId);
        if (state.currentBoard && state.currentBoard.id === deletedBoardId) {
          state.currentBoard = null;
        }
      })
      .addCase(fetchBoardDetails.pending, (state) => {
        state.currentBoardStatus = 'loading';
      })
      .addCase(fetchBoardDetails.fulfilled, (state, action) => {
        state.currentBoardStatus = 'succeeded';
        const boardData = JSON.parse(JSON.stringify(action.payload));
        if (boardData.lists) {
            boardData.lists.forEach(list => {
                if (list.cards) {
                    list.cards.sort((a, b) => a.position - b.position);
                } else {
                    list.cards = [];
                }
            });
            boardData.lists.sort((a,b) => a.position - b.position);
        } else {
            boardData.lists = [];
        }
        state.currentBoard = boardData;
      })
      .addCase(fetchBoardDetails.rejected, (state, action) => {
        state.currentBoardStatus = 'failed';
        state.currentBoardError = action.payload;
      })
      .addCase(addListToBoard.pending, (state) => {
        state.listCardOpStatus = 'loading';
      })
      .addCase(addListToBoard.fulfilled, (state, action) => {
        state.listCardOpStatus = 'succeeded';
        if (state.currentBoard && state.currentBoard.lists) {
          state.currentBoard.lists.push({...action.payload, cards: []});
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
            const existingCards = state.currentBoard.lists[listIndex].cards;
            state.currentBoard.lists[listIndex] = { ...updatedList, cards: existingCards || [] };
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
            } else {
              state.currentBoard.lists[listIndex].cards.push(updatedCard);
            }
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
      .addCase(moveCard.fulfilled, (state, action) => {
        state.listCardOpStatus = 'succeeded';
        const movedCard = action.payload;
        if (state.currentBoard && state.currentBoard.lists) {
            state.currentBoard.lists.forEach(list => {
                list.cards = list.cards.filter(c => c.id !== movedCard.id);
            });
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
      })
      .addCase(updateCardOrderInList.fulfilled, (state, action) => {
          state.listCardOpStatus = 'succeeded';
          const updatedCard = action.payload;
          const listId = action.payload.originalListId;
           if (state.currentBoard && state.currentBoard.lists) {
            const listIndex = state.currentBoard.lists.findIndex(l => l.id === listId);
            if (listIndex !== -1) {
                const cardIndex = state.currentBoard.lists[listIndex].cards.findIndex(c => c.id === updatedCard.id);
                if (cardIndex !== -1) {
                    state.currentBoard.lists[listIndex].cards[cardIndex] = updatedCard;
                } else {
                     state.currentBoard.lists[listIndex].cards.push(updatedCard);
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
      })
      // Card Collaborators
      .addCase(fetchCardCollaborators.pending, (state) => {
        state.listCardOpStatus = 'loading';
      })
      .addCase(fetchCardCollaborators.fulfilled, (state, action) => {
        state.listCardOpStatus = 'succeeded';
        const { cardId, collaborators } = action.payload;
        if (state.currentBoard && state.currentBoard.lists) {
          for (const list of state.currentBoard.lists) {
            const cardIndex = list.cards.findIndex(c => c.id === cardId);
            if (cardIndex !== -1) {
              list.cards[cardIndex].collaborators = collaborators;
              break;
            }
          }
        }
      })
      .addCase(fetchCardCollaborators.rejected, (state, action) => {
        state.listCardOpStatus = 'failed';
        state.listCardOpError = action.payload;
      })
      .addCase(addCardCollaborator.pending, (state) => {
        state.listCardOpStatus = 'loading';
      })
      .addCase(addCardCollaborator.fulfilled, (state, action) => {
        state.listCardOpStatus = 'succeeded';
        const { cardId, collaborator } = action.payload;
        if (state.currentBoard && state.currentBoard.lists) {
          for (const list of state.currentBoard.lists) {
            const card = list.cards.find(c => c.id === cardId);
            if (card) {
              if (!card.collaborators) {
                card.collaborators = [];
              }
              // Avoid adding if already present by ID
              if (!card.collaborators.find(c => c.id === collaborator.id)) {
                card.collaborators.push(collaborator);
              }
              break;
            }
          }
        }
      })
      .addCase(addCardCollaborator.rejected, (state, action) => {
        state.listCardOpStatus = 'failed';
        state.listCardOpError = action.payload;
      })
      .addCase(removeCardCollaborator.pending, (state) => {
        state.listCardOpStatus = 'loading';
      })
      .addCase(removeCardCollaborator.fulfilled, (state, action) => {
        state.listCardOpStatus = 'succeeded';
        const { cardId, userIdRemoved } = action.payload;
        if (state.currentBoard && state.currentBoard.lists) {
          for (const list of state.currentBoard.lists) {
            const card = list.cards.find(c => c.id === cardId);
            if (card && card.collaborators) {
              card.collaborators = card.collaborators.filter(c => c.id !== userIdRemoved);
              break;
            }
          }
        }
      })
      .addCase(removeCardCollaborator.rejected, (state, action) => {
        state.listCardOpStatus = 'failed';
        state.listCardOpError = action.payload;
      });
  },
});

// New reducers for WebSocket events
const webSocketReducers = {
  boardUpdatedByWS: (state, action) => {
    const updatedBoardData = action.payload;
    if (state.currentBoard && state.currentBoard.id === updatedBoardData.id) {
      state.currentBoard.name = updatedBoardData.name;
      state.currentBoard.description = updatedBoardData.description;
      // Potentially other board-level fields
    }
    // Also update in userBoards list if present
    const boardIndex = state.userBoards.findIndex(b => b.id === updatedBoardData.id);
    if (boardIndex !== -1) {
      state.userBoards[boardIndex] = { ...state.userBoards[boardIndex], ...updatedBoardData };
    }
  },
  boardDeletedByWS: (state, action) => {
    const { id: deletedBoardId } = action.payload; // payload is {id: boardID}
    state.userBoards = state.userBoards.filter(board => board.id !== deletedBoardId);
    if (state.currentBoard && state.currentBoard.id === deletedBoardId) {
      state.currentBoard = null; // Or redirect, or show a "board deleted" message
      state.currentBoardStatus = 'idle';
    }
  },
  listCreatedByWS: (state, action) => {
    const newList = action.payload;
    if (state.currentBoard && state.currentBoard.id === newList.boardID) {
      if (!state.currentBoard.lists) state.currentBoard.lists = [];
      // Check if list already exists to prevent duplicates from optimistic + WS update
      if (!state.currentBoard.lists.find(l => l.id === newList.id)) {
        state.currentBoard.lists.push({ ...newList, cards: newList.cards || [] });
        state.currentBoard.lists.sort((a, b) => a.position - b.position);
      }
    }
  },
  listUpdatedByWS: (state, action) => {
    const updatedList = action.payload;
    if (state.currentBoard && state.currentBoard.id === updatedList.boardID) {
      const listIndex = state.currentBoard.lists.findIndex(l => l.id === updatedList.id);
      if (listIndex !== -1) {
        const existingCards = state.currentBoard.lists[listIndex].cards;
        state.currentBoard.lists[listIndex] = { ...state.currentBoard.lists[listIndex], ...updatedList, cards: existingCards };
        state.currentBoard.lists.sort((a, b) => a.position - b.position);
      }
    }
  },
  listDeletedByWS: (state, action) => {
    const { id: deletedListId, boardID } = action.payload; // payload is {id: listID, boardID: boardID}
    if (state.currentBoard && state.currentBoard.id === boardID) {
      state.currentBoard.lists = state.currentBoard.lists.filter(l => l.id !== deletedListId);
    }
  },
  cardCreatedByWS: (state, action) => {
    const newCard = action.payload;
    if (state.currentBoard) {
      const list = state.currentBoard.lists.find(l => l.id === newCard.listID);
      if (list) {
        if (!list.cards) list.cards = [];
         // Check if card already exists
        if (!list.cards.find(c => c.id === newCard.id)) {
          list.cards.push(newCard);
          list.cards.sort((a, b) => a.position - b.position);
        }
      }
    }
  },
  cardUpdatedByWS: (state, action) => {
    const updatedCard = action.payload;
    if (state.currentBoard) {
      const list = state.currentBoard.lists.find(l => l.id === updatedCard.listID);
      if (list) {
        const cardIndex = list.cards.findIndex(c => c.id === updatedCard.id);
        if (cardIndex !== -1) {
          list.cards[cardIndex] = { ...list.cards[cardIndex], ...updatedCard };
        } else { // Card might have moved to this list and updated simultaneously
          list.cards.push(updatedCard);
          list.cards.sort((a, b) => a.position - b.position);
        }
      }
    }
  },
  cardDeletedByWS: (state, action) => {
    const { id: deletedCardId, listID, boardID } = action.payload; // payload is {id: cardID, listID: listID, boardID: boardID}
    if (state.currentBoard && state.currentBoard.id === boardID) {
      const list = state.currentBoard.lists.find(l => l.id === listID);
      if (list) {
        list.cards = list.cards.filter(c => c.id !== deletedCardId);
      }
    }
  },
  cardMovedByWS: (state, action) => {
    const { cardId, oldListId, newListId, oldPosition, newPosition, boardId, updatedCard } = action.payload; // Assuming payload includes the full updated card for simplicity
    if (state.currentBoard && state.currentBoard.id === boardId) {
      // Remove from old list
      const oldList = state.currentBoard.lists.find(l => l.id === oldListId);
      if (oldList) {
        oldList.cards = oldList.cards.filter(c => c.id !== cardId);
        // oldList.cards.sort((a,b) => a.position - b.position); // Re-sort not strictly needed after removal unless positions are dense
      }
      // Add to new list
      const newList = state.currentBoard.lists.find(l => l.id === newListId);
      if (newList) {
        if (!newList.cards.find(c => c.id === cardId)) { // Add if not already there (e.g. from optimistic update)
             newList.cards.push(updatedCard); // The payload should contain the full card object with its new position and listID
        } else { // if already there, make sure it's updated
            const cardIndex = newList.cards.findIndex(c => c.id === cardId);
            newList.cards[cardIndex] = updatedCard;
        }
        newList.cards.sort((a, b) => a.position - b.position);
      }
       // Ensure positions are correctly updated for remaining cards in oldList
      if (oldList && oldListId !== newListId) {
        oldList.cards.forEach(card => {
          if (card.position > oldPosition) {
            // This can be complex as backend should be source of truth for positions.
            // For now, client relies on backend sending correct positions for all affected cards,
            // or a separate message for each card whose position changes.
            // Simplest: rely on the moved card's data and let subsequent updates fix other cards if needed.
          }
        });
      }
    }
  },
  // TODO: Add reducers for BOARD_MEMBER_ADDED, BOARD_MEMBER_REMOVED
  // TODO: Add reducers for CARD_COLLABORATOR_ADDED, CARD_COLLABORATOR_REMOVED
  // These would typically update state.currentBoard.members or state.currentBoard.lists[listIndex].cards[cardIndex].collaborators
};


// Correctly export all actions including the new WebSocket ones
const existingActions = boardsSlice.actions;
export const {
    optimisticallyUpdateCardOrder,
    optimisticallyMoveCardBetweenLists,
    clearCurrentBoard,
    boardUpdatedByWS,
    boardDeletedByWS,
    listCreatedByWS,
    listUpdatedByWS,
    listDeletedByWS,
    cardCreatedByWS,
    cardUpdatedByWS,
    cardDeletedByWS,
    cardMovedByWS,
    // Add other WS actions here as they are created
} = existingActions;

// This re-assignment is for clarity; webSocketActions can be used directly if preferred in components.
export const webSocketActions = existingActions;


export const selectUserBoards = (state) => state.boards.userBoards;
export const selectUserBoardsStatus = (state) => state.boards.userBoardsStatus;
export const selectCurrentBoard = (state) => state.boards.currentBoard;
export const selectCurrentBoardStatus = (state) => state.boards.currentBoardStatus;
export const selectCurrentBoardError = (state) => state.boards.currentBoardError;
export const selectListCardOpStatus = (state) => state.boards.listCardOpStatus;
export const selectListCardOpError = (state) => state.boards.listCardOpError;


export default boardsSlice.reducer;
