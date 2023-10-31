import React, { useRef, useEffect, useState, ChangeEvent } from 'react';
import { useDropzone } from 'react-dropzone';
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
import { ModalTitle } from './style';
import { EditOrgModalProps } from './interface';

const color = colors['light'];

const EditOrgColumns = styled.div`
  display: flex;
  flex-direction: row;
`;

const OrgImageOutline = styled.div`
  width: 142px;
  height: 142px; 
  margin-top: 48px;
  margin-bottom: 10px;
  align-self: center;
  cursor: pointer;
  
  border-style: dashed;
  border-width: 2px;
  border-color: #D0D5D8;
  border-radius: 50%;
`;

const OrgImage = styled.img`
  width: 126px;
  height: 126px;
  flex-shrink: 0;
  border-radius: 50%;
  position: relative;
  left: 50%;
  top: 50%;
  transform: translate(-63px, -63px);
`;

const DragAndDrop = styled.input`
  width: 126px;
  height: 126px;
  flex-shrink: 0;
  border-radius: 50%;
  position: relative;
  left: 50%;
  top: 50%;
  transform: translate(-63px, -63px);
`;

const ResetOrgImage = styled.img`
  width: 38.041px;
  height: 38.041px;
  flex-shrink: 0;

  position: relative;
  transform: translate(100px, -30px);

  cursor: pointer;

  :hover {
    filter: brightness(.9);
    transition: 0.2s;
  }
  :active {
    filter: brightness(.7);
  }
`;

const FormWrapper = styled.div`
  margin-left: 30px;
  padding: 0px;
`;

const EditOrgTitle = styled(ModalTitle)`
  color: var(--Text-2, var(--Hover-Icon-Color, #3C3F41));
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 30px;
  font-style: normal;
  font-weight: 800;
  line-height: 30px; /* 100% */
`;

const ImgImportText = styled.p`
  margin-bottom: 5px;
  color: var(--Main-bottom-icons, var(--Disabled-Icon-color, #5F6368));
  text-align: center;
  leading-trim: both;
  text-edge: cap;
  font-family: Roboto;
  font-size: 13px;
  font-style: normal;
  font-weight: 400;
  line-height: 17px; /* 130.769% */
  letter-spacing: 0.13px;
`;

const FileTypeHint = styled.p`
  margin-bottom: 5px;
  color: var(--Placeholder-Text, var(--Disabled-Icon-color, #B0B7BC));
  text-align: center;
  leading-trim: both;
  text-edge: cap;
  font-family: Roboto;
  font-size: 10px;
  font-style: normal;
  font-weight: 400;
  line-height: 18px; /* 180% */
`;

const ImgBrowse = styled.a`
  color: var(--Primary-blue, var(--Disabled-Icon-color, #618AFF));
  leading-trim: both;
  text-edge: cap;
  font-family: Roboto;
  font-size: 13px;
  font-style: normal;
  font-weight: 400;
  line-height: 17px;
  letter-spacing: 0.13px;
  cursor: pointer;
`;

const EditOrgModal = (props: EditOrgModalProps) => {
  const isMobile = useIsMobile();
  const { isOpen, close, onSubmit, onDelete, org } = props;

  const config = widgetConfigs.organizations;
  const schema = [...config.schema];
  const formRef = useRef(null);
  const initValues = {
    name: org?.name,
    image: org?.img,
    show: org?.show
  };

  const [selectedImage, setSelectedImage] = useState<string>(org?.img || '');
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

  const [files, setFiles] = useState<({ preview: string; })[]>([]);
  const {getRootProps, getInputProps} = useDropzone({
    accept: 'image/*',
    onDrop: (acceptedFiles: any) => {
      setFiles(acceptedFiles.map((file: any) => Object.assign(file, {
        preview: URL.createObjectURL(file)
      })));
      setSelectedImage(files[0].preview);
    }
  });

  const openFileDialog = () => {
    if (fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  const resetImg = (e: any) => {
      setSelectedImage(org?.img || '');
      files.forEach((file: any) => URL.revokeObjectURL(file.preview));
      e.stopPropagation();
  }

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
        width: '551px',
        height: '435px',
        padding: '48px',
        flexShrink: '0',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'flex-start'
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
      <EditOrgTitle>Edit Organization</EditOrgTitle>
      <EditOrgColumns>
        <div style={{ display: 'flex', flexDirection: 'column' }}>
          <OrgImageOutline {...getRootProps()}>
            <DragAndDrop {...getInputProps()} />
            <OrgImage src={selectedImage} />
            <ResetOrgImage 
              onClick={resetImg} 
              src={'/static/badges/ResetOrgProfile.svg'} 
            />
          </OrgImageOutline>
          <div>
            <input
              type="file"
              accept="image/*"
              style={{ display: 'none' }}
              onChange={handleFileInputChange}
              ref={fileInputRef}
            />
            <ImgImportText>Drag and drop or <ImgBrowse onClick={openFileDialog}>Browse</ImgBrowse></ImgImportText>
            <FileTypeHint>PNG, JPG or GIF,  Min. 300 x 300 px</FileTypeHint>
          </div>
        </div>
        <FormWrapper>
          <Formik
            initialValues={initValues || {}}
            onSubmit={onSubmit}
            innerRef={formRef}
            validationSchema={validator(schema)}
          >
            {({
              setFieldTouched,
              handleSubmit,
              values,
              setFieldValue,
              errors,
              initialValues
            }: any) => (
              <Wrap newDesign={true}>
                <div className="SchemaInnerContainer">
                  {schema.map((item: FormField) => {
                    const githubdescription =
                      item.name === 'github_description' && !values.ticket_url
                        ? {
                            display: 'none'
                          }
                        : undefined;
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
                        style={{ ...githubdescription }}
                      />
                    );
                  })}
                </div>
                <Button
                  disabled={false}
                  onClick={() => {
                    handleSubmit();
                  }}
                  loading={false}
                  style={{
                    width: '100%',
                    height: '50px',
                    borderRadius: '5px',
                    alignSelf: 'center'
                  }}
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
          width: '200px',
          height: '50px',
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
    </Modal>
  );
};

export default EditOrgModal;
