/* eslint-disable func-style */
import React, { useCallback, useEffect, useState } from 'react';
import {
  Img,
  RowWrap,
  Sleeve,
  Name,
  Head,
  Tab,
  Tabs,
  Counter,
  Content,
  Panel,
  DBack,
  AboutWrap,
  PeopleScroller,
  PeopleList
} from './personSlim/style';
import { getHost } from '../config/host';
import { useStores } from '../store';

import AboutView from './widgetViews/aboutView';
import BlogView from './widgetViews/blogView';
import OfferView from './widgetViews/offerView';
import PostView from './widgetViews/postView';
import SupportMeView from './widgetViews/supportMeView';
import TwitterView from './widgetViews/twitterView';
import WantedView from './widgetViews/wantedView';

import { useHistory, useLocation } from 'react-router';
import { meSchema } from '../components/form/schema';
import { useIsMobile, usePageScroll } from '../hooks';
import { Button, IconButton, Modal, SearchTextInput } from '../components/common';
import { queryLimit } from '../store/main';
import FocusedView from './main/focusView';
import Person from './person';
import Badges from './utils/badges';
import ConnectCard from './utils/connectCard';
import { widgetConfigs } from './utils/constants';
import NoResults from './utils/noResults';
import NoneSpace from './utils/noneSpace';
import PageLoadSpinner from './utils/pageLoadSpinner';
import { PostBounty } from './widgetViews/postBounty';
import { Widget } from './main/types';
import { LoaderBottom, LoaderTop } from './personSlim/component';
import MobileView from './personSlim/mobileView';

const host = getHost();
function makeQR(pubkey: string) {
  return `sphinx.chat://?action=person&host=${host}&pubkey=${pubkey}`;
}

export default function PersonView(props: any) {
  const { personId, loading, selectPerson, goBack } = props;
  // on this screen, there will always be a pubkey in the url, no need for personId

  const { main, ui } = useStores();
  const { meInfo, peoplePageNumber } = ui || {};

  const [loadingPerson, setLoadingPerson]: any = useState(false);
  const [loadedPerson, setLoadedPerson]: any = useState(null);

  const history = useHistory();
  const location = useLocation();
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

  const { id, img, owner_alias, extras, owner_pubkey } =
    person || {};

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

  const [showQR, setShowQR] = useState(false);
  const [showFocusView, setShowFocusView] = useState(false);
  const qrString = makeQR(owner_pubkey || '');

  async function loadMorePeople(direction) {
    let newPage = peoplePageNumber + direction;
    if (newPage < 1) {
      newPage = 1;
    }
    await main.getPeople({ page: newPage });
  }

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

  function selectPersonWithinFocusView(id, unique_name, pubkey) {
    setShowFocusView(false);
    setFocusIndex(-1);
    selectPerson(id, unique_name, pubkey);
  }

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

  const { loadingTop, loadingBottom, handleScroll } = usePageScroll(
    () => loadMorePeople(1),
    () => loadMorePeople(-1)
  );

  if (loading) return <div>Loading...</div>;

  let widgetSchemas: any = meSchema.find((f) => f.name === 'extras');
  if (widgetSchemas && widgetSchemas.extras) {
    widgetSchemas = widgetSchemas && widgetSchemas.extras;
  }

  const fullSelectedWidget: any = extras && selectedWidget ? extras[selectedWidget] : null;

  console.log(fullSelectedWidget);

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

  function nextIndex() {
    if (focusIndex < 0) {
      console.log('nope!');
      return;
    }
    if (person && person.extras) {
      const g = person.extras[tabs[selectedWidget]?.name];
      const nextindex = focusIndex + 1;
      if (g[nextindex]) setFocusIndex(nextindex);
      else setFocusIndex(0);
    }
  }

  function prevIndex() {
    if (focusIndex < 0) {
      console.log('nope!');
      return;
    }
    if (person && person.extras) {
      const g = person?.extras[tabs[selectedWidget]?.name];
      const previndex = focusIndex - 1;
      if (g[previndex]) setFocusIndex(previndex);
      else setFocusIndex(g.length - 1);
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

  const defaultPic = '/static/person_placeholder.png';
  const mediumPic = img;

  function renderDesktopView() {
    const focusedDesktopModalStyles = newSelectedWidget
      ? {
        ...tabs[newSelectedWidget]?.modalStyle
      }
      : {};

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
                  console.log('handleChange', e);
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
              console.log('success');
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

  return (
    <Content>
      {isMobile ? <MobileView
        logout
        person
        canEdit
        isMobile
        mediumPic
        defaultPic
        tabs
        showFocusView
        focusIndex
        setFocusIndex
        setShowFocusView
        selectedWidget
        extras
        goBack
        switchWidgets
        qrString
        owner_alias
        setShowSupport
        renderEditButton
        newSelectedWidget
        renderWidgets
      /> : renderDesktopView()}

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
                                imgurl="${mediumPic ||
              'https://i.scdn.co/image/28747994a80c78bc2824c2561d101db405926a37'
              }"
                            ></sphinx-widget>`
          }}
        />
      </Modal>
    </Content>
  );
}