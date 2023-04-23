export interface FadeLeftProps {
    drift: number;
    isMounted: boolean;
    dismountCallback: () => void;
    style: React.CSSProperties;
    children?: JSX.Element;
    alwaysRender?: boolean;
    noFadeOnInit?: boolean;
    direction: string;
    withOverlay: boolean;
    overlayClick: () => void;
    noFade?: boolean;
    speed?: number;
    close: () => void; 
}

export interface ButtonProps {
    iconStyle?: React.CSSProperties;
    id?: string;
    color: string;
    leadingIcon?: string;
    endingIcon?: string;
    style?: React.CSSProperties;
    disabled?: boolean;
    height?: number;
    width?: number;
    icon?: string;
    onClick?: (any) => void;
    wideButton?: boolean;
    text: string;
    iconSize?: string;
    imgSize?: string;
    imgStyle?: React.CSSProperties;
    leadingImgUrl?: string;
    ButtonTextStyle?: React.CSSProperties;
    children?: JSX.Element;
    img?: string;
    loading?: boolean;
    submitting?: boolean;
    hovercolor?: string;
    activecolor?: string;
    shadowcolor?: string;
    textStyle?: React.CSSProperties;
}

export interface SelProps {
    onChange: (any) => void;
    value: string;
    style: React.CSSProperties;
    setIsTop: (any) => void;
}

export interface AutoCompleteProps {
    peopleList: [{[key: string]: string}];
    handleAssigneeDetails: (any) => void;
}

export interface IconButtonProps extends ButtonProps {
    iconStyle: React.CSSProperties;
    color: string;
    buttonType?: string;
    leadingImgStyle?: React.CSSProperties;
    endingImgStyle?: React.CSSProperties;
    leadingImg?: string;
    endingImg?: string;
}

export interface ImageButtonProps extends ButtonProps {
    buttonAction: () => void;
    size: number;
    ButtonContainerStyle: React.CSSProperties;
    leadingImageContainerStyle: React.CSSProperties;
    leadingImageSrc: string;
    endImageSrc: string;
    buttonText: string;
    buttonTextStyle: React.CSSProperties;
    endingImageContainerStyle: React.CSSProperties;
}