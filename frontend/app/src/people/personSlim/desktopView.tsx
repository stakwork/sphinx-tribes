import React from 'react';
import {
  Tabs,
  Tab,
  Head,
  Img,
  RowWrap,
  Name,
  Sleeve,
  Counter,
  PeopleList,
  DBack,
  PeopleScroller,
  AboutWrap
} from './style';
import FocusedView from '..//main/focusView';
import { Button, IconButton, Modal, SearchTextInput } from '../../components/common';
import { LoaderBottom, LoaderTop } from '../personSlim/component';
import Person from '../person';
import ConnectCard from '../utils/connectCard';
import NoResults from '../utils/noResults';
import { useStores } from '../../store';

export default function DesktopView(props: any) {
  const {
    isMobile,
    canEdit,
    logout,
    mediumPic,
    defaultPic,
    person,
    personId,
    tabs,
    showFocusView,
    focusIndex,
    setFocusIndex,
    setShowFocusView,
    selectedWidget,
    extras,
    goBack,
    switchWidgets,
    owner_alias,
    setShowSupport,
    renderEditButton,
    newSelectedWidget,
    renderWidgets,
    nextIndex,
    prevIndex,
    setShowQR,
    handleScroll,
    people,
    loadingTop,
    loadingBottom,
    showQR,
    fullSelectedWidget,
    hasWidgets,
    selectPersonWithinFocusView,
    queryLimit
  } = props;

  const focusedDesktopModalStyles = newSelectedWidget
    ? {
        ...tabs[newSelectedWidget]?.modalStyle
      }
    : {};

  const { ui } = useStores();

  return (
    <div
      style={{
        display: 'flex',
        width: '100%',
        height: '100%'
      }}
    >
      {!canEdit && (
        <PeopleList>
          <DBack>
            <Button color="clear" leadingIcon="arrow_back" text="Back" onClick={goBack} />

            <SearchTextInput
              small
              name="search"
              type="search"
              placeholder="Search"
              value={ui.searchText}
              style={{
                width: 120,
                height: 40,
                border: '1px solid #DDE1E5',
                background: '#fff'
              }}
              onChange={(e) => {
                ui.setSearchText(e);
              }}
            />
          </DBack>

          <PeopleScroller
            style={{ width: '100%', overflowY: 'auto', height: '100%' }}
            onScroll={handleScroll}
          >
            <LoaderTop loadingTop={loadingTop} />

            {people?.length ? (
              people.map((t) => (
                <Person
                  {...t}
                  key={t.id}
                  selected={personId === t.id}
                  hideActions={true}
                  small={true}
                  select={selectPersonWithinFocusView}
                />
              ))
            ) : (
              <NoResults />
            )}

            {/* make sure you can always scroll ever with too few people */}
            {people?.length < queryLimit && <div style={{ height: 400 }} />}
          </PeopleScroller>

          <LoaderBottom loadingBottom={loadingBottom} />
        </PeopleList>
      )}

      <AboutWrap
        style={{
          width: 364,
          minWidth: 364,
          background: '#ffffff',
          color: '#000000',
          padding: 40,
          zIndex: 5,
          // height: '100%',
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

          <Img src={mediumPic || defaultPic}>
            <IconButton
              iconStyle={{ color: '#5F6368' }}
              style={{
                zIndex: 2,
                width: '40px',
                height: '40px',
                padding: 0,
                background: '#ffffff',
                border: '1px solid #D0D5D8',
                boxSizing: 'borderBox',
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
                onClick={() => {
                  switchWidgets('about');
                  setShowFocusView(true);
                }}
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

        {renderWidgets('about')}
      </AboutWrap>

      <div
        style={{
          width: canEdit ? 'calc(100% - 365px)' : 'calc(100% - 628px)',
          minWidth: 250,
          zIndex: canEdit ? 6 : 4
        }}
      >
        <Tabs
          style={{
            background: '#fff',
            padding: '0 20px',
            borderBottom: 'solid 1px #ebedef',
            boxShadow: canEdit
              ? '0px 2px 0px rgba(0, 0, 0, 0.07)'
              : '0px 2px 6px rgba(0, 0, 0, 0.07)'
          }}
        >
          {tabs &&
            Object.keys(tabs).map((name, i) => {
              if (name === 'about') return <div key={i} />;
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
                  style={{ height: 64, alignItems: 'center' }}
                  selected={selected}
                  onClick={() => {
                    switchWidgets(name);
                  }}
                >
                  {label}
                  {count > 0 && <Counter>{count}</Counter>}
                </Tab>
              );
            })}
        </Tabs>

        <div
          style={{
            padding: 20,
            height: 'calc(100% - 63px)',
            background: '#F2F3F5',
            overflowY: 'auto',
            position: 'relative'
          }}
        >
          {renderEditButton({ marginBottom: 15 })}
          {/* <div style={{ height: 15 }} /> */}
          <Sleeve
            style={{
              display: 'flex',
              alignItems: 'flex-start',
              justifyContent:
                fullSelectedWidget && fullSelectedWidget.length > 0 ? 'flex-start' : 'center',
              flexWrap: 'wrap',
              height: !hasWidgets() ? 'inherit' : '',
              paddingTop: !hasWidgets() ? 30 : 0
            }}
          >
            {renderWidgets('')}
          </Sleeve>
          <div style={{ height: 60 }} />
        </div>
      </div>

      <ConnectCard
        dismiss={() => setShowQR(false)}
        modalStyle={{ top: -63, height: 'calc(100% + 64px)' }}
        person={person}
        visible={showQR}
      />

      <Modal
        visible={showFocusView}
        style={{
          height: '100%'
        }}
        envStyle={{
          marginTop: isMobile ? 64 : 0,
          borderRadius: 0,
          background: '#fff',
          height: '100%',
          width: '60%',
          minWidth: 500,
          maxWidth: 602,
          zIndex: 20, //minHeight: 300,
          ...focusedDesktopModalStyles
        }}
        nextArrow={nextIndex}
        prevArrow={prevIndex}
        overlayClick={() => {
          setShowFocusView(false);
          setFocusIndex(-1);
          if (selectedWidget === 'about') switchWidgets('badges');
        }}
        bigClose={() => {
          setShowFocusView(false);
          setFocusIndex(-1);
          if (selectedWidget === 'about') switchWidgets('badges');
        }}
      >
        <FocusedView
          person={person}
          canEdit={canEdit}
          selectedIndex={focusIndex}
          config={tabs[selectedWidget] && tabs[selectedWidget]}
          onSuccess={() => {
            setFocusIndex(-1);
            if (selectedWidget === 'about') switchWidgets('badges');
          }}
          goBack={() => {
            setShowFocusView(false);
            setFocusIndex(-1);
            if (selectedWidget === 'about') switchWidgets('badges');
          }}
        />
      </Modal>
    </div>
  );
}
