import React from "react";
import { Formik } from "formik";
import * as Yup from "yup";
import styled from "styled-components";
import Input from "./inputs";
import {EuiButton} from '@elastic/eui'

export default function Form(props:any) {
  return (
    <Formik
      initialValues={props.initialValues || {}}
      onSubmit={props.onSubmit}
      validationSchema={validator(props.schema)}
    >
      {({setFieldTouched, handleSubmit, values, setFieldValue, errors, dirty, isValid, initialValues}) => {
        return (
          <Wrap>
            {props.schema && props.schema.map((item:FormField) => <Input
              {...item}
              key={item.name}
              value={values[item.name]}
              error={errors[item.name]}
              initialValues={initialValues}
              handleChange={(e:any) => {
                setFieldValue(item.name, e);
              }}
              handleBlur={() => setFieldTouched(item.name, false)}
              handleFocus={() => setFieldTouched(item.name, true)}
              extraText={props.extraText && props.extraText[item.name]}
            />)}
            <EuiButton
              isLoading={props.loading}
              onClick={()=> handleSubmit()}
              disabled={!isValid || !dirty}
              style={{fontSize:12, fontWeight:600}}
            >
              {props.buttonText || "Save Changes"}
            </EuiButton>
          </Wrap>
        );
      }}
    </Formik>
  );
}

const Wrap = styled.div`
  display: flex;
  flex-direction: column;
  align-content: center;
  justify-content: space-evenly;
  height: 100%;
`;

type FormFieldType = 'text' | 'img' | 'number' | 'hidden'
export interface FormField {
  name: string
  type: FormFieldType
  label: string
  readOnly?: boolean
  required?: boolean
  validator?: any
  style?: any
  prepend?: string
}

function validator(config: FormField[]) {
  const shape:{[k:string]:any} = {};
  config.forEach((field) => {
    if (typeof field === "object") {
      shape[field.name] = field.validator;
    }
  });
  return Yup.object().shape(shape);
}
