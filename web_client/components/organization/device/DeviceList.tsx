"use client";

import React, { useState } from 'react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { useOrganizationId } from '@/constants';
import { getDevices, getFleets } from '@/queries';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Icon } from '@iconify/react/dist/iconify.js';
import DeleteModal from '@/components/modals/DeleteModal';
import Update from '@/components/modals/Update';
import { toast } from '@/hooks/use-toast';
import { handleDeleteDevice } from '@/mutations';
import DeviceStatus from './DeviceStatus';
import DeviceBattery from './DeviceBattery';

const DeviceList = () => {
  const organization_id = useOrganizationId()
  const [updateDeviceId, setUpdateDeviceId] = useState<string | null>(null);
  const [deleteDeviceId, setDeleteDeviceId] = useState<string | null>(null);
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

  const getFleetName = (fleet_id: string) => {
    const fleet = fleets?.find((f: any) => f.id === fleet_id);
    return fleet ? fleet.name : 'Unassigned';
  };

  const { mutate: deleteDevice } = useMutation({
    mutationKey: ['devices'],
    mutationFn: (deviceId: string) => handleDeleteDevice(deviceId),
    onSuccess: (data: any) => {
      if (data.status === 200) {
        toast({
          description: "Device deleted successfully!"
        })
        queryClient.invalidateQueries({ queryKey: ['devices'] })
      }
    }
  })

  function onDeleteDevice (fleetId: string) {
    deleteDevice(fleetId)
  }


  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Serial Number</TableHead>
          <TableHead>Fleet</TableHead>
          <TableHead>Actions</TableHead>
          <TableHead>Battery</TableHead>
          <TableHead>Status</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {Array.isArray(devices) && devices?.map((device: any) => {
          return (
            <TableRow key={device.id}>
              <TableCell>{device.name}</TableCell>
              <TableCell>{device.serial_number}</TableCell>
              <TableCell>{getFleetName(device.fleet_id)}</TableCell>
              <TableCell>
                <Update
                  key={device.id}
                  shown={updateDeviceId === device.id}
                  action='Update'
                  close={() => setUpdateDeviceId(null)}
                  type='device'
                  object={device}
                />

                <DeleteModal
                  shown={deleteDeviceId === device.id}
                  close={() => setDeleteDeviceId(null)}
                  object={`${device.name}`}
                  onConfirm={() => onDeleteDevice(device.id)}
                />

                <div className='flex gap-4 items-center'>
                  <div onClick={() => setUpdateDeviceId(device.id)}>
                    <Icon
                      icon="lucide:pen"
                      height={20}
                      width={20}
                      className='cursor-pointer'
                    />
                  </div>
                  <div onClick={() => setDeleteDeviceId(device.id)}>
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
              <TableCell>
                <DeviceBattery serialNumber={device.serial_number} />
              </TableCell>
              <TableCell>
                <DeviceStatus serialNumber={device.serial_number} />
              </TableCell>
            </TableRow>
          )})}
      </TableBody>
    </Table>
  );
}

export default DeviceList;
