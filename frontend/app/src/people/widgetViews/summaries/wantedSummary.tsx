import MaterialIcon from '@material/react-material-icon';
import React, { useRef, useState, useLayoutEffect, useEffect, useCallback } from 'react';
import styled from 'styled-components';
import { formatPrice, satToUsd } from '../../../helpers';
import { useIsMobile } from '../../../hooks';
import { Divider, Title, Paragraph, Button, Modal } from '../../../sphinxUI';
import GalleryViewer from '../../utils/galleryViewer';
import NameTag from '../../utils/nameTag';
import FavoriteButton from '../../utils/favoriteButton';
import { extractGithubIssue, extractGithubIssueFromUrl } from '../../../helpers';
import ReactMarkdown from 'react-markdown';
import GithubStatusPill from '../parts/statusPill';
import { useStores } from '../../../store';
import Form from '../../../form';
import { sendBadgeSchema } from '../../../form/schema';
import remarkGfm from 'remark-gfm';
import LoomViewerRecorder from '../../utils/loomViewerRecorder';
import { renderMarkdown } from '../../utils/renderMarkdown';
import { useLocation } from 'react-router-dom';
import { EuiText } from '@elastic/eui';

function useQuery() {
  const { search } = useLocation();

  return React.useMemo(() => new URLSearchParams(search), [search]);
}

export default function WantedSummary(props: any) {
  const {
    title,
    description,
    priceMin,
    priceMax,
    url,
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
    assignee
  } = props;
  let {} = props;
  const [envHeight, setEnvHeight] = useState('100%');
  const imgRef: any = useRef(null);

  const isMobile = useIsMobile();
  const { main, ui } = useStores();
  const { peopleWanteds } = main;

  const [tribeInfo, setTribeInfo]: any = useState(null);
  const [assigneeInfo, setAssigneeInfo]: any = useState(null);
  const [saving, setSaving]: any = useState('');
  const [isCopied, setIsCopied] = useState(false);
  const [owner_idURL, setOwnerIdURL] = useState('');
  const [createdURL, setCreatedURL] = useState('');

  const [showBadgeAwardDialog, setShowBadgeAwardDialog] = useState(false);

  const isMine = ui.meInfo?.owner_pubkey === person?.owner_pubkey;

  const [labels, setLabels] = useState([]);

  useLayoutEffect(() => {
    if (imgRef && imgRef.current) {
      if (imgRef.current?.offsetHeight > 100) {
        setEnvHeight(imgRef.current?.offsetHeight);
      }
    }
  }, [imgRef]);

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
  }, []);

  const searchParams = useQuery();

  useEffect(() => {
    const owner_id = searchParams.get('owner_id');
    const created = searchParams.get('created');
    setOwnerIdURL(owner_id ?? '');
    setCreatedURL(created ?? '');
  }, [owner_idURL, createdURL]);

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

  const handleCopyUrl = useCallback(() => {
    const el = document.createElement('input');
    el.value = window.location.href;
    document.body.appendChild(el);
    el.select();
    document.execCommand('copy');
    document.body.removeChild(el);
    setIsCopied(true);
  }, [isCopied]);

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
      onClick={() => {}}
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
        const twitterLink = `https://twitter.com/intent/tweet?text=Hey, I created a new ticket on Sphinx community.%0A${title} %0A&url=https://community.sphinx.chat/p?owner_id=${owner_idURL}%26created${createdURL} %0A%0A&hashtags=${
          labels && labels.map((x: any) => x.label)
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
              color: '#8E969C',
              marginTop: isMobile ? 20 : 0,
              marginLeft: '-16px'
            }}
          >
            <Img
              src={assigneeInfo.img || '/static/person_placeholder.png'}
              style={{ borderRadius: 30 }}
            />

            <Assignee
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
              color: '#8E969C',
              marginLeft: '16px'
            }}
          >
            <Img
              src={assigneeInfo.img || '/static/person_placeholder.png'}
              style={{ borderRadius: 30 }}
            />

            <Assignee
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
                color: '#8e969c',
                fontWeight: '500'
              }}
            >
              {estimate_session_length && 'Session:'}{' '}
              <span
                style={{
                  fontWeight: '500',
                  color: '#000'
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
              {labels.length > 0 &&
                labels.map((x: any) => {
                  return (
                    <>
                      <div
                        style={{
                          display: 'flex',
                          flexWrap: 'wrap',
                          height: '22px',
                          width: 'fit-content',
                          backgroundColor: '#cfcfcf',
                          border: '1px solid #909090',
                          padding: '3px 10px',
                          borderRadius: '20px',
                          marginRight: '3px',
                          boxShadow: '1px 1px #909090'
                        }}
                      >
                        <div
                          style={{
                            fontSize: '10px',
                            color: '#202020'
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
              <P>
                <B>{formatPrice(price)}</B> SAT / <B>{satToUsd(price)}</B> USD
              </P>
              {heart}
            </Y>
            <Divider style={{ marginBottom: 20 }} />
            <D>{renderMarkdown(description)}</D>
          </Pad>
        </div>
      );
    }

    // desktop view
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
        <Wrap>
          <div
            style={{
              width: 700,
              borderRight: '1px solid #DDE1E5',
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
                    color: '#8e969c',
                    fontWeight: '500'
                  }}
                >
                  {estimate_session_length && 'Session:'}{' '}
                  <span
                    style={{
                      fontWeight: '500',
                      color: '#000'
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
                {labels.length > 0 &&
                  labels.map((x: any) => {
                    return (
                      <>
                        <div
                          style={{
                            display: 'flex',
                            flexWrap: 'wrap',
                            height: '22px',
                            minWidth: 'fit-content',
                            backgroundColor: '#cfcfcf',
                            border: '1px solid #909090',
                            padding: '3px 10px',
                            borderRadius: '20px',
                            marginRight: '3px',
                            boxShadow: '1px 1px #909090'
                          }}
                        >
                          <div
                            style={{
                              fontSize: '10px',
                              color: '#202020'
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
                <P>
                  <B>{formatPrice(price)}</B> SAT / <B>{satToUsd(price)}</B> USD
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
            <P>
              {formatPrice(priceMin) || '0'} <B>SAT</B> - {formatPrice(priceMax)} <B>SAT</B>
            </P>
            {heart}
          </Y>
          <Divider style={{ marginBottom: 22 }} />

          <D>{renderMarkdown(description)}</D>
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
      <Wrap>
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
              <P>
                {formatPrice(priceMin) || '0'} <B>SAT</B> - {formatPrice(priceMax) || '0'}{' '}
                <B>SAT</B>
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

const Wrap = styled.div`
  display: flex;
  width: 100%;
  height: 100%;
  min-width: 800px;
  font-style: normal;
  font-weight: 500;
  font-size: 24px;
  color: #3c3f41;
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
const B = styled.span`
  font-size: 15px;
  font-weight: bold;
  color: #3c3f41;
`;
const P = styled.div`
  font-weight: regular;
  font-size: 15px;
  color: #8e969c;
`;
const D = styled.div`
  color: #5f6368;
  margin: 10px 0 30px;
`;

const Assignee = styled.div`
  margin-left: 3px;
  font-weight: 500;
  cursor: pointer;

  &:hover {
    color: #000;
  }
`;

const ButtonRow = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
`;

const Link = styled.div`
  color: blue;
  overflow-wrap: break-word;
  font-size: 15px;
  font-weight: 300;
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
