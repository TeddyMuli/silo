"use client";

import { fetchOrganization } from '@/queries';
import { useQuery } from '@tanstack/react-query';
import { Settings } from 'lucide-react';
import { useParams } from 'next/navigation';
import React, { useEffect } from 'react';

const Header = () => {
  const params = useParams();
  const { organization_id } = params

  const orgId = Array.isArray(organization_id) ? organization_id[0] : organization_id;

  const { data: organization } = useQuery({
    queryKey: ['organization', organization_id],
    queryFn: () => fetchOrganization(orgId),
    enabled: !!organization_id
  })

  return (
    <div className='flex my-4 text-2xl font-medium justify-between'>
      <p>{organization?.name}</p>
      <Settings />
    </div>
  );
}

export default Header;
