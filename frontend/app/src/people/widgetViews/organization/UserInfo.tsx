import React from 'react';
import styled from 'styled-components';
import { PaymentHistoryUserInfo } from './interface';

const UserInfoWrapper = styled.div`
  display: flex;
  align-items: center;
  position: relative;
  cursor: pointer;

  :hover .tooltipText {
    visibility: visible;
    opacity: 1;
  }
`;

const ToolTipWrapper = styled.div`
  visibility: hidden;
  width: 19rem;
  background-color: #eee;
  color: #1a1a1a;
  overflow-wrap: break-word;
  font-size: 0.75rem;
  text-align: center;
  border-radius: 0.5rem;
  padding: 0.75rem;
  position: absolute;
  box-shadow:
    0px 4px 6px -2px rgba(16, 24, 40, 0.03),
    0px 12px 16px -4px rgba(16, 24, 40, 0.08);
  z-index: 1;
  top: 100%;
  left: 0;
  opacity: 0;
  transition: opacity 0.2s;
`;

const Wrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  margin-right: 0.9rem;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 50%;
`;

const DetailWrapper = styled.div`
  display: flex;
  flex-direction: column;
`;

const Image = styled.img`
  height: 2rem;
  width: 2rem;
  border-radius: 50%;
  object-fit: cover;
`;

const Name = styled.p`
  color: #3c3f41;
  font-family: 'Barlow';
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
  font-family: 'Barlow';
  font-size: 0.6875rem;
  font-style: normal;
  font-weight: 400;
  line-height: normal;
  letter-spacing: 0.01375rem;
  margin-bottom: 0;
`;
const UserInfo = (props: PaymentHistoryUserInfo) => {
  const formatName = (name: string) => {
    if (name.length <= 30) {
      return name;
    }
    return `${name.substring(0, 18)}...`;
  };
  return (
    <UserInfoWrapper>
      <Wrapper>
        <Image src={props.image} />
      </Wrapper>
      <DetailWrapper>
        {props.name.length > 30 ? (
          <ToolTipWrapper className="tooltipText">{props.name}</ToolTipWrapper>
        ) : null}
        <Name>{formatName(props.name)}</Name>
        <Pubkey>{props.pubkey.substring(0, 17)}...</Pubkey>
      </DetailWrapper>
    </UserInfoWrapper>
  );
};

export default UserInfo;
