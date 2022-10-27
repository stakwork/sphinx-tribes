import { EuiText } from '@elastic/eui';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { useIsMobile } from '../../hooks';
import IconButton from '../../sphinxUI/icon_button';
import SearchBar from '../../sphinxUI/search_bar';
import { useStores } from '../../store';

const BountyHeader = ({ selectedWidget, setShowFocusView }) => {
  const { main, ui } = useStores();
  const isMobile = useIsMobile();
  const [peopleList, setPeopleList] = useState<Array<any> | null>(null);
  const [activeBounty, setActiveBounty] = useState<Array<any> | number | null>(0);
  useEffect(() => {
    async function getPeopleList() {
      if (selectedWidget === 'wanted') {
        try {
          const response = await main.getPeople({ page: 1 });
          const bounty = await main.getPeopleWanteds({ page: 1 });
          setPeopleList(response);
          setActiveBounty(bounty?.length);
        } catch (error) {
          console.log(error);
        }
      } else {
        setPeopleList(null);
      }
    }
    getPeopleList();
  }, [main, selectedWidget]);

  return (
    <>
      {!isMobile ? (
        // desktop view
        <BountyHeaderDesk>
          <B>
            <IconButton
              text={'Post a Bounty'}
              endingIcon={'add'}
              width={204}
              height={48}
              color={'success'}
              style={{
                color: '#fff',
                fontSize: '16px',
                fontWeight: '600',
                textDecoration: 'none'
              }}
              hoverColor={'#3CBE88'}
              activeColor={'#2FB379'}
              shadowColor={'rgba(73, 201, 152, 0.5)'}
              iconStyle={{
                fontSize: '16px',
                fontWeight: '400',
                top: '18px',
                right: '24px'
              }}
              onClick={() => {
                if (ui.meInfo && ui.meInfo?.owner_alias) {
                  setShowFocusView(true);
                } else {
                  ui.setShowSignIn(true);
                }
              }}
            />
            <SearchBar
              name="search"
              type="search"
              placeholder="Search"
              value={ui.searchText}
              style={{
                width: 204,
                height: 48,
                background: 'transparent',
                marginLeft: '16px',
                fontFamily: 'Barlow'
              }}
              onChange={(e) => {
                ui.setSearchText(e);
              }}
              iconStyle={{
                top: '13px'
              }}
              TextColor={'#B0B7BC'}
              TextColorHover={'#8E969C'}
              border={'1px solid #D0D5D8'}
              borderHover={'1px solid #BAC1C6'}
              borderActive={'1px solid #A3C1FF'}
              iconColor={'#B0B7BC'}
              iconColorHover={'#8E969C'}
            />
            <IconButton
              text={`${activeBounty} Bounties opened`}
              leadingImg={'/static/copy.svg'}
              width={230}
              height={48}
              color={'transparent'}
              style={{
                color: '#909BAA',
                fontSize: '16px',
                fontWeight: '500',
                cursor: 'default',
                textDecoration: 'none'
              }}
              leadingImgStyle={{
                height: '22px',
                width: '18px',
                marginRight: '10px'
              }}
            />
            <IconButton
              text={'Filter'}
              color={'transparent'}
              leadingIcon={'tune'}
              width={80}
              height={48}
              style={{
                color: '#909BAA',
                fontSize: '16px',
                fontWeight: '500',
                textDecoration: 'none'
              }}
              iconStyle={{
                fontSize: '18px',
                fontWeight: '500'
              }}
              onClick={() => {
                console.log('filter');
              }}
            />
          </B>
          <D>
            <EuiText
              color={'#909BAA'}
              style={{
                fontSize: '16px',
                fontFamily: 'Barlow',
                fontWeight: '500'
              }}>
              Developers
            </EuiText>
            <div
              style={{
                display: 'flex',
                flexDirection: 'row',
                alignItems: 'center',
                color: '#909BAA',
                padding: '0 10px'
              }}>
              {peopleList &&
                peopleList?.slice(0, 3).map((val, index) => {
                  return (
                    <DevelopersImageContainer
                      style={{
                        zIndex: 3 - index,
                        marginLeft: index > 0 ? '-14px' : ''
                      }}>
                      <img
                        height={'23px'}
                        width={'23px'}
                        src={val?.img || '/static/person_placeholder.png'}
                        alt={''}
                        style={{
                          borderRadius: '50%'
                        }}
                      />
                    </DevelopersImageContainer>
                  );
                })}
            </div>
            <EuiText
              style={{
                fontSize: '16px',
                fontWeight: '600',
                fontFamily: 'Barlow'
              }}>
              {peopleList && peopleList?.length}
            </EuiText>
          </D>
        </BountyHeaderDesk>
      ) : (
        <BountyHeaderMobile>
          <LargeActionContainer>
            <SearchBar
              name="search"
              type="search"
              placeholder="Search"
              value={ui.searchText}
              style={{
                width: 240,
                height: 32,
                background: 'transparent',
                marginLeft: '16px',
                fontFamily: 'Barlow'
              }}
              onChange={(e) => {
                ui.setSearchText(e);
              }}
              iconStyle={{
                top: '4px'
              }}
              border={'1px solid #D0D5D8'}
              borderHover={'1px solid #BAC1C6'}
              borderActive={'1px solid #A3C1FF'}
              TextColor={'#B0B7BC'}
              TextColorHover={'#8E969C'}
              iconColor={'#B0B7BC'}
              iconColorHover={'#8E969C'}
            />
            <IconButton
              text={'Filter'}
              color={'transparent'}
              leadingIcon={'tune'}
              width={80}
              height={48}
              style={{
                color: '#909BAA',
                fontSize: '16px',
                fontWeight: '500',
                textDecoration: 'none',
                transform: 'none'
              }}
              iconStyle={{
                fontSize: '16px',
                fontWeight: '500'
              }}
              onClick={() => {
                console.log('filter');
              }}
            />
          </LargeActionContainer>
          <ShortActionContainer>
            <IconButton
              text={'Post a Bounty'}
              endingIcon={'add'}
              width={130}
              height={30}
              color={'success'}
              style={{
                color: '#fff',
                fontSize: '12px',
                fontWeight: '600',
                textDecoration: 'none',
                transform: 'none'
              }}
              hoverColor={'#3CBE88'}
              activeColor={'#2FB379'}
              shadowColor={'rgba(73, 201, 152, 0.5)'}
              iconStyle={{
                fontSize: '12px',
                fontWeight: '600',
                right: '8px',
                top: '9px'
              }}
              onClick={() => {
                if (ui.meInfo && ui.meInfo?.owner_alias) {
                  setShowFocusView(true);
                } else {
                  ui.setShowSignIn(true);
                }
              }}
            />
            <IconButton
              text={`${activeBounty} Bounties opened`}
              leadingImg={'/static/copy.svg'}
              width={'fit-content'}
              height={48}
              color={'transparent'}
              style={{
                color: '#909BAA',
                fontSize: '12px',
                fontWeight: '500',
                cursor: 'default',
                textDecoration: 'none',
                padding: 0
              }}
              leadingImgStyle={{
                height: '21px',
                width: '18px',
                marginRight: '4px'
              }}
            />
            <DevelopersContainerMobile>
              {peopleList &&
                peopleList?.slice(0, 3).map((val, index) => {
                  return (
                    <DevelopersImageContainer
                      style={{
                        zIndex: 3 - index,
                        marginLeft: index > 0 ? '-14px' : ''
                      }}>
                      <img
                        height={'20px'}
                        width={'20px'}
                        src={val?.img || '/static/person_placeholder.png'}
                        alt={''}
                        style={{
                          borderRadius: '50%'
                        }}
                      />
                    </DevelopersImageContainer>
                  );
                })}
              <EuiText
                style={{
                  fontSize: '14px',
                  fontFamily: 'Barlow',
                  fontWeight: '500'
                }}>
                {peopleList && peopleList?.length}
              </EuiText>
            </DevelopersContainerMobile>
          </ShortActionContainer>
        </BountyHeaderMobile>
      )}
    </>
  );
};

export default BountyHeader;

const BountyHeaderDesk = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  padding: 10px 20px;
  align-items: center;
`;

const B = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-evenly;
  align-items: center;
`;

const D = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 0 20px;
`;

const DevelopersImageContainer = styled.div`
  height: 28px;
  width: 28px;
  border-radius: 50%;
  background: #fff;
  overflow: hidden;
  position: static;
  display: flex;
  justify-content: center;
  align-items: center;
`;

const BountyHeaderMobile = styled.div`
  display: flex;
  flex-direction: column;
  padding: 10px 0px;
`;

const ShortActionContainer = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-evenly;
`;

const DevelopersContainerMobile = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
`;

const LargeActionContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
`;
