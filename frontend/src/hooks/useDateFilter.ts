import { useState } from 'react';
import { DateRange } from 'react-day-picker';

export const useDateFilter = (initialRange?: DateRange) => {
  const [dateRange, setDateRange] = useState<DateRange | undefined>(
    initialRange || { from: new Date(), to: new Date() }
  );

  return { dateRange, setDateRange };
};
