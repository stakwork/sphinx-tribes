import React, { useEffect, useState } from 'react';
import { setup, isSupported } from '@loomhq/record-sdk';
import styled from 'styled-components';
import { EuiText } from '@elastic/eui';
import { LoomViewProps } from 'people/interfaces';
import { colors } from '../../config/colors';
import { Button, IconButton } from '../../components/common';

const PUBLIC_APP_ID = 'ded90c8e-92ed-496d-bfe3-f742d7fa9785';

const BUTTON_ID = 'loom-record-sdk-button';

interface styledProps {
  color?: any;
}

const ButtonContainer = styled.div`
  display: flex;
  justify-content: center;
`;

const RemoveButtonContainer = styled.div<styledProps>`
  position: absolute;
  top: 160px;
  right: 45px;
  display: flex;
  justify-content: flex-end;
  align-item: center;
  cursor: pointer;
  .buttonText {
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 500;
    font-size: 14px;
    line-height: 35px;
    display: flex;
    align-items: center;
    text-align: right;
    color: #5f6368;
  }
`;

export default function LoomViewerRecorderNew(props: LoomViewProps) {
  const { loomEmbedUrl, onChange, readOnly, style, setIsVideo } = props;
  const [videoUrl, setVideoUrl] = useState(loomEmbedUrl || '');
  const color = colors['light'];

  useEffect(() => {
    async function setupLoom() {
      const { supported, error } = await isSupported();

      if (!supported) {
        console.warn(`Error setting up Loom: ${error}`);
        return;
      }

      const button = document.getElementById(BUTTON_ID);

      if (!button) {
        return;
      }

      const { configureButton } = await setup({
        publicAppId: PUBLIC_APP_ID
      });

      const sdkButton = configureButton({ element: button });

      sdkButton.on('insert-click', async (video: any) => {
        setVideoUrl(video.embedUrl);
        if (setIsVideo) setIsVideo(true);
        if (onChange) onChange(video.embedUrl);
      });
    }

    setupLoom();
  }, []);

  if (readOnly && !videoUrl) {
    return null;
  }

  const loomViewer = videoUrl && (
    <div
      dangerouslySetInnerHTML={{
        __html: `<div class="lo-emb-vid"
        style="position: relative; padding-bottom: 75%; height: 0; margin-top: -101px; margin-left: -38px">
        <iframe src="${videoUrl}"
            style="position: absolute; top: 0; left: 0; width: 290px; height: 175px; border-radius: 12px; " frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe></div>`
      }}
    />
  );

  return (
    <div style={style}>
      {!readOnly && (
        <ButtonContainer>
          <Button
            text={'Record Loom Video'}
            color={'white'}
            id={BUTTON_ID}
            style={{
              width: 214,
              height: 42,
              display: 'flex',
              justifyContent: 'space-between'
            }}
            img={'loom.png'}
            imgStyle={{
              height: 25,
              width: 25,
              marginRight: '20px'
            }}
            ButtonTextStyle={{
              fontFamily: 'Barlow',
              fontStyle: 'normal',
              fontWeight: '700',
              fontSize: '14px',
              lineHeight: '17px',
              display: 'flex',
              alignItems: 'center',
              textAlign: 'center',
              color: color.grayish.G50
            }}
          />
          {videoUrl ? (
            <RemoveButtonContainer
              onClick={() => {
                setVideoUrl('');
                if (setIsVideo) setIsVideo('');
                if (onChange) onChange('');
              }}
              color={color}
            >
              <IconButton
                color={'widget'}
                icon={'delete'}
                // text={'Delete Video'}

                iconStyle={{
                  fontSize: '20px'
                }}
              />
              <EuiText className="buttonText">Remove</EuiText>
            </RemoveButtonContainer>
          ) : (
            <div />
          )}
        </ButtonContainer>
      )}
      {loomViewer}
    </div>
  );
}
