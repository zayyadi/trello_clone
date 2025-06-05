import React, { useEffect, useState } from 'react'; // Import useEffect and useState
import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material';
import Link from 'next/link';
import { useSelector, useDispatch } from 'react-redux';
import { logout } from '../../features/auth/authSlice';
import { useRouter } from 'next/router';

export default function Navbar() {
  const { user } = useSelector((state) => state.auth);
  const dispatch = useDispatch();
  const router = useRouter();

  const [mounted, setMounted] = useState(false); // State to track if component is mounted

  useEffect(() => {
    setMounted(true); // Set mounted to true after component mounts on client
  }, []);

  const handleLogout = () => {
    dispatch(logout());
    router.push('/login');
  };

  return (
    <AppBar position="static">
      <Toolbar>
        <Link href="/" passHref>
          <Typography variant="h6" sx={{ flexGrow: 1, color: 'white', textDecoration: 'none', cursor: 'pointer' }}>
            Trello Clone
          </Typography>
        </Link>
        {mounted && ( // Only render user-specific content after component is mounted
          user ? (
            <Box>
                <Link href="/calendar" passHref legacyBehavior>
                    <Button color="inherit" component="a" sx={{mr:1}}>Calendar</Button>
                </Link>
              <Typography variant="subtitle1" component="span" sx={{ mr: 2 }}>
                Hi, {user.username}
              </Typography>
              <Button color="inherit" onClick={handleLogout}>Logout</Button>
            </Box>
          ) : (
            <Box>
              <Link href="/login" passHref><Button color="inherit" component="a">Login</Button></Link>
              <Link href="/register" passHref><Button color="inherit" component="a">Register</Button></Link>
            </Box>
          )
        )}
      </Toolbar>
    </AppBar>
  );
}
