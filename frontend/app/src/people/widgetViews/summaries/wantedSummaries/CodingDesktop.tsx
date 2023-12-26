/* eslint-disable func-style */
import React from 'react';
import { EuiText } from '@elastic/eui';
import { CodingViewProps } from 'people/interfaces';
import { Divider, Paragraph, Title } from '../../../../components/common';
import StatusPill from '../../parts/StatusPill';
import LoomViewerRecorder from '../../../utils/LoomViewerRecorder';
import { colors } from '../../../../config/colors';
import { renderMarkdown } from '../../../utils/RenderMarkdown';
import { formatPrice, satToUsd } from '../../../../helpers';
import { AddToFavorites, CopyLink, ShareOnTwitter, ViewTribe, ViewGithub } from './Components';
import { ButtonRow, Y, P, B, Img, Wrap, SectionPad, LoomIcon, GithubIcon } from './style';

export default function DesktopView(props: CodingViewProps) {
  const {
    paid,
    titleString,
    labels,
    price,
    description,
    envHeight,
    estimated_session_length,
    loomEmbedUrl,
    ticketUrl,
    assignee,
    assigneeLabel,
    nametag,
    actionButtons,
    owner_id,
    created
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
            <Title>{titleString}</Title>
            <div style={{ display: 'flex', marginTop: 12 }}>
              <StatusPill assignee={assignee} style={{ marginRight: 25 }} paid={paid} />

              {assigneeLabel}
              {ticketUrl && (
                <GithubIcon
                  onClick={(e: any) => {
                    e.stopPropagation();
                    window.open(ticketUrl, '_blank');
                  }}
                >
                  <img height={'100%'} width={'100%'} src="/static/github_logo.png" alt="github" />
                </GithubIcon>
              )}
              {loomEmbedUrl && (
                <LoomIcon
                  onClick={(e: any) => {
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
            <div
              style={{
                minHeight: '60px',
                width: '100%',
                display: 'flex',
                flexDirection: 'row'
              }}
            >
              {(labels ?? []).map((x: any) => (
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
              ))}
            </div>
          </SectionPad>
          <Divider />
          <SectionPad>
            <Y style={{ padding: 0 }}>
              <P color={color}>
                <B color={color}>{formatPrice(price || 0)}</B> SAT /{' '}
                <B color={color}>{satToUsd(price || 0)}</B> USD
              </P>
            </Y>
          </SectionPad>

          <Divider />

          <SectionPad>
            <ButtonRow>
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

            {actionButtons}
          </SectionPad>
        </div>
      </Wrap>
    </>
  );
}
