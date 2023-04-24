export interface FormProps {
    buttonsOnBottom: React.CSSProperties;
    wrapStyle: React.CSSProperties,
    smallForm: boolean,
    readOnly: boolean,
    scrollDiv: any;
    initialValues: any;
    formRef: any;
    schema: any;
    paged?: boolean;
    onSubmit: () => void;
    isFirstTimeScreen: boolean;
    newDesign?: string;
    extraHTML?: JSX.Element | JSX.Element[];
    close?: () => void;
    loading?: boolean;
    delete: () => void;
}