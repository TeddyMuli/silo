"use client";
import Loading from '@/components/shared/Loading';
import { API_URL } from '@/constants';
import { toast } from '@/hooks/use-toast';
import axios from 'axios';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import React, { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';

const Page = () => {
  const [email, setEmail] = useState<string | null>(null);
  const router = useRouter();

  const [timeLeft, setTimeLeft] = useState(300);

  useEffect(() => {
    const timeOTPSent = localStorage.getItem('timeOTPSent');
    if (timeOTPSent) {
      const elapsedTime = Math.floor((Date.now() - parseInt(timeOTPSent)) / 1000);
      setTimeLeft(Math.max(300 - elapsedTime, 0));
    } else {
      router.push("/login")
    }

    const intervalId = setInterval(() => {
      setTimeLeft((prevTime) => Math.max(prevTime - 1, 0));
    }, 1000);

    return () => clearInterval(intervalId);
  }, []);
  
  const formatTime = (seconds: number) => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds < 10 ? '0' : ''}${remainingSeconds}`;
  };

  useEffect(() => {
    const storedEmail = localStorage.getItem("email");
    if (!storedEmail) router.push("/login")
    setEmail(storedEmail);
  }, []);

  const {
    register,
    reset,
    handleSubmit,
    getValues,
    formState: { isValid, isSubmitting, isDirty }
  } = useForm({
    mode: "onChange",
    defaultValues: {
      otp: ""
    }
  })

  function parseJwt(token: string) {
    if (!token) { return; }
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace('-', '+').replace('_', '/');
    return JSON.parse(window.atob(base64));
  }

  function setTokenCookie(token: string) {
    document.cookie = `auth_token=${token}; path=/; Max-Age=${3 * 24 * 60 * 60}; Secure; HttpOnly; SameSite=Lax`;
  }

  const onSubmit = async() => {
    try {
      const { otp } = getValues()
      const to_send = {
        otp: otp,
        email: email
      }

      const response = await axios.post(`${API_URL}/auth/verify`, to_send)

      if (response.status == 200) {
        const token = response.data.token
        if (token) {
          // Store the token in local storage
          localStorage.setItem('token', token);
          setTokenCookie(token);
        }
        console.log(parseJwt(token))
        reset()
        toast({
          description: "Login successful!"
        })
        localStorage.removeItem("timeOTPSent")
        localStorage.removeItem("email")
        router.push("/")
      } else {
        console.error("Error logging in: ", response.data())
        toast({
          description: "There was an error loggin in.",
          variant: "destructive"
        })  
      }
    } catch(error: any) {
      if (error.response.data.error === "Invalid or expired OTP") {
        toast({
          description: "Invalid or expired OTP",
          variant: "destructive"
        })
      } else {
        console.error("Unknown Error: ", error)
      }
    }
  }

  return (
    <div className='flex flex-col justify-center items-center'>
      <p className='text-2xl py-4 font-semibold'>Verify OTP</p>

      <form onSubmit={handleSubmit(onSubmit)} className='flex flex-col gap-4'>
        <div>
          <input
            type="text"
            {...register("otp")}
            className='p-3 w-64 rounded-lg'
            placeholder='Enter your OTP'
          />
        </div>
        <button
          type='submit'
          disabled={!isValid || !isDirty}
          className='flex items-center justify-center px-4 py-2 text-lg font-medium bg-indigo-600 rounded-lg disabled:cursor-not-allowed disabled:bg-indigo-400'
        >
          {isSubmitting && <Loading className='text-white' />}
          Verify
        </button>
      </form>

      <div className='mt-4'>
        {timeLeft <= 0 ? (
          <p>OTP expired <Link href="/login">Login</Link></p>
        ) : (
          <p>Expires in: <span className='text-red-500'>{formatTime(timeLeft)}</span></p>
        )}
      </div>

      <p className='text-sm mt-4'>Don't have an OTP? <Link href="/login">Log in</Link></p>
    </div>
  );
}

export default Page;
