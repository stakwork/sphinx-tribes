import React, {useState} from 'react'
import styled from 'styled-components'
import {EuiButton} from '@elastic/eui'
import Camera from '../../utils/camera-option-icon.svg'
import Dropzone from "react-dropzone";
import avatarIcon from '../../utils/profile_avatar.svg'
// import {EuiFilePicker} from '@elastic/eui'



export default function ImageInput(props) {
  // const {meme} = useStores();
  const [uploading, setUploading] = useState(false);
  const [picsrc, setPicsrc] = useState("")
  
  console.log("PROPERS === ", props)
  // return <EuiFilePicker value={props.initialValues.img} >
  //   </EuiFilePicker>
  
  async function dropzoneUpload(files) {
    
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

  return <ImageWrap>
            <Dropzone multiple={false} onDrop={dropzoneUpload}>
              {({ getRootProps, getInputProps, isDragActive }) => (
                <DottedCircle {...getRootProps()}>
                  <input {...getInputProps()} />
                    <ImageCircle>
                      <Image style={{backgroundImage: `url(${picsrc ? picsrc + "?thumb=true" : avatarIcon})` }} />
                    </ImageCircle>
              </DottedCircle>
              )}
            </Dropzone>
            <div style={{color: "#6B7A8D", marginTop: 5}}>Drag and drop or</div>
            <EuiButton style={{borderColor: "#6B7A8D", color:"white", fontWeight:400, marginTop: 5, marginBottom: 10}} iconType={Camera} iconSide="right" >
              Change Image
            </EuiButton>
          </ImageWrap>
}

const Image = styled.div`
  background-position: center;
  background-repeat: no-repeat;
  background-size:cover;
  height: 100px;
  width: 100px;
  border-radius: 50%;
`

const ImageWrap = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
`

const DottedCircle=styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 120px;
  width: 120px;
  border-radius: 50%;
  border-style: dashed;
  border-color: #6B7A8D;
  border-width: thin;
`

const ImageCircle=styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100px;
  width: 100px;
  border-radius: 50%;
`