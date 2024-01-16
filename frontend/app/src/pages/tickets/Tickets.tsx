import { EuiGlobalToastList, EuiLoadingSpinner } from '@elastic/eui';
import { observer } from 'mobx-react-lite';
import FirstTimeScreen from 'people/main/FirstTimeScreen';
import BountyHeader from 'people/widgetViews/BountyHeader';
import WidgetSwitchViewer from 'people/widgetViews/WidgetSwitchViewer';
import React, { useCallback, useEffect, useState } from 'react';
import { useHistory } from 'react-router';
import { queryLimit, defaultBountyStatus } from 'store/main';
import { colors } from '../../config/colors';
import { useIsMobile } from '../../hooks';
import { useStores } from '../../store';
import { Body, Backdrop } from './style';

// avoid hook within callback warning by renaming hooks
function BodyComponent() {
  const { main, ui } = useStores();
  const [loading, setLoading] = useState(true);
  const [showDropdown, setShowDropdown] = useState(false);
  const selectedWidget = 'wanted';
  const [scrollValue, setScrollValue] = useState<boolean>(false);
  const [checkboxIdToSelectedMap, setCheckboxIdToSelectedMap] = useState(defaultBountyStatus);
  const [checkboxIdToSelectedMapLanguage, setCheckboxIdToSelectedMapLanguage] = useState({});
  const [languageString, setLanguageString] = useState('');
  const [page, setPage] = useState<number>(1);
  const [currentItems, setCurrentItems] = useState<number>(queryLimit);
  const [totalBounties, setTotalBounties] = useState(0);

  const color = colors['light'];

  const history = useHistory();
  const isMobile = useIsMobile();

  useEffect(() => {
    (async () => {
      await main.getOpenGithubIssues();
      await main.getBadgeList();
      await main.getPeople();
      await main.getPeopleBounties({
        page: 1,
        resetPage: true,
        ...checkboxIdToSelectedMap,
        languages: languageString
      });
      setLoading(false);
    })();
  }, [main, checkboxIdToSelectedMap, languageString]);

  useEffect(() => {
    setCheckboxIdToSelectedMap({
      Open: true,
      Assigned: false,
      Paid: false
    });
  }, [loading]);

  useEffect(() => {
    if (ui.meInfo) {
      main.getTribesByOwner(ui.meInfo.owner_pubkey || '');
    }
  }, [main, ui.meInfo]);

  const getTotalBounties = useCallback(
    async (statusData: any) => {
      const totalBounties = await main.getTotalBountyCount(
        statusData.Open,
        statusData.Assigned,
        statusData.Paid
      );
      setTotalBounties(totalBounties);
    },
    [main]
  );

  useEffect(() => {
    getTotalBounties(checkboxIdToSelectedMap);
  }, [getTotalBounties]);

  const onChangeStatus = (optionId: any) => {
    const newCheckboxIdToSelectedMap = {
      ...checkboxIdToSelectedMap,
      ...{
        [optionId]: !checkboxIdToSelectedMap[optionId]
      }
    };
    setCheckboxIdToSelectedMap(newCheckboxIdToSelectedMap);
    // set to store
    main.setBountiesStatus(newCheckboxIdToSelectedMap);
    getTotalBounties(newCheckboxIdToSelectedMap);
    // set data to default
    setCurrentItems(queryLimit);
    setPage(1);
  };

  const onChangeLanguage = (optionId: any) => {
    const newCheckboxIdToSelectedMapLanguage = {
      ...checkboxIdToSelectedMapLanguage,
      ...{
        [optionId]: !checkboxIdToSelectedMapLanguage[optionId]
      }
    };

    setCheckboxIdToSelectedMapLanguage(newCheckboxIdToSelectedMapLanguage);

    const languageString = Object.keys(newCheckboxIdToSelectedMapLanguage).join(',');
    setLanguageString(languageString);
    main.setBountyLanguages(languageString);
  };

  const onPanelClick = (person: any, item: any) => {
    history.push(`/bounty/${item.id}`);
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

  if (!loading && isMobile) {
    return (
      <Body>
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
          <WidgetSwitchViewer
            checkboxIdToSelectedMap={checkboxIdToSelectedMap}
            checkboxIdToSelectedMapLanguage={checkboxIdToSelectedMapLanguage}
            onPanelClick={onPanelClick}
            fromBountyPage={true}
            selectedWidget={selectedWidget}
            loading={loading}
            totalBounties={totalBounties}
            currentItems={currentItems}
            setCurrentItems={setCurrentItems}
            page={page}
            setPage={setPage}
          />
        </div>

        {toastsEl}
      </Body>
    );
  }

  return (
    !loading && (
      <Body
        onScroll={(e: any) => {
          setScrollValue(e?.currentTarget?.scrollTop >= 20);
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
                onPanelClick={onPanelClick}
                checkboxIdToSelectedMap={checkboxIdToSelectedMap}
                checkboxIdToSelectedMapLanguage={checkboxIdToSelectedMapLanguage}
                fromBountyPage={true}
                selectedWidget={selectedWidget}
                loading={loading}
                totalBounties={totalBounties}
                currentItems={currentItems}
                setCurrentItems={setCurrentItems}
                page={page}
                setPage={setPage}
              />
            </div>
          </div>
        </>
        {toastsEl}
      </Body>
    )
  );
}

export default observer(BodyComponent);
