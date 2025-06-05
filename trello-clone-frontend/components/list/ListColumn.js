import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { Paper, Typography, Box, Button, TextField, Dialog, DialogActions, DialogContent, DialogTitle, IconButton, Menu, MenuItem } from '@mui/material';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import CardItem from '../card/CardItem';
import { addCardToList, updateList, deleteList } from '../../features/boards/boardsSlice';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { useDroppable } from '@dnd-kit/core';

export default function ListColumn({ list, cards, boardId, onCardClick }) { // Removed index prop
  const dispatch = useDispatch();
  const [openCreateCardDialog, setOpenCreateCardDialog] = useState(false);
  const [newCardTitle, setNewCardTitle] = useState('');
  const [newCardDescription, setNewCardDescription] = useState('');

  const [anchorEl, setAnchorEl] = useState(null); // For list options menu
  const [openEditListDialog, setOpenEditListDialog] = useState(false);
  const [editedListName, setEditedListName] = useState(list.name);
  const [openDeleteListConfirm, setOpenDeleteListConfirm] = useState(false);

  const handleOpenCreateCardDialog = () => setOpenCreateCardDialog(true);
  const handleCloseCreateCardDialog = () => {
    setOpenCreateCardDialog(false);
    setNewCardTitle('');
    setNewCardDescription('');
  };

  const handleCreateCard = async () => {
    if (newCardTitle.trim()) {
      await dispatch(addCardToList({ listId: list.id, title: newCardTitle, description: newCardDescription }));
      handleCloseCreateCardDialog();
    }
  };

  // List options menu handlers
  const handleMenuOpen = (event) => setAnchorEl(event.currentTarget);
  const handleMenuClose = () => setAnchorEl(null);

  // Edit List handlers
  const handleOpenEditListDialog = () => {
    setEditedListName(list.name);
    setOpenEditListDialog(true);
    handleMenuClose();
  };
  const handleCloseEditListDialog = () => setOpenEditListDialog(false);
  const handleUpdateListName = async () => {
    if (editedListName.trim() && editedListName !== list.name) {
      await dispatch(updateList({ listId: list.id, listData: { name: editedListName } }));
    }
    handleCloseEditListDialog();
  };

  // Delete List handlers
  const handleOpenDeleteListConfirm = () => {
    setOpenDeleteListConfirm(true);
    handleMenuClose();
  };
  const handleCloseDeleteListConfirm = () => setOpenDeleteListConfirm(false);
  const handleDeleteList = async () => {
    await dispatch(deleteList(list.id));
    handleCloseDeleteListConfirm();
  };

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: `list-${list.id}` }); // Use list.id as the unique ID for sortable context

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    zIndex: isDragging ? 1000 : 'auto',
  };

  const { setNodeRef: setDroppableRef } = useDroppable({
    id: `list-${list.id}`, // Use list.id as the unique ID for droppable context
  });

  return (
    <Paper
      ref={setNodeRef}
      style={style}
      sx={{
        minWidth: 270,
        maxWidth: 270,
        backgroundColor: (theme) => theme.palette.grey[200],
        p: 2,
        mr: 2,
        flexShrink: 0,
        display: 'flex',
        flexDirection: 'column',
        maxHeight: 'calc(100vh - 150px)', // Adjust based on header/footer height
      }}
      {...attributes}
    >
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }} {...listeners}>
        <Typography variant="h6" sx={{ flexGrow: 1 }}>{list.name}</Typography>
        <IconButton size="small" onClick={handleMenuOpen}>
          <MoreVertIcon fontSize="small" />
        </IconButton>
        <Menu
          anchorEl={anchorEl}
          open={Boolean(anchorEl)}
          onClose={handleMenuClose}
        >
          <MenuItem onClick={handleOpenEditListDialog}>Edit List Name</MenuItem>
          <MenuItem onClick={handleOpenDeleteListConfirm}>Delete List</MenuItem>
        </Menu>
      </Box>

      <Box
        ref={setDroppableRef}
        sx={{
          flexGrow: 1,
          minHeight: '50px', // Ensure droppable area is visible
          overflowY: 'auto', // Enable scrolling for cards
          '&::-webkit-scrollbar': {
            width: '8px',
          },
          '&::-webkit-scrollbar-thumb': {
            backgroundColor: (theme) => theme.palette.grey[400],
            borderRadius: '4px',
          },
          '&::-webkit-scrollbar-track': {
            backgroundColor: (theme) => theme.palette.grey[200],
          },
        }}
      >
        <SortableContext items={cards.map(card => `card-${card.id}`)} strategy={verticalListSortingStrategy}>
          {cards.map((card) => (
            <CardItem key={card.id} card={card} onCardClick={onCardClick} />
          ))}
        </SortableContext>
      </Box>

      <Button variant="contained" sx={{ mt: 2 }} onClick={handleOpenCreateCardDialog} fullWidth>
        Add a card
      </Button>

      {/* Create Card Dialog */}
      <Dialog open={openCreateCardDialog} onClose={handleCloseCreateCardDialog}>
        <DialogTitle>Add New Card to {list.name}</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            id="cardTitle"
            label="Card Title"
            type="text"
            fullWidth
            variant="outlined"
            value={newCardTitle}
            onChange={(e) => setNewCardTitle(e.target.value)}
            required
          />
          <TextField
            margin="dense"
            id="cardDescription"
            label="Description (Optional)"
            type="text"
            fullWidth
            multiline
            rows={3}
            variant="outlined"
            value={newCardDescription}
            onChange={(e) => setNewCardDescription(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseCreateCardDialog}>Cancel</Button>
          <Button onClick={handleCreateCard} variant="contained" disabled={!newCardTitle.trim()}>Add Card</Button>
        </DialogActions>
      </Dialog>

      {/* Edit List Dialog */}
      <Dialog open={openEditListDialog} onClose={handleCloseEditListDialog}>
        <DialogTitle>Edit List Name</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            id="editListName"
            label="List Name"
            type="text"
            fullWidth
            variant="outlined"
            value={editedListName}
            onChange={(e) => setEditedListName(e.target.value)}
            required
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseEditListDialog}>Cancel</Button>
          <Button onClick={handleUpdateListName} variant="contained" disabled={!editedListName.trim()}>Save</Button>
        </DialogActions>
      </Dialog>

      {/* Delete List Confirmation Dialog */}
      <Dialog open={openDeleteListConfirm} onClose={handleCloseDeleteListConfirm}>
        <DialogTitle>Confirm Delete List</DialogTitle>
        <DialogContent>
          <Typography>Are you sure you want to delete the list "{list.name}"? This action cannot be undone.</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDeleteListConfirm}>Cancel</Button>
          <Button onClick={handleDeleteList} variant="contained" color="error">Delete</Button>
        </DialogActions>
      </Dialog>
    </Paper>
  );
}
