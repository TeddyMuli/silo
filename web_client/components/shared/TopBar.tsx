"use client";

import React, { useState } from 'react';
import ThemeSwitch from './ThemeSwitch';
import Image from 'next/image';
import { Search } from 'lucide-react';
import Profile from '../modals/Profile';

const TopBar = () => {
  const [showProfile, setShowProfile] = useState(false);

  return (
    <section className='flex gap-3 justify-center items-center'>
      <Profile
        shown={showProfile}
        close={() => setShowProfile(false)}
      />

      <div className='flex gap-3 rounded-full w-2/3 bg-neutral-300 dark:bg-neutral-700 p-3'>
        <Search className='cursor-pointer' />
        <input
          type="text"
          placeholder='Search Aethly'
          className='outline-none bg-transparent w-full dark:placeholder:text-white placeholder:text-black placeholder:font-medium'
        />
      </div>

      <div className='ml-auto flex justify-center items-center gap-3'>
        <ThemeSwitch />
        <div className=''>
          <Image
            src="/assets/profile.png"
            alt="profile picture"
            width={48}
            height={48}
            className='rounded-full cursor-pointer'
            onClick={() => setShowProfile(true)}
          />
        </div>
      </div>
    </section>
  );
}

export default TopBar;
