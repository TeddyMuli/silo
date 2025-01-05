"use client";

import React from 'react';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import Loading from '@/components/shared/Loading';
import Link from 'next/link';
import axios from 'axios';
import { API_URL } from '@/constants';
import { useRouter } from 'next/navigation';
import { useToast } from '@/hooks/use-toast';
import { useGoogleReCaptcha } from 'react-google-recaptcha-v3';

const validationSchema = z.object({
  email: z.string().email(),
  first_name: z.string().min(1, "This field is required!"),
  last_name: z.string().min(1, "This field is required!"),
  phone_number: z.string().min(1, "This field is required!"),
  password: z.string().min(8, "Password too short!")
    .regex(/[a-zA-Z]/, "Password must contain at least one letter")
    .regex(/[0-9]/, "Password must contain at least one number")
    .regex(/[^a-zA-Z0-9]/, "Password must contain at least one special character"),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match!",
  path: ["confirmPassword"],
});

const Page = () => {
  const router = useRouter();
  const { executeRecaptcha } = useGoogleReCaptcha();
  const { toast } = useToast();

  const {
    reset,
    getValues,
    register,
    handleSubmit,
    formState: { errors, isValid, isSubmitting }
  } = useForm({
    mode: "onChange",
    resolver: zodResolver(validationSchema),
    defaultValues: {
      first_name: "",
      last_name: "",
      phone_number: "",
      email: "",
      password: "",
      confirmPassword: ""
    }
  });

  const handleRegistration = async () => {
    if (!executeRecaptcha) {
      console.error('ReCAPTCHA not available');
      return;
    }

    try {
      const gRecaptchaToken = await executeRecaptcha('registerSubmit');

      const recaptchaResponse = await axios.post('/api/recaptchaVerify', {
        gRecaptchaToken,
      }, {
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
        },
      });

      if (!recaptchaResponse.data.success) {
        console.error(`ReCaptcha verification failed with score: ${recaptchaResponse.data.score}`);
        toast({
          description: "ReCaptcha verification failed!",
          variant: "destructive"
        });
        return;
      }

      const data = getValues()
      const { confirmPassword, ...body } = data;
      const response = await axios.post(`${API_URL}/auth/register`, body);

      if (response.status === 201) {
        reset();
        toast({
          description: "Account created successfully!"
        });
        router.push("/login");
      }
    } catch (error: any) {
      if (axios.isAxiosError(error)) {
        if (error.response) {
          if (error.response.data.error === 'Email already exists') {
            toast({
              description: "Email already exists!",
              variant: "destructive"
            });
          } else {
            console.error('Some other 401 error occurred:', error.response.data.error);
          }
        } else {
          console.log('The request was made, but no response was received:', error.message);
        }
      } else {
        console.log('Unknown Error:', error.message);
      }
    }
  };

  return (
    <div className='flex flex-col justify-center items-center'>
      <p className='py-4 text-xl font-semibold'>Register</p>

      <form onSubmit={handleSubmit(handleRegistration)} className='flex flex-col gap-4 text-white'>
        <div>
          <input
            type="text"
            className='p-3 text-white rounded-lg w-64'
            placeholder="Enter firstname"
            {...register("first_name")}
          />
          {errors.first_name && <div className="text-red-500 text-sm font-medium pt-2">{errors.first_name.message}</div>}
        </div>
        <div>
          <input
            type="text"
            className='p-3 text-white rounded-lg w-64'
            placeholder="Enter lastname"
            {...register("last_name")}
          />
          {errors.last_name && <div className="text-red-500 text-sm font-medium pt-2">{errors.last_name.message}</div>}
        </div>
        <div>
          <input
            type="text"
            className='p-3 text-white rounded-lg w-64'
            placeholder="Enter email"
            {...register("email")}
          />
          {errors.email && <div className="text-red-500 text-sm font-medium pt-2">{errors.email.message}</div>}
        </div>
        <div>
          <input
            type="text"
            className='p-3 text-white rounded-lg w-64'
            placeholder="Enter phone number"
            {...register("phone_number")}
          />
          {errors.phone_number && <div className="text-red-500 text-sm font-medium pt-2">{errors.phone_number.message}</div>}
        </div>
        <div>
          <input
            type="password"
            className='p-3 text-white rounded-lg w-64'
            placeholder="Enter password"
            {...register("password")}
          />
          {errors.password && <div className="text-red-500 text-sm font-medium pt-2">{errors.password.message}</div>}
        </div>
        <div>
          <input
            type="password"
            className='p-3 text-white rounded-lg w-64'
            placeholder="Confirm password"
            {...register("confirmPassword")}
            autoComplete='new-password'
          />
          {errors.confirmPassword && <div className="text-red-500 text-sm font-medium pt-2">{errors.confirmPassword.message}</div>}
        </div>
        <button type='submit' disabled={!isValid} className='flex items-center justify-center gap-2 submit px-4 py-2 text-lg font-medium bg-indigo-600 rounded-lg disabled:cursor-not-allowed disabled:bg-indigo-400'>
          {isSubmitting && <Loading />}
          Register
        </button>
      </form>
      <Link href="/login" className='text-sm mt-4'>Already have an account? Login</Link>
    </div>
  );
}

export default Page;
