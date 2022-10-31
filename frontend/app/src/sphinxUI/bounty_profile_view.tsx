import { EuiText } from '@elastic/eui';
import React from 'react';
import styled from 'styled-components';
import IconButton from './icon_button';

const BountyProfileView = (props) => {
  return (
    <>
      <UserProfileContainer>
        <UserImage>
          <img
            width={'100%'}
            height={'100%'}
            src={
              {
                ...props.assignee
              }.img || '/static/default_profile_image.svg'
            }
            alt={''}
          />
        </UserImage>
        <UserInfo>
          <Status
            style={{
              background: props.statusCode
            }}>
            <EuiText className="statusText">{props.status}</EuiText>
          </Status>
          <NameContainer>
            <EuiText
              style={{
                fontSize: '17px',
                fontWeight: '600',
                color: '#3C3F41'
              }}>
              {{ ...props.assignee }.owner_alias || 'dummy'}
            </EuiText>
          </NameContainer>

          <IconButton
            text={'View Profile'}
            endingIcon={'arrow_forward'}
            width={89}
            height={20}
            buttonType={'text'}
            style={{ color: '#83878b', textDecoration: 'none', marginLeft: '-21px' }}
            onClick={(e) => {
              if ({ ...props.assignee }.owner_alias) {
                e.stopPropagation();
                window.open(
                  `/p/${
                    {
                      ...props.assignee
                    }.owner_pubkey
                  }?widget=wanted`,
                  '_blank'
                );
              }
            }}
            textStyle={{
              fontSize: '13px',
              fontWeight: '500',
              color: '#B0B7BC',
              fontFamily: 'Barlow'
            }}
            iconStyle={{
              top: '0px',
              left: '90px',
              color: '#B0B7BC',
              fontSize: '20px'
            }}
            color={''}
          />
        </UserInfo>
      </UserProfileContainer>
    </>
  );
};

export default BountyProfileView;

interface statusProps {
  statusCode?: string;
}

const UserProfileContainer = styled.div`
  min-width: 336px;
  max-width: 336px;
  display: flex;
  padding: 40px 0px 0px 37px;
`;

const UserImage = styled.div`
  width: 80px;
  height: 80px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 50%;
  overflow: hidden;
`;

const UserInfo = styled.div`
  display: flex;
  flex-direction: column;
  margin-left: 28px;
`;

const Status = styled.div`
  width: 63px;
  height: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 2px;
  .statusText {
    font-size: 8px;
    line-height: 9.6px;
    weight: 700;
    color: #fff;
  }
`;

const NameContainer = styled.div`
  width: 100%;
  height: 32px;
  display: flex;
  align-items: center;
`;
