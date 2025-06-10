import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useRouter } from 'next/router';
import { useForm } from 'react-hook-form';
import { registerUser, resetAuthStatus } from '../features/auth/authSlice';
import { TextField, Button, Container, Typography, Paper, Box, Alert } from '@mui/material';
import Link from 'next/link';

export default function RegisterPage() {
  const dispatch = useDispatch();
  const router = useRouter();
  const { user, status, error } = useSelector((state) => state.auth);

  const { register, handleSubmit, formState: { errors } } = useForm();

  useEffect(() => {
    if (user) {
      router.push('/'); // Redirect if already logged in
    } else if (status === 'succeeded') {
      router.push('/login'); // Redirect to login after successful registration
      dispatch(resetAuthStatus()); // Reset status after successful registration
    }
    return () => {
        if (status === 'failed') {
            dispatch(resetAuthStatus());
        }
    }
  }, [user, router, dispatch, status]);

  const onSubmit = (data) => {
    dispatch(registerUser(data));
  };

  return (
    <Container component="main" maxWidth="xs">
      <Paper elevation={3} sx={{ marginTop: 8, padding: 4, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Typography component="h1" variant="h5">
          Sign Up
        </Typography>
        <Box component="form" onSubmit={handleSubmit(onSubmit)} sx={{ mt: 1 }}>
          <TextField
            margin="normal"
            required
            fullWidth
            id="username"
            label="Username"
            autoComplete="username"
            autoFocus
            {...register("username", { required: "Username is required" })}
            error={!!errors.username}
            helperText={errors.username?.message}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            id="email"
            label="Email Address"
            autoComplete="email"
            {...register("email", { 
              required: "Email is required",
              pattern: {
                value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$/i,
                message: "Invalid email address"
              }
            })}
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
            autoComplete="new-password"
            {...register("password", { 
              required: "Password is required",
              minLength: {
                value: 6,
                message: "Password must be at least 6 characters"
              }
            })}
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
            {status === 'loading' ? 'Signing Up...' : 'Sign Up'}
          </Button>
          <Box textAlign="center">
            <Link href="/login" passHref>
              <Typography component="a" variant="body2">
                {"Already have an account? Sign In"}
              </Typography>
            </Link>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
}
