import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { getDeviceData } from '@/queries';

const DeviceBattery = ({ serialNumber } : { serialNumber: string }) => {
  const { data: deviceData, isLoading, isError } = useQuery({
    queryKey: ['device', serialNumber],
    queryFn: () => getDeviceData(serialNumber),
    enabled: !!serialNumber
  })

  return (
    <div>
      <p>
        {deviceData?.battery !== undefined ? `${deviceData?.battery} %` : "offline"}
      </p>
    </div>
  );
}

export default DeviceBattery;
