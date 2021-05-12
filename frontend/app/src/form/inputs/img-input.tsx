import React, { useState } from "react";
import styled from "styled-components";
import { EuiButton } from "@elastic/eui";
import Camera from "../../utils/camera-option-icon.svg";
import Dropzone from "react-dropzone";
import avatarIcon from "../../utils/profile_avatar.svg";
import type {Props} from './propsType'
import { EuiLoadingSpinner } from '@elastic/eui';

export default function ImageInput({label, value, handleChange, handleBlur, handleFocus}:Props) {
  // const {meme} = useStores();
  const [uploading, setUploading] = useState(false);
  const [picsrc, setPicsrc] = useState("");
  // return <EuiFilePicker value={props.initialValues.img} >
  //   </EuiFilePicker>

  async function dropzoneUpload(files:File[]) {
    console.log(files)
    // const file = files[0];
    // const server = meme.getDefaultServer();
    // setUploading(true);
    // const r = await uploadFile(
    //   file,
    //   file.type,
    //   server.host,
    //   server.token,
    //   "Image.jpg",
    //   true
    // );
    // if (r && r.muid) {
    //   // console.log(`https://${server.host}/public/${r.muid}`)
    //   setPicsrc(`https://${server.host}/public/${r.muid}`);
    // }
  }

  return (
    <ImageWrap>
      <Dropzone multiple={false} onDrop={dropzoneUpload}>
        {({ getRootProps, getInputProps, isDragActive, open }) => (
          <DropzoneStuff>
            <DottedCircle {...getRootProps()} isDragActive={isDragActive}>
              <input {...getInputProps()} />
              <ImageCircle>
                <Image style={{
                  backgroundImage: `url(${
                    picsrc ? picsrc + "?thumb=true" : (uploading ? '' : avatarIcon)
                  })`}}
                />
                {uploading && <EuiLoadingSpinner size="xl" style={{marginTop:-14}} />}
              </ImageCircle>
            </DottedCircle>
            <div style={{ color: "#6B7A8D", marginTop: 5 }}>Drag and drop or</div>
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
            </EuiButton>
          </DropzoneStuff>
        )}
      </Dropzone>
    </ImageWrap>
  );
}

const DropzoneStuff = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction:column;
`
const Image = styled.div`
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
  height: 100px;
  width: 100px;
  border-radius: 50%;
`;

const ImageWrap = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
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
  border-color: ${p=> p.isDragActive?'white':'#6b7a8d'};
  border-width: thin;
`;

const ImageCircle = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100px;
  width: 100px;
  border-radius: 50%;
  position:relative;
`;
