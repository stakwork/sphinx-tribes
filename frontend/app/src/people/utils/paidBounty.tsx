import React from 'react';
import styled from 'styled-components';
import { colors } from '../../colors';
import BountyDescription from '../../sphinxUI/bounty_description';
import BountyPrice from '../../sphinxUI/bounty_price';
import BountyProfileView from '../../sphinxUI/bounty_profile_view';

const PaidBounty = (props) => {
  const color = colors['light'];
  return (
    <BountyContainer onClick={props.onPanelClick} Bounty_Container_Background={color.pureWhite}>
      <BountyDescription
        {...props}
        title={props.title}
        codingLanguage={props.codingLanguage}
        isPaid={true}
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
          right: '-2.5px'
        }}
        alt={'paid_ribbon'}
      />
    </BountyContainer>
  );
};

export default PaidBounty;

interface PaidBountyProps {
  Price_User_Container_Border?: string;
  Bounty_Container_Background?: string;
}

const BountyContainer = styled.div<PaidBountyProps>`
  display: flex;
  flex-direction: row;
  width: 100%;
  font-family: Barlow;
  height: 100% !important;
  background: ${(p) => p.Bounty_Container_Background};
`;

const PriceUserContainer = styled.div<PaidBountyProps>`
  display: flex;
  flex-direction: row;
  border: 2px solid ${(p) => p.Price_User_Container_Border};
  border-radius: 10px;
  width: 581px;
`;
