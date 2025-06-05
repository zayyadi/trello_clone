import React, { useState } from 'react';
import { Paper, Typography, Box, Button } from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import TaskCard from '../card/TaskCard';
import AddCardForm from '../card/AddCardForm'; // Import AddCardForm

const TaskList = ({ list, onAddCardToList, onCardClick }) => {
  const [showAddCardForm, setShowAddCardForm] = useState(false);

  const cards = list.cards || []; // Ensure cards is an array

  return (
    <Paper
      elevation={2}
      sx={{
        minWidth: 272, // Trello's default list width
        maxWidth: 272,
        backgroundColor: '#282e33', // Darker background from image
        borderRadius: '12px', // More rounded corners
        p: 1.5,
        display: 'flex',
        flexDirection: 'column',
        maxHeight: 'calc(100vh - 120px)', // Adjust based on header/footer
        overflowY: 'auto',
        mr: 2,
        boxShadow: '0px 1px 0px rgba(9,30,66,.25)', // Subtle shadow
      }}
    >
      <Typography variant="h6" sx={{ fontWeight: 'bold', mb: 1, color: 'white' }}>
        {list.name}
      </Typography>
      <Box sx={{ flexGrow: 1, overflowY: 'auto', pr: 0.5 }}>
        <SortableContext items={cards.map(c => c.id)} strategy={verticalListSortingStrategy}>
          {cards.map((card) => (
            <TaskCard key={card.id} card={card} onClick={() => onCardClick(card)} />
          ))}
        </SortableContext>
      </Box>
      {showAddCardForm ? (
        <AddCardForm
          listId={list.id}
          onAddCardToList={onAddCardToList}
          onCancel={() => setShowAddCardForm(false)}
        />
      ) : (
        <Button
          fullWidth
          startIcon={<AddIcon />}
          onClick={() => setShowAddCardForm(true)}
          sx={{
            mt: 1,
            justifyContent: 'flex-start',
            color: '#a6b0cf', // Lighter text for dark background
            '&:hover': { backgroundColor: 'rgba(170,180,200,0.1)', color: 'white' }, // Subtle hover effect
            textTransform: 'none', // Keep button text normal case
            py: 1, // Vertical padding
          }}
        >
          Add a card
        </Button>
      )}
    </Paper>
  );
};

export default TaskList;
