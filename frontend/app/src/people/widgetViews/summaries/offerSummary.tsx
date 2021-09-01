import React from 'react'
import styled from "styled-components";
import { Offer } from '../../../form/inputs/widgets/interfaces';
import { formatPrice } from '../../../helpers';
import { useIsMobile } from '../../../hooks';
import { Divider, Title, Paragraph } from '../../../sphinxUI';
import GalleryViewer from '../../utils/galleryViewer';

export default function OfferSummary(props: Offer) {
        const { gallery, title, description, price } = props
        const isMobile = useIsMobile()
        if (isMobile) {
                return <>
                        <Pad>
                                <Y>
                                        <div>Price</div>
                                        <P>{formatPrice(price)} <B>sat</B></P>
                                </Y>
                                <Divider style={{
                                        marginBottom: 10
                                }} />


                                <T>{title || 'No title'}</T>
                                <D>{description || 'No description'}</D>
                        </Pad>
                        <Divider />
                        <GalleryViewer gallery={gallery} showAll={true} selectable={false} wrap={false} big={true} />
                </>
        }

        return <Wrap>
                <GalleryViewer
                        style={{ width: 507 }}
                        gallery={gallery} showAll={false} selectable={false} wrap={false} big={true} />
                <div style={{ width: 316, height: '100%', padding: 20, paddingTop: 30 }}>
                        <Pad>
                                <Title>{title}</Title>
                                <Divider style={{ marginTop: 10 }} />
                                <Y>
                                        <P>{formatPrice(price)} <B>sat</B></P>
                                </Y>
                                <Divider style={{ marginBottom: 10 }} />

                                <Paragraph>{description}</Paragraph>
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
        padding:0 20px;
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
        margin: 10px 0;
        `;
