/* eslint-disable func-style */
import React, { useCallback, useEffect, useState } from 'react';
import { EuiText, EuiFieldText, EuiGlobalToastList } from '@elastic/eui';
import { observer } from 'mobx-react-lite';
import moment from 'moment';
import { isInvoiceExpired, userCanManageBounty } from 'helpers';
import { SOCKET_MSG, createSocketInstance } from 'config/socket';
import { Button, Divider, Modal } from '../../../../components/common';
import { colors } from '../../../../config/colors';
import { renderMarkdown } from '../../../utils/RenderMarkdown';
import { satToUsd } from '../../../../helpers';
import { useStores } from '../../../../store';
import IconButton from '../../../../components/common/IconButton2';
import ImageButton from '../../../../components/common/ImageButton';
import BountyProfileView from '../../../../bounties/BountyProfileView';
import ButtonSet from '../../../../bounties/BountyModalButtonSet';
import BountyPrice from '../../../../bounties/BountyPrice';
import InvitePeopleSearch from '../../../../components/form/inputs/widgets/PeopleSearch';
import { CodingBountiesProps } from '../../../interfaces';
import LoomViewerRecorder from '../../../utils/LoomViewerRecorder';
import { paidString, unpaidString } from '../constants';
import Invoice from './Invoice';
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
import { getTwitterLink } from './lib';
import CodingMobile from './CodingMobile';

let interval;

function MobileView(props: CodingBountiesProps) {
  const {
    deliverables,
    description,
    ticket_url,
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
    replitLink,
    assigneeHandlerOpen,
    setCreatorStep,
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
    createdURL,
    created,
    loomEmbedUrl,
    org_uuid,
    id,
    localPaid,
    setLocalPaid,
    isMobile,
    actionButtons,
    assigneeLabel
  } = props;
  const color = colors['light'];

  const { ui, main } = useStores();
  const [invoiceStatus, setInvoiceStatus] = useState(false);
  const [keysendStatus, setKeysendStatus] = useState(false);
  const [lnInvoice, setLnInvoice] = useState('');
  const [toasts, setToasts]: any = useState([]);
  const [updatingPayment, setUpdatingPayment] = useState<boolean>(false);
  const [userBountyRole, setUserBountyRole] = useState(false);

  const [paidStatus, setPaidStatus] = useState(paid);

  const [paymentLoading, setPaymentLoading] = useState(false);

  const userPubkey = ui.meInfo?.owner_pubkey;

  let bountyPaid = paid || invoiceStatus || keysendStatus;

  if (localPaid === 'PAID') {
    bountyPaid = true;
  } else if (localPaid === 'UNPAID') {
    bountyPaid = false;
  }

  const showPayBounty =
    !bountyPaid && !invoiceStatus && assignee && assignee.owner_alias.length < 30;

  const pollMinutes = 2;

  const addToast = (type: string) => {
    switch (type) {
      case SOCKET_MSG.invoice_success: {
        return setToasts([
          {
            id: '1',
            title: 'Invoice has been paid',
            color: 'success'
          }
        ]);
      }
      case SOCKET_MSG.keysend_error: {
        return setToasts([
          {
            id: '2',
            title: 'Keysend payment failed',
            toastLifeTimeMs: 10000,
            color: 'error'
          }
        ]);
      }
      case SOCKET_MSG.keysend_success: {
        return setToasts([
          {
            id: '3',
            title: 'Successful keysend payment',
            color: 'success'
          }
        ]);
      }
    }
  };

  const removeToast = () => {
    setToasts([]);
  };

  const startPolling = useCallback(
    async (paymentRequest: string) => {
      let i = 0;
      interval = setInterval(async () => {
        try {
          const invoiceData = await main.pollInvoice(paymentRequest);
          if (invoiceData) {
            if (invoiceData.success && invoiceData.response.settled) {
              clearInterval(interval);

              setLnInvoice('');
              addToast(SOCKET_MSG.invoice_success);
              main.setKeysendInvoice('');
              setLocalPaid('UNKNOWN');
              setInvoiceStatus(true);
              setKeysendStatus(true);
            }
          }

          i++;
          if (i > 22) {
            if (interval) clearInterval(interval);
          }
        } catch (e) {
          console.warn('CodingBounty Invoice Polling Error', e);
        }
      }, 5000);
    },
    [setLocalPaid, main]
  );

  const generateInvoice = async (price: number) => {
    if (created && ui.meInfo?.websocketToken) {
      const data = await main.getLnInvoice({
        amount: price || 0,
        memo: '',
        owner_pubkey: person.owner_pubkey,
        user_pubkey: assignee.owner_pubkey,
        route_hint: assignee.owner_route_hint ?? '',
        created: created ? created?.toString() : '',
        type: 'KEYSEND'
      });

      const paymentRequest = data.response.invoice;

      if (paymentRequest) {
        setLnInvoice(paymentRequest);
        main.setKeysendInvoice(paymentRequest);
        startPolling(paymentRequest);
      }
    }
  };

  useEffect(() => {
    if (main.keysendInvoice !== '') {
      const expired = isInvoiceExpired(main.keysendInvoice);
      if (!expired) {
        startPolling(main.keysendInvoice);
      } else {
        main.setKeysendInvoice('');
      }
    }

    return () => {
      clearInterval(interval);
    };
  }, [main, startPolling]);

  const recallBounties = async () => {
    await main.getPeopleBounties({ resetPage: true, ...main.bountiesStatus });
  };

  const makePayment = async () => {
    setPaymentLoading(true);
    // If the bounty has a commitment fee, add the fee to the user payment
    const price = Number(props.price);
    // if there is an organization and the organization's
    // buudget is sufficient keysend to the user immediately
    // without generating an invoice, else generate an invoice
    if (org_uuid) {
      const organizationBudget = await main.getOrganizationBudget(org_uuid);
      const budget = organizationBudget.total_budget;

      const bounty = await main.getBountyById(id ?? 0);
      if (bounty.length && Number(budget) >= Number(price)) {
        const b = bounty[0];

        if (!b.body.paid) {
          // make keysend payment
          const body = {
            id: id || 0,
            websocket_token: ui.meInfo?.websocketToken || ''
          };

          await main.makeBountyPayment(body);
          setPaymentLoading(false);
          recallBounties();
        }
      } else {
        setPaymentLoading(false);
        return setToasts([
          {
            id: `${Math.random()}`,
            title: 'Insufficient funds in the organization.',
            color: 'danger',
            toastLifeTimeMs: 10000
          }
        ]);
      }
    } else {
      generateInvoice(price || 0);
    }
  };

  const updatePaymentStatus = async (created: number) => {
    await main.updateBountyPaymentStatus(created);
    recallBounties();
  };

  const handleSetAsPaid = async (e: any) => {
    e.stopPropagation();
    setUpdatingPayment(true);
    setPaidStatus(!paidStatus);
    await updatePaymentStatus(created || 0);
    await setExtrasPropertyAndSaveMultiple('paid', {
      award: awardDetails.name
    });

    setTimeout(() => {
      setCreatorStep(0);
      if (setIsPaidStatusPopOver) setIsPaidStatusPopOver(true);
      if (awardDetails?.name !== '') {
        setIsPaidStatusBadgeInfo(true);
      }
      setUpdatingPayment(false);
    }, 3000);
  };

  const handleSetAsUnpaid = async (e: any) => {
    e.stopPropagation();
    setUpdatingPayment(true);
    await updatePaymentStatus(created || 0);
    setLocalPaid('UNPAID');
    setUpdatingPayment(false);
    setPaidStatus(!paidStatus);
    recallBounties();
  };

  const twitterHandler = () => {
    const twitterLink = getTwitterLink({
      title: titleString,
      labels,
      issueCreated: createdURL,
      ownerPubkey: owner_idURL
    });
    sendToRedirect(twitterLink);
  };

  useEffect(() => {
    const onHandle = (event: any) => {
      const res = JSON.parse(event.data);
      if (res.msg === SOCKET_MSG.user_connect) {
        const user = ui.meInfo;
        if (user) {
          user.websocketToken = res.body;
          ui.setMeInfo(user);
        }
      } else if (res.msg === SOCKET_MSG.invoice_success) {
        setLnInvoice('');
        setLocalPaid('UNKNOWN');
        setInvoiceStatus(true);
        addToast(SOCKET_MSG.invoice_success);
      } else if (res.msg === SOCKET_MSG.keysend_success) {
        setLocalPaid('UNKNOWN');
        setKeysendStatus(true);
        addToast(SOCKET_MSG.keysend_success);
      } else if (res.msg === SOCKET_MSG.keysend_error) {
        addToast(SOCKET_MSG.keysend_error);
      }
    };

    const socket: WebSocket = createSocketInstance();
    socket.onopen = () => {
      console.log('Socket connected');
    };

    socket.onmessage = (event: MessageEvent) => {
      onHandle(event);
    };

    socket.onclose = () => {
      console.log('Socket disconnected');
    };
  }, [setLocalPaid, ui]);

  const checkUserBountyRole = useCallback(async () => {
      const canPayBounty = await userCanManageBounty(org_uuid, userPubkey, main);
      setUserBountyRole(canPayBounty);
  }, [main, org_uuid, userPubkey]);

  useEffect(() => {
    checkUserBountyRole();
  }, [checkUserBountyRole]);

  const isOwner =
    { ...person }?.owner_alias &&
    ui.meInfo?.owner_alias &&
    { ...person }?.owner_alias === ui.meInfo?.owner_alias;

  const hasAccess = isOwner || userBountyRole;
  const payBountyDisable = !isOwner && !userBountyRole;

  useEffect(() => {
    setPaidStatus(paid);
  }, [paid]);

  if (isMobile) {
    return (
      <CodingMobile
        {...props}
        paid={paidStatus}
        labels={labels}
        nametag={nametag}
        actionButtons={actionButtons}
        assigneeLabel={assigneeLabel}
        assignee={assignee}
        handleCopyUrl={handleCopyUrl}
        isCopied={isCopied}
        titleString={titleString}
        showPayBounty={showPayBounty}
        markPaidOrUnpaid={
          hasAccess && (
            <IconButton
              width={'100%'}
              height={48}
              style={{
                bottom: '10px',
                border: `1px solid ${color.primaryColor.P400}`,
                background: paidStatus ? color.green1 : color.pureWhite,
                color: paidStatus ? color.white100 : color.borderGreen1
              }}
              text={paidStatus ? unpaidString : paidString}
              loading={saving === 'paid' || updatingPayment}
              endingImg={'/static/mark_unpaid.svg'}
              textStyle={{
                width: '130px',
                display: 'flex',
                justifyContent: 'center',
                fontFamily: 'Barlow',
                marginLeft: '30px',
                fontSize: '15px'
              }}
              onClick={paidStatus ? handleSetAsUnpaid : handleSetAsPaid}
            />
          )
        }
        payBounty={
          hasAccess && (
            <IconButton
              width={'100%'}
              height={48}
              disabled={paymentLoading || payBountyDisable}
              style={{
                bottom: '10px'
              }}
              text={'Pay Bounty'}
              loading={saving === 'paid' || updatingPayment}
              textStyle={{
                display: 'flex',
                justifyContent: 'center',
                fontFamily: 'Barlow',
                fontSize: '15px',
                marginLeft: '30px'
              }}
              hovercolor={color.button_secondary.hover}
              shadowcolor={color.button_secondary.shadow}
              onClick={makePayment}
            />
          )
        }
      />
    );
  }

  return (
    <div>
      {hasAccess ? (
        /*
         * creator view
         */
        <>
          {creatorStep === 0 && (
            <Creator
              onClick={() => {
                if (setIsPaidStatusPopOver) setIsPaidStatusPopOver(false);
              }}
            >
              <>
                {bountyPaid && (
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
                {bountyPaid && (
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

                <CreatorDescription paid={bountyPaid} color={color}>
                  <div className="CreatorDescriptionOuterContainerCreatorView">
                    <div className="CreatorDescriptionInnerContainerCreatorView">
                      <div>{nametag}</div>
                      {!bountyPaid && hasAccess && (
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
                            buttonTextStyle={{
                              paddingRight: '50px'
                            }}
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
                            disabled={!props?.deleteAction}
                            buttonAction={props?.deleteAction}
                            buttonTextStyle={{
                              paddingRight: '45px'
                            }}
                          />
                        </div>
                      )}
                    </div>
                    <TitleBox color={color}>{titleString}</TitleBox>
                    <LanguageContainer>
                      {dataValue &&
                        dataValue?.length > 0 &&
                        dataValue?.map((lang: any, index: number) => (
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
                          status={bountyPaid ? 'completed' : 'assigned'}
                          canViewProfile={false}
                          statusStyle={{
                            width: '66px',
                            height: '16px',
                            background: bountyPaid ? color.statusCompleted : color.statusAssigned
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
                        {!bountyPaid && (
                          <div
                            className="AssigneeCloseButtonContainer"
                            onClick={() => {
                              changeAssignedPerson();
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
                            width: '114px',
                            paddingLeft: '20px'
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
                      price={props?.price || 0}
                      sessionLength={props?.estimated_session_length}
                      style={{
                        padding: 0,
                        margin: 0
                      }}
                    />

                    {lnInvoice && !invoiceStatus && (
                      <Invoice
                        startDate={
                          new Date(moment().add(pollMinutes, 'minutes').format().toString())
                        }
                        invoiceStatus={invoiceStatus}
                        invoiceTime={pollMinutes}
                        lnInvoice={lnInvoice}
                      />
                    )}
                    {/**
                     * LNURL AUTH users alias are their public keys
                     * which make them so longF
                     * A non LNAUTh user alias is shorter
                     */}
                    {showPayBounty && (
                      <>
                        <Button
                          disabled={paymentLoading || payBountyDisable}
                          iconSize={14}
                          width={220}
                          height={48}
                          onClick={makePayment}
                          style={{ marginTop: '30px', marginBottom: '-20px', textAlign: 'left' }}
                          text="Pay Bounty"
                          ButtonTextStyle={{ padding: 0 }}
                        />
                      </>
                    )}
                  </BountyPriceContainer>
                  <div className="buttonSet">
                    <ButtonSet
                      githubShareAction={() => {
                        const repoUrl = ticket_url
                          ? ticket_url
                          : `https://github.com/${repo}/issues/${issue}`;
                        sendToRedirect(repoUrl);
                      }}
                      copyURLAction={handleCopyUrl}
                      copyStatus={isCopied ? 'Copied' : 'Copy Link'}
                      twitterAction={twitterHandler}
                      replitLink={replitLink}
                      tribe={tribe !== 'none' && tribe}
                      tribeFunction={() => {
                        const profileUrl = `https://community.sphinx.chat/t/${tribe}`;
                        sendToRedirect(profileUrl);
                      }}
                      showGithubBtn={!!ticket_url}
                    />
                  </div>
                  <BottomButtonContainer>
                    {bountyPaid ? (
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
                        text={paid ? unpaidString : paidString}
                        loading={saving === 'paid' || updatingPayment}
                        endingImg={'/static/mark_unpaid.svg'}
                        textStyle={{
                          width: '130px',
                          display: 'flex',
                          justifyContent: 'center',
                          fontFamily: 'Barlow',
                          marginLeft: '30px'
                        }}
                        onClick={handleSetAsUnpaid}
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
                        text={paidString}
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
                        onClick={(e: any) => {
                          e.stopPropagation();
                          setCreatorStep(1);
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
                    onChange={(e: any) => {
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
                  onClick={(e: any) => {
                    e.stopPropagation();
                    setCreatorStep(2);
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
                {awards?.map((award: any, index: number) => (
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
                  text={selectedAward === '' ? 'Skip and Mark Paid' : paidString}
                  loading={isMarkPaidSaved || updatingPayment}
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
                  onClick={handleSetAsPaid}
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
                <EuiText className="autoCompleteHeaderText">Assign Developer</EuiText>
                <InvitePeopleSearch
                  peopleList={peopleList}
                  isProvidingHandler={true}
                  handleAssigneeDetails={(value: any) => {
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
          {bountyPaid && (
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
          <CreatorDescription paid={bountyPaid} color={color}>
            <div className="DescriptionUpperContainerNormalView">
              <div>{nametag}</div>
              <TitleBox color={color}>{titleString}</TitleBox>
              <LanguageContainer>
                {dataValue &&
                  dataValue?.length > 0 &&
                  dataValue?.map((lang: any, index: number) => (
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
              {loomEmbedUrl && (
                <>
                  <div className="loomContainer" />
                  <EuiText className="loomHeading">Video</EuiText>
                  <LoomViewerRecorder
                    readOnly
                    style={{ marginTop: 10 }}
                    loomEmbedUrl={loomEmbedUrl}
                  />
                </>
              )}
            </DescriptionBox>
          </CreatorDescription>

          <AssigneeProfile color={color}>
            {bountyPaid ? (
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
                    price={props?.price || 0}
                    sessionLength={props?.estimated_session_length}
                    style={{
                      padding: 0,
                      margin: 0
                    }}
                  />
                </BountyPriceContainer>
                <ButtonSet
                  showGithubBtn={!!ticket_url}
                  githubShareAction={() => {
                    const repoUrl = ticket_url
                      ? ticket_url
                      : `https://github.com/${repo}/issues/${issue}`;
                    sendToRedirect(repoUrl);
                  }}
                  copyURLAction={handleCopyUrl}
                  copyStatus={isCopied ? 'Copied' : 'Copy Link'}
                  twitterAction={twitterHandler}
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
                    price={props?.price || 0}
                    sessionLength={props?.estimated_session_length}
                    style={{
                      padding: 0,
                      margin: 0
                    }}
                  />
                </BountyPriceContainer>
                <ButtonSet
                  showGithubBtn={!!ticket_url}
                  githubShareAction={() => {
                    const repoUrl = ticket_url
                      ? ticket_url
                      : `https://github.com/${repo}/issues/${issue}`;
                    sendToRedirect(repoUrl);
                  }}
                  copyURLAction={handleCopyUrl}
                  copyStatus={isCopied ? 'Copied' : 'Copy Link'}
                  twitterAction={twitterHandler}
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
                    price={props?.price || 0}
                    sessionLength={props.estimated_session_length}
                    style={{
                      padding: 0,
                      margin: 0
                    }}
                  />
                </BountyPriceContainer>
                <ButtonSet
                  showGithubBtn={!!ticket_url}
                  githubShareAction={() => {
                    const repoUrl = ticket_url
                      ? ticket_url
                      : `https://github.com/${repo}/issues/${issue}`;
                    sendToRedirect(repoUrl);
                  }}
                  copyURLAction={handleCopyUrl}
                  copyStatus={isCopied ? 'Copied' : 'Copy Link'}
                  twitterAction={twitterHandler}
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
      <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={6000} />
    </div>
  );
}
export default observer(MobileView);
