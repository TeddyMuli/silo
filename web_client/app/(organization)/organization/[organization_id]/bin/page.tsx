"use client";

import DeleteModal from '@/components/modals/DeleteModal';
import FileSummary from '@/components/organization/drive/FileSummary';
import { useOrganizationId } from '@/constants';
import { handleEmptyBin } from '@/mutations';
import { fetchDeletedFiles, fetchDeletedFolders } from '@/queries';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import React, { useState } from 'react';

const Page = () => {
  const organization_id = useOrganizationId();
  const [showDelete, setShowDelete] = useState(false);
  const queryClient = useQueryClient();
   
  const { data: folders, isLoading: foldersLoading, isError: foldersError } = useQuery({
    queryKey: ['folders', organization_id],
    queryFn: () => fetchDeletedFolders(organization_id)
  })

  const { data: files, isLoading: filesLoading, isError: filesError } = useQuery({
    queryKey: ['files', organization_id],
    queryFn: () => fetchDeletedFiles(organization_id)
  })
  const mutationKey = ['emptyTrash', organization_id, 'folders', 'files'];

  const { mutate: emptyTrashMutation } = useMutation({
    mutationKey: mutationKey,
    mutationFn: () => handleEmptyBin(organization_id),
    onSuccess: () => {
      queryClient.invalidateQueries({queryKey: ['folders', organization_id]});
      queryClient.invalidateQueries({queryKey: ['files', organization_id]});
    }
  })

  function handleEmptyTrash () {
    emptyTrashMutation();
  }

  return (
    <div>
      <DeleteModal
        shown={showDelete}
        close={() => setShowDelete(false)}
        object='All items in trash'
        onConfirm={handleEmptyTrash}
      />
      <p className='text-2xl font-medium p-3'>Recycle Bin</p>

      <div className='flex p-3 bg-neutral-300 dark:bg-neutral-600 rounded-lg'>
        <p>Items in bin will be deleted after 30 days</p>
        <button
          className='ml-auto'
          onClick={() => setShowDelete(true)}
        >
          <p className='text-red-600 font-semibold'>Empty bin</p>
        </button>
      </div>

      {folders?.map((folder: any, index: number) => (
        <FileSummary key={index} object={folder} type='folder' />
      ))}
      {files?.map((file: any, index: number) => (
        <FileSummary key={index} object={file} type='file' />
      ))}

      {(!folders && !files) && (
        <div className='flex p-3 items-center justify-center text-2xl font-medium'>No items in recycle bin</div>
      )}
    </div>
  );
}

export default Page;
