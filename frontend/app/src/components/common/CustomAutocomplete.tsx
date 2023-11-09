import { EuiText } from '@elastic/eui';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { AutoCompleteProps } from 'components/interfaces';
import { colors } from '../../config/colors';
import ImageButton from './ImageButton';

interface styledProps {
  color?: any;
}

const SearchOuterContainer = styled.div<styledProps>`
  min-height: 347px;
  max-height: 347px;
  min-width: 336px;
  max-width: 336px;
  overflow-x: hidden;
  overflow-y: scroll;
  background: ${(p: any) => p?.color && p.color.pureWhite};
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 25px 0 0 0;

  .SearchInput {
    background: ${(p: any) => p?.color && p.color.pureWhite};
    border: 1px solid ${(p: any) => p?.color && p?.color?.grayish.G600};
    border-radius: 200px;
    width: 292px;
    height: 40px;
    outline: none;
    overflow: hidden;
    caret-color: ${(p: any) => p?.color && p.color.light_blue100};
    margin-bottom: 8px;
    padding: 0px 18px;

    :focus-visible {
      background: ${(p: any) => p?.color && p.color.pureWhite};
      border: 1px solid ${(p: any) => p?.color && p.color.grayish.G600};
      border-radius: 200px;
      outline: none;
    }
    :active {
      .SearchText {
        outline: none;
        background: ${(p: any) => p?.color && p.color.pureWhite};
        border: 1px solid ${(p: any) => p?.color && p.color.grayish.G600};
        border-radius: 200px;
        outline: none;
      }
    }
  }
  .PeopleList {
    background: ${(p: any) => p?.color && p.color.pureWhite};
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
          font-family: 'Barlow';
          font-style: normal;
          font-weight: 500;
          font-size: 13px;
          line-height: 16px;
          color: ${(p: any) => p?.color && p.color.grayish.G10};
          margin-left: 10px;
        }
      }
    }
  }
`;
const AutoComplete = (props: AutoCompleteProps) => {
  const color = colors['light'];
  const [searchValue, setSearchValue] = useState<string>('');
  const [peopleData, setPeopleData] = useState<any>(props.peopleList);

  const handler = (e: any) => {
    setSearchValue(e.target.value);
    const result = props?.peopleList.filter(
      (x: any) => x?.owner_alias.toLowerCase()?.includes(e.target.value.toLowerCase())
    );
    setPeopleData(result);
  };

  useEffect(() => {
    if (searchValue === '') {
      setPeopleData(props.peopleList);
    }
  }, [searchValue, props]);

  return (
    <SearchOuterContainer color={color}>
      <input
        className="SearchInput"
        onChange={handler}
        placeholder={'Search'}
        style={{
          background: color.pureWhite,
          color: color.text1,
          fontFamily: 'Barlow'
        }}
      />
      <div className="PeopleList">
        {peopleData?.slice(0, 5)?.map((value: any, index: number) => (
          <div className="People" key={index}>
            <div className="PeopleDetailContainer">
              <div className="ImageContainer">
                <img
                  src={value.img || '/static/person_placeholder.png'}
                  alt={'user'}
                  height={'100%'}
                  width={'100%'}
                />
              </div>
              <EuiText className="PeopleName">{value.owner_alias}</EuiText>
            </div>
            <ImageButton
              buttonText={'Assign'}
              ButtonContainerStyle={{
                width: '74.58px',
                height: '32px'
              }}
              buttonAction={() => {
                props?.handleAssigneeDetails(value);
              }}
            />
          </div>
        ))}
      </div>
    </SearchOuterContainer>
  );
};

export default AutoComplete;
