import React from 'react';
import { useToast } from '@/hooks/use-toast';
import { TriangleAlert } from 'lucide-react';

const DeleteModal = (
  { shown, close, object, onConfirm } : { shown: boolean, close: () => void, object: string, onConfirm: () => void }
) => {
  const { toast } = useToast();
  const handleConfirm = async (e: any) => {
    e.preventDefault();
    onConfirm();
    toast({
      title: `${object} deleted successfully!`
    })
    close();
  }

  return shown && (
    <div className='flex fixed top-0 bottom-0 left-0 right-0 pt-32 z-[2] bg-white bg-opacity-10 w-full justify-center'>
      <div className='bg-white w-[400px] h-[300px] rounded-xl text-black p-4 flex flex-col gap-3 justify-center items-center'>
        <TriangleAlert className='text-red-600' size={40} />
        <p className='text-xl font-semibold'>Delete {object}</p>
        <p className='text-lg text-center'>
          You are going to delete {object}.<br />
          <span className='text-xl font-medium'>Are you sure?</span>
        </p>
        <div className='flex gap-6'>
          <button
            onClick={close}
            className='px-4 py-2 rounded-full text-black font-medium bg-neutral-200'
          >
            No, Keep It.
          </button>
          <button
            onClick={handleConfirm}
            className='px-4 py-2 rounded-full text-white font-medium bg-red-600'
          >
            Yes, Delete!
          </button>
        </div>
      </div>
    </div>
  );
}

export default DeleteModal;
