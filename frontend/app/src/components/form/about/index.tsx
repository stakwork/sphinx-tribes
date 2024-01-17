import { Formik } from 'formik';
import { observer } from 'mobx-react-lite';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import history from 'config/history';
import { colors } from '../../../config/colors';
import { useStores } from '../../../store';
import { Button, IconButton, Modal } from '../../common';
import Input from '../inputs';
import { dynamicSchemaAutofillFieldsByType, dynamicSchemasByType } from '../schema';
import { AboutSchemaInner, BWrap, SchemaOuterContainer, Wrap } from '../style';
import { FormField, validator } from '../utils';
import { FormProps } from '../interfaces';

function Form(props: FormProps) {
  const { buttonsOnBottom, smallForm, readOnly, scrollDiv: scrollRef, initialValues } = props;
  const page = 1;
  const [loading, setLoading] = useState(true);
  const [dynamicInitialValues, setDynamicInitialValues]: any = useState(null);
  const [dynamicSchema, setDynamicSchema]: any = useState(null);
  const [dynamicSchemaName, setDynamicSchemaName] = useState('');
  const [showSettings, setShowSettings] = useState(false);
  const [showDeleteWarn, setShowDeleteWarn] = useState(false);
  const [disableFormButtons, setDisableFormButtons] = useState(false);
  const refBody: any = useRef(null);
  const { main, ui } = useStores();
  const color = colors['light'];
  const [isFocused, setIsFocused] = useState({});

  const scrollDiv = scrollRef ?? refBody;

  const initValues = dynamicInitialValues || initialValues;

  useEffect(() => {
    const dSchema = props.schema?.find((f: any) => f.defaultSchema);
    const type = initialValues?.type;
    if (dSchema && type) {
      const editSchema = dynamicSchemasByType[type];
      setDynamicSchema(editSchema);
      setDynamicSchemaName(type);
    } else if (dSchema) {
      setDynamicSchema(dSchema.defaultSchema);
      setDynamicSchemaName(dSchema.defaultSchemaName);
    }
    setLoading(false);
  }, [initialValues?.type, props.schema]);

  // this useEffect triggers when the dynamic schema name is updated
  // checks if there are autofill fields that we can pull from local storage

  function reloadForm() {
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
    }, 20);
  }

  useEffect(() => {
    const formRef = props.formRef?.current;
    const vals = formRef && formRef.values;
    if (vals) {
      if (dynamicSchemaAutofillFieldsByType[dynamicSchemaName]) {
        Object.keys(dynamicSchemaAutofillFieldsByType[dynamicSchemaName]).forEach((k: any) => {
          const localStorageKey = dynamicSchemaAutofillFieldsByType[dynamicSchemaName][k];
          const valueToAssign = ui[localStorageKey];
          // if no value exists already
          if (!vals[k] || vals[k] === undefined) {
            if (valueToAssign) {
              setDynamicInitialValues({ ...initValues, [k]: valueToAssign });
              // re-render
              reloadForm();
            }
          }
        });
      }
    }
  }, [dynamicSchemaName, initValues, props.formRef, ui]);

  const scrollToTop = useCallback(() => {
    if (scrollDiv && scrollDiv.current) {
      scrollDiv.current.scrollTop = 0;
    }
  }, [scrollDiv]);

  useEffect(() => {
    scrollToTop();
  }, [scrollToTop]);

  let schema = props.paged ? props.schema?.filter((f: any) => f.page === page) : props.schema;

  // replace schema with dynamic schema if there is one
  schema = dynamicSchema || schema;

  // if no schema, return empty div
  if (loading || !schema) return <div />;

  const buttonAlignment = buttonsOnBottom
    ? { zIndex: 20, bottom: 0, height: 108, justifyContent: 'center' }
    : { top: 0 };

  const buttonStyle = buttonsOnBottom ? { width: '80%', height: 48 } : {};

  const isAboutMeForm = schema?.find((f: any) => f.name === 'owner_alias') ? true : false;

  // inject owner tribes
  const tribesSelectorIndex = schema?.findIndex(
    (f: any) => f.name === 'tribe' || f.name === 'tribes'
  );

  if (tribesSelectorIndex > -1) {
    // give "none" option
    schema[tribesSelectorIndex].options = [{ value: 'none', label: 'None' }];
    // add tribes
    main.ownerTribes?.length &&
      main.ownerTribes.forEach((ot: any) => {
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
      {({ setFieldTouched, handleSubmit, values, setFieldValue, errors, initialValues }: any) => (
        <Wrap ref={refBody} style={{ width: '100%' }}>
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
                        setFieldValue={(e: any, f: any) => {
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
                        style={
                          item.name === 'github_description' && !values.ticket_url
                            ? {
                                display: 'none'
                              }
                            : undefined
                        }
                      />
                    ))}
                </div>

                <div style={{ width: '100%' }}>
                  {schema
                    .filter((item: FormField) => item.type !== 'img')
                    .map((item: FormField) => (
                      <Input
                        {...item}
                        key={item.name}
                        values={values}
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
                        setFieldValue={(e: any, f: any) => {
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
                        style={
                          item.name === 'github_description' && !values.ticket_url
                            ? {
                                display: 'none'
                              }
                            : undefined
                        }
                      />
                    ))}
                </div>
              </div>
            </>
          ) : (
            <AboutSchemaInner>
              <SchemaOuterContainer>
                <div className="SchemaInnerContainer">
                  {schema.map((item: FormField) => (
                    <Input
                      {...item}
                      key={item.name}
                      values={values}
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
                      setFieldValue={(e: any, f: any) => {
                        setFieldValue(e, f);
                      }}
                      setFieldTouched={setFieldTouched}
                      isFocused={isFocused}
                      handleBlur={() => {
                        setFieldTouched(item.name, false);
                        setIsFocused({ [item.label]: false });
                      }}
                      handleFocus={() => {
                        setFieldTouched(item.name, true);
                        setIsFocused({ [item.label]: true });
                      }}
                      setDisableFormButtons={setDisableFormButtons}
                      extraHTML={(props.extraHTML && props.extraHTML[item.name]) || item.extraHTML}
                      style={
                        item.name === 'github_description' && !values.ticket_url
                          ? {
                              display: 'none'
                            }
                          : undefined
                      }
                    />
                  ))}
                </div>
              </SchemaOuterContainer>
            </AboutSchemaInner>
          )}
          {/* make space at bottom for first sign up */}
          {buttonsOnBottom && !smallForm && <div style={{ height: 48, minHeight: 48 }} />}
          {!props?.newDesign && (
            <BWrap style={buttonAlignment} color={color}>
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
                      if (
                        window.location.href.includes('bounty') ||
                        window.location.href.includes('ticket')
                      ) {
                        history.push('/bounties');
                      }
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
                        if (props.delete) props.delete();
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
          {!props.isFirstTimeScreen && isAboutMeForm && ui.meInfo?.id !== 0 && (
            <>
              <SchemaOuterContainer>
                <div
                  className="SchemaInnerContainer"
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
              </SchemaOuterContainer>
              {showSettings && (
                <SchemaOuterContainer>
                  <div style={{ minHeight: 50, height: 50 }} className="SchemaInnerContainer">
                    <Button
                      text={'Delete my account'}
                      color={'link2'}
                      width="fit-content"
                      onClick={() => setShowDeleteWarn(true)}
                    />
                  </div>
                </SchemaOuterContainer>
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
      )}
    </Formik>
  );
}
export default observer(Form);
