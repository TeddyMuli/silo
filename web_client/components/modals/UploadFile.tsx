"use client";

import React, { useEffect, useRef, useState } from 'react';
import { useFolderId, useOrganizationId } from '@/constants';
import { useToast } from '@/hooks/use-toast';
import { handleCreateFolder, handleUploadFile } from '@/mutations';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { FileUp, FolderPlus, FolderUp } from 'lucide-react';

const UploadFile = (
  { shown, close, setShowCreate, position }
  :
  { shown: boolean, close: () => void, setShowCreate: (e: boolean) => void, position: { x: number, y: number } }
) => {
  const { toast } = useToast();
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const folderInputRef = useRef<HTMLInputElement | null>(null);
  const organization_id = useOrganizationId();
  const folder_id = useFolderId();
  const queryClient = useQueryClient()

  const [adjustedPosition, setAdjustedPosition] = useState({ x: 0, y: 0 });
  const menuRef = useRef<HTMLDivElement>(null);

  const fileMutationKey = folder_id ? ['files', folder_id] : ['files', organization_id]

  const { mutate: mutateFile } = useMutation({
    mutationKey: fileMutationKey,
    mutationFn: ({ file }: { file: FormData, isFolderUpload: boolean }) => handleUploadFile(file),
    onSuccess: (data, variables) => {
      if (data?.status === 200) {
        if (!variables.isFolderUpload) {
          toast({
            description: "File uploaded successfully!",
          });
        }
        queryClient.invalidateQueries({ queryKey: fileMutationKey })
      } else {
        if (!variables.isFolderUpload) {
          toast({
            description: "File upload failed",
            variant: "destructive",
          });
        }
      }
    },
    onError: (error, variables) => {
      if (!variables.isFolderUpload) {
        toast({
          description: "File upload failed",
          variant: "destructive",
        });
      }
    },
  });

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files && files.length > 0) {
      const fileArray = Array.from(files);
      fileArray.forEach((file) => {
        const formData = new FormData();
        formData.append('file', file);
        formData.append('organization_id', organization_id);
        formData.append('folder_id', folder_id);

        mutateFile({ file: formData, isFolderUpload: false });
      });
    }
  };

  const createFolderStructure = async (
    folderStructure: any,
    parentFolderId: string | null,
  ) => {
    const promises = Object.entries(folderStructure)?.map(async ([name, content]) => {
      if (content instanceof File) {
        const formData = new FormData();
        formData.append('file', content);
        formData.append('organization_id', organization_id);
        formData.append('folder_id', parentFolderId || '');
  
        formData.append('file_name', content.name);
  
        mutateFile({ file: formData, isFolderUpload: true });
      } else {
        // Handle folder creation
        const folder = {
          name,
          parent_folder_id: parentFolderId,
          organization_id: organization_id,
        };

        // Create the folder in the database and get the folder ID
        const { data: { id: folderId } } = await handleCreateFolder(folder);
  
        // Recursively create the subfolder structure
        await createFolderStructure(content, folderId);
      }
    });
  
    await Promise.all(promises);

    toast({
      description: "Folder uploaded successfully!",
    });
  };

  const handleFolderSelect = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
  
    if (!files || files.length === 0) {
      console.log('No files selected.');
      toast({ description: "You can't upload an empty folder, create one instead." });
      return;
    }
  
    // Organize the folder structure
    if (files && files.length > 0) {
      const folderStructure = {};
  
      Array.from(files).forEach((file) => {
        const pathParts = file.webkitRelativePath.split('/');
        const fileName: any = pathParts.pop();
        let currentLevel: any = folderStructure;
  
        pathParts.forEach(part => {
          if (!currentLevel[part]) {
            currentLevel[part] = {};
          }
          currentLevel = currentLevel[part];
        });
  
        currentLevel[fileName] = file;
      });
    
      try {
        await createFolderStructure(folderStructure, folder_id || null);
        toast({ description: "Folder and files uploaded successfully!" });
      } catch (error) {
        console.error('Error uploading folder structure:', error);
        toast({ description: "An error occurred during upload. Please try again." });
      }
    }
  };

  const handleFileClick = () => {
    if (fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  const handleFolderClick = () => {
    if (folderInputRef.current) {
      folderInputRef.current.click();
    }
  };

  useEffect(() => {
    if (shown) {
      const menuElement = menuRef.current;
      if (menuElement) {
        const menuWidth = menuElement.offsetWidth;
        const menuHeight = menuElement.offsetHeight;

        const adjustedX = (position.x + menuWidth > window.innerWidth) ? (window.innerWidth - menuWidth) : position.x;
        const adjustedY = (position.y + menuHeight > window.innerHeight) ? (window.innerHeight - menuHeight) : position.y;
        setAdjustedPosition({ x: adjustedX, y: adjustedY });
      }
    }
  }, [shown, position]);

  return (
    shown && (
      <div
        className="flex fixed z-[2] top-0 bottom-0 left-0 right-0 w-full h-full"
        onClick={() => close()}
      >
        <div
          ref={menuRef}
          className="w-56 h-[180px]  bg-black/10 bg-neutral-900"
          onClick={(e) => e.stopPropagation()}
          style={{ top: adjustedPosition.y, left: adjustedPosition.x, position: 'absolute' }}
        >
          <div>
            <div
              onClick={() => {
                setShowCreate(true);
                close();
              }}
              className="flex gap-4 border-b border-neutral-400 pb-3 cursor-pointer p-2 hover:bg-neutral-600 transition-all duration-100"
            >
              <FolderPlus />
              <p>Create Folder</p>
            </div>
            <div className="flex flex-col gap-4 py-3">
              <div
                onClick={handleFileClick}
                className="flex gap-4 cursor-pointer p-2 hover:bg-neutral-600 transition-all duration-100"
              >
                <FileUp />
                <p>Upload File</p>
                <input
                  ref={fileInputRef}
                  type="file"
                  multiple
                  style={{ display: 'none' }}
                  onChange={handleFileSelect}
                />
              </div>

              <div
                onClick={handleFolderClick}
                className="flex gap-4 cursor-pointer p-2 hover:bg-neutral-600 transition-all duration-100"
              >
                <FolderUp />
                <p>Upload Folder</p>
                <input
                  ref={folderInputRef}
                  type="file"
                  style={{ display: 'none' }}
                  // @ts-ignore
                  webkitdirectory="true"
                  multiple
                  onChange={handleFolderSelect}
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  );
};

export default UploadFile;
