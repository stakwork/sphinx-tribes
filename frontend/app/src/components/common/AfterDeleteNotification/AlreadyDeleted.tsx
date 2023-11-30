import React from 'react';

import { EuiText } from '@elastic/eui';
import { Box, Stack } from '@mui/system';
import { BaseModal } from '../BaseModal';
import { ReactComponent as CloseIcon } from './close.svg';
import { image } from './image';
export type AlreadyDeletedProps = {
  onClose: () => void;
  bountyLink?: string;
  bountyTitle?: string;
  isDeleted: boolean;
};
export const AlreadyDeleted = ({ onClose, bountyTitle }: AlreadyDeletedProps) => {
  const closeHandler = () => {
    onClose();
  };

  return (
    <BaseModal open onClose={closeHandler}>
      <Stack
        position="relative"
        sx={{
          background:
            'linear-gradient(0deg, rgba(255,255,255,1) 0%, rgba(255,255,255,1) 50%, rgba(245,246,248,1) 50%, rgba(245,246,248,1) 100%)'
        }}
        minWidth={350}
        width="70vw"
        height="100vh"
        overflow="auto"
        p={1}
        alignItems="center"
        justifyContent="center"
      >
        <Box
          onClick={closeHandler}
          position="absolute"
          right="1rem"
          top="1rem"
          sx={{ cursor: 'pointer' }}
        >
          <CloseIcon />
        </Box>
        <Box
          position="absolute"
          sx={{
            top: '50%',
            left: '50%',
            transform: 'translate(-50%, -50%)',
            '& > svg>g': {
              boxShadow: ' 0px 1px 20px 0px rgba(0, 0, 0, 0.15)'
            }
          }}
        >
          {image}
        </Box>
        <Stack
          zIndex={1}
          useFlexGap
          alignItems="center"
          m="auto"
          direction="column"
          justifyContent="space-between"
        >
          <Box
            textAlign="center"
            color="rgba(142, 150, 156, 1)"
            fontSize={36}
            lineHeight={1.2}
            component={EuiText}
            mb={{ xs: '250px', sm: '230px' }}
          >
            This bounty has been <br />
            <Box component="span" fontWeight={700}>
              deleted
            </Box>
          </Box>
          <br />
          <Box color="rgba(142, 150, 156, 1)" component={EuiText}>
            {bountyTitle}
          </Box>
        </Stack>
      </Stack>
    </BaseModal>
  );
};
