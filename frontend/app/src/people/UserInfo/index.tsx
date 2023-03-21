import { useIsMobile } from 'hooks';
import React from 'react';
import { UserInfoDesktopView } from './UserInfoDesktopView';
import { UserInfoMobileView } from './UserInfoMobileView';

type UserInfoProps = { setShowSupport };

export const UserInfo = (props: UserInfoProps) => {
  const isMobile = useIsMobile();

  return isMobile ? <UserInfoMobileView {...props} /> : <UserInfoDesktopView {...props} />;
};
