// components/card/CardDetailModal.js
import React, { useState, useEffect } from 'react';
import {
  Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField,
  Select, MenuItem, FormControl, InputLabel, Grid, Typography, IconButton, Box
} from '@mui/material';
import { useDispatch } from 'react-redux';
import { updateCardDetails } from '../../features/boards/boardsSlice';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import CloseIcon from '@mui/icons-material/Close';

// Dummy users for supervisor/assignee select. In a real app, fetch these.
const DUMMY_USERS = [
  { id: 1, username: 'Alice Smith' },
  { id: 2, username: 'Bob Johnson' },
  { id: 3, username: 'Charlie Brown' },
];


const CardDetailModal = ({ open, onClose, card }) => {
  const dispatch = useDispatch();
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [dueDate, setDueDate] = useState(null);
  const [status, setStatus] = useState('to_do');
  const [assignedUserID, setAssignedUserID] = useState(''); // Store as string for Select, convert to number or null
  const [supervisorID, setSupervisorID] = useState('');   // Store as string

  useEffect(() => {
    if (card) {
      setTitle(card.title || '');
      setDescription(card.description || '');
      setDueDate(card.dueDate ? new Date(card.dueDate) : null);
      setStatus(card.status || 'to_do');
      setAssignedUserID(card.assignedUserID ? String(card.assignedUserID) : '');
      setSupervisorID(card.supervisorID ? String(card.supervisorID) : '');
    }
  }, [card]);

  if (!card) return null;

  const handleSave = () => {
    const payload = {
      cardId: card.id,
      title: title.trim() === '' ? card.title : title.trim(), // Keep original if empty
      description,
      // Ensure dueDate is in ISO string format or null for backend
      dueDate: dueDate ? dueDate.toISOString() : null,
      status,
      assignedUserID: assignedUserID ? parseInt(assignedUserID) : null,
      supervisorID: supervisorID ? parseInt(supervisorID) : null,
    };
    dispatch(updateCardDetails(payload));
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
      <DialogTitle sx={{ m: 0, p: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        Edit Card
        <IconButton aria-label="close" onClick={onClose} sx={{color: (theme) => theme.palette.grey[500]}}>
          <CloseIcon />
        </IconButton>
      </DialogTitle>
      <DialogContent dividers>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <TextField
              label="Title"
              fullWidth
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              variant="outlined"
              margin="dense"
            />
          </Grid>
          <Grid item xs={12}>
            <TextField
              label="Description"
              fullWidth
              multiline
              rows={4}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              variant="outlined"
              margin="dense"
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <LocalizationProvider dateAdapter={AdapterDateFns}>
              <DatePicker
                label="Due Date"
                value={dueDate}
                onChange={(newValue) => setDueDate(newValue)}
                renderInput={(params) => <TextField {...params} fullWidth margin="dense" helperText={params.error ? "Invalid date" : ""} />}
              />
            </LocalizationProvider>
          </Grid>
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth margin="dense">
              <InputLabel id="status-label">Status</InputLabel>
              <Select
                labelId="status-label"
                value={status}
                label="Status"
                onChange={(e) => setStatus(e.target.value)}
              >
                <MenuItem value="to_do">To Do</MenuItem>
                <MenuItem value="pending">Pending</MenuItem>
                <MenuItem value="done">Done</MenuItem>
                {/* <MenuItem value="undone">Undone</MenuItem> */}
              </Select>
            </FormControl>
          </Grid>
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth margin="dense">
              <InputLabel id="assignee-label">Assign To</InputLabel>
              <Select
                labelId="assignee-label"
                value={assignedUserID}
                label="Assign To"
                onChange={(e) => setAssignedUserID(e.target.value)}
              >
                <MenuItem value=""><em>None</em></MenuItem>
                {DUMMY_USERS.map(user => (
                  <MenuItem key={user.id} value={String(user.id)}>{user.username}</MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth margin="dense">
              <InputLabel id="supervisor-label">Supervisor</InputLabel>
              <Select
                labelId="supervisor-label"
                value={supervisorID}
                label="Supervisor"
                onChange={(e) => setSupervisorID(e.target.value)}
              >
                <MenuItem value=""><em>None</em></MenuItem>
                 {DUMMY_USERS.map(user => ( // Assuming supervisors are also users
                  <MenuItem key={user.id} value={String(user.id)}>{user.username}</MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
        </Grid>
      </DialogContent>
      <DialogActions sx={{ p: '16px 24px' }}>
        <Button onClick={onClose}>Cancel</Button>
        <Button onClick={handleSave} variant="contained">Save Changes</Button>
      </DialogActions>
    </Dialog>
  );
};

export default CardDetailModal;