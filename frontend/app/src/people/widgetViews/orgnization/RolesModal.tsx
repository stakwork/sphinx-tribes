import React from 'react';
import { Wrap } from 'components/form/style';
import { useIsMobile } from 'hooks/uiHooks';
import { nonWidgetConfigs } from 'people/utils/Constants';
import { Button, Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import { Check, CheckLabel, CheckLi, CheckUl, ModalTitle } from './style';
import { UserRolesModalProps } from './interface';

const color = colors['light'];

const RolesModal = (props: UserRolesModalProps) => {
    const isMobile = useIsMobile();
    const { isOpen, close, bountyRolesData, roleChange, submitRoles } = props;

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
                <CheckUl>
                    {bountyRolesData.map((role: any, i: number) => (
                        <CheckLi key={i}>
                            <Check
                                checked={role.status}
                                onChange={roleChange}
                                type="checkbox"
                                name={role.name}
                                value={role.name}
                            />
                            <CheckLabel>{role.name}</CheckLabel>
                        </CheckLi>
                    ))}
                </CheckUl>
                <Button
                    onClick={() => submitRoles()}
                    style={{ width: '100%' }}
                    color={'primary'}
                    text={'Add roles'}
                />
            </Wrap>
        </Modal>
    )

};

export default RolesModal;