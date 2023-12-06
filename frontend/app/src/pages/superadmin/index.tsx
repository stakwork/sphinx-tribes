import React from "react";
import styled from 'styled-components';
import { MyTable } from "./TableComponent";
import { bounties } from "./TableComponent/mockBountyData";
import { Header } from "./Header";
import copy from "./Header/icons/copy.svg"
import copygray from "./Header/icons/copygray.svg"
import bountiesposted from "./Header/icons/bountiesposted.svg"
import bountiespaid from "./Header/icons/bountiespaid.svg"
import clockloader from "./Header/icons/clockloder.svg"
import coin from "./Header/icons/coin.svg"
import satoshiesposted from "./Header/icons/satoshiesposted.svg"
import satoshiespaid from "./Header/icons/satoshiespaid.svg"
import calendar from "./Header/icons/calendar.svg"

const Wrapper = styled.section`
  padding: 2em;
  background: var(--Search-bar-background, #F2F3F5);
`;

const Card = styled.div`
  background-color:white;
  margin-top: 2em;
  box-shadow: 0px 0px 10px 0px rgba(219, 219, 219, 0.75);
  padding: 1em;
  position: relative;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1em; /* Adjust the gap as needed */
`;

const VerticaBluelLine = styled.div`
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 6px;
  background-color: #007bff; /* Choose your desired color */
  content: ""; /* Add content property for pseudo-element */
`;


const VerticaGreenlLine = styled.div`
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 6px;
  background-color: #49c998; /* Choose your desired color */
  content: ""; /* Add content property for pseudo-element */
`;

const VerticaGrayLine = styled.div`
  position: absolute;
  left: -25px;
  top: -10px;
  bottom: 0;
  width: 2px;
  height:73px;
  background-color: #edecec; /* Choose your desired color */
  content: ""; /* Add content property for pseudo-element */
`;


const DivWrapper = styled.div`
 position:relative;
 display:flex;

  gap:10px;
`

const LeadingText = styled.h2`
  color: var(--Text-2, var(--Hover-Icon-Color, #3C3F41));
  leading-trim: both; 
  text-edge: cap;
  font-family: Barlow;
  font-size: 20px;
  font-style: normal;
  font-weight: 700;
  line-height: normal;
`;

const Title = styled.h2`
  color: var(--Text-2, var(--Hover-Icon-Color, #3C3F41));
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 24px;
  font-style: normal;
  font-weight: 700;
  line-height: normal;
`;

const TitleBlue = styled.h2`
  color: var(--Primary-Blue-Border, #5078F2);
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 24px;
  font-style: normal;
  font-weight: 700;
  line-height: normal;
`;

const TitleGreen = styled.h2`
  color: var(--Green-Border, #2FB379);
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 24px;
  font-style: normal;
  font-weight: 700;
  line-height: normal;
`;

const Subheading = styled.h3`
  font-size: 0.75em;
  text-align: left;
  color:#bfc5ce;
  font-weight:bold;
`;

const TitleWrapper = styled.div`
  padding-top:0.75em;
  padding-left:5.5em;
   display:flex;
`;

export const SuperAdmin = () => {
    
  console.log("super admin");
  return (
    <>
    <Header/>
    <Wrapper>
      <Card>
      <VerticaBluelLine />
        <TitleWrapper>
          <img src={copy} alt="" width="25px" height="25px"/>
          <LeadingText>Bounties</LeadingText>
        </TitleWrapper>
     
        <DivWrapper>
          <VerticaGrayLine />
          <img src={bountiesposted} alt="" width="35px" height="35px" />
              <div>
                  <Title>
                      200
                  </Title>
                <Subheading>Total Bounties Posted</Subheading>
              </div>
        </DivWrapper>
        <DivWrapper>
          <VerticaGrayLine />
          <img src={copygray} alt="" width="20px" height="35px"/>
          <div>
            <Title>78</Title>
            <Subheading>Bounties assigned</Subheading>
          </div>
        </DivWrapper>
        <DivWrapper>
          <VerticaGrayLine />
          <img src={bountiespaid} alt="" width="24px" height="35px" />
          <div>
          <Title>136</Title>
          <Subheading>Bounties Paid</Subheading>
          </div>
        </DivWrapper>
        <DivWrapper>
        <VerticaGrayLine />
        <img src={clockloader} alt="" width="25px" height="35px" />
        <div>
            <TitleBlue>100%</TitleBlue>
            <Subheading>Bounties Paid</Subheading>
        </div>
        </DivWrapper>
      </Card>
      <Card>
      <VerticaGreenlLine />
        <TitleWrapper>
        <img src={coin} alt="" width="25px" height="25px"/>
          <LeadingText>Satoshies</LeadingText>
        </TitleWrapper>
        <DivWrapper>

          <VerticaGrayLine />
          <img src={satoshiesposted} alt="" width="30px" height="35px"/>
          <div>
         
             <Title>22536</Title>
          <Subheading>Total Sats Posted</Subheading>
          </div>
        </DivWrapper>
        <DivWrapper>
          <VerticaGrayLine />
          <img src={satoshiespaid} alt="" width="25px" height="25px"/>
          <div>
            <Title>13625</Title>
            <Subheading>Sats Paid</Subheading>
          </div>
        </DivWrapper>
        <DivWrapper>
          <VerticaGrayLine />
          <img src={calendar} alt="" width="25px" height="25px"/>
          <div>
          <Title>3 Days</Title>
          <Subheading>Avg Time to Paid</Subheading>
          </div>
        </DivWrapper>
        <DivWrapper>
          <VerticaGrayLine />
          <img src={clockloader} alt="" width="25px" height="25px"/>
          <div>
          <TitleGreen>48%</TitleGreen>
          <Subheading>Paid</Subheading>
          </div>
        </DivWrapper>
      </Card>
    </Wrapper>
 
    <MyTable bounties={bounties} selectedButtonIndex={2} />
   
    </>
  );
};
