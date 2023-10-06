import { Button } from 'components/common';
import { userHasRole } from 'helpers';
import React, { useState, useEffect, useCallback } from 'react';
import { useStores } from 'store';

const ManageButton = (props: { user_pubkey: string, org: any, action: () => void }) => {
    const [userRoles, setUserRoles] = useState<any[]>([]);
    const { main, ui } = useStores();

    const { user_pubkey, org, action } = props;

    const getUserRoles = useCallback(async () => {
        const userRoles = await main.getUserRoles(org.uuid, user_pubkey);
        setUserRoles(userRoles);
    }, [org.uuid])

    useEffect(() => {
        getUserRoles();
    }, [getUserRoles])


    const isOrganizationAdmin = org?.owner_pubkey === ui.meInfo?.owner_pubkey;

    return (
        <>
            {
                (isOrganizationAdmin ||
                    userHasRole(main.bountyRoles, userRoles, 'ADD USER') ||
                    userHasRole(main.bountyRoles, userRoles, 'VIEW REPORT'))
                && (<Button
                    text="Manage"
                    color="white"
                    style={{
                        width: 112,
                        height: 40,
                        color: '#000000',
                        borderRadius: 10
                    }}
                    onClick={action}
                />)
            }
        </>

    )
}

export default ManageButton;