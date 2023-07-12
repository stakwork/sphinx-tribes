import { EuiText } from '@elastic/eui';
import { Formik } from 'formik';
import { observer } from 'mobx-react-lite';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import api from '../../api';
import { colors } from '../../config/colors';
import { BountyDetailsCreationData } from '../../people/utils/BountyCreationConstant';
import { formDropdownOptions } from '../../people/utils/Constants';
import { useStores } from '../../store';
import { Button, Divider, IconButton, Modal } from '../common';
import ImageButton from '../common/ImageButton';
import Input from './inputs';
import { dynamicSchemaAutofillFieldsByType, dynamicSchemasByType } from './schema';
import {
  BWrap,
  BottomContainer,
  BountyContainer,
  ChooseBountyContainer,
  CreateBountyHeaderContainer,
  SchemaOuterContainer,
  SchemaTagsContainer,
  Wrap
} from './style';
import { FormField, validator } from './utils';
import { FormProps } from './interfaces';

function Form(props: FormProps) {
  const {
    buttonsOnBottom,
    wrapStyle,
    smallForm,
    readOnly,
    scrollDiv: scrollRef,
    initialValues
  } = props;
  const page = 1;
  const [loading, setLoading] = useState(true);
  const [dynamicInitialValues, setDynamicInitialValues]: any = useState(null);
  const [dynamicSchema, setDynamicSchema]: any = useState(null);
  const [dynamicSchemaName, setDynamicSchemaName] = useState('');
  const [showSettings, setShowSettings] = useState(false);
  const [showDeleteWarn, setShowDeleteWarn] = useState(false);
  const [disableFormButtons, setDisableFormButtons] = useState(false);
  const [peopleList, setPeopleList] = useState<any>();
  const [assigneeName, setAssigneeName] = useState<string>('');
  const refBody: any = useRef(null);
  const { main, ui } = useStores();
  const color = colors['light'];
  const [isFocused, setIsFocused] = useState({});

  const [schemaData, setSchemaData] = useState(BountyDetailsCreationData.step_1);
  const [stepTracker, setStepTracker] = useState<number>(1);

  let lastPage = 1;
  const scrollDiv = scrollRef ?? refBody;

  const initValues = dynamicInitialValues || initialValues;

  const NextStepHandler = useCallback(() => {
    setStepTracker(stepTracker < 5 ? stepTracker + 1 : stepTracker);
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
      case 4:
        setSchemaData(BountyDetailsCreationData.step_4);
        break;
      case 5:
        setSchemaData(BountyDetailsCreationData.step_5);
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

  if (props.paged) {
    props.schema?.forEach((s: any) => {
      if (s.page > lastPage) lastPage = s.page;
    });
  }

  let schema = props.paged ? props.schema?.filter((f: any) => f.page === page) : props.schema;

  // replace schema with dynamic schema if there is one
  schema = dynamicSchema || schema;

  // if no schema, return empty div
  if (loading || !schema) return <div />;

  const buttonAlignment = buttonsOnBottom
    ? { zIndex: 20, bottom: 0, height: 108, justifyContent: 'center' }
    : { top: 0 };
  const formPad = buttonsOnBottom ? { paddingTop: 30 } : {};

  const buttonStyle = buttonsOnBottom ? { width: '80%', height: 48 } : {};

  const isAboutMeForm = schema?.find((f: any) => f.name === 'owner_alias') ? true : false;

  const dynamicFormOptions =
    (props.schema && props.schema[0] && formDropdownOptions[props.schema[0].dropdownOptions]) || [];

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
      {({ setFieldTouched, handleSubmit, values, setFieldValue, errors, initialValues }: any) => {
        const isDescriptionValid = values.ticket_url
          ? values.github_description || !!values.description
          : !!values.description;

        const valid = schemaData.required.every((key: string) =>
          key === '' ? true : values?.[key]
        );

        const isBtnDisabled = !valid || (stepTracker === 3 && !isDescriptionValid);

        return (
          <Wrap
            ref={refBody}
            style={{
              ...formPad,
              ...wrapStyle,
              ...schemaData.outerContainerStyle
            }}
            newDesign={props?.newDesign}
          >
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
            ) : props?.newDesign ? (
              <>
                <CreateBountyHeaderContainer color={color}>
                  <div className="TopContainer">
                    <EuiText className="stepText">
                      {`STEP ${schemaData.step}`} <span className="stepTextSpan"> / 5</span>
                    </EuiText>
                    <EuiText className="schemaName">{schemaData.schemaName}</EuiText>
                  </div>
                  <EuiText className="HeadingText" style={schemaData.headingStyle}>
                    {schemaData.heading}
                  </EuiText>
                </CreateBountyHeaderContainer>

                {schemaData.step === 1 && dynamicSchema && (
                  <ChooseBountyContainer color={color}>
                    {dynamicFormOptions?.map((v: any) => (
                      <BountyContainer
                        key={v.label}
                        color={color}
                        show={v.value === 'freelance_job_request' ? true : false}
                      >
                        <div className="freelancerContainer">
                          <div
                            style={{
                              minHeight: '134px !important',
                              maxHeight: '134px !important',
                              height: '134px',
                              width: '290px',
                              background: color.white100,
                              borderRadius: '20px 20px 0px 0px'
                            }}
                          >
                            <div
                              style={{
                                height: '100%',
                                width: '100%',
                                display: 'flex',
                                justifyContent: 'center',
                                alignItems: 'flex-end',
                                position: 'relative'
                              }}
                            >
                              <img
                                src={
                                  v.value === 'freelance_job_request'
                                    ? '/static/freelancer_bounty.svg'
                                    : '/static/live_help.svg'
                                }
                                alt="select_type"
                                height={'114%'}
                                width={'114%'}
                                style={{
                                  position: 'absolute',
                                  top: '32px'
                                }}
                              />
                            </div>
                          </div>
                          <div className="TextButtonContainer">
                            <EuiText className="textTop">{v.label}</EuiText>
                            <EuiText className="textBottom">
                              {v.value === 'freelance_job_request'
                                ? 'Choose the right developer'
                                : 'Get instant help for your task'}
                            </EuiText>
                            {v.value === 'freelance_job_request' ? (
                              <div
                                className="StartButton"
                                onClick={() => {
                                  NextStepHandler();
                                  setDynamicSchemaName(v.value);
                                  setDynamicSchema(v.schema);
                                }}
                              >
                                Start
                              </div>
                            ) : (
                              <div className="ComingSoonContainer">
                                <Divider
                                  style={{
                                    width: '26px',
                                    background: color.grayish.G300
                                  }}
                                />
                                <EuiText className="ComingSoonText">Coming soon</EuiText>
                                <Divider
                                  style={{
                                    width: '26px',
                                    background: color.grayish.G300
                                  }}
                                />
                              </div>
                            )}
                          </div>
                        </div>
                      </BountyContainer>
                    ))}
                  </ChooseBountyContainer>
                )}

                {schemaData.step !== 1 && (
                  <>
                    <SchemaTagsContainer>
                      <div className="LeftSchema">
                        {schema
                          .filter((item: any) => schemaData.schema.includes(item.name))
                          .map((item: FormField) => (
                            <Input
                              {...item}
                              key={item.name}
                              newDesign={true}
                              values={values}
                              setAssigneefunction={item.name === 'assignee' && setAssigneeName}
                              peopleList={peopleList}
                              isFocused={isFocused}
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
                              handleBlur={() => {
                                setFieldTouched(item.name, false);
                                setIsFocused({ [item.label]: false });
                              }}
                              handleFocus={() => {
                                setFieldTouched(item.name, true);
                                setIsFocused({ [item.label]: true });
                              }}
                              setDisableFormButtons={setDisableFormButtons}
                              extraHTML={
                                (props.extraHTML && props.extraHTML[item.name]) || item.extraHTML
                              }
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
                      <div className="RightSchema">
                        {schema
                          .filter((item: any) => schemaData.schema2.includes(item.name))
                          .map((item: FormField) => (
                            <Input
                              {...item}
                              peopleList={peopleList}
                              newDesign={true}
                              key={item.name}
                              values={values}
                              testId={item.label}
                              errors={errors}
                              scrollToTop={scrollToTop}
                              value={values[item.name]}
                              error={errors[item.name]}
                              initialValues={initialValues}
                              deleteErrors={() => {
                                if (errors[item.name]) delete errors[item.name];
                              }}
                              isFocused={isFocused}
                              handleChange={(e: any) => {
                                setFieldValue(item.name, e);
                              }}
                              setFieldValue={(e: any, f: any) => {
                                setFieldValue(e, f);
                              }}
                              setFieldTouched={setFieldTouched}
                              handleBlur={() => {
                                setFieldTouched(item.name, false);
                                setIsFocused({ [item.label]: false });
                              }}
                              handleFocus={() => {
                                setFieldTouched(item.name, true);
                                setIsFocused({ [item.label]: true });
                              }}
                              setDisableFormButtons={setDisableFormButtons}
                              extraHTML={
                                (props.extraHTML && props.extraHTML[item.name]) || item.extraHTML
                              }
                              style={
                                item.type === 'loom' && values.ticket_url
                                  ? {
                                      marginTop: '55px'
                                    }
                                  : undefined
                              }
                            />
                          ))}
                      </div>
                    </SchemaTagsContainer>
                    <BottomContainer color={color} assigneeName={assigneeName} valid={valid}>
                      <EuiText className="RequiredText">{schemaData?.extraText}</EuiText>
                      <div
                        className="ButtonContainer"
                        style={{
                          width: stepTracker < 5 ? '45%' : '100%',
                          height: stepTracker < 5 ? '48px' : '48px',
                          marginTop: stepTracker === 5 || stepTracker === 3 ? '20px' : ''
                        }}
                      >
                        {isBtnDisabled && (
                          <div className="nextButtonDisable">
                            <EuiText className="disableText">Next</EuiText>
                          </div>
                        )}
                        {!isBtnDisabled && (
                          <div
                            className="nextButton"
                            onClick={() => {
                              if (schemaData.step === 5 && valid) {
                                if (dynamicSchemaName) {
                                  setFieldValue('type', dynamicSchemaName);
                                }
                                if (assigneeName !== '') {
                                  handleSubmit();
                                } else {
                                  setAssigneeName('a');
                                }
                              } else {
                                if (valid) {
                                  NextStepHandler();
                                }
                              }
                            }}
                            style={{
                              width:
                                schemaData.step === 5
                                  ? assigneeName === ''
                                    ? '145px'
                                    : '120px'
                                  : '120px'
                            }}
                          >
                            {assigneeName === '' ? (
                              <EuiText className="nextText">
                                {schemaData.step === 5 ? 'Decide Later' : 'Next'}
                              </EuiText>
                            ) : (
                              <EuiText className="nextText">Finish</EuiText>
                            )}
                          </div>
                        )}
                        {schemaData.step > 1 && (
                          <>
                            <ImageButton
                              buttonText={'Back'}
                              ButtonContainerStyle={{
                                width: '120px',
                                height: '42px'
                              }}
                              buttonAction={() => {
                                PreviousStepHandler();
                                setAssigneeName('');
                              }}
                            />
                          </>
                        )}
                      </div>
                    </BottomContainer>
                  </>
                )}
              </>
            ) : (
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
        );
      }}
    </Formik>
  );
}
export default observer(Form);
