import React, { useState, useEffect } from 'react'
import styled from "styled-components";
import Input from "../../../form/inputs/";
import { EuiButton } from '@elastic/eui'
import FadeLeft from '../../../animated/fadeLeft';
import WidgetList from './widgetList'

export default function FocusedWidget(props: any) {
    const { name, values, errors, initialValues, setFieldTouched,
        setFieldValue, item, setShowFocused, setDisableFormButtons } = props
    const { single } = item
    const [selectedIndex, setSelectedIndex] = useState(-1)

    let newState = values[name] && values[name][item.name]
    let newErrors = errors[name] && errors[name][item.name]

    if (selectedIndex > -1) {
        newState = newState && newState[selectedIndex]
        newErrors = newErrors && newErrors[selectedIndex]
    }

    const [prevState, setPrevState] = useState(newState)

    useEffect(() => {
        setDisableFormButtons(true)

        return function cleanup() {
            setDisableFormButtons(false)
        }
    }, [])

    function cancel() {
        // new widget cancelled, revert form state
        setFieldValue(`${name}.${item.name}`, prevState);
        if (single) setShowFocused(false)
        else setSelectedIndex(-1)
    }

    function done() {
        // new widget added, keep form state and reset selected index
        // if errors, prevent!
        if (single) setShowFocused(false)
        else setSelectedIndex(-1)


    }

    function startCreate() {
        let formState = (values[name] && values[name][item.name]) || []
        setPrevState(formState)

        const obj = {}
        item.fields.forEach((o) => {
            let val: any = ""
            if (o.type === 'number') val = 0
            obj[o.name] = val
        })

        formState.unshift(obj)

        setFieldValue(`${name}.${item.name}`, formState);
        setSelectedIndex(0)
    }

    function startEdit(obj, i) {
        console.log('startEdit', obj)
        let formState = (values[name] && values[name][item.name]) || []
        let index = i
        if (obj.id) {
            index = formState && formState.findIndex(f => f.id === obj.id)
        }

        setPrevState(formState)
        setSelectedIndex(index)
    }

    function deleteItem(obj, i) {
        console.log('deleteItem', obj)
        let formState = (values[name] && values[name][item.name]) || []
        let index = i
        if (obj.id) {
            index = formState && formState.findIndex(f => f.id === obj.id)
        }

        formState.splice(index, 1)
        setFieldValue(`${name}.${item.name}`, formState);
    }


    function getFieldToUpdate(e) {
        let valueToUpdate = `${name}.${item.name}.${e.name}`
        if (!single) {
            valueToUpdate = `${name}.${item.name}[${selectedIndex}].${e.name}`
        }
        return valueToUpdate
    }

    const showingList = single ? false : (selectedIndex > -1) ? false : true

    return <Wrap>



        {props.icon && <Icon source={props.icon} />}

        {/* single widgets will only show these fields */}
        {/* <FadeLeft isMounted={showInputs} dismountCallback={() =>}> */}
        {!showingList && <>

            {item.fields.map((e, i) => {
                return <Input
                    {...e}
                    key={e.name}
                    value={newState && newState[e.name]}
                    error={newErrors && newErrors[e.name]}
                    initialValues={initialValues}
                    handleChange={(c: any) => {
                        let f = getFieldToUpdate(e)
                        setFieldValue(f, c)
                    }}
                    handleBlur={() => {
                        let f = getFieldToUpdate(e)
                        setFieldTouched(f, false)
                    }}
                    handleFocus={() => {
                        let f = getFieldToUpdate(e)
                        setFieldTouched(f, true)
                    }} />
            })}
            <Nav style={{ marginTop: 20 }}>
                <EuiButton
                    onClick={() => cancel()}
                    style={{ fontSize: 12, fontWeight: 600 }}
                >Cancel</EuiButton>
                <EuiButton
                    onClick={() => done()}
                    style={{ fontSize: 12, fontWeight: 600 }}
                >Save{item.itemLabel && ` ${item.itemLabel}`}</EuiButton>
            </Nav>
        </>}
        {/* </FadeLeft> */}

        {/* only multi widgets will only show the list */}
        {/* <FadeLeft isMounted={showingList} dismountCallback={() =>}> */}

        {showingList && <>
            <Nav >
                <EuiButton
                    onClick={() => setShowFocused(false)}
                    style={{ fontSize: 12, fontWeight: 600 }}
                >Back to Widgets
                </EuiButton>

                <EuiButton
                    onClick={() => startCreate()}
                    style={{ fontSize: 12, fontWeight: 600 }}
                >New{item.itemLabel && ` ${item.itemLabel}`}</EuiButton>
            </Nav>

            <WidgetList schema={item}
                values={values[name] && values[name][item.name]}
                setSelected={startEdit}
                deleteItem={deleteItem} />

        </>}

        {/* </FadeLeft> */}

    </Wrap >

}

const Wrap = styled.div`
    color: #fff;
    display: flex;
    flex-direction: column;
    align-content: center;
    justify-content: space-evenly;
`;

const Nav = styled.div`
    color: #fff;
    display: flex;
    align-content: center;
    justify-content: space-evenly;
    margin-bottom:20px;
    margin-top:2px;
    height:42px;
    min-height:42px;
`;

export interface IconProps {
    source: string;
}

const Icon = styled.img<IconProps>`
    background-image: ${p => `url(${p.source})`};
    width:100px;
    height:100px;
    background-position: center; /* Center the image */
    background-repeat: no-repeat; /* Do not repeat the image */
    background-size: contain; /* Resize the background image to cover the entire container */
`;

