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
            {
                checked: true,
                text: 'hi',
                number: 234
            },
            {
                checked: true,
                text: 'hi',
                number: 234
            },
            {
                checked: true,
                text: 'hi',
                number: 234
            },
            {
                checked: true,
                text: 'hi',
                number: 234
            },
            {
                checked: true,
                text: 'hi',
                number: 234
            },
            {
                checked: true,
                text: 'hi',
                number: 234
            },
        ]

        return <Drawer>
            <Spacer />
            <SearchTextInput
                name='search'
                type='search'
                placeholder='Search'
                value={ui.searchText}
                onChange={e => {
                    console.log('handleChange', e)
                    ui.setSearchText(e)
                }}

            />

            <Spacer />
            <Label>Tags</Label>

            <Tags>
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
  width:260px;
  padding-top:10px;
  overflow:auto;
  display:flex;
  flex-direction:column;
  color:#000;
  margin-right:50px;
//   border-right:1px solid #000;
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
    font-weight: bold;
    font-size: 16px;
    line-height: 31px;
    width:100%;
`
