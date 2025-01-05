"use client";

import React from 'react';
import FileSummary from '@/components/organization/drive/FileSummary';
import Header from '@/components/organization/drive/Header';
import { fetchRootFiles, fetchRootFolders } from '@/queries';
import { useQuery } from '@tanstack/react-query';
import { useOrganizationId } from '@/constants';
import Manifest from '@/components/shared/Manifest';
import ContextMenu from '@/components/shared/ContextMenu';

const Page = () => {
  const organization_id = useOrganizationId();

  const { data: folders, isLoading: foldersLoading, isError: foldersError } = useQuery({
    queryKey: ['folders', organization_id],
    queryFn: () => fetchRootFolders(organization_id)
  })

  const { data: files, isLoading: filesLoading, isError: filesError } = useQuery({
    queryKey: ['files', organization_id],
    queryFn: () => fetchRootFiles(organization_id)
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
