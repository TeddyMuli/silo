"use client";

import { HardDrive, Heart, LayoutDashboard, Plus, Trash2 } from 'lucide-react';
import Link from 'next/link';
import React, { useState } from 'react';
import Logo from '../shared/Logo';
import { Icon } from '@iconify/react';
import FolderTreeComponent from './FolderTree';
import UploadFile from '../modals/UploadFile';
import Create from '../modals/Create';
import { useOrganizationId } from '@/constants';

const LeftSideBar = () => {
  const organization_id = useOrganizationId();
  const orgUrl = `/organization/${organization_id}`;

  const Links = [
    { name: "Dashboard", path: `${orgUrl}`, icon: LayoutDashboard },
    { name: "Favorites", path: `${orgUrl}/favorites`, icon: Heart },
    { name: "Drive", path: `${orgUrl}/drive`, icon: HardDrive },
    { name: "Recycle Bin", path: `${orgUrl}/bin`, icon: Trash2 },
  ];

  const [expandDrive, setExpandDrive] = useState<boolean>(false);
  const [showCreateOptions, setShowCreateOptions] = useState<boolean>(false);
  const[showCreate, setShowCreate] = useState<boolean>(false);

  const toggleExpandDrive = () => {
    setExpandDrive(prevState => !prevState);
  };

  const uploadFilePosition = {
    y: 100,
    x: 20,
  };

  return (
    <section className='overflow-y-auto h-screen w-52 items-center justify-center'>
      <Logo />

      <div>
        <div
          onClick={() => setShowCreateOptions(true)}
          className='flex gap-2 bg-indigo-600 p-2 rounded-xl text-white mb-6 cursor-pointer items-center'
        >
          <Plus />
          <p className='font-medium text-xl'>Create</p>
        </div>
          <UploadFile
            shown={showCreateOptions}
            close={() => setShowCreateOptions(false)}
            setShowCreate={setShowCreate}
            position={uploadFilePosition}
          />
          <Create shown={showCreate} object='Folder' close={() => setShowCreate(false)} />
        {Links?.map((link, index) => (
          <div key={index} className='flex flex-col my-4 items-start'>
          <div className='flex items-center'>
            <Icon
              onClick={toggleExpandDrive}
              icon="material-symbols:arrow-right"
              className={`${expandDrive ? "rotate-90" : ""} ${link.name == "Drive" ? "opacity-100" : "opacity-0"}`}
              width={24}
              height={24}
            />
            <Link href={link.path} className='flex gap-3 items-center'>
              <link.icon />
              <p>{link.name}</p>
            </Link>
          </div>
          {link.name === "Drive" && expandDrive && (
            <div className='ml-4'>
              <FolderTreeComponent />
            </div>
          )}
        </div>
        ))}
      </div>
    </section>
  );
}

export default LeftSideBar;
