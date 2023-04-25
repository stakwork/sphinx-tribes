/* eslint-disable func-style */
import React from 'react';
import {
  AssigneeProfile,
  Creator,
  Img,
  PaidStatusPopover,
  CreatorDescription,
  BountyPriceContainer,
  BottomButtonContainer,
  UnassignedPersonProfile,
  DividerContainer,
  NormalUser,
  LanguageContainer,
  AwardsContainer,
  DescriptionBox,
  AdjustAmountContainer,
  TitleBox,
  CodingLabels,
  AutoCompleteContainer,
  AwardBottomContainer
} from './style';
import { EuiText, EuiFieldText } from '@elastic/eui';
import { Divider, Modal } from '../../../../components/common';
import { colors } from '../../../../config/colors';
import { renderMarkdown } from '../../../utils/renderMarkdown';
import { satToUsd } from '../../../../helpers';
import { useStores } from '../../../../store';
import IconButton from '../../../../components/common/icon_button';
import ImageButton from '../../../../components/common/Image_button';
import BountyProfileView from '../../../../bounties/bounty_profile_view';
import ButtonSet from '../../../../bounties/bountyModal_button_set';
import BountyPrice from '../../../../bounties/bounty_price';
import InvitePeopleSearch from '../../../../components/form/inputs/widgets/PeopleSearch';
import { observer } from 'mobx-react-lite';

export default observer(MobileView);
function MobileView(props: any) {
  const {
    deliverables,
    description,
    ticketUrl,
    assignee,
    titleString,
    nametag,
    labels,
    person,
    setIsPaidStatusPopOver,
    creatorStep,
    paid,
    tribe,
    saving,
    isPaidStatusPopOver,
    isPaidStatusBadgeInfo,
    awardDetails,
    isAssigned,
    dataValue,
    assigneeValue,
    assignedPerson,
    changeAssignedPerson,
    sendToRedirect,
    handleCopyUrl,
    isCopied,
    setExtrasPropertyAndSave,
    setIsModalSideButton,
    replitLink,
    assigneeHandlerOpen,
    setCreatorStep,
    setIsExtraStyle,
    awards,
    setExtrasPropertyAndSaveMultiple,
    handleAssigneeDetails,
    peopleList,
    setIsPaidStatusBadgeInfo,
    bountyPrice,
    selectedAward,
    handleAwards,
    repo,
    issue,
    isMarkPaidSaved,
    setAwardDetails,
    setBountyPrice,
    owner_idURL,
    createdURL
  } = props;
  const color = colors['light'];

  const { ui } = useStores();

  return (
    <div>
      {{ ...person }?.owner_alias === ui.meInfo?.owner_alias ? (
        /*
         * creator view
         */
        <>
          {creatorStep === 0 && (
            <Creator
              onClick={() => {
                setIsPaidStatusPopOver(false);
              }}
            >
              <>
                {paid && (
                  <Img
                    src={'/static/paid_ribbon.svg'}
                    style={{
                      position: 'absolute',
                      right: -4,
                      width: 72.46,
                      height: 71.82,
                      zIndex: 100,
                      pointerEvents: 'none'
                    }}
                  />
                )}
                {paid && (
                  <>
                    <PaidStatusPopover
                      color={color}
                      isPaidStatusPopOver={isPaidStatusPopOver}
                      isPaidStatusBadgeInfo={isPaidStatusBadgeInfo}
                      style={{
                        opacity: isPaidStatusPopOver ? 1 : 0,
                        transition: 'all ease 1s'
                      }}
                    >
                      <div
                        className="PaidStatusContainer"
                        style={{
                          borderRadius: isPaidStatusBadgeInfo ? '6px 6px 0px 0px' : '6px',
                          opacity: isPaidStatusPopOver ? 1 : 0,
                          transition: 'all ease 1s'
                        }}
                      >
                        <div className="imageContainer">
                          <img
                            src="/static/verified_check_icon.svg"
                            alt="check icon"
                            height={'100%'}
                            width={'100%'}
                          />
                        </div>
                        <EuiText className="PaidStatus">Bounty Paid</EuiText>
                      </div>
                      <div
                        className="ExtraBadgeInfo"
                        style={{
                          opacity: isPaidStatusBadgeInfo ? 1 : 0,
                          transition: 'all ease 1s'
                        }}
                      >
                        <div className="imageContainer">
                          <img
                            src="/static/green_checked_icon.svg"
                            alt=""
                            height={'100%'}
                            width={'100%'}
                          />
                        </div>
                        <img
                          src={awardDetails?.image !== '' && awardDetails.image}
                          alt="award_icon"
                          height={'40px'}
                          width={'40px'}
                        />
                        <EuiText className="badgeText">Badge Awarded</EuiText>
                      </div>
                    </PaidStatusPopover>
                  </>
                )}

                <CreatorDescription paid={paid} color={color}>
                  <div className="CreatorDescriptionOuterContainerCreatorView">
                    <div className="CreatorDescriptionInnerContainerCreatorView">
                      <div>{nametag}</div>
                      <div className="CreatorDescriptionExtraButton">
                        <ImageButton
                          buttonText={'Edit'}
                          ButtonContainerStyle={{
                            width: '117px',
                            height: '40px'
                          }}
                          leadingImageSrc={'/static/editIcon.svg'}
                          leadingImageContainerStyle={{
                            left: 320
                          }}
                          buttonAction={props?.editAction}
                        />
                        <ImageButton
                          buttonText={!props.deletingState ? 'Delete' : 'Deleting'}
                          ButtonContainerStyle={{
                            width: '117px',
                            height: '40px'
                          }}
                          leadingImageSrc={'/static/Delete.svg'}
                          leadingImageContainerStyle={{
                            left: 450
                          }}
                          buttonAction={props?.deleteAction}
                        />
                      </div>
                    </div>
                    <TitleBox color={color}>{titleString}</TitleBox>
                    <LanguageContainer>
                      {dataValue &&
                        dataValue?.length > 0 &&
                        dataValue?.map((lang: any, index) => (
                          <CodingLabels
                            key={index}
                            styledColors={color}
                            border={lang?.border}
                            color={lang?.color}
                            background={lang?.background}
                          >
                            <EuiText className="LanguageText">{lang?.label}</EuiText>
                          </CodingLabels>
                        ))}
                    </LanguageContainer>
                  </div>
                  <DescriptionBox color={color}>
                    {renderMarkdown(description)}
                    {deliverables ? (
                      <div className="deliverablesContainer">
                        <EuiText className="deliverablesHeading">Deliverables</EuiText>
                        <EuiText className="deliverablesDesc">{deliverables}</EuiText>
                      </div>
                    ) : null}
                  </DescriptionBox>
                </CreatorDescription>
                <AssigneeProfile color={color}>
                  <UnassignedPersonProfile
                    unassigned_border={color.grayish.G300}
                    grayish_G200={color.grayish.G200}
                    color={color}
                  >
                    {!isAssigned && (
                      <div className="UnassignedPersonContainer">
                        <img
                          src="/static/unassigned_profile.svg"
                          alt=""
                          height={'100%'}
                          width={'100%'}
                        />
                      </div>
                    )}

                    {isAssigned ? (
                      <div className="BountyProfileOuterContainerCreatorView">
                        <BountyProfileView
                          assignee={!assignedPerson ? assignee : assignedPerson}
                          status={paid ? 'completed' : 'assigned'}
                          canViewProfile={false}
                          statusStyle={{
                            width: '66px',
                            height: '16px',
                            background: paid ? color.statusCompleted : color.statusAssigned
                          }}
                          UserProfileContainerStyle={{
                            height: 48,
                            width: 'fit-content',
                            minWidth: 'fit-content',
                            padding: 0
                            // marginTop: '48px'
                          }}
                          isNameClickable={true}
                          UserImageStyle={{
                            width: '48px',
                            height: '48px',
                            display: 'flex',
                            justifyContent: 'center',
                            alignItems: 'center',
                            borderRadius: '200px',
                            overflow: 'hidden'
                          }}
                          NameContainerStyle={{
                            height: '28px',
                            maxWidth: '154px'
                          }}
                          userInfoStyle={{
                            marginLeft: '12px'
                          }}
                        />
                        {!paid && (
                          <div
                            className="AssigneeCloseButtonContainer"
                            onClick={() => {
                              changeAssignedPerson();
                              setIsModalSideButton(false);
                            }}
                          >
                            <img
                              src="/static/assignee_close.png"
                              alt="cross_icon"
                              height={'100%'}
                              width={'100%'}
                            />
                          </div>
                        )}
                      </div>
                    ) : (
                      <div className="UnassignedPersonalDetailContainer">
                        <ImageButton
                          buttonText={'Not Assigned'}
                          ButtonContainerStyle={{
                            width: '159px',
                            height: '48px',
                            background: color.pureWhite,
                            marginLeft: '-12px'
                          }}
                          buttonTextStyle={{
                            color: color.grayish.G50,
                            width: '114px'
                          }}
                          endImageSrc={'/static/addIcon.svg'}
                          endingImageContainerStyle={{
                            right: '34px',
                            fontSize: '12px'
                          }}
                          buttonAction={() => {
                            assigneeHandlerOpen();
                          }}
                        />
                      </div>
                    )}
                  </UnassignedPersonProfile>
                  <DividerContainer>
                    <Divider />
                  </DividerContainer>
                  <BountyPriceContainer margin_top="0px">
                    <BountyPrice
                      priceMin={props?.priceMin}
                      priceMax={props?.priceMax}
                      price={props?.price}
                      sessionLength={props?.estimate_session_length}
                      style={{
                        padding: 0,
                        margin: 0
                      }}
                    />
                  </BountyPriceContainer>
                  <div className="buttonSet">
                    <ButtonSet
                      githubShareAction={() => {
                        const repoUrl = ticketUrl
                          ? ticketUrl
                          : `https://github.com/${repo}/issues/${issue}`;
                        sendToRedirect(repoUrl);
                      }}
                      copyURLAction={handleCopyUrl}
                      copyStatus={isCopied ? 'Copied' : 'Copy Link'}
                      twitterAction={() => {
                        const twitterLink = `https://twitter.com/intent/tweet?text=Hey, I created a new ticket on Sphinx community.%0A${titleString} %0A&url=https://community.sphinx.chat/p?owner_id=${owner_idURL}%26created${createdURL} %0A%0A&hashtags=${
                          labels && labels.map((x: any) => x.label)
                        },sphinxchat`;
                        sendToRedirect(twitterLink);
                      }}
                      replitLink={replitLink}
                      tribe={tribe !== 'none' && tribe}
                      tribeFunction={() => {
                        const profileUrl = `https://community.sphinx.chat/t/${tribe}`;
                        sendToRedirect(profileUrl);
                      }}
                      showGithubBtn={!!ticketUrl}
                    />
                  </div>
                  <BottomButtonContainer>
                    {paid ? (
                      <IconButton
                        width={220}
                        height={48}
                        style={{
                          bottom: '0',
                          marginLeft: '36px',
                          border: `1px solid ${color.primaryColor.P400}`,
                          background: color.pureWhite,
                          color: color.borderGreen1
                        }}
                        text={'Mark Unpaid'}
                        loading={saving === 'paid'}
                        endingImg={'/static/mark_unpaid.svg'}
                        textStyle={{
                          width: '130px',
                          display: 'flex',
                          justifyContent: 'center',
                          fontFamily: 'Barlow',
                          marginLeft: '30px'
                        }}
                        onClick={(e) => {
                          e.stopPropagation();
                          setExtrasPropertyAndSave('paid', !paid);
                          setIsModalSideButton(true);
                        }}
                      />
                    ) : (
                      <IconButton
                        color={'success'}
                        width={220}
                        height={48}
                        style={{
                          bottom: '0',
                          marginLeft: '36px'
                        }}
                        text={'Mark Paid'}
                        loading={saving === 'paid'}
                        endingImg={'/static/mark_paid.svg'}
                        textStyle={{
                          width: '130px',
                          display: 'flex',
                          justifyContent: 'center',
                          fontFamily: 'Barlow',
                          marginLeft: '30px'
                        }}
                        hovercolor={color.button_primary.hover}
                        activecolor={color.button_primary.active}
                        shadowcolor={color.button_primary.shadow}
                        onClick={(e) => {
                          e.stopPropagation();
                          // setExtrasPropertyAndSave('paid', !paid);
                          setCreatorStep(1);
                          setIsModalSideButton(false);
                          setIsExtraStyle(true);
                        }}
                      />
                    )}
                  </BottomButtonContainer>
                </AssigneeProfile>
              </>
            </Creator>
          )}

          {creatorStep === 1 && (
            <AdjustAmountContainer color={color}>
              <div
                className="TopHeader"
                onClick={() => {
                  setCreatorStep(0);
                  setIsModalSideButton(true);
                  setIsExtraStyle(false);
                }}
              >
                <div className="imageContainer">
                  <img
                    height={'12px'}
                    width={'8px'}
                    src={'/static/back_button_image.svg'}
                    alt={'back_button_icon'}
                  />
                </div>
                <EuiText className="TopHeaderText">Back to Bounty</EuiText>
              </div>
              <div className="Header">
                <EuiText className="HeaderText">Adjust the amount</EuiText>
              </div>
              <div className="AssignedProfile">
                <BountyProfileView
                  assignee={assignee}
                  status={'Assigned'}
                  canViewProfile={false}
                  statusStyle={{
                    width: '66px',
                    height: '16px',
                    background: color.statusAssigned
                  }}
                  isNameClickable={true}
                  UserProfileContainerStyle={{
                    height: 80,
                    width: 235,
                    padding: '0px 0px 0px 33px',
                    marginTop: '48px',
                    marginBottom: '27px'
                  }}
                  UserImageStyle={{
                    width: '80px',
                    height: '80px',
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    borderRadius: '200px',
                    overflow: 'hidden'
                  }}
                  NameContainerStyle={{
                    height: '28px'
                  }}
                  userInfoStyle={{
                    marginLeft: '28px',
                    marginTop: '6px'
                  }}
                />
                <div className="InputContainer">
                  <EuiText className="InputContainerLeadingText">$@</EuiText>
                  <EuiFieldText
                    className="InputContainerTextField"
                    type={'number'}
                    value={bountyPrice}
                    onChange={(e) => {
                      setBountyPrice(e.target.value);
                    }}
                  />
                  <EuiText className="InputContainerEndingText">SAT</EuiText>
                </div>
                <EuiText className="USDText">{satToUsd(bountyPrice)} USD</EuiText>
              </div>
              <div className="BottomButton">
                <IconButton
                  color={'primary'}
                  width={120}
                  height={42}
                  text={'Next'}
                  textStyle={{
                    width: '100%',
                    display: 'flex',
                    justifyContent: 'center',
                    fontFamily: 'Barlow'
                  }}
                  hovercolor={color.button_secondary.hover}
                  activecolor={color.button_secondary.active}
                  shadowcolor={color.button_secondary.shadow}
                  onClick={(e) => {
                    e.stopPropagation();
                    // setExtrasPropertyAndSave('paid', !paid);
                    // setExtrasPropertyAndSave('price', bountyPrice);
                    setCreatorStep(2);
                    setIsExtraStyle(false);
                  }}
                />
              </div>
            </AdjustAmountContainer>
          )}
          {creatorStep === 2 && (
            <AwardsContainer color={color}>
              <div className="header">
                <div
                  className="headerTop"
                  onClick={() => {
                    setCreatorStep(1);
                    setIsExtraStyle(true);
                  }}
                >
                  <div className="imageContainer">
                    <img
                      height={'12px'}
                      width={'8px'}
                      src={'/static/back_button_image.svg'}
                      alt={'back_button_icon'}
                    />
                  </div>
                  <EuiText className="TopHeaderText">Back</EuiText>
                </div>
                <EuiText className="headerText">Award Badge</EuiText>
              </div>
              <div className="AwardContainer">
                {awards?.map((award, index) => (
                  <div
                    className="RadioImageContainer"
                    key={index}
                    style={{
                      border: selectedAward === award.id ? `1px solid ${color.blue2}` : ''
                    }}
                    onClick={() => {
                      handleAwards(award.id);
                      setAwardDetails({
                        name: award.label,
                        image: award.label_icon
                      });
                    }}
                  >
                    <input
                      type="radio"
                      id={award.id}
                      name={'award'}
                      value={award.id}
                      checked={selectedAward === award.id}
                      style={{
                        height: '16px',
                        width: '16px',
                        cursor: 'pointer'
                      }}
                    />
                    <div className="awardImageContainer">
                      <img src={award.label_icon} alt="icon" height={'100%'} width={'100%'} />
                    </div>
                    <EuiText className="awardLabelText">{award.label}</EuiText>
                  </div>
                ))}
              </div>
              <AwardBottomContainer color={color}>
                <IconButton
                  color={'success'}
                  width={220}
                  height={48}
                  style={{
                    bottom: '0',
                    marginLeft: '36px'
                  }}
                  text={selectedAward === '' ? 'Skip and Mark Paid' : 'Mark Paid'}
                  loading={isMarkPaidSaved}
                  endingImg={'/static/mark_paid.svg'}
                  textStyle={{
                    width: '130px',
                    display: 'flex',
                    justifyContent: 'center',
                    fontFamily: 'Barlow',
                    marginLeft: '30px',
                    marginRight: '10px'
                  }}
                  hovercolor={color.button_primary.hover}
                  activecolor={color.button_primary.active}
                  shadowcolor={color.button_primary.shadow}
                  onClick={(e) => {
                    e.stopPropagation();
                    setExtrasPropertyAndSaveMultiple('paid', {
                      paid: !paid,
                      price: bountyPrice,
                      award: awardDetails.name
                    });

                    setTimeout(() => {
                      setCreatorStep(0);
                      setIsModalSideButton(true);
                    }, 3000);
                    setTimeout(() => {
                      setIsPaidStatusPopOver(true);
                    }, 4000);
                    setTimeout(() => {
                      if (awardDetails?.name !== '') {
                        setIsPaidStatusBadgeInfo(true);
                      }
                    }, 5500);
                  }}
                />
              </AwardBottomContainer>
            </AwardsContainer>
          )}

          {assigneeValue && (
            <Modal
              visible={true}
              envStyle={{
                borderRadius: '10px',
                background: color.pureWhite,
                maxHeight: '459px',
                width: '44.5%'
              }}
              bigCloseImage={assigneeHandlerOpen}
              bigCloseImageStyle={{
                top: '-18px',
                right: '-18px',
                background: color.pureBlack,
                borderRadius: '50%',
                zIndex: 11
              }}
            >
              <AutoCompleteContainer color={color}>
                <EuiText className="autoCompleteHeaderText">Invite Developer</EuiText>
                <InvitePeopleSearch
                  peopleList={peopleList}
                  isProvidingHandler={true}
                  handleAssigneeDetails={(value) => {
                    handleAssigneeDetails(value);
                  }}
                />
              </AutoCompleteContainer>
            </Modal>
          )}
        </>
      ) : (
        /*
         * normal user view
         */
        <NormalUser>
          {paid && (
            <Img
              src={'/static/paid_ribbon.svg'}
              style={{
                position: 'absolute',
                right: -4,
                width: 72.46,
                height: 71.82,
                zIndex: 100,
                pointerEvents: 'none'
              }}
            />
          )}
          <CreatorDescription paid={paid} color={color}>
            <div className="DescriptionUpperContainerNormalView">
              <div>{nametag}</div>
              <TitleBox color={color}>{titleString}</TitleBox>
              <LanguageContainer>
                {dataValue &&
                  dataValue?.length > 0 &&
                  dataValue?.map((lang: any, index) => (
                    <CodingLabels
                      key={index}
                      styledColors={color}
                      border={lang?.border}
                      color={lang?.color}
                      background={lang?.background}
                    >
                      <EuiText className="LanguageText">{lang?.label}</EuiText>
                    </CodingLabels>
                  ))}
              </LanguageContainer>
            </div>
            <DescriptionBox color={color}>
              {renderMarkdown(description)}
              {deliverables ? (
                <div className="deliverablesContainer">
                  <EuiText className="deliverablesHeading">Deliverables</EuiText>
                  <EuiText className="deliverablesDesc">{deliverables}</EuiText>
                </div>
              ) : null}
            </DescriptionBox>
          </CreatorDescription>

          <AssigneeProfile color={color}>
            {paid ? (
              <>
                <BountyProfileView
                  assignee={assignee}
                  status={'Completed'}
                  canViewProfile={false}
                  statusStyle={{
                    width: '66px',
                    height: '16px',
                    background: color.statusCompleted
                  }}
                  isNameClickable={true}
                  UserProfileContainerStyle={{
                    height: 48,
                    width: 235,
                    padding: '0px 0px 0px 33px',
                    marginTop: '48px'
                  }}
                  UserImageStyle={{
                    width: '48px',
                    height: '48px',
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    borderRadius: '200px',
                    overflow: 'hidden'
                  }}
                  NameContainerStyle={{
                    height: '28px'
                  }}
                  userInfoStyle={{
                    marginLeft: '12px'
                  }}
                />
                <DividerContainer>
                  <Divider />
                </DividerContainer>
                <BountyPriceContainer margin_top="0px">
                  <BountyPrice
                    priceMin={props?.priceMin}
                    priceMax={props?.priceMax}
                    price={props?.price}
                    sessionLength={props?.estimate_session_length}
                    style={{
                      padding: 0,
                      margin: 0
                    }}
                  />
                </BountyPriceContainer>
                <ButtonSet
                  showGithubBtn={!!ticketUrl}
                  githubShareAction={() => {
                    const repoUrl = ticketUrl
                      ? ticketUrl
                      : `https://github.com/${repo}/issues/${issue}`;
                    sendToRedirect(repoUrl);
                  }}
                  copyURLAction={handleCopyUrl}
                  copyStatus={isCopied ? 'Copied' : 'Copy Link'}
                  twitterAction={() => {
                    const twitterLink = `https://twitter.com/intent/tweet?text=Hey, I created a new ticket on Sphinx community.%0A${titleString} %0A&url=https://community.sphinx.chat/p?owner_id=${owner_idURL}%26created${createdURL} %0A%0A&hashtags=${
                      labels && labels.map((x: any) => x.label)
                    },sphinxchat`;
                    sendToRedirect(twitterLink);
                  }}
                  replitLink={replitLink}
                  tribe={tribe !== 'none' && tribe}
                  tribeFunction={() => {
                    const profileUrl = `https://community.sphinx.chat/t/${tribe}`;
                    sendToRedirect(profileUrl);
                  }}
                />
              </>
            ) : assignee?.owner_alias ? (
              <>
                <BountyProfileView
                  assignee={assignee}
                  status={'ASSIGNED'}
                  canViewProfile={false}
                  statusStyle={{
                    width: '55px',
                    height: '16px',
                    background: color.statusAssigned
                  }}
                  isNameClickable={true}
                  UserProfileContainerStyle={{
                    height: 48,
                    width: 235,
                    padding: '0px 0px 0px 33px',
                    marginTop: '48px'
                  }}
                  UserImageStyle={{
                    width: '48px',
                    height: '48px',
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    borderRadius: '200px',
                    overflow: 'hidden'
                  }}
                  NameContainerStyle={{
                    height: '28px'
                  }}
                  userInfoStyle={{
                    marginLeft: '12px'
                  }}
                />
                <DividerContainer>
                  <Divider />
                </DividerContainer>
                <BountyPriceContainer margin_top="0px">
                  <BountyPrice
                    priceMin={props?.priceMin}
                    priceMax={props?.priceMax}
                    price={props?.price}
                    sessionLength={props?.estimate_session_length}
                    style={{
                      padding: 0,
                      margin: 0
                    }}
                  />
                </BountyPriceContainer>
                <ButtonSet
                  showGithubBtn={!!ticketUrl}
                  githubShareAction={() => {
                    const repoUrl = ticketUrl
                      ? ticketUrl
                      : `https://github.com/${repo}/issues/${issue}`;
                    sendToRedirect(repoUrl);
                  }}
                  copyURLAction={handleCopyUrl}
                  copyStatus={isCopied ? 'Copied' : 'Copy Link'}
                  twitterAction={() => {
                    const twitterLink = `https://twitter.com/intent/tweet?text=Hey, I created a new ticket on Sphinx community.%0A${titleString} %0A&url=https://community.sphinx.chat/p?owner_id=${owner_idURL}%26created${createdURL} %0A%0A&hashtags=${
                      labels && labels.map((x: any) => x.label)
                    },sphinxchat`;
                    sendToRedirect(twitterLink);
                  }}
                  replitLink={replitLink}
                  tribe={tribe !== 'none' && tribe}
                  tribeFunction={() => {
                    const profileUrl = `https://community.sphinx.chat/t/${tribe}`;
                    sendToRedirect(profileUrl);
                  }}
                />
              </>
            ) : (
              <>
                <UnassignedPersonProfile
                  unassigned_border={color.grayish.G300}
                  grayish_G200={color.grayish.G200}
                  color={color}
                >
                  <div className="UnassignedPersonContainer">
                    <img
                      src="/static/unassigned_profile.svg"
                      alt=""
                      height={'100%'}
                      width={'100%'}
                    />
                  </div>
                  <div className="UnassignedPersonalDetailContainer">
                    <IconButton
                      text={'I can help'}
                      endingIcon={'arrow_forward'}
                      width={153}
                      height={48}
                      onClick={props.extraModalFunction}
                      color="primary"
                      hovercolor={color.button_secondary.hover}
                      activecolor={color.button_secondary.active}
                      shadowcolor={color.button_secondary.shadow}
                      iconSize={'16px'}
                      iconStyle={{
                        top: '16px',
                        right: '14px'
                      }}
                      textStyle={{
                        width: '106px',
                        display: 'flex',
                        justifyContent: 'flex-start',
                        fontFamily: 'Barlow'
                      }}
                    />
                  </div>
                </UnassignedPersonProfile>
                <DividerContainer>
                  <Divider />
                </DividerContainer>
                <BountyPriceContainer margin_top="0px">
                  <BountyPrice
                    priceMin={props?.priceMin}
                    priceMax={props?.priceMax}
                    price={props?.price}
                    sessionLength={props?.estimate_session_length}
                    style={{
                      padding: 0,
                      margin: 0
                    }}
                  />
                </BountyPriceContainer>
                <ButtonSet
                  showGithubBtn={!!ticketUrl}
                  githubShareAction={() => {
                    const repoUrl = ticketUrl
                      ? ticketUrl
                      : `https://github.com/${repo}/issues/${issue}`;
                    sendToRedirect(repoUrl);
                  }}
                  copyURLAction={handleCopyUrl}
                  copyStatus={isCopied ? 'Copied' : 'Copy Link'}
                  twitterAction={() => {
                    const twitterLink = `https://twitter.com/intent/tweet?text=Hey, I created a new ticket on Sphinx community.%0A${titleString} %0A&url=https://community.sphinx.chat/p?owner_id=${owner_idURL}%26created${createdURL} %0A%0A&hashtags=${
                      labels && labels.map((x: any) => x.label)
                    },sphinxchat`;
                    sendToRedirect(twitterLink);
                  }}
                  replitLink={replitLink}
                  tribe={tribe !== 'none' && tribe}
                  tribeFunction={() => {
                    const profileUrl = `https://community.sphinx.chat/t/${tribe}`;
                    sendToRedirect(profileUrl);
                  }}
                />
              </>
            )}
          </AssigneeProfile>
        </NormalUser>
      )}
    </div>
  );
}
