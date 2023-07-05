import { Button } from 'components/common';
import { observer } from 'mobx-react-lite';
import React, { useState } from 'react';
import { UserInfoProps } from 'people/interfaces';
import ConnectCard from 'people/utils/ConnectCard';
import { useUserInfo } from './hooks';
import { Head, Img, Name, RowWrap } from './styles';
import { HeaderMobile } from './MobileHeader';

export const UserInfoMobileView = observer(({ setShowSupport }: UserInfoProps) => {
  const { canEdit, goBack, userImg, owner_alias, logout, qrString, onEdit, person } = useUserInfo();
  const [showQR, setShowQR] = useState(false);
  return (
    <>
      <HeaderMobile canEdit={canEdit} goBack={goBack} logout={logout} onEdit={onEdit} />
      <Head>
        <Img src={userImg} />
        <RowWrap>
          <Name>{owner_alias}</Name>
        </RowWrap>

        {canEdit ? (
          <div style={{ height: 40 }} />
        ) : (
          <RowWrap style={{ marginBottom: 30, marginTop: 25 }}>
            <a href={qrString}>
              <Button
                text="Connect"
                color="primary"
                height={42}
                width={120}
                onClick={() => setShowQR(true)}
              />
            </a>

            <div style={{ width: 15 }} />

            <Button
              text="Send Tip"
              color="link"
              height={42}
              width={120}
              onClick={() => setShowSupport(true)}
            />
          </RowWrap>
        )}
      </Head>
      <ConnectCard
        dismiss={() => setShowQR(false)}
        modalStyle={{ top: -63, height: 'calc(100% + 64px)' }}
        person={person}
        visible={showQR}
      />
    </>
  );
});
