import React from 'react'
import styled from 'styled-components';
import { useIsMobile } from '../../hooks';
import { badges } from './constants';

export default function Badges(props) {

    // const { ui } = useStores()
    // const { searchText } = ui || {}

    const isMobile = useIsMobile()
    const { person } = props

    const b = Object.keys(badges).map((b, i) => {
        const thisbadge = badges[b]
        return <BWrap key={i + 'badges'} isMobile={isMobile}>
            <Img src={`/static/${thisbadge.src}`} />
            <div style={{ width: '100%', minWidth: 160, paddingRight: 30 }}>
                <T isMobile={isMobile}>{thisbadge.title || 'Title'}</T>
                {/* <D>{thisbadge.desc || 'Desc'}</D> */}
            </div>
        </BWrap>
    })

    // always have at least one badge
    return (<Wrap >
        {/* <div style={{
            display: 'flex', justifyContent: 'center', alignItems: 'center',
            position: 'absolute', top: 20, left: 0, width: '100%', zIndex: 10,
        }}>
            <div style={{
                width: 300, background: '#fff', padding: 20, borderRadius: 4,
                textAlign: 'center',
                boxShadow: '0px 2px 6px rgb(0 0 0 / 15%)',
                fontWeight: 600,
                fontSize: 22,
                letterSpacing: 0,
                color: 'rgb(60,63,65)',

            }}>
                Coming soon!
            </div>
        </div> */}
        {b}
    </Wrap >)

}

const Wrap = styled.div`
display: flex;
position:relative;
flex-wrap:wrap;
width:100%;
overflow-x:hidden;
// justify-content:space-around;
`;

interface BProps {
    readonly isMobile?: boolean;
}
const BWrap = styled.div<BProps>`
display: flex;
align-items: center;
padding: 10px; 
background: #fff;
margin-bottom: 10px;
border-radius: 4px;
box-shadow:0px 1px 2px rgb(0 0 0 / 15%);
min-width:200px;
width:${p => p.isMobile ? '100%' : 'auto'};

padding-right:10px;
margin-right:10px;
overflow-x:hidden;
`;
const T = styled.div<BProps>`
font-size:15px;
margin-left:${p => p.isMobile ? '15px' : '10px'};
width:100%;
text-align:${p => p.isMobile ? '' : 'center'};
`;
const D = styled.div`
font-size:18px;
margin-left:3px;
`;
interface ImageProps {
    readonly src?: string;
}
const Img = styled.div<ImageProps>`
                                        background-image: url("${(p) => p.src}");
                                        background-position: center;
                                        background-size: cover;
                                        position: relative;
                                        min-width:52px;
                                        width:52px;
                                        min-height:52px;
                                        height:52px;
                                        margin-right: 5px;
                                        `;