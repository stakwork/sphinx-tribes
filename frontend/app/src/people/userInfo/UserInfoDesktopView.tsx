import { Button, IconButton } from 'components/common';
import { AboutView } from 'people/widgetViews/AboutView';
import React, { useState } from 'react';
import ConnectCard from 'people/utils/ConnectCard';
import { observer } from 'mobx-react-lite';
import { UserInfoProps } from 'people/interfaces';
import { AboutWrap, Head, Img, Name, RowWrap } from './styles';
import { useUserInfo } from './hooks';

export const UserInfoDesktopView = observer(({ setShowSupport }: UserInfoProps) => {
  const { canEdit, goBack, userImg, owner_alias, logout, person, onEdit } = useUserInfo();
  const [showQR, setShowQR] = useState(false);
  return (
    <AboutWrap
      style={{
        width: 364,
        minWidth: 364,
        background: '#ffffff',
        color: '#000000',
        padding: 40,
        zIndex: 5,
        marginTop: canEdit ? 64 : 0,
        height: canEdit ? 'calc(100% - 64px)' : '100%',
        borderLeft: '1px solid #ebedef',
        borderRight: '1px solid #ebedef',
        boxShadow: '1px 2px 6px -2px rgba(0, 0, 0, 0.07)'
      }}
    >
      {canEdit && (
        <div
          style={{
            position: 'absolute',
            top: 0,
            left: 0,
            display: 'flex',
            background: '#ffffff',
            justifyContent: 'space-between',
            alignItems: 'center',
            width: 364,
            minWidth: 364,
            boxShadow: '0px 2px 6px rgba(0, 0, 0, 0.07)',
            borderBottom: 'solid 1px #ebedef',
            paddingRight: 10,
            height: 64,
            zIndex: 0
          }}
        >
          <Button color="clear" leadingIcon="arrow_back" text="Back" onClick={goBack} />
          <div />
        </div>
      )}

      {/* profile photo */}
      <Head>
        <div style={{ height: 35 }} />

        <Img src={userImg}>
          <IconButton
            iconStyle={{ color: '#5F6368' }}
            style={{
              zIndex: 2,
              width: '40px',
              height: '40px',
              padding: 0,
              background: '#ffffff',
              border: '1px solid #D0D5D8',
              boxSizing: 'border-box',
              borderRadius: 4
            }}
            icon={'qr_code_2'}
            onClick={() => setShowQR(true)}
          />
        </Img>

        <RowWrap>
          <Name>{owner_alias}</Name>
        </RowWrap>

        {/* only see buttons on other people's profile */}
        {canEdit ? (
          <RowWrap
            style={{
              marginBottom: 30,
              marginTop: 25,
              justifyContent: 'space-around'
            }}
          >
            <Button
              text="Edit Profile"
              onClick={onEdit}
              color="widget"
              height={42}
              style={{ fontSize: 13, background: '#f2f3f5' }}
              leadingIcon={'edit'}
              iconSize={15}
            />
            <Button
              text="Sign out"
              onClick={logout}
              height={42}
              style={{ fontSize: 13, color: '#3c3f41' }}
              iconStyle={{ color: '#8e969c' }}
              iconSize={15}
              color="white"
              leadingIcon="logout"
            />
          </RowWrap>
        ) : (
          <RowWrap
            style={{
              marginBottom: 30,
              marginTop: 25,
              justifyContent: 'space-between'
            }}
          >
            <Button
              text="Connect"
              onClick={() => setShowQR(true)}
              color="primary"
              height={42}
              width={120}
            />

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
      <AboutView canEdit={canEdit} {...person} />
      <ConnectCard
        dismiss={() => setShowQR(false)}
        modalStyle={{ top: -63, height: 'calc(100% + 64px)' }}
        person={person}
        visible={showQR}
      />
    </AboutWrap>
  );
});
