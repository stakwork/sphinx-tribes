import { useState, useCallback, useEffect } from 'react';
import { useStores } from 'store';

export const useIsSuperAdmin = () => {
  const { main, ui } = useStores();
  const [isSuperAdmin, setIsSuperAdmin] = useState(true);

  const getIsSuperAdmin = useCallback(async () => {
    const admin = await main.getSuperAdmin();
    setIsSuperAdmin(admin);
  }, [main]);

  useEffect(() => {
    if (ui.meInfo?.tribe_jwt) {
      getIsSuperAdmin();
    }
  }, [main, ui, getIsSuperAdmin]);

  return isSuperAdmin;
};
