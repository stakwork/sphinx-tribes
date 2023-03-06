/* eslint-disable func-style */
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { formatPrice, satToUsd } from '../../helpers';
import { useIsMobile } from '../../hooks';
import { Divider, Title, Button } from '../../components/common';
import NameTag from '../utils/nameTag';
import { extractGithubIssue, extractGithubIssueFromUrl } from '../../helpers';
import GithubStatusPill from './parts/statusPill';
import { useStores } from '../../store';
import { renderMarkdown } from '../utils/renderMarkdown';
import { EuiButtonIcon, EuiText } from '@elastic/eui';
import { getHost } from '../../host';
import PaidBounty from '../utils/paidBounty';
import Bounties from '../utils/assigned_unassigned_bounties';
import { colors } from '../../colors';

export default function WantedView(props: any) {
  const {
    one_sentence_summary,
    title,
    description,
    priceMin,
    priceMax,
    price,
    gallery,
    person,
    created,
    issue,
    ticketUrl,
    repo,
    type,
    codingLanguage,
    assignee,
    estimate_session_length,
    loomEmbedUrl,
    showModal,
    setDeletePayload,
    onPanelClick
  } = props;
  const titleString = title ?? one_sentence_summary;

  let { show, paid } = props;
  const isMobile = useIsMobile();
  const { ui, main } = useStores();
  const [saving, setSaving] = useState(false);
  const [labels, setLabels] = useState([]);
  const { peopleWanteds } = main;
  const color = colors['light'];

  const isMine = ui.meInfo?.owner_pubkey === person?.owner_pubkey;

  if ('show' in props) {
    // show has a value
  } else {
    // if no value default to true
    show = true;
  }

  if ('paid' in props) {
    // show has no value
  } else {
    // if no value default to false
    paid = false;
  }

  async function setExtrasPropertyAndSave(propertyName: string) {
    if (peopleWanteds) {
      setSaving(true);
      try {
        const targetProperty = props[propertyName];
        const [clonedEx, targetIndex] = await main.setExtrasPropertyAndSave(
          'wanted',
          propertyName,
          created,
          !targetProperty
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
            peopleWantedsClone[indexFromPeopleWanted] = {
              person: person,
              body: clonedEx[targetIndex]
            };
          }
          main.setPeopleWanteds(peopleWantedsClone);
        }
      } catch (e) {
        console.log('e', e);
      }

      setSaving(false);
    }
  }

  useEffect(() => {
    if (codingLanguage) {
      const values = codingLanguage.map((value) => ({ ...value }));
      setLabels(values);
    }
  }, [codingLanguage]);

  function renderTickets() {
    const { status } = ticketUrl
      ? extractGithubIssueFromUrl(person, ticketUrl)
      : extractGithubIssue(person, repo, issue);

    const isClosed = status === 'closed' || paid ? true : false;

    const isCodingTask =
      type === 'coding_task' || type === 'wanted_coding_task' || type === 'freelance_job_request';

    // mobile view
    if (isMobile) {
      return (
        <div style={{ position: 'relative' }} onClick={onPanelClick}>
          {paid && (
            <Img
              src={'/static/paid_ribbon.svg'}
              style={{
                position: 'absolute',
                right: '-2.5px',
                width: '80px',
                height: '80px',
                top: 0
              }}
            />
          )}
          <Wrap isClosed={isClosed} style={{ padding: 15 }}>
            <Body style={{ width: '100%' }} color={color}>
              <div
                style={{
                  display: 'flex',
                  width: '100%',
                  justifyContent: 'space-between'
                }}
              >
                <NameTag
                  {...person}
                  created={created}
                  widget={'wanted'}
                  ticketUrl={ticketUrl}
                  loomEmbedUrl={loomEmbedUrl}
                  style={{
                    margin: 0
                  }}
                />
              </div>
              <DT
                style={{
                  margin: '15px 0'
                }}
              >
                {titleString}
              </DT>

              <div
                style={{
                  display: 'flex',
                  flexDirection: 'row',
                  alignItems: 'center'
                }}
              >
                {isCodingTask && (
                  <GithubStatusPill
                    status={status}
                    assignee={assignee}
                    style={{
                      marginTop: 10
                    }}
                  />
                )}
                {{ ...assignee }.owner_alias && (
                  <div
                    style={{
                      marginTop: '8px'
                    }}
                  >
                    <img
                      src={
                        {
                          ...assignee
                        }.img || '/static/person_placeholder.png'
                      }
                      alt="assignee_img"
                      style={{
                        borderRadius: '50%',
                        height: '16px',
                        width: '16px',
                        margin: '0px 8px'
                      }}
                    />
                    <span
                      onClick={(e) => {
                        e.stopPropagation();
                        window.open(
                          `/p/${
                            {
                              ...assignee
                            }.owner_pubkey
                          }?widget=wanted`,
                          '_blank'
                        );
                      }}
                      style={{
                        fontSize: '12px'
                      }}
                    >
                      {
                        {
                          ...assignee
                        }.owner_alias
                      }
                    </span>
                  </div>
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
                  minHeight: '45px',
                  width: '100%',
                  display: 'flex',
                  flexDirection: 'row',
                  marginTop: '10px',
                  flexWrap: 'wrap'
                }}
              >
                {labels.length > 0 ? (
                  labels.map((x: any) => {
                    return (
                      <>
                        <div
                          style={{
                            display: 'flex',
                            flexWrap: 'wrap',
                            height: 'fit-content',
                            width: 'fit-content',
                            backgroundColor: color.grayish.G1000,
                            border: `1px solid ${color.grayish.G70}`,
                            padding: '0px 14px',
                            borderRadius: '20px',
                            marginRight: '3px',
                            marginBottom: '3px'
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
                  })
                ) : (
                  <>
                    <div
                      style={{
                        minHeight: '50px'
                      }}
                    />
                  </>
                )}
              </div>
              <EyeDeleteTextContainerMobile>
                {priceMin ? (
                  <P
                    color={color}
                    style={{
                      margin: '15px 0 0'
                    }}
                  >
                    <B color={color}>{formatPrice(priceMin)}</B>~
                    <B color={color}>{formatPrice(priceMax)}</B> SAT /{' '}
                    <B color={color}>{satToUsd(priceMin)}</B>~
                    <B color={color}>{satToUsd(priceMax)}</B> USD
                  </P>
                ) : (
                  <P
                    color={color}
                    style={{
                      margin: '15px 0 0'
                    }}
                  >
                    <B color={color}>{formatPrice(price)}</B> SAT /{' '}
                    <B color={color}>{satToUsd(price)}</B> USD
                  </P>
                )}
                <EyeDeleteContainerMobile>
                  <div
                    style={{
                      width: '40px'
                    }}
                  >
                    {
                      //  if my own, show this option to show/hide
                      isMine && (
                        <Button
                          icon={show ? 'visibility' : 'visibility_off'}
                          disable={saving}
                          submitting={saving}
                          iconStyle={{
                            color: color.grayish.G20,
                            fontSize: 20
                          }}
                          style={{
                            minWidth: 24,
                            maxWidth: 24,
                            minHeight: 20,
                            height: 20,
                            padding: 0,
                            background: `${color.pureWhite}`
                          }}
                          onClick={(e) => {
                            e.stopPropagation();
                            setExtrasPropertyAndSave('show');
                          }}
                        />
                      )
                    }
                  </div>
                  {ui?.meInfo?.isSuperAdmin && (
                    <EuiButtonIcon
                      onClick={(e) => {
                        e.stopPropagation();
                        showModal();
                        setDeletePayload({
                          created: created,
                          host: getHost(),
                          pubkey: person.owner_pubkey
                        });
                      }}
                      iconType="trash"
                      aria-label="Next"
                      size="s"
                      style={{
                        color: `${color.pureBlack}`,
                        background: `${color.pureWhite}`
                      }}
                    />
                  )}
                </EyeDeleteContainerMobile>
              </EyeDeleteTextContainerMobile>
            </Body>
          </Wrap>
        </div>
      );
    }

    if (props?.fromBountyPage) {
      return (
        <>
          {paid ? (
            <BountyBox color={color}>
              <PaidBounty
                {...person}
                onPanelClick={onPanelClick}
                assignee={assignee}
                created={created}
                ticketUrl={ticketUrl}
                loomEmbedUrl={loomEmbedUrl}
                title={titleString}
                codingLanguage={labels}
                priceMin={priceMin}
                priceMax={priceMax}
                price={price}
                sessionLength={estimate_session_length}
                description={description}
              />
            </BountyBox>
          ) : (
            <BountyBox color={color}>
              <Bounties
                onPanelClick={onPanelClick}
                person={person}
                assignee={assignee}
                created={created}
                ticketUrl={ticketUrl}
                loomEmbedUrl={loomEmbedUrl}
                title={titleString}
                codingLanguage={labels}
                priceMin={priceMin}
                priceMax={priceMax}
                price={price}
                sessionLength={estimate_session_length}
                description={description}
              />
            </BountyBox>
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
              height: 72
            }}
          />
        )}

        <DWrap isClosed={isClosed} color={color}>
          <Pad style={{ padding: 20, minHeight: 410 }}>
            <div
              style={{
                display: 'flex',
                width: '100%',
                justifyContent: 'space-between'
              }}
            >
              <NameTag
                {...person}
                created={created}
                widget={'wanted'}
                ticketUrl={ticketUrl}
                loomEmbedUrl={loomEmbedUrl}
              />
            </div>
            <Divider style={{ margin: '10px 0' }} />
            <DT>{titleString}</DT>
            <div
              style={{
                display: 'flex',
                flexDirection: 'row',
                alignItems: 'center'
              }}
            >
              {isCodingTask ? (
                <GithubStatusPill
                  status={status}
                  assignee={assignee}
                  style={{
                    marginTop: 10
                  }}
                />
              ) : (
                <div
                  style={{
                    minHeight: '36px'
                  }}
                />
              )}
              {{ ...assignee }.owner_alias && (
                <div
                  style={{
                    marginTop: '8px'
                  }}
                >
                  <img
                    src={
                      {
                        ...assignee
                      }.img || '/static/person_placeholder.png'
                    }
                    alt="assignee_img"
                    style={{
                      borderRadius: '50%',
                      height: '16px',
                      width: '16px',
                      margin: '0px 8px'
                    }}
                  />
                  <span
                    onClick={(e) => {
                      e.stopPropagation();
                      window.open(
                        `/p/${
                          {
                            ...assignee
                          }.owner_pubkey
                        }?widget=wanted`,
                        '_blank'
                      );
                    }}
                    style={{
                      fontSize: '12px'
                    }}
                  >
                    {
                      {
                        ...assignee
                      }.owner_alias
                    }
                  </span>
                </div>
              )}
            </div>

            <div
              style={{
                minHeight: '48px',
                width: '100%',
                display: 'flex',
                flexDirection: 'row',
                marginTop: '10px',
                flexWrap: 'wrap'
              }}
            >
              {labels.length > 0 ? (
                labels.map((x: any) => {
                  return (
                    <>
                      <div
                        style={{
                          display: 'flex',
                          flexWrap: 'wrap',
                          height: 'fit-content',
                          width: 'fit-content',
                          backgroundColor: color.grayish.G1000,
                          border: `1px solid ${color.grayish.G70}`,
                          padding: '0px 14px',
                          borderRadius: '20px',
                          marginRight: '3px',
                          marginBottom: '3px'
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
                })
              ) : (
                <>
                  <div
                    style={{
                      minHeight: '50px'
                    }}
                  />
                </>
              )}
            </div>
            <Divider
              style={{
                margin: isCodingTask || gallery ? '22px 0' : '0 0 22px'
              }}
            />
            <DescriptionCodeTask color={color}>
              {renderMarkdown(description)}
              {gallery && (
                <div
                  style={{
                    display: 'flex',
                    flexWrap: 'wrap'
                  }}
                >
                  {gallery.map((val, index) => {
                    return (
                      <div
                        key={index}
                        style={{
                          height: '48px',
                          width: '48px',
                          padding: '0px 2px',
                          borderRadius: '6px',
                          overflow: 'hidden'
                        }}
                      >
                        <img src={val} alt="wanted preview" height={'100%'} width={'100%'} />
                      </div>
                    );
                  })}
                </div>
              )}
            </DescriptionCodeTask>
          </Pad>

          <Divider style={{ margin: 0 }} />

          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              padding: '10px 20px',
              minHeight: '100px'
            }}
          >
            <Pad
              style={{
                flexDirection: 'row',
                justifyContent: 'space-between'
              }}
            >
              {priceMin ? (
                <P color={color}>
                  <B color={color}>{formatPrice(priceMin)}</B>~
                  <B color={color}>{formatPrice(priceMax)}</B> SAT /{' '}
                  <B color={color}>{satToUsd(priceMin)}</B>~
                  <B color={color}>{satToUsd(priceMax)}</B> USD
                </P>
              ) : (
                <P color={color}>
                  <B color={color}>{formatPrice(price)}</B> SAT /{' '}
                  <B color={color}>{satToUsd(price)}</B> USD
                </P>
              )}

              <div
                style={{
                  width: '40px'
                }}
              >
                {
                  //  if my own, show this option to show/hide
                  isMine && (
                    <Button
                      icon={show ? 'visibility' : 'visibility_off'}
                      disable={saving}
                      submitting={saving}
                      iconStyle={{
                        color: color.grayish.G300,
                        fontSize: 20
                      }}
                      style={{
                        minWidth: 24,
                        maxWidth: 24,
                        minHeight: 20,
                        height: 20,
                        padding: 0,
                        background: `${color.pureWhite}`
                      }}
                      onClick={(e) => {
                        e.stopPropagation();
                        setExtrasPropertyAndSave('show');
                      }}
                    />
                  )
                }
              </div>
            </Pad>
            <div
              style={{
                display: 'flex',
                flexDirection: 'row',
                justifyContent: 'space-between',
                alignItems: 'center'
              }}
            >
              <EuiText
                style={{
                  fontSize: '14px',
                  color: color.grayish.G300,
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
              {ui?.meInfo?.isSuperAdmin && (
                <EuiButtonIcon
                  onClick={(e) => {
                    e.stopPropagation();
                    showModal();
                    setDeletePayload({
                      created: created,
                      host: getHost(),
                      pubkey: person.owner_pubkey
                    });
                  }}
                  iconType="trash"
                  aria-label="Next"
                  size="s"
                  style={{
                    color: `${color.pureBlack}`,
                    background: `${color.pureWhite}`
                  }}
                />
              )}
            </div>
          </div>
        </DWrap>
      </>
    );
  }

  return renderTickets();
}

interface WrapProps {
  isClosed?: boolean;
  color?: any;
}

interface styledProps {
  color?: any;
}

const BountyBox = styled.div<styledProps>`
  min-height: 160px;
  max-height: 160px;
  width: 1100px;
  box-shadow: 0px 1px 6px ${(p) => p?.color && p?.color.black100};
  border: none;
`;

const DWrap = styled.div<WrapProps>`
  display: flex;
  flex: 1;
  height: 100%;
  min-height: 510px;
  flex-direction: column;
  width: 100%;
  min-width: 100%;
  max-height: 510px;
  font-style: normal;
  font-weight: 500;
  font-size: 17px;
  line-height: 23px;
  color: ${(p) => p?.color && p?.color.grayish.G10} !important;
  letter-spacing: 0px;
  justify-content: space-between;
  opacity: ${(p) => (p.isClosed ? '0.5' : '1')};
  filter: ${(p) => (p.isClosed ? 'grayscale(1)' : 'grayscale(0)')};
`;

const Wrap = styled.div<WrapProps>`
  display: flex;
  justify-content: flex-start;
  opacity: ${(p) => (p.isClosed ? '0.5' : '1')};
  filter: ${(p) => (p.isClosed ? 'grayscale(1)' : 'grayscale(0)')};
`;

const B = styled.span<styledProps>`
  font-size: 14px;
  font-weight: bold;
  color: ${(p) => p?.color && p?.color.grayish.G10};
`;
const P = styled.div<styledProps>`
  font-weight: regular;
  font-size: 14px;
  color: ${(p) => p?.color && p?.color.grayish.G100};
`;

const Body = styled.div<styledProps>`
  font-size: 15px;
  line-height: 20px;
  padding: 10px;
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  color: ${(p) => p?.color && p?.color.grayish.G05};
  overflow: hidden;
  min-height: 132px;
`;

const Pad = styled.div`
  display: flex;
  flex-direction: column;
`;

const DescriptionCodeTask = styled.div<styledProps>`
  margin-bottom: 10px;

  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 13px;
  line-height: 20px;
  color: ${(p) => p?.color && p?.color.grayish.G50};
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 6;
  -webkit-box-orient: vertical;
  height: 120px;
  max-height: 120px;
`;
const DT = styled(Title)`
  margin-bottom: 9px;
  max-height: 52px;
  min-height: 43.5px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  /* Primary Text 1 */

  font-family: 'Roboto';
  font-style: normal;
  font-weight: 500;
  font-size: 17px;
  line-height: 23px;
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

const EyeDeleteTextContainerMobile = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
`;

const EyeDeleteContainerMobile = styled.div`
  margin-top: 10px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
`;
