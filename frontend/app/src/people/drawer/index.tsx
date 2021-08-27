import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import { useObserver } from 'mobx-react-lite'
import { useStores } from '../../store'
import { useFuse, useScroll } from '../../hooks'
import MaterialIcon from '@material/react-material-icon';
import { colors } from '../../colors'
import { SearchTextInput } from '../../sphinxUI/index'
import Tag from './tag'

// avoid hook within callback warning by renaming hooks
const getFuse = useFuse
const getScroll = useScroll

export default function DrawerComponent() {
    const { main, ui } = useStores()
    const [loading, setLoading] = useState(false)
    const c = colors['light']

    return useObserver(() => {
        const peeps = getFuse(main.people, ["owner_alias"])
        const { handleScroll, n, loadingMore } = getScroll()
        let people = peeps.slice(0, n)
        people = [...people, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}]

        const tags = [
            // {
            //     checked: true,
            //     text: 'hi',
            //     number: 234
            // },
            // {
            //     checked: true,
            //     text: 'hi',
            //     number: 234
            // },
            // {
            //     checked: true,
            //     text: 'hi',
            //     number: 234
            // },
            // {
            //     checked: true,
            //     text: 'hi',
            //     number: 234
            // },
            // {
            //     checked: true,
            //     text: 'hi',
            //     number: 234
            // },
            // {
            //     checked: true,
            //     text: 'hi',
            //     number: 234
            // },
        ]

        const width = 150

        return <Drawer>
            <Spacer />
            <SearchTextInput
                name='search'
                type='search'
                placeholder='Search'
                value={ui.searchText}
                style={{ width }}
                onChange={e => {
                    console.log('handleChange', e)
                    ui.setSearchText(e)
                }}

            />

            <Spacer />

            <Tags style={{ width }}>
                <Label>Tags</Label>
                <Spacer />
                {tags.map((t, i) => {
                    return <Tag {...t}
                        handleChange={(e) => console.log(e)}
                        key={i} />
                })
                }
            </Tags>

        </Drawer>
    }
    )
}
const Drawer = styled.div`
  height:100%;
  width:230px;
  padding-top:10px;
  overflow:auto;
  display:flex;
  flex-direction:column;
  align-items:center;
  color:#000;
`

const Tags = styled.div`
  display:flex;
  flex-direction:column;
  width:100%;
  
`

const Spacer = styled.div`
    height:20px;
`

const Label = styled.div`
    font-style: normal;
    font-weight: 500;
    font-size: 16px;
    line-height: 31px;
    width:100%;
`
