import React, { useState, useEffect } from 'react'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../../store'
import styled from 'styled-components'
import {
    EuiHeader,
    EuiHeaderSection,
    EuiFieldSearch,
} from '@elastic/eui';
import { useFuse, useIsMobile } from '../../hooks'
import { colors } from '../../colors'
import { useHistory, useLocation } from 'react-router-dom'
import { Modal, Button, Divider } from '../../sphinxUI';

// import EditInfo from '../edit/editInfo'
import SignIn from '../auth/signIn';

import api from '../../api';


export default function Header() {
    const { main, ui } = useStores()
    const location = useLocation()
    const history = useHistory()
    const isMobile = useIsMobile()

    const c = colors['light']

    const tabs = [
        {
            label: 'Tribes',
            name: 'tribes',
            path: '/t'
        },
        {
            label: 'People',
            name: 'people',
            path: '/p'
        },
        {
            label: 'Bots',
            name: 'bots',
            path: '/b'
        },
    ]

    const [showWelcome, setShowWelcome] = useState(false)

    async function testChallenge(chal: string) {
        try {
            const me: any = await api.get(`poll/${chal}`)
            if (me && me.pubkey) {
                ui.setMeInfo(me)
                ui.setShowSignIn(false)
                setShowWelcome(true)
            }
        } catch (e) {
            console.log(e)
        }
    }

    function urlRedirect() {
        // if route not supported, redirect
        let pass = false
        let path = location.pathname
        tabs.forEach((t => {
            if (path.includes(t.path)) pass = true
        }))
        if (!pass) history.push('/p')
    }

    useEffect(() => {
        let path = location.pathname
        if (!path.includes('/p') && (ui.selectedPerson || ui.selectingPerson)) {
            ui.setSelectedPerson(0)
            ui.setSelectingPerson(0)
        }
    }, [location.pathname])

    useEffect(() => {
        (async () => {
            try {
                var urlObject = new URL(window.location.href);
                var params = urlObject.searchParams;
                const chal = params.get('challenge')
                if (chal) {
                    await testChallenge(chal)
                }
                urlRedirect()
            } catch (e) { }
        })()

    }, [])

    function goToEditSelf() {
        // if page is not /p, go to /p (people)
        let path = location.pathname
        if (!path.includes('/p')) history.push('/p')

        if (ui.meInfo?.id) {
            ui.setSelectedPerson(ui.meInfo.id)
            ui.setSelectingPerson(ui.meInfo.id)
        }
    }

    const headerBackground = '#1A242E'

    function renderHeader() {
        if (isMobile) {
            return <EuiHeader id="header" style={{
                color: '#fff', background: headerBackground, paddingBottom: 0,
            }}>
                < div className="container" >
                    <Row style={{ justifyContent: 'space-between' }}>
                        <EuiHeaderSection grow={false}>
                            <Img src="/static/people_logo.svg" />
                        </EuiHeaderSection>

                        <Corner>
                            {ui.meInfo ?
                                <Imgg
                                    style={{ height: 30, width: 30, marginRight: 10, border: '1px solid #ffffff55' }}
                                    src={ui.meInfo?.img || '/static/sphinx.png'}
                                    onClick={() => {
                                        goToEditSelf()
                                    }} />
                                :
                                <Button
                                    icon={'account_circle'}
                                    text={'Sign in'}
                                    color='primary'
                                    onClick={() => ui.setShowSignIn(true)}
                                />
                            }
                        </Corner>
                    </Row>

                    {/* <EuiHeaderSection id="hide-icons" style={{ margin: '10px 10px', borderRadius: 50, overflow: 'hidden' }} >
                        <EuiFieldSearch id="search-input"
                            placeholder="Search"
                            value={ui.searchText}
                            onChange={e => ui.setSearchText(e.target.value)}
                            style={{ width: '100%', height: '100%' }}
                            aria-label="search"

                        />
                    </EuiHeaderSection> */}

                    <MTabs>
                        {tabs && tabs.map((t, i) => {
                            const label = t.label
                            const selected = location.pathname.includes(t.path)

                            return <MTab key={i}
                                selected={selected}
                                onClick={() => {
                                    history.push(t.path)
                                }}>
                                {label}
                            </MTab>
                        })}

                    </MTabs>
                </div>
            </ EuiHeader >
        }

        // desktop version
        return <EuiHeader style={{ color: '#fff', width: '100%', height: 64, padding: '0 20px', background: headerBackground }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
                <Row>
                    <EuiHeaderSection grow={false}>
                        <Img src="/static/people_logo.svg" />
                    </EuiHeaderSection>

                    <Tabs>
                        {tabs && tabs.map((t, i) => {
                            const label = t.label
                            const selected = location.pathname.includes(t.path)

                            return <Tab key={i}
                                selected={selected}
                                onClick={() => {
                                    history.push(t.path)
                                }}>
                                {label}
                            </Tab>
                        })}

                    </Tabs>

                    {/* <EuiHeaderSection id="hide-icons" style={{ margin: '10px 10px', borderRadius: 50, overflow: 'hidden', width: 295 }} >
                        <EuiFieldSearch id="search-input"
                            placeholder="Search for People"
                            value={ui.searchText}
                            onChange={e => ui.setSearchText(e.target.value)}
                            style={{ width: 295, height: '100%' }}
                            aria-label="search"

                        />
                    </EuiHeaderSection> */}

                    {/* <Button
                        text='Tags'
                    /> */}
                </Row>

                <Corner>
                    {ui.meInfo ?
                        <Button onClick={() => {
                            goToEditSelf()
                        }}>
                            <div style={{ display: 'flex', alignItems: 'center' }}>
                                <Imgg
                                    style={{ height: 30, width: 30, marginRight: 10 }}
                                    src={ui.meInfo?.img || '/static/sphinx.png'} />
                                <div style={{ color: '#fff' }}>
                                    {ui.meInfo?.owner_alias}
                                </div>
                            </div>
                        </Button>

                        :
                        <Button
                            icon={'account_circle'}
                            text={'Sign in'}
                            color='primary'
                            onClick={() => ui.setShowSignIn(true)}
                        />
                    }
                </Corner>
            </div>
        </EuiHeader >
    }


    return useObserver(() => {
        return <>

            {renderHeader()}

            {/* you wanna login modal  */}
            <Modal
                visible={ui.showSignIn}
                close={() => ui.setShowSignIn(false)}
                overlayClick={() => ui.setShowSignIn(false)}
            >
                <SignIn
                    onSuccess={() => {
                        ui.setShowSignIn(false)
                        setShowWelcome(true)

                        // if page is not /p, go to /p (people)
                        let path = location.pathname
                        if (!path.includes('/p')) history.push('/p')
                    }} />
            </Modal >

            {/* you logged in modal  */}
            < Modal
                visible={(ui.meInfo && showWelcome) ? true : false}>
                <div>
                    <Column>
                        <Imgg
                            style={{ height: 128, width: 128, marginBottom: 40 }}
                            src={ui.meInfo?.img || '/static/sphinx.png'} />

                        <T>
                            <div style={{ lineHeight: '26px' }}>Welcome <Name>{ui.meInfo?.owner_alias}</Name></div>

                        </T>

                        <Welcome>
                            Your profile is now public.
                            Connect with other people, join tribes and listen your favorite podcast!
                        </Welcome>

                        <Button
                            text={'Continue'}
                            height={48}
                            width={'100%'}
                            color={'primary'}
                            onClick={() => {
                                // switch from welcome modal to edit modal
                                setShowWelcome(false)
                                goToEditSelf()
                            }}
                        />
                    </Column>
                </div>
            </Modal>
        </>
    })
}

const Row = styled.div`
                display:flex;
                align-items:center;
                width:100%;
                `
const Corner = styled.div`
                display:flex;
                align-items:center;
                `
const T = styled.div`
                display:flex;
                font-size: 26px;
                line-height: 19px;
                `

interface ImageProps {
    readonly src: string;
}
const Img = styled.div<ImageProps>`
                    background-image: url("${(p) => p.src}");
                    background-position: center;
                    background-size: cover;
                    height:37px;
                    width:232px;

                    position: relative;
                    `;


const Name = styled.span`
                    font-style: normal;
                    font-weight: 500;
                    font-size: 26px;
                    line-height: 19px;
                    /* or 73% */

                    /* Text 2 */

                    color: #292C33;
                    `;
const Welcome = styled.div`
                    font-size: 15px;
                    line-height: 24px;
                    margin:20px 0 50px;
                    text-align: center;

                    /* Text 2 */

                    color: #3C3F41;
                    `


const Column = styled.div`
                    width:100%;
                    display:flex;
                    flex-direction:column;
                    justify-content:center;
                    align-items:center;
                    padding: 25px;

                    `
const Imgg = styled.div<ImageProps>`
                        background-image: url("${(p) => p.src}");
                        background-position: center;
                        background-size: cover;
                        width:90px;
                        height:90px;
                        border-radius: 50%;
                        `;

const Tabs = styled.div`
                        display:flex;
                        margin-left:20px;
                        `;

const MTabs = styled.div`
                        display:flex;
                        margin:0 20px;
                        justify-content:space-around;
                        `;
interface TagProps {
    selected: boolean;
}
const Tab = styled.div<TagProps>`
                            display:flex;
                            padding:10px 25px;
                            margin-right:10px;
                            color:${p => p.selected ? '#fff' : '#6B7A8D'};
                            cursor:pointer;
                            font-weight: 600;
                            font-size: 15px;
                            line-height: 19px;
                            background:${p => p.selected ? 'rgba(255,255,255,0.07)' : '#3C3F4100'};
                            border-radius:25px;
                            `;

const MTab = styled.div<TagProps>`
                            display:flex;
                            margin:25px 5px 0;
                            color:${p => p.selected ? '#fff' : '#ffffff99'};
                            cursor:pointer;
                            height:30px;
                            min-width:65px;
                            font-weight: 500;
                            font-size: 15px;
                            line-height: 19px;
                            justify-content:center;
                            border-bottom:${p => p.selected ? '3px solid #618AFF' : 'none'};
                            `;