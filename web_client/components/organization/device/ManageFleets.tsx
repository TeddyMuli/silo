"use client";

import React, { useState } from 'react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { getDevices, getFleets } from '@/queries';
import { useOrganizationId } from '@/constants';
import { handleDeleteFleet, handleUpdateDevice } from '@/mutations';
import { toast } from '@/hooks/use-toast';
import { Icon } from '@iconify/react/dist/iconify.js';
import Update from '@/components/modals/Update';
import DeleteModal from '@/components/modals/DeleteModal';

const ManageFleets = () => {
  const organization_id = useOrganizationId()
  const [selectedDevice, setSelectedDevice] = useState<{ [key: string]: string | null }>({})
  const [renameFleetId, setRenameFleetId] = useState<string | null>(null);
  const [deleteFleetId, setDeleteFleetId] = useState<string | null>(null);
  const queryClient = useQueryClient()

  const { data: fleets } = useQuery({
    queryKey: ['fleets'],
    queryFn: () => getFleets(organization_id),
    enabled: !!organization_id
  })

  const { data: devices } = useQuery({
    queryKey: ['devices'],
    queryFn: () => getDevices(organization_id, ""),
    enabled: !!organization_id
  })

  const { mutate: updateDevice } = useMutation({
    mutationKey: ['devices'],
    mutationFn: ({ deviceId, device } : { deviceId: string, device: any }) => handleUpdateDevice(deviceId, device),
    onSuccess: (data: any) => {
      if (data.status === 200) {
        toast({
          description: "Device updated successfully!"
        })
        queryClient.invalidateQueries({ queryKey: ['devices'] })
      }
    }
  })

  const { mutate: deleteFleet } = useMutation({
    mutationKey: ['fleets'],
    mutationFn: (fleetId: string) => handleDeleteFleet(fleetId),
    onSuccess: (data: any) => {
      if (data.status === 200) {
        toast({
          description: "Fleet deleted successfully!"
        })
        queryClient.invalidateQueries({ queryKey: ['fleets'] })
      }
    }
  })

  function onUpdateDevice (e: any) {
    e.preventDefault()
    const formData = new FormData(e.target as HTMLFormElement);
    const fleetId = formData.get('fleetId');
    const deviceId = formData.get('deviceSelect') as string;
    const device = {
      fleet_id: fleetId
    }
    updateDevice({ deviceId, device })
  }

  function onDeleteFleet (fleetId: string) {
    deleteFleet(fleetId)
  }

  const handleDeviceChange = (fleetId: string, deviceId: string) => {
    setSelectedDevice((prevState) => ({
      ...prevState,
      [fleetId]: deviceId,
    }));
  };

  return (
    <Table>      
      <TableHeader>
        <TableRow>
          <TableHead>Fleet Name</TableHead>
          <TableHead>Devices</TableHead>
          <TableHead>Devices</TableHead>
          <TableHead>Actions</TableHead>
        </TableRow>
      </TableHeader>

      <TableBody>
        {fleets?.map((fleet: any) => (
          <TableRow key={fleet.id}>
            <TableCell>{fleet.name}</TableCell>
            <TableCell>{devices?.filter((device: any) => device.fleet_id === fleet.id).length || 0}</TableCell>
            <TableCell>
              <form onSubmit={onUpdateDevice} className="flex items-center space-x-2">
                <input type="hidden" name="fleetId" value={fleet.id} />
                <Select name="deviceSelect" onValueChange={(value) => handleDeviceChange(fleet.id, value)}>
                  <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder="Add device" />
                  </SelectTrigger>
                  <SelectContent>
                    {devices?.filter((device: any) => !device.fleet_id)?.map((device: any) => (
                      <SelectItem key={device.id} value={device.id.toString()}>
                        {device.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <Button
                  type="submit"
                  size="sm"
                  disabled={!selectedDevice[fleet.id]}
                >
                  Add
                </Button>
              </form>
            </TableCell>
            <TableCell>
              <Update
                key={fleet.id}
                shown={renameFleetId === fleet.id}
                action='Rename'
                close={() => setRenameFleetId(null)}
                type='fleet'
                object={fleet}
              />

              <DeleteModal
                shown={deleteFleetId === fleet.id}
                close={() => setDeleteFleetId(null)}
                object={`${fleet.name}`}
                onConfirm={() => onDeleteFleet(fleet.id)}
              />

              <div className='flex gap-4 items-center'>
                <div onClick={() => setRenameFleetId(fleet.id)}>
                  <Icon
                    icon="lucide:pen"
                    height={20}
                    width={20}
                    className='cursor-pointer'
                  />
                </div>
                <div onClick={() => setDeleteFleetId(fleet.id)}>
                  <Icon
                    icon="material-symbols:delete"
                    height={24}
                    width={24}
                    style={{ color: "#ff0000" }}
                    className='cursor-pointer text-#ff0000'
                  />
                </div>
              </div>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

export default ManageFleets;
