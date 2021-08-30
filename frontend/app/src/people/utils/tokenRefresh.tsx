import React, { useEffect } from 'react'
import { useStores } from '../../store'

export default function TokenRefresh() {
    const { main, ui } = useStores()

    useEffect(() => {
        (async () => {
            if (ui.meInfo) {
                const res = await main.refreshJwt()
                if (res && res.jwt) {
                    console.log('token refresh!')
                    ui.setMeInfo({ ...ui.meInfo, jwt: res.jwt })
                }
                else {
                    console.log('kick!')
                }
            }
        })()
    }, [])

    return null
}
