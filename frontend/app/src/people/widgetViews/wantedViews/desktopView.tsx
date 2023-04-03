/* eslint-disable func-style */
import React from 'react';
import { Img, P, B, DT, DWrap, DescriptionCodeTask, Pad } from './style';
import { EuiButtonIcon, EuiText } from '@elastic/eui';
import GithubStatusPill from '../parts/statusPill';
import { colors } from '../../../config/colors';
import NameTag from '../../utils/nameTag';
import { useStores } from '../../../store';
import { formatPrice, satToUsd } from '../../../helpers';
import { Button, Divider } from '../../../components/common';
import { getHost } from '../../../config/host';
import { renderMarkdown } from '../../utils/renderMarkdown';
import { observer } from 'mobx-react-lite';

export default observer(DesktopView);

function DesktopView(props: any) {
  const {
    description,
    priceMin,
    priceMax,
    price,
    person,
    created,
    ticketUrl,
    gallery,
    assignee,
    estimate_session_length,
    loomEmbedUrl,
    showModal,
    setDeletePayload,
    key,
    setExtrasPropertyAndSave,
    saving,
    labels,
    isClosed,
    onPanelClick,
    status,
    isCodingTask,
    show,
    paid,
    isMine,
    titleString
  } = props;

  const { ui } = useStores();
  const color = colors['light'];

  return (
    <div key={key} onClick={onPanelClick}>
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
                <B color={color}>{satToUsd(priceMin)}</B>~<B color={color}>{satToUsd(priceMax)}</B>{' '}
                USD
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
    </div>
  );
}
