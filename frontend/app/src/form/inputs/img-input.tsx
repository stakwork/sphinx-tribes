import React, { useState } from 'react';
import styled from 'styled-components';
import { EuiButton } from '@elastic/eui';
import Camera from '../../utils/camera-option-icon.svg';
import Dropzone from 'react-dropzone';
import avatarIcon from '../../utils/profile_avatar.svg';
import backgroundIcon from '../../utils/background_icon.svg';

import type { Props } from './propsType';
import { EuiLoadingSpinner } from '@elastic/eui';
import { useStores } from '../../store';
import { Button, Modal } from '../../sphinxUI';
import { MAX_UPLOAD_SIZE } from '../../people/utils/constants';
import { Note } from './index';

export default function ImageInput({
  label,
  note,
  value,
  handleChange,
  handleBlur,
  handleFocus,
  notProfilePic,
  imageIcon
}: Props) {
  const { ui } = useStores();
  const [uploading, setUploading] = useState(false);
  const [showError, setShowError] = useState('');
  const [picsrc, setPicsrc] = useState('');

  async function uploadBase64Pic(img_base64: string, img_type: string) {
    console.log('uploadBase64Pic', img_type);
    try {
      const info = ui.meInfo as any;
      if (!info) {
        alert('You are not logged in.');
        return console.log('no meInfo');
      }
      const URL = info.url.startsWith('http') ? info.url : `https://${info.url}`;
      const r = await fetch(`${URL}/public_pic`, {
        method: 'POST',
        body: JSON.stringify({
          img_base64,
          img_type
        }),
        headers: {
          'x-jwt': info.jwt,
          'Content-Type': 'application/json'
        }
      });
      const j = await r.json();

      if (j.success && j.response && j.response.img) {
        setPicsrc(img_base64);
        handleChange(j.response.img);
      }
    } catch (e) {
      console.log('ERROR UPLOADING IMAGE', e);
    }
  }

  async function dropzoneUpload(files: File[], fileRejections) {
    console.log('fileRejections', fileRejections);
    if (fileRejections.length) {
      fileRejections.forEach((file) => {
        file.errors.forEach((err) => {
          if (err.code === 'file-too-large') {
            setShowError(`Error: ${err.message}`);
          }
          if (err.code === 'file-invalid-type') {
            setShowError(`Error: ${err.message}`);
          }
        });
      });
      console.log('upload error');
      return;
    }

    console.log(files);
    const file = files[0];
    setUploading(true);
    const reader = new FileReader();
    reader.onload = async (event: any) => {
      await uploadBase64Pic(event.target.result, file.type);
      setUploading(false);
    };
    reader.readAsDataURL(file);
  }

  const addedStyle = notProfilePic
    ? {
        borderRadius: 0
      }
    : {};

  const defaultIcon = notProfilePic ? backgroundIcon : avatarIcon;

  return (
    <ImageWrap>
      <Dropzone multiple={false} onDrop={dropzoneUpload} maxSize={MAX_UPLOAD_SIZE}>
        {({ getRootProps, getInputProps, isDragActive, open }) => (
          <DropzoneStuff>
            {imageIcon ? (
              <DottedCircle isDragActive={isDragActive} style={addedStyle}>
                <input {...getInputProps()} />
                <ImageCircle style={addedStyle}>
                  {!uploading ? (
                    <Image
                      style={{
                        backgroundImage: `url(${
                          picsrc ? picsrc : value ? value : uploading ? '' : defaultIcon
                        })`,
                        ...addedStyle
                      }}
                    />
                  ) : (
                    <EuiLoadingSpinner size="xl" />
                  )}
                </ImageCircle>
                <div
                  style={{
                    position: 'absolute',
                    height: '38px',
                    width: '38px',
                    top: '260px',
                    left: '125px'
                  }}
                  {...getRootProps()}
                >
                  <img
                    src="/static/badges/Camera.png"
                    height={'100%'}
                    width={'100%'}
                    alt="camera_icon"
                  />
                </div>
              </DottedCircle>
            ) : (
              <DottedCircle {...getRootProps()} isDragActive={isDragActive} style={addedStyle}>
                <input {...getInputProps()} />
                <ImageCircle style={addedStyle}>
                  {!uploading ? (
                    <Image
                      style={{
                        backgroundImage: `url(${
                          picsrc ? picsrc : value ? value : uploading ? '' : defaultIcon
                        })`,
                        ...addedStyle
                      }}
                    />
                  ) : (
                    <EuiLoadingSpinner size="xl" />
                  )}
                </ImageCircle>
              </DottedCircle>
            )}

            {/* <div style={{ color: "#6B7A8D", marginTop: 5 }}>Drag and drop or</div>
            <EuiButton onClick={open}
              style={{
                borderColor: "#6B7A8D",
                color: "white",
                fontWeight: 400,
                fontSize:12, 
                marginTop: 5,
                marginBottom: 10,
              }}
              iconType={Camera}
              iconSide="right"
            >
              Change Image
            </EuiButton> */}
          </DropzoneStuff>
        )}
      </Dropzone>

      <Modal visible={showError ? true : false}>
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center',
            padding: 20
          }}
        >
          <div style={{ marginBottom: 20 }}>{showError}</div>
          <Button onClick={() => setShowError('')} text={'Okay'} color={'primary'} />
        </div>
      </Modal>

      {note && <Note>*{note}</Note>}
    </ImageWrap>
  );
}

const DropzoneStuff = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
`;
const Image = styled.div`
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
  height: 100px;
  width: 100px;
  border-radius: 50%;
`;

const ImageSq = styled.div`
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
  height: 100px;
  width: 100px;
`;

const ImageWrap = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
`;
export interface DottedCircleProps {
  isDragActive: boolean;
}
const DottedCircle = styled.div<DottedCircleProps>`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 120px;
  width: 120px;
  border-radius: 50%;
  border-style: dashed;
  border-color: ${(p) => (p.isDragActive ? 'white' : '#6b7a8d')};
  border-width: thin;
  cursor: pointer;
`;

const ImageCircle = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100px;
  width: 100px;
  border-radius: 50%;
  position: relative;
`;

const ImageSquare = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100px;
  width: 100px;
  position: relative;
`;
