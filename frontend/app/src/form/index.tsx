import React, { useState, useRef, useEffect } from "react";
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
  const refBody: any = useRef(null)

  let lastPage = 1

  const scrollDiv = props.scrollDiv ? props.scrollDiv : refBody

  useEffect(() => {
    scrollToTop()
  }, [page])

  if (props.paged) {
    props.schema.forEach((s) => {
      if (s.page > lastPage) lastPage = s.page
    })
  }
  function scrollToTop() {
    if (scrollDiv && scrollDiv.current) {
      scrollDiv.current.scrollTop = 0
    }
  }

  const schema = props.paged ? props.schema.filter(f => f.page === page) : props.schema

  let buttonText = props.buttonText || "Save Changes"
  if (lastPage !== page) buttonText = 'Next'

  return (
    <Formik
      initialValues={props.initialValues || {}}
      onSubmit={props.onSubmit}
      innerRef={props.formRef}
      validationSchema={validator(props.schema)}
      style={{ height: 'inherit' }}
      innerStyle={{ height: 'inherit' }}
    >
      {({ setFieldTouched, handleSubmit, values, setFieldValue, errors, dirty, isValid, initialValues }) => {

        console.log('errors', errors)
        return (
          <FadeLeft
            alwaysRender
            noFadeOnInit
            isMounted={formMounted}
            dismountCallback={() => setFormMounted(true)}
          >
            <Wrap ref={refBody}>
              {schema && schema.map((item: FormField) => <Input
                {...item}
                key={item.name}
                values={values}
                errors={errors}
                scrollToTop={scrollToTop}
                value={values[item.name]}
                error={errors[item.name]}
                initialValues={initialValues}
                deleteErrors={() => {
                  if (errors[item.name]) delete errors[item.name]
                }}
                handleChange={(e: any) => {
                  setFieldValue(item.name, e);
                }}
                setFieldValue={(e, f) => {
                  setFieldValue(e, f)
                }}
                setFieldTouched={setFieldTouched}
                handleBlur={() => setFieldTouched(item.name, false)}
                handleFocus={() => setFieldTouched(item.name, true)}
                setDisableFormButtons={setDisableFormButtons}
                extraHTML={props.extraHTML && props.extraHTML[item.name]}
              />)}

              <FadeLeft isMounted={!disableFormButtons}>
                <BWrap floatingButtons={props.floatingButtons}>
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
              </FadeLeft>

            </Wrap >

          </FadeLeft >
        );
      }}
    </Formik >
  );
}

const Wrap = styled.div`
  display: flex;
  height:inherit;
  flex-direction: column;
  align-content: center;
  justify-content: space-between;
`;

interface BWrapProps {
  readonly floatingButtons: boolean;
}

const BWrap = styled.div<BWrapProps>`
  display: flex;
  justify-content: space-evenly;
  align-items:center;
  width:100%;
  height:42px;
  min-height:42px;
  margin-top:20px;
  position:${p => p.floatingButtons && 'absolute'};
  bottom:${p => p.floatingButtons && '0px'};
  left:${p => p.floatingButtons && '0px'};
`;

type FormFieldType = 'text' | 'textarea' | 'img' | 'gallery' | 'number' | 'hidden' | 'widgets' | 'widget' | 'switch'

type FormFieldClass = 'twitter' | 'blog' | 'offer' | 'wanted' | 'supportme'

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