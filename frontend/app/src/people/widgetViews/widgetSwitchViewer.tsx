import React from 'react'
import BlogView from "../widgetViews/blogView";
import OfferView from "../widgetViews/offerView";
import TwitterView from "../widgetViews/twitterView";
import SupportMeView from "../widgetViews/supportMeView";
import WantedView from "../widgetViews/wantedView";
import PostView from "../widgetViews/postView";
import styled from 'styled-components';
import { useIsMobile } from '../../hooks';
import { useStores } from '../../store';
import { useObserver } from 'mobx-react-lite';
import { useFuse, useScroll } from '../../hooks';
import { widgetConfigs } from '../utils/constants';
import moment from 'moment';
const getFuse = useFuse
const getScroll = useScroll

export default function WidgetSwitchViewer(props) {

    const { main, ui } = useStores()
    const isMobile = useIsMobile()

    return useObserver(() => {

        let { selectedWidget, onPanelClick } = props
        const peeps = [...main.people]
        const { handleScroll, n, loadingMore } = getScroll()
        let people = peeps.slice(0, n)
        people = (people && people.filter(f => !f.hide)) || []

        if (props.people) {
            // props.people overrides ui
            people = props.people
        }

        if (!selectedWidget) {
            return <div style={{ height: 200 }} />
        }

        const renderView = {
            post: (p, i) => <PostView showName key={i + p.owner_pubkey + 'view'} person={p} />,
            offer: (p, i) => <OfferView showName key={i + p.owner_pubkey + 'view'} person={p} />,
            wanted: (p, i) => <WantedView showName key={i + p.owner_pubkey + 'view'} person={p} />,
        }

        let allElements: any = []

        const searchKeys = widgetConfigs[selectedWidget] && widgetConfigs[selectedWidget].schema?.map(s => s.name) || []

        let peopleClone = [...people]
        peopleClone && peopleClone.sort((a: any, b: any) => {
            return moment(a.updated).valueOf() - moment(b.updated).valueOf()
        }).reverse().forEach((p, i) => {
            // if this person has entries for this widget
            const thisWidget = p.extras && p.extras[selectedWidget]
            if (thisWidget && thisWidget.length) {
                const theseExtras = getFuse(thisWidget, searchKeys)
                if (renderView[selectedWidget]) {
                    allElements = [...allElements, ...wrapIt(renderView[selectedWidget](p, i), theseExtras, p)]
                }
            }
        })

        function wrapIt(child, fullSelectedWidget, person) {
            const elementArray: any = []

            const panelStyles = isMobile ? {
                minHeight: 132
            } : {
                maxWidth: 291, minWidth: 291,
                marginRight: 20, marginBottom: 20, minHeight: 472
            }

            fullSelectedWidget && fullSelectedWidget.forEach((s, i) => {
                elementArray.push(<Panel key={i + person.owner_pubkey}
                    onClick={() => {
                        if (onPanelClick) onPanelClick(person, s, i)
                    }}
                    style={{
                        ...panelStyles,
                        cursor: 'pointer',
                        padding: 0, overflow: 'hidden'
                    }}
                >
                    {React.cloneElement(child, { ...s })}
                </Panel>)
            })

            return elementArray
        }

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

