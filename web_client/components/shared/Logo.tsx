import Link from 'next/link';
import React from 'react';

const Logo = () => {

  return (
    <Link href="/" className='flex gap-1 cursor-pointer w-36 mb-10'>
      <div className=''>
        <h1 className='font-bold text-4xl'>Storify<span className='text-blue-600'>.</span></h1>
      </div>
    </Link>
  );
}

export default Logo;
