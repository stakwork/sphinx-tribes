/* eslint-disable func-style */
import { EuiFieldText, EuiText } from '@elastic/eui';
import React, { useCallback, useEffect, useLayoutEffect, useRef, useState } from 'react';
import { useLocation } from 'react-router-dom';
import styled from 'styled-components';
import api from '../../../api';
import { colors } from '../../../colors';
import Form from '../../../components/form';
import InvitePeopleSearch from '../../../components/form/inputs/widgets/PeopleSearch';
import { sendBadgeSchema } from '../../../components/form/schema';
import { extractGithubIssue, extractGithubIssueFromUrl, formatPrice, satToUsd } from '../../../helpers';
import { useIsMobile } from '../../../hooks';
import { Button, Divider, Modal, Paragraph, Title } from '../../../components/common';
import ImageButton from '../../../components/common/Image_button';
import ButtonSet from '../../../bounties/bountyModal_button_set';
import BountyPrice from '../../../bounties/bounty_price';
import BountyProfileView from '../../../bounties/bounty_profile_view';
import IconButton from '../../../components/common/icon_button';
import { useStores } from '../../../store';
import FavoriteButton from '../../utils/favoriteButton';
import GalleryViewer from '../../utils/galleryViewer';
import { LanguageObject, awards } from '../../utils/language_label_style';
import LoomViewerRecorder from '../../utils/loomViewerRecorder';
import NameTag from '../../utils/nameTag';
import { renderMarkdown } from '../../utils/renderMarkdown';
import GithubStatusPill from '../parts/statusPill';

function useQuery() {
  const { search } = useLocation();

  return React.useMemo(() => new URLSearchParams(search), [search]);
}

export default function WantedSummary(props: any) {
  const {
    description,
    deliverables,
    priceMin,
    priceMax,
    ticketUrl,
    gallery,
    person,
    created,
    repo,
    issue,
    price,
    type,
    tribe,
    paid,
    badgeRecipient,
    loomEmbedUrl,
    codingLanguage,
    estimate_session_length,
    assignee,
    fromBountyPage,
    wanted_type,
    one_sentence_summary,
    github_description,
    show,
    setIsModalSideButton,
    setIsExtraStyle,
    formSubmit
  } = props;
  const title = one_sentence_summary;
  const [envHeight, setEnvHeight] = useState('100%');
  const imgRef: any = useRef(null);

  const isMobile = useIsMobile();
  const { main, ui } = useStores();
  const { peopleWanteds } = main;
  const color = colors['light'];

  const [tribeInfo, setTribeInfo]: any = useState(null);
  const [assigneeInfo, setAssigneeInfo]: any = useState(null);
  const [saving, setSaving]: any = useState('');
  const [isCopied, setIsCopied] = useState(false);
  const [owner_idURL, setOwnerIdURL] = useState('');
  const [createdURL, setCreatedURL] = useState('');
  const [dataValue, setDataValue] = useState([]);
  const [peopleList, setPeopleList] = useState<any>();
  const [isAssigned, setIsAssigned] = useState<boolean>(false);
  const [assignedPerson, setAssignedPerson] = useState<any>();
  const [replitLink, setReplitLink] = useState('');
  const [creatorStep, setCreatorStep] = useState<number>(0);
  const [bountyPrice, setBountyPrice] = useState<any>(price ?? priceMin ?? 0);
  const [selectedAward, setSelectedAward] = useState('');
  const [isPaidStatusPopOver, setIsPaidStatusPopOver] = useState<boolean>(false);
  const [isPaidStatusBadgeInfo, setIsPaidStatusBadgeInfo] = useState<boolean>(false);
  const [isMarkPaidSaved, setIsMarkPaidSaved] = useState<boolean>(false);
  const [awardDetails, setAwardDetails] = useState<any>({
    name: '',
    image: ''
  });

  useEffect(() => {
    if (description) {
      setReplitLink(
        description.match(
          /https?:\/\/(www\.)?[replit]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)/
        )
      );
    }
  }, [description]);

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsPaidStatusPopOver(false);
    }, 7000);

    return () => {
      clearTimeout(timer);
    };
  }, [isPaidStatusPopOver]);

  const handleAwards = (optionId) => {
    setSelectedAward(optionId);
  };

  const [showBadgeAwardDialog, setShowBadgeAwardDialog] = useState(false);

  const isMine = ui.meInfo?.owner_pubkey === person?.owner_pubkey;

  const [labels, setLabels] = useState([]);
  const [assigneeValue, setAssigneeValue] = useState(false);

  const assigneeHandlerOpen = () => setAssigneeValue((assigneeValue) => !assigneeValue);

  useLayoutEffect(() => {
    if (imgRef && imgRef.current) {
      if (imgRef.current?.offsetHeight > 100) {
        setEnvHeight(imgRef.current?.offsetHeight);
      }
    }
  }, [imgRef]);

  useEffect(() => {
    if (assignee?.owner_alias) {
      setIsAssigned(true);
    }
  }, [assignee]);

  useEffect(() => {
    (async () => {
      try {
        const response = await api.get(`people?page=1&search=&sortBy=last_login&limit=100`);
        setPeopleList(response);
      } catch (error) {
        console.log(error);
      }
    })();
  }, []);

  const handleAssigneeDetails = useCallback(
    (value) => {
      setIsAssigned(true);
      setAssignedPerson(value);
      assigneeHandlerOpen();
      const newValue = {
        title: title,
        wanted_type: wanted_type,
        one_sentence_summary: one_sentence_summary,
        ticketUrl: ticketUrl,
        github_description: github_description,
        description: description,
        price: price,
        assignee: {
          owner_alias: value?.owner_alias || '',
          owner_pubkey: value?.owner_pubkey || '',
          img: value?.img || '',
          value: value?.owner_pubkey || '',
          label: `${value.owner_alias} (${value.owner_alias.toLowerCase().replace(' ', '')})` || ''
        },
        codingLanguage: codingLanguage?.map((x) => {
          return { ...x };
        }),
        estimate_session_length: estimate_session_length,
        show: show,
        type: type,
        created: created
      };
      formSubmit && formSubmit(newValue);
    },
    [codingLanguage, created, description, estimate_session_length, formSubmit, github_description, one_sentence_summary, price, show, ticketUrl, title, type, wanted_type]
  );

  const changeAssignedPerson = useCallback(() => {
    setIsAssigned(false);
  }, []);

  useEffect(() => {
    (async () => {
      if (props.assignee) {
        try {
          const p = await main.getPersonByPubkey(props.assignee.owner_pubkey);
          setAssigneeInfo(p);
        } catch (e) {
          console.log('e', e);
        }
      }
      if (tribe) {
        try {
          const t = await main.getSingleTribeByUn(tribe);
          setTribeInfo(t);
        } catch (e) {
          console.log('e', e);
        }
      }
    })();
  }, [main, props.assignee, tribe]);

  useEffect(() => {
    let res;
    if (codingLanguage?.length > 0) {
      res = LanguageObject?.filter((value) => {
        return codingLanguage?.find((val) => {
          return val.label === value.label;
        });
      });
    }
    setDataValue(res);
    setLabels(res);
  }, [codingLanguage]);

  const searchParams = useQuery();

  useEffect(() => {
    const owner_id = searchParams.get('owner_id');
    const created = searchParams.get('created');
    setOwnerIdURL(owner_id ?? '');
    setCreatedURL(created ?? '');
  }, [owner_idURL, createdURL, searchParams]);

  useEffect(() => {
    if (codingLanguage) {
      const values = codingLanguage.map((value) => ({ ...value }));
      setLabels(values);
    }
  }, [codingLanguage]);

  async function setExtrasPropertyAndSave(propertyName: string, value: any) {
    if (peopleWanteds) {
      setSaving(propertyName);
      try {
        const [clonedEx, targetIndex] = await main.setExtrasPropertyAndSave(
          'wanted',
          propertyName,
          created,
          value
        );

        // saved? ok update in wanted list if found
        const peopleWantedsClone: any = [...peopleWanteds];
        const indexFromPeopleWanted = peopleWantedsClone.findIndex((f) => {
          const val = f.body || {};
          return f.person.owner_pubkey === ui.meInfo?.owner_pubkey && val.created === created;
        });

        // if we found it in the wanted list, update in people wanted list
        if (indexFromPeopleWanted > -1) {
          // if it should be hidden now, remove it from the list
          if ('show' in clonedEx[targetIndex] && clonedEx[targetIndex].show === false) {
            peopleWantedsClone.splice(indexFromPeopleWanted, 1);
          } else {
            // gotta update person extras! this is what is used for summary viewer
            const personClone: any = person;
            personClone.extras['wanted'][targetIndex] = clonedEx[targetIndex];

            peopleWantedsClone[indexFromPeopleWanted] = {
              person: personClone,
              body: clonedEx[targetIndex]
            };
          }
          main.setPeopleWanteds(peopleWantedsClone);
        }
      } catch (e) {
        console.log('e', e);
      }

      setSaving('');
    }
  }

  async function setExtrasPropertyAndSaveMultiple(propertyName: any, dataObject: any) {
    if (peopleWanteds) {
      setIsMarkPaidSaved(true);
      try {
        const [clonedEx, targetIndex] = await main.setExtrasMultipleProperty(
          dataObject,
          'wanted',
          created
        );

        // saved? ok update in wanted list if found
        const peopleWantedsClone: any = [...peopleWanteds];
        const indexFromPeopleWanted = peopleWantedsClone.findIndex((f) => {
          const val = f.body || {};
          return f.person.owner_pubkey === ui.meInfo?.owner_pubkey && val.created === created;
        });

        // if we found it in the wanted list, update in people wanted list
        if (indexFromPeopleWanted > -1) {
          // if it should be hidden now, remove it from the list
          if ('show' in clonedEx[targetIndex] && clonedEx[targetIndex].show === false) {
            peopleWantedsClone.splice(indexFromPeopleWanted, 1);
          } else {
            // gotta update person extras! this is what is used for summary viewer
            const personClone: any = person;
            personClone.extras['wanted'][targetIndex] = clonedEx[targetIndex];

            peopleWantedsClone[indexFromPeopleWanted] = {
              person: personClone,
              body: clonedEx[targetIndex]
            };
          }

          main.setPeopleWanteds(peopleWantedsClone);
        }
      } catch (e) {
        console.log('e', e);
      }

      setIsMarkPaidSaved(false);
    }
  }

  const handleCopyUrl = useCallback(() => {
    const el = document.createElement('input');
    el.value = window.location.href;
    document.body.appendChild(el);
    el.select();
    document.execCommand('copy');
    document.body.removeChild(el);
    setIsCopied(true);
  }, []);

  async function sendBadge(body: any) {
    const { recipient, badge } = body;

    setSaving('badgeRecipient');
    try {
      if (badge?.amount < 1) {
        alert("You don't have any of the selected badge");
        throw new Error("You don't have any of the selected badge");
      }

      // first get the user's liquid address
      const recipientDetails = await main.getPersonByPubkey(recipient.owner_pubkey);

      const liquidAddress =
        recipientDetails?.extras?.liquid && recipientDetails?.extras?.liquid[0]?.value;

      if (!liquidAddress) {
        alert('This user has not provided an L-BTC address');
        throw new Error('This user has not provided an L-BTC address');
      }

      // asset: number
      // to: string
      // amount?: number
      // memo: string
      const pack = {
        asset: badge.id,
        to: liquidAddress,
        amount: 1,
        memo: props.ticketUrl
      };

      const r = await main.sendBadgeOnLiquid(pack);

      if (r.ok) {
        await setExtrasPropertyAndSave('badgeRecipient', recipient.owner_pubkey);
        setShowBadgeAwardDialog(false);
      } else {
        alert(r.statusText || 'Operation failed! Contact support.');
        throw new Error(r.statusText);
      }
    } catch (e) {
      console.log(e);
    }

    setSaving('');
  }

  const heart = <FavoriteButton />;

  const viewGithub = (
    ticketUrl && (
      <Button
        text={'Original Ticket'}
        color={'white'}
        endingIcon={'launch'}
        iconSize={14}
        style={{ fontSize: 14, height: 48, width: '100%', marginBottom: 20 }}
        onClick={() => {
          const repoUrl = ticketUrl ? ticketUrl : `https://github.com/${repo}/issues/${issue}`;
          sendToRedirect(repoUrl);
        }}
      />
    )
  );

  const viewTribe = tribe && tribe !== 'none' && (
    <Button
      text={'View Tribe'}
      color={'white'}
      leadingImgUrl={tribeInfo?.img || ' '}
      endingIcon={'launch'}
      iconSize={14}
      imgStyle={{ position: 'absolute', left: 10 }}
      style={{ fontSize: 14, height: 48, width: '100%', marginBottom: 20 }}
      onClick={() => {
        const profileUrl = `https://community.sphinx.chat/t/${tribe}`;
        sendToRedirect(profileUrl);
      }}
    />
  );

  const addToFavorites = tribe && tribe !== 'none' && (
    <Button
      text={'Add to Favorites'}
      color={'white'}
      icon={'favorite_outline'}
      iconSize={18}
      iconStyle={{ left: 14 }}
      style={{
        fontSize: 14,
        height: 48,
        width: '100%',
        marginBottom: 20,
        paddingLeft: 5
      }}
      onClick={() => { }}
    />
  );

  const copyLink = (
    <Button
      text={isCopied ? 'Copied' : 'Copy Link'}
      color={'white'}
      icon={'content_copy'}
      iconSize={18}
      iconStyle={{ left: 14 }}
      style={{
        fontSize: 14,
        height: 48,
        width: '100%',
        marginBottom: 20,
        paddingLeft: 5
      }}
      onClick={handleCopyUrl}
    />
  );

  const shareOnTwitter = (
    <Button
      text={'Share to Twitter'}
      color={'white'}
      icon={'share'}
      iconSize={18}
      iconStyle={{ left: 14 }}
      style={{
        fontSize: 14,
        height: 48,
        width: '100%',
        marginBottom: 20,
        paddingLeft: 5
      }}
      onClick={() => {
        const twitterLink = `https://twitter.com/intent/tweet?text=Hey, I created a new ticket on Sphinx community.%0A${title} %0A&url=https://community.sphinx.chat/p?owner_id=${owner_idURL}%26created${createdURL} %0A%0A&hashtags=${labels && labels.map((x: any) => x.label)
          },sphinxchat`;
        sendToRedirect(twitterLink);
      }}
    />
  );

  //  if my own, show this option to show/hide
  const markPaidButton = (
    <Button
      color={'primary'}
      iconSize={14}
      style={{ fontSize: 14, height: 48, width: '100%', marginBottom: 20 }}
      endingIcon={'paid'}
      text={paid ? 'Mark Unpaid' : 'Mark Paid'}
      loading={saving === 'paid'}
      onClick={(e) => {
        e.stopPropagation();
        setExtrasPropertyAndSave('paid', !paid);
      }}
    />
  );

  const awardBadgeButton = !badgeRecipient && (
    <Button
      color={'primary'}
      iconSize={14}
      endingIcon={'offline_bolt'}
      style={{ fontSize: 14, height: 48, width: '100%', marginBottom: 20 }}
      text={badgeRecipient ? 'Badge Awarded' : 'Award Badge'}
      disabled={badgeRecipient ? true : false}
      loading={saving === 'badgeRecipient'}
      onClick={(e) => {
        e.stopPropagation();
        if (!badgeRecipient) {
          setShowBadgeAwardDialog(true);
        }
      }}
    />
  );

  const actionButtons = isMine && (
    <ButtonRow>
      {showBadgeAwardDialog ? (
        <>
          <Form
            loading={saving === 'badgeRecipient'}
            smallForm
            buttonsOnBottom
            wrapStyle={{ padding: 0, margin: 0, maxWidth: '100%' }}
            close={() => setShowBadgeAwardDialog(false)}
            onSubmit={(e) => {
              sendBadge(e);
            }}
            submitText={'Send Badge'}
            schema={sendBadgeSchema}
          />
          <div style={{ height: 100 }} />
        </>
      ) : (
        <>
          {markPaidButton}
          {awardBadgeButton}
        </>
      )}
    </ButtonRow>
  );

  function sendToRedirect(url) {
    const el = document.createElement('a');
    el.href = url;
    el.target = '_blank';
    el.click();
  }

  const nametag = (
    <NameTag
      iconSize={24}
      textSize={13}
      style={{ marginBottom: 10 }}
      {...person}
      created={created}
      widget={'wanted'}
    />
  );

  function renderCodingTask() {
    const { status } = ticketUrl
      ? extractGithubIssueFromUrl(person, ticketUrl)
      : extractGithubIssue(person, repo, issue);

    let assigneeLabel: any = null;
    if (assigneeInfo) {
      if (!isMobile) {
        assigneeLabel = (
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              fontSize: 12,
              color: color.grayish.G100,
              marginTop: isMobile ? 20 : 0,
              marginLeft: '-16px'
            }}
          >
            <Img
              src={assigneeInfo.img || '/static/person_placeholder.png'}
              style={{ borderRadius: 30 }}
            />

            <Assignee
              color={color}
              onClick={() => {
                const profileUrl = `https://community.sphinx.chat/p/${assigneeInfo.owner_pubkey}`;
                sendToRedirect(profileUrl);
              }}
              style={{ marginLeft: 3, fontWeight: 500, cursor: 'pointer' }}
            >
              {assigneeInfo.owner_alias}
            </Assignee>
          </div>
        );
      } else {
        assigneeLabel = (
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              fontSize: 12,
              color: color.grayish.G100,
              marginLeft: '16px'
            }}
          >
            <Img
              src={assigneeInfo.img || '/static/person_placeholder.png'}
              style={{ borderRadius: 30 }}
            />

            <Assignee
              color={color}
              onClick={() => {
                const profileUrl = `https://community.sphinx.chat/p/${assigneeInfo.owner_pubkey}`;
                sendToRedirect(profileUrl);
              }}
              style={{ marginLeft: 3, fontWeight: 500, cursor: 'pointer' }}
            >
              {assigneeInfo.owner_alias}
            </Assignee>
          </div>
        );
      }
    }

    if (isMobile) {
      return (
        <div style={{ padding: 20, overflow: 'auto' }}>
          <Pad>
            {nametag}

            <T>{title}</T>

            <div
              style={{
                display: 'flex',
                flexDirection: 'row'
              }}
            >
              <GithubStatusPill status={status} assignee={assignee} />
              {assigneeLabel}
              {ticketUrl && (
                <GithubIconMobile
                  onClick={(e) => {
                    e.stopPropagation();
                    window.open(ticketUrl, '_blank');
                  }}
                >
                  <img height={'100%'} width={'100%'} src="/static/github_logo.png" alt="github" />
                </GithubIconMobile>
              )}
              {loomEmbedUrl && (
                <LoomIconMobile
                  onClick={(e) => {
                    e.stopPropagation();
                    window.open(loomEmbedUrl, '_blank');
                  }}
                >
                  <img height={'100%'} width={'100%'} src="/static/loom.png" alt="loomVideo" />
                </LoomIconMobile>
              )}
            </div>

            <EuiText
              style={{
                fontSize: '13px',
                color: color.grayish.G100,
                fontWeight: '500'
              }}
            >
              {estimate_session_length && 'Session:'}{' '}
              <span
                style={{
                  fontWeight: '500',
                  color: color.pureBlack
                }}
              >
                {estimate_session_length ?? ''}
              </span>
            </EuiText>
            <div
              style={{
                width: '100%',
                display: 'flex',
                flexDirection: 'row',
                marginTop: '10px',
                minHeight: '60px'
              }}
            >
              {labels?.length > 0 &&
                labels?.map((x: any) => {
                  return (
                    <>
                      <div
                        style={{
                          display: 'flex',
                          flexWrap: 'wrap',
                          height: '22px',
                          width: 'fit-content',
                          backgroundColor: color.grayish.G1000,
                          border: `1px solid ${color.grayish.G70}`,
                          padding: '3px 10px',
                          borderRadius: '20px',
                          marginRight: '3px',
                          boxShadow: `1px 1px ${color.grayish.G70}`
                        }}
                      >
                        <div
                          style={{
                            fontSize: '10px',
                            color: color.black300
                          }}
                        >
                          {x.label}
                        </div>
                      </div>
                    </>
                  );
                })}
            </div>

            <div style={{ height: 10 }} />
            <ButtonRow style={{ margin: '10px 0' }}>
              {viewGithub}
              {viewTribe}
              {addToFavorites}
              {copyLink}
              {shareOnTwitter}
            </ButtonRow>

            {actionButtons}

            <LoomViewerRecorder readOnly loomEmbedUrl={loomEmbedUrl} style={{ marginBottom: 20 }} />

            <Divider />
            <Y>
              <P color={color}>
                <B color={color}>{formatPrice(price)}</B> SAT /{' '}
                <B color={color}>{satToUsd(price)}</B> USD
              </P>
              {heart}
            </Y>
            <Divider style={{ marginBottom: 20 }} />
            <D color={color}>{renderMarkdown(description)}</D>
          </Pad>
        </div>
      );
    }

    // desktop view
    if (fromBountyPage) {
      return (
        <>
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
                          <Profile>{nametag}</Profile>
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
                        <TitleBox color={color}>{title}</TitleBox>
                        <LanguageContainer>
                          {dataValue &&
                            dataValue?.length > 0 &&
                            dataValue?.map((lang: any, index) => {
                              return (
                                <CodingLabels
                                  key={index}
                                  styledColors={color}
                                  border={lang?.border}
                                  color={lang?.color}
                                  background={lang?.background}
                                >
                                  <EuiText className="LanguageText">{lang?.label}</EuiText>
                                </CodingLabels>
                              );
                            })}
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
                            const twitterLink = `https://twitter.com/intent/tweet?text=Hey, I created a new ticket on Sphinx community.%0A${title} %0A&url=https://community.sphinx.chat/p?owner_id=${owner_idURL}%26created${createdURL} %0A%0A&hashtags=${labels && labels.map((x: any) => x.label)
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
                            hoverColor={color.button_primary.hover}
                            activeColor={color.button_primary.active}
                            shadowColor={color.button_primary.shadow}
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
                      hoverColor={color.button_secondary.hover}
                      activeColor={color.button_secondary.active}
                      shadowColor={color.button_secondary.shadow}
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
                      hoverColor={color.button_primary.hover}
                      activeColor={color.button_primary.active}
                      shadowColor={color.button_primary.shadow}
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
                  <Profile>{nametag}</Profile>
                  <TitleBox color={color}>{title}</TitleBox>
                  <LanguageContainer>
                    {dataValue &&
                      dataValue?.length > 0 &&
                      dataValue?.map((lang: any, index) => {
                        return (
                          <CodingLabels
                            key={index}
                            styledColors={color}
                            border={lang?.border}
                            color={lang?.color}
                            background={lang?.background}
                          >
                            <EuiText className="LanguageText">{lang?.label}</EuiText>
                          </CodingLabels>
                        );
                      })}
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
                        const twitterLink = `https://twitter.com/intent/tweet?text=Hey, I created a new ticket on Sphinx community.%0A${title} %0A&url=https://community.sphinx.chat/p?owner_id=${owner_idURL}%26created${createdURL} %0A%0A&hashtags=${labels && labels.map((x: any) => x.label)
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
                        const twitterLink = `https://twitter.com/intent/tweet?text=Hey, I created a new ticket on Sphinx community.%0A${title} %0A&url=https://community.sphinx.chat/p?owner_id=${owner_idURL}%26created${createdURL} %0A%0A&hashtags=${labels && labels.map((x: any) => x.label)
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
                          // style={{ marginTop: 5 }}
                          onClick={props.extraModalFunction}
                          color="primary"
                          hoverColor={color.button_secondary.hover}
                          activeColor={color.button_secondary.active}
                          shadowColor={color.button_secondary.shadow}
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
                        const twitterLink = `https://twitter.com/intent/tweet?text=Hey, I created a new ticket on Sphinx community.%0A${title} %0A&url=https://community.sphinx.chat/p?owner_id=${owner_idURL}%26created${createdURL} %0A%0A&hashtags=${labels && labels.map((x: any) => x.label)
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
        </>
      );
    }
    return (
      <>
        {paid && (
          <Img
            src={'/static/paid_ribbon.svg'}
            style={{
              position: 'absolute',
              top: -1,
              right: 0,
              width: 64,
              height: 72,
              zIndex: 100,
              pointerEvents: 'none'
            }}
          />
        )}
        <Wrap color={color}>
          <div
            style={{
              width: 700,
              borderRight: `1px solid ${color.grayish.G600}`,
              minHeight: '100%',
              overflow: 'auto'
            }}
          >
            <SectionPad style={{ minHeight: 160, maxHeight: 160 }}>
              <Title>{title}</Title>
              <div style={{ display: 'flex', marginTop: 12 }}>
                <GithubStatusPill status={status} assignee={assignee} style={{ marginRight: 25 }} />
                {assigneeLabel}
                {ticketUrl && (
                  <GithubIcon
                    onClick={(e) => {
                      e.stopPropagation();
                      window.open(ticketUrl, '_blank');
                    }}
                  >
                    <img
                      height={'100%'}
                      width={'100%'}
                      src="/static/github_logo.png"
                      alt="github"
                    />
                  </GithubIcon>
                )}
                {loomEmbedUrl && (
                  <LoomIcon
                    onClick={(e) => {
                      e.stopPropagation();
                      window.open(loomEmbedUrl, '_blank');
                    }}
                  >
                    <img height={'100%'} width={'100%'} src="/static/loom.png" alt="loomVideo" />
                  </LoomIcon>
                )}
              </div>
              <div
                style={{
                  marginTop: '2px'
                }}
              >
                <EuiText
                  style={{
                    fontSize: '13px',
                    color: color.text2_4,
                    fontWeight: '500'
                  }}
                >
                  {estimate_session_length && 'Session:'}{' '}
                  <span
                    style={{
                      fontWeight: '500',
                      color: color.pureBlack
                    }}
                  >
                    {estimate_session_length ?? ''}
                  </span>
                </EuiText>
              </div>
            </SectionPad>
            <Divider />

            <SectionPad>
              <Paragraph
                style={{
                  overflow: 'hidden',
                  wordBreak: 'normal'
                }}
              >
                {renderMarkdown(description)}
              </Paragraph>

              <LoomViewerRecorder readOnly style={{ marginTop: 10 }} loomEmbedUrl={loomEmbedUrl} />
            </SectionPad>
          </div>

          <div style={{ width: 320, height: envHeight, overflow: 'auto' }}>
            <SectionPad style={{ minHeight: 160, maxHeight: 160 }}>
              <div
                style={{
                  display: 'flex',
                  width: '100%',
                  justifyContent: 'space-between'
                }}
              >
                {nametag}
              </div>
              {/* <Img
                src={'/static/github_logo2.png'}
                style={{ width: 77, height: 43 }}
              /> */}

              <div
                style={{
                  minHeight: '60px',
                  width: '100%',
                  display: 'flex',
                  flexDirection: 'row'
                }}
              >
                {labels?.length > 0 &&
                  labels?.map((x: any) => {
                    return (
                      <>
                        <div
                          style={{
                            display: 'flex',
                            flexWrap: 'wrap',
                            height: '22px',
                            minWidth: 'fit-content',
                            backgroundColor: color.grayish.G1000,
                            border: `1px solid ${color.grayish.G70}`,
                            padding: '3px 10px',
                            borderRadius: '20px',
                            marginRight: '3px',
                            boxShadow: `1px 1px ${color.grayish.G70}`
                          }}
                        >
                          <div
                            style={{
                              fontSize: '10px',
                              color: color.black300
                            }}
                          >
                            {x.label}
                          </div>
                        </div>
                      </>
                    );
                  })}
              </div>
            </SectionPad>
            <Divider />
            <SectionPad>
              <Y style={{ padding: 0 }}>
                <P color={color}>
                  <B color={color}>{formatPrice(price)}</B> SAT /{' '}
                  <B color={color}>{satToUsd(price)}</B> USD
                </P>
              </Y>
            </SectionPad>

            <Divider />

            <SectionPad>
              <ButtonRow>
                {viewGithub}
                {viewTribe}
                {addToFavorites}
                {copyLink}
                {shareOnTwitter}
              </ButtonRow>

              {actionButtons}
            </SectionPad>
          </div>
        </Wrap>
      </>
    );
  }

  if (type === 'coding_task' || type === 'wanted_coding_task' || type === 'freelance_job_request') {
    return renderCodingTask();
  }

  if (isMobile) {
    return (
      <div style={{ padding: 20, overflow: 'auto' }}>
        <Pad>
          {nametag}

          <T>{title || 'No title'}</T>
          <Divider
            style={{
              marginTop: 22
            }}
          />
          <Y>
            <P color={color}>
              {formatPrice(priceMin) || '0'} <B color={color}>SAT</B> - {formatPrice(priceMax)}{' '}
              <B color={color}>SAT</B>
            </P>
            {heart}
          </Y>
          <Divider style={{ marginBottom: 22 }} />

          <D color={color}>{renderMarkdown(description)}</D>
          <GalleryViewer
            gallery={gallery}
            showAll={true}
            selectable={false}
            wrap={false}
            big={true}
          />
        </Pad>
      </div>
    );
  }

  return (
    <div
      style={{
        paddingTop: gallery && '40px'
      }}
    >
      <Wrap color={color}>
        <div>
          <GalleryViewer
            innerRef={imgRef}
            style={{ width: 507, height: 'fit-content' }}
            gallery={gallery}
            showAll={false}
            selectable={false}
            wrap={false}
            big={true}
          />
        </div>
        <div
          style={{
            width: 316,
            padding: '40px 20px',
            overflowY: 'auto',
            height: envHeight
          }}
        >
          <Pad>
            {nametag}

            <Title>{title}</Title>

            <Divider style={{ marginTop: 22 }} />
            <Y>
              <P color={color}>
                {formatPrice(priceMin) || '0'} <B color={color}>SAT</B> -{' '}
                {formatPrice(priceMax) || '0'} <B color={color}>SAT</B>
              </P>
              {heart}
            </Y>
            <Divider style={{ marginBottom: 22 }} />

            <Paragraph>{renderMarkdown(description)}</Paragraph>
          </Pad>
        </div>
      </Wrap>
    </div>
  );
}

interface colorProps {
  color?: any;
  isPaidStatusPopOver?: any;
  isPaidStatusBadgeInfo?: any;
}
interface styleProps extends colorProps {
  paid?: string;
}

const Wrap = styled.div<colorProps>`
  display: flex;
  width: 100%;
  height: 100%;
  min-width: 800px;
  font-style: normal;
  font-weight: 500;
  font-size: 24px;
  color: ${(p) => p?.color && p.color.grayish.G10};
  justify-content: space-between;
`;

const SectionPad = styled.div`
  padding: 38px;
  word-break: break-word;
`;

const Pad = styled.div`
  padding: 0 20px;
  word-break: break-word;
`;
const Y = styled.div`
  display: flex;
  justify-content: space-between;
  width: 100%;
  padding: 20px 0;
  align-items: center;
`;
const T = styled.div`
  font-weight: 500;
  font-size: 20px;
  margin: 10px 0;
`;
const B = styled.span<colorProps>`
  font-size: 15px;
  font-weight: bold;
  color: ${(p) => p?.color && p.color.grayish.G10};
`;
const P = styled.div<colorProps>`
  font-weight: regular;
  font-size: 15px;
  color: ${(p) => p?.color && p.color.grayish.G100};
`;
const D = styled.div<colorProps>`
  color: ${(p) => p?.color && p.color.grayish.G50};
  margin: 10px 0 30px;
`;

const Assignee = styled.div<colorProps>`
  margin-left: 3px;
  font-weight: 500;
  cursor: pointer;

  &:hover {
    color: ${(p) => p?.color && p.color.pureBlack};
  }
`;

const ButtonRow = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
`;

const GithubIcon = styled.div`
  height: 20px;
  width: 20px;
  position: relative;
  top: -6px;
  margin-left: 20px;
  cursor: pointer;
`;

const LoomIcon = styled.div`
  height: 20px;
  width: 20px;
  position: relative;
  top: -6px;
  margin-left: 20px;
  cursor: pointer;
`;

const GithubIconMobile = styled.div`
  height: 20px;
  width: 20px;
  position: relative;
  margin-left: 20px;
  cursor: pointer;
`;

const LoomIconMobile = styled.div`
  height: 20px;
  width: 20px;
  position: relative;
  margin-left: 20px;
  cursor: pointer;
`;

interface ImageProps {
  readonly src?: string;
}
const Img = styled.div<ImageProps>`
  background-image: url('${(p) => p.src}');
  background-position: center;
  background-size: cover;
  position: relative;
  width: 22px;
  height: 22px;
`;

const Creator = styled.div`
  min-width: 892px;
  max-width: 892px;
  min-height: 768px;
  max-height: 100vh;
  display: flex;
  justify-content: space-between;
`;

const NormalUser = styled.div`
  min-width: 892px;
  max-width: 892px;
  min-height: 100vh;
  display: flex;
  justify-content: space-between;
`;

const CreatorDescription = styled.div<styleProps>`
  min-width: 600px;
  max-width: 600px;
  min-height: 768px;
  height: 100vh;
  border-right: ${(p) =>
    p?.paid ? `3px solid ${p?.color?.primaryColor.P400}` : `1px solid ${p?.color.grayish.G700}`};
  background: ${(p) => p?.color && p.color.pureWhite};
  padding: 48px 0px 0px 48px;
  .DescriptionUpperContainerNormalView {
    padding-right: 28px;
  }
  .CreatorDescriptionOuterContainerCreatorView {
    padding-right: 28px;
  }
  .CreatorDescriptionInnerContainerCreatorView {
    display: flex;
    justify-content: space-between;
    .CreatorDescriptionExtraButton {
      min-width: 250px;
      max-width: 250px;
      min-height: 40px;
      max-height: 40px;
      display: flex;
      justify-content: space-between;
    }
  }
`;

const Profile = styled.div`
  // padding-top: 48px;
`;

const TitleBox = styled.div<colorProps>`
  margin-top: 24px;
  font-family: 'Barlow';
  font-style: normal;
  font-weight: 600;
  font-size: 22px;
  line-height: 26px;
  display: flex;
  align-items: center;
  color: ${(p) => p?.color && p.color.text1};
`;

const DescriptionBox = styled.div<colorProps>`
  padding-right: 44px;
  margin-right: 5px;
  min-height: 540px;
  max-height: 560px;
  overflow-y: scroll;
  overflow-wrap: anywhere;
  font-family: Barlow;
  font-weight: 400;
  font-size: 15px;
  line-height: 25px;
  color: ${(p) => p?.color && p.color.black500};
  .deliverablesContainer {
    margin-top: 23px;
    .deliverablesHeading {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 700;
      font-size: 13px;
      line-height: 25px;
      letter-spacing: 0.06em;
      text-transform: uppercase;
      color: ${(p) => p?.color && p.color.black500};
    }
    .deliverablesDesc {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 400;
      font-size: 15px;
      line-height: 25px;
      color: ${(p) => p?.color && p.color.black500};
    }
  }
  ::-webkit-scrollbar {
    width: 6px;
    height: 6px;
  }
  ::-webkit-scrollbar-thumb {
    background: ${(p) => p?.color && p.color.grayish.G300} !important;
    height: 80px;
  }
  ::-webkit-scrollbar-track-piece {
    height: 80px;
  }
`;

const AssigneeProfile = styled.div<colorProps>`
  min-width: 292px;
  max-width: 292px;
  min-height: 768px;
  max-height: 100vh;
  background: ${(p) => p?.color && p.color.pureWhite};
  display: flex;
  flex-direction: column;
  .buttonSet {
    display: flex;
    flex-direction: column;
    flex: 1;
  }
`;

interface BountyPriceContainerProps {
  margin_top?: string;
}

const BountyPriceContainer = styled.div<BountyPriceContainerProps>`
  padding-left: 37px;
  margin-top: ${(p) => p.margin_top};
`;

interface codingLangProps {
  background?: string;
  border?: string;
  color?: string;
  styledColors?: any;
}

const LanguageContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  width: 80%;
  margin-top: 16px;
  margin-bottom: 23.25px;
`;

const CodingLabels = styled.div<codingLangProps>`
  padding: 0px 8px;
  border: ${(p) => (p.border ? p?.border : `1px solid ${p?.styledColors.pureBlack}`)};
  color: ${(p) => (p.color ? p?.color : `${p?.styledColors.pureBlack}`)};
  background: ${(p) => (p.background ? p?.background : `${p?.styledColors.pureWhite}`)};
  border-radius: 4px;
  overflow: hidden;
  max-height: 22.75px;
  min-height: 22.75px;
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-right: 4px;
  .LanguageText {
    font-size: 13px;
    fontweight: 500;
    text-align: center;
    font-family: Barlow;
    line-height: 16px;
  }
`;

const DividerContainer = styled.div`
  padding: 32px 36.5px;
`;

interface containerProps {
  color?: any;
  unAssignedBackgroundImage?: string;
  assignedBackgroundImage?: string;
  unassigned_border?: string;
  grayish_G200?: string;
}

const UnassignedPersonProfile = styled.div<containerProps>`
  min-width: 228px;
  min-height: 57.6px;
  display: flex;
  padding-top: 0px;
  padding-left: 28px;
  margin-top: 43px;
  .UnassignedPersonContainer {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 57.6px;
    width: 57.6px;
    border-radius: 50%;
  }
  .UnassignedPersonalDetailContainer {
    margin-left: 25px;
    display: flex;
    align-items: center;
  }
  .BountyProfileOuterContainerCreatorView {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
  }
  .AssigneeCloseButtonContainer {
    margin-left: 6px;
    margin-top: 5px;
    align-self: center;
    height: 22px;
    width: 22px;
    cursor: pointer;
  }
`;

const AutoCompleteContainer = styled.div<colorProps>`
  overflow: hidden;
  z-index: 10;
  padding: 25px 53px 6px 53px;
  .autoCompleteHeaderText {
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 800;
    font-size: 26px;
    line-height: 36px;
    color: ${(p) => p.color && p.color.text2};
    height: 44px;
    margin-bottom: 11px;
  }
`;

const BottomButtonContainer = styled.div`
  margin-bottom: 20px;
`;

const AdjustAmountContainer = styled.div<colorProps>`
  min-height: 460px;
  max-height: 460px;
  min-width: 440px;
  max-width: 440px;
  border-radius: 10px;
  background: ${(p) => p.color && p.color.pureWhite};
  padding-top: 32px;
  .TopHeader {
    max-height: 48px;
    height: 100%;
    display: flex;
    align-items: center;
    margin-left: 25px;
    cursor: pointer;
    .imageContainer {
      height: 48px;
      width: 48px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    .TopHeaderText {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 500;
      font-size: 15px;
      line-height: 18px;
      letter-spacing: 0.06em;
      text-transform: uppercase;
      color: ${(p) => p.color && p.color.black500};
    }
  }
  .Header {
    height: 32px;
    margin-left: 70px;
    .HeaderText {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 800;
      font-size: 36px;
      line-height: 43px;
      display: flex;
      align-items: center;
      text-align: center;
      color: ${(p) => p.color && p.color.black500};
    }
  }
  .AssignedProfile {
    height: 184px;
    margin-top: 30px;
    padding: 0px 31px 0px 38px;
    .InputContainer {
      display: flex;
      flex-direction: row;
      align-items: center;
      .InputContainerLeadingText {
        font-family: Barlow;
        font-style: normal;
        font-weight: 400;
        font-size: 14px;
        line-height: 17px;
        display: flex;
        align-items: center;
        color: ${(p) => (p.color ? p.color.grayish.G100 : '')};
        margin-right: 7px;
      }
      .InputContainerTextField {
        width: 296px;
        background: ${(p) => p?.color && p?.color?.pureWhite};
        border: 1px solid ${(p) => p.color && p.color.grayish.G600};
        color: ${(p) => p.color && p.color.pureBlack};
      }
      .InputContainerEndingText {
        font-family: Barlow;
        font-style: normal;
        font-weight: 400;
        font-size: 14px;
        line-height: 17px;
        display: flex;
        align-items: center;
        color: ${(p) => p.color && p.color.grayish.G100};
        margin-left: 14px;
      }
    }
    .USDText {
      font-family: Barlow;
      font-style: normal;
      font-weight: 500;
      font-size: 13px;
      line-height: 16px !important;
      display: flex;
      align-items: center;
      color: ${(p) => p.color && p.color.grayish.G100};
      margin-left: 42px;
      height: 32px;
    }
  }
  .BottomButton {
    margin-top: 20px;
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    padding-right: 36px;
  }
`;

const AwardsContainer = styled.div<colorProps>`
  width: 622px;
  height: 100vh;
  min-height: 768px;
  max-height: 100vh;
  background: ${(p) => p.color && p.color.pureWhite};
  display: flex;
  flex-direction: column;
  .header {
    min-height: 159px;
    max-height: 159px;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    border-bottom: 1px solid ${(p) => p.color && p.color.grayish.G600};
    box-shadow: 0px 1px 4px ${(p) => p.color && p.color.black80};
    .headerTop {
      height: 48px;
      display: flex;
      align-items: center;
      margin: 32px 0px 0px 25px;
      cursor: pointer;
      .imageContainer {
        height: 48px;
        width: 48px;
        display: flex;
        justify-content: center;
        align-items: center;
      }
      .TopHeaderText {
        font-family: 'Barlow';
        font-style: normal;
        font-weight: 500;
        font-size: 15px;
        line-height: 18px;
        letter-spacing: 0.06em;
        text-transform: uppercase;
        color: ${(p) => p.color && p.color.black500};
      }
    }
    .headerText {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 800;
      font-size: 36px;
      line-height: 43px;
      display: flex;
      align-items: center;
      color: ${(p) => p.color && p.color.black500};
      margin-left: 73px;
      margin-bottom: 48px;
    }
  }
  .AwardContainer {
    min-height: 481px;
    height: 100%;
    display: grid;
    grid-template-columns: 0.5fr 0.5fr;
    gap: 10px;
    place-content: flex-start;
    overflow-y: scroll;
    margin-left: 63px;
    user-select: none;
    cursor: pointer;
    .RadioImageContainer {
      display: flex;
      flex-direction: row;
      height: 65px;
      width: 248px;
      align-items: center;
      padding-left: 9px;
      margin-top: 14px;
      border-radius: 6px;
      input[type='radio'] {
        border: 1px solid ${(p) => p.color && p.color.grayish.G500};
        border-radius: 2px;
        -webkit-appearance: none;
      }
      input[type='radio']:checked {
        background: url('/static/Checked.svg');
        background-repeat: no-repeat;
        border-radius: 2px;
        border: none;
      }
    }
    .awardImageContainer {
      height: 40px;
      width: 40px;
      margin-left: 13px;
    }
    .awardLabelText {
      margin-left: 15px;
      font-family: Barlow;
      font-weight: 500;
      font-size: 13px;
      line-height: 15px;
      color: ${(p) => p.color && p.color.grayish.G05};
    }
  }
`;

const AwardBottomContainer = styled.div<colorProps>`
  height: 129px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-top: 1px solid ${(p) => p.color && p.color.grayish.G600};
  box-shadow: 0px -1px 4px ${(p) => p.color && p.color.black80};
`;

const PaidStatusPopover = styled.div<colorProps>`
  position: absolute;
  background: transparent;
  height: 70px;
  width: 222px;
  right: 54px;
  top: 120px;
  background-image: url('/static/paid_popover_triangle.svg');
  background-size: 16px 16px;
  background-repeat: no-repeat;
  background-position: 16% 0%;
  filter: drop-shadow(0px 1px 20px rgba(0, 0, 0, 0.15));

  .PaidStatusContainer {
    height: 65px;
    width: 222px;
    background: ${(p) => p.color && p.color.green1};
    margin-top: 5px;
    padding: 18px 0px 0px 21px;
    display: flex;
    flex-direction: row;
    cursor: pointer;
    .imageContainer {
      width: 31px;
      height: 31px;
    }
    .PaidStatus {
      font-family: Barlow;
      font-style: normal;
      font-weight: 700;
      font-size: 17px;
      line-height: 15px;
      color: ${(p) => p.color && p.color.pureWhite};
      margin-top: 6px;
      margin-left: 18px;
      user-select: none;
    }
  }
  .ExtraBadgeInfo {
    display: flex;
    flex-direction: row;
    align-items: center;
    background: ${(p) => p.color && p.color.black400};
    height: 75px;
    width: 222px;
    padding: 14px 0px 0px 19px;
    object-fit: cover;
    border-radius: 0px 0px 6px 6px;
    opacity: ${(p) => (p?.isPaidStatusBadgeInfo ? 1 : 0)};
    transition: all ease 4s;
    .imageContainer {
      position: absolute;
      top: 96px;
      left: 14px;
      height: 15px;
      width: 15px;
      background: ${(p) => p.color && p.color.pureWhite};
      display: flex;
      justify-content: center;
      align-items: center;
      border-radius: 50%;
      border: none;
    }
    .badgeText {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 700;
      font-size: 17px;
      line-height: 15px;
      display: flex;
      align-items: center;
      color: ${(p) => p.color && p.color.pureWhite};
      margin-left: 11px;
    }
  }
`;
