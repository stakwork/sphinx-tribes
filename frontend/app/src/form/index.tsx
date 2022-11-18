import React, { useState, useRef, useEffect, useCallback } from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import styled from 'styled-components';
import Input from './inputs';
import { Button, IconButton, Modal } from '../sphinxUI';
import { useStores } from '../store';
import Select from '../sphinxUI/select';
import { dynamicSchemasByType, dynamicSchemaAutofillFieldsByType } from './schema';
import { formDropdownOptions } from '../people/utils/constants';
import { EuiText } from '@elastic/eui';
import api from '../api';

const BountyDetailsCreationData = {
  step_1: {
    step: 1,
    heading: 'Basic info',
    sub_heading: 'Nemo enim ipsam voluptatem quia voluptas sit magni voluptatem sequi.',
    schema: ['wanted_type', 'one_sentence_summary'],
    schema2: ['ticketUrl', 'github_description', 'description']
  },
  step_2: {
    step: 2,
    heading: 'Price and Estimate',
    sub_heading: 'Nemo enim ipsam voluptatem quia voluptas sit magni voluptatem sequi.',
    schema: ['price', 'codingLanguage', 'tribe', 'estimate_session_length'],
    schema2: ['estimate_complete_date', 'deliverables', 'show']
  },
  step_3: {
    step: 3,
    heading: 'Invite Developer',
    sub_heading: 'Nemo enim ipsam quia voluptas sit',
    schema: ['assignee'],
    schema2: ['']
  }
};

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
  const [peopleList, setPeopleList] = useState<any>();
  const refBody: any = useRef(null);
  const { main, ui } = useStores();

  const [schemaData, setSchemaData] = useState(BountyDetailsCreationData.step_1);
  const [stepTracker, setStepTracker] = useState<number>(1);

  let lastPage = 1;
  const { readOnly } = props;
  const scrollDiv = props.scrollDiv ? props.scrollDiv : refBody;

  const initValues = dynamicInitialValues || props.initialValues;

  const NextStepHandler = useCallback(() => {
    setStepTracker(stepTracker < 3 ? stepTracker + 1 : stepTracker);
  }, [stepTracker]);

  const PreviousStepHandler = useCallback(() => {
    setStepTracker(stepTracker > 1 ? stepTracker - 1 : stepTracker);
  }, [stepTracker]);

  useEffect(() => {
    switch (stepTracker) {
      case 1:
        setSchemaData(BountyDetailsCreationData.step_1);
        break;
      case 2:
        setSchemaData(BountyDetailsCreationData.step_2);
        break;
      case 3:
        setSchemaData(BountyDetailsCreationData.step_3);
        break;
      default:
        return;
    }
  }, [stepTracker]);

  useEffect(() => {
    (async () => {
      try {
        const response = await api.get(`people?page=1&search=&sortBy=last_login&limit=100`);
        setPeopleList(response);
      } catch (error) {
        console.log(error);
      }
    })();
  }, []);

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
      validationSchema={validator(schema)}>
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
          <Wrap
            ref={refBody}
            style={{
              ...formPad,
              ...wrapStyle,
              height: stepTracker === 3 ? '592px' : '560px',
              minWidth: stepTracker === 3 ? '388px' : '712px',
              maxWidth: stepTracker === 3 ? '388px' : '712px'
            }}
            newDesign={props?.newDesign}>
            {/* schema flipping dropdown */}
            {/* {dynamicSchema && (
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
            )} */}

            {props.isFirstTimeScreen && schema ? (
              <>
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    width: '100%'
                  }}>
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
                          label={''}
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
            ) : props?.newDesign ? (
              <>
                <div
                  style={{
                    display: 'flex',
                    flexDirection: 'column',
                    justifyContent: 'center',
                    padding: '0 48px'
                  }}>
                  <EuiText
                    style={{
                      fontFamily: 'Barlow',
                      fontSize: '15px',
                      fontWeight: '500',
                      lineHeight: '18px',
                      letterSpacing: '0.06em'
                    }}>{`STEP ${schemaData.step}/3`}</EuiText>
                  <EuiText
                    style={{
                      fontFamily: 'Barlow',
                      fontSize: '36px',
                      fontWeight: '800',
                      lineHeight: '43px',
                      color: '#3C3F41',
                      marginBottom: '26px',
                      marginTop: '16px'
                    }}>
                    {schemaData.heading}
                  </EuiText>
                  <EuiText
                    style={{
                      fontFamily: 'Barlow',
                      fontSize: '17px',
                      fontWeight: '400',
                      lineHeight: '20px',
                      color: '#292C33',
                      marginBottom: schemaData.step === 1 ? '29px' : '37px'
                    }}>
                    {schemaData.sub_heading}
                  </EuiText>
                </div>

                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    height: '100%',
                    padding: '0px 48px'
                  }}>
                  <div
                    style={{
                      width: '292px'
                    }}>
                    {schemaData.step === 1 && dynamicSchema && (
                      <Select
                        style={{ marginBottom: 24 }}
                        onChange={(v) => {
                          console.log('v', v);
                          const selectedOption = dynamicFormOptions?.find((f) => f.value === v);
                          if (selectedOption) {
                            setDynamicSchemaName(v);
                            setDynamicSchema(selectedOption.schema);
                          }
                        }}
                        handleActive={() => {}}
                        options={dynamicFormOptions}
                        value={dynamicSchemaName}
                      />
                    )}
                    {schema
                      .filter((item) => schemaData.schema.includes(item.name))
                      .map((item: FormField) => (
                        <Input
                          {...item}
                          key={item.name}
                          newDesign={true}
                          values={values}
                          peopleList={peopleList}
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
                        />
                      ))}
                  </div>
                  <div
                    style={{
                      width: '292px'
                    }}>
                    {schema
                      .filter((item) => schemaData.schema2.includes(item.name))
                      .map((item: FormField) => (
                        <Input
                          {...item}
                          peopleList={peopleList}
                          newDesign={true}
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
                        />
                      ))}
                  </div>
                </div>
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    padding: '0px 48px'
                  }}>
                  {stepTracker < 3 && (
                    <EuiText
                      style={{
                        fontSize: '13px',
                        fontFamily: 'Barlow',
                        fontWeight: '400',
                        lineHeight: '35px',
                        color: '#B0B7BC'
                      }}>
                      * Required
                    </EuiText>
                  )}
                  <div
                    style={{
                      display: 'flex',
                      flexDirection: 'row-reverse',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      width: stepTracker < 3 ? '45%' : '100%'
                    }}>
                    <div
                      style={{
                        width: '120px',
                        height: '42px',
                        display: 'flex',
                        justifyContent: 'center',
                        alignItems: 'center',
                        cursor: 'pointer',
                        background: '#618AFF',
                        boxShadow: '0px 2px 10px rgba(97, 138, 255, 0.5)',

                        borderRadius: '32px',
                        color: '#fff'
                      }}
                      onClick={() => {
                        if (schemaData.step === 3) {
                          if (dynamicSchemaName) {
                            // inject type in body
                            setFieldValue('type', dynamicSchemaName);
                          }
                          handleSubmit();
                        } else {
                          NextStepHandler();
                        }
                      }}>
                      <EuiText
                        style={{
                          fontFamily: 'Barlow',
                          fontSize: '16px',
                          fontWeight: '600',
                          lineHeight: '19px',
                          userSelect: 'none'
                        }}>
                        {schemaData.step === 3 ? 'Skip' : 'Next'}
                      </EuiText>
                    </div>
                    {schemaData.step > 1 && (
                      <div
                        style={{
                          width: '120px',
                          height: '42px',
                          display: 'flex',
                          justifyContent: 'center',
                          alignItems: 'center',
                          cursor: 'pointer',
                          background: '#fff',
                          border: '1px solid #DDE1E5',
                          borderRadius: '32px',
                          color: '#5F6368'
                        }}
                        onClick={PreviousStepHandler}>
                        <EuiText
                          style={{
                            fontFamily: 'Barlow',
                            fontSize: '16px',
                            fontWeight: '600',
                            lineHeight: '19px',
                            userSelect: 'none'
                          }}>
                          Back
                        </EuiText>
                      </div>
                    )}
                  </div>
                </div>
              </>
            ) : (
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'center',
                  width: '100%'
                }}>
                <div
                  style={{
                    width: '100%'
                  }}>
                  {schema.map((item: FormField) => (
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
                  ))}
                </div>
              </div>
            )}

            {/* make space at bottom for first sign up */}
            {buttonsOnBottom && !smallForm && <div style={{ height: 48, minHeight: 48 }} />}
            {!props?.newDesign && (
              <BWrap style={buttonAlignment}>
                {props?.close && buttonsOnBottom ? (
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
                      text={'Save'}
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
            )}
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
                  onClick={() => setShowSettings(!showSettings)}>
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
                      }}>
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

interface WrapProps {
  newDesign?: string;
}

const Wrap = styled.div<WrapProps>`
  padding: ${(p) => (p?.newDesign ? '28px 0px' : '80px 0px 0px 0px')};
  margin-bottom: ${(p) => !p?.newDesign && '100px'};
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
