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
  if (Object.keys(filterLanguage)?.every((key: any) => !filterLanguage[key])) {
    return true;
  } else return codingLanguage?.some(({ value }: any) => filterLanguage[value]) ?? false;
};
