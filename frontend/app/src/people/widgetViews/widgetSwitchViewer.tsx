import React from 'react'
import OfferView from "../widgetViews/offerView";
import WantedView from "../widgetViews/wantedView";
import PostView from "../widgetViews/postView";
import styled from 'styled-components';
import { useIsMobile } from '../../hooks';
import { useStores } from '../../store';
import { useObserver } from 'mobx-react-lite';
import { useFuse, useScroll } from '../../hooks';
import { widgetConfigs } from '../utils/constants';
import moment from 'moment';
import { Spacer } from '../main/body';

// const getFuse = useFuse
// const getScroll = useScroll

export default function WidgetSwitchViewer(props) {

    const { main } = useStores()
    const isMobile = useIsMobile()

    return useObserver(() => {
        const { peoplePosts, peopleWanteds, peopleOffers } = main

        const listSource = {
            'post': peoplePosts,
            'wanted': peopleWanteds,
            'offer': peopleOffers
        }

        let { selectedWidget, onPanelClick } = props

        if (!selectedWidget) {
            return <div style={{ height: 200 }} />
        }

        const activeList = listSource[selectedWidget]

        const renderView = {
            post: (p, i) => <PostView showName key={i + p.owner_pubkey + 'view'} person={p} />,
            offer: (p, i) => <OfferView showName key={i + p.owner_pubkey + 'view'} person={p} />,
            wanted: (p, i) => <WantedView showName key={i + p.owner_pubkey + 'view'} person={p} />,
        }

        let allElements: any = []

        let searchKeys: any = widgetConfigs[selectedWidget]?.schema?.map(s => s.name) || []
        let foundDynamicSchema = widgetConfigs[selectedWidget]?.schema?.find(f => f.dynamicSchemas)
        // if dynamic schema, get all those fields
        if (foundDynamicSchema) {
            let dynamicFields: any = []
            foundDynamicSchema.dynamicSchemas?.forEach(ds => {
                ds.forEach(f => {
                    if (!dynamicFields.includes(f.name)) dynamicFields.push(f.name)
                })
            })
            searchKeys = dynamicFields
        }


        activeList && activeList.forEach((item, i) => {
            // if this person has entries for this widget
            if (item) {
                if (renderView[selectedWidget]) {
                    allElements.push(wrapIt(renderView[selectedWidget](item.person, i), item, i))
                }
            }
        })

        function wrapIt(child, item, i) {
            const { person, body } = item
            const panelStyles = isMobile ? {
                minHeight: 132
            } : {
                maxWidth: 291, minWidth: 291,
                marginRight: 20, marginBottom: 20, minHeight: 472
            }

            return <Panel key={person?.owner_pubkey + i}
                onClick={() => {
                    if (onPanelClick) onPanelClick(person, body)
                }}
                style={{
                    ...panelStyles,
                    cursor: 'pointer',
                    padding: 0, overflow: 'hidden'
                }}
            >
                {React.cloneElement(child, { ...body })}
            </Panel>
        }

        allElements = [...allElements, <Spacer key={'spacer'} />]

        return allElements
    })
}


const Panel = styled.div`
            position:relative;
            background:#ffffff;
            color:#000000;
            margin-bottom:10px;
            padding:20px;
            box-shadow:0px 0px 3px rgb(0 0 0 / 29%);
            `;

