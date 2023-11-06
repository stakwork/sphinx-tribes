/* eslint-disable @typescript-eslint/typedef */
import LighningDecoder from 'light-bolt11-decoder';
import { getHost } from '../config/host';
import { uiStore } from '../store/ui';

export const formatPrice = (amount = 0) => amount;

export const satToUsd = (amount = 0) => {
  if (!amount) amount = 0;
  const satExchange = uiStore.usdToSatsExchangeRate ?? 0;
  const returnValue = (amount / satExchange).toFixed(2);

  if (returnValue === 'Infinity' || isNaN(parseFloat(returnValue))) {
    return '. . .';
  }

  return returnValue;
};

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
  let numString = localeString.replace(/\D/g, '');

  let num = parseInt(numString);
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
  type: 'minutes' | 'days'
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

export const userHasRole = (bountyRoles: any[], userRoles: any[], role: Roles): boolean => {
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
