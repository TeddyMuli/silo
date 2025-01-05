"use client";

import { useTheme } from 'next-themes';
import Image from 'next/image';
import Link from 'next/link';
import React, { useEffect, useState } from 'react';

const Logo = () => {
  const { setTheme, resolvedTheme } = useTheme()
  const [image, setImage] = useState("/assets/logo_dark.png")

  useEffect(() => {
    if (resolvedTheme === "dark") {
      setImage("/assets/logo_dark.png")
    } else {
      setImage("/assets/logo_light.png")
    }
  }, [resolvedTheme])

  return (
    <Link href="/" className='flex gap-1 cursor-pointer w-36 mb-10'>
      <div className=''>
        <Image
          src={image}
          alt="logo" 
          width={150}
          height={150}
          className='object-cover'
        />
      </div>
    </Link>
  );
}

export default Logo;
