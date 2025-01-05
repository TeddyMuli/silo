import React, { useState } from 'react';
import Modal from './Modal';
import { useToast } from '@/hooks/use-toast';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useUser } from '@/context/UserContext';
import { handleCreateFolder, handleCreateOrganization } from '@/mutations';
import { useOrganizationId } from '@/constants';

const Create = ({ object, shown, close, parent_folder_id } : { object: string, shown: boolean, close: () => void, parent_folder_id?: string | null }) => {
  const { toast } = useToast();
  const[name, setName] = useState('');
  const queryClient = useQueryClient();
  const { user } = useUser();
  const organization_id = useOrganizationId();
  
  const { mutate: createOrganization } = useMutation({
    mutationKey: ['organizations', user?.user_id],
    mutationFn: (organization: any) => handleCreateOrganization(organization),
    onSuccess: (data) => {
      if (data.status === 200) {
        queryClient.invalidateQueries({ queryKey: ['organizations'] });
        toast({
          description: `${object} created successfully!`
        })
        setName("")
        close();
      } else {
        if (data?.data.error == "Organization with the same name exists") {
          toast({
            description: "Organization with the same name exists",
            variant: "destructive"
          })  
        }
      }
    },
    onError: () => {
      toast({
        description: `Error creating ${object}`,
        variant: "destructive"
      });
    }
  });

  const { mutate: createFolder, variables } = useMutation({
    mutationKey: ['folders', organization_id],
    mutationFn: (folder: any) => handleCreateFolder(folder),
    onSuccess: (data) => {
      if (data?.status === 201) {
        const queryKey = variables.parent_folder_id ? ['folders', variables.parent_folder_id] : ['folders', organization_id];
        queryClient.invalidateQueries({ queryKey });
        toast({
          description: `${object} created successfully!`
        })  
        setName("")
        close();    
      } else {
        if (data?.data.error == "Folder with the same name exists") {
          toast({
            description: "Folder with the same name exists",
            variant: "destructive"
          })  
        }
      }
    },
    onError: () => {
      toast({
        description: `Error creating ${object}`,
        variant: "destructive"
      });
    }
  });

  const handleCreate = () => {
    if (object === 'Organization') {
      const organization = {
        user_id: user?.user_id,
        name: name
      };  
      createOrganization(organization);
    } else if (object === 'Folder') {
      const folder = {
        organization_id: organization_id,
        name: name,
        parent_folder_id: parent_folder_id || null
      };  
      createFolder(folder);
    }
  }

  return (
    <Modal
      shown={shown}
      close={close}
    >
      <div className='flex flex-col p-4 gap-4 justify-center items-center h-[200px] w-[400px] bg-neutral-800 rounded-xl'>
        <p className='text-xl font-semibold'>New {object}</p>
        <input
          type="text"
          placeholder={`Enter ${object} name`}
          className='p-3 outline-none border-2 border-indigo-600 rounded-xl w-64'
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <div className='flex gap-4 ml-auto mr-16 text-lg font-medium'>
          <button className='text-neutral-500' onClick={close}>Cancel</button>
          <button disabled={!name} className='text-blue-500' onClick={handleCreate}>Create</button>
        </div>
      </div>
    </Modal>
  );
}

export default Create;
