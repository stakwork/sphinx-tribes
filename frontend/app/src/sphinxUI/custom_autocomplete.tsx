import { EuiText } from '@elastic/eui';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';

const AutoComplete = (props) => {
  const [searchValue, setSearchValue] = useState<string>('');
  const [peopleData, setPeopleData] = useState<any>(props.peopleList);

  const handler = (e) => {
    setSearchValue(e.target.value);
    let result = props?.peopleList.filter((x) =>
      x?.owner_alias.toLowerCase()?.includes(e.target.value.toLowerCase())
    );
    setPeopleData(result);
  };

  useEffect(() => {
    if (searchValue === '') {
      setPeopleData(props.peopleList);
    }
  }, [searchValue, props]);

  return (
    <SearchOuterContainer>
      <input
        className="SearchInput"
        onChange={handler}
        placeholder={'Search'}
        style={{
          background: '#fff',
          color: '#292C33',
          fontFamily: 'Barlow'
        }}
      />
      <div className="PeopleList">
        {peopleData?.slice(0, 5)?.map((value, index) => {
          return (
            <div className="People" key={index}>
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'center',
                  alignItems: 'center'
                }}
              >
                <div
                  style={{
                    height: '32px',
                    width: '32px',
                    borderRadius: '50%',
                    overflow: 'hidden',
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center'
                  }}
                >
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
                  }}
                >
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
                onClick={() => {
                  props?.handleAssigneeDetails(value);
                }}
              >
                <EuiText
                  style={{
                    fontFamily: 'Barlow',
                    fontStyle: 'normal',
                    fontWeight: '500',
                    fontSize: '13px',
                    lineHeight: '16px',
                    letterSpacing: '0.01em',
                    color: '#5F6368'
                  }}
                >
                  Assign
                </EuiText>
              </div>
            </div>
          );
        })}
      </div>
    </SearchOuterContainer>
  );
};

export default AutoComplete;

const SearchOuterContainer = styled.div`
  min-height: 347px;
  max-height: 347px;
  min-width: 336px;
  max-width: 336px;
  overflow-x: hidden;
  overflow-y: scroll;
  background: #fff;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 25px 0 0 0;

  .SearchInput {
    background: #ffffff;
    border: 1px solid #dde1e5;
    border-radius: 200px;
    width: 292px;
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
    background: #ffffff;
    .People {
      height: 32px;
      min-width: 291.5813903808594px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-top: 16px;
    }
  }
`;
