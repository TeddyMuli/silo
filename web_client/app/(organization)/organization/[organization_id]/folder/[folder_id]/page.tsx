"use client";

import React from 'react';
import FileSummary from '@/components/organization/drive/FileSummary';
import Header from '@/components/organization/drive/Header';
import ContextMenu from '@/components/shared/ContextMenu';
import Manifest from '@/components/shared/Manifest';
import { useFolderId, useOrganizationId } from '@/constants';
import { fetchChildFiles, fetchChildrenFolders } from '@/queries';
import { useQuery } from '@tanstack/react-query';

const Page = () => {
  const organization_id = useOrganizationId();
  const parent_folder_id = useFolderId();

  const { data: folders, isLoading: foldersLoading, error: foldersError } = useQuery({
    queryKey: ['folders', parent_folder_id],
    queryFn: () => fetchChildrenFolders(organization_id, parent_folder_id)
  });
  
  const { data: files, isLoading: filesLoading, error: filesError } = useQuery({
    queryKey: ['files', parent_folder_id],
    queryFn: () => fetchChildFiles(organization_id, parent_folder_id)
  })

  return (
    <div className='h-full'>
      <Header />
      {folders?.map((folder: any, index: number) => (
        <FileSummary key={index} object={folder} type='folder' />
      ))}
      {files?.map((file: any, index: number) => (
        <FileSummary key={index} object={file} type='file' />
      ))}

      <Manifest folders={folders} files={files} />
      <ContextMenu />
    </div>
  );
}

export default Page;
