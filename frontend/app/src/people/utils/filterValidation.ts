export const bountyHeaderFilter = ({ Paid, Assigned, Open }, bodyPaid, bodyAssignee) => {
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

export const bountyHeaderLanguageFilter = (codingLanguage, filterLanguage) => {
  if (Object.keys(filterLanguage)?.every((key) => !filterLanguage[key])) {
    return true;
  } else return codingLanguage?.some(({ value }) => filterLanguage[value]) ?? false;
};
