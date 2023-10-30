import React from 'react';
import styled from 'styled-components';
import { PaymentHistoryUserInfo } from './interface';

const UserInfoWrapper = styled.div`
  display: flex;
  align-items: center;
`;

const Wrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  margin-right: 0.9rem;
`;

const DetailWrapper = styled.div`
  display: flex;
  flex-direction: column;
`;

const Image = styled.img`
  height: 2rem;
  width: 2rem;
  border-radius: 50%;
`;

const Name = styled.p`
  color: #3c3f41;
  font-family: Barlow;
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 500;
  line-height: 1rem;
  margin-bottom: 0;
  text-transform: capitalize;
`;

const Pubkey = styled.p`
  overflow: hidden;
  color: #8e969c;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: Barlow;
  font-size: 0.6875rem;
  font-style: normal;
  font-weight: 400;
  line-height: normal;
  letter-spacing: 0.01375rem;
  margin-bottom: 0;
`;
const UserInfo = (props: PaymentHistoryUserInfo) => (
    <UserInfoWrapper>
      <Wrapper>
        <Image src={props.image} />
      </Wrapper>
      <DetailWrapper>
        <Name>{props.name}</Name>
        <Pubkey>{props.pubkey.substring(0, 17)}...</Pubkey>
      </DetailWrapper>
    </UserInfoWrapper>
  );

export default UserInfo;
