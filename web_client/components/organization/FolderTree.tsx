"use client";

import { useOrganizationId } from '@/constants';
import { useUser } from '@/context/UserContext';
import { fetchFolders } from '@/queries';
import { Icon } from '@iconify/react/dist/iconify.js';
import { useQuery } from '@tanstack/react-query';
import React, { useEffect, useState } from 'react';
import Loading from '../shared/Loading';
import Link from 'next/link';

function buildTree(folders: any, parentId: string = "") {
  return folders
    ?.filter((folder: any) => folder.parent_folder_id === parentId)
    ?.map((folder: any) => ({
      ...folder,
      subfolders: buildTree(folders, folder.id)
    }));
}

const FolderTree = ({ tree, level = 0 } : { tree: any, level: number }) => {
  const organization_id = useOrganizationId();
  const orgUrl = `/organization/${organization_id}`
  const [expanded, setExpanded] = useState<{ [key: string]: boolean }>({});

  const toggleExpand = (id: string) => {
    setExpanded(prevState => ({
      ...prevState,
      [id]: !prevState[id]
    }));
  };

  useEffect(() => {
    console.log("Folder Tree: ", tree)
  }, [tree])
  
  return (
    <ul className={level === 0 ? "ml-0" : "ml-4"}>
      {tree?.map((folder: any) => (
        <li key={folder.id}>
          <Link
            href={`${orgUrl}/folder/${folder?.id}`}
            className='flex gap-2'
          >
            <div className='flex gap-1'>
              <Icon
                onClick={() => toggleExpand(folder?.id)}
                icon="material-symbols:arrow-right"
                className={`${expanded[folder.id] && "rotate-90"} ${folder.subfolders.length > 0 ? "opacity-100" : "opacity-0"}`}
                width={24} 
                height={24}
              />
              <Icon icon="flat-color-icons:folder" width={24} height={24} />
            </div>
            {folder?.name}
          </Link>
          {(expanded[folder.id] && folder.subfolders.length > 0) && (
            <FolderTree tree={folder.subfolders} level={level + 1} />
          )}
        </li>
      ))}
    </ul>
  );
};

const FolderTreeComponent = () => {
  const organization_id = useOrganizationId();
  const { data: folders, isLoading: foldersLoading, isError: foldersError } = useQuery({
    queryKey: ['folders', organization_id],
    queryFn: () => fetchFolders(organization_id),
    enabled: !!organization_id
  })

  if (foldersLoading) {
    return <div className='flex justify-center items-center'>
      <Loading />
    </div>
  }

  if (!folders) {
    return <div className='flex justify-center items-center'>
      <p>No folders</p>
    </div>
  }

  const folderTree = buildTree(folders);

  return (
    <div className='py-2'>
      <FolderTree tree={folderTree} level={0} />
    </div>
  );
}

export default FolderTreeComponent;
