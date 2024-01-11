import React from "react"
import styled from "styled-components"
import addBounty from "./Icons/addBounty.svg"
import searchIcon from "./Icons/searchIcon.svg"
import file from "./Icons/file.svg"


const SuperContainer = styled.div`
  display:flex;
  flex-direction:column;
  background: var(--Search-bar-background, #F2F3F5);
  height: 100vh;
`
const Container =styled.div`
    display:flex;
    flex-direction:column;
    justify-content:center;
    overflow: auto; 
`
const Header = styled.div`
    display: flex;
    height: 130px;
    padding: 45px 133px 45px 1089px;
    justify-content: flex-end;
    align-items: center;
    align-self: stretch;
    border-bottom: 1px solid var(--Input-BG-1, #F2F3F5);
    background: #FFF;
   
`
const Filters = styled.div`
    display: flex;
    padding: 10px 133px;
    justify-content: center;
    align-items: center;
    gap: 198px;
    align-self: stretch;
    background: #FFF;
`
const FiltersRight =styled.span`
    display: flex;
    height: 40px;
    padding-right: 122px;
    align-items: flex-start;
    gap: 52px;
    flex: 1 0 0;
`
const StatusContainer = styled.span`
    padding: 10px 0px;
    align-items: center;
    gap: 4px;
`
const Status = styled.select`
    border:none;
`
const SkillContainer = styled.span`
    padding: 10px 0px;
    align-items: center;
    gap: 4px;
`
const Skill = styled.select`
    border:none;
`

const Button = styled.button`
    border-radius: 6px;
    background: var(--Primary-Green, #49C998);
    box-shadow: 0px 2px 10px 0px rgba(73, 201, 152, 0.50);
    border:none;
    display: flex;
    width: 144px;
    height: 40px;
    padding: 8px 16px;
    justify-content: flex-end;
    align-items: center;
    gap: 6px;
    color: var(--White, #FFF);
    text-align: center;
    font-family: Barlow;
    font-size: 14px;
    font-style: normal;
    font-weight: 500;
    line-height: 0px; /* 0% */
    letter-spacing: 0.14px; 
`;

const Label = styled.label`
    color: var(--Main-bottom-icons, #5F6368);
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
    display:flex;
    position:relative;
`;

const Icon = styled.img`
    position:absolute;
    right:30px;
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
    background: var(--Input-BG-1, #F2F3F5);
    outline: none;
    border: none;
`;

const SoryByContainer = styled.span`
    justify-content: center;
    align-items: center;
    gap: 4px;
`
const SortBy = styled.select`
    border:none;
`
const NumberOfBounties = styled.div`
    display: flex;
    width: 1101px;
    height: 23px;
    padding: 1.5px 983.492px 1.5px 0px;
    align-items: center;
    flex-shrink: 0;
    margin:23px 133px;
`
const BountyNumber = styled.span`
    display: flex;
    justify-content:center;
    align-items: center;
    gap: 11px;
`
const PrimaryText = styled.p`
    color: var(--Primary-Text-1, var(--Press-Icon-Color, #292C33));
    font-family: Barlow;
    font-size: 15px;
    font-style: normal;
    font-weight: 600;
    line-height: normal;
`
const SecondaryText = styled.p`
    color: var(--Main-bottom-icons, #5F6368);
    font-family: Barlow;
    font-size: 15px;
    font-style: normal;
    font-weight: 500;
    line-height: normal;
`
const Img = styled.img`
    padding-bottom:10px;
`;

console.log("render", "rendere");

export const OrgBounty = () => (
      <SuperContainer>
        <Container>
            <Header>
                <Button><img src={addBounty} alt="" />Post a Bounty</Button>
            </Header>
            <Filters>
                <FiltersRight>
                    <StatusContainer>
                        <Label htmlFor="statusSelect">Status</Label>
                        <Status id="statusSelect"/>
                    </StatusContainer>
                    <SkillContainer>
                        <Label htmlFor="statusSelect">Skill</Label>
                        <Skill id="statusSelect"/>
                    </SkillContainer>
                    <SearchWrapper>
                        <SearchBar placeholder="Search" />
                        <Icon src={searchIcon} alt="Search" />
                    </SearchWrapper>
                    
                </FiltersRight>
                <SoryByContainer>
                        <Label htmlFor="statusSelect">Sort by:Newest First</Label>
                        <SortBy id="statusSelect"/>
                </SoryByContainer>
            </Filters>
        </Container>
         <NumberOfBounties>
             <BountyNumber>
                 <Img src={file} alt="" />
                <PrimaryText>284</PrimaryText>
                <SecondaryText>Bounties</SecondaryText>
             </BountyNumber>
         </NumberOfBounties>
      </SuperContainer>
    )