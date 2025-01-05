"use client";

import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@radix-ui/react-label';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { getFleets } from '@/queries';
import { useOrganizationId } from '@/constants';
import { useForm } from 'react-hook-form';
import { handleCreateDevice } from '@/mutations';
import { toast } from '@/hooks/use-toast';
import Loading from '@/components/shared/Loading';

const inputs = [
  { id: "name", name: "Device Name", placeHolder: "Enter device name" },
  { id: "serial_number", name: "Serial Number", placeHolder: "Enter serial number" },
];

const RegisterDevice = () => {
  const [selectedFleet, setSelectedFleet] = useState<string>("");

  const organization_id = useOrganizationId()
  const queryClient = useQueryClient()

  const { data: fleets } = useQuery({
    queryKey: ['fleets'],
    queryFn: () => getFleets(organization_id),
    enabled: !!organization_id
  })

  const {
    handleSubmit,
    register,
    getValues,
    reset,
    formState: { isDirty, isSubmitting }
  } = useForm({
    mode: "onChange",
    defaultValues: {
      name: "",
      serial_number: "",
      fleet_id: ""
    }
  })

  const { mutate } = useMutation({
    mutationKey: ['devices'],
    mutationFn: (device: any) => handleCreateDevice(device),
    onSuccess: (data: any) => {
      if (data.status === 200 || data.status === 201) {
        toast({
          description: "Device created successsfully!"
        })

        reset()
        queryClient.invalidateQueries({ queryKey: ['devices'] })
      } else if (data.data.error === "Device with the same name exists!") {
        toast({
          description: "Device with the same name exists!",
          variant: "destructive"
        })
      }else if (data.data.error === "Device with the same serial number exists!") {
        toast({
          description: "Device with the same serial number exists!",
          variant: "destructive"
        })
      }
    }
  })

  function onSubmit () {
    const { name, serial_number } = getValues()

    const device = {
      organization_id: organization_id,
      name: name,
      serial_number: serial_number,
      fleet_id: selectedFleet
    }

    mutate(device)
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      {inputs?.map((input, index) => (
        <div key={index} className="space-y-2">
          <Label htmlFor={input.id}>{input.name}</Label>
          <Input
            id={input.id}
            // @ts-ignore
            {...register(input.id)}
            placeholder={input.placeHolder}
            required
          />
        </div>
      ))}
      <div className="space-y-2">
        <Label htmlFor="fleetId">Fleet</Label>
        <Select onValueChange={setSelectedFleet}>
          <SelectTrigger className="w-full">
            <SelectValue placeholder="Select a fleet" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              {fleets?.map((fleet: any, index: number) => (
                <SelectItem key={index} value={fleet.id}>
                  {fleet.name}
                </SelectItem>
              ))}
            </SelectGroup>
          </SelectContent>
        </Select>
      </div>
      <Button
        type="submit"
        disabled={!isDirty}
        className='flex'
      >
        {isSubmitting && <Loading className='text-white' />}
        Register Device
      </Button>
    </form>
  );
}

export default RegisterDevice;
