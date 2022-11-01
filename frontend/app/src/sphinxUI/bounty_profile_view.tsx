import { EuiText } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';
import React from 'react';
import styled from 'styled-components';

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
              ...props.statusStyle
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
              {{ ...props.assignee }.owner_alias || 'Guest Developer  '}
            </EuiText>
          </NameContainer>

          <div
            style={{
              height: '20px',
              width: '92px',
              left: '909px',
              top: '96px',
              display: 'flex',
              flexDirection: 'row'
            }}
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
            }}>
            <EuiText
              style={{
                fontFamily: 'Barlow',
                fontStyle: 'normal',
                fontWeight: '500',
                fontSize: '13px',
                lineHeight: '16px',
                display: 'flex',
                alignItems: 'center',
                color: '#B0B7BC'
              }}>
              View Profile
            </EuiText>
            <div
              style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                height: '20px',
                width: '24px'
              }}>
              <MaterialIcon
                icon={'arrow_forward'}
                style={{
                  color: '#B0B7BC',
                  fontStyle: 'normal',
                  fontWeight: '400',
                  fontSize: '12px',
                  lineHeight: '12px',
                  display: 'flex',
                  alignItems: 'center',
                  textAlign: 'center',
                  letterSpacing: '0.01em'
                }}
              />
            </div>
          </div>
        </UserInfo>
      </UserProfileContainer>
    </>
  );
};

export default BountyProfileView;

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
  border-radius: 200px;
  overflow: hidden;
`;

const UserInfo = styled.div`
  display: flex;
  flex-direction: column;
  margin-left: 28px;
  margin-top: 3px;
`;

const Status = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 2px;
  .statusText {
    color: #ffffff;
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 700;
    font-size: 8px;
    line-height: 10px;
    display: flex;
    align-items: center;
    text-align: center;
    letter-spacing: 0.1em;
    text-transform: uppercase;
  }
`;

const NameContainer = styled.div`
  width: 100%;
  height: 32px;
  font-family: 'Barlow';
  font-style: normal;
  font-weight: 600;
  font-size: 17px;
  line-height: 23px;
  display: flex;
  align-items: center;
  color: #3c3f41;
  margin-top: 3px;
  margin-bottom: 1px;
`;
