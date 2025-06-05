// pages/calendar.js
import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useRouter } from 'next/router';
import { fetchUserBoards, selectUserBoards, selectUserBoardsStatus } from '../features/boards/boardsSlice';
import { selectAuthToken } from '../features/auth/authSlice';
import { Container, Typography, List, ListItem, ListItemText, CircularProgress, Paper, Box, Alert } from '@mui/material';
import Link from 'next/link';

const formatDate = (dateString) => {
    if (!dateString) return 'No due date';
    return new Date(dateString).toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' });
};

export default function CalendarPage() {
    const dispatch = useDispatch();
    const router = useRouter();
    const token = useSelector(selectAuthToken);
    const boards = useSelector(selectUserBoards);
    const boardsStatus = useSelector(selectUserBoardsStatus);

    useEffect(() => {
        if (!token) {
            router.push('/login');
        } else if (boardsStatus === 'idle' || boards.length === 0) {
            dispatch(fetchUserBoards());
        }
    }, [token, router, dispatch, boardsStatus, boards.length]);

    const cardsWithDueDates = [];
    if (boardsStatus === 'succeeded') {
        boards.forEach(board => {
            (board.lists || []).forEach(list => {
                (list.cards || []).forEach(card => {
                    if (card.dueDate) {
                        cardsWithDueDates.push({ ...card, boardName: board.name, boardId: board.id, listName: list.name });
                    }
                });
            });
        });
        cardsWithDueDates.sort((a, b) => new Date(a.dueDate) - new Date(b.dueDate));
    }

    if (!token) return <Box sx={{display: 'flex', justifyContent: 'center', mt:5}}><CircularProgress /></Box>;

    return (
        <Container maxWidth="md" sx={{ mt: 4 }}>
            <Typography variant="h4" component="h1" gutterBottom>
                Task Calendar (Due Dates)
            </Typography>

            {boardsStatus === 'loading' && <CircularProgress />}
            {boardsStatus === 'failed' && <Alert severity="error">Could not load board data.</Alert>}
            
            {boardsStatus === 'succeeded' && (
                cardsWithDueDates.length === 0 ? (
                    <Typography>No tasks with due dates found.</Typography>
                ) : (
                    <Paper elevation={2}>
                        <List>
                            {cardsWithDueDates.map(card => (
                                <ListItem key={card.id} divider>
                                    <ListItemText
                                        primary={<Link href={`/board/${card.boardId}?card=${card.id}`} passHref legacyBehavior><a>{card.title}</a></Link>}
                                        secondary={
                                            <>
                                                <Typography component="span" variant="body2" color="textPrimary">
                                                    Due: {formatDate(card.dueDate)}
                                                </Typography>
                                                <br />
                                                <Typography component="span" variant="caption" color="textSecondary">
                                                    In "{card.listName}" on board "{card.boardName}" - Status: {card.status}
                                                </Typography>
                                            </>
                                        }
                                    />
                                </ListItem>
                            ))}
                        </List>
                    </Paper>
                )
            )}
        </Container>
    );
}