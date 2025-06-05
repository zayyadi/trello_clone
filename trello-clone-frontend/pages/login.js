import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useRouter } from 'next/router';
import { useForm } from 'react-hook-form';
import { loginUser, resetAuthStatus } from '../features/auth/authSlice';
import { TextField, Button, Container, Typography, Paper, Box, Alert } from '@mui/material';
import Link from 'next/link';

export default function LoginPage() {
  const dispatch = useDispatch();
  const router = useRouter();
  const { user, status, error } = useSelector((state) => state.auth);

  const { register, handleSubmit, formState: { errors } } = useForm();

  useEffect(() => {
    if (user) {
      router.push('/'); // Redirect if already logged in
    }
    // Reset status on component mount or when error changes to allow re-submission
    return () => {
        if (status === 'failed') {
            dispatch(resetAuthStatus());
        }
    }
  }, [user, router, dispatch, status]);

  const onSubmit = (data) => {
    dispatch(loginUser(data));
  };

  return (
    <Container component="main" maxWidth="xs">
      <Paper elevation={3} sx={{ marginTop: 8, padding: 4, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Typography component="h1" variant="h5">
          Sign In
        </Typography>
        <Box component="form" onSubmit={handleSubmit(onSubmit)} sx={{ mt: 1 }}>
          <TextField
            margin="normal"
            required
            fullWidth
            id="email"
            label="Email Address"
            autoComplete="email"
            autoFocus
            {...register("email", { required: "Email is required" })}
            error={!!errors.email || (status === 'failed' && error?.toLowerCase().includes('email'))}
            helperText={errors.email?.message}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            label="Password"
            type="password"
            id="password"
            autoComplete="current-password"
            {...register("password", { required: "Password is required" })}
            error={!!errors.password || (status === 'failed' && error?.toLowerCase().includes('password'))}
            helperText={errors.password?.message}
          />
           {status === 'failed' && error && (
             <Alert severity="error" sx={{ width: '100%', mt: 2 }}>{error}</Alert>
           )}
          <Button
            type="submit"
            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 2 }}
            disabled={status === 'loading'}
          >
            {status === 'loading' ? 'Signing In...' : 'Sign In'}
          </Button>
          <Box textAlign="center">
            <Link href="/register" passHref>
              <Typography component="a" variant="body2">
                {"Don't have an account? Sign Up"}
              </Typography>
            </Link>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
}