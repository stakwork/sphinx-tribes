import { uiStore } from '../store/ui';
import { RolesCategory } from './helpers-extended';

export const satToUsd = (amount: number = 0) => {
  if (!amount) amount = 0;
  const satExchange = uiStore.usdToSatsExchangeRate;
  const returnValue = (amount / satExchange).toFixed(2);

  if (returnValue === 'Infinity' || isNaN(parseFloat(returnValue))) {
    return '. . .';
  }

  return returnValue;
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

export * from './helpers-extended';
