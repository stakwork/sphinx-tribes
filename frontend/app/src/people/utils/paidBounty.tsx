import React from 'react';
import styled from 'styled-components';
import BountyDescription from '../../sphinxUI/bounty_description';
import BountyPrice from '../../sphinxUI/bounty_price';
import BountyProfileView from '../../sphinxUI/bounty_profile_view';

const PaidBounty = (props) => {
  return (
    <BountyContainer>
      {/* left part */}
      <BountyDescription
        {...props}
        title={props.title}
        codingLanguage={props.codingLanguage}
        isPaid={true}
      />
      {/* right part */}
      <PriceUserContainer>
        <BountyPrice
          priceMin={props.priceMin}
          priceMax={props.priceMax}
          price={props.price}
          sessionLength={props.sessionLength}
          style={{
            borderRight: '1px solid rgba(73, 201, 152, 0.2)',
            maxWidth: '245px',
            minWidth: '245px'
          }}
        />
        <BountyProfileView
          assignee={props.assignee}
          status={'COMPLETED'}
          statusStyle={{
            width: '63px',
            height: '16px',
            background: '#8256D0'
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

const BountyContainer = styled.div`
  display: flex;
  flex-direction: row;
  width: 100%;
  font-family: Barlow;
  height: 100% !important;
  background: #fff;
`;

const PriceUserContainer = styled.div`
  display: flex;
  flex-direction: row;
  border: 2px solid #86d9b9;
  border-radius: 10px;
  width: 581px;
`;
