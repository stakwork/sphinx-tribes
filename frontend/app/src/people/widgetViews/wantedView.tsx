import React, { useState } from 'react'
import styled from "styled-components";
import { formatPrice, satToUsd } from '../../helpers';
import { useIsMobile } from '../../hooks';
import GalleryViewer from '../utils/galleryViewer';
import { Divider, Title, Button } from '../../sphinxUI';
import NameTag from '../utils/nameTag';
import { extractGithubIssue } from '../../helpers';
import GithubStatusPill from './parts/statusPill';
import { useStores } from '../../store';

export default function WantedView(props: any) {
    let { title, description, priceMin, priceMax, price, url, gallery, person, created, issue, repo, type, show, paid } = props
    const isMobile = useIsMobile()
    const { ui, main } = useStores()
    const [saving, setSaving] = useState(false)
    const { peopleWanteds } = main

    const isMine = ui.meInfo?.owner_pubkey === person?.owner_pubkey

    if ('show' in props) {
        // show has a value
    } else {
        // if no value default to true
        show = true
    }

    if ('paid' in props) {
        // show has no value
    } else {
        // if no value default to false
        paid = false
    }

    async function setExtrasPropertyAndSave(propertyName: string) {
        if (peopleWanteds) {
            setSaving(true)
            try {
                const targetProperty = props[propertyName]
                const [clonedEx, targetIndex] = await main.setExtrasPropertyAndSave(
                    'wanted',
                    propertyName,
                    created,
                    !targetProperty)

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
                        peopleWantedsClone[indexFromPeopleWanted] = {
                            person: person,
                            body: clonedEx[targetIndex]
                        }
                    }
                    main.setPeopleWanteds(peopleWantedsClone)
                }
            } catch (e) {
                console.log('e', e)
            }

            setSaving(false)
        }
    }


    function renderTickets() {
        const { assignee, status } = extractGithubIssue(person, repo, issue)

        const isClosed = ((status === 'closed') || paid) ? true : false

        console.log('assignee', assignee)

        const isCodingTask = type === 'coding_task' || type === 'wanted_coding_task'

        if (isMobile) {
            return <>
                {paid && <Img src={'/static/paid_ribbon.svg'} style={{
                    position: 'absolute', top: -1,
                    right: 0, width: 64, height: 72
                }} />}
                <Wrap isClosed={isClosed} style={{ padding: 15 }}>
                    <Body style={{ width: '100%' }}>
                        <div style={{ display: 'flex', width: '100%', justifyContent: 'space-between', }}>
                            <NameTag {...person} created={created} widget={'wanted'} style={{ margin: 0 }} />
                        </div>
                        <DT style={{ margin: '15px 0' }}>{title}</DT>
                        <div style={{ width: '100%', display: 'flex', justifyContent: 'space-between', margin: '5px 0' }}>
                            {isCodingTask && <GithubStatusPill status={status} assignee={assignee} />}
                        </div>
                        {priceMin ?
                            <P style={{ margin: '15px 0 0' }}><B>{formatPrice(priceMin)}</B>~<B>{formatPrice(priceMax)}</B> SAT / <B>{satToUsd(priceMin)}</B>~<B>{satToUsd(priceMax)}</B> USD</P>
                            : <P style={{ margin: '15px 0 0' }}><B>{formatPrice(price)}</B> SAT / <B>{satToUsd(price)}</B> USD</P>
                        }

                    </Body>
                </Wrap>
            </>
        }

        return <>
            {paid && <Img src={'/static/paid_ribbon.svg'} style={{
                position: 'absolute', top: -1,
                right: 0, width: 64, height: 72
            }} />}

            <DWrap isClosed={isClosed}>
                <Pad style={{ padding: 20, height: 410 }}>
                    <div style={{ display: 'flex', width: '100%', justifyContent: 'space-between' }}>
                        <NameTag {...person} created={created} widget={'wanted'} />
                    </div>

                    <Divider style={{ margin: '10px 0' }} />

                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        {isCodingTask &&
                            <Img src={'/static/github_logo2.png'} style={{ width: 77, height: 43 }} />
                        }
                        {isMine && <Button
                            style={{
                                height: 30,
                            }}
                            color={'primary'}
                            text={paid ? 'Mark Unpaid' : 'Mark Paid'}
                            onClick={e => {
                                e.stopPropagation()
                                setExtrasPropertyAndSave('paid')
                            }} />
                        }
                    </div>


                    <DT>{title}</DT>

                    {isCodingTask &&
                        <GithubStatusPill status={status} assignee={assignee} style={{ marginTop: 10 }} />
                    }

                    <Divider style={{ margin: isCodingTask ? '22px 0' : '0 0 22px' }} />

                    <DescriptionCodeTask>{description}</DescriptionCodeTask>

                </Pad>
                <Divider style={{ margin: 0 }} />
                <Pad style={{ padding: 20, flexDirection: 'row', justifyContent: 'space-between' }}>
                    {priceMin ?
                        <P><B>{formatPrice(priceMin)}</B>~<B>{formatPrice(priceMax)}</B> SAT / <B>{satToUsd(priceMin)}</B>~<B>{satToUsd(priceMax)}</B> USD</P>
                        : <P><B>{formatPrice(price)}</B> SAT / <B>{satToUsd(price)}</B> USD</P>
                    }


                    <div>
                        {
                            //  if my own, show this option to show/hide
                            isMine &&
                            <Button
                                icon={show ? 'visibility' : 'visibility_off'}
                                disable={saving}
                                submitting={saving}
                                iconStyle={{ color: '#555', fontSize: 20 }}
                                style={{
                                    minWidth: 24, width: 24, minHeight: 20,
                                    height: 20, padding: 0, background: '#fff',
                                }}
                                onClick={(e) => {
                                    e.stopPropagation()
                                    setExtrasPropertyAndSave('show')
                                }}
                            />
                        }

                    </div>
                </Pad>
            </DWrap>
        </>
    }

    function getMobileView() {
        return <Wrap>
            <GalleryViewer cover gallery={gallery} selectable={false} wrap={false} big={false} showAll={false} />
            <Body>
                <NameTag {...person} created={created} widget={'wanted'} style={{ margin: 0 }} />
                <T>{title}</T>
                <D>{description}</D>
                <P><B>{formatPrice(priceMin) || '0'} - {formatPrice(priceMax)}</B> SAT</P>
            </Body>
        </Wrap>
    }

    function getDesktopView() {
        return <DWrap>
            <GalleryViewer
                cover
                showAll={false}
                big={true}
                wrap={false}
                selectable={true}
                gallery={gallery}
                style={{ maxHeight: 276, overflow: 'hidden' }} />

            <Pad style={{ padding: 20 }}>
                <NameTag {...person} created={created} widget={'wanted'} />
                <DT>{title}</DT>
                <DD style={{ maxHeight: gallery ? 40 : '' }}>{description}</DD>
            </Pad>
            <Divider style={{ margin: 0 }} />
            <Pad style={{ padding: 20, }}>
                <P><B>{formatPrice(priceMin) || '0'} - {formatPrice(priceMax)}</B> SAT</P>
            </Pad>
        </DWrap>
    }

    return renderTickets()

    // muted for now
    // if (isMobile) {
    //     return getMobileView()
    // }

    // return getDesktopView()
}

interface WrapProps {
    isClosed?: boolean;
}

const DWrap = styled.div<WrapProps>`
display: flex;
flex:1;
height:100%;
min-height:100%;
flex-direction:column;
width:100%;
min-width:100%;
max-height:471px;
font-style: normal;
font-weight: 500;
font-size: 17px;
line-height: 23px;
color: #3C3F41 !important;
letter-spacing:0px;
justify-content:space-between;
opacity:${p => p.isClosed ? '0.5' : '1'};
filter: ${p => p.isClosed ? 'grayscale(1)' : 'grayscale(0)'};


`;

const Wrap = styled.div<WrapProps>`
display: flex;
justify-content:flex-start;
opacity:${p => p.isClosed ? '0.5' : '1'};
filter: ${p => p.isClosed ? 'grayscale(1)' : 'grayscale(0)'};
`;


const T = styled.div`
font-weight:bold;
overflow:hidden;
line-height: 20px;
text-overflow: ellipsis;
display: -webkit-box;
-webkit-line-clamp: 2;
-webkit-box-orient: vertical;

font-family: 'Roboto';
font-style: normal;
font-weight: 500;
font-size: 17px;
line-height: 23px;
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
overflow:hidden;
line-height:18px;
text-overflow: ellipsis;
display: -webkit-box;
-webkit-line-clamp: 2;
-webkit-box-orient: vertical;
`;


const Body = styled.div`
font-size: 15px;
line-height: 20px;
/* or 133% */
padding:10px;
display: flex;
flex-direction:column;
justify-content: space-around;

/* Primary Text 1 */

color: #292C33;
overflow:hidden;
min-height:132px;
`;

const Pad = styled.div`
display:flex;
flex-direction:column;
padding:10px;
`;



const DD = styled.div`
margin-bottom:10px;
overflow:hidden;

font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 13px;
line-height: 20px;
/* or 154% */

/* Main bottom icons */

color: #5F6368;

overflow: hidden;
text-overflow: ellipsis;
display: -webkit-box;
-webkit-line-clamp: 2;
-webkit-box-orient: vertical;

`;

const DescriptionCodeTask = styled.div`
margin-bottom:10px;

font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 13px;
line-height: 20px;
color: #5F6368;
overflow: hidden;
text-overflow: ellipsis;
display: -webkit-box;
-webkit-line-clamp: 6;
-webkit-box-orient: vertical;
height: 120px;
`
const DT = styled(Title)`
margin-bottom:9px;
max-height:52px;
overflow:hidden;
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
                        background-image: url("${(p) => p.src}");
                        background-position: center;
                        background-size: cover;
                        position: relative;
                        width:22px;
                        height:22px;
                        `;