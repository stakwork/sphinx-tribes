import MaterialIcon from '@material/react-material-icon';
import React, { useRef, useState, useLayoutEffect, useEffect } from 'react'
import styled from "styled-components";
import { formatPrice, satToUsd } from '../../../helpers';
import { useIsMobile } from '../../../hooks';
import { Divider, Title, Paragraph, Button } from '../../../sphinxUI';
import GalleryViewer from '../../utils/galleryViewer';
import NameTag from '../../utils/nameTag';
import FavoriteButton from '../../utils/favoriteButton'
import { extractGithubIssue } from '../../../helpers';
import ReactMarkdown from 'react-markdown'
import GithubStatusPill from '../parts/statusPill';
import { useHistory } from 'react-router';
import { useStores } from '../../../store';

export function renderMarkdown(str) {
        return <ReactMarkdown>{str}</ReactMarkdown>
}

export default function WantedSummary(props: any) {
        const { title, description, priceMin, priceMax, url, gallery, person, created, repo, issue, price, type, tribe, paid } = props
        const [envHeight, setEnvHeight] = useState('100%')
        const imgRef: any = useRef(null)
        const history = useHistory()
        const heart = <FavoriteButton />
        const isMobile = useIsMobile()
        const { main } = useStores()

        const [tribeInfo, setTribeInfo]: any = useState(null)
        const [assigneeInfo, setAssigneeInfo]: any = useState(null)

        useLayoutEffect(() => {
                if (imgRef && imgRef.current) {
                        if (imgRef.current?.offsetHeight > 100) {
                                setEnvHeight(imgRef.current?.offsetHeight)
                        }
                }
        }, [imgRef])

        useEffect(() => {
                (async () => {
                        if (props.assignee) {
                                try {
                                        const p = await main.getPersonByPubkey(props.assignee.owner_pubkey)
                                        setAssigneeInfo(p)
                                } catch (e) {
                                        console.log('e', e)
                                }
                        }
                        if (tribe) {
                                try {
                                        const t = await main.getSingleTribeByUn(tribe)
                                        setTribeInfo(t)
                                } catch (e) {
                                        console.log('e', e)
                                }
                        }


                })()
        }, [])

        const viewGithub = <Button
                text={'Original Ticket'}
                color={'white'}
                endingIcon={'launch'}
                iconSize={14}
                style={{ fontSize: 14, height: 48 }}
                onClick={() => {
                        const repoUrl = `https://github.com/${repo}/issues/${issue}`
                        sendToRedirect(repoUrl)
                }}
        />

        const viewTribe = <Button
                text={'View Tribe'}
                color={'white'}
                leadingImgUrl={tribeInfo?.img || ' '}
                endingIcon={'launch'}
                iconSize={14}
                imgStyle={{ marginRight: 10 }}
                style={{ fontSize: 14, height: 48 }}
                onClick={() => {
                        const profileUrl = `https://community.sphinx.chat/t/${tribe}`
                        sendToRedirect(profileUrl)
                }}
        />

        function sendToRedirect(url) {
                let el = document.createElement("a");
                el.href = url;
                el.target = '_blank';
                el.click();
        }


        function renderCodingTask() {
                const { assignee, status } = extractGithubIssue(person, repo, issue)

                let assigneeLabel: any = null

                if (assigneeInfo) {
                        assigneeLabel = (<div style={{ display: 'flex', alignItems: 'center', fontSize: 12, color: '#8E969C', }}>
                                <Img src={assigneeInfo.img || '/static/person_placeholder.png'} style={{ borderRadius: 30 }} />
                                <div style={{ marginLeft: 5, fontWeight: 300 }}>
                                        Owner assigned to
                                </div>
                                <Assignee
                                        onClick={() => {
                                                const profileUrl = `https://community.sphinx.chat/p/${assigneeInfo.owner_pubkey}`
                                                sendToRedirect(profileUrl)
                                        }}
                                        style={{ marginLeft: 3, fontWeight: 500, cursor: 'pointer' }}>
                                        {assigneeInfo.owner_alias}
                                </Assignee>
                        </div>)
                }

                if (isMobile) {
                        return <Pad>
                                <NameTag {...person}
                                        created={created}
                                        widget={'wanted'} />

                                <T>{title}</T>

                                <GithubStatusPill status={status} assignee={assignee} />
                                {assigneeLabel}

                                <ButtonRow style={{ marginTop: 22 }}>
                                        {viewGithub}
                                        {viewTribe}
                                </ButtonRow>
                                <Divider style={{
                                        marginTop: 22
                                }} />
                                <Y>
                                        <P><B>{formatPrice(price)}</B> SAT / <B>{satToUsd(price)}</B> USD</P>
                                        {heart}
                                </Y>
                                <Divider style={{ marginBottom: 22 }} />
                                {/* <Img src={'/static/github_logo.png'} /> */}
                                <D>{renderMarkdown(description)}</D>

                        </Pad>
                }

                return <>
                        {paid && <Img src={'/static/paid_ribbon.svg'} style={{
                                position: 'absolute', top: -1,
                                right: 0, width: 64, height: 72
                        }} />}<Wrap>
                                <div style={{ width: 500, padding: 20, borderRight: '1px solid #DDE1E5', minHeight: '100%' }}>
                                        <MaterialIcon icon={'code'} style={{ marginBottom: 5 }} />
                                        <Paragraph style={{
                                                overflow: 'hidden',
                                                wordBreak: 'normal'
                                        }}>{renderMarkdown(description)}</Paragraph>
                                </div>

                                <div style={{ width: 410, padding: 20, height: envHeight }}>
                                        <Pad>
                                                <div style={{ display: 'flex', width: '100%', justifyContent: 'space-between' }}>
                                                        <NameTag
                                                                style={{ marginBottom: 14 }}
                                                                {...person}
                                                                created={created}
                                                                widget={'wanted'} />
                                                </div>
                                                <Divider style={{ margin: '14px 0 20px' }} />

                                                <Title>{title}</Title>
                                                <div style={{ height: 10 }} />
                                                {assigneeLabel}
                                                <div style={{ height: 10 }} />
                                                <GithubStatusPill status={status} assignee={assignee} style={{ marginTop: 10 }} />
                                                <div style={{ height: 30 }} />
                                                <ButtonRow>
                                                        {viewGithub}
                                                        {viewTribe}
                                                </ButtonRow>

                                                <Divider style={{ margin: '20px 0 0' }} />
                                                <Y>
                                                        <P><B>{formatPrice(price)}</B> SAT / <B>{satToUsd(price)}</B> USD</P>
                                                        {heart}
                                                </Y>
                                                <Divider />

                                        </Pad>
                                </div>

                        </Wrap>
                </>
        }



        if (type === 'coding_task' || type === 'wanted_coding_task') {
                return renderCodingTask()
        }

        if (isMobile) {
                return <>
                        <Pad>
                                <NameTag {...person}
                                        created={created}
                                        widget={'wanted'} />

                                <T>{title || 'No title'}</T>
                                <Divider style={{
                                        marginTop: 22
                                }} />
                                <Y>
                                        <P>{formatPrice(priceMin) || '0'} <B>SAT</B> - {formatPrice(priceMax)} <B>SAT</B></P>
                                        {heart}
                                </Y>
                                <Divider style={{ marginBottom: 22 }} />

                                <D>{renderMarkdown(description)}</D>

                                <GalleryViewer gallery={gallery} showAll={true} selectable={false} wrap={false} big={true} />
                        </Pad>
                </>
        }

        return <Wrap>
                <GalleryViewer
                        innerRef={imgRef}
                        style={{ width: 507, height: 'fit-content' }}
                        gallery={gallery} showAll={false} selectable={false} wrap={false} big={true} />
                <div style={{ width: 316, padding: 20, overflowY: 'auto', height: envHeight }}>
                        <Pad>
                                <NameTag
                                        style={{ marginBottom: 14 }}
                                        {...person}
                                        created={created}
                                        widget={'wanted'} />

                                <Title>{title}</Title>

                                <Divider style={{ marginTop: 22 }} />
                                <Y>
                                        <P>{formatPrice(priceMin) || '0'} <B>SAT</B> - {formatPrice(priceMax) || '0'} <B>SAT</B></P>
                                        {heart}
                                </Y>
                                <Divider style={{ marginBottom: 22 }} />

                                <Paragraph>{renderMarkdown(description)}</Paragraph>
                        </Pad>
                </div>

        </Wrap>

}

const Wrap = styled.div`
display: flex;
width:100%;
height:100%;
min-width:800px;
font-style: normal;
font-weight: 500;
font-size: 24px;

color: #3C3F41;
justify-content:space-between;

`;
const Pad = styled.div`
        padding: 0 20px;
        `;
const Y = styled.div`
        display: flex;
        justify-content:space-between;
        width:100%;
        padding: 20px 0;
        align-items:center;
        `;
const T = styled.div`
        font-weight:bold;
        font-size:20px;
        margin: 10px 0;
        `;
const B = styled.span`
        font-size:15px;
        font-weight:bold;
        color:#3c3f41;
        `;
const P = styled.div`
        font-weight:regular;
        font-size:15px;
        color:#8e969c;
        `;
const D = styled.div`
        color:#5F6368;
        margin: 10px 0 30px;
        `;


const Assignee = styled.div`
        margin-left: 3px;
        font-weight: 500; 
        cursor: pointer;

        &:hover{
                color:#000;
        }
        `;

const ButtonRow = styled.div`
        display: flex;
        justify-content:space-around;
        `;

const Link = styled.div`
        color:blue;
        overflow-wrap:break-word;
        font-size:15px;
        font-weight:300;
        `;


interface ImageProps {
        readonly src?: string;
}
const Img = styled.div<ImageProps>`
                            background-image: url("${(p) => p.src}");
                            background-position: center;
                            background-size: cover;
                            position: relative;
                            width:22px;
                            height:22px;
                            `;