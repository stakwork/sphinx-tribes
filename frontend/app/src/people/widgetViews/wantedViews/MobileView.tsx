/* eslint-disable func-style */
import React from 'react';
import { EuiButtonIcon, EuiText } from '@elastic/eui';
import { observer } from 'mobx-react-lite';
import StatusPill from '../parts/StatusPill';
import { colors } from '../../../config/colors';
import NameTag from '../../utils/NameTag';
import { useStores } from '../../../store';
import { formatPrice, satToUsd } from '../../../helpers';
import { Button } from '../../../components/common';
import { getHost } from '../../../config/host';
import {
  Img,
  EyeDeleteContainerMobile,
  Wrap,
  Body,
  P,
  B,
  DT,
  EyeDeleteTextContainerMobile
} from './style';

function MobileView(props: any) {
  const {
    priceMin,
    priceMax,
    price,
    person,
    created,
    ticket_url,
    assignee,
    estimated_session_length,
    loomEmbedUrl,
    showModal,
    setDeletePayload,
    onPanelClick,
    setExtrasPropertyAndSave,
    saving,
    labels,
    isCodingTask,
    show,
    paid,
    isMine,
    titleString
  } = props;

  const { ui } = useStores();
  const color = colors['light'];

  return (
    <div
      style={{ borderBottom: '2px solid #EBEDEF', position: 'relative' }}
      onClick={onPanelClick}
      key={titleString}
    >
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
      <Wrap style={{ padding: 15 }}>
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
              ticketUrl={ticket_url}
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
              <StatusPill
                assignee={assignee}
                style={{
                  marginTop: 10
                }}
                paid={paid}
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
                  onClick={(e: any) => {
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
            {estimated_session_length && 'Session:'}{' '}
            <span
              style={{
                fontWeight: '500',
                color: color.pureBlack
              }}
            >
              {estimated_session_length ?? ''}
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
              labels.map((x: any) => (
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
              ))
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
                <B color={color}>{satToUsd(priceMin)}</B>~<B color={color}>{satToUsd(priceMax)}</B>{' '}
                USD
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
                      disabled={saving}
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
                      onClick={(e: any) => {
                        e.stopPropagation();
                        setExtrasPropertyAndSave('show');
                      }}
                    />
                  )
                }
              </div>
              {ui?.meInfo?.isSuperAdmin && (
                <EuiButtonIcon
                  onClick={(e: any) => {
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
export default observer(MobileView);
