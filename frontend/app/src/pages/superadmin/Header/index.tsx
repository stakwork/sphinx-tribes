import styled from "styled-components"
import React from "react";
import signout from "./icons/signout.svg"
import arrowback from "./icons/arrowback.svg"
import arrowforward from "./icons/arrowforward.svg"

export const Header = () => {

  const NavWrapper = styled.div`
    display: flex;
    padding: 13px 37px 0px 37px;
    justify-content: space-between;
    align-items: flex-start;
    border-bottom: 1px solid var(--Divider-2, #DDE1E5);
    background: var(--Body, #FFF);

  `
  const AlternateWrapper = styled.div`
  background: #FFF;
  box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.15);
  display: flex;
  height: 72px;
  padding-left: 1em;
  padding-right: 2em;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;

`

  const ButtonWrapper = styled.div`
  padding-left:1em;
  padding-right:1em;
  display:flex;
  gap:4px;

`
  
  const Title = styled.h4`
    color: var(--Text-2, var(--Hover-Icon-Color, #3C3F41));
    leading-trim: both;
    text-edge: cap;
    font-family: Barlow;
    font-size: 24px;
    font-style: normal;
    font-weight: 900;
    line-height: 14px; /* 58.333% */
    display:flex;
    gap:6px;
  `
  const Button = styled.h5`
    color: var(--Text-2, var(--Hover-Icon-Color, #3C3F41));
    text-align: center;
    font-family: Barlow;
    font-size: 14px;
    font-style: normal;
    font-weight: 600;
    line-height: normal;
    cursor:pointer;
  `

  const AlternateTitle = styled.h4`
      color: var(--Text-2, var(--Hover-Icon-Color, #3C3F41));
      leading-trim: both;
      text-edge: cap;
      font-family: Barlow;
      font-size: 24px;
      font-style: normal;
      font-weight: 400;
      line-height: 14px;
`
const ExportButton =styled.button`
  border-radius: 6px;
  border: 1px solid var(--Input-Outline-1, #D0D5D8);
  background: var(--White, #FFF);
  box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.06);
  height:2em;
  
`

const Month = styled.h4`
padding-top:8px ;
font-size: 18px;
font-weight:400;
` 
 

  const ArrowButton =styled.button`
    border-radius: 6px;
    border: 1px solid var(--Input-Outline-1, #D0D5D8);
    background: var(--White, #FFF);
    box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.06);
    width:3em;
    height:2em;
  `
  const DropDown = styled.select`
      border:none;
      border-radius: 6px;
      margin-left:12px;
      background: var(--Primary-blue, #618AFF);
      box-shadow: 0px 2px 10px 0px rgba(97, 138, 255, 0.50);
      width:6em;
      color:white;
      height:2em;
      outline:none;
      padding: 0px 10px;
  `
  

  const DropDownOption = styled.option`
    color:white;
    width:3em;
    height:2em;
`  

  const Flex = styled.div`
    display:flex;
  `
  // const daysOfWeek = ['1 Day', '2 Days ', '3 Days ', '4 Days ', '5 Days ', '6 Days ', '7 Days '];
  // const [isOpen, setIsOpen] = useState(false);

  // const handleButtonClick = () => {
  //   setIsOpen(!isOpen);
  // };

  // const handleItemClick = (day:any) => {
  //   console.log(`Selected day: ${day}`);
  //   setIsOpen(false);
  // };

  return (
    <>
  
          <NavWrapper>
            <Title>
              Bounties
              <AlternateTitle>Super Admin</AlternateTitle>
            </Title>
            <Button>
              Sign out 
              <img src={signout} alt="Sign Out" />

            </Button>
          </NavWrapper>
          <AlternateWrapper>
            <Flex>
              <ButtonWrapper>
                <ArrowButton>
                  <img src={arrowback} alt=""  />
                </ArrowButton>
                <ArrowButton><img src={arrowforward} alt=""  /></ArrowButton>
              </ButtonWrapper>
              <Month>01 Oct-31 Dec 2023</Month>
            </Flex>
            <div>
              <ExportButton>
                Export CSV
              </ExportButton>
              <DropDown>
                <DropDownOption>7 Days</DropDownOption>
                <DropDownOption>30 Days</DropDownOption>
                <DropDownOption>45 Days</DropDownOption>
              </DropDown>
            </div>
          </AlternateWrapper>
      
    </>
  )
}

