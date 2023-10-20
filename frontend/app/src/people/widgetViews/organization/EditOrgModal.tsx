import React, { useRef, useState, ChangeEvent } from 'react';
import { Wrap } from 'components/form/style';
import { useIsMobile } from 'hooks/uiHooks';
import styled from 'styled-components';
import { Formik } from 'formik';
import { validator } from 'components/form/utils';
import { widgetConfigs } from 'people/utils/Constants';
import { FormField } from 'components/form/utils';
import Input from '../../../components/form/inputs';
import { Button, Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import avatarIcon from '../../../public/static/profile_avatar.svg';
import { ModalTitle } from './style';
import { EditOrgModalProps } from './interface';

const color = colors['light'];

const EditOrgColumns = styled.div`
  display: flex;
  flex-direction: row;
`;

const EditOrgName = styled.p`
  color: #8e969c;
  margin: 5px 0px;
`;

const EditOrgWrap = styled(Wrap)`
  width: 100%;
`;

const OrgImage = styled.img`
  width: 150px;
  height: 150px;
  margin-top: 28px;
  border-radius: 50%;
  align-self: center;
`;

const FormWrapper = styled.div`
  margin-left: 30px;
  padding: 0px;
`;

const EditOrgTitle = styled(ModalTitle)`
  font-weight: 800;
  font-size: 30px
`;

const EditOrgModal = (props: EditOrgModalProps) => {
  const isMobile = useIsMobile();
  const { isOpen, close, onSubmit, onDelete, org } = props;

  const config = widgetConfigs.organizations;
  const schema = [...config.schema];
  const formRef = useRef(null);
  const initValues = {
    name: org?.name,
    image: org?.img
  }

  const [selectedImage, setSelectedImage] = useState<string>('');
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const handleFileInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files && e.target.files[0];
    if (file) {
      // Display the selected image
      const imageUrl = URL.createObjectURL(file);
      setSelectedImage(imageUrl);
    } else {
      // Handle the case where the user cancels the file dialog
      setSelectedImage('');
    }
  };

  const openFileDialog = () => {
    if (fileInputRef.current) {
      fileInputRef.current.click();
    }
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
      borderRadius: '10px',
      padding: '20px 60px 10px 60px'
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
      <EditOrgWrap newDesign={true}>
        <EditOrgTitle>Edit Organization</EditOrgTitle>
        <EditOrgColumns>
          <div style={{display:'flex', flexDirection: 'column'}}>
            <OrgImage src={selectedImage}/>
            <div>
              <input
                type="file"
                accept="image/*"
                style={{ display: 'none' }}
                onChange={handleFileInputChange}
                ref={fileInputRef}
              />
              <Button
                disabled={false}
                onClick={() => {
                  openFileDialog();
                }}
                loading={false}
                style={{ width: 'auto', height: '40px', borderRadius: '5px', alignSelf: 'center', margin: '10px' }}
                color={'secondary'}
                text={'Upload org image'}
              />
            </div>
          </div>
          <FormWrapper>
            <Formik
              initialValues={initValues || {}}
              onSubmit={onSubmit}
              innerRef={formRef}
              validationSchema={validator(schema)}
            >
              {({ setFieldTouched, handleSubmit, values, setFieldValue, errors, initialValues }: any) => (
                <Wrap newDesign={true}>
                  <div className="SchemaInnerContainer">
                    {schema.map((item: FormField) => {
                      const githubdescription = item.name === 'github_description' && !values.ticket_url
                        ? {
                            display: 'none'
                          }
                        : undefined
                      return (
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
                            setFieldValue(item.name, e);
                          }}
                          setFieldValue={(e: any, f: any) => {
                            setFieldValue(e, f);
                          }}
                          setFieldTouched={setFieldTouched}
                          handleBlur={() => setFieldTouched(item.name, false)}
                          handleFocus={() => setFieldTouched(item.name, true)}
                          borderType={'bottom'}
                          imageIcon={true}
                          style={
                            {...githubdescription}
                          }
                        />
                    )})}
                  </div>
                  <Button
                    disabled={false}
                    onClick={() => {
                      handleSubmit();
                    }}
                    loading={false}
                    style={{ width: '100%', height: '50px', borderRadius: '5px', alignSelf: 'center' }}
                    color={'primary'}
                    text={'Save changes'}
                  />
                </Wrap>
              )}
            </Formik>
          </FormWrapper>
        </EditOrgColumns>
        <Button
          disabled={false}
          onClick={() => {
            onDelete();
          }}
          loading={false}
          style={{ 
            width: '200px', height: '50px', 
            borderRadius: '5px', 
            borderStyle: 'solid',
            alignSelf: 'center',
            borderWidth: '2px',
            backgroundColor: 'white',
            borderColor: '#ED7474',
            color: '#ED7474'
          }}
          color={'#ED7474'}
          text={'Delete organization'}
        />
      </EditOrgWrap>
    </Modal>
  );
};

export default EditOrgModal;