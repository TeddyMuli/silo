"use client";

import { LayoutDashboard, Plus } from 'lucide-react';
import Link from 'next/link';
import React, { useState } from 'react';
import Logo from './Logo';
import Create from '../modals/Create';
import { useUser } from '@/context/UserContext';
import { useQuery } from '@tanstack/react-query';
import { fetchOrganizations } from '@/queries';

const LeftSideBar = () => {
  const[showCreate, setShowCreate] = useState<boolean>(false);

  const { user } = useUser()  
  const { data: organizations, isLoading, isError } = useQuery({
    queryKey: ['organizations', user?.user_id],
    queryFn: () => fetchOrganizations(user),
    enabled: !!user?.user_id
  });

  return (
    <section className='overflow-y-auto h-screen w-52 items-center justify-center'>
      <Logo />
      <Create
        shown={showCreate}
        object='Organization'
        close={() => setShowCreate(false)}
      />
      <div>
        <div
          onClick={() => setShowCreate(true)}
          className='flex gap-2 bg-indigo-600 p-2 rounded-xl text-white cursor-pointer'
        >
          <Plus />
          <p className='font-medium'>Create Organization</p>
        </div>

        <div className='my-4'>
          <div className='flex gap-3'>
            <LayoutDashboard />
            <p>Dashboard</p>
          </div>
        </div>

        <p className='mt-4'>Organizations</p>
        {organizations?.map((organization: any, index: number) => (
          <div key={index} className='mt-3 p-2 rounded-xl bg-neutral-300'>
            <Link
              key={index}
              href={`/organization/${organization.organization_id}`}
              className='text-black font-medium'
            >
              {organization.name}
            </Link>
          </div>
        ))}
      </div>
    </section>
  );
}

export default LeftSideBar;
