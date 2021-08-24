import React from 'react'
import styled from "styled-components";
import { Offer } from '../../../form/inputs/widgets/interfaces';
import { formatPrice } from '../../../helpers';
import { Divider } from '../../../sphinxUI';
import GalleryViewer from '../../utils/galleryViewer';

export default function OfferSummary(props: Offer) {
        const { gallery, title, description, price } = props

        return <Wrap>

                <GalleryViewer gallery={gallery} selectable={false} wrap={false} big={true} />
                <Pad>
                        <Y>
                                <div>Price</div>
                                <P>{formatPrice(price)} <B>sat</B></P>
                        </Y>
                        HIII
                        <Divider style={{
                                marginBottom: 10
                        }} />


                        <T>{title || 'No title'}</T>
                        <D>{description || 'No description'}</D>
                </Pad>

        </Wrap>

}
const Pad = styled.div`
        padding:20px;
        `;
const Wrap = styled.div`
    display: flex;
    flex-direction:column;
        `;
const Y = styled.div`
        display: flex;
        justify-content:space-between;
        width:100%;
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

const Body = styled.div`
        font-size:14px;
        margin-left:10px;
        font-size: 15px;
        line-height: 20px;
        /* or 133% */

        display: flex;
        flex-direction:column;
        justify-content: space-around;

        /* Primary Text 1 */

        color: #292C33;
        overflow:hidden;
        `;

