import { Button, SearchTextInput } from 'components/common';
import { usePageScroll } from 'hooks';
import NoResults from 'people/utils/NoResults';
import PageLoadSpinner from 'people/utils/PageLoadSpinner';
import React from 'react';
import { useHistory } from 'react-router-dom';
import { useStores } from 'store';
import { queryLimit } from 'store/main';
import styled from 'styled-components';
import { observer } from 'mobx-react-lite';
import Person from '../Person';
const PeopleScroller = styled.div`
  overflow-y: overlay !important;
  width: 100%;
  height: 100%;
`;

const DBack = styled.div`
  min-height: 64px;
  height: 64px;
  display: flex;
  padding-right: 10px;
  align-items: center;
  justify-content: space-between;
  background: #ffffff;
  box-shadow: 0px 1px 6px rgba(0, 0, 0, 0.07);
  z-index: 0;
`;

const PeopleListContainer = styled.div`
  position: relative;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  background: #ffffff;
  width: 265px;
  overflow-y: overlay !important;

  * {
    scrollbar-width: 6px;
    scrollbar-color: rgba(176, 183, 188, 0.25);
  }

  /* Works on Chrome, Edge, and Safari */
  *::-webkit-scrollbar {
    width: 6px;
  }

  *::-webkit-scrollbar-thumb {
    background-color: rgba(176, 183, 188, 0.25);
    background: rgba(176, 183, 188, 0.25);
    width: 6px;
    border-radius: 10px;
    background-clip: padding-box;
  }

  ::-webkit-scrollbar-track-piece:start {
    background: transparent url('images/backgrounds/scrollbar.png') repeat-y !important;
  }

  ::-webkit-scrollbar-track-piece:end {
    background: transparent url('images/backgrounds/scrollbar.png') repeat-y !important;
  }
`;

export const PeopleList = observer(() => {
  const { main, ui } = useStores();
  const { peoplePageNumber } = ui || {};
  const history = useHistory();
  const personId = ui.selectedPerson;

  async function loadMorePeople(direction: number) {
    let newPage = peoplePageNumber + direction;
    if (newPage < 1) {
      newPage = 1;
    }
    await main.getPeople({ page: newPage });
  }

  function goBack() {
    ui.setSelectingPerson(0);
    history.goBack();
  }

  function selectPerson(id: number, unique_name: string, pubkey: string) {
    ui.setSelectedPerson(id);
    ui.setSelectingPerson(id);

    history.replace(`/p/${pubkey}`);
  }

  const people: any = (main.people && main.people.filter((f: any) => !f.hide)) || [];

  const { loadingTop, loadingBottom, handleScroll } = usePageScroll(
    () => loadMorePeople(1),
    () => loadMorePeople(-1)
  );

  return (
    <PeopleListContainer>
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
          onChange={(e: any) => {
            ui.setSearchText(e);
          }}
        />
      </DBack>

      <PeopleScroller
        style={{ width: '100%', overflowY: 'auto', height: '100%' }}
        onScroll={handleScroll}
      >
        <PageLoadSpinner show={loadingTop} />;
        {people?.length ? (
          people.map((t: any, i: number) => (
            <Person
              key={`${t.id}_${i}`}
              {...t}
              selected={personId === t.id}
              hideActions={true}
              small={true}
              select={selectPerson}
            />
          ))
        ) : (
          <NoResults />
        )}
        {/* make sure you can always scroll ever with too few people */}
        {people?.length < queryLimit && <div style={{ height: 400 }} />}
      </PeopleScroller>

      <PageLoadSpinner
        noAnimate
        show={loadingBottom}
        style={{ position: 'absolute', bottom: 0, left: 0 }}
      />
    </PeopleListContainer>
  );
});
