export const filterCount = (filterValues) => {
  let count = 0;
  for (const [key, value] of Object.entries(filterValues)) {
    if (value) {
      count += 1;
    }
  }
  return count;
};
