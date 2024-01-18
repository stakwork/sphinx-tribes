import React, { useState } from 'react';
import styled from 'styled-components';
import { PostModal } from 'people/widgetViews/postBounty/PostModal';
import { SearchBar } from '../../../../components/common/index.tsx';
import { useStores } from '../../../../store';
import { colors } from '../../../../config';
import addBounty from './Icons/addBounty.svg';
import file from './Icons/file.svg';

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
const Status = styled.select`
  background-color: transparent;
  border: none;
`;
const SkillContainer = styled.span`
  padding: 10px 0px;
  align-items: center;
  gap: 4px;
`;
const Skill = styled.select`
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

export const OrgHeader = () => {
  const [isPostBountyModalOpen, setIsPostBountyModalOpen] = useState(false);
  const selectedWidget = 'wanted';
  const handlePostBountyClick = () => {
    setIsPostBountyModalOpen(true);
  };
  const handlePostBountyClose = () => {
    setIsPostBountyModalOpen(false);
  };
  const { main, ui } = useStores();
  const color = colors['light'];

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
              <Status id="statusSelect" />
            </StatusContainer>
            <SkillContainer>
              <Label htmlFor="statusSelect">Skill</Label>
              <Skill id="statusSelect" />
            </SkillContainer>
            <SearchBar
                name="search"
                type="search"
                placeholder="Search"
                value={ui.searchText}
                style={{
                  width: '384px',
                  height: '40px',
                  background: color.grayish.G950,
                  fontFamily: 'Barlow',
                  color: color.text2,
                  gap: '10px',
                  flex: '1 0 0',
                  display: 'flex',
                  position: 'relative',
                  marginLeft: '20px',
                  alignItems: 'center',
                  borderRadius: '6px',
                  border: 'none',
                  outline: 'none',
                }}
                onChange={(e: any) => {
                  ui.setSearchText(e);
                }}
                onKeyUp={(e: any) => {
                  if (e.key === 'Enter' || e.keyCode === 13) {
                    main.getOrgBounties({ page: 1, resetPage: true });
                  }
                }}
                TextColor={color.grayish.G100}
                TextColorHover={color.grayish.G50}
                iconColor={color.grayish.G300}
                iconColorHover={color.grayish.G50}
            />
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
