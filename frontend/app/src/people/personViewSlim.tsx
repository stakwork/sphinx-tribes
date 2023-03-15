/* eslint-disable func-style */
import React, { useCallback, useEffect, useState } from 'react';
import styled from 'styled-components';
import { useStores } from '../store';

import AboutView from './widgetViews/aboutView';
import BlogView from './widgetViews/blogView';
import OfferView from './widgetViews/offerView';
import PostView from './widgetViews/postView';
import SupportMeView from './widgetViews/supportMeView';
import TwitterView from './widgetViews/twitterView';
import WantedView from './widgetViews/wantedView';

import { useHistory, useLocation } from 'react-router';
import { Button, IconButton, Modal } from '../components/common';
import { meSchema } from '../components/form/schema';
import { useIsMobile } from '../hooks';
import { PeopleList } from './PeopleList';
import { UserInfo } from './UserInfo';
import { Widget } from './main/types';
import Badges from './utils/badges';
import { widgetConfigs } from './utils/constants';
import NoneSpace from './utils/noneSpace';
import PageLoadSpinner from './utils/pageLoadSpinner';
import { PostBounty } from './widgetViews/postBounty';
import { useModalsVisibility } from 'store/modals';

export default function PersonView({ loading = false }) {
  // on this screen, there will always be a pubkey in the url, no need for personId
  const { main, ui } = useStores();
  const modals = useModalsVisibility();
  const { meInfo } = ui || {};
  const history = useHistory();
  const location = useLocation();

  const personId = ui.selectedPerson;
  function goBack() {
    ui.setSelectingPerson(0);
    history.goBack();
  }

  const [loadingPerson, setLoadingPerson]: any = useState(false);
  const [loadedPerson, setLoadedPerson]: any = useState(null);

  const pathname = history?.location?.pathname;
  // FOR PEOPLE VIEW
  let person: any = main.people && main.people.length && main.people.find((f) => f.id === personId);

  // migrating to loading person on person view load
  if (loadedPerson) {
    person = loadedPerson;
  }

  // if i select myself, fill person with meInfo
  if (personId === ui.meInfo?.id) {
    person = {
      ...ui.meInfo
    };
  }

  const people: any = (main.people && main.people.filter((f) => !f.hide)) || [];

  const { id, img, owner_alias, extras, owner_pubkey } = person || {};

  let { description } = person || {};

  // backend is adding 'description' to empty descriptions, short term fix
  if (description === 'description') description = '';

  const canEdit = id === meInfo?.id;
  const isMobile = useIsMobile();

  const initialWidget = !isMobile || canEdit ? 'badges' : 'about';

  const [selectedWidget, setSelectedWidget] = useState<Widget>(initialWidget);
  const [newSelectedWidget, setNewSelectedWidget] = useState<Widget>(initialWidget);
  const [focusIndex, setFocusIndex] = useState(-1);
  const [showSupport, setShowSupport] = useState(false);

  
  const [showFocusView, setShowFocusView] = useState(false);

  // if no people, load people on mount
  useEffect(() => {
    if (!people.length) main.getPeople({ page: 1, resetPage: true });
  }, [main, people.length]);

  // fill state from url
  const doDeeplink = useCallback(async () => {
    console.log('personviewslim: doDeeplink', pathname);
    if (pathname) {
      const splitPathname = pathname?.split('/');
      // eslint-disable-next-line prefer-destructuring
      const personPubkey: string = splitPathname[2];
      if (personPubkey) {
        setLoadingPerson(true);
        const p = await main.getPersonByPubkey(personPubkey);
        setLoadedPerson(p);
        setLoadingPerson(false);

        const search = location?.search;

        // deeplink for widgets
        const widgetName: any = new URLSearchParams(search).get('widget');
        const widgetTimestamp: any = new URLSearchParams(search).get('timestamp');

        if (widgetName) {
          setNewSelectedWidget(widgetName);
          setSelectedWidget(widgetName);
          if (widgetTimestamp) {
            const thisExtra = p?.extras && p?.extras[widgetName];
            const thisItemIndex =
              thisExtra &&
              thisExtra.length &&
              thisExtra.findIndex((f) => f.created === parseInt(widgetTimestamp));
            if (thisItemIndex > -1) {
              // select it!
              setFocusIndex(thisItemIndex);
              setShowFocusView(true);
            }
          }
        }
      }
    }
  }, [location?.search, main, pathname]);

  // deeplink load person
  useEffect(() => {
    doDeeplink();
  }, [doDeeplink, pathname]);

  const updatePath = useCallback(
    (name) => {
      history.push(`${location.pathname}?widget=${name}`);
    },
    [history, location.pathname]
  );

  function updatePathIndex(timestamp) {
    history.push(`${location.pathname}?widget=${selectedWidget}&timestamp=${timestamp}`);
  }

  const switchWidgets = useCallback(
    (name) => {
      setNewSelectedWidget(name);
      setSelectedWidget(name);
      updatePath(name);
      setShowFocusView(false);
      setFocusIndex(-1);
    },

    [updatePath]
  );

  useEffect(() => {
    if (ui.personViewOpenTab) {
      switchWidgets(ui.personViewOpenTab);
      ui.setPersonViewOpenTab('');
    }
  }, [switchWidgets, ui, ui.personViewOpenTab]);

  function logout() {
    ui.setEditMe(false);
    ui.setMeInfo(null);
    main.getPeople({ resetPage: true });
    goBack();
  }

  if (loading) return <div>Loading...</div>;

  let widgetSchemas: any = meSchema.find((f) => f.name === 'extras');
  if (widgetSchemas && widgetSchemas.extras) {
    widgetSchemas = widgetSchemas && widgetSchemas.extras;
  }

  const fullSelectedWidget: any = extras && selectedWidget ? extras[selectedWidget] : null;

  // we do this because sometimes the widgets are empty arrays
  const filteredExtras = extras && { ...extras };
  if (filteredExtras) {
    const emptyArrayKeys = [''];

    Object.keys(filteredExtras).forEach((name) => {
      const p = extras && extras[name];
      if (Array.isArray(p) && !p.length) {
        emptyArrayKeys.push(name);
      }
      const thisSchema = widgetSchemas && widgetSchemas.find((e) => e.name === name);
      if (filteredExtras && thisSchema && thisSchema.single) {
        delete filteredExtras[name];
      }
    });

    emptyArrayKeys.forEach((e) => {
      if (filteredExtras && e) delete filteredExtras[e];
    });
  }

  const tabs = widgetConfigs;

  function hasWidgets() {
    let has = false;
    if (fullSelectedWidget && fullSelectedWidget.length) {
      has = true;
    }
    if (selectedWidget === 'badges') {
      has = true;
    }
    return has;
  }

  function renderWidgets(name: string) {
    if (name) {
      switch (name) {
        case 'about':
          return <AboutView canEdit={canEdit} {...person} />;
        case 'post':
          return wrapIt(<PostView {...fullSelectedWidget} person={person} />);
        case 'twitter':
          return wrapIt(<TwitterView {...fullSelectedWidget} person={person} />);
        case 'supportme':
          return wrapIt(<SupportMeView {...fullSelectedWidget} person={person} />);
        case 'offer':
          return wrapIt(<OfferView {...fullSelectedWidget} person={person} />);
        case 'wanted':
          return wrapIt(<WantedView {...fullSelectedWidget} person={person} />);
        case 'blog':
          return wrapIt(<BlogView {...fullSelectedWidget} person={person} />);
        default:
          return wrapIt(<></>);
      }
    }
    if (!selectedWidget) {
      return <div style={{ height: 200 }} />;
    }

    if (selectedWidget === 'badges') {
      return <Badges person={person} />;
    }

    const widgetSchema: any =
      (widgetSchemas && widgetSchemas.find((f) => f.name === selectedWidget)) || {};
    const { single } = widgetSchema;

    function wrapIt(child) {
      if (single) {
        return <Panel isMobile={isMobile}>{child}</Panel>;
      }

      const elementArray: any = [];

      const panelStyles = isMobile
        ? {
            minHeight: 132
          }
        : {
            maxWidth: 291,
            minWidth: 291,
            marginRight: 20,
            marginBottom: 20,
            minHeight: 472
          };

      fullSelectedWidget &&
        fullSelectedWidget.forEach((s, i) => {
          if (!canEdit && 'show' in s && s.show === false) {
            // skip hidden items
            return;
          }

          const conditionalStyles =
            !isMobile && s?.paid
              ? {
                  border: '1px solid #dde1e5',
                  boxShadow: 'none'
                }
              : {};

          elementArray.push(
            <Panel
              isMobile={isMobile}
              key={i}
              onClick={() => {
                setShowFocusView(true);
                setFocusIndex(i);
                if (s.created) updatePathIndex(s.created);
              }}
              style={{
                ...panelStyles,
                ...conditionalStyles,
                cursor: 'pointer',
                padding: 0,
                overflow: 'hidden'
              }}
            >
              {React.cloneElement(child, { ...s })}
            </Panel>
          );
        });
      const noneKey = canEdit ? 'me' : 'otherUser';
      const noneSpaceProps = tabs[selectedWidget]?.noneSpace[noneKey];

      const panels: any = elementArray.length ? (
        <div style={{ width: '100%', display: 'flex', flexDirection: 'column' }}>
          {person?.owner_pubkey === ui?.meInfo?.pubkey && selectedWidget === 'wanted' && (
            <div
              style={{
                width: '100%',
                display: 'flex',
                justifyContent: 'flex-end',
                paddingBottom: '16px'
              }}
            >
              <PostBounty widget={selectedWidget} />
            </div>
          )}
          <div style={{ width: '100%', display: 'flex', flexDirection: 'row', flexWrap: 'wrap' }}>
            {elementArray}
          </div>
        </div>
      ) : (
        <div
          style={{
            width: '100%'
          }}
        >
          <NoneSpace
            small
            Button={
              canEdit && (
                <PostBounty
                  title={noneSpaceProps.buttonText}
                  buttonProps={{
                    leadingIcon: noneSpaceProps.buttonIcon,
                    color: 'secondary'
                  }}
                  widget={selectedWidget}
                  onSucces={() => {
                    if (selectedWidget === 'about') switchWidgets('badges');
                  }}
                  onGoBack={() => {
                    if (selectedWidget === 'about') switchWidgets('badges');
                  }}
                />
              )
            }
            {...tabs[selectedWidget]?.noneSpace[noneKey]}
          />
        </div>
      );

      console.log('elementArray', elementArray.length);

      return (
        <>
          <PageLoadSpinner show={loadingPerson} />
          {panels}
        </>
      );
    }

    switch (selectedWidget) {
      case 'about':
        return (
          <Panel isMobile={isMobile}>
            <AboutView {...person} />
          </Panel>
        );
      case 'post':
        return wrapIt(<PostView {...fullSelectedWidget} person={person} />);
      case 'twitter':
        return wrapIt(<TwitterView {...fullSelectedWidget} person={person} />);
      case 'supportme':
        return wrapIt(<SupportMeView {...fullSelectedWidget} person={person} />);
      case 'offer':
        return wrapIt(<OfferView {...fullSelectedWidget} person={person} />);
      case 'wanted':
        return wrapIt(<WantedView {...fullSelectedWidget} person={person} />);
      case 'blog':
        return wrapIt(<BlogView {...fullSelectedWidget} person={person} />);
      default:
        return wrapIt(<></>);
    }
  }

  function renderEditButton(style: any) {
    if (!canEdit || !selectedWidget) return <div />;

    if (selectedWidget === 'badges') return <div />;

    // don't return button if there are no items in list, the button is returned elsewhere
    if (selectedWidget !== 'about') {
      if (!fullSelectedWidget || (fullSelectedWidget && fullSelectedWidget.length < 1))
        return <div />;
    }
  }

  const mediumPic = img;

  function renderMobileView() {
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
                    modals.setUserEditModal(true);
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
          <UserInfo  setShowSupport={setShowSupport} />
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
      </div>
    );
  }

  function renderDesktopView() {
    return (
      <div
        style={{
          display: 'flex',
          width: '100%',
          height: '100%'
        }}
      >
        {!canEdit && <PeopleList />}
        <UserInfo setShowSupport={setShowSupport} />

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
      </div>
    );
  }

  return (
    <Content>
      {isMobile ? renderMobileView() : renderDesktopView()}

      <Modal
        visible={showSupport}
        close={() => setShowSupport(false)}
        style={{
          height: '100%'
        }}
        envStyle={{
          marginTop: isMobile || canEdit ? 64 : 123,
          borderRadius: 0
        }}
      >
        <div
          dangerouslySetInnerHTML={{
            __html: `<sphinx-widget
                                pubkey=${owner_pubkey}
                                amount="500"
                                title="Support Me"
                                subtitle="Because I'm awesome"
                                buttonlabel="Donate"
                                defaultinterval="weekly"
                                imgurl="${
                                  mediumPic ||
                                  'https://i.scdn.co/image/28747994a80c78bc2824c2561d101db405926a37'
                                }"
                            ></sphinx-widget>`
          }}
        />
      </Modal>
    </Content>
  );
}

interface PanelProps {
  isMobile: boolean;
}

const Panel = styled.div<PanelProps>`
  position: relative;
  background: #ffffff;
  color: #000000;
  padding: 20px;
  box-shadow: ${(p) => (p.isMobile ? 'none' : '0px 0px 6px rgb(0 0 0 / 7%)')};
  border-bottom: ${(p) => (p.isMobile ? '2px solid #EBEDEF' : 'none')};
`;
const Content = styled.div`
  display: flex;
  flex-direction: column;

  width: 100%;
  height: 100%;
  align-items: center;
  color: #000000;
  background: #f0f1f3;
`;

const Counter = styled.div`
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 11px;
  line-height: 19px;
  margin-bottom: -3px;
  /* or 173% */
  margin-left: 8px;

  display: flex;
  align-items: center;

  /* Placeholder Text */

  color: #b0b7bc;
`;

const Tabs = styled.div`
  display: flex;
  width: 100%;
  align-items: center;
  // justify-content:center;
  overflow-x: auto;
  ::-webkit-scrollbar {
    display: none;
  }
`;

interface TagProps {
  selected: boolean;
}
const Tab = styled.div<TagProps>`
  display: flex;
  padding: 10px;
  margin-right: 25px;
  color: ${(p) => (p.selected ? '#292C33' : '#8E969C')};
  border-bottom: ${(p) => p.selected && '4px solid #618AFF'};
  cursor: hover;
  font-weight: 500;
  font-size: 15px;
  line-height: 19px;
  cursor: pointer;
`;

const Sleeve = styled.div``;
