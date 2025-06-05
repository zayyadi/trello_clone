import React from 'react';
import { Paper, Typography, Box } from '@mui/material';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';

export default function CardItem({ card, onCardClick }) { // Removed index prop
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: `card-${card.id}` }); // Use card.id as the unique ID for sortable context

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    zIndex: isDragging ? 1000 : 'auto',
  };

  return (
    <Paper
      ref={setNodeRef}
      style={style}
      sx={{
        p: 1.5,
        mb: 1,
        backgroundColor: (theme) => theme.palette.background.paper,
        boxShadow: 1,
        '&:hover': {
          boxShadow: 3,
          cursor: 'pointer',
        },
      }}
      {...attributes}
      {...listeners}
      onClick={() => onCardClick(card)}
    >
      <Typography variant="subtitle1" sx={{ fontWeight: 'bold' }}>{card.title}</Typography>
          {card.description && (
            <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
              {card.description}
            </Typography>
          )}
          {/* Add more card details or actions here */}
        </Paper>
  );
}
