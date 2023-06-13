import { getHost } from '../config/host';
import { uiStore } from '../store/ui';

export function formatPrice(amount: number) {
  return amount;
}

export function satToUsd(amount: number) {
  if (!amount) amount = 0;
  const satExchange = uiStore.usdToSatsExchangeRate ?? 0;
  const returnValue = (amount / satExchange).toFixed(2);

  if (returnValue === 'Infinity' || isNaN(parseFloat(returnValue))) {
    return '. . .';
  }

  return returnValue;
}

export const DollarConverter = (e) => {
  const dollarUSLocale = Intl.NumberFormat('en-US');
  return dollarUSLocale.format(formatPrice(e)).split(',').join(' ');
};

const host = getHost();
export function makeConnectQR(pubkey: string) {
  return `sphinx.chat://?action=person&host=${host}&pubkey=${pubkey}`;
}

export function extractRepoAndIssueFromIssueUrl(url: string) {
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
}

export function extractGithubIssue(
  person: { github_issues: Record<string, any> },
  repo: string,
  issue: string
) {
  const { github_issues } = person;
  const keyname = `${repo}/${issue}`;
  return (github_issues && github_issues[keyname]) || {};
}

export function extractGithubIssueFromUrl(
  person: { github_issues: Record<string, any> },
  url: string
) {
  try {
    const { repo, issue } = extractRepoAndIssueFromIssueUrl(url);
    return extractGithubIssue(person, repo, issue);
  } catch (e) {
    return {};
  }
}

export const randomString = (l: number): string =>
  Array.from(crypto.getRandomValues(new Uint8Array(l)), (byte) =>
    `0${(byte & 0xff).toString(16)}`.slice(-2)
  ).join('');

export function sendToRedirect(url: string) {
  const el = document.createElement('a');
  el.href = url;
  el.target = '_blank';
  el.click();
}

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
