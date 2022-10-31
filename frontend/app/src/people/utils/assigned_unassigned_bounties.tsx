import { EuiText } from '@elastic/eui';
import React, { useState } from 'react';
import styled from 'styled-components';
import BountyDescription from '../../sphinxUI/bounty_description';
import BountyPrice from '../../sphinxUI/bounty_price';
import BountyProfileView from '../../sphinxUI/bounty_profile_view';
import IconButton from '../../sphinxUI/icon_button';
import StartUpModal from './start_up_modal';

const Bounties = (props) => {
  const [openStartUpModel, setOpenStartUpModel] = useState<boolean>(false);
  const closeModal = () => setOpenStartUpModel(false);
  const showModal = () => setOpenStartUpModel(true);
  return (
    <>
      {{ ...props.assignee }.owner_alias ? (
        <BountyContainer assignedBackgroundImage={'url("/static/assigned_bounty_bg.svg")'}>
          <div
            style={{
              width: '553px'
            }}>
            <BountyDescription
              {...props}
              title={props.title}
              codingLanguage={props.codingLanguage}
            />
          </div>
          <div
            style={{
              display: 'flex',
              flexDirection: 'row',
              width: '545px'
            }}>
            <BountyPrice
              priceMin={props.priceMin}
              priceMax={props.priceMax}
              price={props.price}
              sessionLength={props.sessionLength}
              style={{
                minWidth: '213px',
                maxWidth: '213px',
                borderRight: '1px solid #49C998'
              }}
            />

            <BountyProfileView
              assignee={props.assignee}
              status={'ASSIGNED'}
              statusCode={'#49C998'}
            />
          </div>
        </BountyContainer>
      ) : (
        <BountyContainer>
          <DescriptionPriceContainer unAssignedBackgroundImage='url("/static/unassigned_bounty_bg.svg")'>
            <BountyDescription
              {...props}
              title={props.title}
              codingLanguage={props.codingLanguage}
            />
            <BountyPrice
              priceMin={props.priceMin}
              priceMax={props.priceMax}
              price={props.price}
              sessionLength={props.sessionLength}
              style={{
                borderLeft: '1px solid #EBEDEF',
                maxWidth: '245px',
                minWidth: '245px'
              }}
            />
            <UnassignedPersonProfile>
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'center',
                  alignItems: 'center',
                  height: '80px',
                  width: '80px',
                  borderRadius: '50%'
                }}>
                <img src="/static/unassigned_profile.svg" alt="" height={'100%'} width={'100%'} />
              </div>
              <div
                style={{
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  marginLeft: '16px'
                }}>
                <EuiText
                  style={{
                    fontSize: '15px',
                    fontWeight: '500',
                    fontFamily: 'Barlow',
                    color: '#909BAA',
                    marginBottom: '-16px'
                  }}>
                  Do your skills match?
                </EuiText>
                <IconButton
                  text={'I can help'}
                  endingIcon={'arrow_forward'}
                  width={166}
                  height={48}
                  style={{ marginTop: 20 }}
                  onClick={(e) => {
                    e.stopPropagation();
                    showModal();
                  }}
                  color="primary"
                  hoverColor={'#5881F8'}
                  activeColor={'#5078F2'}
                  shadowColor={'rgba(97, 138, 255, 0.5)'}
                  iconStyle={{
                    top: '13px',
                    right: '14px'
                  }}
                />
              </div>
            </UnassignedPersonProfile>
          </DescriptionPriceContainer>
        </BountyContainer>
      )}
      {openStartUpModel && (
        <StartUpModal closeModal={closeModal} dataObject={'getWork'} buttonColor={'primary'} />
      )}
    </>
  );
};

export default Bounties;

interface containerProps {
  unAssignedBackgroundImage?: string;
  assignedBackgroundImage?: string;
}

const BountyContainer = styled.div<containerProps>`
  display: flex;
  flex-direction: row;
  width: 1100px !important;
  font-family: Barlow;
  height: 160px;
  background: transparent;
  background: ${(p) => (p.assignedBackgroundImage ? p.assignedBackgroundImage : '')};
  background-repeat: no-repeat;
  background-size: cover;
`;

const DescriptionPriceContainer = styled.div<containerProps>`
  display: flex;
  flex-direction: row;
  width: 758px;
  min-height: 160px !important;
  height: 100%;
  background: ${(p) => (p.unAssignedBackgroundImage ? p.unAssignedBackgroundImage : '')};
  background-repeat: no-repeat;
  background-size: cover;
`;

const UnassignedPersonProfile = styled.div`
  min-width: 336px;
  min-height: 160px;
  border: 1px dashed #b0b7bc;
  border-radius: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
`;
