import React from 'react';

const Offline = ({
  children, shown, className }: 
  { children:  React.ReactNode, shown: boolean, className?: string }) => {
  return shown && (
    <div
      className='flex fixed justify-center items-center z-[2] top-0 bottom-0 left-0 right-0 bg-white bg-opacity-5'
    >
      {children}
    </div>
  );
}

export default Offline;
