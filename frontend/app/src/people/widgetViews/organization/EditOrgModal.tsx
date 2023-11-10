import React, { useRef, useState, ChangeEvent } from 'react';
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
import DeleteOrgWindow from './DeleteOrgWindow';

const color = colors['light'];

const EditOrgColumns = styled.div`
  display: flex;
  flex-direction: row;
  width: 100%;
`;

const OrgEditImageWrapper = styled.div`
  display: flex;
  flex-direction: column;
  flex: 40%;
`;

const OrgImageOutline = styled.div`
  width: 142px;
  height: 142px;
  margin-top: 28px;
  margin-bottom: 10px;
  align-self: center;
  cursor: pointer;

  border-style: dashed;
  border-width: 2px;
  border-color: #d0d5d8;
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
    filter: brightness(0.9);
    transition: 0.2s;
  }
  :active {
    filter: brightness(0.7);
  }
`;

const FormWrapper = styled.div`
  flex: 60%;
`;

const EditOrgTitle = styled(ModalTitle)`
  color: var(--Text-2, var(--Hover-Icon-Color, #3c3f41));
  font-family: 'Barlow';
  font-size: 30px;
  font-style: normal;
  font-weight: 800;
  line-height: 30px; /* 100% */
`;

const ImgImportText = styled.p`
  margin-bottom: 5px;
  color: var(--Main-bottom-icons, var(--Disabled-Icon-color, #5f6368));
  text-align: center;
  font-family: 'Roboto';
  font-size: 13px;
  font-style: normal;
  font-weight: 400;
  line-height: 17px; /* 130.769% */
  letter-spacing: 0.13px;
`;

const FileTypeHint = styled.p`
  margin-bottom: 5px;
  color: var(--Placeholder-Text, var(--Disabled-Icon-color, #b0b7bc));
  text-align: center;
  font-family: 'Roboto';
  font-size: 10px;
  font-style: normal;
  font-weight: 400;
  line-height: 18px; /* 180% */
`;

const ImgBrowse = styled.a`
  color: var(--Primary-blue, var(--Disabled-Icon-color, #618aff));
  font-family: 'Roboto';
  font-size: 13px;
  font-style: normal;
  font-weight: 400;
  line-height: 17px;
  letter-spacing: 0.13px;
  cursor: pointer;
`;

const HLine = styled.div`
  background-color: #ebedef;
  height: 1px;
  width: 100%;
  margin: 5px 0px 20px;
`;

const EditOrgModal = (props: EditOrgModalProps) => {
  const isMobile = useIsMobile();
  const { isOpen, close, onSubmit, onDelete, org } = props;
  const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);

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

  const [files, setFiles] = useState<{ preview: string }[]>([]);
  const { getRootProps, getInputProps } = useDropzone({
    onDrop: (acceptedFiles: any) => {
      setFiles(
        acceptedFiles.map((file: any) =>
          Object.assign(file, {
            preview: URL.createObjectURL(file)
          })
        )
      );
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
  };

  return (
    <>
      <Modal
        visible={isOpen}
        style={{
          height: '100%',
          flexDirection: 'column'
        }}
        envStyle={{
          marginTop: isMobile ? 64 : 0,
          background: color.pureWhite,
          zIndex: 0,
          ...(config?.modalStyle ?? {}),
          maxHeight: '100%',
          borderRadius: '10px',
          width: isMobile ? '100%' : '551px',
          minWidth: '20px',
          height: isMobile ? 'auto' : '435px',
          padding: '48px',
          flexShrink: '0',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'flex-start'
        }}
        overlayClick={close}
        bigCloseImage={close}
        bigCloseImageStyle={{
          top: isMobile ? '26px' : '-18px',
          right: isMobile ? '26px' : '-18px',
          background: '#000',
          borderRadius: '50%'
        }}
      >
        <EditOrgTitle>Edit Organization</EditOrgTitle>
        <EditOrgColumns style={{ flexDirection: isMobile ? 'column' : 'row' }}>
          <OrgEditImageWrapper>
            <OrgImageOutline {...getRootProps()}>
              <DragAndDrop type="file" accept="image/*" {...getInputProps()} />
              <OrgImage src={selectedImage} />
              <ResetOrgImage onClick={resetImg} src={'/static/badges/ResetOrgProfile.svg'} />
            </OrgImageOutline>
            <div>
              <input
                type="file"
                accept="image/*"
                style={{ display: 'none' }}
                onChange={handleFileInputChange}
                ref={fileInputRef}
              />
              <ImgImportText>
                Drag and drop or <ImgBrowse onClick={openFileDialog}>Browse</ImgBrowse>
              </ImgImportText>
              <FileTypeHint>PNG, JPG or GIF, Min. 300 x 300 px</FileTypeHint>
            </div>
          </OrgEditImageWrapper>
          <FormWrapper style={{ marginLeft: isMobile ? '0px' : '28px' }}>
            <Formik
              initialValues={initValues || {}}
              onSubmit={onSubmit}
              innerRef={formRef}
              validationSchema={validator(schema)}
              style={{ width: '100%' }}
            >
              {({
                setFieldTouched,
                handleSubmit,
                values,
                setFieldValue,
                errors,
                initialValues
              }: any) => (
                <Wrap style={{ width: '100%' }} newDesign={true}>
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
                        style={{ width: '100%' }}
                      />
                    ))}
                  </div>
                  <Button
                    disabled={false}
                    onClick={() => handleSubmit()}
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
        <HLine style={{ width: '551px', transform: 'translate(-48px, 0px' }} />
        <Button
          disabled={false}
          onClick={() => {
            setShowDeleteModal(true);
          }}
          loading={false}
          style={{
            width: isMobile ? '100%' : 'calc(60% - 18px)',
            height: '50px',
            borderRadius: '5px',
            borderStyle: 'solid',
            alignSelf: 'flex-end',
            borderWidth: '2px',
            backgroundColor: 'white',
            borderColor: '#ED7474',
            color: '#ED7474'
          }}
          color={'#ED7474'}
          text={'Delete organization'}
        />
        {showDeleteModal ? (
          <DeleteOrgWindow onDeleteOrg={onDelete} close={() => setShowDeleteModal(false)} />
        ) : (
          <></>
        )}
      </Modal>
    </>
  );
};

export default EditOrgModal;
