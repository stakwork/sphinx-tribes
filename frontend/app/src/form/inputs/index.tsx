import React from 'react'
import TextInput from './text-input'
import TextAreaInput from './text-area-input'
import ImageInput from './img-input'
import GalleryInput from './gallery-input'
import NumberInput from './number-input'
import Widgets from './widgets/index'
import SwitchInput from './switch-input'

export default function Input(props: any) {
    switch (props.type) {
        case 'text':
            return <TextInput {...props} />
        case 'textarea':
            return <TextAreaInput {...props} />
        case 'img':
            return <ImageInput {...props} />
        case 'gallery':
            return <GalleryInput {...props} />
        case 'number':
            return <NumberInput {...props} />
        case 'switch':
            return <SwitchInput {...props} />
        case 'widgets':
            return <Widgets {...props} />
        case 'hidden':
            return <></>
        default:
            return <></>
    }
}