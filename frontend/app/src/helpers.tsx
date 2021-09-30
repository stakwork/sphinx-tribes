
import { getHost } from "./host";

export function formatPrice(amount) {
    return amount
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

