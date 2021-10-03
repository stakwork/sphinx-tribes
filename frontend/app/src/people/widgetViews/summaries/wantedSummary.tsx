import MaterialIcon from '@material/react-material-icon';
import React, { useRef, useState, useLayoutEffect } from 'react'
import styled from "styled-components";
import { formatPrice } from '../../../helpers';
import { useIsMobile } from '../../../hooks';
import { Divider, Title, Paragraph, Button } from '../../../sphinxUI';
import GalleryViewer from '../../utils/galleryViewer';
import NameTag from '../../utils/nameTag';
import FavoriteButton from '../../utils/favoriteButton'
import { extractGithubIssue } from '../../../helpers';
import ReactMarkdown from 'react-markdown'
import GithubStatusPill from '../parts/statusPill';
import { useHistory } from 'react-router';

export function renderMarkdown(str) {
        return <ReactMarkdown>{str}</ReactMarkdown>
}

export default function WantedSummary(props: any) {
        const { title, description, priceMin, priceMax, url, gallery, person, created, repo, issue, price, type, tribe } = props
        const [envHeight, setEnvHeight] = useState('100%')
        const imgRef: any = useRef(null)
        const history = useHistory()
        const heart = <FavoriteButton />
        const isMobile = useIsMobile()

        useLayoutEffect(() => {
                if (imgRef && imgRef.current) {
                        if (imgRef.current?.offsetHeight > 100) {
                                setEnvHeight(imgRef.current?.offsetHeight)
                        }
                }
        }, [imgRef])

        const repoUrl = `github.com/${repo}/issues/${issue}`
        const githubLink = <a href={'https://' + repoUrl} target='_blank'><Link >{repoUrl}</Link></a>
        const addTribe = <Button
                text={'View Tribe'}
                color={'primary'}
                onClick={() => history.push(`/t/${tribe}`)}
        />

        function renderCodingTask() {
                const { assignee, status } = extractGithubIssue(person, repo, issue)

                if (isMobile) {
                        return <Pad>
                                <NameTag {...person}
                                        created={created}
                                        widget={'wanted'} />

                                <T>{title}</T>
                                <div style={{ margin: '5px 0 10px' }}>
                                        {githubLink}
                                </div>
                                <GithubStatusPill status={status} assignee={assignee} />

                                <Status style={{ marginTop: 22 }}>{addTribe}</Status>
                                <Divider style={{
                                        marginTop: 22
                                }} />
                                <Y>
                                        <P>{formatPrice(price)} <B>SAT</B></P>
                                        {heart}
                                </Y>
                                <Divider style={{ marginBottom: 22 }} />
                                {/* <Img src={'/static/github_logo.png'} /> */}
                                <D>{renderMarkdown(description)}</D>

                        </Pad>
                }

                return <Wrap>
                        <div style={{ width: 500, padding: 20, borderRight: '1px solid #DDE1E5', minHeight: '100%' }}>
                                <MaterialIcon icon={'code'} style={{ marginBottom: 5 }} />
                                <Paragraph>{renderMarkdown(description)}</Paragraph>
                        </div>
                        <div style={{ width: 316, padding: 20, overflowY: 'auto', height: envHeight }}>
                                <Pad>
                                        <div style={{ display: 'flex', width: '100%', justifyContent: 'space-between' }}>
                                                <NameTag
                                                        style={{ marginBottom: 14 }}
                                                        {...person}
                                                        created={created}
                                                        widget={'wanted'} />
                                                <Img src={'/static/github_logo.png'} />
                                        </div>


                                        <Title>{title}</Title>
                                        {githubLink}
                                        <GithubStatusPill status={status} assignee={assignee} style={{ marginTop: 10 }} />


                                        <Status style={{ marginTop: 22 }}>{addTribe}</Status>

                                        <Divider style={{ marginTop: 22 }} />
                                        <Y>
                                                <P>{formatPrice(price) || '0'} <B>SAT</B></P>
                                                {heart}
                                        </Y>


                                </Pad>
                        </div>

                </Wrap>
        }



        if (type === 'coding_task') {
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
line-height: 20px;
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
        height:50px;
        align-items:center;
        `;
const T = styled.div`
        font-weight:bold;
        font-size:20px;
        margin: 10px 0;
        `;
const B = styled.span`
        font-weight:300;
        `;
const P = styled.div`
        font-weight:500;
        `;
const D = styled.div`
        color:#5F6368;
        margin: 10px 0 30px;
        `;


const Assignee = styled.div`
        display: flex;
        font-size:12px;
        font-weight:300;
        `;

const Status = styled.div`
        display: flex;
        font-size:12px;
        margin-right:4px;
        font-weight:300;
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