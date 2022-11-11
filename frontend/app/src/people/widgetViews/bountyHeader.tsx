import { EuiText } from '@elastic/eui';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import api from '../../api';
import { colors } from '../../colors';
import { useIsMobile } from '../../hooks';
import IconButton from '../../sphinxUI/icon_button';
import SearchBar from '../../sphinxUI/search_bar';
import { useStores } from '../../store';
import StartUpModal from '../utils/start_up_modal';

const BountyHeader = ({ selectedWidget, setShowFocusView, scrollValue }) => {
  const { main, ui } = useStores();
  const isMobile = useIsMobile();
  const [peopleList, setPeopleList] = useState<Array<any> | null>(null);
  const [activeBounty, setActiveBounty] = useState<Array<any> | number | null>(0);
  const [openStartUpModel, setOpenStartUpModel] = useState<boolean>(false);
  const closeModal = () => setOpenStartUpModel(false);
  const showModal = () => setOpenStartUpModel(true);
  const color = colors['light'];
  useEffect(() => {
    async function getPeopleList() {
      if (selectedWidget === 'wanted') {
        try {
          const bounty = await api.get(
            `people/wanteds?page=1&resetPage=true&search=&sortBy=created&limit=100`
          );
          const response = await api.get(`people?page=1&search=&sortBy=last_login&limit=100`);
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
        <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            width: '100%',
            minHeight: '80px',
            alignItems: 'center',
            position: 'sticky',
            top: 0,
            zIndex: '1',
            background: 'inherit',
            boxShadow: scrollValue ? ' 0px 1px 6px rgba(0, 0, 0, 0.07)' : '',
            borderBottom: scrollValue
              ? `1px solid ${color.grayish.G600}`
              : `0px solid ${color.grayish.G600}`
          }}
        >
          <BountyHeaderDesk>
            <B>
              <IconButton
                text={'Post a Bounty'}
                endingIcon={'add'}
                width={204}
                height={48}
                color={'success'}
                style={{
                  color: color.pureWhite,
                  fontSize: '16px',
                  fontWeight: '600',
                  textDecoration: 'none'
                }}
                hoverColor={color.button_primary.hover}
                activeColor={color.button_primary.active}
                shadowColor={color.button_primary.shadow}
                iconStyle={{
                  fontSize: '16px',
                  fontWeight: '400',
                  top: '17px',
                  right: '18px'
                }}
                onClick={() => {
                  if (ui.meInfo && ui.meInfo?.owner_alias) {
                    setShowFocusView(true);
                  } else {
                    showModal();
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
                  fontFamily: 'Barlow',
                  color: color.text2
                }}
                onChange={(e) => {
                  ui.setSearchText(e);
                }}
                iconStyle={{
                  top: '13px'
                }}
                TextColor={color.grayish.G400}
                TextColorHover={color.grayish.G100}
                border={`1px solid ${color.grayish.G500}`}
                borderHover={`1px solid ${color.grayish.G400}`}
                borderActive={'1px solid #A3C1FF'}
                iconColor={color.grayish.G300}
                iconColorHover={color.grayish.G100}
              />
              <div
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  marginLeft: '33px'
                }}
              >
                <img src="/static/copy.svg" alt="" height={22} width={18} />
                <EuiText
                  style={{
                    color: color.grayish.G200,
                    fontWeight: '500',
                    fontSize: '16px',
                    lineHeight: '19px',
                    fontFamily: 'Barlow',
                    marginLeft: '10px',
                    height: '51px',
                    width: '153px',
                    display: 'flex',
                    alignItems: 'center'
                  }}
                >
                  <span
                    style={{
                      color: color.pureBlack
                    }}
                  >
                    {activeBounty}
                  </span>
                  &nbsp; Bounties opened
                </EuiText>
              </div>
            </B>
            <D>
              <EuiText className="DText" color={color.grayish.G200}>
                Developers
              </EuiText>
              <div className="ImageOuterContainer">
                {peopleList &&
                  peopleList?.slice(0, 3).map((val, index) => {
                    return (
                      <DevelopersImageContainer
                        style={{
                          zIndex: 3 - index,
                          marginLeft: index > 0 ? '-14px' : ''
                        }}
                      >
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
                  fontFamily: 'Barlow',
                  color: '#222E3A'
                }}
              >
                {peopleList && peopleList?.length}
              </EuiText>
            </D>
          </BountyHeaderDesk>
        </div>
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
              TextColor={color.grayish.G400}
              TextColorHover={color.grayish.G100}
              border={`1px solid ${color.grayish.G500}`}
              borderHover={`1px solid ${color.grayish.G400}`}
              borderActive={'1px solid #A3C1FF'}
              iconColor={color.grayish.G300}
              iconColorHover={color.grayish.G100}
            />
            {/*
            
            // TODO: add filter when have functionality.

            <IconButton
              text={'Filter'}
              color={'transparent'}
              leadingIcon={'tune'}
              width={80}
              height={48}
              style={{
                color: color.grayish.G200,
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
            /> */}
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
              hoverColor={color.button_primary.hover}
              activeColor={color.button_primary.active}
              shadowColor={color.button_primary.shadow}
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
                      }}
                    >
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
                  fontWeight: '500',
                  color: '#222E3A'
                }}
              >
                {peopleList && peopleList?.length}
              </EuiText>
            </DevelopersContainerMobile>
          </ShortActionContainer>
        </BountyHeaderMobile>
      )}
      {openStartUpModel && (
        <StartUpModal closeModal={closeModal} dataObject={'createWork'} buttonColor={'success'} />
      )}
    </>
  );
};

export default BountyHeader;

const BountyHeaderDesk = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  min-width: 1100px;
  max-width: 1100px;
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
  .DText {
    font-size: 16px;
    font-family: Barlow;
    font-weight: 500;
  }
  .ImageOuterContainer {
    display: flex;
    flex-direction: row;
    align-items: center;
    color: #909baa;
    padding: 0 10px;
  }
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
