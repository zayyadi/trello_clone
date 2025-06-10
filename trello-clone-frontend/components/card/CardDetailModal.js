// components/card/CardDetailModal.js
import React, { useState, useEffect } from 'react';
import {
  Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField,
  Select, MenuItem, FormControl, InputLabel, Grid, Typography, IconButton, Box,
  List, ListItem, ListItemText, Chip, CircularProgress, Alert
} from '@mui/material';
import { useDispatch, useSelector } from 'react-redux';
import {
    updateCardDetails,
    fetchCardCollaborators,
    addCardCollaborator,
    removeCardCollaborator,
    moveCard // Import moveCard thunk
} from '../../features/boards/boardsSlice';
import { selectCurrentUser } from '../../features/auth/authSlice';
import { selectCurrentBoard, selectListCardOpStatus, selectListCardOpError } from '../../features/boards/boardsSlice';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import CloseIcon from '@mui/icons-material/Close';
import DeleteIcon from '@mui/icons-material/Delete';


const CardDetailModal = ({ open, onClose, card }) => {
  const dispatch = useDispatch();
  const currentUser = useSelector(selectCurrentUser);
  const currentBoard = useSelector(selectCurrentBoard);
  const listCardOpStatus = useSelector(selectListCardOpStatus);
  const listCardOpError = useSelector(selectListCardOpError);

  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [dueDate, setDueDate] = useState(null);
  const [status, setStatus] = useState('TO_DO'); // Default to uppercase
  const [assignedUserID, setAssignedUserID] = useState('');
  const [supervisorID, setSupervisorID] = useState('');
  const [collaboratorInput, setCollaboratorInput] = useState('');
  const [targetListId, setTargetListId] = useState(''); // State for target list for move

  const isBoardOwner = currentUser?.id === currentBoard?.ownerID;
  // Card details (including collaborators) will come from the `card` prop, which is updated by Redux state
  const isCollaborator = card?.collaborators?.some(c => c.id === currentUser?.id);
  const isAssignee = card?.assignedUserID === currentUser?.id;
  const isCollaboratorOrAssignee = isCollaborator || isAssignee;

  // Dummy users for supervisor/assignee select. In a real app, fetch these or get from board members.
  // For collaborators, we will use the actual collaborator list from the card.
  const boardMembersForSelect = currentBoard?.members?.map(member => member.user) || [];
  const availableListsForMove = currentBoard?.lists?.filter(l => l.id !== card?.listID) || [];


  useEffect(() => {
    if (card) {
      setTitle(card.title || '');
      setDescription(card.description || '');
      setDueDate(card.dueDate ? new Date(card.dueDate) : null);
      setStatus(card.status?.toUpperCase() || 'TO_DO'); // Convert to uppercase on load, default to TO_DO
      setAssignedUserID(card.assignedUserID ? String(card.assignedUserID) : '');
      setSupervisorID(card.supervisorID ? String(card.supervisorID) : '');
      setTargetListId(''); // Reset target list when card changes or modal opens

      // Fetch collaborators if card exists and collaborators are not yet fetched/present
      // (Or if they might have changed, though for simplicity, just fetch if undefined)
      if (card.id && card.collaborators === undefined) { // Check if undefined to fetch initially
        dispatch(fetchCardCollaborators(card.id));
      }
    }
  }, [card, dispatch]);

  if (!card) return null;

  const handleSave = () => {
    const payload = {
      cardId: card.id,
      title: title.trim() === '' ? card.title : title.trim(),
      description,
      dueDate: dueDate ? dueDate.toISOString() : null,
      status,
      assignedUserID: assignedUserID ? parseInt(assignedUserID) : null,
      supervisorID: supervisorID ? parseInt(supervisorID) : null,
    };
    dispatch(updateCardDetails(payload));
    onClose();
  };

  const handleAddCollaborator = () => {
    if (collaboratorInput.trim() === '') return;
    // Simple check if input is numeric for ID, otherwise assume email
    const isNumeric = /^\d+$/.test(collaboratorInput);
    let userIdentifier;
    if (isNumeric) {
      userIdentifier = { userID: parseInt(collaboratorInput) };
    } else {
      userIdentifier = { email: collaboratorInput };
    }
    dispatch(addCardCollaborator({ cardId: card.id, userIdentifier }));
    setCollaboratorInput('');
  };

  const handleRemoveCollaborator = (userIdToRemove) => {
    dispatch(removeCardCollaborator({ cardId: card.id, userIdToRemove }));
  };

  const handleMoveCard = () => {
    if (!targetListId) {
      console.warn("No target list selected for move.");
      // Optionally, show an alert to the user
      return;
    }
    const targetList = currentBoard?.lists?.find(l => l.id === parseInt(targetListId));
    const newPosition = targetList?.cards?.length + 1 || 1;

    dispatch(moveCard({
      cardId: card.id,
      targetListId: parseInt(targetListId),
      newPosition: newPosition
    }));
    onClose(); // Close modal after dispatching move
  };

  // Logging for debugging member population
  console.log('[CardDetailModal] currentBoard from selector:', currentBoard);
  console.log('[CardDetailModal] currentBoard.members:', currentBoard?.members);
  console.log('[CardDetailModal] boardMembersForSelect for dropdowns:', boardMembersForSelect);


  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="md"> {/* Changed to md for more space */}
      <DialogTitle sx={{ m: 0, p: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        Edit Card
        <IconButton aria-label="close" onClick={onClose} sx={{color: (theme) => theme.palette.grey[500]}}>
          <CloseIcon />
        </IconButton>
      </DialogTitle>
      <DialogContent dividers>
        <Grid container spacing={3}> {/* Increased spacing */}
          <Grid item xs={12} md={8}> {/* Main content area */}
            <TextField
              label="Title"
              fullWidth
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              variant="outlined"
              margin="dense"
              disabled={!isBoardOwner}
            />
            <TextField
              label="Description"
              fullWidth
              multiline
              rows={4}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              variant="outlined"
              margin="dense"
              disabled={!isBoardOwner && !isCollaboratorOrAssignee}
            />
          </Grid>
          <Grid item xs={12} md={4}> {/* Sidebar area for details */}
            <LocalizationProvider dateAdapter={AdapterDateFns}>
              <DatePicker
                label="Due Date"
                value={dueDate}
                onChange={(newValue) => setDueDate(newValue)}
                renderInput={(params) => <TextField {...params} fullWidth margin="dense" helperText={params.error ? "Invalid date" : ""} />}
                disabled={!isBoardOwner && !isCollaboratorOrAssignee}
              />
            </LocalizationProvider>
            <FormControl fullWidth margin="dense" disabled={!isBoardOwner && !isCollaboratorOrAssignee}>
              <InputLabel id="status-label">Status</InputLabel>
              <Select
                labelId="status-label"
                value={status}
                label="Status"
                onChange={(e) => setStatus(e.target.value)}
              >
                <MenuItem value="TO_DO">To Do</MenuItem>
                <MenuItem value="PENDING">Pending</MenuItem>
                <MenuItem value="DONE">Done</MenuItem>
                <MenuItem value="UNDONE">Undone</MenuItem>
              </Select>
            </FormControl>
            <FormControl fullWidth margin="dense" disabled={!isBoardOwner && !isCollaboratorOrAssignee}>
              <InputLabel id="assignee-label">Assign To</InputLabel>
              <Select
                labelId="assignee-label"
                value={assignedUserID}
                label="Assign To"
                onChange={(e) => setAssignedUserID(e.target.value)}
              >
                <MenuItem value=""><em>None</em></MenuItem>
                {boardMembersForSelect.map(user => (
                  <MenuItem key={user.id} value={String(user.id)}>{user.username}</MenuItem>
                ))}
              </Select>
            </FormControl>

            {/* Move Card Dropdown */}
            <FormControl fullWidth margin="dense">
              <InputLabel id="move-card-target-list-label">Move to List</InputLabel>
              <Select
                labelId="move-card-target-list-label"
                value={targetListId}
                label="Move to List"
                onChange={(e) => setTargetListId(e.target.value)}
                disabled={availableListsForMove.length === 0}
              >
                <MenuItem value=""><em>Select a list...</em></MenuItem>
                {availableListsForMove.map(list => (
                  <MenuItem key={list.id} value={list.id}>{list.name}</MenuItem>
                ))}
              </Select>
            </FormControl>
            <FormControl fullWidth margin="dense" disabled={!isBoardOwner && !isCollaboratorOrAssignee}>
              <InputLabel id="supervisor-label">Supervisor</InputLabel>
              <Select
                labelId="supervisor-label"
                value={supervisorID}
                label="Supervisor"
                onChange={(e) => setSupervisorID(e.target.value)}
              >
                <MenuItem value=""><em>None</em></MenuItem>
                 {boardMembersForSelect.map(user => (
                  <MenuItem key={user.id} value={String(user.id)}>{user.username}</MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>

          {/* Collaborators Section */}
          <Grid item xs={12}>
            <Typography variant="h6" gutterBottom>Collaborators</Typography>
            {listCardOpStatus === 'loading_collaborators' && <CircularProgress size={20} />}
            {listCardOpError && <Alert severity="error">{listCardOpError}</Alert>}

            <List dense>
              {card.collaborators && card.collaborators.map(collab => (
                <ListItem
                  key={collab.id}
                  secondaryAction={
                    isBoardOwner && (
                      <IconButton edge="end" aria-label="delete" onClick={() => handleRemoveCollaborator(collab.id)}>
                        <DeleteIcon />
                      </IconButton>
                    )
                  }
                >
                  <ListItemText primary={collab.username} secondary={collab.email} />
                </ListItem>
              ))}
              {(!card.collaborators || card.collaborators.length === 0) && listCardOpStatus !== 'loading_collaborators' && (
                <Typography variant="body2" color="textSecondary">No collaborators yet.</Typography>
              )}
            </List>

            {isBoardOwner && (
              <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
                <TextField
                  label="Add Collaborator (Email or ID)"
                  size="small"
                  value={collaboratorInput}
                  onChange={(e) => setCollaboratorInput(e.target.value)}
                  variant="outlined"
                  sx={{ flexGrow: 1, mr: 1 }}
                />
                <Button onClick={handleAddCollaborator} variant="outlined" size="small" disabled={listCardOpStatus === 'loading_add_collaborator'}>
                  {listCardOpStatus === 'loading_add_collaborator' ? <CircularProgress size={20}/> : "Add"}
                </Button>
              </Box>
            )}
          </Grid>
        </Grid>
      </DialogContent>
      <DialogActions sx={{ p: '16px 24px', justifyContent: 'space-between' }}> {/* Adjusted for more buttons */}
        <Box> {/* Group action buttons */}
          <Button onClick={handleMoveCard} color="secondary" disabled={!targetListId || listCardOpStatus.startsWith('loading')}>
            Move Card
          </Button>
        </Box>
        <Box> {/* Group save/cancel */}
          <Button onClick={onClose} sx={{ mr: 1 }}>Cancel</Button>
          <Button onClick={handleSave} variant="contained" disabled={listCardOpStatus.startsWith('loading')}>Save Changes</Button>
        </Box>
      </DialogActions>
    </Dialog>
  );
};

export default CardDetailModal;