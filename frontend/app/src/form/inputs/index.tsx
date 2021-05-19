import React from 'react'
import TextInput from './text-input'
import ImageInput from './img-input'
import NumberInput from './number-input'

export default function Input(props:any){
    switch(props.type){
        case 'text':
            return <TextInput {...props}/>
        case 'img':
            return <ImageInput {...props}/>
        case 'number':
            return <NumberInput {...props}/>
        default:
            return <></>
    }
}