import React from 'react';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Paper, Typography, Box, Chip, Avatar } from '@mui/material';
import EventIcon from '@mui/icons-material/Event';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import RadioButtonUncheckedIcon from '@mui/icons-material/RadioButtonUnchecked';
import HourglassEmptyIcon from '@mui/icons-material/HourglassEmpty';
import DescriptionOutlinedIcon from '@mui/icons-material/DescriptionOutlined'; // For description
import ChatBubbleOutlineIcon from '@mui/icons-material/ChatBubbleOutline'; // For comments
import AttachFileIcon from '@mui/icons-material/AttachFile'; // For attachments

// Helper to format date
const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
}; 

const getStatusChipProps = (status) => {
    switch (status) {
        case 'done':
            return { label: 'Done', color: 'success', icon: <CheckCircleIcon fontSize="small"/> };
        case 'pending':
            return { label: 'Pending', color: 'warning', icon: <HourglassEmptyIcon fontSize="small"/> };
        case 'to_do':
        default:
            return { label: 'To Do', color: 'default', icon: <RadioButtonUncheckedIcon fontSize="small"/> };
        // case 'undone': // If you use this distinct from to_do
        //     return { label: 'Undone', color: 'error', icon: <HighlightOffIcon /> };
    }
};


const TaskCard = ({ card, onClick }) => {
  const { id, title, listId, dueDate, status, description, comments, attachments, assignedUser } = card; // Destructure new fields

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id, data: { type: 'CARD', parentListId: listId, cardId: id, cardTitle: title, cardData: card } });

  const style = {
    transform: CSS.Translate.toString(transform),
    transition,
    opacity: isDragging ? 0.7 : 1,
    zIndex: isDragging ? 100 : 'auto',
    cursor: 'grab',
  };

  const chipProps = getStatusChipProps(status);

  return (
    <Paper
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      onClick={() => onClick(card)} // Propagate click for modal
      sx={{
        p: 1.5,
        mb: 1,
        backgroundColor: 'background.paper',
        boxShadow: isDragging ? '0px 5px 15px rgba(0,0,0,0.3)' : '0px 1px 2px rgba(0,0,0,0.1)',
        '&:hover': { backgroundColor: '#f8f8f8', cursor: 'pointer' }
      }}
    >
      {/* Labels/Tags */}
      {card.labels && card.labels.length > 0 && (
        <Box sx={{ display: 'flex', gap: 0.5, mb: 1, flexWrap: 'wrap' }}>
          {card.labels.map((label, index) => (
            <Box
              key={index}
              sx={{
                height: 8,
                width: 40,
                borderRadius: '3px',
                backgroundColor: label.color || '#ccc', // Use label color or a default
              }}
            />
          ))}
        </Box>
      )}

      <Typography variant="body1" sx={{ mb: 1, fontWeight: 'bold' }}>{title}</Typography>

      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, flexWrap: 'wrap', mt: 0.5 }}>
        {dueDate && (
          <Chip
            icon={<EventIcon fontSize="small" />}
            label={formatDate(dueDate)}
            size="small"
            variant="outlined"
            color={new Date(dueDate) < new Date() && status !== 'done' ? 'error' : 'default'}
            sx={{ height: 24 }}
          />
        )}
        {description && description.trim() !== '' && (
          <DescriptionOutlinedIcon fontSize="small" sx={{ color: 'text.secondary' }} />
        )}
        {comments && comments.length > 0 && (
          <Box sx={{ display: 'flex', alignItems: 'center', color: 'text.secondary', fontSize: '0.875rem' }}>
            <ChatBubbleOutlineIcon fontSize="small" sx={{ mr: 0.5 }} />
            {comments.length}
          </Box>
        )}
        {attachments && attachments.length > 0 && (
          <Box sx={{ display: 'flex', alignItems: 'center', color: 'text.secondary', fontSize: '0.875rem' }}>
            <AttachFileIcon fontSize="small" sx={{ mr: 0.5 }} />
            {attachments.length}
          </Box>
        )}
      </Box>

      {assignedUser && (
        <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 1 }}>
          <Avatar sx={{ width: 28, height: 28, bgcolor: 'primary.main', fontSize: '0.8rem' }}>
            {assignedUser.username ? assignedUser.username.charAt(0).toUpperCase() : ''}
          </Avatar>
        </Box>
      )}
    </Paper>
  );
};

export default TaskCard;
