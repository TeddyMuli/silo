import React, { useState } from 'react';
import UploadFile from '../modals/UploadFile';
import Create from '../modals/Create';

const ContextMenu = () => {
  const [menuVisible, setMenuVisible] = useState(false);
  const [menuPosition, setMenuPosition] = useState({ x: 0, y: 0 });
  const [showCreate, setShowCreate] = useState(false);

  const handleContextMenu = (event: any) => {
    event.preventDefault();
    setMenuPosition({ x: event.pageX, y: event.pageY });
    setMenuVisible(true);
  };

  const handleClick = () => {
    setMenuVisible(false);
  };

  return (
    <div
      onContextMenu={handleContextMenu}
      onClick={handleClick}
      style={{ height: '100%' }}
    >
      {menuVisible && (
        <>
          <UploadFile
            shown={menuVisible}
            close={() => setMenuVisible(false)}
            setShowCreate={setShowCreate}
            position={menuPosition}
          />
          <Create shown={showCreate} object='Folder' close={() => setShowCreate(false)} />
        </>
      )}
    </div>
  );
}

export default ContextMenu;
