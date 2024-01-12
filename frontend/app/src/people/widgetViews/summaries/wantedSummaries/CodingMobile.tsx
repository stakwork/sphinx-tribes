/* eslint-disable func-style */
import React from 'react';
import { EuiText } from '@elastic/eui';
import { CodingViewProps } from 'people/interfaces';
import StatusPill from '../../parts/StatusPill';
import { Divider } from '../../../../components/common';
import LoomViewerRecorder from '../../../utils/LoomViewerRecorder';
import { colors } from '../../../../config/colors';
import { renderMarkdown } from '../../../utils/RenderMarkdown';
import { formatPrice, satToUsd } from '../../../../helpers';
import {
  Heart,
  AddToFavorites,
  CopyLink,
  ShareOnTwitter,
  ViewTribe,
  ViewGithub
} from './Components';
import { ButtonRow, Pad, Img, GithubIconMobile, T, Y, P, D, B, LoomIconMobile } from './style';

export default function MobileView(props: CodingViewProps) {
  const {
    description,
    ticket_url,
    price,
    loomEmbedUrl,
    estimated_session_length,
    assignee,
    titleString,
    nametag,
    assigneeLabel,
    labels,
    payBounty,
    showPayBounty,
    owner_id,
    created,
    markPaidOrUnpaid,
    paid
  } = props;

  const color = colors['light'];

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
      <div style={{ padding: 20, overflow: 'auto', height: 'calc(100% - 60px)' }}>
        <Pad>
          {nametag}
          <T>{titleString}</T>

          <div
            style={{
              display: 'flex',
              flexDirection: 'row'
            }}
          >
            <StatusPill assignee={assignee} paid={paid} />
            {assigneeLabel}
            {ticket_url && (
              <GithubIconMobile
                onClick={(e: any) => {
                  e.stopPropagation();
                  window.open(ticket_url, '_blank');
                }}
              >
                <img height={'100%'} width={'100%'} src="/static/github_logo.png" alt="github" />
              </GithubIconMobile>
            )}
            {loomEmbedUrl && (
              <LoomIconMobile
                onClick={(e: any) => {
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
              width: '100%',
              display: 'flex',
              flexDirection: 'row',
              marginTop: '10px',
              minHeight: '60px'
            }}
          >
            {(labels ?? []).map((x: any) => (
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
            ))}
          </div>

          <div style={{ height: 10 }} />
          {showPayBounty && payBounty}
          <ButtonRow style={{ margin: '10px 0' }}>
            <ViewGithub {...props} />
            <ViewTribe {...props} />
            <AddToFavorites {...props} />
            <CopyLink {...props} />
            <ShareOnTwitter
              issueCreated={created}
              ownerPubkey={owner_id}
              labels={labels}
              titleString={titleString}
            />
          </ButtonRow>

          {markPaidOrUnpaid}
          <LoomViewerRecorder readOnly loomEmbedUrl={loomEmbedUrl} style={{ marginBottom: 20 }} />

          <Divider />
          <Y>
            <P color={color}>
              <B color={color}>{formatPrice(price || 0)}</B> SAT /{' '}
              <B color={color}>{satToUsd(price || 0)}</B> USD
            </P>
            <Heart />
          </Y>
          <Divider style={{ marginBottom: 20 }} />
          <D color={color}>{renderMarkdown(description)}</D>
        </Pad>
      </div>
    </>
  );
}
