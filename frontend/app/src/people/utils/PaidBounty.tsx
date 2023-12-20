import React from 'react';
import styled from 'styled-components';
import { PaidBountiesProps } from 'people/interfaces';
import BountyDescription from '../../bounties/BountyDescription';
import BountyPrice from '../../bounties/BountyPrice';
import BountyProfileView from '../../bounties/BountyProfileView';
import { colors } from '../../config/colors';

interface PaidBountyProps {
  Price_User_Container_Border?: string;
  Bounty_Container_Background?: string;
  color?: any;
}

const BountyContainer = styled.div<PaidBountyProps>`
  display: flex;
  flex-direction: row;
  width: 100%;
  font-family: 'Barlow';
  height: 160px !important;
  background: ${(p: any) => p.Bounty_Container_Background};
  border: 2px solid ${(p: any) => p.color.grayish.G950};
  border-radius: 10px;
  :hover {
    border: 2px solid ${(p: any) => p.color && p.color.borderGreen2};
    border-radius: 10px;
  }
`;

const PriceUserContainer = styled.div<PaidBountyProps>`
  display: flex;
  flex-direction: row;
  border: 2px solid ${(p: any) => p.Price_User_Container_Border};
  border-radius: 10px;
  width: 579px;
  margin: -0.5px -1.1px;
`;
const PaidBounty = (props: PaidBountiesProps) => {
  const color = colors['light'];

  return (
    <>
      <BountyContainer
        onClick={props.onPanelClick}
        Bounty_Container_Background={color.pureWhite}
        color={color}
      >
        <BountyDescription
          {...props}
          title={props.title}
          codingLanguage={props.codingLanguage}
          isPaid={true}
          org_img={props.org_img}
        />
        <PriceUserContainer Price_User_Container_Border={color.primaryColor.P400}>
          <BountyPrice
            priceMin={props.priceMin}
            priceMax={props.priceMax}
            price={props.price}
            sessionLength={props.sessionLength}
            style={{
              borderRight: `1px solid ${color.primaryColor.P200}`,
              maxWidth: '245px',
              minWidth: '245px'
            }}
          />
          <BountyProfileView
            assignee={props.assignee}
            status={'COMPLETED'}
            canViewProfile={true}
            statusStyle={{
              width: '63px',
              height: '16px',
              background: color.statusCompleted
            }}
          />
        </PriceUserContainer>
        <img
          src={'/static/paid_ribbon.svg'}
          style={{
            position: 'sticky',
            width: '80px',
            height: '80px',
            right: '-1.5px'
          }}
          alt={'paid_ribbon'}
        />
      </BountyContainer>
    </>
  );
};

export default PaidBounty;
