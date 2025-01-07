"use client";
import Loading from '@/components/shared/Loading';
import { API_URL } from '@/constants';
import { useToast } from '@/hooks/use-toast';
import { zodResolver } from '@hookform/resolvers/zod';
import axios from 'axios';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import React from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

const validationSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8, "Password too short!"),
})

const Page = () => {
  const router = useRouter()
  const { toast } = useToast()

  const {
    register,
    reset,
    getValues,
    handleSubmit,
    formState: { errors, isValid, isSubmitting }
  } = useForm({
    mode: "onChange",
    resolver: zodResolver(validationSchema),
    defaultValues: {
      email: "",
      password: ""
    }
  });

  const handleLogin = async () => {
    try {
      const body = getValues()
      localStorage.setItem('email', body.email);

      
      const response = await axios.post(`${API_URL}/auth/login`, body)
      if (response.status == 200) {
        reset()
        router.push("/verify")
        toast({
          description: "Enter OTP sent to your email"
        })
      } else {
        console.error("Error logging in: ", response.data())
        toast({
          description: "There was an error loggin in.",
          variant: "destructive"
        })  
      }
    } catch (error: any) {
      if (axios.isAxiosError(error)) {
        if (error.response) {
          if (error.response.status === 401 && error.response.data.error === 'Invalid credentials') {
            toast({
              description: "Invalid credentials!",
              variant: "destructive"
            });
          } else {
            console.error('Some other 401 error occurred:', error.response.data.error);
          }

          if (error.response.status === 400 && error.response.data.error === "User doesn't exist!") {
            toast({
              description: "User doesn't exist!",
              variant: "destructive"
            });
          }
        } else {
          console.log('The request was made, but no response was received:', error.message);
        }
      } else {
        console.log('Unknown Error:', error.message);
      }
    }
  }
  
  return (
    <div className='flex flex-col justify-center items-center'>
      <p className='text-2xl py-4 font-semibold'>Login</p>

      <form onSubmit={handleSubmit(handleLogin)} className='flex flex-col gap-4'>
        <div>
          <input
            type="email"
            {...register("email")}
            className='p-3 w-64 rounded-lg'
            placeholder='Enter your email'
          />
          {errors.email && <div className="text-red-500 text-sm font-medium pt-2">{errors.email.message}</div>}
        </div>
        <div>
          <input
            type="password"
            {...register("password")}
            className='p-3 w-64 rounded-lg'
            placeholder='Enter your password'
            autoComplete='current-password'
          />
          {errors.password && <div className="text-red-500 text-sm font-medium pt-2">{errors.password.message}</div>}
        </div>

        <button
          type='submit'
          disabled={!isValid}
          className='flex items-center justify-center px-4 py-2 text-lg font-medium bg-indigo-600 rounded-lg disabled:cursor-not-allowed disabled:bg-indigo-400'
        >
          {isSubmitting && <Loading className='text-white' />}
          Login
        </button>
      </form>
      <Link href="/register" className='text-sm mt-4'>Don't have an account? Sign up</Link>
    </div>
  );
}

export default Page;
