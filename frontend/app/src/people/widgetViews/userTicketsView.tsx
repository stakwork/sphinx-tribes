import React, { useEffect, useState } from "react";
import WantedView from '../widgetViews/wantedView';
import { useHistory, useParams } from "react-router-dom";
import { useStores } from "store";
import { bountyHeaderFilter, bountyHeaderLanguageFilter } from '../utils/filterValidation';
import NoResults from "people/utils/noResults";
import { Panel } from "people/personSlim/style";
import { useIsMobile } from "hooks";
import { colors } from '../../config/colors';
import DeleteTicketModal from "./deleteModal";
import { Spacer } from "people/main/body";

const UserTickets = () => {
    const color = colors['light'];
    const { personPubkey } = useParams<{ personPubkey: string }>();
    const { main, ui } = useStores();
    const isMobile = useIsMobile();
    const history = useHistory();
    const [userTickets, setUserTickets] = useState<any>([]);
    const [checkboxIdToSelectedMap, setCheckboxIdToSelectedMap] = useState<any>({});
    const [checkboxIdToSelectedMapLanguage, setCheckboxIdToSelectedMapLanguage] = useState({});
    const [currentItems] = useState<number>(10);
    const [deletePayload, setDeletePayload] = useState<object>({});
    const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
    const closeModal = () => setShowDeleteModal(false);
    const showModal = () => setShowDeleteModal(true);

    const data = {
        checkboxIdToSelectedMap,
    };

    const panelStyles = isMobile
        ? {
            minHeight: 132
        }
        : {
            minWidth: '1100px',
            maxWidth: '1100px',
            marginBottom: 16,
            borderRadius: '10px',
            display: 'flex',
            justifyContent: 'center'
        };

    const activeList = userTickets.filter(({ body }) => {
        const value = { ...body };
        return (
            bountyHeaderFilter(data?.checkboxIdToSelectedMap, value?.paid, !!value?.assignee) &&
            bountyHeaderLanguageFilter(value?.codingLanguage, checkboxIdToSelectedMapLanguage)
        );
    });

    async function getUserTickets() {
        const tickets = await main.getPersonWanteds({}, personPubkey);
        console.log("User Tickets ===", tickets);
        setUserTickets(tickets);
    }

    function onPanelClick(person, item) {
        history.replace({
            pathname: history?.location?.pathname,
            search: `?owner_id=${person.owner_pubkey}&created=${item.created}`,
            state: {
                owner_id: person.owner_pubkey,
                created: item.created
            }
        });
    }

    const confirmDelete = async () => {
        try {
            if (deletePayload) {
                await deleteTicket(deletePayload);
            }
        } catch (error) {
            console.log(error);
        }
        closeModal();
    };

    const deleteTicket = async (payload: any) => {
        const info = ui.meInfo as any;
        const URL = info.url.startsWith('http') ? info.url : `https://${info.url}`;
        try {
            await fetch(`${URL}/delete_ticket`, {
                method: 'POST',
                body: JSON.stringify(payload),
                headers: {
                    'x-jwt': info.jwt,
                    'Content-Type': 'application/json'
                }
            });
        } catch (error) {
            console.log(error);
        }
    };

    useEffect(() => {
        getUserTickets();
    }, [])

    const listItems =
        activeList && activeList.length ? (
            activeList.slice(0, currentItems).map((item, i) => {
                const { person, body } = item;
                const conditionalStyles = body?.paid
                    ? {
                        border: isMobile ? `2px 0 0 0 solid ${color.grayish.G600}` : '',
                        boxShadow: 'none'
                    }
                    : {};

                // if this person has entries for this widget
                return (
                    <Panel
                        isMobile={isMobile}
                        key={person?.owner_pubkey + i + body?.created}
                        style={{
                            ...panelStyles,
                            ...conditionalStyles,
                            cursor: 'pointer',
                            padding: 0,
                            overflow: 'hidden',
                            background: 'transparent',
                            minHeight: !isMobile ? '160px' : '',
                            boxShadow: 'none'
                        }}
                    >

                        <WantedView
                            showName
                            onPanelClick={() => {
                                if (onPanelClick) onPanelClick(person, body);
                            }}
                            person={person}
                            showModal={showModal}
                            setDeletePayload={setDeletePayload}
                            fromBountyPage={true}
                            {...body}
                        />
                    </Panel>
                );
            })
        ) : (
            <NoResults />
        );

    return (
        <>
            {listItems}
            <Spacer key={'spacer2'} />
            {showDeleteModal && (
                <DeleteTicketModal closeModal={closeModal} confirmDelete={confirmDelete} />
            )}
        </>
    )
}

export default UserTickets;

