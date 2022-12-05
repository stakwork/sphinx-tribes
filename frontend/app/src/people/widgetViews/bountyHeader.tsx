/* eslint-disable @typescript-eslint/no-unused-vars */
import { EuiCheckboxGroup, EuiPopover, EuiText } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import api from '../../api';
import { colors } from '../../colors';
import { useIsMobile } from '../../hooks';
import IconButton from '../../sphinxUI/icon_button';
import ImageButton from '../../sphinxUI/Image_button';
import SearchBar from '../../sphinxUI/search_bar';
import { useStores } from '../../store';
import { coding_languages, GetValue, status } from '../utils/language_label_style';
import StartUpModal from '../utils/start_up_modal';

const Status = GetValue(status);
const Coding_Languages = GetValue(coding_languages);

const BountyHeader = ({
  selectedWidget,
  setShowFocusView,
  scrollValue,
  onChangeStatus,
  onChangeLanguage,
  checkboxIdToSelectedMap,
  checkboxIdToSelectedMapLanguage
}) => {
  const color = colors['light'];
  const { main, ui } = useStores();
  const isMobile = useIsMobile();
  const [peopleList, setPeopleList] = useState<Array<any> | null>(null);
  const [developerCount, setDeveloperCount] = useState<number>(0);
  const [activeBounty, setActiveBounty] = useState<Array<any> | number | null>(0);
  const [openStartUpModel, setOpenStartUpModel] = useState<boolean>(false);
  const closeModal = () => setOpenStartUpModel(false);
  const showModal = () => setOpenStartUpModel(true);
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);

  const onButtonClick = () => setIsPopoverOpen((isPopoverOpen) => !isPopoverOpen);
  const closePopover = () => setIsPopoverOpen(false);
  useEffect(() => {
    // eslint-disable-next-line func-style
    async function getPeopleList() {
      if (selectedWidget === 'wanted') {
        try {
          /*
          
           * TODO : Since this PR is merged only in people-test we will be using this api, when it will be merge in master then remove this fetch and use api.get() function.

           */

          const responseNew = await fetch(
            'https://people-test.sphinx.chat/people/wanteds/header'
          ).then((response) => response.json());
          setPeopleList(responseNew.people);
          setDeveloperCount(responseNew?.developer_count || 0);
          setActiveBounty(responseNew?.bounties_count);
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
            boxShadow: scrollValue ? `0px 1px 6px ${color.black100}` : '',
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
                placeholder={`Search across ${activeBounty} Bounties`}
                value={ui.searchText}
                style={{
                  width: 298,
                  height: 48,
                  background: color.grayish.G600,
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
                TextColor={color.grayish.G100}
                TextColorHover={color.grayish.G50}
                border={`1px solid ${color.grayish.G600}`}
                borderHover={`1px solid ${color.grayish.G400}`}
                borderActive={`1px solid ${color.light_blue100}`}
                iconColor={color.grayish.G300}
                iconColorHover={color.grayish.G50}
              />

              <EuiPopover
                button={
                  <FilterContainer onClick={onButtonClick} color={color}>
                    <div className="filterImageContainer">
                      <MaterialIcon
                        className="materialIconImage"
                        icon="tune"
                        style={{
                          color: isPopoverOpen ? color.grayish.G10 : ''
                        }}
                      />
                    </div>
                    <EuiText
                      className="filterText"
                      style={{
                        color: isPopoverOpen ? color.grayish.G10 : ''
                      }}
                    >
                      Filter
                    </EuiText>
                  </FilterContainer>
                }
                panelStyle={{
                  border: 'none',
                  boxShadow: `0px 1px 20px ${color.black90}`,
                  background: `${color.pureWhite}`,
                  borderRadius: '6px',
                  minWidth: '432px',
                  minHeight: '304px',
                  marginTop: '0px',
                  marginLeft: '20px'
                }}
                isOpen={isPopoverOpen}
                closePopover={closePopover}
                panelClassName="yourClassNameHere"
                panelPaddingSize="none"
                anchorPosition="downLeft"
              >
                <div
                  style={{
                    display: 'flex',
                    flexDirection: 'row'
                  }}
                >
                  <EuiPopOverCheckboxLeft className="CheckboxOuter" color={color}>
                    <EuiText className="leftBoxHeading">STATUS</EuiText>
                    <EuiCheckboxGroup
                      options={Status}
                      idToSelectedMap={checkboxIdToSelectedMap}
                      onChange={(id) => {
                        onChangeStatus(id);
                      }}
                    />
                  </EuiPopOverCheckboxLeft>
                  <PopOverRightBox color={color}>
                    <EuiText className="rightBoxHeading">Languages</EuiText>
                    <EuiPopOverCheckboxRight className="CheckboxOuter" color={color}>
                      <EuiCheckboxGroup
                        options={Coding_Languages}
                        idToSelectedMap={checkboxIdToSelectedMapLanguage}
                        onChange={(id) => {
                          onChangeLanguage(id);
                        }}
                      />
                    </EuiPopOverCheckboxRight>
                  </PopOverRightBox>
                </div>
              </EuiPopover>
            </B>
            <D color={color}>
              <EuiText className="DText" color={color.grayish.G200}>
                Developers
              </EuiText>
              <div className="ImageOuterContainer">
                {peopleList &&
                  peopleList?.slice(0, 3).map((val, index) => {
                    return (
                      <DevelopersImageContainer
                        color={color}
                        style={{
                          zIndex: 3 - index,
                          marginLeft: index > 0 ? '-14px' : '',
                          objectFit: 'cover'
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
                  color: color.black400
                }}
              >
                {developerCount}
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
              borderActive={`1px solid ${color.light_blue100}`}
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
                color: color.pureWhite,
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
                color: color.grayish.G200,
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
                      color={color}
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
                  color: color.black400
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

interface styledProps {
  color?: any;
}

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

const D = styled.div<styledProps>`
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
    color: ${(p) => p.color && p.color.grayish.G200};
    padding: 0 10px;
  }
`;

const DevelopersImageContainer = styled.div<styledProps>`
  height: 28px;
  width: 28px;
  border-radius: 50%;
  background: ${(p) => p.color && p.color.pureWhite};
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

const FilterContainer = styled.div<styledProps>`
  width: 78px;
  height: 48px;
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  margin-left: 19px;
  cursor: pointer;
  user-select: none;
  .filterImageContainer {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 48px;
    width: 36px;
    .materialIconImage {
      color: ${(p) => p.color && p.color.grayish.G200};
      cursor: pointer;
      font-size: 18px;
      margin-top: 4px;
    }
  }
  .filterText {
    font-family: Barlow;
    font-style: normal;
    font-weight: 500;
    font-size: 16px;
    line-height: 19px;
    display: flex;
    align-items: center;
    color: ${(p) => p.color && p.color.grayish.G200};
  }
  &:hover {
    .filterImageContainer {
      .materialIconImage {
        color: ${(p) => p.color && p.color.grayish.G50} !important;
        cursor: pointer;
        font-size: 18px;
        margin-top: 4px;
      }
    }
    .filterText {
      color: ${(p) => p.color && p.color.grayish.G50};
    }
  }
  &:active {
    .filterImageContainer {
      .materialIconImage {
        color: ${(p) => p.color && p.color.grayish.G10} !important;
        cursor: pointer;
        font-size: 18px;
        margin-top: 4px;
      }
    }
    .filterText {
      color: ${(p) => p.color && p.color.grayish.G10};
    }
  }
`;

const EuiPopOverCheckboxLeft = styled.div<styledProps>`
  width: 147px;
  height: 304px;
  padding: 15px 18px;
  border-right: 1px solid ${(p) => p.color && p.color.grayish.G700};
  user-select: none;
  .leftBoxHeading {
    font-family: Barlow;
    font-style: normal;
    font-weight: 700;
    font-size: 12px;
    line-height: 32px;
    text-transform: uppercase;
    color: ${(p) => p.color && p.color.grayish.G100};
    margin-bottom: 10px;
  }

  &.CheckboxOuter > div {
    display: flex;
    flex-direction: column;
    flex-wrap: wrap;

    .euiCheckboxGroup__item {
      .euiCheckbox__square {
        top: 5px;
        border: 1px solid ${(p) => p?.color && p?.color?.grayish.G500};
        border-radius: 2px;
      }
      .euiCheckbox__input + .euiCheckbox__square {
        background: ${(p) => p?.color && p?.color?.pureWhite} no-repeat center;
      }
      .euiCheckbox__input:checked + .euiCheckbox__square {
        border: 1px solid ${(p) => p?.color && p?.color?.blue1};
        background: ${(p) => p?.color && p?.color?.blue1} no-repeat center;
        background-image: url('static/checkboxImage.svg');
      }
      .euiCheckbox__label {
        font-family: 'Barlow';
        font-style: normal;
        font-weight: 500;
        font-size: 13px;
        line-height: 16px;
        color: ${(p) => p?.color && p?.color?.grayish.G50};
        &:hover {
          color: ${(p) => p?.color && p?.color?.grayish.G05};
        }
      }
      input.euiCheckbox__input:checked ~ label {
        color: ${(p) => p?.color && p?.color?.blue1};
      }
    }
  }
`;

const PopOverRightBox = styled.div<styledProps>`
  display: flex;
  flex-direction: column;
  max-height: 304px;
  padding: 15px 0px 20px 21px;
  .rightBoxHeading {
    font-family: Barlow;
    font-style: normal;
    font-weight: 700;
    font-size: 12px;
    line-height: 32px;
    text-transform: uppercase;
    color: ${(p) => p.color && p.color.grayish.G100};
  }
`;

const EuiPopOverCheckboxRight = styled.div<styledProps>`
  min-width: 285px;
  max-width: 285px;
  height: 240px;
  user-select: none;

  &.CheckboxOuter > div {
    height: 100%;
    display: flex;
    flex-direction: column;
    flex-wrap: wrap;
    justify-content: center;
    .euiCheckboxGroup__item {
      .euiCheckbox__square {
        top: 5px;
        border: 1px solid ${(p) => p?.color && p?.color?.grayish.G500};
        border-radius: 2px;
      }
      .euiCheckbox__input + .euiCheckbox__square {
        background: ${(p) => p?.color && p?.color?.pureWhite} no-repeat center;
      }
      .euiCheckbox__input:checked + .euiCheckbox__square {
        border: 1px solid ${(p) => p?.color && p?.color?.blue1};
        background: ${(p) => p?.color && p?.color?.blue1} no-repeat center;
        background-image: url('static/checkboxImage.svg');
      }
      .euiCheckbox__label {
        font-family: 'Barlow';
        font-style: normal;
        font-weight: 500;
        font-size: 13px;
        line-height: 16px;
        color: ${(p) => p?.color && p?.color?.grayish.G50};
        &:hover {
          color: ${(p) => p?.color && p?.color?.grayish.G05};
        }
      }
      input.euiCheckbox__input:checked ~ label {
        color: ${(p) => p?.color && p?.color?.blue1};
      }
    }
  }
`;
