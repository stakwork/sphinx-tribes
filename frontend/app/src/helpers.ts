import { getHost } from './host';
import { uiStore } from './store/ui';

export function formatPrice(amount) {
  return amount;
}

export function satToUsd(amount) {
  if (!amount) amount = 0;
  const satExchange = uiStore.usdToSatsExchangeRate ? uiStore.usdToSatsExchangeRate : 0;
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

export function extractGithubIssue(person, repo, issue) {
  const { github_issues } = person;
  const keyname = `${repo}/${issue}`;
  return (github_issues && github_issues[keyname]) || {};
}

export function extractRepoAndIssueFromIssueUrl(url: string) {
  let orgName = '';
  let repoName = '';
  let repo = '';
  let issue = '';

  // example: https://github.com/stakwork/sphinx-tribes/issues/206

  try {
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
    repo = `${orgName}/${repoName}`;
  } catch (e) {
    console.log('e', e);
  }

  return { repo, issue };
}

export function extractGithubIssueFromUrl(person, url) {
  const { github_issues } = person;
  const { repo, issue } = extractRepoAndIssueFromIssueUrl(url);
  const keyname = `${repo}/${issue}`;
  return (github_issues && github_issues[keyname]) || {};
}

export const randomString = (l: number): string => {
  return Array.from(crypto.getRandomValues(new Uint8Array(l)), (byte) => {
    return `0${(byte & 0xff).toString(16)}`.slice(-2);
  }).join('');
};
