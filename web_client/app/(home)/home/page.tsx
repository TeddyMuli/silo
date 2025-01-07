import Link from 'next/link';
import React from 'react';

const links = [
  { name: "login", path: "/login" },
  { name: "register", path: "/register" }
];

const Page = () => {
  return (
    <div className='flex flex-col'>
      <div className='ml-auto'>
        {links?.map((link, index) => (
          <Link key={index} href={link.path} className='capitalize mx-2'>
            {link.name}
          </Link>
        ))}
      </div>
      <div className='justify-center items-center'>
        <p className='text-xl font-semibold'>Silo Homepage</p>
      </div>
    </div>
  );
}

export default Page;
