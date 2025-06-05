import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useRouter } from 'next/router';
import { fetchUserBoards, createNewBoard } from '../features/boards/boardsSlice';
import { Container, Typography, Grid, Paper, Button, CircularProgress, Box, TextField, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle } from '@mui/material';
import Link from 'next/link'; // For linking to individual boards

// Simple Board Card Component (can be moved to components/board/BoardCard.js)
const BoardCard = ({ board }) => (
  <Grid item xs={12} sm={6} md={4}>
    <Link href={`/board/${board.id}`} passHref>
      <Paper component="a" elevation={2} sx={{ padding: 2, textDecoration: 'none', display: 'block', '&:hover': { boxShadow: 6 } }}>
        <Typography variant="h6">{board.name}</Typography>
        <Typography variant="body2" color="text.secondary">{board.description || 'No description'}</Typography>
      </Paper>
    </Link>
  </Grid>
);


export default function HomePage() {
  const dispatch = useDispatch();
  const router = useRouter();
  const { user, token } = useSelector((state) => state.auth);
  const { userBoards, userBoardsStatus, userBoardsError } = useSelector((state) => state.boards);

  const [openCreateDialog, setOpenCreateDialog] = useState(false);
  const [newBoardName, setNewBoardName] = useState('');
  const [newBoardDescription, setNewBoardDescription] = useState('');
  const [mounted, setMounted] = useState(false); // State to track if component is mounted

  useEffect(() => {
    setMounted(true); // Set mounted to true after component mounts on client
  }, []);

  useEffect(() => {
    if (mounted) { // Only run this effect on the client side after mounting
      if (!token) {
        router.push('/login');
      } else {
        if (userBoardsStatus === 'idle') {
          dispatch(fetchUserBoards());
        }
      }
    }
  }, [token, router, dispatch, userBoardsStatus, mounted]);

  const handleOpenCreateDialog = () => setOpenCreateDialog(true);
  const handleCloseCreateDialog = () => {
    setOpenCreateDialog(false);
    setNewBoardName('');
    setNewBoardDescription('');
  };

  const handleCreateBoard = async () => {
    if (newBoardName.trim()) {
      await dispatch(createNewBoard({ name: newBoardName, description: newBoardDescription }));
      handleCloseCreateDialog();
      // Optionally, re-fetch boards or rely on the slice to update the list
      // dispatch(fetchUserBoards());
    }
  };

  // Render a loading state or nothing during SSR, or if not mounted yet
  if (!mounted || !user) { // Check mounted first, then user
    return (
      <Container sx={{display: 'flex', justifyContent: 'center', alignItems: 'center', height: '80vh'}}>
         <CircularProgress />
      </Container>
    )
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" component="h1">
          Your Boards
        </Typography>
        <Button variant="contained" onClick={handleOpenCreateDialog}>Create New Board</Button>
      </Box>

      {userBoardsStatus === 'loading' && <CircularProgress />}
      {userBoardsStatus === 'failed' && <Typography color="error">Error: {userBoardsError}</Typography>}
      {userBoardsStatus === 'succeeded' && (userBoards || []).length === 0 && (
        <Typography>No boards found. Create one!</Typography>
      )}
      {userBoardsStatus === 'succeeded' && (userBoards || []).length > 0 && (
        <Grid container spacing={3}>
          {(userBoards || []).map((board) => (
            <BoardCard key={board.id} board={board} />
          ))}
        </Grid>
      )}

      {/* Create Board Dialog */}
      <Dialog open={openCreateDialog} onClose={handleCloseCreateDialog}>
        <DialogTitle>Create New Board</DialogTitle>
        <DialogContent>
          <DialogContentText sx={{mb: 2}}>
            Enter a name and an optional description for your new board.
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            id="name"
            label="Board Name"
            type="text"
            fullWidth
            variant="outlined"
            value={newBoardName}
            onChange={(e) => setNewBoardName(e.target.value)}
            required
          />
          <TextField
            margin="dense"
            id="description"
            label="Board Description (Optional)"
            type="text"
            fullWidth
            multiline
            rows={3}
            variant="outlined"
            value={newBoardDescription}
            onChange={(e) => setNewBoardDescription(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseCreateDialog}>Cancel</Button>
          <Button onClick={handleCreateBoard} variant="contained" disabled={!newBoardName.trim()}>Create</Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
}
