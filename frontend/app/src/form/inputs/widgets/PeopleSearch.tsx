import { EuiText } from '@elastic/eui';
import React, { useCallback, useEffect, useState } from 'react';
import styled from 'styled-components';
import { colors } from '../../../colors';
import ImageButton from '../../../sphinxUI/Image_button';

const InvitePeopleSearch = (props) => {
  const color = colors['light'];
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
    <SearchOuterContainer color={color}>
      <input
        value={searchValue}
        className="SearchInput"
        onChange={(e) => {
          handler(e, '');
        }}
        placeholder={'Search'}
        style={{
          background: color.pureWhite,
          color: color.text1,
          fontFamily: 'Barlow'
        }}
      />
      <div className="PeopleList">
        {peopleData?.slice(0, 50)?.map((value) => {
          return (
            <div className="People" key={value.id}>
              <div className="PeopleDetailContainer">
                <div className="ImageContainer">
                  <img
                    src={value.img || '/static/person_placeholder.png'}
                    alt={'user-image'}
                    height={'100%'}
                    width={'100%'}
                  />
                </div>
                <EuiText className="PeopleName">{value.owner_alias}</EuiText>
              </div>
              <ImageButton
                buttonText={inviteNameId === value?.id ? ' Invited' : ' Invite'}
                ButtonContainerStyle={{
                  width: '74.58px',
                  height: '32px'
                }}
                buttonAction={(e) => {
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
                }}
              />
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

interface styledProps {
  color?: any;
}

const SearchOuterContainer = styled.div<styledProps>`
  min-height: 256x;
  max-height: 256x;
  min-width: 302px;
  max-width: 302px;
  background: ${(p) => p?.color && p?.color?.pureWhite};
  display: flex;
  flex-direction: column;
  align-items: center;

  .SearchInput {
    background: ${(p) => p?.color && p?.color?.pureWhite};
    border: 1px solid ${(p) => p?.color && p?.color?.grayish.G600};
    border-radius: 200px;
    width: 302px;
    height: 40px;
    outline: none;
    overflow: hidden;
    caret-color: ${(p) => p?.color && p?.color?.textBlue1};
    margin-bottom: 8px;
    padding: 0px 18px;

    :focus-visible {
      background: ${(p) => p?.color && p?.color?.pureWhite};
      border: 1px solid ${(p) => p?.color && p?.color?.grayish.G600};
      border-radius: 200px;
      outline: none;
    }
    :active {
      .SearchText {
        outline: none;
        background: ${(p) => p?.color && p?.color?.pureWhite};
        border: 1px solid ${(p) => p?.color && p?.color?.grayish.G600};
        border-radius: 200px;
        outline: none;
      }
    }
  }
  .PeopleList {
    background: ${(p) => p?.color && p?.color?.grayish.G950};
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
      .PeopleDetailContainer {
        display: flex;
        justify-content: center;
        align-items: center;
        .ImageContainer {
          height: 32px;
          width: 32px;
          border-radius: 50%;
          overflow: hidden;
          display: flex;
          justify-content: center;
          align-items: center;
          object-fit: cover;
        }
        .PeopleName {
          font-family: Barlow;
          font-style: normal;
          font-weight: 500;
          font-size: 13px;
          line-height: 16px;
          color: ${(p) => p?.color && p?.color?.grayish.G10};
          margin-left: 10px;
        }
      }
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
        color: ${(p) => p?.color && p?.color?.grayish.G50};
        word-spacing: 0.08em;
      }
    }
  }
`;
