import React, { useState, useRef, useEffect } from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import styled from 'styled-components';
import Input from './inputs';
import { Button, IconButton, Modal } from '../sphinxUI';
import { useStores } from '../store';
import Select from '../sphinxUI/select';
import { dynamicSchemasByType, dynamicSchemaAutofillFieldsByType } from './schema';
import { formDropdownOptions } from '../people/utils/constants';

export default function Form(props: any) {
  const { buttonsOnBottom, wrapStyle, smallForm } = props;
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [dynamicInitialValues, setDynamicInitialValues]: any = useState(null);
  const [dynamicSchema, setDynamicSchema]: any = useState(null);
  const [dynamicSchemaName, setDynamicSchemaName] = useState('');
  const [showSettings, setShowSettings] = useState(false);
  const [showDeleteWarn, setShowDeleteWarn] = useState(false);
  const [disableFormButtons, setDisableFormButtons] = useState(false);
  const refBody: any = useRef(null);
  const { main, ui } = useStores();

  const firstTimeScreenData = ['pubkey', 'owner_alias', 'description', 'price_to_meet', 'twitter'];

  let lastPage = 1;
  const { readOnly } = props;
  const scrollDiv = props.scrollDiv ? props.scrollDiv : refBody;

  const initValues = dynamicInitialValues || props.initialValues;

  useEffect(() => {
    const dSchema = props.schema?.find((f) => f.defaultSchema);
    const type = props.initialValues?.type;
    if (dSchema && type) {
      const editSchema = dynamicSchemasByType[type];
      setDynamicSchema(editSchema);
      setDynamicSchemaName(type);
    } else if (dSchema) {
      setDynamicSchema(dSchema.defaultSchema);
      setDynamicSchemaName(dSchema.defaultSchemaName);
    }
    setLoading(false);
  }, []);

  // this useEffect triggers when the dynamic schema name is updated
  // checks if there are autofill fields that we can pull from local storage

  useEffect(() => {
    const formRef = props.formRef?.current;
    const vals = formRef && formRef.values;
    if (vals) {
      if (dynamicSchemaAutofillFieldsByType[dynamicSchemaName]) {
        Object.keys(dynamicSchemaAutofillFieldsByType[dynamicSchemaName]).forEach((k) => {
          const localStorageKey = dynamicSchemaAutofillFieldsByType[dynamicSchemaName][k];
          const valueToAssign = ui[localStorageKey];
          // if no value exists already
          if (!vals[k] || vals[k] == undefined) {
            if (valueToAssign) {
              setDynamicInitialValues({ ...initValues, [k]: valueToAssign });
              // re-render
              reloadForm();
            }
          }
        });
      }
    }
  }, [dynamicSchemaName]);

  useEffect(() => {
    scrollToTop();
  }, [page]);

  function reloadForm() {
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
    }, 20);
  }

  function scrollToTop() {
    if (scrollDiv && scrollDiv.current) {
      scrollDiv.current.scrollTop = 0;
    }
  }

  if (props.paged) {
    props.schema?.forEach((s) => {
      if (s.page > lastPage) lastPage = s.page;
    });
  }

  let schema = props.paged ? props.schema?.filter((f) => f.page === page) : props.schema;

  // replace schema with dynamic schema if there is one
  schema = dynamicSchema || schema;

  // if no schema, return empty div
  if (loading || !schema) return <div />;

  const buttonAlignment = buttonsOnBottom
    ? { zIndex: 20, bottom: 0, height: 108, justifyContent: 'center' }
    : { top: 0 };
  const formPad = buttonsOnBottom ? { paddingTop: 30 } : {};

  const buttonStyle = buttonsOnBottom ? { width: '80%', height: 48 } : {};

  const isAboutMeForm = schema?.find((f) => f.name === 'owner_alias') ? true : false;

  const dynamicFormOptions =
    (props.schema && props.schema[0] && formDropdownOptions[props.schema[0].dropdownOptions]) || [];

  // inject owner tribes
  const tribesSelectorIndex = schema?.findIndex((f) => f.name === 'tribe' || f.name === 'tribes');
  if (tribesSelectorIndex > -1) {
    // give "none" option
    schema[tribesSelectorIndex].options = [{ value: 'none', label: 'None' }];
    // add tribes
    main.ownerTribes?.length &&
      main.ownerTribes.forEach((ot) => {
        schema[tribesSelectorIndex].options.push({
          ...ot,
          value: ot.unique_name,
          label: ot.name
        });
      });
  }

  return (
    <Formik
      initialValues={initValues || {}}
      onSubmit={props.onSubmit}
      innerRef={props.formRef}
      validationSchema={validator(schema)}
    >
      {({
        setFieldTouched,
        handleSubmit,
        values,
        setFieldValue,
        errors,
        dirty,
        isValid,
        initialValues
      }) => {
        return (
          <Wrap ref={refBody} style={{ ...formPad, ...wrapStyle }}>
            {/* schema flipping dropdown */}
            {dynamicSchema && (
              <Select
                style={{ marginBottom: 14 }}
                onChange={(v) => {
                  console.log('v', v);
                  const selectedOption = dynamicFormOptions?.find((f) => f.value === v);
                  if (selectedOption) {
                    setDynamicSchemaName(v);
                    setDynamicSchema(selectedOption.schema);
                  }
                }}
                options={dynamicFormOptions}
                value={dynamicSchemaName}
              />
            )}

            {props.isFirstTimeScreen && schema ? (
              <>
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    width: '100%'
                  }}
                >
                  <div style={{ marginRight: '40px' }}>
                    {schema
                      .filter((item: FormField) => item.type === 'img')
                      .map((item: FormField) => (
                        <Input
                          {...item}
                          key={item.name}
                          values={values}
                          // disabled={readOnly}
                          // readOnly={readOnly}
                          errors={errors}
                          scrollToTop={scrollToTop}
                          value={values[item.name]}
                          error={errors[item.name]}
                          initialValues={initialValues}
                          deleteErrors={() => {
                            if (errors[item.name]) delete errors[item.name];
                          }}
                          handleChange={(e: any) => {
                            setFieldValue(item.name, e);
                          }}
                          setFieldValue={(e, f) => {
                            setFieldValue(e, f);
                          }}
                          setFieldTouched={setFieldTouched}
                          handleBlur={() => setFieldTouched(item.name, false)}
                          handleFocus={() => setFieldTouched(item.name, true)}
                          setDisableFormButtons={setDisableFormButtons}
                          extraHTML={
                            (props.extraHTML && props.extraHTML[item.name]) || item.extraHTML
                          }
                          borderType={'bottom'}
                          imageIcon={true}
                        />
                      ))}
                  </div>

                  <div style={{ width: '100%' }}>
                    {schema
                      .filter((item: FormField) => item.type !== 'img')
                      .map((item: FormField) => {
                        return (
                          <Input
                            {...item}
                            key={item.name}
                            values={values}
                            // disabled={readOnly}
                            // readOnly={readOnly}
                            errors={errors}
                            scrollToTop={scrollToTop}
                            value={values[item.name]}
                            error={errors[item.name]}
                            initialValues={initialValues}
                            deleteErrors={() => {
                              if (errors[item.name]) delete errors[item.name];
                            }}
                            handleChange={(e: any) => {
                              setFieldValue(item.name, e);
                            }}
                            setFieldValue={(e, f) => {
                              setFieldValue(e, f);
                            }}
                            setFieldTouched={setFieldTouched}
                            handleBlur={() => setFieldTouched(item.name, false)}
                            handleFocus={() => setFieldTouched(item.name, true)}
                            setDisableFormButtons={setDisableFormButtons}
                            extraHTML={
                              (props.extraHTML && props.extraHTML[item.name]) || item.extraHTML
                            }
                            borderType={'bottom'}
                          />
                        );
                      })}
                  </div>
                </div>
              </>
            ) : (
              schema.map((item: FormField) => (
                <Input
                  {...item}
                  key={item.name}
                  values={values}
                  // disabled={readOnly}
                  // readOnly={readOnly}
                  errors={errors}
                  scrollToTop={scrollToTop}
                  value={values[item.name]}
                  error={errors[item.name]}
                  initialValues={initialValues}
                  deleteErrors={() => {
                    if (errors[item.name]) delete errors[item.name];
                  }}
                  handleChange={(e: any) => {
                    setFieldValue(item.name, e);
                  }}
                  setFieldValue={(e, f) => {
                    setFieldValue(e, f);
                  }}
                  setFieldTouched={setFieldTouched}
                  handleBlur={() => setFieldTouched(item.name, false)}
                  handleFocus={() => setFieldTouched(item.name, true)}
                  setDisableFormButtons={setDisableFormButtons}
                  extraHTML={(props.extraHTML && props.extraHTML[item.name]) || item.extraHTML}
                />
              ))
            )}

            {/* make space at bottom for first sign up */}
            {buttonsOnBottom && !smallForm && <div style={{ height: 48, minHeight: 48 }} />}
            <BWrap style={buttonAlignment}>
              {props.close && buttonsOnBottom ? (
                <Button
                  disabled={disableFormButtons || props.loading}
                  onClick={() => {
                    if (props.close) props.close();
                  }}
                  style={{ ...buttonStyle, marginRight: 10, width: '140px' }}
                  color={'white'}
                  text={'Cancel'}
                />
              ) : (
                <IconButton
                  icon="arrow_back"
                  onClick={() => {
                    if (props.close) props.close();
                  }}
                  style={{ fontSize: 12, fontWeight: 600 }}
                />
              )}

              {readOnly ? (
                <div />
              ) : (
                <div style={{ display: 'flex', alignItems: 'center' }}>
                  <Button
                    disabled={disableFormButtons || props.loading}
                    onClick={() => {
                      if (dynamicSchemaName) {
                        // inject type in body
                        setFieldValue('type', dynamicSchemaName);
                      }
                      handleSubmit();
                      // if (lastPage === page) handleSubmit()
                      // else {
                      //   // this does form animation between pages
                      //   setFormMounted(false)
                      //   await sleep(200)
                      //   //background
                      //   setPage(page + 1)
                      // }
                    }}
                    loading={props.loading}
                    style={{ ...buttonStyle, width: '140px' }}
                    color={'primary'}
                    text={props.submitText || 'Save'}
                  />

                  {props.delete && (
                    <IconButton
                      disabled={disableFormButtons || props.loading}
                      onClick={() => {
                        props.delete();
                      }}
                      icon={'delete'}
                      loading={props.loading}
                      style={{ marginLeft: 10 }}
                      color={'clear'}
                    />
                  )}
                </div>
              )}
            </BWrap>
            {/*  if schema is AboutMe */}
            {!props.isFirstTimeScreen && isAboutMeForm && ui.meInfo?.id != 0 && (
              <>
                <div
                  style={{
                    cursor: 'pointer',
                    marginTop: 20,
                    fontSize: 12,
                    minHeight: 30,
                    height: 30
                  }}
                  onClick={() => setShowSettings(!showSettings)}
                >
                  Advanced Settings {showSettings ? '-' : '+'}
                </div>

                {showSettings && (
                  <div style={{ minHeight: 50, height: 50 }}>
                    <Button
                      text={'Delete my account'}
                      color={'link2'}
                      width="fit-content"
                      onClick={() => setShowDeleteWarn(true)}
                    />
                  </div>
                )}

                <Modal visible={showDeleteWarn}>
                  <div style={{ padding: 40, textAlign: 'center' }}>
                    <div style={{ fontSize: 30, marginBottom: 10 }}>Danger zone</div>
                    <p>
                      Are you sure? Doing so will delete your profile and <b>all of your posts.</b>
                    </p>

                    <div
                      style={{
                        width: '100%',
                        display: 'flex',
                        flexDirection: 'column',
                        justifyContent: 'center',
                        alignItems: 'center',
                        marginTop: 20
                      }}
                    >
                      <Button
                        text={'Nevermind'}
                        color={'white'}
                        onClick={() => {
                          setShowSettings(false);
                          setShowDeleteWarn(false);
                        }}
                      />
                      <div style={{ height: 20 }} />
                      <Button
                        text={'Delete everything'}
                        color={'danger'}
                        onClick={() => main.deleteProfile()}
                      />
                    </div>
                  </div>
                </Modal>
              </>
            )}
          </Wrap>
        );
      }}
    </Formik>
  );
}

const Wrap = styled.div`
  padding: 10px;
  padding-top: 80px;
  margin-bottom: 100px;
  display: flex;
  height: inherit;
  flex-direction: column;
  align-content: center;
  // max-width:400px;
  min-width: 230px;
`;

interface BWrapProps {
  readonly floatingButtons: boolean;
}

const BWrap = styled.div`
  display: flex;
  justify-content: space-between !important;
  align-items: center;
  width: 100%;
  padding: 10px;
  min-height: 42px;
  position: absolute;
  left: 0px;
  background: #ffffff;
  z-index: 10;
  box-shadow: 0px 1px 6px rgba(0, 0, 0, 0.07);
`;

type FormFieldType =
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
}

function validator(config: FormField[]) {
  const shape: { [k: string]: any } = {};
  config.forEach((field) => {
    if (typeof field === 'object') {
      shape[field.name] = field.validator;
    }
  });
  return Yup.object().shape(shape);
}
