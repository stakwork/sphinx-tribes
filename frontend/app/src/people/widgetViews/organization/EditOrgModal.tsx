import React from 'react';
import { Wrap } from 'components/form/style';
import { useIsMobile } from 'hooks/uiHooks';
import { nonWidgetConfigs } from 'people/utils/Constants';
import { Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import { ModalTitle } from './style';
import { ModalProps } from './interface';

const color = colors['light'];

const EditOrgModal = (props: ModalProps) => {
  const isMobile = useIsMobile();
  const { isOpen, close } = props;

  const config = nonWidgetConfigs['organizationusers'];

  return (
    <Modal
      visible={isOpen}
      style={{
        height: '100%',
        flexDirection: 'column'
      }}
      envStyle={{
        marginTop: isMobile ? 64 : 0,
        background: color.pureWhite,
        zIndex: 20,
        ...(config?.modalStyle ?? {}),
        maxHeight: '100%',
        borderRadius: '10px'
      }}
      overlayClick={close}
      bigCloseImage={close}
      bigCloseImageStyle={{
        top: '-18px',
        right: '-18px',
        background: '#000',
        borderRadius: '50%'
      }}
    >
      <Wrap newDesign={true}>
        <ModalTitle>Add user roles</ModalTitle>
        
        
      </Wrap>
    </Modal>
  );
};

export default EditOrgModal;
