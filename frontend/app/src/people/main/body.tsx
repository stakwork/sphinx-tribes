import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { useObserver } from 'mobx-react-lite';
import { useStores } from '../../store';
import { EuiGlobalToastList, EuiLoadingSpinner } from '@elastic/eui';
import Person from '../person';
import PersonViewSlim from '../personViewSlim';
import { useFuse, usePageScroll, useIsMobile, useScreenWidth } from '../../hooks';
import FadeLeft from '../../animated/fadeLeft';
import FirstTimeScreen from './firstTimeScreen';
import { Divider, SearchTextInput, Modal } from '../../sphinxUI';
import WidgetSwitchViewer from '../widgetViews/widgetSwitchViewer';
import FocusedView from './focusView';
import { widgetConfigs } from '../utils/constants';
import { useHistory } from 'react-router';
import { useLocation } from 'react-router-dom';
import NoResults from '../utils/noResults';
import PageLoadSpinner from '../utils/pageLoadSpinner';
import BountyHeader from '../widgetViews/bountyHeader';
import { colors } from '../../colors';
// import { SearchTextInput } from '../../sphinxUI/index'
// avoid hook within callback warning by renaming hooks

const getFuse = useFuse;
const getPageScroll = usePageScroll;

function useQuery() {
  const { search } = useLocation();

  return React.useMemo(() => new URLSearchParams(search), [search]);
}

export default function BodyComponent({ selectedWidget }) {
  const { main, ui } = useStores();
  const [loading, setLoading] = useState(true);
  const [showDropdown, setShowDropdown] = useState(false);
  const screenWidth = useScreenWidth();
  const [publicFocusPerson, setPublicFocusPerson]: any = useState(null);
  const [publicFocusIndex, setPublicFocusIndex] = useState(-1);
  const [showFocusView, setShowFocusView] = useState(false);
  const [focusIndex, setFocusIndex] = useState(-1);
  const [isMobileViewTicketModal, setIsMobileViewTicketModal] = useState(false);
  const [scrollValue, setScrollValue] = useState<boolean>(false);

  const color = colors['light'];

  const {
    peoplePageNumber,
    peopleWantedsPageNumber,
    peoplePostsPageNumber,
    peopleOffersPageNumber,
    openGithubIssues
  } = ui;

  const { peoplePosts, peopleWanteds, peopleOffers } = main;

  const listSource = {
    post: peoplePosts,
    wanted: peopleWanteds,
    offer: peopleOffers
  };

  const person: any =
    main.people && main.people.length && main.people.find((f) => f.id === ui.selectedPerson);

  const { id } = person || {};

  const canEdit = id === ui.meInfo?.id;

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

  const activeList = listSource[selectedWidget];

  const history = useHistory();

  interface ITabAction {
    text?: string;
    icon?: string;
  }
  interface ITab {
    label: string;
    name: string;
    description: string;
    action?: ITabAction;
  }

  const tabs: Array<ITab> = [
    {
      label: 'People',
      name: 'people',
      description: 'Find fellow Sphinx Chatters',
      action: {
        text: 'Add New Ticket',
        icon: 'group'
      }
    },
    // // widgetConfigs['post'],
    // {
    //   ...widgetConfigs['offer'],
    //   label: 'Portfolios',
    // },
    {
      ...widgetConfigs['wanted'],
      description: 'Earn sats for completing tickets'
    }
  ];

  const tabsModal = widgetConfigs;

  const isMobile = useIsMobile();
  const pathname = history?.location?.pathname;

  const loadMethods = {
    people: (a) => main.getPeople(a),
    post: (a) => main.getPeoplePosts(a),
    offer: (a) => main.getPeopleOffers(a),
    wanted: (a) => main.getPeopleWanteds(a)
  };

  // deeplink page navigation
  useEffect(() => {
    doDeeplink();
  }, []);

  const searchParams = useQuery();

  async function publicPanelClick(person, item) {
    // migrating to load widgets separate from person
    console.log('person', { person }, 'and items', { item });
    const itemIndex = person[selectedWidget]?.findIndex((f) => f.created === item.created);
    if (itemIndex > -1) {
      // make person into proper structure (derived from widget)
      const p = {
        ...person,
        extras: {
          [selectedWidget]: person[selectedWidget]
        }
      };
      setPublicFocusPerson(p);
      setPublicFocusIndex(itemIndex);
    }
  }

  useEffect(() => {
    const owner_id = searchParams.get('owner_id');
    const created = searchParams.get('created');
    if (owner_id && created) {
      const value =
        activeList && activeList.length
          ? activeList.find((item) => {
              const { person, body } = item;
              return owner_id === person.owner_pubkey && created === `${body.created}`;
            })
          : {};
      if (value.person && value.body) {
        publicPanelClick(value.person, value.body);
      }
    }
  }, [searchParams, main, activeList]);

  useEffect(() => {
    // clear public focus is selected person
    if (ui.selectedPerson) {
      setPublicFocusPerson(null);
      setPublicFocusIndex(-1);
    } else {
      //pull list again, we came back from focus view
      const loadMethod = loadMethods[selectedWidget];
      loadMethod({ page: 1, resetPage: true });
    }
  }, [ui.selectedPerson]);

  useEffect(() => {
    // loadPeople()
    // loadPeopleExtras()
    main.getOpenGithubIssues();
    main.getBadgeList();
  }, []);

  useEffect(() => {
    if (ui.meInfo) {
      main.getTribesByOwner(ui.meInfo.owner_pubkey || '');
    }
  }, [ui.meInfo]);

  // do search update
  useEffect(() => {
    (async () => {
      // selectedWidget
      // get assets page 1, by widget
      console.log('refresh list for search');
      let loadMethod = loadMethods[selectedWidget];

      // if person is selected, always searching people
      if (ui.selectingPerson) {
        loadMethod = loadMethods['people'];
      }

      // reset page will replace all results, this is good for a new search!
      await loadMethod({ page: 1, resetPage: true });

      setLoading(false);
    })();
  }, [ui.searchText, selectedWidget]);

  async function doDeeplink() {
    if (pathname) {
      const splitPathname = pathname?.split('/');
      const personPubkey: string = splitPathname[2];
      if (personPubkey) {
        const p = await main.getPersonByPubkey(personPubkey);
        ui.setSelectedPerson(p?.id);
        ui.setSelectingPerson(p?.id);
        // make sure to load people in a person deeplink
        const loadMethod = loadMethods['people'];
        await loadMethod({ page: 1, resetPage: true });
      }
    }
  }

  function selectPerson(id: number, unique_name: string, pubkey: string) {
    console.log('selectPerson', id, unique_name, pubkey);
    ui.setSelectedPerson(id);
    ui.setSelectingPerson(id);

    history.push(`/p/${pubkey}`);
  }

  function goBack() {
    ui.setSelectingPerson(0);
    history.push('/tickets');
  }

  return useObserver(() => {
    let people = getFuse(main.people, ['owner_alias']);

    const loadForwardFunc = () => loadMore(1);
    const loadBackwardFunc = () => loadMore(-1);
    const { loadingTop, loadingBottom, handleScroll } = getPageScroll(
      loadForwardFunc,
      loadBackwardFunc
    );

    people = (people && people.filter((f) => !f.hide)) || [];

    async function loadMore(direction) {
      let currentPage = 1;

      switch (selectedWidget) {
        case 'people':
          currentPage = peoplePageNumber;
          break;
        case 'wanted':
          currentPage = peopleWantedsPageNumber;
          break;
        case 'offer':
          currentPage = peopleOffersPageNumber;
          break;
        case 'post':
          currentPage = peoplePostsPageNumber;
          break;
        default:
          console.log('scroll', direction);
      }

      let newPage = currentPage + direction;
      if (newPage < 1) newPage = 1;

      try {
        switch (selectedWidget) {
          case 'people':
            await main.getPeople({ page: newPage });
            break;
          case 'wanted':
            await main.getPeopleWanteds({ page: newPage });
            break;
          case 'offer':
            await main.getPeopleOffers({ page: newPage });
            break;
          case 'post':
            await main.getPeoplePosts({ page: newPage });
            break;
          default:
            console.log('scroll', direction);
        }
      } catch (e) {
        console.log('load failed', e);
      }
    }

    function renderPeople() {
      let p;
      if (people.length) {
        p = people?.map((t) => (
          <Person
            {...t}
            key={t.id}
            small={isMobile}
            squeeze={screenWidth < 1420}
            selected={ui.selectedPerson === t.id}
            select={selectPerson}
          />
        ));

        // add space at bottom
        p = [...p, <Spacer key={'spacer'} />];
      } else {
        p = <NoResults />;
      }

      return p;
    }

    const listContent =
      selectedWidget === 'people' ? (
        renderPeople()
      ) : !isMobile ? (
        <div
          style={{
            width: '100%',
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center'
          }}
        >
          <WidgetSwitchViewer
            onPanelClick={(person, item) => {
              history.replace({
                pathname: history?.location?.pathname,
                search: `?owner_id=${person.owner_pubkey}&created=${item.created}`,
                state: {
                  owner_id: person.owner_pubkey,
                  created: item.created
                }
              });
              publicPanelClick(person, item);
            }}
            fromBountyPage={true}
            selectedWidget={selectedWidget}
            loading={loading}
          />
        </div>
      ) : (
        <WidgetSwitchViewer
          onPanelClick={(person, item) => {
            history.replace({
              pathname: history?.location?.pathname,
              search: `?owner_id=${person.owner_pubkey}&created=${item.created}`,
              state: {
                owner_id: person.owner_pubkey,
                created: item.created
              }
            });
            publicPanelClick(person, item);
          }}
          fromBountyPage={true}
          selectedWidget={selectedWidget}
          loading={loading}
        />
      );

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

    const widgetLabel = selectedWidget && tabs.find((f) => f.name === selectedWidget);

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
            {selectedWidget === 'wanted' && (
              <BountyHeader
                selectedWidget={selectedWidget}
                setShowFocusView={setIsMobileViewTicketModal}
                scrollValue={scrollValue}
              />
            )}
            {selectedWidget === 'people' && (
              <div
                style={{
                  padding: '0 20px'
                }}
              >
                <SearchTextInput
                  small
                  name="search"
                  type="search"
                  placeholder="Search"
                  value={ui.searchText}
                  style={{
                    width: '100%',
                    height: 40,
                    border: `1px solid ${color.grayish.G600}`,
                    background: color.pureWhite
                  }}
                  onChange={(e) => {
                    console.log('handleChange', e);
                    ui.setSearchText(e);
                  }}
                />
              </div>
            )}
          </div>

          {showDropdown && <Backdrop onClick={() => setShowDropdown(false)} />}
          <div style={{ width: '100%' }}>
            <PageLoadSpinner show={loadingTop} />
            {listContent}
            <PageLoadSpinner noAnimate show={loadingBottom} />
          </div>

          <FadeLeft
            withOverlay
            drift={40}
            overlayClick={() => goBack()}
            style={{
              position: 'absolute',
              top: 0,
              right: 0,
              zIndex: 10000,
              width: '100%'
            }}
            isMounted={ui.selectingPerson ? true : false}
            dismountCallback={() => ui.setSelectedPerson(0)}
          >
            <PersonViewSlim
              goBack={goBack}
              personId={ui.selectedPerson}
              selectPerson={selectPerson}
              loading={loading}
            />
          </FadeLeft>
          {publicFocusPerson && (
            <Modal visible={publicFocusPerson ? true : false} fill={true}>
              <FocusedView
                person={publicFocusPerson}
                canEdit={false}
                selectedIndex={publicFocusIndex}
                config={widgetConfigs[selectedWidget] && widgetConfigs[selectedWidget]}
                onSuccess={() => {
                  console.log('success');
                  setPublicFocusPerson(null);
                  setPublicFocusIndex(-1);
                }}
                goBack={() => {
                  setPublicFocusPerson(null);
                  setPublicFocusIndex(-1);
                  history.push('/tickets');
                }}
              />
            </Modal>
          )}

          {isMobileViewTicketModal && (
            <Modal visible={isMobileViewTicketModal} fill={true}>
              <FocusedView
                person={person}
                canEdit={!canEdit}
                selectedIndex={focusIndex}
                config={tabsModal[selectedWidget] && tabsModal[selectedWidget]}
                onSuccess={() => {
                  console.log('success');
                  setFocusIndex(-1);
                  // if (selectedWidget === 'about') switchWidgets('badges');
                }}
                goBack={() => {
                  setIsMobileViewTicketModal(false);
                  setFocusIndex(-1);
                  history.push('/tickets');
                  // if (selectedWidget === 'about') switchWidgets('badges');
                }}
              />
            </Modal>
          )}

          {toastsEl}
        </Body>
      );
    }

    const focusedDesktopModalStyles =
      selectedWidget && widgetConfigs[selectedWidget]
        ? {
            ...widgetConfigs[selectedWidget].modalStyle
          }
        : {};

    // desktop mode
    return (
      <Body
        onScroll={(e) => {
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
        {selectedWidget === 'wanted' && (
          <BountyHeader
            selectedWidget={selectedWidget}
            setShowFocusView={setShowFocusView}
            scrollValue={scrollValue}
          />
        )}
        {selectedWidget === 'people' && (
          <div
            style={{
              display: 'flex',
              justifyContent: 'flex-end',
              padding: '10px 0'
            }}
          >
            <SearchTextInput
              small
              name="search"
              type="search"
              placeholder="Search"
              value={ui.searchText}
              style={{
                width: 204,
                height: 40,
                border: `1px solid ${color.grayish.G600}`,
                background: color.grayish.G600
              }}
              onChange={(e) => {
                console.log('handleChange', e);
                ui.setSearchText(e);
              }}
            />
          </div>
        )}
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
            {listContent}
            <PageLoadSpinner noAnimate show={loadingBottom} />
          </div>
        </>
        {/* selected view */}
        <FadeLeft
          withOverlay={isMobile}
          drift={40}
          overlayClick={() => goBack()}
          style={{
            position: 'absolute',
            top: isMobile ? 0 : 64,
            right: 0,
            zIndex: 10000,
            width: '100%'
          }}
          isMounted={ui.selectingPerson ? true : false}
          dismountCallback={() => ui.setSelectedPerson(0)}
        >
          <PersonViewSlim
            goBack={goBack}
            personId={ui.selectedPerson}
            loading={loading}
            peopleView={true}
            selectPerson={selectPerson}
          />
        </FadeLeft>
        {/* modal onClick on tickets */}
        {publicFocusPerson && (
          <Modal
            visible={publicFocusPerson ? true : false}
            envStyle={{
              borderRadius: 0,
              background: color.pureWhite,
              height: '100%',
              width: '60%',
              minWidth: 500,
              maxWidth: 602,
              ...focusedDesktopModalStyles
            }}
            bigClose={() => {
              setPublicFocusPerson(null);
              setPublicFocusIndex(-1);
              history.push('/tickets');
            }}
          >
            <FocusedView
              person={publicFocusPerson}
              canEdit={false}
              selectedIndex={publicFocusIndex}
              config={widgetConfigs[selectedWidget] && widgetConfigs[selectedWidget]}
              onSuccess={() => {
                console.log('success');
                setPublicFocusPerson(null);
                setPublicFocusIndex(-1);
              }}
              goBack={() => {
                setPublicFocusPerson(null);
                setPublicFocusIndex(-1);
              }}
            />
          </Modal>
        )}
        {toastsEl}
        {/* modal create ticket */}
        {showFocusView && (
          <Modal
            visible={showFocusView}
            style={{
              // top: -64,
              // height: 'calc(100% + 64px)'
              height: '100%'
            }}
            envStyle={{
              marginTop: isMobile ? 64 : 0,
              borderRadius: 0,
              background: color.pureWhite,
              height: '100%',
              width: '60%',
              minWidth: 500,
              maxWidth: 602,
              zIndex: 20,
              ...focusedDesktopModalStyles
            }}
            nextArrow={nextIndex}
            prevArrow={prevIndex}
            overlayClick={() => {
              setShowFocusView(false);
              setFocusIndex(-1);
              // if (selectedWidget === 'about') switchWidgets('badges');
            }}
            bigClose={() => {
              setShowFocusView(false);
              setFocusIndex(-1);
              // if (selectedWidget === 'about') switchWidgets('badges');
            }}
          >
            <FocusedView
              person={person}
              canEdit={!canEdit}
              selectedIndex={focusIndex}
              config={tabsModal[selectedWidget] && tabsModal[selectedWidget]}
              onSuccess={() => {
                console.log('success');
                setFocusIndex(-1);
                // if (selectedWidget === 'about') switchWidgets('badges');
              }}
              goBack={() => {
                setShowFocusView(false);
                setFocusIndex(-1);
                // if (selectedWidget === 'about') switchWidgets('badges');
              }}
            />
          </Modal>
        )}
      </Body>
    );
  });
}

const Body = styled.div`
  flex: 1;
  height: calc(100% - 105px);
  // padding-bottom:80px;
  width: 100%;
  overflow: auto;
  display: flex;
  flex-direction: column;
`;
const Label = styled.div`
  font-family: Roboto;
  font-style: normal;
  font-weight: bold;
  font-size: 26px;
  line-height: 40px;
  /* or 154% */
  width: 204px;

  display: flex;
  align-items: center;

  /* Text 2 */

  color: #3c3f41;
`;

const Tabs = styled.div`
  display: flex;
`;

interface TagProps {
  selected: boolean;
}
const Tab = styled.div<TagProps>`
  display: flex;
  padding: 10px 25px;
  margin-right: 35px;
  height: 42px;
  color: ${(p) => (p.selected ? '#5D8FDD' : '#5F6368')};
  border: 2px solid #5f636800;
  border-color: ${(p) => (p.selected ? '#CDE0FF' : '#5F636800')};
  // border-bottom: ${(p) => p.selected && '4px solid #618AFF'};
  cursor: pointer;
  font-weight: 400;
  font-size: 15px;
  line-height: 19px;
  background: ${(p) => (p.selected ? '#DCEDFE' : '#3C3F4100')};
  border-radius: 25px;
`;

const TabMobile = styled(Tab)`
  margin: 0;
  border-radius: 0;
  height: auto;
  background: #fff;
  border: none;
  padding: 20px 12px 0 20px;
  display: flex;
  align-items: top;

  .tab-icon {
    margin-right: 20px;
    color: ${(p) => p.selected && '#A2C0FD'};
  }

  .tab-details {
    display: flex;
    flex: 1;
    flex-direction: column;
    border-bottom: 1px solid #eee;
    padding-bottom: 20px;
    &__title {
      color: ${(p) => p.selected && '#618AFF'};
      font-size: 16px;
      font-weight: 600;
    }
    &__subtitle {
      color: ${(p) => p.selected && '#A2C0FD'};
      font-size: 12px;
    }
  }
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

const Link = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  margin-left: 6px;
  color: #618aff;
  cursor: pointer;
`;

const Loader = styled.div`
  position: absolute;
  width: 100%;
  display: flex;
  justify-content: center;
  padding: 10px;
  left: 0px;
  z-index: 20;
`;

export const Spacer = styled.div`
  display: flex;
  min-height: 10px;
  min-width: 100%;
  height: 10px;
  width: 100%;
`;

const MobileDropdown = styled.div`
  position: absolute;
  top: calc(100% + 1px);
  left: 0;
  right: 0;
  z-index: 10;
  background: #fff;
`;
