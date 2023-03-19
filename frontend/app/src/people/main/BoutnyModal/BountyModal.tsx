import { Modal } from 'components/common';
import React from 'react'
import FocusedView from '../focusView';
import { useStores } from 'store';
import { usePerson } from 'hooks';
import { widgetConfigs } from 'people/utils/constants';
import { useHistory, useParams, useRouteMatch } from 'react-router-dom';

const config = widgetConfigs.wanted;
export const BountyModal = ({basePath}) => {
  const history = useHistory();
  const {path, url} = useRouteMatch();

  console.log(path, url)
  const {wantedId} = useParams<{wantedId: string}>();
  const {ui} = useStores()
  const {canEdit, person} = usePerson(ui.selectedPerson)

  const wantedLength = person?.extras.wanted?.length;

  const changeWanted = (step) => {
    if(!wantedLength) return;
    const currentStep = Number(wantedId);
    const newStep = currentStep + step;
    if (step === 1) {
      if(newStep < wantedLength) {
        history.replace({
          pathname: `${basePath}/${newStep}`
        })
      }
    }
    if(step === -1) {
      if(newStep >= 0) {
        history.replace({
          pathname: `${basePath}/${newStep}`
        })
      }
    }
  }
  const onGoBack = ()=> {
    history.push({
      pathname: basePath
    });
  }

  return (
<Modal
          visible={true}
          style={{
            height: '100%'
          }}
          envStyle={{
            marginTop: 0,
            borderRadius: 0,
            background: '#fff',
            height: '100%',
            width: 'auto',
            minWidth: 500,
            maxWidth: '90%',
            zIndex: 20
          }}
          nextArrow={() =>changeWanted(1)}
          prevArrow={() =>changeWanted(-1)
          }
          overlayClick={() => {
onGoBack()

          }}
          bigClose={() => {
onGoBack()
          }}
>
          <FocusedView
            person={person}
            canEdit={canEdit}
            selectedIndex={Number(wantedId)}
            config={config}
            goBack={() => {
              onGoBack();
            }}
          />
</Modal>
  )
}
