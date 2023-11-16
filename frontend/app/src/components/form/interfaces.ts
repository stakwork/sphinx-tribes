export interface FormProps {
  buttonsOnBottom?: React.CSSProperties | any;
  wrapStyle?: React.CSSProperties;
  smallForm?: boolean;
  readOnly?: boolean;
  scrollDiv?: any;
  initialValues?: any;
  formRef?: any;
  schema: any;
  paged?: boolean;
  onSubmit: (v: any) => Promise<void> | void;
  isFirstTimeScreen?: boolean;
  newDesign?: string | boolean;
  extraHTML?: JSX.Element | JSX.Element[] | any;
  close?: () => void;
  loading?: boolean;
  delete?: () => void;
  submitText?: string;
  onEditSuccess?: () => void;
  setLoading?: (value: boolean) => void;
}
