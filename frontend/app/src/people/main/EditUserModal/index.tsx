import { useIsMobile } from 'hooks';
import { observer } from 'mobx-react-lite';
import React from 'react';
import { useStores } from 'store';
import { EditUserDesktopView } from './EditUserDesktopView';
import { EditUserMobileView } from './EditUserMobileView';

export const EditUserModal = observer(() => {
  const isMobile = useIsMobile();
  const { modals } = useStores();

  if (!modals.userEditModal) {
    return null;
  }

  return <>{isMobile ? <EditUserMobileView /> : <EditUserDesktopView />}</>;
});
