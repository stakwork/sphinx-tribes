import React, { useRef, useState, ChangeEvent, DragEvent } from 'react';
import { useStores } from 'store';
import { useIsMobile } from 'hooks/uiHooks';
import styled from 'styled-components';
import { Formik } from 'formik';
import { validator } from 'components/form/utils';
import { widgetConfigs } from 'people/utils/Constants';
import { FormField } from 'components/form/utils';
import Input from '../../../components/form/inputs';
import { Button, Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import {
  ImgContainer,
  ImgDashContainer,
  ImgDetailInfo,
  ImgInstructionSpan,
  ImgInstructionText,
  ImgText,
  ImgTextContainer,
  InputFile,
  ModalTitle,
  OrgInputContainer,
  SelectedImg,
  UploadImageContainer
} from './style';
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

const FormWrapper = styled.div`
  flex: 60%;
`;

const EditOrgTitle = styled(ModalTitle)`
  color: #3c3f41;
  font-family: 'Barlow';
  font-size: 1.875rem;
  font-style: normal;
  font-weight: 800;
  line-height: 1.875rem;
  margin-bottom: 3rem;

  @media only screen and (max-width: 500px) {
    text-align: center;
    font-size: 1.4rem;
    margin-bottom: 1.2rem;
  }
`;

const HLine = styled.div`
  background-color: #ebedef;
  height: 1px;
  width: 100%;
  margin: 5px 0px 20px;
`;

const EditOrgModal = (props: EditOrgModalProps) => {
  const isMobile = useIsMobile();
  const { isOpen, close, onDelete, org, addToast } = props;
  const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
  const { main } = useStores();
  const [loading, setLoading] = useState(false);

  const config = widgetConfigs.organizations;
  const schema = [...config.schema];
  const formRef = useRef(null);
  const initValues = {
    name: org?.name,
    image: org?.img,
    show: org?.show
  };

  const [selectedImage, setSelectedImage] = useState<string>(org?.img || '');
  const [rawSelectedFile, setRawSelectedFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const onSubmitEditOrg = async (body: any) => {
    if (!org) {
      addToast('Invalid organization update', 'danger');
      return;
    }
    setLoading(true);
    try {
      let img = '';
      const formData = new FormData();
      if (rawSelectedFile) {
        formData.append('file', rawSelectedFile);
        const file = await main.uploadFile(formData);
        if (file && file.ok) {
          img = await file.json();
        }
      }

      const newOrg = {
        id: org.id,
        uuid: org.uuid,
        name: body.name || org.name,
        owner_pubkey: org.owner_pubkey,
        img: img || org.img,
        created: org.created,
        updated: org.updated,
        show: body?.show !== undefined ? body.show : org.show,
        bounty_count: org.bounty_count,
        budget: org.budget
      };

      const res = await main.updateOrganization(newOrg);
      if (res.status === 200) {
        addToast('Sucessfully updated organization', 'success');
        // update the org ui
        props.resetOrg(newOrg);
        close();
      } else {
        addToast('Error: could not update organization', 'danger');
      }
    } catch (error) {
      addToast('Error: could not update organization', 'danger');
    }
    setLoading(false);
  };

  const handleFileInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files && e.target.files[0];
    if (file) {
      // Display the selected image
      const imageUrl = URL.createObjectURL(file);
      setSelectedImage(imageUrl);
      setRawSelectedFile(file);
    } else {
      // Handle the case where the user cancels the file dialog
      setSelectedImage('');
    }
  };

  const handleDrop = (event: DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    const file = event.dataTransfer.files[0];
    if (file) {
      const imageUrl = URL.createObjectURL(file);
      setSelectedImage(imageUrl);
      setRawSelectedFile(file);
    }
  };

  const handleDragOver = (event: DragEvent<HTMLDivElement>) => {
    event.preventDefault();
  };

  // const [files, setFiles] = useState<{ preview: string }[]>([]);

  const openFileDialog = () => {
    if (fileInputRef.current) {
      fileInputRef.current.click();
    }
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
          width: isMobile ? '100%' : '34.4375rem',
          minWidth: '20px',
          height: isMobile ? 'auto' : '27.1875rem',
          padding: '3rem 3rem 1.5rem 3rem',
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
            <ImgDashContainer onDragOver={handleDragOver} onDrop={handleDrop}>
              <UploadImageContainer onClick={openFileDialog}>
                <img src="/static/badges/ResetOrgProfile.svg" alt="upload" />
              </UploadImageContainer>
              <ImgContainer>
                {selectedImage ? (
                  <SelectedImg src={selectedImage} alt="selected file" />
                ) : (
                  <ImgText>LOGO</ImgText>
                )}
              </ImgContainer>
            </ImgDashContainer>
            <ImgTextContainer>
              <InputFile
                type="file"
                id="file-input"
                accept=".jpg, .jpeg, .png, .gif"
                onChange={handleFileInputChange}
                ref={fileInputRef}
              />
              <ImgInstructionText>
                Drag and drop or{' '}
                <ImgInstructionSpan onClick={openFileDialog}>Browse</ImgInstructionSpan>
              </ImgInstructionText>
              <ImgDetailInfo>PNG, JPG or GIF, Min. 300 x 300 px</ImgDetailInfo>
            </ImgTextContainer>
          </OrgEditImageWrapper>
          <FormWrapper style={{ marginLeft: isMobile ? '0px' : '28px' }}>
            <Formik
              initialValues={initValues || {}}
              onSubmit={onSubmitEditOrg}
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
                <OrgInputContainer>
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
                    loading={loading}
                    style={{
                      width: '100%',
                      height: '50px',
                      borderRadius: '5px',
                      alignSelf: 'center'
                    }}
                    color={'primary'}
                    text={'Save changes'}
                  />
                </OrgInputContainer>
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
            borderWidth: '1px',
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
