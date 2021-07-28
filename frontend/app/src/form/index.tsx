import React, { useState } from "react";
import { Formik } from "formik";
import * as Yup from "yup";
import styled from "styled-components";
import Input from "./inputs";
import { EuiButton } from '@elastic/eui'

export default function Form(props: any) {

  const [page, setPage] = useState(1)

  let lastPage = 1

  if (props.paged) {
    props.schema.forEach((s) => {
      if (s.page > lastPage) lastPage = s.page
    })
  }

  const schema = props.paged ? props.schema.filter(f => f.page === page) : props.schema

  let buttonText = props.buttonText || "Save Changes"
  if (lastPage !== page) buttonText = 'Next'

  return (
    <Formik
      initialValues={props.initialValues || {}}
      onSubmit={props.onSubmit}
      validationSchema={validator(props.schema)}
    >
      {({ setFieldTouched, handleSubmit, values, setFieldValue, errors, dirty, isValid, initialValues }) => {

        return (
          <Wrap>
            {schema && schema.map((item: FormField) => <Input
              {...item}
              key={item.name}
              values={values}
              errors={errors}
              value={values[item.name]}
              error={errors[item.name]}
              initialValues={initialValues}
              handleChange={(e: any) => {
                setFieldValue(item.name, e);
              }}
              setFieldValue={setFieldValue}
              setFieldTouched={setFieldTouched}
              handleBlur={() => setFieldTouched(item.name, false)}
              handleFocus={() => setFieldTouched(item.name, true)}
              extraHTML={props.extraHTML && props.extraHTML[item.name]}
            />)}

            <BWrap>

              {page > 1 &&
                <EuiButton
                  disabled={props.loading}
                  onClick={() => {
                    setPage(page - 1)
                  }}
                  style={{ fontSize: 12, fontWeight: 600 }}
                >
                  Back
                </EuiButton>
              }

              <EuiButton
                isLoading={props.loading}
                onClick={() => {
                  if (lastPage === page) handleSubmit()
                  else setPage(page + 1)
                }}
                disabled={!isValid}
                style={{ fontSize: 12, fontWeight: 600 }}
              >
                {buttonText}
              </EuiButton>

            </BWrap>
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

const BWrap = styled.div`
  display: flex;
  justify-content: space-evenly;
  margin-top:10px;
`;

type FormFieldType = 'text' | 'img' | 'number' | 'hidden' | 'widgets' | 'widget'

export interface FormField {
  name: string
  type: FormFieldType
  label: string
  readOnly?: boolean
  required?: boolean
  validator?: any
  style?: any
  prepend?: string
  page?: number
  extras?: FormField[]
  fields?: FormField[]
  icon?: string
}

function validator(config: FormField[]) {
  const shape: { [k: string]: any } = {};
  config.forEach((field) => {
    if (typeof field === "object") {
      shape[field.name] = field.validator;
    }
  });
  return Yup.object().shape(shape);
}