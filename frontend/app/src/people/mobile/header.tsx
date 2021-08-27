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
import FadeLeft from '../../animated/fadeLeft';
// import EditInfo from '../edit/editInfo'
import SignIn from '../auth/signIn';

import PersonViewSlim from '../personViewSlim';
import { MeInfo } from '../../store/ui';
import api from '../../api';


export default function Header() {
    const { main, ui } = useStores()

    const people = useFuse(main.people, ["owner_alias"])
    const location = useLocation()
    const isMobile = useIsMobile()

    // function selectPerson(id: number, unique_name: string) {
    //   console.log('selectPerson', id, unique_name)
    //   setSelectedPerson(id)
    //   if (unique_name && window.history.pushState) {
    //     window.history.pushState({}, 'Sphinx Tribes', '/p/' + unique_name);
    //   }
    // }
    const c = colors['light']

    const tabs = [
        {
            label: 'Tribes',
            name: 'tribes',
            path: '/t/'
        },
        {
            label: 'People',
            name: 'people',
            path: '/p/'
        },
        {
            label: 'Bots',
            name: 'bots',
            path: '/b/'
        },
    ]

    const [showSignIn, setShowSignIn] = useState(false)
    const [showWelcome, setShowWelcome] = useState(false)
    const [showInitEditSelf, setShowInitEditSelf] = useState(false)
    const [showEditSelf, setShowEditSelf] = useState(false)

    async function testChallenge(chal: string) {
        try {
            const me: any = await api.get(`poll/${chal}`)
            if (me && me.pubkey) {
                ui.setMeInfo(me)
                setShowSignIn(false)
                setShowWelcome(true)
            }
        } catch (e) {
            console.log(e)
        }
    }

    useEffect(() => {
        try {
            var urlObject = new URL(window.location.href);
            var params = urlObject.searchParams;
            const chal = params.get('challenge')
            if (chal) {
                testChallenge(chal)
            }
        } catch (e) { }
    }, [])

    function renderHeader() {
        if (isMobile) {
            return <EuiHeader id="header" style={{ color: '#fff' }}>
                <div className="container">
                    <Row style={{ justifyContent: 'space-between' }}>
                        <EuiHeaderSection grow={false}>
                            <Img src="/static/people_logo.svg" />
                        </EuiHeaderSection>

                        <Corner>
                            {ui.meInfo ?
                                <Imgg
                                    style={{ height: 30, width: 30, marginRight: 10 }}
                                    src={(ui.meInfo.img && ui.meInfo.img + '?thumb=true') || '/static/sphinx.png'}
                                    onClick={() => setShowEditSelf(true)} />
                                :
                                <Button
                                    icon={'account_circle'}
                                    text={'Sign in'}
                                    color='primary'
                                    onClick={() => setShowSignIn(true)}
                                />
                            }
                        </Corner>

                    </Row>

                    <EuiHeaderSection id="header-right" side="right" style={{
                        background: '#000000',
                        boxShadow: 'inset 0px 1px 2px rgba(0, 0, 0, 0.15)',
                        borderRadius: 50, overflow: 'hidden'
                    }}>
                        <EuiFieldSearch id="search-input"
                            placeholder="Search for People"
                            value={ui.searchText}
                            onChange={e => ui.setSearchText(e.target.value)}
                            style={{ width: '100%', height: '100%' }}
                            aria-label="search"

                        />
                    </EuiHeaderSection>
                </div>
            </EuiHeader >
        }

        // desktop version
        return <EuiHeader style={{ color: '#fff', width: '100%', height: 64, padding: '0 20px' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
                <Row>
                    <EuiHeaderSection grow={false}>
                        <Img src="/static/people_logo.svg" />
                    </EuiHeaderSection>

                    <Tabs>
                        {tabs && tabs.map((t, i) => {
                            const label = t.label
                            const selected = t.name === 'people'

                            return <Tab key={i}
                                selected={selected}
                                onClick={() => {

                                }}>
                                {label}
                            </Tab>
                        })}

                    </Tabs>
                </Row>

                <Corner>
                    {ui.meInfo ?
                        <Button>
                            <div style={{ display: 'flex', alignItems: 'center' }} onClick={() => setShowEditSelf(true)}>
                                <Imgg
                                    style={{ height: 30, width: 30, marginRight: 10 }}
                                    src={(ui.meInfo.img && ui.meInfo.img + '?thumb=true') || '/static/sphinx.png'}
                                    onClick={() => setShowEditSelf(true)} />
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
                            onClick={() => setShowSignIn(true)}
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
                visible={showSignIn}
                close={() => setShowSignIn(false)}
                overlayClick={() => setShowSignIn(false)}
            >
                <SignIn
                    onSuccess={() => {
                        setShowSignIn(false)
                        setShowWelcome(true)
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
                            <div style={{ marginRight: 6 }}>Welcome</div>
                            <Name>{ui.meInfo?.owner_alias}</Name>
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
                                setShowEditSelf(true)
                            }}
                        />
                    </Column>
                </div>
            </Modal>

            {/* ONLY FOR FIRST TIME USER edit your info modal  */}
            {/* < Modal visible={showInitEditSelf}
                drift={40}
                fill
                close={() => setShowInitEditSelf(false)}
            >
                <div style={{
                    background: '#fff',
                    height: '100%',
                    width: '100%',
                    overflow: 'auto'
                }}>
                    
                </div>
            </Modal> */}

            {/* personSelect should just be kept in state */}
            < Modal visible={(ui.meInfo && showEditSelf) ? true : false}
                drift={40}
                fill
                close={() => setShowEditSelf(false)}
            >
                <PersonViewSlim goBack={() => setShowEditSelf(false)}
                    personId={ui.meInfo?.id}
                />
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


const Name = styled.div`
                font-style: normal;
                font-weight: 500;
                font-size: 26px;
                line-height: 19px;
                /* or 73% */

                text-align: center;

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

interface TagProps {
    selected: boolean;
}
const Tab = styled.div<TagProps>`
                        display:flex;
                        padding:10px;
                        margin-right:25px;
                        color:${p => p.selected ? '#fff' : '#8E969C'};
                        border-bottom: ${p => p.selected && '4px solid #618AFF'};                        
                        cursor:pointer;
                        font-weight: 500;
                        font-size: 13px;
                        line-height: 19px;
                        `;