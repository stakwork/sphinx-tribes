import React, { useState, useRef, useEffect } from "react";
import { Formik } from "formik";
import * as Yup from "yup";
import styled from "styled-components";
import Input from "./inputs";

import FadeLeft from '../animated/fadeLeft';
import { Button, IconButton } from "../sphinxUI";

const sleep = ms => new Promise(resolve => setTimeout(resolve, ms))

export default function Form(props: any) {

  const [page, setPage] = useState(1)
  const [formMounted, setFormMounted] = useState(true)
  const [disableFormButtons, setDisableFormButtons] = useState(false)
  const refBody: any = useRef(null)

  let lastPage = 1

  const readOnly = props.readOnly

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
    // style={{ height: 'inherit' }}
    // innerStyle={{ height: 'inherit' }}
    >
      {({ setFieldTouched, handleSubmit, values, setFieldValue, errors, dirty, isValid, initialValues }) => {

        return (
          <Wrap ref={refBody}>
            {schema && schema.map((item: FormField) => <Input
              {...item}
              key={item.name}
              values={values}
              disabled={readOnly}
              readOnly={readOnly}
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

            <BWrap >

              <IconButton
                isLoading={props.loading}
                icon='arrow_back'
                onClick={() => {
                  if (props.close) props.close()
                }}
                disabled={disableFormButtons || !isValid}
                style={{ fontSize: 12, fontWeight: 600 }}
              />

              {readOnly ? <div /> :
                <Button
                  disabled={disableFormButtons || props.loading}
                  onClick={() => {
                    handleSubmit()
                    // if (lastPage === page) handleSubmit()
                    // else {
                    //   // this does form animation between pages
                    //   setFormMounted(false)
                    //   await sleep(200)
                    //   //
                    //   setPage(page + 1)
                    // }
                  }}
                  color={'primary'}
                  text={props.submitText || 'Save'}
                />
              }



            </BWrap>

          </Wrap >

        );
      }}
    </Formik >
  );
}

const Wrap = styled.div`
  padding:10px;
  padding-top:80px;
  margin-bottom:100px;
  display: flex;
  height:inherit;
  flex-direction: column;
  align-content: center;
`;

interface BWrapProps {
  readonly floatingButtons: boolean;
}

const BWrap = styled.div`
  display: flex;
  justify-content: space-between;
  align-items:center;
  width:100%;
  padding:10px;
  min-height:42px;
  position: absolute;
  top:0px;
  left:0px;
  background:#ffffff;
  box-shadow: 0px 1px 6px rgba(0, 0, 0, 0.07);
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