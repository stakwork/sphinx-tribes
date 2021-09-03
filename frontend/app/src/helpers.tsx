
import { getHost } from "./host";

export function formatPrice(amount) {
    return amount
}

const host = getHost();
export function makeConnectQR(pubkey: string) {
    return `sphinx.chat://?action=person&host=${host}&pubkey=${pubkey}`;
}