import React, { useState, useEffect } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { handleDeleteFile, handleDeleteFolder, handleMoveFileToTrash, handleMoveFolderToTrash, handleRestoreFile, handleRestoreFolder } from '@/mutations';
import { toast } from '@/hooks/use-toast';
import { downloadFile, downloadFolderAsZip, useFolderId, useOrganizationId } from '@/constants';
import { Icon } from '@iconify/react/dist/iconify.js';
import { usePathname } from 'next/navigation';
import DeleteModal from './DeleteModal';
import Update from './Update';

const options = [
  { icon: "tabler:download", name: "Download" },
  { icon: "lucide:pen", name: "Rename" },
  { icon: "material-symbols:delete", name: "Move to Trash" },
  { icon: "ic:twotone-restore", name: "Restore" },
  { icon: "material-symbols:delete", name: "Delete" },
];

const Options = (
  { shown, close, type, object, triggerRef }
  :
  { shown: boolean, close: () => void, type: string, object: any, triggerRef: any }
) => {
  const organization_id = useOrganizationId()
  const folder_id = useFolderId();
  const pathName = usePathname()
  const segments = pathName.split('/');
  const lastSegment = segments.pop() || '';

  const queryClient = useQueryClient()

  const [showRename, setShowRename] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [position, setPosition] = useState({ top: 0, left: 0 });
  const folderMutationKey = folder_id ? ['folders', folder_id] : ['folders', organization_id]
  const fileMutationKey = folder_id ? ['files', folder_id] : ['files', organization_id]

  const { mutate: moveFolderTrash } = useMutation({
    mutationKey: folderMutationKey,
    mutationFn: (folderId: string) => handleMoveFolderToTrash(folderId),
    onSuccess: (data: any) => {
      if (data.status === 200) {
        close()
        toast({
          description: "Folder moved to trash!"
        })
        queryClient.invalidateQueries({ queryKey: folderMutationKey })
      }
    }
  });

  const { mutate: moveFileTrash } = useMutation({
    mutationKey: fileMutationKey,
    mutationFn: (fileId: string) => handleMoveFileToTrash(fileId),
    onSuccess: (data: any) => {
      if (data.status === 200) {
        close()
        toast({
          description: "File moved to trash!"
        })
        queryClient.invalidateQueries({ queryKey: fileMutationKey })
      }
    }
  });

  const { mutate: deleteFolder } = useMutation({
    mutationKey: folderMutationKey,
    mutationFn: (folderId: string) => handleDeleteFolder(folderId),
    onSuccess: (data: any) => {
      if (data.status === 200) {
        close()
        toast({
          description: "Folder deleted!"
        })
        queryClient.invalidateQueries({ queryKey: folderMutationKey })
      }
    }
  });

  const { mutate: deleteFile } = useMutation({
    mutationKey: fileMutationKey,
    mutationFn: (fileId: string) => handleDeleteFile(fileId),
    onSuccess: (data: any) => {
      if (data.status === 200) {
        close()
        toast({
          description: "File deleted!"
        })
        queryClient.invalidateQueries({ queryKey: fileMutationKey })
      }
    }
  });

  const { mutate: restoreFolder } = useMutation({
    mutationKey: folderMutationKey,
    mutationFn: (folderId: string) => handleRestoreFolder(folderId),
    onSuccess: (data: any) => {
      if (data.status === 200) {
        close()
        toast({
          description: "Folder restored!"
        })
        queryClient.invalidateQueries({ queryKey: folderMutationKey })
      }
    }
  });

  const { mutate: restoreFile } = useMutation({
    mutationKey: fileMutationKey,
    mutationFn: (fileId: string) => handleRestoreFile(fileId),
    onSuccess: (data: any) => {
      if (data.status === 200) {
        close()
        toast({
          description: "File restored!"
        })
        queryClient.invalidateQueries({ queryKey: fileMutationKey })
      }
    }
  });

  const handleRestore = () => {
    if (type === "folder") {
      restoreFolder(object?.id);
    } else if (type === "file") {
      restoreFile(object?.id);
    }
  }

  const handleMoveToTrash = () => {
    if (type === "folder") {
      moveFolderTrash(object?.id);
    } else if (type === "file") {
      moveFileTrash(object?.id);
    }
  };

  const handleDelete = () => {
    if (type === "folder") {
      deleteFolder(object?.id);
    } else if (type === "file") {
      deleteFile(object?.id);
    }
  }

  const handleMoveOrDelete = () => {
    if (lastSegment === "bin") {
      setShowDelete(true)
    } else {
      handleMoveToTrash()
    }
  }

  useEffect(() => {
    if (shown && triggerRef.current) {
      const rect = triggerRef.current.getBoundingClientRect();

      const containerHeight = 100;
      const viewportHeight = window.innerHeight;

      let topPosition = rect.bottom + window.scrollY + 20;
      if (topPosition + containerHeight > viewportHeight) {
        topPosition = rect.top + window.scrollY - containerHeight - 20;
      }
      setPosition({
        top: topPosition,
        left: rect.right + window.scrollX - 150,
      });
    }
  }, [shown, triggerRef]);

  const filteredOptions = options.filter((option) => {
    if (lastSegment === 'bin') {
      return option.name === "Delete" || option.name === "Restore";
    }
    return option.name !== "Delete" && option.name !== "Restore";
  });

  const handleOptionClick = (option: any) => {
    if (option.name === "Rename") {
      setShowRename(true)
    } else if (option.name === "Download") {
      if (type === "folder") {
        downloadFolderAsZip(organization_id, object?.id)
      } else if (type === "file") {
        downloadFile(object?.file_path, object.name)
      }
    } else if (option.name === "Restore") {
      handleRestore()
    } else {
      handleMoveOrDelete()
    }
  }

  return shown && (
    <div
      className='fixed top-0 bottom-0 left-0 right-0 z-[2] bg-transparent w-full h-full translate-all duration-200'
      onClick={() => close()}
    >
      <div
        className="absolute bg-neutral-300 dark:bg-neutral-900 p-3 rounded-lg flex flex-col gap-2 mr-4"
        style={{ top: position.top, left: position.left }}
        onClick={(e) => e.stopPropagation()}
      >
        <DeleteModal
          shown={showDelete}
          close={() => setShowDelete(false)}
          object={`${object.name}`}
          onConfirm={handleDelete}
        />
        <Update
          action='Rename'
          type={type}
          object={object}
          shown={showRename}
          close={() => setShowRename(false)}
        />
        {filteredOptions?.map((option, index) => (
          <div
            key={index}
            className='flex gap-2 items-center cursor-pointer'
            onClick={() => handleOptionClick(option)}
          >
            <Icon
              icon={option.icon}
              height={20}
              width={20}
              style={option.icon === "material-symbols:delete" ? { color: "#ff0000" } : {}}
            />
            <p>{option.name}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Options;
