import React from 'react'
import styled from "styled-components";
import Blog from './listItems/blog'


export default function WidgetList(props: any) {
    function renderByType(v, i) {
        function wrap(child) {
            return <IWrap
                key={i + 'listItem'}
                onClick={() => props.setSelected(v, i)}>
                <Eraser onClick={(e) => {
                    e.stopPropagation()
                    props.deleteItem(v, i)
                }}>X</Eraser>

                {child}
            </IWrap>
        }

        switch (props.schema.class) {
            case 'blog':
                return wrap(<Blog {...v} />)
            case 'twitter':
                return wrap(<Blog {...v} />)
            case 'donations':
                return wrap(<Blog {...v} />)
            case 'offer':
                return wrap(<Blog {...v} />)
            case 'wanted':
                return <></>
            default:
                return <></>
        }
    }

    return <Wrap >
        {props.values && props.values.map((v, i) => {
            return renderByType(v, i)
        })}

        {(!props.values || (props.values.length < 1)) && props.schema.itemLabel &&
            <IWrap>
                {props.schema.label} is empty
            </IWrap>
        }
    </Wrap>

}


export interface IconProps {
    source: string;
}

const Wrap = styled.div`
        color: #fff;
        width: 100%;
        margin-bottom:10px;
        display: flex;
        flex-direction: column-reverse;
        align-content: center;
        justify-content: space-evenly;
        `;

const IWrap = styled.div`
        position:relative;
        `;

const Eraser = styled.div`
        position:absolute;
        top:10px;
        right:10px;
        `;

const Icon = styled.img<IconProps>`
            background-image: ${p => `url(${p.source})`};
            width:100px;
            height:100px;
            background-position: center; /* Center the image */
            background-repeat: no-repeat; /* Do not repeat the image */
            background-size: contain; /* Resize the background image to cover the entire container */
            `;

