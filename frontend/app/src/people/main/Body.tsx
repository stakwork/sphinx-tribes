import { observer } from 'mobx-react-lite';
import React, { useEffect, useState } from 'react';
import { useHistory } from 'react-router';
import styled from 'styled-components';
import { EuiLoadingSpinner, EuiGlobalToastList } from '@elastic/eui';
import PeopleHeader from 'people/widgetViews/PeopleHeader';
import { Person as PersonType } from 'store/main';
import filterByCodingLanguage from 'people/utils/filterPeople';
import { SearchTextInput } from '../../components/common';
import { colors } from '../../config/colors';
import { useFuse, useIsMobile, usePageScroll, useScreenWidth } from '../../hooks';
import { useStores } from '../../store';
import Person from '../../pages/people/Person';
import NoResults from '../utils/NoResults';
import PageLoadSpinner from '../utils/PageLoadSpinner';
import StartUpModal from '../utils/StartUpModal';

const color = colors['light'];
const Body = styled.div<{ isMobile: boolean }>`
  flex: 1;
  height: ${(p: any) => (p.isMobile ? 'calc(100% - 105px)' : 'calc(100% - 65px)')};
  background: ${(p: any) => (p.isMobile ? undefined : color.grayish.G950)};
  width: 100%;
  overflow: auto;
  display: flex;
  flex-direction: column;
  & > .header {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    padding: 10px 0;
  }
  & > .content {
    width: 100%;
    display: flex;
    flex-wrap: wrap;
    height: 100%;
    justify-content: flex-start;
    align-items: flex-start;
    padding: 0px 20px 20px 20px;
  }
`;

export const Spacer = styled.div`
  display: flex;
  min-height: 10px;
  min-width: 100%;
  height: 10px;
  width: 100%;
`;

function BodyComponent() {
  const { main, ui } = useStores();
  const [loading, setLoading] = useState(true);
  const screenWidth = useScreenWidth();
  const [openStartUpModel, setOpenStartUpModel] = useState<boolean>(false);
  const [checkboxIdToSelectedMapLanguage, setCheckboxIdToSelectedMapLanguage] = useState({});
  const [filterResult, setFilterResult] = useState<PersonType[]>(main.people);
  const closeModal = () => setOpenStartUpModel(false);
  const { peoplePageNumber } = ui;
  const history = useHistory();
  const isMobile = useIsMobile();
  const people = useFuse(main.people, ['owner_alias']).filter((f: any) => !f.hide) || [];
  async function loadMore(direction: number) {
    let currentPage = 1;
    currentPage = peoplePageNumber;

    let newPage = currentPage + direction;
    if (newPage < 1) newPage = 1;
    try {
      await main.getPeople({ page: newPage });
    } catch (e: any) {
      console.log('load failed', e);
    }
  }
  const loadForwardFunc = () => loadMore(1);
  const loadBackwardFunc = () => loadMore(-1);
  const { loadingBottom, handleScroll } = usePageScroll(loadForwardFunc, loadBackwardFunc);

  const onChangeLanguage = (optionId: any) => {
    const newCheckboxIdToSelectedMapLanguage = {
      ...checkboxIdToSelectedMapLanguage,
      ...{
        [optionId]: !checkboxIdToSelectedMapLanguage[optionId]
      }
    };
    setCheckboxIdToSelectedMapLanguage(newCheckboxIdToSelectedMapLanguage);
  };

  const toastsEl = (
    <EuiGlobalToastList
      toasts={ui.toasts}
      dismissToast={() => ui.setToasts([])}
      toastLifeTimeMs={3000}
    />
  );

  useEffect(() => {
    if (ui.meInfo) {
      main.getTribesByOwner(ui.meInfo.owner_pubkey || '');
    }
  }, [main, ui.meInfo]);

  useEffect(() => {
    setFilterResult(filterByCodingLanguage(main.people, checkboxIdToSelectedMapLanguage));
  }, [checkboxIdToSelectedMapLanguage]);

  // update search
  useEffect(() => {
    (async () => {
      await main.getPeople({ page: 1, resetPage: true });
      setLoading(false);
    })();
  }, [ui.searchText, ui.selectingPerson, main]);

  function selectPerson(id: number, unique_name: string, pubkey: string) {
    ui.setSelectedPerson(id);
    ui.setSelectingPerson(id);

    history.push(`/p/${pubkey}`);
  }

  if (loading) {
    return (
      <Body isMobile={isMobile} style={{ justifyContent: 'center', alignItems: 'center' }}>
        <EuiLoadingSpinner size="xl" />
      </Body>
    );
  }

  return (
    <Body
      isMobile={isMobile}
      onScroll={(e: any) => {
        handleScroll(e);
      }}
    >
      <div className="header">
        <PeopleHeader
          onChangeLanguage={onChangeLanguage}
          checkboxIdToSelectedMapLanguage={checkboxIdToSelectedMapLanguage}
        />

        <SearchTextInput
          small
          name="search"
          type="search"
          placeholder="Search"
          value={ui.searchText}
          style={{
            width: isMobile ? '95vw' : 240,
            height: 40,
            border: `1px solid ${color.grayish.G600}`,
            background: color.grayish.G600
          }}
          onChange={(e: any) => {
            ui.setSearchText(e);
          }}
        />
      </div>
      <div className="content">
        {(ui.searchText ? people : filterResult).map((t: any) => (
          <Person
            {...t}
            key={t.id}
            small={isMobile}
            squeeze={screenWidth < 1420}
            selected={ui.selectedPerson === t.id}
            select={selectPerson}
          />
        ))}
        {!(ui.searchText ? people : filterResult)?.length && <NoResults />}
        <PageLoadSpinner noAnimate show={loadingBottom} />
      </div>

      {openStartUpModel && (
        <StartUpModal closeModal={closeModal} dataObject={'getWork'} buttonColor={'primary'} />
      )}
      {toastsEl}
    </Body>
  );
}

export default observer(BodyComponent);
