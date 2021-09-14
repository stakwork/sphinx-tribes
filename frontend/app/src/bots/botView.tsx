import React, { useState } from "react";
import styled from "styled-components";
import { useStores } from '../store'

import { Button, Divider, IconButton } from "../sphinxUI";
import { useIsMobile } from "../hooks";
import Bot from "./bot";
import BotBar from "./utils/botBar";

export default function BotView(props: any) {

    const {
        botUniqueName,
        selectBot,
        loading,
        goBack,
    } = props

    const { main } = useStores()

    const bot: any = (main.bots && main.bots.length && main.bots.find(f => f.unique_name === botUniqueName))

    const {
        name,
        unique_name,
        description,
        img
    } = bot || {}

    // FOR BOT VIEW
    const bots: any = (main.bots && main.bots.length && main.bots.filter(f => !f.hide))

    const isMobile = useIsMobile()

    if (loading) return <div>Loading...</div>

    function renderMobileView() {
        return <div style={{
            display: 'flex', flexDirection: 'column',
            width: '100%', overflow: 'auto', height: '100%'
        }}>
            <Panel style={{ paddingBottom: 0, paddingTop: 80 }}>
                <div style={{
                    position: 'absolute',
                    top: 20, left: 0,
                    display: 'flex',
                    justifyContent: 'space-between', width: '100%',
                    padding: '0 20px'
                }}>
                    <IconButton
                        onClick={goBack}
                        icon='arrow_back'
                    />
                    <div />
                </div>

                {/* profile photo */}
                <Head>
                    <Img src={img || '/static/sphinx.png'} />
                    <RowWrap>
                        <Name>{name}</Name>
                    </RowWrap>
                    <RowWrap>
                        <BotBar value={unique_name} />
                    </RowWrap>
                </Head>
            </Panel>

            <Sleeve style={{ padding: 20 }}>
                {description}
                <div style={{ height: 60 }} />
            </Sleeve>
        </div>
    }

    function renderDesktopView() {

        return <div style={{
            display: 'flex',
            width: '100%', height: '100%'
        }}>

            <PeopleList>
                <DBack >
                    <Button
                        color='clear'
                        leadingIcon='arrow_back'
                        text='Back'
                        onClick={goBack}
                    />
                </DBack>

                <div style={{ width: '100%', overflowY: 'auto' }} >
                    {bots.map(t => <Bot {...t} key={t.id}
                        selected={botUniqueName === t.unique_name}
                        hideActions={true}
                        small={true}
                        select={() => selectBot(t.unique_name)}
                    />)}
                </div>

            </PeopleList>

            <div style={{
                width: 364,
                minWidth: 364, overflowY: 'auto',
                position: 'relative',
                background: '#ffffff',
                color: '#000000',
                padding: 40,
                height: '100%',
                borderLeft: '1px solid #F2F3F5',
                borderRight: '1px solid #F2F3F5',
                boxShadow: '1px 0px 6px -2px rgba(0, 0, 0, 0.07)'
            }}>

                {/* profile photo */}
                <Head>
                    <div style={{ height: 35 }} />

                    <Img src={img || '/static/sphinx.png'} />

                    <RowWrap>
                        <Name>{name}</Name>
                    </RowWrap>
                    <RowWrap>
                        <BotBar value={unique_name} />
                    </RowWrap>

                    {/* only see buttons on other people's profile */}

                </Head>

                {/* Here's where the details go */}

            </div>

            <div style={{
                width: 'calc(100% - 628px)',
                minWidth: 250
            }}>
                <div style={{
                    padding: 62, height: 'calc(100% - 63px)',
                    overflowY: 'auto',
                    position: 'relative',
                }}>
                    <Sleeve style={{
                        display: 'flex',
                        alignItems: 'flex-start',
                        flexWrap: 'wrap',
                    }}>
                        {description}
                    </Sleeve>
                    <div style={{ height: 60 }} />
                </div>

            </div>

        </div >
    }


    return (
        <Content>
            {isMobile ? renderMobileView() : renderDesktopView()}
        </Content >

    );

}

const PeopleList = styled.div`
            display:flex;
            flex-direction:column;
            background:#ffffff;
            width: 265px;
            `;

const DBack = styled.div`
            height:64px;
            display:flex;
            align-items:center;
            background: #FFFFFF;
            box-shadow: 0px 0px 6px rgba(0, 0, 0, 0.07);
            z-index:2;
            `

const Panel = styled.div`
            position:relative;
            background:#ffffff;
            color:#000000;
            margin-bottom:10px;
            padding:20px;
            box-shadow:0px 0px 3px rgb(0 0 0 / 29%);
            `;
const Content = styled.div`
            display: flex;
            flex-direction:column;

            width:100%;
            height: 100%;
            align-items:center;
            color:#000000;
            background:#fff;
            `;

const Head = styled.div`
                display:flex;
                flex-direction:column;
                justify-content:center;
                align-items:center;
                width:100%;
                `;

const Name = styled.div`
                    font-style: normal;
                    font-weight: 500;
                    font-size: 30px;
                    line-height: 28px;
                    /* or 73% */

                    text-align: center;

                    /* Text 2 */

                    color: #3C3F41;
                    margin-bottom:20px;
                    `;


const Sleeve = styled.div`

                    `;

const RowWrap = styled.div`
                    display:flex;
                    justify-content:center;

                    width:100%`;


interface ImageProps {
    readonly src: string;
}
const Img = styled.div<ImageProps>`
                        background-image: url("${(p) => p.src}");
                        background-position: center;
                        background-size: cover;
                        margin-bottom:20px;
                        width:150px;
                        height:150px;
                        border-radius: 50%;
                        position: relative;
                        display:flex;
                        align-items:flex-end;
                        justify-content:flex-end;
                        `;
