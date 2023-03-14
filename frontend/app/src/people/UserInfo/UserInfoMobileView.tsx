import React from 'react';
import { Head, AboutWrap, Img, Name, RowWrap } from './styles';
import { useUserInfo } from './hooks';
import { Button } from 'components/common';

export const UserInfoMobileView = ({ setShowSupport }) => {
  const { canEdit, goBack, userImg, owner_alias, logout, person, qrString } = useUserInfo();
  return (
    <Head>
      <Img src={userImg} />
      <RowWrap>
        <Name>{owner_alias}</Name>
      </RowWrap>

      {/* only see buttons on other people's profile */}
      {canEdit ? (
        <div style={{ height: 40 }} />
      ) : (
        <RowWrap style={{ marginBottom: 30, marginTop: 25 }}>
          <a href={qrString}>
            <Button
              text="Connect"
              onClick={(e) => e.stopPropagation()}
              color="primary"
              height={42}
              width={120}
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
  );
};
