"use client";

import React from 'react';
import { Icon } from '@iconify/react';
import { EllipsisVertical } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { fetchFolderHierarchy } from '@/queries';
import { useFolderId, useOrganizationId } from '@/constants';
import Loading from '@/components/shared/Loading';
import { useRouter } from 'next/navigation';

const Header = () => {
  const router = useRouter()
  const folderId = useFolderId()

  const organization_id = useOrganizationId();
  const orgUrl = `/organization/${organization_id}`

  const { data: folderHierarchy, isLoading: folderLoading, error: folderError } = useQuery({
    queryKey: ['folderHierarchy', folderId],
    queryFn: () => fetchFolderHierarchy(folderId),
    enabled: !!folderId
  });

  if (folderLoading) return <div className='flex items-center justify-center'>
    <Loading />
  </div>;

  if (folderError) {
    console.error("Error: ", folderError.message)
    return
  };

  return (
    <div className='flex flex-col gap-8 overflow-hidden'>
      <div className='flex items-center'>
        <div className='flex cursor-pointer py-1 px-3 rounded-full hover:bg-black/10 dark:hover:bg-white/30 transition-all duration-200'>
          <p className='text-2xl text-black/30 dark:text-white/70' onClick={() => router.push(`${orgUrl}/drive`)}>Drive</p>
          <div className='mt-1'>
            {!folderHierarchy ? (
              <Icon icon="material-symbols:arrow-right" className='rotate-90' width={24} height={24} />
            ) : (
              <Icon icon="ic:round-chevron-right" className='mt-[2px]' height={24} />
            )}
          </div>
        </div>
        {folderHierarchy?.map((folder, index) => (
          <div key={folder.id} className='flex gap-1 place-items-center'>
            {index < folderHierarchy?.length - 1 ? (
              <div className='flex justify-center items-center cursor-pointer py-1 px-3 rounded-full hover:bg-black/10 dark:hover:bg-white/30 transition-all duration-200'>
                <p
                  className='text-2xl text-black/30 dark:text-white/70'
                  onClick={() => router.push(`${orgUrl}/folder/${folder?.id}`)}
                >
                  {folder.name}
                </p>
                <Icon icon="ic:round-chevron-right" height={24} className='mt-1' />
              </div>
            ) : (
              <div className='flex justify-center items-center cursor-pointer py-1 px-3 rounded-full hover:bg-black/10 dark:hover:bg-white/30 transition-all duration-200'>
                <p className='text-2xl'>{folder.name}</p>
                <Icon icon="material-symbols:arrow-right" className='rotate-90 mt-1' width={24} height={24} />
              </div>
            )}
          </div>
        ))}
      </div>

      <div className='flex border-b border-neutral-500 pb-3'>
        <div className='flex items-center gap-4'>
          <p className='text-lg'>Name</p>
          <Icon icon="charm:arrow-down" width={20} height={20} />
        </div>
        <div className='flex items-center gap-24 ml-auto'>
          <p className='w-24'>Owner</p>
          <div className='flex items-center w-32'>
            <p>Last Modified</p>
            <Icon icon="material-symbols:arrow-right" className='rotate-90 cursor-pointer' width={24} height={24} />
          </div>
          <p className='w-20'>Size</p>
          <EllipsisVertical size={16} />
        </div>
      </div>
    </div>
  );
}

export default Header;
