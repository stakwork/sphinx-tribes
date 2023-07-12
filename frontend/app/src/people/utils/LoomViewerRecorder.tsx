import React, { useEffect, useState } from 'react';
import { setup, isSupported } from '@loomhq/record-sdk';
import { LoomViewProps } from 'people/interfaces';
import { Button, IconButton } from '../../components/common';

const PUBLIC_APP_ID = 'ded90c8e-92ed-496d-bfe3-f742d7fa9785';

const BUTTON_ID = 'loom-record-sdk-button';

export default function LoomViewerRecorder(props: LoomViewProps) {
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

      sdkButton.on('insert-click', async (video: any) => {
        setVideoUrl(video.embedUrl);
        if (onChange) onChange(video.embedUrl);
      });
    }

    setupLoom();
  }, [onChange]);

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
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 10 }}>
          <Button text={'Record Loom Video'} color={'primary'} id={BUTTON_ID} />
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
        </div>
      )}
      {loomViewer}
    </div>
  );
}
