import { withCreateModal } from 'components/common';
import compose from 'compose-function';
import { withStores } from 'store';

export const withProviders = compose(withStores, withCreateModal);
