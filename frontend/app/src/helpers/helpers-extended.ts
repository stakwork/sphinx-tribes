/* eslint-disable @typescript-eslint/typedef */
import LighningDecoder from 'light-bolt11-decoder';
import { MainStore } from 'store/main';
import { getHost } from '../config/host';

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

export const formatPrice = (amount = 0) => amount;

export const formatSatPrice = (amount = 0) => {
  const dollarUSLocale = Intl.NumberFormat('en-US');
  return dollarUSLocale.format(amount);
};

export const convertToLocaleString = (value: number): string => {
  if (value) {
    const formattedValue = Number(value).toLocaleString();
    return formattedValue;
  } else {
    return '0';
  }
};

export const convertLocaleToNumber = (localeString: string): number => {
  const numString = localeString.replace(/\D/g, '');

  const num = parseInt(numString);
  return num;
};

export const DollarConverter = (e: any) => {
  const dollarUSLocale = Intl.NumberFormat('en-US');
  return dollarUSLocale.format(formatPrice(e)).split(',').join(' ');
};

const host = getHost();
export function makeConnectQR(pubkey: string) {
  return `sphinx.chat://?action=person&host=${host}&pubkey=${pubkey}`;
}

export const extractRepoAndIssueFromIssueUrl = (url: string) => {
  let orgName = '';
  let repoName = '';
  let repo = '';
  let issue = '';

  const splitString = url.split('/');
  const issueIndex = splitString.length - 1;
  const repoNameIndex = splitString.length - 3;
  const orgNameIndex = splitString.length - 4;

  // pop last element if not a number (page focus could be "commits", "checks", "files")
  if (isNaN(parseInt(splitString[issueIndex]))) {
    splitString.pop();
  }

  issue = splitString[issueIndex];
  orgName = splitString[orgNameIndex];
  repoName = splitString[repoNameIndex];
  if (!issue || !orgName || !repoName) {
    throw new Error('Invalid github url');
  }
  repo = `${orgName}/${repoName}`;

  return { repo, issue };
};

export const extractGithubIssue = (
  person: { github_issues: Record<string, any> },
  repo: string,
  issue: string
) => {
  const { github_issues } = person;
  const keyname = `${repo}/${issue}`;
  return (github_issues && github_issues[keyname]) || {};
};

export const extractGithubIssueFromUrl = (
  person: { github_issues: Record<string, any> },
  url: string
) => {
  try {
    const { repo, issue } = extractRepoAndIssueFromIssueUrl(url);
    return extractGithubIssue(person, repo, issue);
  } catch (e) {
    return {};
  }
};

export const randomString = (l: number): string =>
  Array.from(crypto.getRandomValues(new Uint8Array(l)), (byte: any) =>
    `0${(byte & 0xff).toString(16)}`.slice(-2)
  ).join('');

export const sendToRedirect = (url: string) => {
  const el = document.createElement('a');
  el.href = url;
  el.target = '_blank';
  el.click();
};

export const calculateTimeLeft = (
  timeLimit: Date,
  type: 'minutes' | 'days' | 'hours'
): {
  days?: number;
  hours?: number;
  minutes: number;
  seconds: number;
} => {
  const difference = new Date(timeLimit).getTime() - new Date().getTime();

  let timeLeft: any = {};

  if (difference > 0 && type === 'days') {
    timeLeft = {
      // Time calculations for days, hours, minutes and seconds
      days: Math.floor(difference / (1000 * 60 * 60 * 24)),
      hours: Math.floor((difference % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)),
      minutes: Math.floor((difference % (1000 * 60 * 60)) / (1000 * 60)),
      seconds: Math.floor((difference % (1000 * 60)) / 1000)
    };
  } else if (difference > 0 && type === 'hours') {
    const timeInSecs = Math.floor(difference / 1000);

    timeLeft = {
      hours: Math.floor(timeInSecs / 3600),
      minutes: Math.floor((timeInSecs % 3600) / 60),
      seconds: Math.floor((timeInSecs % 3600) % 60)
    };
  } else {
    timeLeft = {
      minutes: Math.floor((difference / 1000 / 60) % 60),
      seconds: Math.floor((difference / 1000) % 60)
    };
  }

  return timeLeft;
};

export const formatRelayPerson = (person: any): any => ({
  owner_pubkey: person.owner_pubkey,
  owner_alias: person.alias,
  owner_contact_key: person.contact_key,
  owner_route_hint: person.route_hint ?? '',
  description: person.description,
  extras: person.extras,
  price_to_meet: person.price_to_meet,
  img: person.img,
  tags: [],
  route_hint: person.route_hint
});

export type Roles =
  | 'EDIT ORGANIZATION'
  | 'ADD BOUNTY'
  | 'UPDATE BOUNTY'
  | 'DELETE BOUNTY'
  | 'PAY BOUNTY'
  | 'ADD USER'
  | 'UPDATE USER'
  | 'DELETE USER'
  | 'ADD ROLES'
  | 'ADD BUDGET'
  | 'WITHDRAW BUDGET'
  | 'VIEW REPORT';

export const ManageBountiesGroup = ['ADD BOUNTY', 'UPDATE BOUNTY', 'DELETE BOUNTY', 'PAY BOUNTY'];

export interface RolesCategory {
  name: string;
  roles: string[];
  status: boolean;
}

export const s_RolesCategories = [
  {
    name: 'Manage organization',
    roles: ['EDIT ORGANIZATION'],
    status: false
  },
  {
    name: 'Manage bounties',
    roles: ManageBountiesGroup,
    status: false
  },
  {
    name: 'Fund organization',
    roles: ['ADD BUDGET'],
    status: false
  },
  {
    name: 'Withdraw from organization',
    roles: ['WITHDRAW BUDGET'],
    status: false
  },
  {
    name: 'View transaction history',
    roles: ['VIEW REPORT'],
    status: false
  },
  {
    name: 'Update members',
    roles: ['ADD USER', 'UPDATE USER', 'DELETE USER', 'ADD ROLES'],
    status: false
  }
];

export const userHasRole = (
  bountyRoles: any[],
  userRoles: any[],
  role: Roles | string
): boolean => {
  if (bountyRoles.length) {
    const bountyRolesMap = {};
    const userRolesMap = {};

    bountyRoles.forEach((role: any) => {
      bountyRolesMap[role.name] = role.name;
    });

    userRoles.forEach((user: any) => {
      userRolesMap[user.role] = user.role;
    });

    if (bountyRolesMap.hasOwnProperty(role) && userRolesMap.hasOwnProperty(role)) {
      return true;
    }

    return false;
  }
  return false;
};

export const userHasManageBountyRoles = (bountyRoles: any[], userRoles: any[]): boolean => {
  let manageLength = ManageBountiesGroup.length;
  if (bountyRoles.length) {
    ManageBountiesGroup.forEach((role: string) => {
      const hasRole = userHasRole(bountyRoles, userRoles, role);
      if (hasRole) {
        manageLength--;
      }
    });
  }
  if (manageLength !== 0) {
    return false;
  }
  return true;
};

export const toCapitalize = (word: string): string => {
  if (!word.length) return word;

  const wordString = word.split(' ');
  const capitalizeStrings = wordString.map((w: string) => w[0].toUpperCase() + w.slice(1));

  const result = capitalizeStrings.join(' ');
  return result;
};

export const isInvoiceExpired = (paymentRequest: string): boolean => {
  // decode invoice to see if it has expired
  const decoded = LighningDecoder.decode(paymentRequest);
  const invoiceTimestamp = decoded.sections[4].value;
  const expiry = decoded.sections[8].value;
  const expired = invoiceTimestamp + expiry;

  if (expired * 1000 > Date.now()) {
    return false;
  }
  return true;
};

export const spliceOutPubkey = (userAddress: string): string => {
  if (userAddress.includes(':')) {
    const addArray = userAddress.split(':');
    const pubkey = addArray[0];

    return pubkey;
  }

  return userAddress;
};

export function handleDisplayRole(displayedRoles: RolesCategory[]) {
  // Set default data roles for first assign user
  const defaultRole = {
    'Manage bounties': true,
    'Fund organization': true,
    'Withdraw from organization': true,
    'View transaction history': true
  };

  const tempDataRole: { [id: string]: boolean } = {};
  const newDisplayedRoles = displayedRoles.map((role: RolesCategory) => {
    if (defaultRole[role.name]) {
      role.status = true;
      role.roles.forEach((dataRole: string) => (tempDataRole[dataRole] = true));
    }
    return role;
  });

  return { newDisplayedRoles, tempDataRole };
}

export async function userCanManageBounty(
  org_uuid: string | undefined,
  userPubkey: string | undefined,
  main: MainStore
): Promise<boolean> {
  if (org_uuid && userPubkey) {
    const userRoles = await main.getUserRoles(org_uuid, userPubkey);
    const org = await main.getUserOrganizationByUuid(org_uuid);
    if (org) {
      const isOrganizationAdmin = org.owner_pubkey === userPubkey;
      const userAccess =
        userHasManageBountyRoles(main.bountyRoles, userRoles) &&
        userHasRole(main.bountyRoles, userRoles, 'VIEW REPORT');
      return isOrganizationAdmin || userAccess;
    }
  }
  return false;
}
