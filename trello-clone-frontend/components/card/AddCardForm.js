import React, { useState } from 'react';
import { Box, TextField, Button } from '@mui/material';

const AddCardForm = ({ listId, onAddCardToList, onCancel }) => {
  const [title, setTitle] = useState('');

  const handleSubmit = () => {
    if (title.trim()) {
      console.log('AddCardForm: listId =', listId, 'title =', title.trim());
      onAddCardToList(listId, title.trim());
      setTitle(''); // Clear input after submission
    }
  };

  return (
    <Box sx={{ mt: 2 }}>
      <TextField
        fullWidth
        multiline
        minRows={2}
        placeholder="Enter a title for this card..."
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        variant="outlined"
        size="small"
        sx={{ mb: 1, backgroundColor: 'white' }}
      />
      <Button
        variant="contained"
        onClick={handleSubmit}
        sx={{ mr: 1 }}
      >
        Add card
      </Button>
      <Button onClick={onCancel}>
        Cancel
      </Button>
    </Box>
  );
};

export default AddCardForm;
