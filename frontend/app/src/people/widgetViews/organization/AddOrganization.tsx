import React, { useState, DragEvent, ChangeEvent } from 'react';
import styled from 'styled-components';
import { useStores } from 'store';
import { EuiGlobalToastList, EuiLoadingSpinner } from '@elastic/eui';
import { Toast } from './interface';

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

const OrgInputContainer = styled.div`
  width: 16rem;
  display: flex;
  flex-direction: column;
  @media only screen and (max-width: 500px) {
    width: 100%;
    margin-top: 1rem;
  }
`;
const OrgImgOutterContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const UploadImageContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.37756rem;
  height: 2.37756rem;
  position: absolute;
  bottom: 0;
  right: 0;
  cursor: pointer;
`;

const ImgDashContainer = styled.div`
  width: 8.875rem;
  height: 8.875rem;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  border: 1px dashed #d0d5d8;
  padding: 0.5rem;
  position: relative;
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
const ImgText = styled.h3`
  color: #b0b7bc;
  text-align: center;
  font-family: 'Barlow';
  font-size: 1.875rem;
  font-style: normal;
  font-weight: 800;
  line-height: 1.0625rem;
  letter-spacing: 0.01875rem;
  text-transform: uppercase;
  opacity: 0.5;
  margin-bottom: 0;
`;

const ImgTextContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  margin-top: 1rem;
`;

const ImgInstructionText = styled.p`
  color: #5f6368;
  text-align: center;
  font-family: 'Roboto';
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 400;
  line-height: 1.0625rem;
  letter-spacing: 0.00813rem;
  margin-bottom: 0;
`;

const ImgInstructionSpan = styled.span`
  color: #618aff;
  cursor: pointer;
`;

const ImgDetailInfo = styled.p`
  color: #b0b7bc;
  text-align: center;
  font-family: 'Roboto';
  font-size: 0.625rem;
  font-style: normal;
  font-weight: 400;
  line-height: 1.125rem;
  margin-bottom: 0;
  margin-top: 1rem;
`;

const OrgLabel = styled.label`
  color: #5f6368;
  font-family: 'Barlow';
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 500;
  /* line-height: 2.1875rem; */
  margin-bottom: 0.75rem;
`;

const OrgInput = styled.input`
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  border: 2px solid #82b4ff;
  outline: none;
  caret-color: #618aff;
  color: #3c3f41;
  font-family: 'Barlow';
  font-size: 0.9375rem;
  font-style: normal;
  font-weight: 500;
  line-height: 2.1875rem;
  width: 100%;

  ::placeholder {
    color: #b0b7bc;
    font-family: 'Barlow';
    font-size: 0.9375rem;
    font-style: normal;
    font-weight: 400;
    line-height: 2.1875rem;
  }
`;

const SelectedImg = styled.img`
  width: 7.875rem;
  height: 7.875rem;
  border-radius: 50%;
  object-fit: cover;
`;

const OrgButton = styled.button`
  width: 100%;
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

const InputFile = styled.input`
  display: none;
`;

const AddOrganization = (props: {
  closeHandler: () => void;
  getUserOrganizations: () => void;
  owner_pubkey: string | undefined;
}) => {
  const [orgName, setOrgName] = useState('');
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const { main } = useStores();
  const [toasts, setToasts] = useState<Toast[]>([]);

  const handleOrgNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setOrgName(e.target.value);
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
          <OrgLabel>Organization Name</OrgLabel>
          <OrgInput
            placeholder="My Organization..."
            value={orgName}
            onChange={handleOrgNameChange}
          />
          <OrgButton disabled={!orgName} onClick={addOrganization}>
            {isLoading ? <EuiLoadingSpinner size="m" /> : 'Add Organization'}
          </OrgButton>
        </OrgInputContainer>
      </OrgDetailsContainer>
      <EuiGlobalToastList toasts={toasts} dismissToast={removeToast} toastLifeTimeMs={3000} />
    </AddOrgWrapper>
  );
};

export default AddOrganization;
