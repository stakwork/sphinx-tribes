import { BountyStatus, defaultBountyStatus } from 'store/main';

type SortCriterion = 'open' | 'in-progress' | 'completed' | 'assignee';

export const getBountyStatus = (sortOrder: SortCriterion) => {
  let newStatus: BountyStatus;

  switch (sortOrder) {
    case 'open': {
      newStatus = { ...defaultBountyStatus, Open: true };
      break;
    }
    case 'in-progress': {
      newStatus = {
        ...defaultBountyStatus,
        Open: false,
        Assigned: true
      };
      break;
    }
    case 'completed': {
      newStatus = {
        ...defaultBountyStatus,
        Open: false,
        Paid: true
      };
      break;
    }
    default: {
      newStatus = {
        ...defaultBountyStatus,
        Open: false
      };
      break;
    }
  }

  return newStatus;
};

export const dateFilterOptions = Object.freeze([
  { id: 'Newest', label: 'Newest', value: 'desc' },
  { id: 'Oldest', label: 'Oldest', value: 'asc' }
]);
