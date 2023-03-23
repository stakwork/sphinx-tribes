import React from 'react';
import { Panel, Tabs, Tab, Head, Img, RowWrap, Name, Sleeve, Counter } from './style';
import FocusedView from '..//main/focusView';
import { Button, IconButton, Modal } from '../../components/common';

export default function MobileView(props: any) {
  const {
    isMobile,
    canEdit,
    logout,
    mediumPic,
    defaultPic,
    person,
    tabs,
    showFocusView,
    focusIndex,
    setFocusIndex,
    setShowFocusView,
    selectedWidget,
    extras,
    goBack,
    switchWidgets,
    qrString,
    owner_alias,
    setShowSupport,
    renderEditButton,
    newSelectedWidget,
    renderWidgets
  } = props;
  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        width: '100%',
        overflow: 'auto',
        height: '100%'
      }}
    >
      <Panel isMobile={isMobile} style={{ paddingBottom: 0, paddingTop: 80 }}>
        <div
          style={{
            position: 'absolute',
            top: 20,
            left: 0,
            display: 'flex',
            justifyContent: 'space-between',
            width: '100%',
            padding: '0 20px'
          }}
        >
          <IconButton onClick={goBack} icon="arrow_back" />
          {canEdit ? (
            <>
              <Button
                text="Edit Profile"
                onClick={() => {
                  switchWidgets('about');
                  setShowFocusView(true);
                }}
                color="white"
                height={42}
                style={{
                  fontSize: 13,
                  color: '#3c3f41',
                  border: 'none',
                  marginLeft: 'auto'
                }}
                leadingIcon={'edit'}
                iconSize={15}
              />
              <Button
                text="Sign out"
                onClick={logout}
                height={42}
                style={{
                  fontSize: 13,
                  color: '#3c3f41',
                  border: 'none',
                  margin: 0,
                  padding: 0
                }}
                iconStyle={{ color: '#8e969c' }}
                iconSize={20}
                color="white"
                leadingIcon="logout"
              />
            </>
          ) : (
            <div />
          )}
        </div>

        {/* profile photo */}
        <Head>
          <Img src={mediumPic || defaultPic} />
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

        <Tabs>
          {tabs &&
            Object.keys(tabs).map((name, i) => {
              const t = tabs[name];
              const { label } = t;
              const selected = name === newSelectedWidget;
              const hasExtras = extras && extras[name] && extras[name].length > 0;
              const count: any = hasExtras
                ? extras[name].filter((f) => {
                  if ('show' in f) {
                    // show has a value
                    if (!f.show) return false;
                  }
                  // if no value default to true
                  return true;
                }).length
                : null;

              return (
                <Tab
                  key={i}
                  selected={selected}
                  onClick={() => {
                    switchWidgets(name);
                  }}
                >
                  {label} {count && <Counter>{count}</Counter>}
                </Tab>
              );
            })}
        </Tabs>
      </Panel>

      <Sleeve>
        {renderEditButton({})}
        {renderWidgets('')}
        <div style={{ height: 60 }} />
      </Sleeve>

      <Modal fill visible={showFocusView}>
        <FocusedView
          person={person}
          canEdit={canEdit}
          selectedIndex={focusIndex}
          config={tabs[selectedWidget] && tabs[selectedWidget]}
          onSuccess={() => {
            console.log('success');
            setFocusIndex(-1);
          }}
          goBack={() => {
            setShowFocusView(false);
            setFocusIndex(-1);
          }}
        />
      </Modal>
    </div>
  );
}
