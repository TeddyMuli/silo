"use client";

import React, { useEffect, useState } from 'react';
import { formatFileSize, generateManifestForDownload, transformHierarchy, useFolderId, useOrganizationId } from '@/constants';
import { fetchFolder, fetchFolderHierarchyRecursively } from '@/queries';
import { Icon } from '@iconify/react/dist/iconify.js';

const Manifest = ({ folders, files } : { folders: any, files: any }) => {
  const [hierarchy, setHierarchy] = useState<any>(null);
  const [manifestSize, setManifestSize] = useState<number>(0);

  const organization_id = useOrganizationId();
  const folderId = useFolderId() || "root"
  const getName = async () => {
    if (folderId) {
      const folderData: any = await fetchFolder(folderId);
      const folderName: string = folderData.name || "root";
      return folderName
    } else {
      return null
    }
  }

  useEffect(() => {
    const fetchHierarchy = async () => {
      try {
        const fetchedHierarchy = await fetchFolderHierarchyRecursively(organization_id, folderId);
        const transformedHierarchy = transformHierarchy(fetchedHierarchy);
        setHierarchy(fetchedHierarchy);
        
        const jsonString = JSON.stringify(transformedHierarchy);
        const sizeInBytes = new Blob([jsonString]).size;
        setManifestSize(sizeInBytes);
      } catch (error) {
        console.error("Error fetching hierarchy: ", error);
      }
    };

    fetchHierarchy()
  }, [folders, files])

  const handleGenerateManifest = async () => {
    if (hierarchy) {
      const folderName = await getName() || "drive"
      generateManifestForDownload(hierarchy, folderName);
    } else {
      console.error("Hierarchy is not available");
    }
  };

  return (
    <div
        className='flex p-3 gap-4 dark:hover:bg-white/10 hover:bg-white/60 cursor-pointer select-none'
        onDoubleClick={handleGenerateManifest}
      >
        <div className='flex gap-4'>
          <Icon icon="vscode-icons:file-type-light-json" height={24} width={24} />
          <p>Manifest</p>
        </div>
        <p className='ml-auto mr-24 w-20'>{formatFileSize(manifestSize)}</p>
      </div>

  );
}

export default Manifest;
