import React from 'react'
import styled from "styled-components";
import Input from "../../../form/inputs/";

export default function FocusedWidget(props: any) {
    const { name, values, errors, initialValues, setFieldTouched, setFieldValue, item } = props

    console.log('props', props)
    return <Wrap>
        {props.icon && <Icon source={props.icon} />}
        {item.fields.map((e, i) => {
            return <Input
                {...e}
                key={e.name}
                value={values[name] && values[name][item.name] && values[name][item.name][e.name]}
                error={errors[name] && errors[name][item.name] && errors[name][item.name][e.name]}
                initialValues={initialValues}
                handleChange={(c: any) => {
                    setFieldValue(`${name}.${item.name}.${e.name}`, c);
                }}
                handleBlur={() => setFieldTouched(`${name}.${item.name}.${e.name}`, false)}
                handleFocus={() => setFieldTouched(`${name}.${item.name}.${e.name}`, true)} />
        })}
    </Wrap>

}

const Wrap = styled.div`
    color: #fff;
    display: flex;
    flex-direction: column;
    align-content: center;
    justify-content: space-evenly;
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

