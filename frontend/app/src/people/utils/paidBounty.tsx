import { EuiText } from '@elastic/eui';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import NameTag from './nameTag';

const PaidBounty = (props) => {
  const [codingLabels, setCodingLabels] = useState([]);

  useEffect(() => {
    if (props.codingLanguage && props.codingLanguage.length > 0) {
      const values = props.codingLanguage.map((value) => ({ ...value }));
      setCodingLabels(values);
    }
  }, [props.codingLanguage]);

  return (
    <BountyContainer>
      {/* left part */}
      <BountyDescription>
        <Header>
          <div
            style={{
              display: 'flex',
              flexDirection: 'column'
            }}>
            <NameTag {...props} />
          </div>
        </Header>
        <Description>
          <div
            style={{
              width: '481px',
              height: '64px',
              padding: '8px 0px'
            }}>
            <EuiText
              style={{
                fontSize: '17px',
                lineHeight: '20px'
              }}>
              {props.title}
            </EuiText>
          </div>
        </Description>
        <LanguageContainer>
          {props.codingLanguage?.length > 0 &&
            props.codingLanguage?.map((lang, index) => {
              return (
                <CodingLabels key={index}>
                  <EuiText
                    style={{
                      fontSize: '13px',
                      fontWeight: '500',
                      textAlign: 'center'
                    }}>
                    {lang.label}
                  </EuiText>
                </CodingLabels>
              );
            })}
        </LanguageContainer>
      </BountyDescription>
      {/* right part */}
      <PriceUserContainer>
        <PriceContainer>
          <div style={{}}>
            <span
              style={{
                fontSize: '14px',
                fontWeight: '400',
                lineHeight: '16.8px'
              }}>
              $@
            </span>
            <span
              style={{
                fontSize: '17px',
                fontWeight: '700'
              }}>
              25000SAT
            </span>
            <span>110.96 USD</span>
          </div>
          <div>{'Est: > 4 hrs'}</div>
        </PriceContainer>
        <UserProfileContainer>
          <UserImage>
            <img
              width={'100%'}
              height={'100%'}
              src={
                {
                  ...props.assignee
                }.img
              }
              alt={''}
            />
          </UserImage>
          <UserInfo>
            <EuiText>Completed</EuiText>
            <EuiText>Name</EuiText>
            <EuiText>{'View Profile(button)'}</EuiText>
          </UserInfo>
        </UserProfileContainer>
      </PriceUserContainer>
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
  padding: 16px;
`;

const Header = styled.div`
  display: flex;
  flex-direction: row;
  align-item: center;
  height: 32px;
`;

const OwnerImage = styled.div`
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
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

const PriceUserContainer = styled.div`
  display: flex;
  flex-direction: row;
  border: 2px solid #86d9b9;
  border-radius: 10px;
  min-width: 581px;
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

const CodingLabels = styled.div`
  padding: 0px 8px;
  border: 1px solid #000;
  border-radius: 4px;
  overflow: hidden;
  max-height: 22px;
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-right: 4px;
`;
