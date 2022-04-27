import MaterialIcon from '@material/react-material-icon';
import React, { useRef, useState, useLayoutEffect, useEffect } from 'react'
import styled from "styled-components";
import { formatPrice, satToUsd } from '../../../helpers';
import { useIsMobile } from '../../../hooks';
import { Divider, Title, Paragraph, Button, Modal } from '../../../sphinxUI';
import GalleryViewer from '../../utils/galleryViewer';
import NameTag from '../../utils/nameTag';
import FavoriteButton from '../../utils/favoriteButton'
import { extractGithubIssue, extractGithubIssueFromUrl } from '../../../helpers';
import ReactMarkdown from 'react-markdown'
import GithubStatusPill from '../parts/statusPill';
import { useStores } from '../../../store';
import Form from '../../../form';
import { sendBadgeSchema } from '../../../form/schema';
import remarkGfm from 'remark-gfm'
import LoomViewerRecorder from '../../utils/loomViewerRecorder'

export function renderMarkdown(markdown) {

        return <ReactMarkdown children={markdown} remarkPlugins={[remarkGfm]}
                components={{
                        code({ node, inline, className, children, ...props }) {
                                return (
                                        <code className={className} {...props}>
                                                {children}
                                        </code>
                                )
                        },
                        img({ className, ...props }) {
                                return (
                                        <img className={className}
                                                style={{ width: '100%' }}
                                                {...props} />
                                )
                        }
                }} />

}

export default function WantedSummary(props: any) {
        const { title, description, priceMin, priceMax, url, ticketUrl, gallery, person, created, repo, issue, price, type, tribe, paid, badgeRecipient, loomEmbedUrl } = props
        let { } = props
        const [envHeight, setEnvHeight] = useState('100%')
        const imgRef: any = useRef(null)

        const isMobile = useIsMobile()
        const { main, ui } = useStores()
        const { peopleWanteds } = main

        const [tribeInfo, setTribeInfo]: any = useState(null)
        const [assigneeInfo, setAssigneeInfo]: any = useState(null)
        const [saving, setSaving]: any = useState('')

        const [showBadgeAwardDialog, setShowBadgeAwardDialog] = useState(false)

        const isMine = ui.meInfo?.owner_pubkey === person?.owner_pubkey

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



        async function setExtrasPropertyAndSave(propertyName: string, value: any) {
                if (peopleWanteds) {
                        setSaving(propertyName)
                        try {
                                const [clonedEx, targetIndex] = await main.setExtrasPropertyAndSave(
                                        'wanted',
                                        propertyName,
                                        created,
                                        value)

                                // saved? ok update in wanted list if found
                                const peopleWantedsClone: any = [...peopleWanteds]
                                const indexFromPeopleWanted = peopleWantedsClone.findIndex(f => {
                                        let val = f.body || {}
                                        return ((f.person.owner_pubkey === ui.meInfo?.owner_pubkey) && val.created === created)
                                })

                                // if we found it in the wanted list, update in people wanted list
                                if (indexFromPeopleWanted > -1) {
                                        // if it should be hidden now, remove it from the list
                                        if ('show' in clonedEx[targetIndex] && clonedEx[targetIndex].show === false) {
                                                peopleWantedsClone.splice(indexFromPeopleWanted, 1)
                                        } else {

                                                // gotta update person extras! this is what is used for summary viewer
                                                const personClone: any = person
                                                personClone.extras['wanted'][targetIndex] = clonedEx[targetIndex]

                                                peopleWantedsClone[indexFromPeopleWanted] = {
                                                        person: personClone,
                                                        body: clonedEx[targetIndex]
                                                }
                                        }

                                        main.setPeopleWanteds(peopleWantedsClone)
                                }
                        } catch (e) {
                                console.log('e', e)
                        }

                        setSaving('')
                }
        }

        async function sendBadge(body: any) {
                const { recipient, badge } = body

                setSaving('badgeRecipient')
                try {
                        if (badge?.amount < 1) {
                                alert("You don't have any of the selected badge")
                                throw new Error("You don't have any of the selected badge")
                        }

                        // first get the user's liquid address
                        const recipientDetails = await main.getPersonByPubkey(recipient.owner_pubkey)

                        const liquidAddress = recipientDetails?.extras?.liquid && recipientDetails?.extras?.liquid[0]?.value

                        if (!liquidAddress) {
                                alert('This user has not provided an L-BTC address')
                                throw new Error('This user has not provided an L-BTC address')
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
                        }

                        const r = await main.sendBadgeOnLiquid(pack)

                        if (r.ok) {
                                await setExtrasPropertyAndSave('badgeRecipient', recipient.owner_pubkey)
                                setShowBadgeAwardDialog(false)
                        } else {
                                alert(r.statusText)
                                throw new Error(r.statusText)
                        }


                } catch (e) {
                        console.log(e)
                }

                setSaving('')

        }

        const heart = <FavoriteButton />

        const viewGithub = <Button
                text={'Original Ticket'}
                color={'white'}
                endingIcon={'launch'}
                iconSize={14}
                style={{ fontSize: 14, height: 48, marginRight: 10 }}
                onClick={() => {
                        const repoUrl = ticketUrl ? ticketUrl : `https://github.com/${repo}/issues/${issue}`
                        sendToRedirect(repoUrl)
                }}
        />

        const viewTribe = tribeInfo ? <Button
                text={'View Tribe'}
                color={'white'}
                leadingImgUrl={tribeInfo?.img || ' '}
                endingIcon={'launch'}
                iconSize={14}
                imgStyle={{ marginRight: 10 }}
                style={{ fontSize: 14, height: 48, marginRight: 10 }}
                onClick={() => {
                        const profileUrl = `https://community.sphinx.chat/t/${tribe}`
                        sendToRedirect(profileUrl)
                }}
        /> : <div />


        //  if my own, show this option to show/hide
        const markPaidButton = <Button
                color={'primary'}
                iconSize={14}
                style={{ fontSize: 14, height: 48, minWidth: 130, marginRight: 10 }}
                endingIcon={'paid'}
                text={paid ? 'Mark Unpaid' : 'Mark Paid'}
                loading={saving === 'paid'}
                onClick={e => {
                        e.stopPropagation()
                        setExtrasPropertyAndSave('paid', !paid)
                }} />

        const awardBadgeButton = !badgeRecipient && <Button
                color={'primary'}
                iconSize={14}
                endingIcon={'offline_bolt'}
                style={{ fontSize: 14, height: 48, minWidth: 130, marginRight: 10 }}
                text={'Award Badge'}
                loading={saving === 'badgeRecipient'}
                onClick={e => {
                        e.stopPropagation()
                        if (!badgeRecipient) {
                                setShowBadgeAwardDialog(true)
                        }

                }} />

        const actionButtons = isMine && (
                <ButtonRow style={{
                        marginBottom: 20
                }}>
                        {showBadgeAwardDialog ?
                                <>
                                        <Form
                                                loading={saving === 'badgeRecipient'}
                                                smallForm
                                                buttonsOnBottom
                                                wrapStyle={{ padding: 0, margin: 0 }}
                                                close={() => setShowBadgeAwardDialog(false)}
                                                onSubmit={(e) => {
                                                        sendBadge(e)
                                                }}
                                                submitText={'Send Badge'}
                                                schema={sendBadgeSchema}
                                        />
                                </> :
                                <>
                                        {markPaidButton}
                                        {awardBadgeButton}
                                </>
                        }
                </ButtonRow>
        )


        function sendToRedirect(url) {
                let el = document.createElement("a");
                el.href = url;
                el.target = '_blank';
                el.click();
        }


        function renderCodingTask() {
                const { assignee, status } = ticketUrl ? extractGithubIssueFromUrl(person, ticketUrl) : extractGithubIssue(person, repo, issue)

                let assigneeLabel: any = null

                if (assigneeInfo) {
                        assigneeLabel = (<div style={{ display: 'flex', alignItems: 'center', fontSize: 12, color: '#8E969C', marginTop: 20 }}>
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
                        return <div style={{ padding: 20, overflow: 'auto' }}>
                                <Pad>
                                        <NameTag {...person}
                                                created={created}
                                                widget={'wanted'} />

                                        <T>{title}</T>

                                        <GithubStatusPill status={status} assignee={assignee} />
                                        {assigneeLabel}

                                        <ButtonRow style={{ margin: '20px 0' }}>
                                                {viewGithub}
                                                {viewTribe}
                                        </ButtonRow>

                                        {actionButtons}

                                        <LoomViewerRecorder
                                                readOnly
                                                style={{ marginTop: 20 }}
                                                loomEmbedUrl={loomEmbedUrl} />

                                        <Divider style={{
                                                marginTop: 22
                                        }} />
                                        <Y>
                                                <P><B>{formatPrice(price)}</B> SAT / <B>{satToUsd(price)}</B> USD</P>
                                                {heart}
                                        </Y>
                                        <Divider style={{ marginBottom: 22 }} />
                                        <D>{renderMarkdown(description)}</D>
                                </Pad>
                        </div>
                }

                return <>
                        {paid && <Img src={'/static/paid_ribbon.svg'} style={{
                                position: 'absolute', top: -1,
                                right: 0, width: 64, height: 72, zIndex: 100, pointerEvents: 'none'
                        }} />}<Wrap>
                                <div style={{ width: 500, padding: 40, borderRight: '1px solid #DDE1E5', minHeight: '100%', overflow: 'auto' }}>
                                        <MaterialIcon icon={'code'} style={{ marginBottom: 5 }} />
                                        <Paragraph style={{
                                                overflow: 'hidden',
                                                wordBreak: 'normal'
                                        }}>{renderMarkdown(description)}</Paragraph>
                                </div>

                                <div style={{ width: 410, padding: '40px 20px', height: envHeight, overflow: 'auto' }}>
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
                                                <GithubStatusPill status={status} assignee={assignee} style={{ marginTop: 25 }} />
                                                {assigneeLabel}

                                                <div style={{ height: 10 }} />
                                                <ButtonRow style={{ margin: '20px 0' }}>
                                                        {viewGithub}
                                                        {viewTribe}
                                                </ButtonRow>

                                                {actionButtons}

                                                <LoomViewerRecorder
                                                        readOnly
                                                        style={{ marginTop: 20 }}
                                                        loomEmbedUrl={loomEmbedUrl} />

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
                return <div style={{ padding: 20, overflow: 'auto' }}>
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
                </div>
        }

        return <Wrap>
                <GalleryViewer
                        innerRef={imgRef}
                        style={{ width: 507, height: 'fit-content' }}
                        gallery={gallery} showAll={false} selectable={false} wrap={false} big={true} />
                <div style={{ width: 316, padding: '40px 20px', overflowY: 'auto', height: envHeight }}>
                        <Pad>
                                < NameTag
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
                </div >

        </Wrap >

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
        word-break: break-word;
        `;
const Y = styled.div`
        display: flex;
        justify-content:space-between;
        width:100%;
        padding: 20px 0;
        align-items:center;
        `;
const T = styled.div`
        font-weight:500;
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
        // justify-content:space-around;
        flex-wrap:wrap;
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