"use client";

import React, { useEffect, useRef, useState } from 'react';
import { Icon } from "@iconify/react"
import Image from 'next/image';
import { EllipsisVertical } from 'lucide-react';
import { usePathname, useRouter } from 'next/navigation';
import { downloadFile, formatFileSize, useOrganizationId } from '@/constants';
import axios from 'axios';
import Options from '@/components/modals/Options';
import { fetchFolderHierarchyRecursively } from '@/queries';

const FileSummary = ({ object, type } : { object: any, type: string }) => {
  const organization_id = useOrganizationId();
  const orgUrl = `/organization/${organization_id}`
  const router = useRouter();
  const[showOptions, setShowOptions] = useState(false);
  const triggerRef = useRef(null);
  const [folderSize, setFolderSize] = useState("")
  const pathName = usePathname()

  function calculateFolderSize (folder: any): number {
    let totalSize = folder.files?.reduce((acc: number, file: any) => acc + (file.file_size || 0), 0) || 0;

    // Recursively sum the sizes of all subfolders
    folder.subfolders?.forEach((subfolder: any) => {
      totalSize += calculateFolderSize(subfolder);
    });

    return totalSize;
  }

  useEffect(() => {
    const fetchHierarchy = async () => {
      if (type === "folder") {
        try {
          const fetchedHierarchy = await fetchFolderHierarchyRecursively(organization_id, object?.id);
          const folderSize = calculateFolderSize(fetchedHierarchy);
          setFolderSize(formatFileSize(folderSize))
        } catch (error) {
          console.error("Error fetching hierarchy: ", error);
        }  
      } else {
        return
      }
    };

    fetchHierarchy()
  }, [object])
  
  const handleDoubleClick = () => {
    const segments = pathName.split('/');
    const lastSegment = segments.pop() || '';

    if (lastSegment !== "bin") {
      if (type === "folder") {
        router.push(`${orgUrl}/folder/${object?.id}`)
      } else {
        downloadFile(object?.file_path, object?.name)
      }
    } else {
      return;
    }
  }

  return (
    <div
      onDoubleClick={handleDoubleClick}
      className='flex relative py-3 pl-4 border-b justify-between border-neutral-500 items-center dark:hover:bg-white/10 hover:bg-white/60 transition-all duration-100 select-none'
    >
      <Options
        shown={showOptions}
        type={type}
        close={() => setShowOptions(false)}
        object={object}
        triggerRef={triggerRef}
      />
      <div className='flex gap-3 cursor-pointer'>
        <Icon
          icon={`${type === "folder" ? "flat-color-icons:folder" : "tabler:file-filled"}`}
          width={24}
          height={24}
        />
        <p>{object?.name}</p>
      </div>

      <div className='flex gap-24 items-center'>
        <div className='flex gap-2 items-center mr-1 w-24'>
          <Image src="/assets/profile.png" alt='profile' width={24} height={24} />
          <p>me</p>
        </div>

        <div className='mr-8 w-24'>
          <p>Sept 1, 2024</p>
        </div>

        <div className='w-20'>
          <p>{type === "file" ? (formatFileSize(object?.file_size)) : (folderSize)}</p>
        </div>

        <EllipsisVertical
          size={16}
          onClick={() => setShowOptions(!showOptions)}
          className='cursor-pointer'
          ref={triggerRef}
        />
      </div>
    </div>
  );
}

export default FileSummary;
