import { Button } from 'components/common';
import { observer } from 'mobx-react-lite';
import React from 'react';
import { useUserInfo } from './hooks';
import { Head, Img, Name, RowWrap } from './styles';

export const UserInfoMobileView = observer(({ setShowSupport }: any) => {
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
});
