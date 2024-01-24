import React, { useEffect, useState, useRef } from 'react';
import moment from 'moment';
import { cloneDeep } from 'lodash';
import { observer } from 'mobx-react-lite';
import { FocusViewProps } from 'people/interfaces';
import { EuiGlobalToastList } from '@elastic/eui';
import { Box } from '@mui/system';
import { useStores } from '../../store';
import Form from '../../components/form/bounty';
import {
  Button,
  IconButton,
  useAfterDeleteNotification,
  useDeleteConfirmationModal
} from '../../components/common';
import WantedSummary from '../widgetViews/summaries/WantedSummary';
import { useIsMobile } from '../../hooks';
import { dynamicSchemasByType } from '../../components/form/schema';
import { convertLocaleToNumber, extractRepoAndIssueFromIssueUrl } from '../../helpers';
import { B, BWrap } from './style';

// selected bounty popup window
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
    bounty,
    setRemoveNextAndPrev,
    setAfterEdit
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

  //this code block can be used to get the real org values
  // const [userOrganizations,setUserOrganizations] = useState<any>({});
  // const getOrg =async()=>{
  //   const response = await main.getOrganizations()
  //   const mappedResponse = response.map((org: Organization) => ({
  //           label: org.name,
  //           value: org.uuid
  //         }))
  //   setUserOrganizations(mappedResponse)
  // }
  // getOrg()

  const urlParts = window.location.href.split('/bounties/');
  const uuid = urlParts.length > 1 ? urlParts[1] : null;
  const [orgName, setOrgName] = useState<any>([]);
  const getOrganization = async () => {
    if (!uuid) {
      console.error('No UUID found in the URL');
      return;
    }

    try {
      const response = await main.getUserOrganizationByUuid(uuid);
      console.log(response, 'res');

      if (response && response.uuid && response.name) {
        // Creating an array with an object containing label and value
        const orgDetails = [
          {
            label: response.name,
            value: response.uuid
          }
        ];
        setOrgName(orgDetails);
      } else {
        console.error('Response does not have the expected structure');
      }
    } catch (error) {
      console.error('Error fetching organization:', error);
    }
  };

  useEffect(() => {
    getOrganization();
  });

  function isNotHttps(url: string | undefined) {
    if (main.isTorSave() || url?.startsWith('http://')) {
      return true;
    }
    return false;
  }

  // close bounty popup window
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

  const canDeleteBounty =
    bounty && bounty.length ? !(bounty[0]?.body?.paid || bounty[0]?.body?.assignee.id) : false;

  const { openAfterDeleteNotification } = useAfterDeleteNotification();

  const afterDeleteHandler = (title?: string, link?: string) => {
    openAfterDeleteNotification({
      bountyTitle: title,
      bountyLink: link
    });
  };

  // callback for deleting the open bounty
  async function deleteIt() {
    if (bounty && bounty.length) {
      const delBounty = bounty[0];
      setDeleting(true);
      try {
        if (delBounty.body.created) {
          await main.deleteBounty(delBounty.body.created, delBounty.body.owner_id);
          afterDeleteHandler(delBounty.body.title, delBounty.body.ticket_url);
          closeModal();
          if (props?.deleteExtraFunction) props?.deleteExtraFunction();
        }
      } catch (e) {
        console.log('e', e);
      }
      setDeleting(false);
      if (!isNotHttps(ui?.meInfo?.url) && props.ReCallBounties) props.ReCallBounties();
    }
  }

  const { openDeleteConfirmation } = useDeleteConfirmationModal();

  const deleteHandler = () => {
    openDeleteConfirmation({
      onDelete: deleteIt,
      children: (
        <Box fontSize={20} textAlign="center">
          Are you sure you want to <br />
          <Box component="span" fontWeight="500">
            Delete this Bounty?
          </Box>
        </Box>
      )
    });
  };

  async function preSubmitFunctions(body: any) {
    const newBody = cloneDeep(body);

    // if github repo
    const githubError = "Couldn't locate this Github issue. Make sure this repo is public.";
    try {
      // convert the amount from string to number
      if (newBody.price) {
        newBody.price = convertLocaleToNumber(String(newBody.price));
      }

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

        if (newBody.price) {
          newBody.price = Number(newBody.price);
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

  // eslint-disable-next-line @typescript-eslint/no-inferrable-types
  async function submitForm(body: any, notEdit?: boolean) {
    let newBody = cloneDeep(body);

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
      if (typeof newBody?.assignee !== 'string' || !newBody?.assignee) {
        newBody.assignee = newBody.assignee?.owner_pubkey ?? '';
      }

      if (body.one_sentence_summary !== '') {
        newBody.title = body.one_sentence_summary;
      } else {
        newBody.title = body.title;
      }
      newBody.one_sentence_summary = '';

      if (!newBody.id && !newBody.owner_id) {
        newBody.owner_id = info.pubkey;
      }

      // For editing a bounty, get the pubkey of the bounty creator
      const bounty = await main.getBountyById(Number(newBody.id));
      if (newBody.id && bounty.length) {
        const b = bounty[0];
        newBody.owner_id = b.body.owner_id;
      }

      await main.saveBounty(newBody);

      // Refresh the tickets page if a user eidts from the tickets tab
      if (window.location.href.includes('wanted')) {
        await main.getPersonCreatedBounties({}, info.pubkey);
      }
    } catch (e) {
      console.log('e', e);
    }

    if (props?.onSuccess) props.onSuccess();

    if (notEdit === true) {
      setLoading(false);
    }
    if (ui?.meInfo?.hasOwnProperty('url') && !isNotHttps(ui?.meInfo?.url) && props?.ReCallBounties)
      props?.ReCallBounties();
  }

  // set user organizations
  if (config?.schema?.[0]?.['defaultSchema']?.[0]?.['options']) {
    config.schema[0]['defaultSchema'][0]['options'] = orgName;
  }

  let initialValues: any = {};

  let altInitialValue = {};
  if (window.location.href.includes('/org')) {
    altInitialValue = {
      org_uuid: uuid
    };
  }

  const personInfo = canEdit ? ui.meInfo : person;

  // set initials here
  if (personInfo) {
    // if there is a selected index, fill in values
    if (bounty && bounty.length && selectedIndex >= 0) {
      const selectedBounty = bounty[0];
      const wanted = selectedBounty.body;
      initialValues.estimated_completion_date = wanted?.estimated_completion_date
        ? moment(wanted?.estimated_completion_date)
        : '';

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
    } else {
      initialValues = altInitialValue;
    }
  }

  const noShadow: any = !isMobile ? { boxShadow: '0px 0px 0px rgba(0, 0, 0, 0)' } : {};

  function getExtras(): any {
    if (bounty) {
      const selectedBounty = bounty[0];

      if (selectedIndex >= 0 && selectedBounty && selectedBounty.body) {
        return selectedBounty.body;
      }
    }
    return null;
  }

  function handleEditAction() {
    setEditable(false);
    setEditMode(true);
    setRemoveNextAndPrev && setRemoveNextAndPrev(true);
  }

  function handleEditFinish() {
    setEditable(true);
    setEditMode(false);
    setRemoveNextAndPrev && setRemoveNextAndPrev(false);
    setAfterEdit && setAfterEdit(true);
  }

  function handleFormClose() {
    if (skipEditLayer && goBack) goBack();
    else {
      setEditMode(false);
      setRemoveNextAndPrev && setRemoveNextAndPrev(false);
    }
  }

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
              close={handleFormClose}
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
              onEditSuccess={handleEditFinish}
              setLoading={setLoading}
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
                  color="noColor"
                  onClick={() => {
                    if (goBack) goBack();
                  }}
                  style={{
                    fontSize: 3,
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
                    onClick={deleteHandler}
                    color={'white'}
                    loading={deleting}
                    disabled={!canDeleteBounty}
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
            deleteAction={canDeleteBounty ? deleteHandler : undefined}
            deletingState={deleting}
            editAction={handleEditAction}
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
