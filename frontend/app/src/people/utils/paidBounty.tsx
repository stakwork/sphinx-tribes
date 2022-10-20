import { EuiText } from '@elastic/eui';
import React from 'react';
import styled from 'styled-components';

const PaidBounty = (props) => {
  return (
    <BountyContainer>
      {/* left part */}
      <BountyDescription>
        <Header>
          <OwnerImage>
            <img
              height={'100%'}
              width={'100%'}
              src={props.img || `/static/person_placeholder.png`}
              alt=""
            />
          </OwnerImage>
          <div
            style={{
              display: 'flex',
              flexDirection: 'column'
            }}>
            <EuiText>{props.owner_alias}</EuiText>
            <EuiText>{props.lastSeen}</EuiText>
          </div>
        </Header>
        <Description>
          <EuiText>Description</EuiText>
        </Description>
        <LanguageContainer>
          {props.codingLanguage.map((x) => {
            return <div>{x.label}</div>;
          })}
        </LanguageContainer>
      </BountyDescription>
      {/* right part */}
      <UserPriceContainer>
        <PriceContainer>
          <div>
            <span>$@</span>
            <span>25000SAT</span>
            <span>110.96 USD</span>
          </div>
          <div>{'Est: > 4 hrs'}</div>
        </PriceContainer>
        <UserProfileContainer>
          <UserImage>
            <img width={'100%'} height={'100%'} src={''} alt={''} />
          </UserImage>
          <UserInfo>
            <EuiText>Completed</EuiText>
            <EuiText>Name</EuiText>
            <EuiText>{'View Profile(button)'}</EuiText>
          </UserInfo>
        </UserProfileContainer>
      </UserPriceContainer>
    </BountyContainer>
  );
};

export default PaidBounty;

const BountyContainer = styled.div`
  display: flex;
  flex-direction: row;
  max-width: 1100px;
`;

const BountyDescription = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  min-width: 519px;
  max-width: 519px;
`;

const Header = styled.div`
  display: flex;
  flex-direction: row;
  align-item: center;
`;

const OwnerImage = styled.div`
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
`;

const Description = styled.div`
  display: flex;
  flex-direction: row;
  align-item: center;
`;

const LanguageContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  width: 100%;
`;

const UserPriceContainer = styled.div`
  display: flex;
  flex-direction: row;
  border: 2px solid #86d9b9;
  border-radius: 10px;
`;

const PriceContainer = styled.div`
  max-width: 160px;
  min-width: 160px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: flex-start;
`;

const UserProfileContainer = styled.div`
  min-width: 336px;
  max-width: 336px;
  display: flex;
  justify-content: center;
  align-items: center;
`;

const UserImage = styled.div`
  width: 80px;
  height: 80px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 50%;
`;

const UserInfo = styled.div`
  display: flex;
  flex-direction: column;
`;
