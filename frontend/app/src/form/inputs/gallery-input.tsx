import React, { useState } from "react";
import styled from "styled-components";
import { EuiButton, EuiFormRow } from "@elastic/eui";
import Camera from "../../utils/camera-option-icon.svg";
import Dropzone from "react-dropzone";
import type { Props } from './propsType'
import { EuiLoadingSpinner } from '@elastic/eui';
import { useStores } from '../../store'
import MaterialIcon from "@material/react-material-icon";
import { FieldEnv, FieldTextArea } from './index'


export default function GalleryInput({ label, value, handleChange, handleBlur, handleFocus }: Props) {
    const { ui } = useStores();
    const [uploading, setUploading] = useState(false);
    // return <EuiFilePicker value={props.initialValues.img} >
    //   </EuiFilePicker>

    const picsrcArray = value || []

    async function uploadBase64Pic(img_base64: string, img_type: string) {
        console.log('uploadBase64Pic', img_type)
        try {
            const info = ui.meInfo as any;
            if (!info) return console.log("no meInfo");
            const URL = info.url.startsWith('http') ? info.url : `https://${info.url}`
            const r = await fetch(URL + "/public_pic", {
                method: "POST",
                body: JSON.stringify({
                    img_base64, img_type
                }),
                headers: {
                    "x-jwt": info.jwt,
                    "Content-Type": "application/json"
                },
            });
            const j = await r.json()

            if (j.success && j.response && j.response.img) {
                addImg(j.response.img)
            }
        } catch (e) {
            console.log('ERROR UPLOADING IMAGE', e)
        }
    }

    async function dropzoneUpload(files: File[]) {
        console.log(files)
        const file = files[0];
        setUploading(true)
        const reader = new FileReader();
        reader.onload = async (event: any) => {
            await uploadBase64Pic(event.target.result, file.type)
            setUploading(false)
        }
        reader.readAsDataURL(file);
    }

    async function addImg(img) {
        let picsClone = [...picsrcArray]
        picsClone.push(img)
        handleChange(picsClone)
    }

    async function deleteImg(index) {
        let picsClone = [...picsrcArray]
        picsClone.splice(index, 1)
        handleChange(picsClone)
    }

    const MAX_SIZE = 4194304 // 4MB

    return (
        <>
            <Wrapper>
                {picsrcArray && picsrcArray.map((v, i) => {
                    return <ImageWrap key={i}>
                        <Close onClick={() => deleteImg(i)}>
                            <MaterialIcon icon={'close'} style={{ color: '#000', fontSize: 12 }} />
                        </Close>
                        <Sq>
                            <ImageCircle>
                                <Image style={{
                                    backgroundImage: `url(${v})`
                                }}
                                />
                            </ImageCircle>
                        </Sq>
                    </ImageWrap>
                })}

                <ImageWrap>
                    <Dropzone multiple={false} onDrop={dropzoneUpload} maxSize={MAX_SIZE}>
                        {({ getRootProps, getInputProps, isDragActive, open }) => (
                            <DropzoneStuff>
                                <DottedSq {...getRootProps()} isDragActive={isDragActive}>
                                    <input {...getInputProps()} />
                                    <ImageCircle>

                                        <MaterialIcon icon={'add'} />
                                        {uploading && <EuiLoadingSpinner size="xl" style={{ marginTop: -14 }} />}
                                    </ImageCircle>
                                </DottedSq>
                            </DropzoneStuff>
                        )}
                    </Dropzone>

                </ImageWrap>
            </Wrapper>
        </>
    );
}

const Wrapper = styled.div`
  display: flex;
  align-items: center;
`

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
  height:55px;
  width: 80px;
`;
const Close = styled.div`
  display:flex;
  align-items:center;
  justify-content:center;
  position:absolute;
  top:0px;
  right:0px;
  height:20px;
  width:20px;
  background:white;
  border-radius:50%;
  z-index:10;
  cursor: pointer;
`;
const ImageWrap = styled.div`
  display: flex;
  margin:2px;
  position:relative;
`;
export interface DottedCircleProps {
    isDragActive?: boolean;
}
const DottedSq = styled.div<DottedCircleProps>`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 60px;
  width: 90px;
  border-style: dashed;
  border-color: ${p => p.isDragActive ? 'white' : '#6b7a8d'};
  border-width: thin;
  cursor:pointer;
`;
const Sq = styled.div<DottedCircleProps>`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 60px;
  width: 90px;
  border-color: ${p => p.isDragActive ? 'white' : '#6b7a8d'};
  border-width: thin;
`;

const ImageCircle = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 80px;
  width: 80px;
  position:relative;
`;
