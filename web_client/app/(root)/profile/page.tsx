"use client";

import { useUser } from '@/context/UserContext';
import { toast } from '@/hooks/use-toast';
import { handleUpdateUser } from '@/mutations';
import { getUser } from '@/queries';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import Image from 'next/image';
import React, { useEffect } from 'react';
import { useForm } from 'react-hook-form';

type FormData = {
  firstName: string;
  lastName: string;
  phoneNumber: string;
};

const inputs: { name: string; id: keyof FormData }[] = [
  { name: "First Name", id: "firstName" },
  { name: "Last Name", id: "lastName" },
  { name: "Phone Number", id: "phoneNumber" },
];

const Page = () => {
  const { user } = useUser();
  const queryClient = useQueryClient()

  const { data: userData } = useQuery({
    queryKey: ['user'],
    queryFn: () => getUser(user?.user_id),
    enabled: !!user?.user_id
  });

  const { mutate: updateUser } = useMutation({
    mutationKey: ['user'],
    mutationFn: (userData: any) => handleUpdateUser(user?.email, userData),
    onSuccess: (data: any) => {
      if (data.status === 200) {
        toast({
          description: "User updated successfully"
        })
        queryClient.invalidateQueries({ queryKey: ['user'] })
      }
    }
  })

  const {
    register,
    handleSubmit,
    reset,
    getValues,
    formState: { errors, isValid, isDirty }
  } = useForm<FormData>({
    mode: "onChange",
    defaultValues: {
      firstName: "",
      lastName: "",
      phoneNumber: ""
    }
  });

  useEffect(() => {
    reset({
      firstName: userData?.first_name || "",
      lastName: userData?.last_name || "",
      phoneNumber: userData?.phone_number || ""
    });
  }, [userData, reset]);

  const onSubmit = () => {
    const data = getValues()
    const user = {
      first_name: data.firstName,
      last_name: data.lastName,
      phone_number: data.phoneNumber
    }
    updateUser(user)
  };

  return (
    <div>
      <p className='text-2xl font-medium mb-8'>Profile</p>

      <div className='flex flex-col lg:flex-row gap-3 justify-center'>
        <div className='flex flex-col gap-3 justify-center items-center'>
          <Image
            src="/assets/profile.png"
            alt="profile picture"
            width={100}
            height={100}
            className='rounded-full cursor-pointer'
          />
          <label htmlFor="imageUpload">Upload Image</label>
          <input
            type="file"
            id="imageUpload"
            accept="image/*"
            onChange={(e: any) => {
              const file = e?.target?.files[0];
              if (file) {
                // Handle the file upload logic here
                console.log("Selected file:", file);
              }
            }}
          />
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className='flex flex-col gap-3 justify-center items-center'>
          {inputs?.map((input, index) => (
            <div key={index} className='flex flex-col gap-2'>
              <label htmlFor={input.id}>{input.name}</label>
              <input
                type="text"
                id={input.id}
                placeholder={`Enter your ${input.name}`}
                {...register(input.id, { required: `${input.name} is required` })}
                className='p-3 rounded-xl w-64'
              />
              {errors[input.id] && <p className='text-red-500'>{errors[input.id]?.message}</p>}
            </div>
          ))}

          <button
            type="submit"
            className={`px-3 py-2 bg-blue-500 rounded-xl text-xl opacity-100 disabled:opacity-0`}
            disabled={!isValid || !isDirty}
          >
            Save
          </button>
        </form>
      </div>
    </div>
  );
}

export default Page;
