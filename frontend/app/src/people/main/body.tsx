/* eslint-disable func-style */
import { EuiGlobalToastList, EuiLoadingSpinner } from '@elastic/eui';
import { useObserver } from 'mobx-react-lite';
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { useHistory } from 'react-router';
import { useLocation } from 'react-router-dom';
import styled from 'styled-components';
import FadeLeft from '../../animated/fadeLeft';
import { colors } from '../../colors';
import { useFuse, useIsMobile, usePageScroll, useScreenWidth } from '../../hooks';
import { Modal, SearchTextInput } from '../../sphinxUI';
import { useStores } from '../../store';
import Person from '../person';
import PersonViewSlim from '../personViewSlim';
import ConnectCard from '../utils/connectCard';
import { widgetConfigs } from '../utils/constants';
import NoResults from '../utils/noResults';
import PageLoadSpinner from '../utils/pageLoadSpinner';
import StartUpModal from '../utils/start_up_modal';
import BountyHeader from '../widgetViews/bountyHeader';
import WidgetSwitchViewer from '../widgetViews/widgetSwitchViewer';
import FirstTimeScreen from './firstTimeScreen';
import FocusedView from './focusView';
import { Widget } from './types';
// avoid hook within callback warning by renaming hooks

const getFuse = useFuse;
const getPageScroll = usePageScroll;

function useQuery() {
  const { search } = useLocation();

  return React.useMemo(() => new URLSearchParams(search), [search]);
}

export default function BodyComponent({ selectedWidget }: { selectedWidget: Widget }) {
  const { main, ui } = useStores();
  const [loading, setLoading] = useState(true);
  const [showDropdown, setShowDropdown] = useState(false);
  const screenWidth = useScreenWidth();
  const [publicFocusPerson, setPublicFocusPerson]: any = useState(null);
  const [publicFocusIndex, setPublicFocusIndex] = useState(-1);
  const [scrollValue, setScrollValue] = useState<boolean>(false);
  const [openStartUpModel, setOpenStartUpModel] = useState<boolean>(false);
  const closeModal = () => setOpenStartUpModel(false);
  const showModal = () => setOpenStartUpModel(true);
  const [openConnectModal, setConnectModal] = useState<boolean>(false);
  const closeConnectModal = () => setConnectModal(false);
  const showConnectModal = () => setConnectModal(true);
  const [connectPerson, setConnectPerson] = useState<any>();
  const [connectPersonBody, setConnectPersonBody] = useState<any>();
  const [activeListIndex, setActiveListIndex] = useState<number>(0);
  const [checkboxIdToSelectedMap, setCheckboxIdToSelectedMap] = useState({});
  const [checkboxIdToSelectedMapLanguage, setCheckboxIdToSelectedMapLanguage] = useState({});
  const [isModalSideButton, setIsModalSideButton] = useState<boolean>(true);
  const [isExtraStyle, setIsExtraStyle] = useState<boolean>(false);

  const color = colors['light'];

  const {
    peoplePageNumber,
    peopleWantedsPageNumber,
    peoplePostsPageNumber,
    peopleOffersPageNumber
  } = ui;

  const { peoplePosts, peopleWanteds, peopleOffers } = main;

  const listSource = {
    post: peoplePosts,
    wanted: peopleWanteds,
    offer: peopleOffers
  };

  const ReCallBounties = () => {
    (async () => {
      /*
       TODO : after getting the better way to reload the bounty, this code will be removed.
       */
      history.push('/tickets');
      await window.location.reload();
    })();
  };

  const activeList = listSource[selectedWidget];

  const history = useHistory();

  const isMobile = useIsMobile();
  const pathname = history?.location?.pathname;

  const loadMethods = useMemo(
    () => ({
      people: (a) => main.getPeople(a),
      post: (a) => main.getPeoplePosts(a),
      offer: (a) => main.getPeopleOffers(a),
      wanted: (a) => main.getPeopleWanteds(a)
    }),
    [main]
  );
  const doDeeplink = useCallback(async () => {
    if (pathname) {
      const splitPathname = pathname?.split('/');
      // eslint-disable-next-line prefer-destructuring
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
  }, [loadMethods, main, pathname, ui]);
  // deeplink page navigation
  useEffect(() => {
    doDeeplink();
  }, [doDeeplink]);

  const searchParams = useQuery();

  const publicPanelClick = useCallback(
    async (person, item) => {
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
        setConnectPerson({ ...person });
        setConnectPersonBody({ ...item });
      }
    },
    [selectedWidget]
  );

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
      setActiveListIndex(
        activeList && activeList.length
          ? activeList.findIndex((item) => {
              const { person, body } = item;
              return owner_id === person.owner_pubkey && created === `${body.created}`;
            })
          : {}
      );

      if (value.person && value.body) {
        publicPanelClick(value.person, value.body);
      }
    }
  }, [searchParams, main, activeList, publicPanelClick]);

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
  }, [loadMethods, selectedWidget, ui.selectedPerson]);

  useEffect(() => {
    main.getOpenGithubIssues();
    main.getBadgeList();
  }, [main]);

  useEffect(() => {
    if (ui.meInfo) {
      main.getTribesByOwner(ui.meInfo.owner_pubkey || '');
    }
  }, [main, ui.meInfo]);

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
  }, [ui.searchText, selectedWidget, loadMethods, ui.selectingPerson]);

  const onChangeStatus = (optionId) => {
    const newCheckboxIdToSelectedMap = {
      ...checkboxIdToSelectedMap,
      ...{
        [optionId]: !checkboxIdToSelectedMap[optionId]
      }
    };
    setCheckboxIdToSelectedMap(newCheckboxIdToSelectedMap);
  };

  const onChangeLanguage = (optionId) => {
    const newCheckboxIdToSelectedMapLanguage = {
      ...checkboxIdToSelectedMapLanguage,
      ...{
        [optionId]: !checkboxIdToSelectedMapLanguage[optionId]
      }
    };
    setCheckboxIdToSelectedMapLanguage(newCheckboxIdToSelectedMapLanguage);
  };

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
        p = [...p, <Spacer key={'spacer1'} />];
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
            alignItems: 'center',
            height: '100%'
          }}
        >
          <WidgetSwitchViewer
            checkboxIdToSelectedMap={checkboxIdToSelectedMap}
            checkboxIdToSelectedMapLanguage={checkboxIdToSelectedMapLanguage}
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
          checkboxIdToSelectedMap={checkboxIdToSelectedMap}
          checkboxIdToSelectedMapLanguage={checkboxIdToSelectedMapLanguage}
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
                scrollValue={scrollValue}
                onChangeStatus={onChangeStatus}
                onChangeLanguage={onChangeLanguage}
                checkboxIdToSelectedMap={checkboxIdToSelectedMap}
                checkboxIdToSelectedMapLanguage={checkboxIdToSelectedMapLanguage}
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
            scrollValue={scrollValue}
            onChangeStatus={onChangeStatus}
            onChangeLanguage={onChangeLanguage}
            checkboxIdToSelectedMap={checkboxIdToSelectedMap}
            checkboxIdToSelectedMapLanguage={checkboxIdToSelectedMapLanguage}
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
              borderRadius: isExtraStyle ? '10px' : 0,
              background: color.pureWhite,
              ...focusedDesktopModalStyles,
              maxHeight: '100vh',
              zIndex: 20
            }}
            style={{
              background: 'rgba( 0 0 0 /75% )'
            }}
            bigCloseImage={() => {
              setPublicFocusPerson(null);
              setPublicFocusIndex(-1);
              history.push('/tickets');
              setIsExtraStyle(false);
              setIsModalSideButton(true);
            }}
            bigCloseImageStyle={{
              top: isExtraStyle ? '-18px' : '18px',
              right: isExtraStyle ? '-18px' : '-50px',
              // background: isExtraStyle ? '#000' : '#000',
              borderRadius: isExtraStyle ? '50%' : '50%'
            }}
            prevArrowNew={
              activeListIndex === 0
                ? null
                : isModalSideButton
                ? () => {
                    const { person, body } = activeList[activeListIndex - 1];
                    if (person && body) {
                      history.replace({
                        pathname: history?.location?.pathname,
                        search: `?owner_id=${person?.owner_pubkey}&created=${body?.created}`,
                        state: {
                          owner_id: person?.owner_pubkey,
                          created: body?.created
                        }
                      });
                    }
                  }
                : null
            }
            nextArrowNew={
              activeListIndex + 1 > activeList?.length
                ? null
                : isModalSideButton
                ? () => {
                    const { person, body } = activeList[activeListIndex + 1];
                    if (person && body) {
                      history.replace({
                        pathname: history?.location?.pathname,
                        search: `?owner_id=${person?.owner_pubkey}&created=${body?.created}`,
                        state: {
                          owner_id: person?.owner_pubkey,
                          created: body?.created
                        }
                      });
                    }
                  }
                : null
            }
          >
            <FocusedView
              ReCallBounties={ReCallBounties}
              person={publicFocusPerson}
              personBody={connectPersonBody}
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
              fromBountyPage={true}
              extraModalFunction={() => {
                setPublicFocusPerson(null);
                setPublicFocusIndex(-1);
                history.push('/tickets');
                if (ui.meInfo) {
                  showConnectModal();
                } else {
                  showModal();
                }
              }}
              deleteExtraFunction={() => {
                setPublicFocusPerson(null);
                setPublicFocusIndex(-1);
              }}
              setIsModalSideButton={setIsModalSideButton}
              setIsExtraStyle={setIsExtraStyle}
            />
          </Modal>
        )}
        {openStartUpModel && (
          <StartUpModal closeModal={closeModal} dataObject={'getWork'} buttonColor={'primary'} />
        )}
        <ConnectCard
          dismiss={() => closeConnectModal()}
          modalStyle={{
            top: '-64px',
            height: 'calc(100% + 64px)'
          }}
          person={connectPerson}
          visible={openConnectModal}
        />
        {toastsEl}
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
