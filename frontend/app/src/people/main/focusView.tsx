import React, { useEffect, useState, useRef } from 'react';
import { useStores } from '../../store';
import { useObserver } from 'mobx-react-lite';
import Form from '../../form';
import styled, { css } from 'styled-components';
import { Button, IconButton } from '../../sphinxUI';
import moment from 'moment';
import SummaryViewer from '../widgetViews/summaryViewer';
import { useIsMobile } from '../../hooks';
import { dynamicSchemasByType } from '../../form/schema';
import { extractRepoAndIssueFromIssueUrl } from '../../helpers';

// this is where we see others posts (etc) and edit our own
export default function FocusedView(props: any) {
  const {
    onSuccess,
    goBack,
    config,
    selectedIndex,
    canEdit,
    person,
    buttonsOnBottom,
    formHeader,
    manualGoBackOnly,
    isFirstTimeScreen
  } = props;
  const { ui, main } = useStores();
  const { ownerTribes } = main;

  const skipEditLayer = selectedIndex < 0 || config.skipEditLayer ? true : false;

  const [loading, setLoading] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [editMode, setEditMode] = useState(skipEditLayer);

  const scrollDiv: any = useRef(null);
  const formRef: any = useRef(null);

  const isMobile = useIsMobile();

  // console.log({ ...person });

  const torSave = canEdit && ui?.meInfo?.url?.includes('.onion');

  function closeModal(override) {
    if (!manualGoBackOnly) {
      console.log('close modal');
      ui.setEditMe(false);
      if (props.goBack) props.goBack();
    }
  }

  // get self on unmount if tor user
  useEffect(() => {
    return function cleanup() {
      if (torSave) {
        main.getSelf(null);
      }
    };
  }, []);

  function mergeFormWithMeData(v) {
    let fullMeData: any = null;

    if (ui.meInfo) {
      fullMeData = { ...ui.meInfo };

      // add extras if doesnt exist, for brand new users
      if (!fullMeData.extras) fullMeData.extras = {};
      // if about
      if (config.name === 'about') {
        config?.schema?.forEach((s) => {
          if (s.widget && fullMeData.extras) {
            // this allows the link widgets to be edited as a part of about me,
            // when really they are stored as extras

            // include full tribe info from ownerTribes data
            if (s.name === 'tribes') {
              const submitTribes: any = [];

              v[s.name] &&
                v[s.name].forEach((t) => {
                  const fullTribeInfo =
                    ownerTribes && ownerTribes?.find((f) => f.unique_name === t.value);

                  // disclude sensitive details
                  if (fullTribeInfo)
                    submitTribes.push({
                      name: fullTribeInfo.name,
                      unique_name: fullTribeInfo.unique_name,
                      img: fullTribeInfo.img,
                      description: fullTribeInfo.description,
                      ...t
                    });
                });

              fullMeData.extras[s.name] = submitTribes;
            } else if (s.name === 'repos' || s.name === 'coding_languages') {
              // multiples, so we don't need a wrapper
              fullMeData.extras[s.name] = v[s.name];
            } else {
              fullMeData.extras[s.name] = [{ value: v[s.name] }];
            }
          } else {
            fullMeData[s.name] = v[s.name];
          }
        });
      }
      // if extras
      else {
        // add timestamp if not there
        if (!v.created) v.created = moment().unix();

        if (!fullMeData.extras) fullMeData.extras = {};
        // if editing widget
        if (selectedIndex > -1) {
          // mutate it
          fullMeData.extras[config.name][selectedIndex] = v;
        } else {
          // if creating new widget
          if (fullMeData.extras[config.name]) {
            //if not first of its kind
            fullMeData.extras[config.name].unshift(v);
          } else {
            //if first of its kind
            fullMeData.extras[config.name] = [v];
          }
        }
      }
    }
    return fullMeData;
  }

  async function deleteIt() {
    let body: any = null;
    body = { ...ui.meInfo };

    // mutates
    body.extras[config.name].splice(selectedIndex, 1);

    const info = ui.meInfo as any;
    if (!info) return console.log('no meInfo');

    setDeleting(true);
    try {
      await main.saveProfile(body);
      await main.getPeople();
      closeModal(true);
    } catch (e) {
      console.log('e', e);
    }
    setDeleting(false);
  }

  async function preSubmitFunctions(body) {
    // if github repo

    const githubError = "Couldn't locate this Github issue. Make sure this repo is public.";
    try {
      if (
        body.type === 'wanted_coding_task' ||
        body.type === 'coding_task' ||
        body.type === 'freelance_job_request'
      ) {
        const { repo, issue } = extractRepoAndIssueFromIssueUrl(body.ticketUrl);
        const splitString = repo.split('/');
        const ownerName = splitString[0];
        const repoName = splitString[1];
        const res = await main.getGithubIssueData(ownerName, repoName, `${issue}`);

        if (!res) {
          throw githubError;
        }
        const { description, title } = res;
        console.log(description, { ...body });
        if (body.github_description) {
          body.description = description;
        }
        // body.description = description;
        body.title = body.one_sentence_summary;

        // save repo to cookies for autofill in form
        ui.setLastGithubRepo(body.ticketUrl);
      }
    } catch (e) {
      throw githubError;
    }

    return body;
  }

  async function submitForm(body) {
    console.log('START SUBMIT FORM', body);

    try {
      body = await preSubmitFunctions(body);
    } catch (e) {
      console.log('e', e);
      alert(e);
      return;
    }

    body = mergeFormWithMeData(body);

    if (!body) return; // avoid saving bad state
    const info = ui.meInfo as any;
    if (!info) return console.log('no meInfo');

    const date = new Date();
    const unixTimestamp = Math.floor(date.getTime() / 1000);
    setLoading(true);
    try {
      const newBody = {
        ...body,
        alert: undefined,
        new_ticket_time: unixTimestamp,
        extras: {
          ...body?.extras,
          alert: body.alert
        }
      };

      await main.saveProfile(
        config.name === 'about' || config.name === 'wanted' ? { ...newBody } : body
      );
      closeModal(true);
    } catch (e) {
      console.log('e', e);
    }
    setLoading(false);
  }

  return useObserver(() => {
    // let initialValues: MeData = emptyMeInfo;
    const initialValues: any = {};

    const personInfo = canEdit ? ui.meInfo : person;

    // // console.log({...personInfo}.extras.wanted.map((value) => value.estimated_completion_date));
    // personInfo?.extras?.wanted?.map((value: any) => {
    //   console.log(typeof new Date({ ...value }?.estimated_completion_date));
    // });
    // set initials here
    if (personInfo) {
      if (config && config.name === 'about') {
        initialValues.id = personInfo.id || 0;
        initialValues.pubkey = personInfo.pubkey;
        initialValues.alert = personInfo.extras?.alert || false;
        initialValues.owner_alias = personInfo.owner_alias || '';
        initialValues.img = personInfo.img || '';
        initialValues.price_to_meet = personInfo.price_to_meet || 0;
        initialValues.description = personInfo.description || '';
        initialValues.loomEmbedUrl = personInfo.loomEmbedUrl || '';
        initialValues.estimated_completion_date =
          personInfo.extras?.wanted?.map((value) => {
            return moment(value?.estimated_completion_date);
          }) || '';
        // below are extras,
        initialValues.twitter =
          (personInfo.extras?.twitter && personInfo.extras?.twitter[0]?.value) || '';
        initialValues.github =
          (personInfo.extras?.github && personInfo.extras?.github[0]?.value) || '';
        initialValues.facebook =
          (personInfo.extras?.facebook && personInfo.extras?.facebook[0]?.value) || '';
        // extras with multiple items
        initialValues.coding_languages = personInfo.extras?.coding_languages || [];
        initialValues.tribes = personInfo.extras?.tribes || [];
        initialValues.repos = personInfo.extras?.repos || [];
        initialValues.lightning =
          (personInfo.extras?.lightning && personInfo.extras?.lightning[0]?.value) || '';
        initialValues.amboss =
          (personInfo.extras?.amboss && personInfo.extras?.amboss[0]?.value) || '';
      } else {
        // if there is a selected index, fill in values
        if (selectedIndex > -1) {
          const extras = { ...personInfo.extras };
          const sel =
            extras[config.name] &&
            extras[config.name].length > selectedIndex - 1 &&
            extras[config.name][selectedIndex];

          if (sel) {
            // if dynamic, find right schema
            const dynamicSchema = config?.schema?.find((f) => f.defaultSchema);
            if (dynamicSchema) {
              if (sel.type) {
                const thisDynamicSchema = dynamicSchemasByType[sel.type];
                thisDynamicSchema?.forEach((s) => {
                  initialValues[s.name] = sel[s.name];
                });
              } else {
                // use default schema
                dynamicSchema?.defaultSchema?.forEach((s) => {
                  initialValues[s.name] = sel[s.name];
                });
              }
            } else {
              config?.schema?.forEach((s) => {
                initialValues[s.name] = sel[s.name];
              });
            }
          }
        }
      }
    }

    const noShadow: any = !isMobile ? { boxShadow: '0px 0px 0px rgba(0, 0, 0, 0)' } : {};

    return (
      <div
        style={{
          ...props.style,
          width: '100%',
          height: '100%'
        }}
      >
        {editMode ? (
          <B ref={scrollDiv} hide={false}>
            {formHeader && formHeader}
            {ui.meInfo && (
              <Form
                buttonsOnBottom={buttonsOnBottom}
                isFirstTimeScreen={isFirstTimeScreen}
                readOnly={!canEdit}
                formRef={formRef}
                submitText={config && config.submitText}
                loading={loading}
                close={() => {
                  if (skipEditLayer && goBack) goBack();
                  else setEditMode(false);
                }}
                onSubmit={submitForm}
                scrollDiv={scrollDiv}
                schema={config && config.schema}
                initialValues={initialValues}
                extraHTML={
                  ui.meInfo.verification_signature
                    ? {
                        twitter: `<span>Post this to your twitter account to verify:</span><br/><strong>Sphinx Verification: ${ui.meInfo.verification_signature}</strong>`
                      }
                    : {}
                }
              />
            )}
          </B>
        ) : (
          <>
            {(isMobile || canEdit) && (
              <BWrap
                style={{
                  ...noShadow
                }}
              >
                {goBack ? (
                  <IconButton
                    icon="arrow_back"
                    onClick={() => {
                      if (goBack) goBack();
                    }}
                    style={{
                      fontSize: 12,
                      fontWeight: 600
                    }}
                  />
                ) : (
                  <div />
                )}
                {canEdit ? (
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'center',
                      alignItems: 'center'
                    }}
                  >
                    <Button
                      onClick={() => setEditMode(true)}
                      color={'widget'}
                      leadingIcon={'edit'}
                      iconSize={18}
                      width={100}
                      text={'Edit'}
                    />
                    <Button
                      onClick={() => deleteIt()}
                      color={'white'}
                      loading={deleting}
                      leadingIcon={'delete_outline'}
                      text={'Delete'}
                      style={{
                        marginLeft: 10
                      }}
                    />
                  </div>
                ) : (
                  <div />
                )}
              </BWrap>
            )}

            {(isMobile || canEdit) && <div style={{ height: 60 }} />}

            {/* display item */}
            <SummaryViewer
              person={person}
              item={person?.extras && person.extras[config?.name][selectedIndex]}
              config={config}
            />
          </>
        )}
      </div>
    );
  });
}

const BWrap = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding: 10px;
  min-height: 42px;
  position: absolute;
  left: 0px;
  background: #ffffff;
  box-shadow: 0px 1px 6px rgba(0, 0, 0, 0.07);
  z-index: 100;
`;

const EnvWithScrollBar = ({ thumbColor, trackBackgroundColor }) => css`
                scrollbar-color: ${thumbColor} ${trackBackgroundColor}; // Firefox support
                scrollbar-width: thin;

                &::-webkit-scrollbar {
                    width: 6px;
                height: 100%;
  }

                &::-webkit-scrollbar-thumb {
                    background - color: ${thumbColor};
                background-clip: content-box;
                border-radius: 5px;
                border: 1px solid ${trackBackgroundColor};
  }

                &::-webkit-scrollbar-corner,
                &::-webkit-scrollbar-track {
                    background - color: ${trackBackgroundColor};
  }
}

                `;
interface BProps {
  hide: boolean;
}
const B = styled.div<BProps>`
  display: ${(p) => p.hide && 'none'};
  height: 100%;
  width: 100%;
  overflow-y: auto;
  box-sizing: border-box;
  ${EnvWithScrollBar({
    thumbColor: '#5a606c',
    trackBackgroundColor: 'rgba(0,0,0,0)'
  })}
`;
