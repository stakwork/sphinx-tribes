import { EuiLoadingSpinner } from '@elastic/eui';
import React, { useEffect, useState } from 'react'
import styled from 'styled-components';
import { useIsMobile } from '../../hooks';
import { useStores } from '../../store';
import PageLoadSpinner from './pageLoadSpinner';
// import { badges } from './constants';

export default function Badges(props) {

    const { main, ui } = useStores()
    const { badgeList } = ui || {}

    const [balancesTxns, setBalancesTxns]: any = useState({})
    const [loading, setLoading] = useState(true)

    const isMobile = useIsMobile()
    const { person } = props

    // console.log('badgeList', badgeList)
    // console.log('balancesTxns', balancesTxns)

    useEffect(() => {

        (async () => {
            setLoading(true)
            if (person?.owner_pubkey) {
                const b = await main.getBalances(person?.owner_pubkey)
                setBalancesTxns(b)
            }
            setLoading(false)
        })()

    }, [person?.owner_pubkey])

    const topLevelBadges = balancesTxns?.balances?.map((b, i) => {

        const badgeDetails = badgeList?.find(f => f.id === b.asset_id)
        // if early adopter badge
        let counter = ''
        if (b.asset_id === 8) {
            counter = balancesTxns?.txs?.find(f => f.asset_id === b.asset_id)?.metadata
        }

        // let status = 'Pending'
        // console.log('b', b)

        return <BWrap key={i + 'badges'} isMobile={isMobile}>
            <Img src={`${badgeDetails?.icon}`} />
            <div style={{ width: '100%', minWidth: 160, paddingRight: 30 }}>
                <T isMobile={isMobile}>{badgeDetails?.name} {b.balance > 1 && `(${b.balance})`}</T>
                {counter && <D><Counter>{counter}</Counter></D>}
            </div>

        </BWrap>
    })

    // always have at least one badge
    return (<Wrap >
        <PageLoadSpinner show={loading} />
        {topLevelBadges}
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
position:relative;
align-items: center;
padding: 10px; 
background: #fff;
margin-bottom: 10px;
border-radius: 4px;
box-shadow:0px 1px 2px rgb(0 0 0 / 15%);
min-width:200px;
width:${p => p.isMobile ? '100%' : 'auto'};

padding-right:20px;
margin-right:20px;

`;
const T = styled.div<BProps>`
font-size:15px;
margin-left:${p => p.isMobile ? '15px' : '10px'};
width:100%;
text-align:${p => p.isMobile ? '' : 'center'};
`;
const D = styled.div`
position:absolute;
top:0;
left:0;
width:100%;
font-size:12px;
display:flex;
justify-content:flex-end;
`;
const Counter = styled.div`
padding:3px 8px;
background:#D4AF37;
border-bottom-left-radius:4px;
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