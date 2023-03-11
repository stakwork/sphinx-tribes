/* eslint-disable func-style */
import React from "react";
import { ButtonRow, Pad, GithubIconMobile, T, Y, P, D, B, LoomIconMobile } from "./style";
import { Heart, AddToFavorites, CopyLink, ShareOnTwitter, ViewTribe, ViewGithub } from "./components";
import GithubStatusPill from '../../parts/statusPill';
import { EuiText } from '@elastic/eui';
import { Divider } from '../../../../components/common';
import LoomViewerRecorder from '../../../utils/loomViewerRecorder';
import { colors } from '../../../../config/colors';
import { renderMarkdown } from '../../../utils/renderMarkdown';
import { formatPrice, satToUsd } from '../../../../helpers';

export default function MobileView(props: any) {
    const {
        description,
        ticketUrl,
        price,
        loomEmbedUrl,
        estimate_session_length,
        assignee,
        titleString,
        nametag,
        assigneeLabel,
        labels,
        actionButtons,
        status
    } = props;
    const color = colors['light'];

    return (
        <div style={{ padding: 20, overflow: 'auto' }}>
            <Pad>
                {nametag}

                <T>{titleString}</T>

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
                    <ViewGithub {...props} />
                    <ViewTribe   {...props} />
                    <AddToFavorites  {...props} />
                    <CopyLink   {...props} />
                    <ShareOnTwitter  {...props} />
                </ButtonRow>

                {actionButtons}

                <LoomViewerRecorder readOnly loomEmbedUrl={loomEmbedUrl} style={{ marginBottom: 20 }} />

                <Divider />
                <Y>
                    <P color={color}>
                        <B color={color}>{formatPrice(price)}</B> SAT /{' '}
                        <B color={color}>{satToUsd(price)}</B> USD
                    </P>
                    <Heart />
                </Y>
                <Divider style={{ marginBottom: 20 }} />
                <D color={color}>{renderMarkdown(description)}</D>
            </Pad>
        </div>
    );
}