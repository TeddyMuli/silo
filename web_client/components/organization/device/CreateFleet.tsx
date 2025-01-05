"use client";

import Loading from '@/components/shared/Loading';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useOrganizationId } from '@/constants';
import { toast } from '@/hooks/use-toast';
import { handleCreateFleet } from '@/mutations';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import React from 'react';
import { useForm } from 'react-hook-form';

const CreateFleet = () => {
  const organization_id = useOrganizationId()
  const queryClient = useQueryClient()

  const {
    register,
    reset,
    getValues,
    handleSubmit,
    formState: { isDirty, isSubmitting }
  } = useForm({
    mode: "onChange",
    defaultValues: {
      name: ""
    }
  })

  const { mutate } = useMutation({
    mutationKey: ['fleets'],
    mutationFn: (fleet: any) => handleCreateFleet(fleet),
    onSuccess: (data: any) => {
      if (data.status === 200 || data.status === 201) {
        toast({
          description: "Fleet created successfully!"
        })

        reset()

        queryClient.invalidateQueries({ queryKey: ["fleets"] })
      } else if (data.data.error === "Fleet with the same name exists!") {
        toast({
          description: "Fleet with the same name exists!",
          variant: "destructive"
        })
      }
    }
  })

  function onSubmit () {
    const fleet = {
      organization_id: organization_id,
      name: getValues().name
    }

    mutate(fleet)
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="fleetName">Fleet Name</Label>
        <Input
          id="fleetName"
          {...register("name")}
          placeholder='Enter fleet name'
          required
        />
      </div>
      <Button
        type="submit"
        disabled={!isDirty}
        className='flex disabled:cursor-not-allowed'
      >
        {isSubmitting && <Loading className='text-white' />}
        Create Fleet
      </Button>
    </form>
  );
}

export default CreateFleet;
