import React, { useEffect, useState } from 'react'
import styled from 'styled-components';
import NoneSpace from '../people/utils/noneSpace';
import { Button } from '../sphinxUI';
import { useStores } from '../store';
import { useObserver } from 'mobx-react-lite'
import { EuiLoadingSpinner } from '@elastic/eui';
import { useFuse, useScroll } from '../hooks'
import { colors } from '../colors'
import FadeLeft from '../animated/fadeLeft';
import { useIsMobile } from '../hooks';
import Bot from './bot'

// avoid hook within callback warning by renaming hooks
const getFuse = useFuse
const getScroll = useScroll

export default function BotBody() {
    const { main, ui } = useStores()
    const [loading, setLoading] = useState(false)

    const c = colors['light']
    const isMobile = useIsMobile()

    function selectBot(id: number, unique_name: string) {
        console.log('selectBot', id, unique_name)
        ui.setSelectedBot(id)
        ui.setSelectingBot(id)
        if (unique_name && window.history.pushState) {
            window.history.pushState({}, 'Sphinx Tribes', '/b/' + unique_name);
        }
    }

    async function loadBots() {
        setLoading(true)
        let un = ''
        if (window.location.pathname.startsWith('/b/')) {
            un = window.location.pathname.substr(3)
        }
        const ps = await main.getBots(un)
        if (un) {
            const initial = ps[0]
            if (initial && initial.unique_name === un) ui.setSelectedBot(initial.id || 0)
        }
        setLoading(false)
    }

    useEffect(() => {
        loadBots()
    }, [])

    return useObserver(() => {
        const bs = getFuse(main.bots, ["owner_alias"])
        const { handleScroll, n, loadingMore } = getScroll()
        let bots = bs.slice(0, n)

        bots = (bots && bots.filter(f => !f.hide)) || []

        if (loading) {
            return <Body style={{ justifyContent: 'center', alignItems: 'center' }}>
                <EuiLoadingSpinner size="xl" />
            </Body>
        }

        if (!bots.length) {
            return (<>
                {/* <Button
                    text={'Make bot'}
                    color={'primary'}
                    onClick={() => main.makeBot()}
                /> */}
                <NoneSpace
                    img={'coming_soon.png'}
                    text={'COMING SOON'}
                    sub={'Stay tuned for something amazing!'}
                />
            </>)
        }

        if (isMobile) {
            return <Body>
                <div style={{ width: '100%' }} >
                    {bots.map(t => <Bot
                        {...t} key={t.id}
                        selected={ui.selectedBot === t.id}
                        small={isMobile}
                        select={selectBot}
                    />)}
                </div>
                <FadeLeft
                    withOverlay
                    drift={40}
                    overlayClick={() => ui.setSelectingBot(0)}
                    style={{ position: 'absolute', top: 0, right: 0, zIndex: 10000, width: '100%' }}
                    isMounted={ui.selectingPerson ? true : false}
                    dismountCallback={() => ui.setSelectedBot(0)}
                >
                    <BotView
                    // goBack={() => ui.setSelectingBot(0)}
                    // botId={ui.selectedBot}
                    // selectBot={selectBot}
                    // loading={loading}
                    />
                </FadeLeft>
            </Body >
        }

        // desktop mode
        return <Body style={{
            background: '#f0f1f3',
            height: 'calc(100% - 65px)'
        }}>
            <>
                <div style={{
                    width: '100%', display: 'flex', flexWrap: 'wrap', height: '100%',
                    justifyContent: 'flex-start', alignItems: 'flex-start', padding: 20
                }}>
                    {bots.map(t => <Bot
                        {...t} key={t.id}
                        small={false}
                        selected={ui.selectedBot === t.id}
                        select={selectBot}
                    />
                    )}
                </div>
                <div style={{ height: 100 }} />
            </>


            {/* selected view */}
            <FadeLeft
                withOverlay={isMobile}
                drift={40}
                overlayClick={() => ui.setSelectingBot(0)}
                style={{ position: 'absolute', top: isMobile ? 0 : 65, right: 0, zIndex: 10000, width: '100%' }}
                isMounted={ui.selectingPerson ? true : false}
                dismountCallback={() => ui.setSelectedBot(0)}
            >
                <BotView
                // goBack={() => ui.setSelectingPerson(0)}
                // personId={ui.selectedPerson}
                // loading={loading}
                // selectPerson={selectPerson}
                />
            </FadeLeft>
        </Body >
    }
    )
}

const BotView = styled.div`
            flex:1;
            height:calc(100% - 105px);
            padding-bottom:80px;
            width:100%;
            overflow:auto;
            display:flex;
            flex-direction:column;
            `
const Body = styled.div`
            flex:1;
            height:calc(100% - 105px);
            padding-bottom:80px;
            width:100%;
            overflow:auto;
            display:flex;
            flex-direction:column;
            `
const Label = styled.div`
            font-family: Roboto;
            font-style: normal;
            font-weight: bold;
            font-size: 26px;
            line-height: 40px;
            /* or 154% */
            
            display: flex;
            align-items: center;
            
            /* Text 2 */
            
            color: #3C3F41;`