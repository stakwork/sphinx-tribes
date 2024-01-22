import React, { useRef, useState, ChangeEvent, DragEvent } from 'react';
import { useIsMobile } from 'hooks/uiHooks';
import styled from 'styled-components';
import { Formik } from 'formik';
import { validator } from 'components/form/utils';
import { widgetConfigs } from 'people/utils/Constants';
import { FormField } from 'components/form/utils';
import { useStores } from 'store';
import Input from '../../../components/form/inputs';
import { Button, Modal } from '../../../components/common';
import { colors } from '../../../config/colors';
import {
  ImgContainer,
  ImgDashContainer,
  ImgDetailInfo,
  ImgInstructionSpan,
  ImgInstructionText,
  ImgTextContainer,
  InputFile,
  ModalTitle,
  SelectedImg,
  UploadImageContainer
} from './style';
import { EditOrgModalProps } from './interface';
import DeleteOrgWindow from './DeleteOrgWindow';

const color = colors['light'];

const EditOrgWrapper = styled.div`
  padding: 2.375rem 3rem 3rem 3rem;
  display: flex;
  gap: 38px;
  flex-direction: column;
  position: relative;
  width: 100%;

  @media only screen and (max-width: 500px) {
    padding: 1rem 1.2rem 1.2rem 1.2rem;
    justify-content: center;
    gap: 0;
  }
`;

const InputWrapper = styled.div`
  display: grid;
  grid-template-columns: 242px 256px;
  grid-template-rows: repeat(3, 61px);
  grid-column-gap: 32px;
  grid-row-gap: 20px;

  @media only screen and (max-width: 500px) {
    display: flex;
    flex-direction: column;
    width: 100%;
    gap: 0;
  }
`;

const LabelWrapper = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const Label = styled.span`
  font-family: Barlow;
  font-size: 13px;
  font-weight: 500;
  line-height: 35px;
  letter-spacing: 0px;
  text-align: left;
  color: #5f6368;
`;

const SecondaryText = styled.span`
  color: #b0b7bc;
  font-family: Roboto;
  font-size: 13px;
  font-weight: 400;
  line-height: 35px;
  letter-spacing: 0px;
  text-align: left;
  vertical-align: center;
`;

const EditOrgRow = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  height: max-content;
`;

const OrgEditImageWrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;

  @media only screen and (max-width: 500px) {
    margin: auto;
  }
`;

const EditOrgTitle = styled(ModalTitle)`
  color: #3c3f41;
  font-family: 'Barlow';
  font-size: 1.875rem;
  font-style: normal;
  font-weight: 800;
  line-height: 1.875rem;
  margin-bottom: 0;

  @media only screen and (max-width: 500px) {
    text-align: center;
    font-size: 1.4rem;
  }
`;

const EditOrgModal = (props: EditOrgModalProps) => {
  const { ui } = useStores();

  const isMobile = useIsMobile();
  const { isOpen, close, onDelete, org, addToast } = props;
  const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
  const { main } = useStores();
  const [loading, setLoading] = useState(false);
  const [characterCount] = useState(0);

  const config = widgetConfigs.organizations;
  const schema = [...config.schema];
  const formRef = useRef(null);
  const initValues = {
    name: org?.name,
    image: org?.img,
    description: org?.description,
    github: org?.github,
    website: org?.website,
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
        description: body.description || org?.description,
        github: body.github || org?.github,
        website: body.website || org?.website,
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
  const isOrganizationAdmin = props.org?.owner_pubkey === ui.meInfo?.owner_pubkey;

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
          minWidth: isMobile ? '100%' : '51.9375rem',
          minHeight: isMobile ? '100vh' : '29rem'
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
        <EditOrgWrapper>
          <EditOrgRow>
            <EditOrgTitle>Edit Organization</EditOrgTitle>
            <Button
              disabled={!isOrganizationAdmin}
              onClick={() => {
                setShowDeleteModal(true);
              }}
              loading={false}
              style={{
                width: isMobile ? '100%' : '170px',
                height: '40px',
                borderRadius: '6px',
                padding: '8px, 16px, 8px, 16px',
                borderStyle: 'solid',
                alignSelf: 'flex-end',
                borderWidth: '1px',
                backgroundColor: 'white',
                borderColor: '#ED7474',
                color: '#ED7474',
                position: !isMobile ? 'initial' : 'absolute',
                bottom: '3px',
                maxWidth: 'calc(100% - 2.4rem)'
              }}
              color={'#ED7474'}
              text={'Delete Organization'}
            />
          </EditOrgRow>
          <EditOrgRow>
            <OrgEditImageWrapper>
              <ImgDashContainer onDragOver={handleDragOver} onDrop={handleDrop}>
                <UploadImageContainer onClick={openFileDialog}>
                  <img src="/static/badges/ResetOrgProfile.svg" alt="upload" />
                </UploadImageContainer>
                <ImgContainer>
                  <SelectedImg
                    src={
                      selectedImage === ''
                        ? '/static/badges/editOrganisationImage.svg'
                        : selectedImage
                    }
                    alt="selected file"
                  />
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
                <InputWrapper>
                  {schema.map((item: FormField) => (
                    <div key={item.name} style={item.style}>
                      <LabelWrapper>
                        <Label>{item.label}</Label>
                        {item.maxCharacterLimit ? (
                          <SecondaryText>
                            {characterCount}/{item.maxCharacterLimit}
                          </SecondaryText>
                        ) : null}
                      </LabelWrapper>
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
                        style={{
                          width: '100%',
                          ...item.style,
                          maxHeight: isMobile ? '145px' : 'auto'
                        }}
                        newDesign
                      />
                    </div>
                  ))}
                  <Button
                    disabled={false}
                    onClick={() => handleSubmit()}
                    loading={loading}
                    style={{
                      width: '100%',
                      maxWidth: isMobile ? '100%' : '256px',
                      height: '50px',
                      borderRadius: '5px',
                      alignSelf: 'center',
                      position: isMobile ? 'initial' : 'absolute',
                      top: '368px',
                      left: '527px'
                    }}
                    color={'primary'}
                    text={'Save changes'}
                  />
                </InputWrapper>
              )}
            </Formik>
          </EditOrgRow>
          <EditOrgRow>
            <SecondaryText>* Required fields</SecondaryText>
          </EditOrgRow>
          {showDeleteModal ? (
            <DeleteOrgWindow onDeleteOrg={onDelete} close={() => setShowDeleteModal(false)} />
          ) : (
            <></>
          )}
        </EditOrgWrapper>
      </Modal>
    </>
  );
};

export default EditOrgModal;
