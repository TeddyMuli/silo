import React from 'react';

const Modal = ({
  children, shown, close, className }: 
  { children:  React.ReactNode, shown: boolean, close: () => void, className?: string }) => {
  return shown && (
    <div
      className='flex fixed justify-center items-center z-[2] top-0 bottom-0 left-0 right-0 bg-white bg-opacity-5'
      onClick={() => close()}
    >
      <div
        className={className}
        onClick={(e) => {
          e.stopPropagation();
        }}
      >
        {children}
      </div>
    </div>
  );
}

export default Modal;
