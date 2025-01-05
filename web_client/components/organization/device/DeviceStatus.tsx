import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { getDeviceData } from '@/queries';

const DeviceStatus = ({ serialNumber } : { serialNumber: string }) => {
  const { data: deviceData, isLoading, isError } = useQuery({
    queryKey: ['device', serialNumber],
    queryFn: () => getDeviceData(serialNumber),
    enabled: !!serialNumber
  })

  return (
    <div className="relative group group-hover:cursor-pointer">
      <div
        className={`h-3 w-3 rounded-full ${deviceData?.online ? 'bg-green-600' : 'bg-red-600'}`}
      ></div>
      <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 opacity-0 group-hover:opacity-100 transition-opacity bg-gray-700 text-white text-xs px-2 py-1 rounded">
        {deviceData?.online ? 'Online' : 'Offline'}
      </div>
    </div>
  );
}

export default DeviceStatus;
