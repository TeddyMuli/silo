import React from 'react';

const Loading = ({ className } : { className?: string }) => {
  return (
    <div>
      <div
        className={
          `inline-block h-6 w-6 animate-spin rounded-full mx-3
          border-4 border-solid border-current border-r-transparent
          align-[-0.125em] motion-reduce:animate-[spin_1.5s_linear_infinite] ${className}`
        }
        role="status"></div>
    </div>
  );
}

export default Loading;
