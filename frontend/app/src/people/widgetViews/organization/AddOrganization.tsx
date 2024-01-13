import React, { useState, DragEvent, ChangeEvent } from 'react';
import styled from 'styled-components';
import { useStores } from 'store';
import { EuiGlobalToastList, EuiLoadingSpinner } from '@elastic/eui';
import { Toast } from './interface';
import {
  ImgDashContainer,
  ImgDetailInfo,
  ImgInstructionSpan,
  ImgInstructionText,
  ImgText,
  ImgTextContainer,
  InputFile,
  TextInput,
  OrgInputContainer,
  OrgLabel,
  SelectedImg,
  UploadImageContainer,
  TextAreaInput,
  SecondaryText
} from './style';

const AddOrgWrapper = styled.div`
  padding: 3rem;
  display: flex;
  flex-direction: column;

  @media only screen and (max-width: 500px) {
    padding: 1.2rem;
    width: 100%;
  }
`;

const AddOrgHeader = styled.h2`
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
const OrgDetailsContainer = styled.div`
  margin-top: 3rem;
  display: flex;
  gap: 3.56rem;
  @media only screen and (max-width: 500px) {
    flex-direction: column;
    gap: 0.5rem;
  }
`;
const FooterContainer = styled.div`
  display: flex;
  gap: 3.56rem;
  align-items: end;
  justify-content: space-between;

  @media only screen and (max-width: 500px) {
    flex-direction: column;
    gap: 0.5rem;
  }
`;

const OrgImgOutterContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const ImgContainer = styled.div`
  width: 7.875rem;
  height: 7.875rem;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #ebedf1;
`;

const OrgButton = styled.button`
  width: 16rem;
  height: 3rem;
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  font-family: 'Barlow';
  font-size: 0.9375rem;
  font-style: normal;
  font-weight: 500;
  line-height: 0rem;
  letter-spacing: 0.00938rem;
  margin-top: 1.5rem;
  border: none;
  background: var(--Primary-blue, #618aff);
  box-shadow: 0px 2px 10px 0px rgba(97, 138, 255, 0.5);
  color: #fff;

  :disabled {
    border: 1px solid rgba(0, 0, 0, 0.07);
    background: rgba(0, 0, 0, 0.04);
    color: rgba(142, 150, 156, 0.85);
    box-shadow: none;
  }
`;

const AddOrganization = (props: {
  closeHandler: () => void;
  getUserOrganizations: () => void;
  owner_pubkey: string | undefined;
}) => {
  const [orgName, setOrgName] = useState('');
  const [websiteName, setWebsiteName] = useState('');
  const [githubRepo, setGithubRepo] = useState('');
  const [description, setDescription] = useState('');
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const { main } = useStores();
  const [toasts, setToasts] = useState<Toast[]>([]);

  const handleOrgNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setOrgName(e.target.value);
  };

  const handleWebsiteNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setWebsiteName(e.target.value);
  };

  const handleGithubRepoChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setGithubRepo(e.target.value);
  };
  const handleDescriptionChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setDescription(e.target.value);
  };

  const handleDrop = (event: DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    const file = event.dataTransfer.files[0];
    if (file) {
      setSelectedFile(file);
    }
  };

  const handleDragOver = (event: DragEvent<HTMLDivElement>) => {
    event.preventDefault();
  };

  const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    const fileList = event.target.files;

    if (fileList) {
      setSelectedFile(fileList[0]);
    }
  };

  function addSuccessToast() {
    setToasts([
      {
        id: '1',
        title: 'Create Organization',
        color: 'success',
        text: 'Organization created successfully'
      }
    ]);
  }

  function addErrorToast(text: string) {
    setToasts([
      {
        id: '2',
        title: 'Create Organization',
        color: 'danger',
        text
      }
    ]);
  }

  function removeToast() {
    setToasts([]);
  }

  const handleBrowse = () => {
    const fileInput = document.getElementById('file-input');
    fileInput?.click();
  };

  const addOrganization = async () => {
    try {
      setIsLoading(true);
      let img_url = '';
      const formData = new FormData();
      if (selectedFile) {
        formData.append('file', selectedFile);
        const file = await main.uploadFile(formData);
        if (file && file.ok) {
          img_url = await file.json();
        }
      }
      const body = { owner_pubkey: props.owner_pubkey || '', name: orgName, img: img_url };

      const res = await main.addOrganization(body);
      if (res.status === 200) {
        addSuccessToast();
        setTimeout(async () => {
          await props.getUserOrganizations();
          setIsLoading(false);
          props.closeHandler();
        }, 500);
      } else {
        addErrorToast(await res.json());
        setIsLoading(false);
      }
    } catch (error) {
      addErrorToast('Error occured while creating organization');
      console.error('Error occured', error);
      setIsLoading(false);
    }
  };

  return (
    <AddOrgWrapper>
      <AddOrgHeader>Add New Organization</AddOrgHeader>
      <OrgDetailsContainer>
        <OrgImgOutterContainer>
          <ImgDashContainer onDragOver={handleDragOver} onDrop={handleDrop}>
            <UploadImageContainer onClick={handleBrowse}>
              <img src="/static/upload.svg" alt="upload" />
            </UploadImageContainer>
            <ImgContainer>
              {selectedFile ? (
                <SelectedImg src={URL.createObjectURL(selectedFile)} alt="selected file" />
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
              onChange={handleFileChange}
            />
            <ImgInstructionText>
              Drag and drop or{' '}
              <ImgInstructionSpan onClick={handleBrowse}>Browse</ImgInstructionSpan>
            </ImgInstructionText>
            <ImgDetailInfo>PNG, JPG or GIF, Min. 300 x 300 px</ImgDetailInfo>
          </ImgTextContainer>
        </OrgImgOutterContainer>
        <OrgInputContainer>
          <OrgLabel>Organization Name *</OrgLabel>
          <TextInput
            placeholder="My Organization..."
            value={orgName}
            onChange={handleOrgNameChange}
          />
          <OrgLabel>Website</OrgLabel>
          <TextInput
            placeholder="Website URL..."
            value={websiteName}
            onChange={handleWebsiteNameChange}
          />
          <OrgLabel>Github repo</OrgLabel>
          <TextInput
            placeholder="Github link..."
            value={githubRepo}
            onChange={handleGithubRepoChange}
          />
        </OrgInputContainer>
        <OrgInputContainer>
          <OrgLabel>Description</OrgLabel>
          <TextAreaInput
            placeholder="Description Text..."
            rows={7}
            value={description}
            onChange={handleDescriptionChange}
          />
        </OrgInputContainer>
      </OrgDetailsContainer>
      <FooterContainer>
        <SecondaryText>* Required fields</SecondaryText>
        <OrgButton disabled={!orgName} onClick={addOrganization}>
          {isLoading ? <EuiLoadingSpinner size="m" /> : 'Add Organization'}
        </OrgButton>
      </FooterContainer>
      <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={3000} />
    </AddOrgWrapper>
  );
};

export default AddOrganization;
