import * as Yup from 'yup';

type FormFieldType =
  | 'header'
  | 'text'
  | 'textarea'
  | 'img'
  | 'imgcanvas'
  | 'gallery'
  | 'number'
  | 'hidden'
  | 'widgets'
  | 'widget'
  | 'switch'
  | 'select'
  | 'multiselect'
  | 'creatablemultiselect'
  | 'searchableselect'
  | 'loom'
  | 'space'
  | 'hide'
  | 'date';

type FormFieldClass = 'twitter' | 'blog' | 'offer' | 'wanted' | 'supportme';

export interface FormField {
  name: string;
  type: FormFieldType;
  class?: FormFieldClass;
  label: string;
  itemLabel?: string;
  single?: boolean;
  readOnly?: boolean;
  required?: boolean;
  validator?: any;
  style?: any;
  prepend?: string;
  widget?: boolean;
  page?: number;
  extras?: FormField[];
  fields?: FormField[];
  icon?: string;
  note?: string;
  extraHTML?: string;
  options?: any[];
  defaultSchema?: FormField[];
  defaultSchemaName?: string;
  dropdownOptions?: string;
  dynamicSchemas?: any[];
  testId?: string;
}

export function validator(config: FormField[]) {
  const shape: { [k: string]: any } = {};
  config.forEach((field: any) => {
    if (typeof field === 'object') {
      shape[field.name] = field.validator;
    }
  });
  return Yup.object().shape(shape);
}
