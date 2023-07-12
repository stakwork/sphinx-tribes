import React, { useState } from 'react';
import styled from 'styled-components';
import Dropzone from 'react-dropzone';
import MaterialIcon from '@material/react-material-icon';
import { observer } from 'mobx-react-lite';
import { useStores } from '../../../store';

import { Button, Modal } from '../../common';
import { MAX_UPLOAD_SIZE } from '../../../people/utils/Constants';
import { colors } from '../../../config/colors';
import type { Props } from './propsType';

const Wrapper = styled.div`
  display: flex;
  align-items: center;
  flex-wrap: wrap;
`;

const DropzoneStuff = styled.div``;
const Image = styled.div`
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
  height: 105px;
  width: 105px;
`;
const Close = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  position: absolute;
  top: 2px;
  right: 2px;
  height: 20px;
  width: 20px;
  background: white;
  border-radius: 50%;
  z-index: 10;
  cursor: pointer;
`;
const ImageWrap = styled.div`
  display: flex;
  margin: 2px;
  position: relative;
`;
const Sq = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  border-width: thin;
`;

const ImageCircle = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
`;

function GalleryInput({ value, handleChange }: Props) {
  const color = colors['light'];
  const { ui } = useStores();
  const [uploading, setUploading] = useState(false);
  const [showError, setShowError] = useState('');
  const picsrcArray = value || [];
  async function addImg(img: any) {
    const picsClone = [...picsrcArray];
    picsClone.push(img);
    handleChange(picsClone);
  }

  async function uploadBase64Pic(img_base64: string, img_type: string) {
    try {
      const info = ui.meInfo as any;
      if (!info) return;
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
        // addNewImg(img_base64)
        addImg(j.response.img);
      }
    } catch (e) {
      console.log('ERROR UPLOADING IMAGE', e);
    }
  }

  async function dropzoneUpload(files: File[], fileRejections: any) {
    if (fileRejections.length) {
      fileRejections.forEach((file: any) => {
        file.errors.forEach((err: any) => {
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

    const file = files[0];

    setUploading(true);
    const reader = new FileReader();
    reader.onload = async (event: any) => {
      await uploadBase64Pic(event.target.result, file.type);
      setUploading(false);
    };

    reader.readAsDataURL(file);
  }

  async function deleteImg(index: any) {
    const picsClone = [...picsrcArray];
    picsClone.splice(index, 1);
    handleChange(picsClone);
  }

  return (
    <>
      <Wrapper>
        {picsrcArray &&
          picsrcArray.map((v: any, i: number) => (
            <ImageWrap key={i}>
              <Close onClick={() => deleteImg(i)}>
                <MaterialIcon icon={'close'} style={{ color: color.pureBlack, fontSize: 12 }} />
              </Close>
              <Sq>
                <ImageCircle>
                  <Image
                    style={{
                      backgroundImage: `url(${v})`
                    }}
                  />
                </ImageCircle>
              </Sq>
            </ImageWrap>
          ))}
      </Wrapper>

      <div style={{ marginTop: 5 }}>
        <Dropzone multiple={false} onDrop={dropzoneUpload} maxSize={MAX_UPLOAD_SIZE}>
          {({ getRootProps, getInputProps }: any) => (
            <DropzoneStuff>
              <div>
                <input {...getInputProps()} />
                <Button
                  {...getRootProps()}
                  leadingIcon={'add'}
                  style={{
                    width: 154,
                    paddingRight: 20
                  }}
                  height={48}
                  text={'Add Media'}
                  color="widget"
                  loading={uploading}
                />
              </div>
            </DropzoneStuff>
          )}
        </Dropzone>
      </div>

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
    </>
  );
}

export default observer(GalleryInput);
