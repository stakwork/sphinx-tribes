
import React from 'react'
import styled from 'styled-components';
import moment from 'moment'
export default function NameTag(props) {
    const { owner_alias, img, created, style } = props

    let diffHours = 1
    if (created) {
        let a = moment()
        let b = moment.unix(created)
        diffHours = a.diff(b, 'hours')
    }


    return <div style={{
        display: 'flex', alignItems: 'center',
        width: 'fit-content', marginBottom: 10, ...style
    }}>
        <Img
            src={img || `/static/sphinx`}
        />
        <Name>
            {owner_alias}
        </Name>

        <div style={{ height: 4, width: 4, borderRadius: '50%', margin: '0 5px' }} />

        <Date>
            {!diffHours ? 'just now' : `${diffHours} h`}
        </Date>
    </div>

}

interface ImageProps {
    readonly src: string;
}
const Img = styled.div<ImageProps>`
            background-image: url("${(p) => p.src}");
            background-position: center;
            background-size: cover;
            height: 16px;
            width: 16px;
            border-radius: 50%;
            position: relative;
            `;

const Name = styled.div`
font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 12px;
line-height: 19px;
/* or 158% */
margin-left:5px;

display: flex;
align-items: center;

/* Secondary Text 4 */

color: #8E969C;

            `;

const Date = styled.div`
font-family: Roboto;
font-style: normal;
font-weight: normal;
font-size: 12px;
line-height: 19px;
/* or 158% */
margin-left:5px;

display: flex;
align-items: center;

/* Secondary Text 4 */

color: #8E969C;

            `;