import React, { useState } from "react";
import { Formik } from "formik";
import * as Yup from "yup";
import styled from "styled-components";
import Input from "./inputs";
import { EuiButton } from '@elastic/eui'
import FadeLeft from '../animated/fadeLeft';

const sleep = ms => new Promise(resolve => setTimeout(resolve, ms))

export default function Form(props: any) {

  const [page, setPage] = useState(1)
  const [formMounted, setFormMounted] = useState(true)
  const [disableFormButtons, setDisableFormButtons] = useState(false)

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
      style={{ height: '100%' }}
    >
      {({ setFieldTouched, handleSubmit, values, setFieldValue, errors, dirty, isValid, initialValues }) => {

        return (
          <FadeLeft
            alwaysRender
            noFadeOnInit
            isMounted={formMounted}
            dismountCallback={() => setFormMounted(true)}
          >
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
                setDisableFormButtons={setDisableFormButtons}
                extraHTML={props.extraHTML && props.extraHTML[item.name]}
              />)}

              <BWrap>

                {page > 1 &&
                  <EuiButton
                    disabled={disableFormButtons || props.loading}
                    onClick={async () => {
                      // this does form animation between pages
                      setFormMounted(false)
                      await sleep(200)
                      //
                      setPage(page - 1)
                    }}
                    style={{ fontSize: 12, fontWeight: 600 }}
                  >
                    Back
                  </EuiButton>
                }

                <EuiButton
                  isLoading={props.loading}
                  onClick={async () => {
                    if (lastPage === page) handleSubmit()
                    else {
                      // this does form animation between pages
                      setFormMounted(false)
                      await sleep(200)
                      //
                      setPage(page + 1)
                    }
                  }}
                  disabled={disableFormButtons || !isValid}
                  style={{ fontSize: 12, fontWeight: 600 }}
                >
                  {buttonText}
                </EuiButton>
              </BWrap>
            </Wrap >

          </FadeLeft >
        );
      }}
    </Formik >
  );
}

const Wrap = styled.div`
  display: flex;
  flex:1;
  flex-direction: column;
  align-content: center;
  justify-content: space-between;
`;

const BWrap = styled.div`
  display: flex;
  justify-content: space-evenly;
  align-items:flex-end;
  width:100%;
  flex:1;
  margin-top:20px;
`;

type FormFieldType = 'text' | 'img' | 'number' | 'hidden' | 'widgets' | 'widget'

type FormFieldClass = 'twitter' | 'blog' | 'offer' | 'wanted' | 'donations'

export interface FormField {
  name: string
  type: FormFieldType
  class?: FormFieldClass
  label: string
  itemLabel?: string
  single?: boolean
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