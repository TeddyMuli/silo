import { useToast } from '@/hooks/use-toast';
import React, { useEffect } from 'react';

const Error = ({ action, description } : { action: string, description: string }) => {
  const { toast } = useToast();

  useEffect(() => {
    toast({
      title: `Error ${action}`,
      description: `${description}`,
      variant: "destructive"
    })
  }, []);

  return <div></div>
}

export default Error;
