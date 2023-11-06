export const filterCount = (filterValues: any) => {
  let count = 0;
  for (const [, value] of Object.entries(filterValues)) {
    if (value) {
      count += 1;
    }
  }
  return count;
};

export const formatSat = (sat: number) => {
  if (sat === 0 || !sat) {
    return '0';
  }
  const satsWithComma = sat.toLocaleString();
  const splittedSat = satsWithComma.split(',');
  return splittedSat.join(' ');
};
