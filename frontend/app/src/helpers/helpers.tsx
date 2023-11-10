import { uiStore } from 'stores/uiStore';
import * as extendedHelpers from './helpers-extended';

const satToUsd = (amount: number = 0) => {
  if (!amount) amount = 0;
  const satExchange = uiStore.usdToSatsExchangeRate;
  const returnValue = (amount / satExchange).toFixed(2);

  if (returnValue === 'Infinity' || isNaN(parseFloat(returnValue))) {
    return '. . .';
  }

  return returnValue;
};

export default { ...extendedHelpers, satToUsd };
