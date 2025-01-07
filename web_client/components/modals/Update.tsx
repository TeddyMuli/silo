import React, { useEffect, useState } from 'react';
import Modal from './Modal';
import { useToast } from '@/hooks/use-toast';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { handleRenameFile, handleRenameFolder } from '@/mutations';
import { useFolderId, useOrganizationId } from '@/constants';

const Update = ({ type, action, object, shown, close } : { type: string, action: string, object: any, shown: boolean, close: () => void }) => {
  const { toast } = useToast();
  const organization_id = useOrganizationId()
  const folder_id = useFolderId()

  const queryClient = useQueryClient()

  const [name, setName] = useState("");
  const [serial_number, setSerialNumber] = useState("");

  useEffect(() => {
    if (object) {
      setName(object?.name)
      if (type === "device") {
        setSerialNumber(object?.serial_number)  
      }
    }
  }, [shown]);

  const folderMutationKey = folder_id ? ['folders', folder_id] : ['folders', organization_id]
  const fileMutationKey = folder_id ? ['files', folder_id] : ['files', organization_id]

  const { mutate: updateFolder } = useMutation({
    mutationKey: folderMutationKey,
    mutationFn: ({ folder, folderId }: { folder: any, folderId: string }) => handleRenameFolder(folder, folderId),
    onSuccess: (data: any) => {
      if (data.status === 200) {
        close()
        toast({
          description: "Folder Updated!"
        })
        queryClient.invalidateQueries({ queryKey: folderMutationKey })
      }
    }
  })

  const { mutate: updateFile } = useMutation({
    mutationKey: fileMutationKey,
    mutationFn: ({ file, fileId }: { file: any, fileId: string }) => handleRenameFile(file, fileId),
    onSuccess: (data: any) => {
      if (data.status === 200) {
        close()
        toast({
          description: "File Updated!"
        })
        queryClient.invalidateQueries({ queryKey: fileMutationKey })
      }
    }
  })

  const handleRename = () => {
    if (type === "folder") {
      const folder = {
        name: name
      }

      updateFolder({ folder, folderId: object?.id })
    } else if (type === "file") {
      const file = {
        name: name
      }

      updateFile({ file, fileId: object?.id})
    } else {
      return
    }
  }

  return (
    <Modal
      shown={shown}
      close={close}
    >
      <div className='flex flex-col p-4 gap-4 justify-center items-center w-[400px] bg-neutral-800 rounded-xl'>
        <p className='text-xl font-semibold'>{action} {type}</p>
        <div className='flex flex-col'>
          <label className='p-2 font-medium'>Name</label>
          <input
            type="text"
            placeholder={`Enter ${type} name`}
            className='p-3 outline-none border-2 border-indigo-600 rounded-xl w-64'
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>

        <div className='flex gap-4 ml-auto mr-16 text-lg font-medium'>
          <button
            className='text-neutral-500'
            onClick={() => {
              close();
              setName("")
              setSerialNumber("")
            }}
          >
            Cancel
          </button>
          <button
            disabled={
              !name.trim() ||
              (name.trim() === object?.name &&
                (type !== "device" || serial_number.trim() === object?.serial_number)
              )
            }
            className='disabled:text-blue-300 cursor-pointer disabled:cursor-not-allowed text-blue-500'
            onClick={handleRename}
          >
            {action}
          </button>
        </div>
      </div>
    </Modal>
  );
}

export default Update;
