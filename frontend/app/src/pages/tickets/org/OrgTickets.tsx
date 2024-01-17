import { EuiGlobalToastList, EuiLoadingSpinner } from '@elastic/eui';
import { observer } from 'mobx-react-lite';
import FirstTimeScreen from 'people/main/FirstTimeScreen';
import BountyHeader from 'people/widgetViews/BountyHeader';
import WidgetSwitchViewer from 'people/widgetViews/WidgetSwitchViewer';
import React, { useEffect, useState } from 'react';
import { useHistory } from 'react-router';
import { useParams } from 'react-router-dom';
import { colors } from '../../../config/colors';
import { useIsMobile } from '../../../hooks';
import { useStores } from '../../../store';
import { Body, Backdrop } from '../style';

function OrgBodyComponent() {
  const { main, ui } = useStores();
  const [loading, setLoading] = useState(true);
  const [showDropdown, setShowDropdown] = useState(false);
  const selectedWidget = 'wanted';
  const [scrollValue, setScrollValue] = useState<boolean>(false);
  const [checkboxIdToSelectedMap, setCheckboxIdToSelectedMap] = useState({});
  const [checkboxIdToSelectedMapLanguage, setCheckboxIdToSelectedMapLanguage] = useState({});
  const { uuid } = useParams<{ uuid: string; bountyId: string }>();

  const color = colors['light'];

  const history = useHistory();
  const isMobile = useIsMobile();

  useEffect(() => {
    (async () => {
      await main.getOpenGithubIssues();
      await main.getBadgeList();
      await main.getPeople();
      if (uuid) {
        await main.getOrganizationBounties(uuid, { page: 1, resetPage: true });
      }
      setLoading(false);
    })();
  }, [main, uuid]);

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

  const onPanelClick = (person: any, item: any) => {
    history.push(`/bounty/${item.id}`, { uuid });
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
                checkboxIdToSelectedMap={checkboxIdToSelectedMap}
                checkboxIdToSelectedMapLanguage={checkboxIdToSelectedMapLanguage}
                onPanelClick={onPanelClick}
                fromBountyPage={true}
                selectedWidget={selectedWidget}
                loading={loading}
              />
            </div>
          </div>
        </>
        {toastsEl}
      </Body>
    )
  );
}

export default observer(OrgBodyComponent);
