import React, { useEffect, useState } from 'react';
import { setup, isSupported } from '@loomhq/record-sdk';
import { Button, IconButton } from '../../sphinxUI';
import styled from 'styled-components';

const PUBLIC_APP_DEVELOPMENT_ID = 'beec6b9b-d84c-44f4-ba70-f63f32f9e603';

const PUBLIC_APP_ID = 'ded90c8e-92ed-496d-bfe3-f742d7fa9785';

const BUTTON_ID = 'loom-record-sdk-button';

export default function LoomViewerRecorderNew(props) {
  const { loomEmbedUrl, onChange, readOnly, style } = props;
  const [videoUrl, setVideoUrl] = useState(loomEmbedUrl || '');

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

      sdkButton.on('insert-click', async (video) => {
        setVideoUrl(video.embedUrl);
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
        style="position: relative; padding-bottom: 75%; height: 0;">
        <iframe src="${videoUrl}"
            style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe></div>`
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
              color: '#5F6368'
            }}
          />
          {videoUrl ? (
            <IconButton
              color={'widget'}
              icon={'close'}
              // text={'Delete Video'}
              onClick={() => {
                setVideoUrl('');
                if (onChange) onChange('');
              }}
            />
          ) : (
            <div />
          )}
        </ButtonContainer>
      )}
      {loomViewer}
    </div>
  );
}

interface styleProps {
  color?: any;
}

const ButtonContainer = styled.div`
  display: flex;
  justify-content: center;
`;
