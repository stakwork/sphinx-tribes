import React from 'react'
import styled from "styled-components";
import { Divider } from '../../sphinxUI';
import QrBar from '../utils/QrBar'
import ReactMarkdown from 'react-markdown'
import { useStores } from '../../store';
import { useHistory, useLocation } from 'react-router-dom'

export function renderMarkdown(str) {
    return <ReactMarkdown>{str}</ReactMarkdown>
}

export default function AboutView(props: any) {
    const { ui, main } = useStores()
    const history = useHistory()
    let { ownerTribes } = main
    const { price_to_meet, description, extras, twitter_confirmed, owner_pubkey } = props
    const { twitter, github, coding_languages, tribes } = extras || {}
    let tag = ''
    let githubTag = ''
    let codingLanguages = ''
    let myTribes: any = null
    if (twitter && twitter[0] && twitter[0].value) tag = twitter[0].value
    if (github && github[0] && github[0].value) githubTag = github[0].value
    if (coding_languages && coding_languages[0] && coding_languages[0].value) codingLanguages = coding_languages[0].value
    if (tribes && tribes[0] && tribes[0].value) {
        myTribes = tribes[0].value
    }

    // for testing
    // let ownerTribes: any = [{
    //     unique_name: 'none',
    //     name: 'None',
    //     img: '/static/twitter2.png'
    // },
    // {
    //     unique_name: 'this',
    //     name: 'This',
    //     img: '/static/twitter2.png'
    // },
    // {
    //     unique_name: 'that',
    //     name: 'That',
    //     img: '/static/twitter2.png'
    // }]

    return <Wrap>
        <Row>
            <div>Price to Connect:</div>
            <div style={{ fontWeight: 'bold', color: '#000' }}>{price_to_meet}</div>
        </Row>

        <Divider />

        <QrBar value={owner_pubkey} />

        {githubTag &&
            <>
                <Divider />

                <Row style={{ justifyContent: 'flex-start', fontSize: 14 }}>
                    <Img src={'/static/github_logo.png'} />
                    <a href={`https://github.com/${githubTag}`} target='_blank'>https://github.com/{githubTag}</a>
                </Row>
            </>
        }

        {tag && <>
            <Divider />
            <Row>
                {/* <T>For Normies</T> */}
                <I>
                    <div style={{ width: 4 }} />
                    <Icon source={`/static/twitter2.png`} />
                    <Tag>@{tag}</Tag>
                    {twitter_confirmed ?
                        <Badge>VERIFIED</Badge> :
                        <Badge style={{ background: '#b0b7bc' }}>PENDING</Badge>
                    }

                </I>
            </Row>
        </>}

        {codingLanguages && <>
            <Divider />
            <Row style={{ justifyContent: 'flex-start', fontSize: 14 }}>
                {codingLanguages}
            </Row>
        </>}

        {myTribes &&
            <>
                <Divider />
                <T style={{ height: 20 }}>My Tribes</T>
                <Grow >
                    {myTribes.map((t, i) => {
                        const thisTribe = ownerTribes && ownerTribes.find(f => f.unique_name === t.value)
                        if (!thisTribe) return <div key={i} />
                        return (<TribeRow key={i + 'mytribe'}
                            onClick={() => history.push(`/t/${thisTribe?.unique_name}`)}>
                            <Img src={thisTribe?.img} />
                            <div>{thisTribe?.name}</div>
                        </TribeRow>)
                    })}
                </Grow>
            </>
        }

        <Divider />



        <D>{renderMarkdown(description)}</D>

    </Wrap>

}
const Badge = styled.div`
display:flex;
justify-content:center;
align-items:center;
margin-left:10px;
height:20px;
color:#ffffff;
background: #1DA1F2;
border-radius: 32px;
font-weight: bold;
font-size: 8px;
line-height: 9px;
padding 0 10px;
`;
const TribeRow = styled.div`
display: flex;
align-items: center;
cursor: pointer;
&:hover{
    color:#000;
}
`;
const Wrap = styled.div`
display: flex;
flex-direction:column;
width:100%;
`;
const I = styled.div`
display:flex;
align-items:center;

`;

const Tag = styled.div`
font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 14px;
line-height: 26px;
/* or 173% */

display: flex;
align-items: center;

/* Main bottom icons */

color: #5F6368;
`

const Row = styled.div`
display:flex;
justify-content:space-between;
height:48px;
align-items:center;
font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 15px;
line-height: 48px;
/* identical to box height, or 320% */

display: flex;
align-items: center;

/* Secondary Text 4 */

color: #8E969C;

`;

const Grow = styled.div`
display:flex;
justify-content:flex-start;
flex-direction:column;
min-height:28px;
font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 14px;
line-height: 28px;
margin-bottom:8px;
color: #8E969C;

`;
const T = styled.div`
font-family: Roboto;
font-style: normal;
font-weight: bold;
font-size: 10px;
line-height: 26px;
/* or 260% */

letter-spacing: 0.3px;
text-transform: uppercase;

/* Text 2 */

color: #3C3F41;
margin-top:5px;
margin-bottom:5px;
`;

const D = styled.div`

margin:35px 0 10px 0;
font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 15px;
line-height: 20px;
/* or 133% */


/* Main bottom icons */

color: #5F6368;

`;

interface IconProps {
    source: string;
}

const Icon = styled.div<IconProps>`
                    background-image: ${p => `url(${p.source})`};
                    width:16px;
                    height:13px;
                    margin-right:8px;
                    background-position: center; /* Center the image */
                    background-repeat: no-repeat; /* Do not repeat the image */
                    background-size: contain; /* Resize the background image to cover the entire container */
                    border-radius:5px;
                    overflow:hidden;
                `;

interface ImageProps {
    readonly src?: string;
}
const Img = styled.div<ImageProps>`
                                        background-image: url("${(p) => p.src}");
                                        background-position: center;
                                        background-size: cover;
                                        position: relative;
                                        width:18px;
                                        height:18px;
                                        margin-left:2px;
                                        margin-right: 10px;
                                        border-radius:5px;
                                        `;