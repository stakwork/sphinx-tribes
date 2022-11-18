import { EuiText } from '@elastic/eui';
import React, { useCallback, useEffect, useState } from 'react';
import styled from 'styled-components';

const InvitePeopleSearch = (props) => {
  const [searchValue, setSearchValue] = useState<string>('');
  const [peopleData, setPeopleData] = useState<any>(props?.peopleList);
  const [inviteNameId, setInviteNameId] = useState<number>(0);

  const handler = useCallback((e, value) => {
    if (value === '') {
      setSearchValue(e.target.value);
      const result = props?.peopleList.filter((x) =>
        x?.owner_alias.toLowerCase()?.includes(e.target.value.toLowerCase())
      );
      setPeopleData(result);
    } else {
      const result = props?.peopleList.filter((x) =>
        x?.owner_alias.toLowerCase()?.includes(value.toLowerCase())
      );
      setPeopleData(result);
    }
  }, []);

  useEffect(() => {
    if (searchValue === '') {
      setPeopleData(props.peopleList);
    }
  }, [searchValue, props]);

  return (
    <SearchOuterContainer>
      <input
        value={searchValue}
        className="SearchInput"
        onChange={(e) => {
          handler(e, '');
        }}
        placeholder={'Search'}
        style={{
          background: '#fff',
          color: '#292C33',
          fontFamily: 'Barlow'
        }}
      />
      <div className="PeopleList">
        {peopleData?.slice(0, 10)?.map((value) => {
          return (
            <div className="People" key={value.id}>
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'center',
                  alignItems: 'center'
                }}>
                <div
                  style={{
                    height: '32px',
                    width: '32px',
                    borderRadius: '50%',
                    overflow: 'hidden',
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    objectFit: 'cover'
                  }}>
                  <img
                    src={value.img || '/static/person_placeholder.png'}
                    alt={'user-image'}
                    height={'100%'}
                    width={'100%'}
                  />
                </div>
                <EuiText
                  style={{
                    fontFamily: 'Barlow',
                    fontStyle: 'normal',
                    fontWeight: '500',
                    fontSize: '13px',
                    lineHeight: '16px',
                    color: '#3C3F41',
                    marginLeft: '10px'
                  }}>
                  {value.owner_alias}
                </EuiText>
              </div>
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'center',
                  alignItems: 'center',
                  minHeight: '32px',
                  minWidth: '74.58px',
                  borderRadius: '50px',
                  background: '#FFFFFF',
                  border: ' 1px solid #DDE1E5',
                  cursor: 'pointer'
                }}
                onClick={(e) => {
                  handler('', value.owner_alias);
                  setInviteNameId(value.id);
                  props?.handleChange({
                    owner_alias: value.owner_alias,
                    owner_pubkey: value.owner_pubkey,
                    img: value.img,
                    value: value.owner_pubkey,
                    label: `${value.owner_alias} (${value.owner_alias
                      .toLowerCase()
                      .replace(' ', '')})`
                  });
                  setSearchValue(value.owner_alias);
                }}>
                <EuiText
                  style={{
                    fontFamily: 'Barlow',
                    fontStyle: 'normal',
                    fontWeight: '500',
                    fontSize: '13px',
                    lineHeight: '16px',
                    letterSpacing: '0.01em',
                    color: '#5F6368'
                  }}>
                  {inviteNameId === value?.id ? ' Invited' : ' Invite'}
                </EuiText>
              </div>
            </div>
          );
        })}
        {peopleData?.length === 0 && (
          <div className="no_result_container">
            <EuiText className="no_result_text">No Result Found</EuiText>
          </div>
        )}
      </div>
    </SearchOuterContainer>
  );
};

export default InvitePeopleSearch;

const SearchOuterContainer = styled.div`
  min-height: 256x;
  max-height: 256x;
  min-width: 302px;
  max-width: 302px;
  //   overflow: visible;
  background: #fff;
  display: flex;
  flex-direction: column;
  align-items: center;

  .SearchInput {
    background: #ffffff;
    border: 1px solid #dde1e5;
    border-radius: 200px;
    width: 302px;
    height: 40px;
    outline: none;
    overflow: hidden;
    caret-color: #a3c1ff;
    margin-bottom: 8px;
    padding: 0px 18px;

    :focus-visible {
      background: #ffffff;
      border: 1px solid #dde1e5;
      border-radius: 200px;
      outline: none;
    }
    :active {
      .SearchText {
        outline: none;
        background: #ffffff;
        border: 1px solid #dde1e5;
        border-radius: 200px;
        outline: none;
      }
    }
  }
  .PeopleList {
    background: #f2f3f5;
    width: 400px;
    padding: 0 49px 16px;
    min-height: 256px;
    max-height: 256px;
    overflow-y: scroll;
    .People {
      height: 32px;
      min-width: 291.5813903808594px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-top: 16px;
    }
    .no_result_container {
      display: flex;
      height: 210px;
      justify-content: center;
      align-items: center;
      .no_result_text {
        font-family: Barlow;
        font-size: 16px;
        font-weight: 600;
        color: #5f6368;
        word-spacing: 0.08em;
      }
    }
  }
`;
