import { EuiText } from '@elastic/eui';
import React, { useState } from 'react';
import styled from 'styled-components';
import { colors } from '../../colors';
import BountyDescription from '../../sphinxUI/bounty_description';
import BountyPrice from '../../sphinxUI/bounty_price';
import BountyProfileView from '../../sphinxUI/bounty_profile_view';
import IconButton from '../../sphinxUI/icon_button';
import StartUpModal from './start_up_modal';
import ConnectCard from '../utils/connectCard';
import { useStores } from '../../store';

const Bounties = (props) => {
  const { assignee, price, sessionLength, priceMin, priceMax, codingLanguage, title, person,onPanelClick } =
    props;

  const color = colors['light'];
  const [openStartUpModel, setOpenStartUpModel] = useState<boolean>(false);
  const closeModal = () => setOpenStartUpModel(false);
  const showModal = () => setOpenStartUpModel(true);
  const [openConnectModal, setConnectModal] = useState<boolean>(false);
  const closeConnectModal = () => setConnectModal(false);
  const showConnectModal = () => setConnectModal(true);

  const { ui } = useStores();
  return (
    <>
      {{ ...assignee }.owner_alias ? (
        <BountyContainer onClick={onPanelClick} assignedBackgroundImage={'url("/static/assigned_bounty_bg.svg")'}>
          <div className="BountyDescriptionContainer">
            <BountyDescription
              {...person}
              {...props}
              title={title}
              codingLanguage={codingLanguage}
            />
          </div>
          <div className="BountyPriceContainer">
            <BountyPrice
              priceMin={priceMin}
              priceMax={priceMax}
              price={price}
              sessionLength={sessionLength}
              style={{
                minWidth: '213px',
                maxWidth: '213px',
                borderRight: `1px solid ${color.primaryColor.P200}`
              }}
            />
            <BountyProfileView
              assignee={assignee}
              status={'ASSIGNED'}
              canViewProfile={true}
              statusStyle={{
                width: '55px',
                height: '16px',
                background: color.statusAssigned
              }}
            />
          </div>
        </BountyContainer>
      ) : (
        <BountyContainer>
          <DescriptionPriceContainer unAssignedBackgroundImage='url("/static/unassigned_bounty_bg.svg")'>
            <div style={{display:'flex',flexDirection:'row'}} onClick={onPanelClick} >
            <BountyDescription
              {...person}
              {...props}
              title={title}
              codingLanguage={codingLanguage}
            />
            <BountyPrice
              priceMin={priceMin}
              priceMax={priceMax}
              price={price}
              sessionLength={sessionLength}
              style={{
                borderLeft: `1px solid ${color.grayish.G700}`,
                maxWidth: '245px',
                minWidth: '245px'
              }}
            />
            </div>
            <UnassignedPersonProfile
              unassigned_border={color.grayish.G300}
              grayish_G200={color.grayish.G200}
            >
              <div className="UnassignedPersonContainer">
                <img src="/static/unassigned_profile.svg" alt="" height={'100%'} width={'100%'} />
              </div>
              <div className="UnassignedPersonalDetailContainer">
                <EuiText className="ProfileText">Do your skills match?</EuiText>
                <IconButton
                  text={'I can help'}
                  endingIcon={'arrow_forward'}
                  width={166}
                  height={48}
                  style={{ marginTop: 20 }}
                  onClick={(e) => {
                    if (ui.meInfo) {
                      showConnectModal();
                      e.stopPropagation();
                    } else {
                      e.stopPropagation();
                      showModal();
                    }
                  }}
                  color="primary"
                  hoverColor={color.button_secondary.hover}
                  activeColor={color.button_secondary.active}
                  shadowColor={color.button_secondary.shadow}
                  iconSize={'16px'}
                  iconStyle={{
                    top: '17px',
                    right: '14px'
                  }}
                  textStyle={{
                    width: '108px',
                    display: 'flex',
                    justifyContent: 'flex-start',
                    fontFamily: 'Barlow'
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
      <ConnectCard
        dismiss={() => closeConnectModal()}
        modalStyle={{ top: -64, height: 'calc(100% + 64px)' }}
        person={person}
        visible={openConnectModal}
      />
    </>
  );
};

export default Bounties;

interface containerProps {
  unAssignedBackgroundImage?: string;
  assignedBackgroundImage?: string;
  unassigned_border?: string;
  grayish_G200?: string;
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
  .BountyDescriptionContainer {
    min-width: 553px;
    max-width: 553px;
  }
  .BountyPriceContainer {
    display: flex;
    flex-direction: row;
    width: 545px;
  }
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

const UnassignedPersonProfile = styled.div<containerProps>`
  min-width: 336px;
  min-height: 160px;
  background-image: url("data:image/svg+xml,%3csvg width='100%25' height='100%25' xmlns='http://www.w3.org/2000/svg'%3e%3crect width='100%25' height='100%25' fill='none' rx='10' ry='10' stroke='%23B0B7BCFF' stroke-width='3' stroke-dasharray='4' stroke-dashoffset='0' stroke-linecap='butt'/%3e%3c/svg%3e");
  border-radius: 10px;
  display: flex;
  padding-top: 32px;
  padding-left: 37px;
  .UnassignedPersonContainer {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 80px;
    width: 80px;
    border-radius: 50%;
    margin-top: 5px;
  }
  .UnassignedPersonalDetailContainer {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-left: 25px;
    margin-bottom: 2px;
  }
  .ProfileText {
    font-size: 15px;
    font-weight: 500;
    font-family: Barlow;
    color: ${(p) => (p.grayish_G200 ? p.grayish_G200 : '')};
    margin-bottom: -13px;
    line-height: 18px;
    display: flex;
    align-items: center;
  }
`;
