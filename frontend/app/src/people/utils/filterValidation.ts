export const bountyHeaderFilter = ({ Paid, Assigned, Opened }, bodyPaid, bodyAssignee) => {
  if (Paid) {
    if (Assigned) {
      if (Opened) {
        return true;
      } else {
        return bodyAssignee || bodyPaid;
      }
    } else {
      if (Opened) {
        return bodyPaid || !bodyAssignee;
      } else {
        return bodyPaid;
      }
    }
  } else {
    if (Assigned) {
      if (Opened) {
        return !bodyPaid;
      } else {
        return !bodyPaid && bodyAssignee;
      }
    } else {
      if (Opened) {
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
