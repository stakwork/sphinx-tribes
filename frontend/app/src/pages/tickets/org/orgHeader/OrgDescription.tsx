import React from 'react';
import { Organization } from 'store/main';
import styled from 'styled-components';
import addBounty from './Icons/addBounty.svg';
import Globe from './Icons/Globe.svg';
import GithubIcon from './Icons/GithubIcon.svg';

const HeaderContainer = styled.div`
  width: 100%;
  border-bottom: 1px solid var(--Input-BG-1, #f2f3f5);
  background: #fff;
`;

const Header = styled.div`
  max-width: 80%;
  margin: 29px auto;
  display: flex;
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const RightHeader = styled.div`
  max-width: 593px;
  display: flex;
  gap: 46px;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
`;
const Leftheader = styled.div`
  display: flex;
  gap: 20px;
  align-items: center;
  justify-content: center;
  justify-content: space-between;
`;
const OrgLabel = styled.div`
  max-width: 403px;
  text-align: right;
  color: var(--Main-bottom-icons, #5f6368);
  font-family: Barlow;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
`;
const ImageContainer = styled.div`
  width: 72px;
  height: 72px;
`;
const StyledImage = styled.img`
  border-radius: 66px;
`;
const CompanyContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

const CompanyLabel = styled.p`
  color: var(--Text-2, #3c3f41);
  font-family: Barlow;
  font-size: 28px;
  font-style: normal;
  font-weight: 700;
`;

const ButtonContainer = styled.div`
  display: flex;
  flex-direction: row;
  gap: 15px;
`;
const SmallButton = styled.a`
  border-radius: 4px;
  border: 1px solid var(--Divider-2, #dde1e5);
  display: flex;
  width: 86px;
  height: 28px;
  padding: 0px 10px 0px 7px;
  align-items: center;
  gap: 6px;
  color: var(--Main-bottom-icons, #5f6368);
  text-decoration: none;
  cursor: pointer;
  font-family: Barlow;
  font-size: 13px;
  font-style: normal;
  font-weight: 500;
  outline: none;
  &:hover {
    text-decoration: none;
    color: var(--Main-bottom-icons, #5f6368);
  }
  &:focus {
    outline: none;
  }
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

function OrgDescription({
  updateIsPostBountyModalOpen,
  orgData
}: {
  updateIsPostBountyModalOpen: (value: boolean) => void;
  orgData: Organization | undefined;
}) {
  if (!orgData) return null;

  const { name, img, website, github, description } = orgData as Organization;

  const handlePostBountyClick = () => {
    updateIsPostBountyModalOpen(true);
  };

  return (
    <HeaderContainer>
      <Header>
        <Leftheader>
          <ImageContainer>
            <StyledImage src={img} width="72px" height="72px" alt="organization icon" />
          </ImageContainer>
          <CompanyContainer>
            <CompanyLabel>{name || ''}</CompanyLabel>
            <ButtonContainer>
              <SmallButton href={website} target="_blank">
                <img src={Globe} alt="globe-website icon" />
                Website
              </SmallButton>
              <SmallButton href={github} target="_blank">
                {' '}
                <img src={GithubIcon} alt="github icon" />
                Github
              </SmallButton>
            </ButtonContainer>
          </CompanyContainer>
        </Leftheader>
        <RightHeader>
          <OrgLabel>{description}</OrgLabel>
          <Button onClick={handlePostBountyClick}>
            <img src={addBounty} alt="" />
            Post a Bounty
          </Button>
        </RightHeader>
      </Header>
    </HeaderContainer>
  );
}

export default OrgDescription;
