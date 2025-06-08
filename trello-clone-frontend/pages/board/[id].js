import React, { useEffect, useState, useRef } from 'react'; // Import useRef
import { useRouter } from 'next/router';
import { useDispatch, useSelector } from 'react-redux';
import {
  fetchBoardDetails, addListToBoard, updateList, deleteList, moveCard,
  updateCardOrderInList, updateCard, deleteCard, addCommentToCard,
  deleteCommentFromCard, updateBoard, deleteBoard, webSocketActions
} from '../../features/boards/boardsSlice'; // Import webSocketActions
import { Container, Typography, CircularProgress, Box, Button, TextField, Dialog, DialogActions, DialogContent, DialogTitle, Grid, IconButton, Menu, MenuItem, Paper } from '@mui/material';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import AddIcon from '@mui/icons-material/Add'; // Import AddIcon
import CardDetailModal from '../../components/card/CardDetailModal';
import { addCardToList, /* other actions */ } from '../../features/boards/boardsSlice';
import TaskList from '../../components/list/TaskList'; // Import TaskList
import { DndContext, closestCorners, PointerSensor, useSensor, useSensors } from '@dnd-kit/core';
import { SortableContext, horizontalListSortingStrategy } from '@dnd-kit/sortable';
import { arrayMoveImmutable } from 'array-move';

export default function BoardPage() {
  const router = useRouter();
  const { id } = router.query;
  const dispatch = useDispatch();

  const { user, token } = useSelector((state) => state.auth);
  const { currentBoard, currentBoardStatus, currentBoardError } = useSelector((state) => state.boards);

  const [openCreateListDialog, setOpenCreateListDialog] = useState(false);
  const [newListName, setNewListName] = useState('');
  const [selectedCard, setSelectedCard] = useState(null); // For modal
  const [isModalOpen, setIsModalOpen] = useState(false);

  const [anchorElBoardMenu, setAnchorElBoardMenu] = useState(null); // For board options menu
  const [openEditBoardDialog, setOpenEditBoardDialog] = useState(false);
  const [editedBoardName, setEditedBoardName] = useState(currentBoard?.name || '');
  const [editedBoardDescription, setEditedBoardDescription] = useState(currentBoard?.description || '');
  const [openDeleteBoardConfirm, setOpenDeleteBoardConfirm] = useState(false);

  const [mounted, setMounted] = useState(false); // State to track if component is mounted
  const wsRef = useRef(null); // Ref to hold the WebSocket instance

  useEffect(() => {
    setMounted(true); // Set mounted to true after component mounts on client
  }, []);

  // Effect for WebSocket connection
  useEffect(() => {
    if (!mounted || !id || !token) { // Ensure component is mounted, boardId and token are available
      return;
    }

    // Construct WebSocket URL
    // TODO: Use wss:// for production and configure the domain/port appropriately
    const wsHostname = window.location.hostname;
    // Assuming backend runs on port 8080 during development, adjust if necessary
    const wsPort = process.env.NEXT_PUBLIC_WS_PORT || '8080'; // Allow overriding port via env var
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';

    // Use actual hostname for WS connection, not necessarily localhost if accessing remotely
    const wsUrl = `${wsProtocol}//${wsHostname}:${wsPort}/ws?boardID=${id}&token=${encodeURIComponent(token)}`;

    console.log('Attempting to connect WebSocket to:', wsUrl);
    const ws = new WebSocket(wsUrl);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log('WebSocket connection established for board:', id);
      // Optionally send a ping or initial message if required by backend
    };

    ws.onmessage = (event) => {
      // TODO: Handle incoming messages from the server
      console.log('WebSocket message received for board:', id, event.data);
      try {
        const message = JSON.parse(event.data);
        console.log('Parsed WebSocket message:', message);

        // Ensure currentBoard.id is a string for comparison if id from router is a string
        const currentBoardIdStr = currentBoard?.id?.toString();
        const messageBoardIdStr = message.payload?.boardID?.toString() || message.payload?.boardId?.toString();


        // For board-specific messages, ensure it's for the current board
        // BOARD_DELETED is special, payload is {id: boardID}
        if (message.type === 'BOARD_DELETED' && message.payload?.id?.toString() === id) {
           dispatch(webSocketActions.boardDeletedByWS(message.payload));
           // Consider redirecting or showing a message that the board was deleted
           router.push('/'); // Example: redirect to home
           return;
        }

        // For other messages, check if the message's boardID matches the current board.
        // Some payloads might have boardID directly, others might be nested (e.g. list.boardID).
        // The backend should ideally ensure only relevant messages are sent.
        // Here, we add a basic check if messageBoardIdStr is available and matches.
        // If messageBoardIdStr is not available in payload, we might assume it's for current board OR refine this check.
        if (messageBoardIdStr && messageBoardIdStr !== id) {
            console.log(`WebSocket message for board ${messageBoardIdStr} ignored on board ${id} view.`);
            return;
        }


        switch (message.type) {
          case 'BOARD_UPDATED':
            if (message.payload?.id?.toString() === id) {
              dispatch(webSocketActions.boardUpdatedByWS(message.payload));
            }
            break;
          case 'LIST_CREATED':
            // Ensure boardID in payload matches current board if provided, or assume it's for current board
             if (message.payload?.boardID?.toString() === id) {
                dispatch(webSocketActions.listCreatedByWS(message.payload));
             }
            break;
          case 'LIST_UPDATED':
             if (message.payload?.boardID?.toString() === id) {
                dispatch(webSocketActions.listUpdatedByWS(message.payload));
             }
            break;
          case 'LIST_DELETED':
            // Payload for LIST_DELETED is {id: listID, boardID: boardID}
            if (message.payload?.boardID?.toString() === id) {
                dispatch(webSocketActions.listDeletedByWS(message.payload));
            }
            break;
          case 'CARD_CREATED':
            // Assuming card payload includes listID, and we can infer boardID or it's in payload
            // For now, assume it's for the current board if it passed the initial boardID check or if no specific boardID in card payload
            dispatch(webSocketActions.cardCreatedByWS(message.payload));
            break;
          case 'CARD_UPDATED':
            dispatch(webSocketActions.cardUpdatedByWS(message.payload));
            break;
          case 'CARD_DELETED':
             // Payload for CARD_DELETED is {id: cardID, listID: listID, boardID: boardID}
            if (message.payload?.boardID?.toString() === id) {
                dispatch(webSocketActions.cardDeletedByWS(message.payload));
            }
            break;
          case 'CARD_MOVED':
            // Payload for CARD_MOVED is { cardId, oldListId, newListId, newPosition, boardId, updatedCard }
            if (message.payload?.boardId?.toString() === id) {
                dispatch(webSocketActions.cardMovedByWS(message.payload));
            }
            break;
          // TODO: Add cases for BOARD_MEMBER_ADDED, BOARD_MEMBER_REMOVED
          // TODO: Add cases for CARD_COLLABORATOR_ADDED, CARD_COLLABORATOR_REMOVED
          default:
            console.log('Received unhandled WebSocket message type:', message.type);
        }
      } catch (error) {
        console.error('Failed to parse or process WebSocket message:', error, event.data);
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error for board:', id, error);
    };

    ws.onclose = (event) => {
      console.log('WebSocket connection closed for board:', id, 'Code:', event.code, 'Reason:', event.reason);
      // Optionally implement reconnection logic here if desired
    };

    // Cleanup on component unmount or if id/token changes
    return () => {
      if (wsRef.current) {
        console.log('Closing WebSocket connection for board:', id);
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, [id, token, mounted, dispatch, currentBoard, router]); // Added dispatch, currentBoard, router to dependencies

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 5, // Optional tweak — requires slight movement before activating
      },
    })
  );

  useEffect(() => {
    if (mounted) { // Only run this effect on the client side after mounting
      if (!token) {
        router.push('/login');
      } else if (id && currentBoardStatus === 'idle' && !currentBoard) { // Fetch only if not already loaded
        dispatch(fetchBoardDetails(id));
      }
    }
  }, [id, token, router, dispatch, currentBoardStatus, mounted, currentBoard]); // Added currentBoard to dependencies

  const handleOpenCreateListDialog = () => setOpenCreateListDialog(true);
  const handleCloseCreateListDialog = () => {
    setOpenCreateListDialog(false);
    setNewListName('');
  };

  const handleCreateList = async () => {
    if (newListName.trim() && id) {
      await dispatch(addListToBoard({ boardId: id, name: newListName }));
      handleCloseCreateListDialog();
    }
  };

  const handleAddCard = (listId, title) => {
    // For now, only title is passed from AddCardForm.
    // Other fields (dueDate, assignedUser, etc.) will be default or null.
    // They can be edited via the CardDetailModal.
    dispatch(addCardToList({ listId, title }));
  };

  const handleCardClick = (card) => {
    setSelectedCard(card);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setSelectedCard(null); // Clear selected card
  };

  const [activeId, setActiveId] = useState(null);

  function handleDragStart(event) {
    setActiveId(event.active.id);
  }

  function handleDragEnd(event) {
    const { active, over } = event;
    if (!over) return;

    const lists = currentBoard?.lists || [];

    // Handle list reordering
    if (active.id.toString().startsWith('list-') && over.id.toString().startsWith('list-')) {
      const oldIndex = lists.findIndex(list => list.id === parseInt(active.id.replace('list-', '')));
      const newIndex = lists.findIndex(list => list.id === parseInt(over.id.replace('list-', '')));

      if (oldIndex !== newIndex) {
        const newListOrder = arrayMoveImmutable(lists, oldIndex, newIndex);
        newListOrder.forEach((list, index) => {
          if (list.position !== index + 1) {
            dispatch(updateList({ listId: list.id, listData: { position: index + 1 } }));
          }
        });
      }
    }

    // Handle card reordering or moving between lists
    if (active.id.toString().startsWith('card-')) {
      const activeCardId = parseInt(active.id.replace('card-', ''));
      const overId = over.id; // This could be a list ID or another card ID

      const activeList = lists.find(list => list.cards.some(card => card.id === activeCardId));
      const overList = lists.find(list => list.id === parseInt(overId.replace('list-', '')) || list.cards.some(card => card.id === parseInt(overId.replace('card-', ''))));

      if (!activeList || !overList) return;

      const activeCard = activeList.cards.find(card => card.id === activeCardId);
      if (!activeCard) return;

      const sourceListId = activeList.id;
      let destinationListId = overList.id;

      // Determine the new position within the destination list
      let newPosition;
      if (over.id.toString().startsWith('card-')) {
        // Dropped over another card
        const overCardIndex = overList.cards.findIndex(card => card.id === parseInt(over.id.replace('card-', '')));
        newPosition = overCardIndex + 1; // 1-based index
      } else {
        // Dropped over an empty list or at the end of a list
        newPosition = overList.cards.length + 1; // Append to end
      }

      if (sourceListId === destinationListId) {
        // Reordering within the same list
        const oldIndex = activeList.cards.findIndex(card => card.id === activeCardId);
        const newIndex = newPosition - 1; // Convert 1-based to 0-based for arrayMoveImmutable

        const newOrderedCards = arrayMoveImmutable(activeList.cards, oldIndex, newIndex);
        dispatch(optimisticallyUpdateCardOrder({ listId: sourceListId, orderedCards: newOrderedCards }));
        dispatch(updateCardOrderInList({ cardId: activeCardId, newPosition: newPosition, listId: sourceListId }));
      } else {
        // Moving between different lists
        const newSourceCards = activeList.cards.filter(card => card.id !== activeCardId);
        const newDestinationCards = [...overList.cards];
        newDestinationCards.splice(newPosition - 1, 0, { ...activeCard, listID: destinationListId, position: newPosition });

        dispatch(optimisticallyMoveCardBetweenLists({
          sourceListId: sourceListId,
          destinationListId: destinationListId,
          sourceCards: newSourceCards,
          destinationCards: newDestinationCards,
        }));
        dispatch(moveCard({ cardId: activeCardId, targetListId: destinationListId, newPosition: newPosition }));
      }
    }
    setActiveId(null);
  }

  function handleDragCancel() {
    setActiveId(null);
  }

  // Board options menu handlers
  const handleBoardMenuOpen = (event) => setAnchorElBoardMenu(event.currentTarget);
  const handleBoardMenuClose = () => setAnchorElBoardMenu(null);

  // Edit Board handlers
  const handleOpenEditBoardDialog = () => {
    setEditedBoardName(currentBoard.name);
    setEditedBoardDescription(currentBoard.description || '');
    setOpenEditBoardDialog(true);
    handleBoardMenuClose();
  };
  const handleCloseEditBoardDialog = () => setOpenEditBoardDialog(false);
  const handleUpdateBoard = async () => {
    if (editedBoardName.trim()) {
      await dispatch(updateBoard({ boardId: id, boardData: { name: editedBoardName, description: editedBoardDescription } }));
      handleCloseEditBoardDialog();
    }
  };

  // Delete Board handlers
  const handleOpenDeleteBoardConfirm = () => {
    setOpenDeleteBoardConfirm(true);
    handleBoardMenuClose();
  };
  const handleCloseDeleteBoardConfirm = () => setOpenDeleteBoardConfirm(false);
  const handleDeleteBoard = async () => {
    await dispatch(deleteBoard(id));
    handleCloseDeleteBoardConfirm();
    router.push('/'); // Redirect to dashboard after deleting board
  };

  // Render a loading state or nothing during SSR, or if not mounted yet
  if (!mounted || !user) { // Check mounted first, then user
    return (
      <Container sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '80vh' }}>
         <CircularProgress />
      </Container>
    )
  }

  if (currentBoardStatus === 'loading') {
    return (
      <Container sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '80vh' }}>
        <CircularProgress />
      </Container>
    );
  }

  if (currentBoardStatus === 'failed') {
    return (
      <Container sx={{ mt: 4 }}>
        <Typography color="error">Error loading board: {currentBoardError}</Typography>
      </Container>
    );
  }

  if (!currentBoard) {
    return (
      <Container sx={{ mt: 4 }}>
        <Typography>Board not found.</Typography>
      </Container>
    );
  }

  // currentBoard.lists and currentBoard.lists[].cards are already sorted by position due to extraReducers
  const lists = currentBoard?.lists || [];

  return (
    <Container
      maxWidth={false} // Use maxWidth={false} to allow full width
      disableGutters // Remove default padding
      sx={{
        backgroundColor: '#1d2125', // Dark background for the entire board area
        minHeight: 'calc(100vh - 64px)', // Full height minus Navbar height
        pt: 2, // Padding top for board header
        pb: 2, // Padding bottom
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2, px: 3 }}> {/* Add horizontal padding */}
        <Box>
          <Typography variant="h4" component="h1" gutterBottom sx={{ color: 'white', fontWeight: 'bold' }}>
            {currentBoard.name}
          </Typography>
          <Typography variant="body1" sx={{ color: '#a6b0cf' }}> {/* Lighter text for description */}
            {currentBoard.description}
          </Typography>
        </Box>
        <IconButton size="small" onClick={handleBoardMenuOpen} sx={{ color: '#a6b0cf' }}> {/* Icon color */}
          <MoreVertIcon fontSize="small" />
        </IconButton>
        <Menu
          anchorEl={anchorElBoardMenu}
          open={Boolean(anchorElBoardMenu)}
          onClose={handleBoardMenuClose}
        >
          <MenuItem onClick={handleOpenEditBoardDialog}>Edit Board Details</MenuItem>
          <MenuItem onClick={handleOpenDeleteBoardConfirm}>Delete Board</MenuItem>
        </Menu>
      </Box>

      <DndContext
        sensors={sensors}
        collisionDetection={closestCorners}
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
        onDragCancel={handleDragCancel}
      >
        <SortableContext items={lists.map(list => `list-${list.id}`)} strategy={horizontalListSortingStrategy}>
          <Box
            sx={{
              display: 'flex',
              overflowX: 'auto',
              pb: 2, // Padding bottom for scrollbar
              px: 3, // Horizontal padding for lists container
              flexGrow: 1, // Allow lists container to grow
              alignItems: 'flex-start', // Align lists to the top
              '&::-webkit-scrollbar': {
                height: '8px',
              },
              '&::-webkit-scrollbar-thumb': {
                backgroundColor: 'rgba(255,255,255,0.3)', // Lighter scrollbar for dark theme
                borderRadius: '4px',
              },
              '&::-webkit-scrollbar-track': {
                backgroundColor: 'rgba(0,0,0,0.1)', // Darker track
              },
            }}
          >
            {lists.map((list) => (
              <TaskList
                key={list.id}
                list={list}
                onAddCardToList={handleAddCard}
                onCardClick={handleCardClick} // Pass handler to TaskList
              />
            ))}
            <Paper
              sx={{
                minWidth: 272,
                maxWidth: 272,
                backgroundColor: 'rgba(255,255,255,0.1)', // Semi-transparent background for "Add list"
                borderRadius: '12px',
                p: 1.5,
                ml: 2,
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center',
                alignItems: 'center',
                flexShrink: 0,
                boxShadow: 'none', // No shadow for this button
              }}
            >
              <Button
                fullWidth
                startIcon={<AddIcon sx={{ color: 'white' }} />} // White icon
                onClick={handleOpenCreateListDialog}
                sx={{
                  color: 'white', // White text
                  textTransform: 'none',
                  justifyContent: 'flex-start',
                  py: 1,
                  '&:hover': { backgroundColor: 'rgba(255,255,255,0.2)' }, // Lighter hover
                }}
              >
                Add another list
              </Button>
            </Paper>
          </Box>
        </SortableContext>
      </DndContext>

      {/* Create List Dialog */}
      <Dialog open={openCreateListDialog} onClose={handleCloseCreateListDialog}>
        <DialogTitle>Create New List</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            id="listName"
            label="List Name"
            type="text"
            fullWidth
            variant="outlined"
            value={newListName}
            onChange={(e) => setNewListName(e.target.value)}
            required
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseCreateListDialog}>Cancel</Button>
          <Button onClick={handleCreateList} variant="contained" disabled={!newListName.trim()}>Create</Button>
        </DialogActions>
      </Dialog>

      {selectedCard && (
        <CardDetailModal
          open={isModalOpen}
          onClose={handleCloseModal}
          card={selectedCard}
          boardId={id} // Pass boardId if needed by the modal
        />
      )}

      {/* Edit Board Dialog */}
      <Dialog open={openEditBoardDialog} onClose={handleCloseEditBoardDialog}>
        <DialogTitle>Edit Board Details</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            id="boardName"
            label="Board Name"
            type="text"
            fullWidth
            variant="outlined"
            value={editedBoardName}
            onChange={(e) => setEditedBoardName(e.target.value)}
            required
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            id="boardDescription"
            label="Board Description (Optional)"
            type="text"
            fullWidth
            multiline
            rows={3}
            variant="outlined"
            value={editedBoardDescription}
            onChange={(e) => setEditedBoardDescription(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseEditBoardDialog}>Cancel</Button>
          <Button onClick={handleUpdateBoard} variant="contained" disabled={!editedBoardName.trim()}>Save</Button>
        </DialogActions>
      </Dialog>

      {/* Delete Board Confirmation Dialog */}
      <Dialog open={openDeleteBoardConfirm} onClose={handleCloseDeleteBoardConfirm}>
        <DialogTitle>Confirm Delete Board</DialogTitle>
        <DialogContent>
          <Typography>Are you sure you want to delete the board "{currentBoard?.name}"? This action cannot be undone.</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDeleteBoardConfirm}>Cancel</Button>
          <Button onClick={handleDeleteBoard} variant="contained" color="error">Delete</Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
}
