import React, { useEffect, useState, useRef } from 'react';
import { cloneDeep } from 'lodash';
import { observer } from 'mobx-react-lite';
import { FocusViewProps } from 'people/interfaces';
import { useStores } from '../../store';
import Form from '../../components/form/about';

import { B } from './style';

const AboutFocusView = (props: FocusViewProps) => {
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
    newDesign,
    setRemoveNextAndPrev
  } = props;
  const { ui, main } = useStores();

  const skipEditLayer = selectedIndex < 0 || config.skipEditLayer ? true : false;

  const [loading] = useState(false);
  const [editMode, setEditMode] = useState(skipEditLayer);
  const [editable] = useState<boolean>(!canEdit);

  const scrollDiv: any = useRef(null);
  const formRef: any = useRef(null);

  const isTorSave = canEdit && main.isTorSave();

  // close bounty popup window
  function closeModal() {
    if (!manualGoBackOnly) {
      ui.setEditMe(false);
      if (props.goBack) props.goBack();
    }
  }

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

  function formatAboutBody(newBody: any): any {
    const requestData = {
      ...newBody,
      extras: {
        ...newBody?.extras,
        alert: newBody.alert,
        tribes: newBody.tribes,
        coding_languages: newBody.coding_languages,
        lightning: [
          {
            value: newBody.lightning,
            label: newBody.lightning
          }
        ],
        amboss: [
          {
            value: newBody.amboss,
            label: newBody.amboss
          }
        ],
        email: [
          {
            value: newBody.email,
            label: newBody.email
          }
        ],
        facebook: [
          {
            value: newBody.facebook,
            label: newBody.facebook
          }
        ],
        twitter: [
          {
            value: newBody.twitter,
            label: newBody.twitter
          }
        ],
        github: [
          {
            value: newBody.github,
            label: newBody.github
          }
        ]
      }
    };
    return requestData;
  }

  // eslint-disable-next-line @typescript-eslint/no-inferrable-types
  async function submitForm(body: any, shouldCloseModal: boolean = true) {
    const newBody = cloneDeep(body);

    if (config && config.name === 'about') {
      const requestData = formatAboutBody(newBody);
      await main.saveProfile(requestData);
      if (shouldCloseModal) {
        closeModal();
      }
      return;
    }
  }

  const initialValues: any = {};

  const personInfo = canEdit ? ui.meInfo : person;

  // set initials here
  if (personInfo) {
    if (config && config.name === 'about') {
      initialValues.id = personInfo.id || 0;
      initialValues.pubkey = personInfo.pubkey;
      initialValues.owner_pubkey = personInfo.pubkey;
      initialValues.alert = personInfo.extras?.alert || false;
      initialValues.owner_alias = personInfo.owner_alias || '';
      initialValues.img = personInfo.img || '';
      initialValues.price_to_meet = personInfo.price_to_meet || 0;
      initialValues.description = personInfo.description || '';
      initialValues.loomEmbedUrl = personInfo.loomEmbedUrl || '';
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
      initialValues.lightning =
        (personInfo.extras?.lightning && personInfo.extras?.lightning[0]?.value) || '';
      initialValues.amboss =
        (personInfo.extras?.amboss && personInfo.extras?.amboss[0]?.value) || '';
    }
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
      {editMode && (
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
            />
          )}
        </B>
      )}
    </div>
  );
};

export default observer(AboutFocusView);
