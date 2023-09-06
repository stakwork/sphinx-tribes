import React, { useEffect, useState, useRef } from 'react';
import styled, { css } from 'styled-components';
import moment from 'moment';
import { cloneDeep } from 'lodash';
import { observer } from 'mobx-react-lite';
import { FocusViewProps } from 'people/interfaces';
import { EuiGlobalToastList } from '@elastic/eui';
import { Organization } from 'store/main';
import { useStores } from '../../store';
import Form from '../../components/form';
import { Button, IconButton } from '../../components/common';
import WantedSummary from '../widgetViews/summaries/WantedSummary';
import { useIsMobile } from '../../hooks';
import { dynamicSchemasByType } from '../../components/form/schema';
import { extractRepoAndIssueFromIssueUrl, toCapitalize } from '../../helpers';

// this is where we see others posts (etc) and edit our own
const BWrap = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding: 10px;
  min-height: 42px;
  position: absolute;
  left: 0px;
  border-bottom: 1px solid rgb(221, 225, 229);
  background: #ffffff;
  box-shadow: 0px 1px 6px rgba(0, 0, 0, 0.07);
  z-index: 100;
`;

const EnvWithScrollBar = ({ thumbColor, trackBackgroundColor }: any) => css`
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
  display: ${(p: any) => (p.hide ? 'none' : 'flex')};
  justify-content: ${(p: any) => (p.hide ? 'none' : 'center')};
  height: 100%;
  width: 100%;
  overflow-y: auto;
  box-sizing: border-box;
  ${EnvWithScrollBar({
    thumbColor: '#5a606c',
    trackBackgroundColor: 'rgba(0,0,0,0)'
  })}
`;
function FocusedView(props: FocusViewProps) {
  const {
    goBack,
    config,
    selectedIndex,
    canEdit,
    person,
    buttonsOnBottom,
    formHeader,
    manualGoBackOnly,
    isFirstTimeScreen,
    fromBountyPage,
    newDesign,
    setIsModalSideButton,
    bounty
  } = props;
  const { ui, main } = useStores();

  const skipEditLayer = selectedIndex < 0 || config.skipEditLayer ? true : false;

  const [loading, setLoading] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [editMode, setEditMode] = useState(skipEditLayer);
  const [editable, setEditable] = useState<boolean>(!canEdit);
  const [toasts, setToasts]: any = useState([]);

  const scrollDiv: any = useRef(null);
  const formRef: any = useRef(null);

  const isMobile = useIsMobile();

  const isTorSave = canEdit && main.isTorSave();

  const userOrganizations = main.organizations.map((org: Organization) => ({
    label: toCapitalize(org.name),
    value: org.uuid
  }));

  function isNotHttps(url: string | undefined) {
    if (main.isTorSave() || url?.startsWith('http://')) {
      return true;
    }
    return false;
  }

  function closeModal() {
    if (!manualGoBackOnly) {
      ui.setEditMe(false);
      if (props.goBack) props.goBack();
    }
  }

  const addToast = () => {
    setToasts([
      {
        id: '1',
        title: 'Add a description to your bounty'
      }
    ]);
  };

  const removeToast = () => {
    setToasts([]);
  };

  // get self on unmount if tor user
  useEffect(
    () =>
      function cleanup() {
        if (isTorSave) {
          main.getSelf(null);
        }
      },
    [main, isTorSave]
  );

  async function deleteIt() {
    const delBounty = bounty && bounty.length ? bounty[0] : main.peopleWanteds[selectedIndex];
    if (!delBounty) return;
    setDeleting(true);
    try {
      if (delBounty.body.created) {
        await main.deleteBounty(delBounty.body.created, delBounty.body.owner_id);
        closeModal();
        if (props?.deleteExtraFunction) props?.deleteExtraFunction();
      }
    } catch (e) {
      console.log('e', e);
    }
    setDeleting(false);
    if (!isNotHttps(ui?.meInfo?.url) && props.ReCallBounties) props.ReCallBounties();
  }

  async function preSubmitFunctions(body: any) {
    const newBody = cloneDeep(body);

    // if github repo
    const githubError = "Couldn't locate this Github issue. Make sure this repo is public.";
    try {
      if (
        newBody.ticket_url &&
        (newBody.type === 'wanted_coding_task' ||
          newBody.type === 'coding_task' ||
          newBody.type === 'freelance_job_request')
      ) {
        const { repo, issue } = extractRepoAndIssueFromIssueUrl(newBody.ticket_url);
        const splitString = repo.split('/');
        const [ownerName, repoName] = splitString;
        const res = await main.getGithubIssueData(ownerName, repoName, `${issue}`);

        if (!res) {
          throw githubError;
        }

        const { description } = res;

        if (newBody.github_description) {
          newBody.description = description;
        }

        // body.description = description;
        newBody.title = newBody.one_sentence_summary;

        // save repo to cookies for autofill in form
        ui.setLastGithubRepo(newBody.ticket_url);
      }
    } catch (e) {
      throw githubError;
    }

    return newBody;
  }

  async function submitForm(body: any) {
    let newBody = cloneDeep(body);
    delete newBody.assignee;
    try {
      newBody = await preSubmitFunctions(newBody);
    } catch (e) {
      console.log('e', e);
      alert(e);
      return;
    }

    if (!newBody) return; // avoid saving bad state
    if (!newBody.description) {
      addToast();
    }
    const info = ui.meInfo as any;
    if (!info) return console.log('no meInfo');
    setLoading(true);
    try {
      if (body?.assignee?.owner_pubkey) {
        newBody.assignee = body.assignee.owner_pubkey;
      }
      if (body.one_sentence_summary !== '') {
        newBody.title = body.one_sentence_summary;
      } else {
        newBody.title = body.title;
      }
      newBody.one_sentence_summary = '';
      newBody.owner_id = info.pubkey;

      await main.saveBounty(newBody);
      // Refresh the tickets page if a user eidts from the tickets tab
      if (window.location.href.includes('wanted')) {
        await main.getPersonCreatedWanteds({}, info.pubkey);
      }
      closeModal();
    } catch (e) {
      console.log('e', e);
    }
    if (props?.onSuccess) props.onSuccess();
    setLoading(false);
    if (ui?.meInfo?.hasOwnProperty('url') && !isNotHttps(ui?.meInfo?.url) && props?.ReCallBounties)
      props?.ReCallBounties();
  }

  let initialValues: any = {};

  const personInfo = canEdit ? ui.meInfo : person;
  const selectedBounty = bounty && bounty.length ? bounty[0] : main.peopleWanteds[selectedIndex];

  // set initials here
  if (personInfo && selectedBounty && selectedIndex >= 0) {
    const wanted = selectedBounty.body;
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
        wanted?.map((value: any) => moment(value?.estimated_completion_date)) || '';
      // below are extras,
      initialValues.twitter =
        (personInfo.extras?.twitter && personInfo.extras?.twitter[0]?.value) || '';
      initialValues.email = (personInfo.extras?.email && personInfo.extras?.email[0]?.value) || '';
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
      if (selectedIndex >= 0) {
        if (wanted.type) {
          const thisDynamicSchema = dynamicSchemasByType[wanted.type];
          const newValues = thisDynamicSchema.map((s: any) => {
            if (s.name === 'estimated_completion_date') {
              return {
                [s.name]: wanted['estimated_completion_date'] || new Date()
              };
            } else if (s.name === 'one_sentence_summary') {
              return {
                [s.name]: wanted['one_sentence_summary'] || wanted['title']
              };
            } else if (s.name === 'coding_languages') {
              const coding_languages =
                wanted['coding_languages'] && wanted['coding_languages'].length
                  ? wanted['coding_languages'].map((lang: any) => ({ value: lang, label: lang }))
                  : [];
              return {
                [s.name]: coding_languages
              };
            }
            return {
              [s.name]: wanted[s.name]
            };
          });

          const valueMap = Object.assign({}, ...newValues);
          initialValues = { ...initialValues, ...valueMap };
        } else {
          const dynamicSchema = config?.schema?.find((f: any) => f.defaultSchema);
          dynamicSchema?.defaultSchema?.forEach((s: any) => {
            initialValues[s.name] = wanted[s.name];
          });
        }
      }
    }
  }

  const noShadow: any = !isMobile ? { boxShadow: '0px 0px 0px rgba(0, 0, 0, 0)' } : {};

  function getExtras(): any {
    const selectedBounty = bounty && bounty.length ? bounty[0] : main.peopleWanteds[selectedIndex];
    if (selectedIndex >= 0) {
      if (selectedBounty) {
        return selectedBounty.body;
      } else {
        return null;
      }
    }
    return null;
  }

  // set user organizations
  config.schema[0]['defaultSchema'][0]['options'] = userOrganizations;

  return (
    <div
      style={{
        ...props?.style,
        width: '100%',
        height: '100%'
      }}
    >
      {editMode ? (
        <B ref={scrollDiv} hide={false}>
          {formHeader && formHeader}
          {ui.meInfo && (
            <Form
              newDesign={newDesign}
              buttonsOnBottom={buttonsOnBottom}
              isFirstTimeScreen={isFirstTimeScreen}
              readOnly={editable}
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
          <WantedSummary
            {...getExtras()}
            ReCallBounties={props?.ReCallBounties}
            formSubmit={submitForm}
            person={person}
            personBody={props?.personBody}
            item={getExtras()}
            config={config}
            fromBountyPage={fromBountyPage}
            extraModalFunction={props?.extraModalFunction}
            deleteAction={deleteIt}
            deletingState={deleting}
            editAction={() => {
              setEditable(false);
              setEditMode(true);
            }}
            setIsModalSideButton={setIsModalSideButton}
            setIsExtraStyle={props?.setIsExtraStyle}
          />
        </>
      )}
      <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={6000} />
    </div>
  );
}

export default observer(FocusedView);
