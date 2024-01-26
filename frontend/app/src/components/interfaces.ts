export interface FadeLeftProps {
  drift?: number;
  isMounted: boolean;
  dismountCallback?: () => void;
  style?: React.CSSProperties;
  children?: JSX.Element;
  alwaysRender?: boolean;
  noFadeOnInit?: boolean;
  direction?: string;
  withOverlay?: boolean;
  overlayClick?: () => void;
  noFade?: boolean;
  speed?: number;
  close?: () => void;
}

export interface ButtonProps {
  iconStyle?: React.CSSProperties;
  id?: string;
  color?: string;
  leadingIcon?: string;
  endingIcon?: string;
  style?: React.CSSProperties;
  disabled?: boolean;
  height?: string | number;
  width?: string | number;
  icon?: string;
  onClick?: (any) => void;
  wideButton?: boolean;
  text?: string;
  iconSize?: number | string;
  imgSize?: number | string;
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
  style?: React.CSSProperties;
  setIsTop?: (any) => void;
  options: any;
  onInputChange?: (any) => void;
  loading?: boolean;
  selectStyle: React.CSSProperties;
  handleActive?: (any) => any;
  testId?: string;
  writeMode?: boolean;
  name?: string;
}

export interface AutoCompleteProps {
  peopleList: [{ [key: string]: string }];
  handleAssigneeDetails: (any) => void;
}

export interface IconButtonProps extends ButtonProps {
  iconStyle?: React.CSSProperties;
  color?: string;
  buttonType?: string;
  leadingImgStyle?: React.CSSProperties;
  endingImgStyle?: React.CSSProperties;
  leadingImg?: string;
  endingImg?: string;
  style?: React.CSSProperties | any;
  size?: string | number;
  onClick?: (any) => void;
  disabled?: boolean;
}

export interface ImageButtonProps extends ButtonProps {
  buttonAction?: (any) => void;
  size?: number;
  ButtonContainerStyle?: React.CSSProperties;
  leadingImageContainerStyle?: React.CSSProperties;
  leadingImageSrc?: string;
  endImageSrc?: string;
  buttonText: string;
  buttonTextStyle?: React.CSSProperties;
  endingImageContainerStyle?: React.CSSProperties;
}

export interface ModalProps {
  visible?: any;
  fill?: boolean;
  overlayClick?: () => void;
  dismountCallback?: () => void;
  children?: JSX.Element | JSX.Element[];
  close?: () => void;
  style?: React.CSSProperties;
  hideOverlay?: boolean;
  envStyle?: React.CSSProperties;
  nextArrow?: () => void;
  prevArrow?: () => void;
  nextArrowNew?: () => void;
  prevArrowNew?: () => void;
  bigClose?: () => void;
  bigCloseImage?: () => void;
  bigCloseImageStyle?: React.CSSProperties;
}

export interface SearchTextInputProps {
  small?: boolean;
  onChange: (any) => void;
  iconStyle?: React.CSSProperties;
  style: React.CSSProperties;
  name?: string;
  type?: string;
  placeholder?: string;
  value?: string;
}

export interface TextInputProps {
  label: string;
  value: string;
  onChange: (any) => void;
  handleBlur?: () => void;
  handleFocus?: () => void;
  readOnly?: boolean;
  prepend?: string;
  style?: React.CSSProperties;
}

export interface ElementProps {
  style?: React.CSSProperties;
  children: JSX.Element | JSX.Element[] | string;
  onClick?: (any) => void;
}
