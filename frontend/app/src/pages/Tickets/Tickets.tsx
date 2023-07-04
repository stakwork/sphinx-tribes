import { EuiGlobalToastList, EuiLoadingSpinner } from '@elastic/eui';
import FirstTimeScreen from 'people/main/firstTimeScreen';
import PageLoadSpinner from 'people/utils/pageLoadSpinner';
import BountyHeader from 'people/widgetViews/bountyHeader';
import WidgetSwitchViewer from 'people/widgetViews/widgetSwitchViewer';
import React, { useEffect, useState } from 'react';
import { useHistory } from 'react-router';
import styled from 'styled-components';
import { observer } from 'mobx-react-lite';
import { colors } from '../../config/colors';
import { useIsMobile, usePageScroll } from '../../hooks';
import { useStores } from '../../store';

// avoid hook within callback warning by renaming hooks
const Body = styled.div`
  flex: 1;
  height: calc(100% - 105px);
  width: 100%;
  overflow: auto;
  display: flex;
  flex-direction: column;
`;

const Backdrop = styled.div`
  position: fixed;
  z-index: 1;
  background: rgba(0, 0, 0, 70%);
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
`;

export const Spacer = styled.div`
  display: flex;
  min-height: 10px;
  min-width: 100%;
  height: 10px;
  width: 100%;
`;

const getPageScroll = usePageScroll;

function BodyComponent() {
  const { main, ui } = useStores();
  const [loading, setLoading] = useState(true);
  const [showDropdown, setShowDropdown] = useState(false);
  const selectedWidget = 'wanted';
  const [scrollValue, setScrollValue] = useState<boolean>(false);
  const [checkboxIdToSelectedMap, setCheckboxIdToSelectedMap] = useState({});
  const [checkboxIdToSelectedMapLanguage, setCheckboxIdToSelectedMapLanguage] = useState({});
  const color = colors['light'];
  const { peopleWantedsPageNumber } = ui;

  const history = useHistory();
  const isMobile = useIsMobile();

  useEffect(() => {
    main.getOpenGithubIssues();
    main.getBadgeList();
    main.getPeople();
  }, [main]);

  useEffect(() => {
    if (ui.meInfo) {
      main.getTribesByOwner(ui.meInfo.owner_pubkey || '');
    }
  }, [main, ui.meInfo]);

  useEffect(() => {
    (async () => {
      await main.getPeopleWanteds({ page: 1, resetPage: true });
      setLoading(false);
    })();
  }, [main]);

  const onChangeStatus = (optionId: any) => {
    const newCheckboxIdToSelectedMap = {
      ...checkboxIdToSelectedMap,
      ...{
        [optionId]: !checkboxIdToSelectedMap[optionId]
      }
    };
    setCheckboxIdToSelectedMap(newCheckboxIdToSelectedMap);
  };

  const onChangeLanguage = (optionId: any) => {
    const newCheckboxIdToSelectedMapLanguage = {
      ...checkboxIdToSelectedMapLanguage,
      ...{
        [optionId]: !checkboxIdToSelectedMapLanguage[optionId]
      }
    };
    setCheckboxIdToSelectedMapLanguage(newCheckboxIdToSelectedMapLanguage);
  };

  async function loadMore(direction: number) {
    let currentPage = 1;
    currentPage = peopleWantedsPageNumber;
    let newPage = currentPage + direction;
    if (newPage < 1) newPage = 1;
    try {
      await main.getPeopleWanteds({ page: newPage });
    } catch (e) {
      console.log('load failed', e);
    }
  }
  const loadForwardFunc = () => loadMore(1);
  const loadBackwardFunc = () => loadMore(-1);
  const { loadingTop, loadingBottom, handleScroll } = getPageScroll(
    loadForwardFunc,
    loadBackwardFunc
  );

  const onPanelClick = (person: any, item: any) => {
    history.replace({
      pathname: history?.location?.pathname,
      search: `?owner_id=${person.owner_pubkey}&created=${item.created}`,
      state: {
        owner_id: person.owner_pubkey,
        created: item.created
      }
    });
  };

  if (loading) {
    return (
      <Body style={{ justifyContent: 'center', alignItems: 'center' }}>
        <EuiLoadingSpinner size="xl" />
      </Body>
    );
  }

  const showFirstTime = ui.meInfo && ui.meInfo.id === 0;

  if (showFirstTime) {
    return <FirstTimeScreen />;
  }

  const toastsEl = (
    <EuiGlobalToastList
      toasts={ui.toasts}
      dismissToast={() => ui.setToasts([])}
      toastLifeTimeMs={3000}
    />
  );

  if (isMobile) {
    return (
      <Body onScroll={handleScroll}>
        <div
          style={{
            width: '100%',
            padding: '8px 0px',
            boxShadow: `0 0 6px 0 ${color.black100}`,
            zIndex: 2,
            position: 'relative',
            background: color.pureWhite,
            borderBottom: `1px solid ${color.black100}`
          }}
        >
          <BountyHeader
            selectedWidget={selectedWidget}
            scrollValue={scrollValue}
            onChangeStatus={onChangeStatus}
            onChangeLanguage={onChangeLanguage}
            checkboxIdToSelectedMap={checkboxIdToSelectedMap}
            checkboxIdToSelectedMapLanguage={checkboxIdToSelectedMapLanguage}
          />
        </div>

        {showDropdown && <Backdrop onClick={() => setShowDropdown(false)} />}
        <div style={{ width: '100%' }}>
          <PageLoadSpinner show={loadingTop} />
          <WidgetSwitchViewer
            checkboxIdToSelectedMap={checkboxIdToSelectedMap}
            checkboxIdToSelectedMapLanguage={checkboxIdToSelectedMapLanguage}
            onPanelClick={onPanelClick}
            fromBountyPage={true}
            selectedWidget={selectedWidget}
            loading={loading}
          />
          <PageLoadSpinner noAnimate show={loadingBottom} />
        </div>

        {toastsEl}
      </Body>
    );
  }
  return (
    <Body
      onScroll={(e: any) => {
        setScrollValue(e?.currentTarget?.scrollTop >= 20);
        handleScroll(e);
      }}
      style={{
        background: color.grayish.G950,
        height: 'calc(100% - 65px)'
      }}
    >
      <div
        style={{
          minHeight: '32px'
        }}
      />

      <BountyHeader
        selectedWidget={selectedWidget}
        scrollValue={scrollValue}
        onChangeStatus={onChangeStatus}
        onChangeLanguage={onChangeLanguage}
        checkboxIdToSelectedMap={checkboxIdToSelectedMap}
        checkboxIdToSelectedMapLanguage={checkboxIdToSelectedMapLanguage}
      />

      <>
        <div
          style={{
            width: '100%',
            display: 'flex',
            flexWrap: 'wrap',
            height: '100%',
            justifyContent: 'flex-start',
            alignItems: 'flex-start',
            padding: '0px 20px 20px 20px'
          }}
        >
          <PageLoadSpinner show={loadingTop} />
          <div
            style={{
              width: '100%',
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              height: '100%'
            }}
          >
            <WidgetSwitchViewer
              checkboxIdToSelectedMap={checkboxIdToSelectedMap}
              checkboxIdToSelectedMapLanguage={checkboxIdToSelectedMapLanguage}
              onPanelClick={onPanelClick}
              fromBountyPage={true}
              selectedWidget={selectedWidget}
              loading={loading}
            />
          </div>
          <PageLoadSpinner noAnimate show={loadingBottom} />
        </div>
      </>
      {toastsEl}
    </Body>
  );
}

export default observer(BodyComponent);
