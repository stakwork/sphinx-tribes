import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import {
  EuiFormFieldset,
  EuiLoadingSpinner,
  EuiPopover,
  EuiSelectable,
  EuiButton,
  EuiHighlight
} from '@elastic/eui';
import { observer } from 'mobx-react-lite';
import { useStores } from '../../store';
import { useIsMobile, usePageScroll } from '../../hooks';
import { SearchTextInput } from '../../components/common';
import NoResults from '../../people/utils/NoResults';
import PageLoadSpinner from '../../people/utils/PageLoadSpinner';
import { colors } from '../../config/colors';
import tags from './tags';
import Tag from './Tag';
import Tribe from './Tribe';

const Body = styled.div`
  flex: 1;
  height: calc(100vh - 60px);
  // padding-bottom:80px;
  width: 100%;
  overflow: auto;
  background: ${colors.dark.tribesBackground};
  display: flex;
  flex-direction: column;
  align-items: center;
`;
const Column = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 10px;
  // max-width:900px;
  width: 100%;
`;

function BodyComponent() {
  const { main, ui } = useStores();
  const [selected, setSelected] = useState('');
  const [tagsPop, setTagsPop] = useState(false);
  const [tagOptions, setTagOptions] = useState(ui.tags);
  const [loading, setLoading] = useState(true);
  const [loadingList, setLoadingList] = useState(true);

  const { tribesPageNumber } = ui;

  const isMobile = useIsMobile();

  const selectedTags = tagOptions.filter((t: any) => t.checked === 'on');
  const showTagCount = selectedTags.length > 0 ? true : false;

  function selectTribe(uuid: string, unique_name: string) {
    setSelected(uuid);
    if (!uuid) {
      window.history.pushState({}, 'Sphinx Tribes', '/t');
    } else if (unique_name && window.history.pushState) {
      window.history.pushState({}, 'Sphinx Tribes', `/t/${unique_name}`);
    }
  }

  async function loadMore(direction: number) {
    if (tagsPop) return;

    const currentPage = tribesPageNumber;
    let newPage = currentPage + direction;
    if (newPage < 1) newPage = 1;

    try {
      await main.getTribes({ page: newPage });
    } catch (e) {
      console.log(e);
    }
  }

  async function refreshList() {
    setLoadingList(true);

    // reset page will replace all results, this is good for a new search!
    await main.getTribes({ page: 1, resetPage: true });

    // do deeplink
    let deeplinkUn = '';
    if (window.location.pathname.startsWith('/t/')) {
      deeplinkUn = window.location.pathname.substr(3);
    }
    if (deeplinkUn) {
      const t = await main.getTribeByUn(deeplinkUn);
      setSelected(t.uuid);
      window.history.pushState({}, 'Sphinx Tribes', '/t');
    }
    setLoadingList(false);
    setLoading(false);
  }

  // do search update
  useEffect(() => {
    refreshList();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [ui.searchText, ui.tags]);

  let { tribes } = main;

  const loadForwardFunc = () => loadMore(1);
  const loadBackwardFunc = () => loadMore(-1);
  const { loadingTop, loadingBottom, handleScroll } = usePageScroll(
    loadForwardFunc,
    loadBackwardFunc
  );

  if (loading) {
    return (
      <Body style={{ justifyContent: 'center', alignItems: 'center' }}>
        <EuiLoadingSpinner size="xl" />
      </Body>
    );
  }

  // if NSFW not selected, filter out NSFW
  if (!selectedTags.find((f: any) => f.label === 'NSFW')) {
    tribes = tribes.filter((f: any) => !f.tags.includes('NSFW'));
  }

  const button = (
    <EuiButton
      iconType="arrowDown"
      iconSide="right"
      size="s"
      onClick={(e: any) => {
        e.stopPropagation();
        setTagsPop(!tagsPop);
      }}
    >
      {`Tags ${showTagCount ? `(${selectedTags.length})` : ''}`}
    </EuiButton>
  );

  return (
    <Body id="main" onScroll={handleScroll} style={{ paddingTop: 0 }}>
      <div
        style={{
          width: '100%',
          display: 'flex',
          justifyContent: 'flex-end',
          alignItems: 'flex-start',
          padding: 20,
          height: 62
        }}
      >
        <div style={{ display: 'flex', alignItems: 'baseline' }}>
          <EuiPopover
            panelPaddingSize="none"
            button={button}
            isOpen={tagsPop}
            closePopover={() => setTagsPop(false)}
          >
            <EuiSelectable
              searchable
              options={tagOptions}
              renderOption={(option: any, searchValue: any) => (
                <div
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    opacity: loadingList ? 0.5 : 1
                  }}
                >
                  <Tag type={option.label} iconOnly />
                  <EuiHighlight
                    search={searchValue}
                    style={{
                      fontSize: 11,
                      marginLeft: 5,
                      color: tags[option.label].color
                    }}
                  >
                    {option.label}
                  </EuiHighlight>
                </div>
              )}
              listProps={{ rowHeight: 30 }} // showIcons:false
              onChange={(opts: any) => {
                if (!loadingList) {
                  setTagOptions(opts);
                  ui.setTags(opts);
                }
              }}
            >
              {(list: any, search: any) => (
                <div style={{ width: 220 }}>
                  {search}
                  {list}
                </div>
              )}
            </EuiSelectable>
          </EuiPopover>

          <SearchTextInput
            name="search"
            type="search"
            small={isMobile}
            placeholder="Search"
            value={ui.searchText}
            style={{
              width: 204,
              height: 40,
              background: '#111',
              color: '#fff',
              border: 'none',
              marginLeft: 20
            }}
            onChange={(e: any) => {
              ui.setSearchText(e);
            }}
          />
        </div>
      </div>
      <Column className="main-wrap">
        <PageLoadSpinner show={loadingTop} />
        <EuiFormFieldset style={{ width: '100%', paddingBottom: 0 }} className="container">
          <div style={{ justifyContent: 'center' }} className="row">
            {tribes.length ? (
              tribes.map((t: any) => (
                <Tribe {...t} key={t.uuid} selected={selected === t.uuid} select={selectTribe} />
              ))
            ) : (
              <NoResults />
            )}
          </div>
        </EuiFormFieldset>
        <PageLoadSpinner noAnimate show={loadingBottom} />
      </Column>
    </Body>
  );
}
export default observer(BodyComponent);
