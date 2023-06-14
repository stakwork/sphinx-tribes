export const filterCount = (filterValues: any) => {
  let count = 0;
  for (const [_key, value] of Object.entries(filterValues)) {
    if (value) {
      count += 1;
    }
  }
  return count;
};
