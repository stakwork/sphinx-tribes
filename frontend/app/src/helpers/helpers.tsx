import { uiStore } from '../store/ui';

// eslint-disable-next-line @typescript-eslint/no-inferrable-types
export const satToUsd = (amount: number = 0) => {
  if (!amount) amount = 0;
  const satExchange = uiStore.usdToSatsExchangeRate;
  const returnValue = (amount / satExchange).toFixed(2);

  if (returnValue === 'Infinity' || isNaN(parseFloat(returnValue))) {
    return '. . .';
  }

  return returnValue;
};

export * from './helpers-extended';
