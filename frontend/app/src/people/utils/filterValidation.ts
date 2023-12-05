export const bountyHeaderFilter = (
  { Paid, Assigned, Open }: any,
  bodyPaid: any,
  bodyAssignee: any
) => {
  if (Paid) {
    if (Assigned) {
      if (Open) {
        return true;
      } else {
        return bodyAssignee || bodyPaid;
      }
    } else {
      if (Open) {
        return bodyPaid || !bodyAssignee;
      } else {
        return bodyPaid;
      }
    }
  } else {
    if (Assigned) {
      if (Open) {
        return !bodyPaid;
      } else {
        return !bodyPaid && bodyAssignee;
      }
    } else {
      if (Open) {
        return !bodyPaid && !bodyAssignee;
      } else {
        return true;
      }
    }
  }
};

export const bountyHeaderLanguageFilter = (codingLanguage: any, filterLanguage: any) => {
  const selectedLanguages: string[] = Object.keys(filterLanguage).filter(
    (key: string) => filterLanguage[key]
  );
  if (!Array.isArray(selectedLanguages) || selectedLanguages.length === 0) {
    return true; // No filter selected, show all bounties
  } else {
    // Use "and" logic - all selected skills must match- thus using every
    return selectedLanguages.every((selectedLanguage: string) =>
      codingLanguage.includes(selectedLanguage)
    );
  }
};
