import React from "react";
import copy from "../Header/icons/copy.svg"
import copygray from "../Header/icons/copygray.svg"
import bountiesposted from "../Header/icons/bountiesposted.svg"
import bountiespaid from "../Header/icons/bountiespaid.svg"
import clockloader from "../Header/icons/clockloder.svg"
import coin from "../Header/icons/coin.svg"
import satoshiesposted from "../Header/icons/satoshiesposted.svg"
import satoshiespaid from "../Header/icons/satoshiespaid.svg"
import calendar from "../Header/icons/calendar.svg"
import {
  Wrapper,
  Card,
  VerticaBluelLine,
  VerticaGreenlLine,
  VerticaGrayLine,
  DivWrapper,
  LeadingText,
  Title,
  TitleBlue,
  TitleGreen,
  Subheading,
  TitleWrapper,
} from './StatisticsStyles';

export const Statistics = () => {
    
  console.log("super admin");
  return (
    <>
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
   
    </>
  );
};
