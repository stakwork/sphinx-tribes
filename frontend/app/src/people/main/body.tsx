import React, { useEffect, useState, useRef } from 'react';
import styled from 'styled-components';
import { useObserver } from 'mobx-react-lite';
import { useStores } from '../../store';
import { EuiGlobalToastList, EuiLoadingSpinner } from '@elastic/eui';
import Person from '../person';
import PersonViewSlim from '../personViewSlim';
import {
  useFuse,
  usePageScroll,
  useIsMobile,
  useScreenWidth,
} from '../../hooks';
import FadeLeft from '../../animated/fadeLeft';
import FirstTimeScreen from './firstTimeScreen';
import NoneSpace from '../utils/noneSpace';
import { Divider, SearchTextInput, Modal, Button } from '../../sphinxUI';
import WidgetSwitchViewer from '../widgetViews/widgetSwitchViewer';
import MaterialIcon from '@material/react-material-icon';
import FocusedView from './focusView';

import { widgetConfigs } from '../utils/constants';
import { useHistory } from 'react-router';
import { useLocation, useParams } from 'react-router-dom';
import NoResults from '../utils/noResults';
import PageLoadSpinner from '../utils/pageLoadSpinner';
// import { SearchTextInput } from '../../sphinxUI/index'
// avoid hook within callback warning by renaming hooks

const getFuse = useFuse;
const getPageScroll = usePageScroll;

function useQuery() {
  const { search } = useLocation();

  return React.useMemo(() => new URLSearchParams(search), [search]);
}

export default function BodyComponent() {
  const { main, ui } = useStores();
  const [loading, setLoading] = useState(true);
  const [showDropdown, setShowDropdown] = useState(false);
  const screenWidth = useScreenWidth();
  const [publicFocusPerson, setPublicFocusPerson]: any = useState(null);
  const [publicFocusIndex, setPublicFocusIndex] = useState(-1);
  const [showFocusView, setShowFocusView] = useState(false);
  const [focusIndex, setFocusIndex] = useState(-1);
  const [isMobileViewTicketModal, setIsMobileViewTicketModal] = useState(false);

  const {
    peoplePageNumber,
    peopleWantedsPageNumber,
    peoplePostsPageNumber,
    peopleOffersPageNumber,
    openGithubIssues,
  } = ui;

  const [selectedWidget, setSelectedWidget] = useState('wanted');
  const { peoplePosts, peopleWanteds, peopleOffers } = main;

  const listSource = {
    post: peoplePosts,
    wanted: peopleWanteds,
    offer: peopleOffers,
  };

  let person: any =
    main.people &&
    main.people.length &&
    main.people.find((f) => f.id === ui.selectedPerson);

  const { id } = person || {};

  const canEdit = id === ui.meInfo?.id;

  function nextIndex() {
    if (focusIndex < 0) {
      console.log('nope!');
      return;
    }
    if (person && person.extras) {
      let g = person.extras[tabs[selectedWidget]?.name];
      let nextindex = focusIndex + 1;
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
      let g = person?.extras[tabs[selectedWidget]?.name];
      let previndex = focusIndex - 1;
      if (g[previndex]) setFocusIndex(previndex);
      else setFocusIndex(g.length - 1);
    }
  }

  const activeList = listSource[selectedWidget];

  const history = useHistory();

  const tabs = [
    {
      label: 'People',
      name: 'people',
    },
    // widgetConfigs['post'],
    {
      ...widgetConfigs['offer'],
      label: 'Portfolios',
    },
    widgetConfigs['wanted'],
  ];

  const tabsModal = widgetConfigs;

  const isMobile = useIsMobile();
  const pathname = history?.location?.pathname;

  const loadMethods = {
    people: (a) => main.getPeople(a),
    post: (a) => main.getPeoplePosts(a),
    offer: (a) => main.getPeopleOffers(a),
    wanted: (a) => main.getPeopleWanteds(a),
  };

  // deeplink page navigation
  useEffect(() => {
    doDeeplink();
  }, []);

  const searchParams = useQuery();

  async function publicPanelClick(person, item) {
    // migrating to load widgets separate from person
    console.log('person', { person }, 'and items', { item });
    const itemIndex = person[selectedWidget]?.findIndex(
      (f) => f.created === item.created
    );
    if (itemIndex > -1) {
      // make person into proper structure (derived from widget)
      let p = {
        ...person,
        extras: {
          [selectedWidget]: person[selectedWidget],
        },
      };
      //   console.log(p, itemIndex);
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
              return (
                owner_id === person.owner_pubkey &&
                created === body.created + ''
              );
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
      console.log('pull list');
      let loadMethod = loadMethods[selectedWidget];
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
      let splitPathname = pathname?.split('/');
      let personPubkey: string = splitPathname[2];
      if (personPubkey) {
        let p = await main.getPersonByPubkey(personPubkey);
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
    history.push('/p');
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
      ) : (
        <WidgetSwitchViewer
          onPanelClick={(person, item) => {
            history.replace({
              pathname: history?.location?.pathname,
              search: `?owner_id=${person.owner_pubkey}&created=${item.created}`,
              state: {
                owner_id: person.owner_pubkey,
                created: item.created,
              },
            });
            publicPanelClick(person, item);
          }}
          selectedWidget={selectedWidget}
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

    const widgetLabel =
      selectedWidget && tabs.find((f) => f.name === selectedWidget);

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
          {!ui.meInfo && (
            <div style={{ marginTop: 60 }}>
              <NoneSpace
                buttonText={'Get Started'}
                buttonIcon={'arrow_forward'}
                action={() => ui.setShowSignIn(true)}
                img={'explore.png'}
                text={'Start your own profile'}
                style={{ height: 320, background: '#fff' }}
              />
              <Divider />
            </div>
          )}

          <div
            style={{
              width: '100%',
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'flex-start',
              padding: 20,
              height: 82,
              boxShadow: '0 0 6px 0 rgba(0, 0, 0, 0.07)',
              zIndex: 2,
            }}
          >
            <Label style={{ fontSize: 20 }}>
              Explore
              <Link onClick={() => setShowDropdown(true)}>
                <div>{widgetLabel && widgetLabel.label}</div>
                <MaterialIcon
                  icon={'expand_more'}
                  style={{ fontSize: 18, marginLeft: 5 }}
                />

                {showDropdown && (
                  <div
                    style={{
                      position: 'absolute',
                      top: 0,
                      left: 0,
                      zIndex: 10,
                      background: '#fff',
                    }}
                  >
                    {tabs &&
                      tabs.map((t, i) => {
                        const label = t.label;
                        const selected = selectedWidget === t.name;

                        return (
                          <Tab
                            key={i}
                            style={{ borderRadius: 0, margin: 0 }}
                            selected={selected}
                            onClick={(e) => {
                              e.stopPropagation();
                              setShowDropdown(false);
                              setSelectedWidget(t.name);
                            }}
                          >
                            {label}
                          </Tab>
                        );
                      })}
                  </div>
                )}
              </Link>
            </Label>

							<div style={{display: 'flex'}} >
            {selectedWidget === 'wanted' &&
              ui.meInfo &&
              ui.meInfo?.owner_alias && (
                <>
                  <div
                    style={{
                      fontSize: '15px',
                      fontWeight: '400',
                      marginRight: '10px',
                      cursor: 'pointer',
                      borderRadius: '20px',
                      userSelect: 'none',
                      background: '#dcedfe',
                      border: '2px solid #cddffd',
                      padding: '8px 20px',
                      color: '#5d92df',
                    }}
                    onClick={() => {
                      // setShowFocusView(true);
                      setIsMobileViewTicketModal(true);
                      console.log('hi');
                    }}
                  >
                    +
                  </div>
											</>
              )}

            <SearchTextInput
              small
              name="search"
              type="search"
              placeholder="Search"
              value={ui.searchText}
              style={{
                width: 164,
                height: 40,
                border: '1px solid #DDE1E5',
                background: '#fff',
              }}
              onChange={(e) => {
                console.log('handleChange', e);
                ui.setSearchText(e);
              }}
            />
							</div>
          </div>

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
              width: '100%',
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
                config={
                  widgetConfigs[selectedWidget] && widgetConfigs[selectedWidget]
                }
                onSuccess={() => {
                  console.log('success');
                  setPublicFocusPerson(null);
                  setPublicFocusIndex(-1);
                }}
                goBack={() => {
                  setPublicFocusPerson(null);
                  setPublicFocusIndex(-1);
                  history.push('/p');
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
                  history.push('/p');
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
            ...widgetConfigs[selectedWidget].modalStyle,
          }
        : {};

    // desktop mode
    return (
      <Body
        onScroll={handleScroll}
        style={{
          background: '#f0f1f3',
          height: 'calc(100% - 65px)',
        }}
      >
        {!ui.meInfo && (
          <div>
            <NoneSpace
              banner
              buttonText={'Get Started'}
              buttonIcon={'arrow_forward'}
              action={() => ui.setShowSignIn(true)}
              img={'explore.png'}
              text={'Start your own profile'}
              style={{ height: 320 }}
            />
            <Divider />
          </div>
        )}
        <div
          style={{
            width: '100%',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'flex-start',
            padding: 20,
            height: 62,
          }}
        >
          <Label>Explore</Label>

          <Tabs>
            {tabs &&
              tabs.map((t, i) => {
                const label = t.label;
                const selected = selectedWidget === t.name;
                const isWanted = 'wanted' === t.name;

                return (
                  <Tab
                    key={i}
                    selected={selected}
                    onClick={() => {
                      setSelectedWidget(t.name);
                    }}
                  >
                    {label}
                  </Tab>
                );
              })}
          </Tabs>

          <div
            style={{
              display: 'flex',
              flexDirection: 'row',
              alignItems: 'center',
              justifyContent: 'space-between',
            }}
          >
            {selectedWidget === 'wanted' &&
              (ui.meInfo && ui.meInfo?.owner_alias ? (

                  <div
                    style={{
                      fontSize: '15px',
                      fontWeight: '400',
                      marginRight: '10px',
                      cursor: 'pointer',
                      borderRadius: '20px',
                      userSelect: 'none',
                      background: '#dcedfe',
                      border: '2px solid #cddffd',
                      padding: '8px 20px',
                      color: '#5d92df',
                    }}
                    onClick={() => {
                      setShowFocusView(true);
                    }}
                  >
                    Create Ticket
                  </div>
              ) : (
                  <div
                    style={{
                      padding: '10px 20px',
                      borderRadius: '20px',
                      userSelect: 'none',
                      cursor: 'not-allowed',
                      color: '#83737d',
                      backgroundColor: '#dde0e5',
                      fontSize: '14px',
                      marginRight: '10px',
                    }}
                  >
                    Login to Create Tickets
                  </div>
              ))}

            <SearchTextInput
              name="search"
              type="search"
              placeholder="Search"
              value={ui.searchText}
              style={{ width: 204, height: 40, background: '#DDE1E5' }}
              onChange={(e) => {
                console.log('handleChange', e);
                ui.setSearchText(e);
              }}
            />
          </div>
        </div>
        <>
          <div
            style={{
              width: '100%',
              display: 'flex',
              flexWrap: 'wrap',
              height: '100%',
              justifyContent: 'flex-start',
              alignItems: 'flex-start',
              padding: 20,
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
            width: '100%',
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

        {publicFocusPerson && (
          <Modal
            visible={publicFocusPerson ? true : false}
            envStyle={{
              borderRadius: 0,
              background: '#fff',
              height: '100%',
              width: '60%',
              minWidth: 500,
              maxWidth: 602,
              ...focusedDesktopModalStyles,
            }}
            bigClose={() => {
              setPublicFocusPerson(null);
              setPublicFocusIndex(-1);
              history.push('/p');
            }}
          >
            <FocusedView
              person={publicFocusPerson}
              canEdit={false}
              selectedIndex={publicFocusIndex}
              config={
                widgetConfigs[selectedWidget] && widgetConfigs[selectedWidget]
              }
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

        {showFocusView && (
          <Modal
            visible={showFocusView}
            style={{
              // top: -64,
              // height: 'calc(100% + 64px)'
              height: '100%',
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
              ...focusedDesktopModalStyles,
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

const Link = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  margin-left: 6px;
  color: #618aff;
  cursor: pointer;
  position: relative;
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
