
import { getHost } from "./host";
import { uiStore } from "./store/ui";

export function formatPrice(amount) {
    return amount
}

export function satToUsd(amount) {
    if (!amount) amount = 0
    return (amount / (uiStore.usdToSatsExchangeRate || 0)).toFixed(2) + ' USD'
}

const host = getHost();
export function makeConnectQR(pubkey: string) {
    return `sphinx.chat://?action=person&host=${host}&pubkey=${pubkey}`;
}

export function extractGithubIssue(person, repo, issue) {
    const { github_issues } = person
    const keyname = repo + '/' + issue
    return (github_issues && github_issues[keyname]) || {}
}

