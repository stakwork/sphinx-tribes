import React, { useState } from 'react';
import styled from 'styled-components';
import { PostModal } from 'people/widgetViews/postBounty/PostModal';
import { EuiCheckboxGroup } from '@elastic/eui';
import {GetValue, coding_languages, status } from 'people/utils/languageLabelStyle';
import { BountyHeaderProps } from 'people/interfaces';
import { colors } from 'config';
import addBounty from './Icons/addBounty.svg';
import searchIcon from './Icons/searchIcon.svg';
import file from './Icons/file.svg';


const Status = GetValue(status);
const Coding_Languages = GetValue(coding_languages);

interface styledProps {
  color?: any;
}

const color = colors['light'];

const Header = styled.div`
  width: 1366px;
  height: 130px;
  padding: 45px 132px 45px 1089px;
  justify-content: flex-end;
  align-items: center;
  align-self: stretch;
  border-bottom: 1px solid var(--Input-BG-1, #f2f3f5);
  background: #fff;
  margin-left: auto;
  margin-right: auto;
`;

const FillContainer = styled.div`
  width: 100vw;
  align-self: stretch;
  background: #fff;
`;

const Filters = styled.div`
  display: flex;
  width: 1366px;
  padding: 10px 130px;
  justify-content: center;
  align-items: center;
  gap: 198px;
  align-self: stretch;
  background: #fff;
  margin-left: auto;
  margin-right: auto;
`;
const FiltersRight = styled.span`
  display: flex;
  height: 40px;
  padding-right: 122px;
  align-items: flex-start;
  gap: 52px;
  flex: 1 0 0;
  width: 1366px;
`;
const StatusContainer = styled.span`
  padding: 10px 0px;
  align-items: center;
  gap: 4px;
`;
const StatusSelector = styled.select`
  background-color: transparent;
  border: none;
  z-index:10;
`;
const SkillContainer = styled.span`
  padding: 10px 0px;
  align-items: center;
  gap: 4px;
  display:flex;
  position:relative;
`;
const Skill = styled.button`
  border: none;
  background-color: transparent;
`;

const Button = styled.button`
  border-radius: 6px;
  background: var(--Primary-Green, #49c998);
  box-shadow: 0px 2px 10px 0px rgba(73, 201, 152, 0.5);
  border: none;
  display: flex;
  width: 144px;
  height: 40px;
  padding: 8px 16px;
  justify-content: flex-end;
  align-items: center;
  gap: 6px;
  color: var(--White, #fff);
  text-align: center;
  font-family: Barlow;
  font-size: 14px;
  font-style: normal;
  font-weight: 500;
  line-height: 0px; /* 0% */
  letter-spacing: 0.14px;
`;

const Label = styled.label`
  color: var(--Main-bottom-icons, #5f6368);
  font-family: Barlow;
  font-size: 15px;
  font-style: normal;
  font-weight: 500;
  line-height: 17px; /* 113.333% */
  letter-spacing: 0.15px;
`;

const SearchWrapper = styled.div`
  height: 40px;
  padding: 0px 16px;
  align-items: center;
  gap: 10px;
  flex: 1 0 0;
  display: flex;
  position: relative;
`;

const Icon = styled.img`
  position: absolute;
  right: 30px;
`;

const SearchBar = styled.input`
  display: flex;
  height: 40px;
  padding: 0px 16px;
  padding-left: 30px;
  align-items: center;
  gap: 10px;
  flex: 1 0 0;
  border-radius: 6px;
  background: var(--Input-BG-1, #f2f3f5);
  outline: none;
  border: none;
`;

const SoryByContainer = styled.span`
  justify-content: center;
  align-items: center;
  gap: 4px;
`;
const SortBy = styled.select`
  background-color: transparent;
  border: none;
`;
const NumberOfBounties = styled.div`
  height: 23px;
  padding: 1.5px 983.492px 1.5px 10px;
  align-items: center;
  flex-shrink: 0;
  margin: 23px 133px;
  width: 1366px;
  margin-left: auto;
  margin-right: auto;
`;
const BountyNumber = styled.span`
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 11px;
`;
const PrimaryText = styled.p`
  color: var(--Primary-Text-1, var(--Press-Icon-Color, #292c33));
  font-family: Barlow;
  font-size: 15px;
  font-style: normal;
  font-weight: 600;
  line-height: normal;
`;
const SecondaryText = styled.p`
  color: var(--Main-bottom-icons, #5f6368);
  font-family: Barlow;
  font-size: 15px;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
`;
const Img = styled.img`
  padding-bottom: 10px;
`;

const SkillFilter = styled.div`
  width:448px;
  height:228px;
  background-color:white;
  position:absolute;
  top:50px;
  z-index:999;
  /* border-top: 3px solid var(--Primary-blue, #618AFF);
  border-top-height: 20px; */

  ::after {
    content: '';
    position: absolute;
    left:0;
    right: 380px;
    top: 0;
    height: 3px; 
    background:var(--Primary-blue, #618AFF); 
}
`
const InternalContainer = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: flex-start;
  gap: 30px;
  padding:24px 28px 28px 32px;
`

const EuiPopOverCheckboxRight = styled.div<styledProps>`

  height: auto;
  user-select: none;
  &.CheckboxOuter > div {
    height: 100%;
    display: grid;
    grid-template-columns: 1fr 1fr 1fr; 
    column-gap:56px;
    justify-content: center;
    .euiCheckboxGroup__item {
      .euiCheckbox__square {
        top: 5px;
        border: 1px solid ${(p: any) => p?.color && p?.color?.grayish.G500};
        border-radius: 2px;
      }
      .euiCheckbox__input + .euiCheckbox__square {
        background: ${(p: any) => p?.color && p?.color?.pureWhite} no-repeat center;
      }
      .euiCheckbox__input:checked + .euiCheckbox__square {
        border: 1px solid ${(p: any) => p?.color && p?.color?.blue1};
        background: ${(p: any) => p?.color && p?.color?.blue1} no-repeat center;
        background-image: url('static/checkboxImage.svg');
      }
      .euiCheckbox__label {
        font-family: 'Barlow';
        font-style: normal;
        font-weight: 500;
        font-size: 13px;
        line-height: 16px;
        color: ${(p: any) => p?.color && p?.color?.grayish.G50};
        &:hover {
          color: ${(p: any) => p?.color && p?.color?.grayish.G05};
        }
      }
      input.euiCheckbox__input:checked ~ label {
        color: ${(p: any) => p?.color && p?.color?.blue1};
      }
    }
  }
`;


export const OrgHeader = ({
  onChangeLanguage,
  checkboxIdToSelectedMapLanguage
}: BountyHeaderProps) => {

  const [isPostBountyModalOpen, setIsPostBountyModalOpen] = useState(false);
  const [filterClick,setFilterClick] = useState(false);
  const selectedWidget = 'wanted';
  const handlePostBountyClick = () => {
    setIsPostBountyModalOpen(true);
  };
  const handlePostBountyClose = () => {
    setIsPostBountyModalOpen(false);
  };

  const handleClick =()=>{
    console.log("callsd");
    setFilterClick(!filterClick)
  }
  return (
    <>
      <FillContainer>
        <Header>
          <Button onClick={handlePostBountyClick}>
            <img src={addBounty} alt="" />
            Post a Bounty
          </Button>
        </Header>
      </FillContainer>
      <FillContainer>
        <Filters>
          <FiltersRight>
            <StatusContainer>
              <Label htmlFor="statusSelect">Status</Label>
              <StatusSelector id="statusSelect" />
            </StatusContainer>
            <SkillContainer>
              <Label htmlFor="skillSelect">Skill</Label>
              <Skill id="skillSelect"/>
               <button onClick={handleClick}> a</button>
               {filterClick ? 
               <SkillFilter>
                <InternalContainer>
                <EuiPopOverCheckboxRight className="CheckboxOuter" color={color}>
                      <EuiCheckboxGroup
                        options={Coding_Languages}
                        idToSelectedMap={checkboxIdToSelectedMapLanguage}
                        onChange={(id: any) => {
                          onChangeLanguage(id);
                        }}
                      />
                </EuiPopOverCheckboxRight>
                </InternalContainer>
               </SkillFilter> : null}
            </SkillContainer>
            <SearchWrapper>
              <SearchBar placeholder="Search" disabled />
              <Icon src={searchIcon} alt="Search" />
            </SearchWrapper>
          </FiltersRight>
          <SoryByContainer>
            <Label htmlFor="statusSelect">Sort by:Newest First</Label>
            <SortBy id="statusSelect" />
          </SoryByContainer>
        </Filters>
      </FillContainer>
      <NumberOfBounties>
        <BountyNumber>
          <Img src={file} alt="" />
          <PrimaryText>284</PrimaryText>
          <SecondaryText>Bounties</SecondaryText>
        </BountyNumber>
      </NumberOfBounties>
      <PostModal
        widget={selectedWidget}
        isOpen={isPostBountyModalOpen}
        onClose={handlePostBountyClose}
      />
    </>
  );
};
