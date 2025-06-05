import { configureStore } from '@reduxjs/toolkit';
import authReducer from '../features/auth/authSlice';
import boardsReducer from '../features/boards/boardsSlice';
// Import other reducers (lists, cards) here

export const store = configureStore({
  reducer: {
    auth: authReducer,
    boards: boardsReducer,
    // lists: listsReducer,
    // cards: cardsReducer,
  },
  // Middleware can be added here if needed
});

export default store;