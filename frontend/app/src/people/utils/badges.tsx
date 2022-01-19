
import React, { useEffect, useState } from 'react'
import styled from 'styled-components';
import { useIsMobile } from '../../hooks';
import { useStores } from '../../store';
import PageLoadSpinner from './pageLoadSpinner';
import { Modal, Button, Divider } from '../../sphinxUI';
// import { badges } from './constants';

export default function Badges(props) {

    const { main, ui } = useStores()
    const { badgeList, meInfo } = ui || {}

    const [balancesTxns, setBalancesTxns]: any = useState({})
    const [loading, setLoading] = useState(true)
    const [selectedBadge, setSelectedBadge]: any = useState(null)
    const [badgeToPush, setBadgeToPush]: any = useState(null)

    const isMobile = useIsMobile()
    const { person } = props

    // console.log('badgeList', badgeList)
    // console.log('balancesTxns', balancesTxns)

    const thisIsMe = meInfo?.owner_pubkey === person?.owner_pubkey

    useEffect(() => {

        (async () => {
            setLoading(true)
            setSelectedBadge(null)
            setBadgeToPush(null)
            if (person?.owner_pubkey) {
                const b = await main.getBalances(person?.owner_pubkey)
                setBalancesTxns(b)
            }
            setLoading(false)
        })()

    }, [person?.owner_pubkey])


    // metadata should be json to support badge details

    const topLevelBadges = balancesTxns?.balances?.map((b, i) => {

        const badgeDetails = badgeList?.find(f => f.id === b.asset_id)
        // if early adopter badge
        let counter = ''
        let metadata = balancesTxns?.txs?.find(f => f.asset_id === b.asset_id)?.metadata
        if (b.asset_id === 8) {
            counter = metadata
        }

        // let status = 'Pending'
        // console.log('b', b)


        const packedBadge = {
            ...badgeDetails,
            counter,
            metadata,
            deck: balancesTxns?.txs?.filter(f => f.asset_id === b.asset_id) || []
        }

        return <BWrap key={i + 'badges'} isMobile={isMobile} onClick={() => {
            setSelectedBadge(packedBadge)
        }}>
            <Img src={`${badgeDetails?.icon}`} isMobile={isMobile} />
            <div style={{ width: '100%', minWidth: 160 }}>
                <T isMobile={isMobile}>{badgeDetails?.name} {b.balance > 1 && `(${b.balance})`}</T>
                {badgeDetails?.description && <S isMobile={isMobile}>{badgeDetails?.description}</S>}
                {counter && <D><Counter>{counter} / {badgeDetails?.amount}</Counter></D>}
            </div>

            <Status>
                <StatusText>{'Off-chain'}</StatusText>
            </Status>

        </BWrap>
    })

    return (<Wrap >
        <PageLoadSpinner show={loading} />
        {selectedBadge ?
            <div style={{ width: '100%' }}>
                <Button
                    color='noColor'
                    leadingIcon='arrow_back'
                    text='Back to all badges'
                    onClick={() => setSelectedBadge(null)}
                    style={{ marginBottom: 20 }}
                />

                {selectedBadge.deck?.map((badge, i) => {
                    return <BWrap key={i + 'badges'} isMobile={isMobile} style={{ height: 'auto', minHeight: 'auto', cursor: 'default' }}>
                        <SmallImg src={`${selectedBadge?.icon}`} isMobile={isMobile} />
                        <div style={{
                            width: '100%', minWidth: 160, display: 'flex', flexDirection: 'column',
                            justifyContent: 'center'
                        }}>
                            <T isMobile={isMobile}>{selectedBadge?.name} {selectedBadge?.balance > 1 && `(${selectedBadge?.balance})`}</T>
                            {selectedBadge?.counter && <D><Counter>{selectedBadge?.counter} / {selectedBadge?.amount}</Counter></D>}


                            <div style={{ marginTop: 20, width: '100%' }}>
                                {thisIsMe ?
                                    <div style={{
                                        display: 'flex', flexDirection: 'column',
                                        justifyContent: 'center', width: '100%', textAlign: 'center'
                                    }}>
                                        <Divider />
                                        <Button
                                            style={{ margin: 0, marginTop: 2, padding: 0, minHeight: 40, border: 'none', }}
                                            color='link'
                                            text='Claim on Liquid'
                                            onClick={() => setBadgeToPush(badge)}
                                        />
                                    </div>
                                    :
                                    <Status >
                                        <StatusText>{'Off-chain'}</StatusText>
                                    </Status>}
                            </div>
                        </div>
                    </BWrap>
                })}
            </div>
            : topLevelBadges}

        <Modal
            visible={badgeToPush ? true : false}
            close={() => setBadgeToPush(null)}
        >
            <div style={{ padding: 20, display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
                <CodeText>
                    {JSON.stringify(badgeToPush)}
                </CodeText>
                <Button
                    color='primary'
                    text='Claim on Liquid'
                    disabled={true}
                />
            </div>
        </Modal>
    </Wrap >)

}


interface BProps {
    readonly isMobile?: boolean;
}

const Wrap = styled.div<BProps>`
            display: flex;
            flex-wrap:${p => p.isMobile ? '' : 'wrap'};;
            width:100%;
            overflow-x:hidden;
            `;


const BWrap = styled.div<BProps>`
                display: flex;
                cursor:pointer;
                flex-direction:${p => p.isMobile ? 'row' : 'column'};
                position:relative;
                width: ${p => p.isMobile ? '100%' : '192px'};
                min-width: ${p => p.isMobile ? '100%' : '192px'};
                height: ${p => p.isMobile ? '' : '272px'};
                min-height:${p => p.isMobile ? '' : '272px'};
                max-width: ${p => p.isMobile ? '100%' : '192px'};
                align-items: center;
                padding: ${p => p.isMobile ? '10px' : '20px 10px 10px'};
                background: #fff;
                margin-bottom: 10px;
                border-radius: 4px;
                box-shadow:0px 1px 2px rgb(0 0 0 / 15%);

                width:${p => p.isMobile ? '100%' : 'auto'};
                margin-right:${p => p.isMobile ? '0px' : '20px'};

                `;
const T = styled.div<BProps>`
                    font-size:15px;
                    width:100%;
                    text-align:${p => p.isMobile ? 'left' : 'center'};

                    font-family: Roboto;
                    font-style: normal;
                    font-weight: 600;
                    font-size: ${p => p.isMobile ? '20px' : '15px'};
                    line-height: 20px;
                    /* or 133% */

                    /* Primary Text 1 */

                    color: #292C33;
                    `;
const S = styled.div<BProps>`
                        font-size:15px;
                        margin-left:${p => p.isMobile ? '15px' : '10px'};
                        width:100%;
                        text-align:${p => p.isMobile ? '' : 'center'};

                        font-family: Roboto;
                        font-style: normal;
                        font-weight: 400;
                        font-size: 15px;
                        line-height: 15px;
                        /* or 133% */

                        text-align: center;

                        /* Primary Text 1 */

                        color: #5F6368;
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

const Status = styled.div`
                        margin: 15px 20px 25px;
                        // width:100%;
                        font-size:12px;
                        display:flex;
                        justify-content:center;
                        `;
const StatusText = styled.div`
                        display:flex;
                        justify-content:center;
                        align-items:center;
                        height:26px;
                        color:#5078F2;
                        background: #DCEDFE;
                        border-radius: 32px;
                        font-weight: bold;
                        font-size: 12px;
                        line-height: 13px;
                        padding 0 10px;
                        `;
const Counter = styled.div`
                        padding:3px 8px;
                        background:#D4AF3799;
                        border-bottom-left-radius:4px;
                        `;
const CodeText = styled.div`
    padding:20px;
    background:#A3C1FF55;
    word-break:break-word;
    margin:20px;
`;
interface ImageProps {
    readonly src?: string;
    readonly isMobile?: boolean;
}
const Img = styled.div<ImageProps>`
                            background-image: url("${(p) => p.src}");
                            background-position: center;
                            background-size: cover;
                            position: relative;
                            min-width:${p => p.isMobile ? '108px' : '132px'};
                            width:${p => p.isMobile ? '108px' : '132px'};
                            min-height:${p => p.isMobile ? '108px' : '132px'};
                            height:${p => p.isMobile ? '108px' : '132px'};
                            margin:${p => p.isMobile ? '24px' : '20px 30px'};
                            `;
const SmallImg = styled.div<ImageProps>`
                            background-image: url("${(p) => p.src}");
                            background-position: center;
                            background-size: cover;
                            position: relative;
                            min-width:90px;
                            width:90px;
                            min-height:90px;
                            height:90px;
                            margin:20px 30px;
                            `;
