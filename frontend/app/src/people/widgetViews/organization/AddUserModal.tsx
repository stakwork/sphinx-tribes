import React, { useRef, useState } from 'react';
import { Wrap } from 'components/form/style';
import { useIsMobile } from 'hooks/uiHooks';
import { nonWidgetConfigs } from 'people/utils/Constants';
import { Formik } from 'formik';
import { FormField, validator } from 'components/form/utils';
import { spliceOutPubkey } from 'helpers';
import { Button, Modal } from '../../../components/common';
import Input from '../../../components/form/inputs';
import { colors } from '../../../config/colors';
import { ModalTitle, RouteHintText } from './style';
import { AddUserModalProps } from './interface';

const color = colors['light'];

const AddUserModal = (props: AddUserModalProps) => {
  const isMobile = useIsMobile();
  const { isOpen, close, onSubmit, loading, disableFormButtons, setDisableFormButtons } = props;
  const [displayHint, setDisplayHint] = useState(false);

  const hintText = 'Route hint detected and removed';

  const checkDisplayHint = (address: string) => {
    if (address.includes(':')) {
      setDisplayHint(true);
    } else {
      setDisplayHint(false);
    }
  };

  const config = nonWidgetConfigs['organizationusers'];

  const schema = [...config.schema];

  const formRef = useRef(null);

  const initValues = {
    owner_pubkey: ''
  };

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
      <Formik
        initialValues={initValues || {}}
        onSubmit={onSubmit}
        innerRef={formRef}
        validationSchema={validator(schema)}
      >
        {({ setFieldTouched, handleSubmit, values, setFieldValue, errors, initialValues }: any) => (
          <Wrap newDesign={true}>
            <ModalTitle>Add new user</ModalTitle>
            {displayHint && <RouteHintText>{hintText}</RouteHintText>}
            <div className="SchemaInnerContainer">
              {schema.map((item: FormField) => (
                <Input
                  {...item}
                  key={item.name}
                  values={values}
                  errors={errors}
                  value={values[item.name]}
                  error={errors[item.name]}
                  initialValues={initialValues}
                  deleteErrors={() => {
                    if (errors[item.name]) delete errors[item.name];
                  }}
                  handleChange={(e: any) => {
                    checkDisplayHint(e);
                    const pubkey = spliceOutPubkey(e);
                    setFieldValue(item.name, pubkey);
                  }}
                  setFieldValue={(e: any, f: any) => {
                    setFieldValue(e, f);
                  }}
                  setFieldTouched={setFieldTouched}
                  handleBlur={() => setFieldTouched(item.name, false)}
                  handleFocus={() => setFieldTouched(item.name, true)}
                  setDisableFormButtons={setDisableFormButtons}
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
              <Button
                disabled={disableFormButtons || loading}
                onClick={() => {
                  handleSubmit();
                }}
                loading={loading}
                style={{ width: '100%' }}
                color={'primary'}
                text={'Add user'}
              />
            </div>
          </Wrap>
        )}
      </Formik>
    </Modal>
  );
};

export default AddUserModal;
