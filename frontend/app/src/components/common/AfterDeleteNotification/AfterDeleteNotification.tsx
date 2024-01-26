import React from 'react';

import { EuiText } from '@elastic/eui';
import { Box, Stack } from '@mui/system';
import { colors } from '../../../config';
import { BaseModal } from '../BaseModal';
import { ButtonContainer } from '../ButtonContainer';
import { useCreateModal } from '../useCreateModal';
import { ReactComponent as CloseIcon } from './close.svg';
import { image } from './image';

export type AfterDeleteNotificationProps = {
  onClose: () => void;
  bountyTitle?: string;
  bountyLink?: string;
};
export const AfterDeleteNotification = ({
  onClose,
  bountyLink,
  bountyTitle
}: AfterDeleteNotificationProps) => {
  const closeHandler = () => {
    onClose();
  };

  const linkHandler = () => {
    window.open(bountyLink, '_blank');
  };

  const color = colors['light'];

  return (
    <BaseModal backdrop="white" open onClose={closeHandler}>
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
            Your Bounty is <br />
            <Box component="span" fontWeight={700}>
              Successfully Deleted
            </Box>
          </Box>

          <Box color="rgba(142, 150, 156, 1)" component={EuiText}>
            {bountyTitle}
          </Box>
        </Stack>
        {bountyLink && (
          <Stack
            position="absolute"
            bottom={{ xs: '0.5rem', sm: '2rem' }}
            alignItems="center"
            justifyContent="center"
            spacing={1}
          >
            <Box
              // mt={{ xs: '50px', sm: '100px' }}
              color="rgba(142, 150, 156, 1)"
              component={EuiText}
            >
              Your original Github Ticket is still available
            </Box>
            <ButtonContainer onClick={linkHandler} color={color}>
              <div className="LeadingImageContainer">
                <img
                  className="buttonImage"
                  src={'/static/github_icon.svg'}
                  alt={'github_ticket'}
                  height={'20px'}
                  width={'20px'}
                />
              </div>
              <EuiText className="ButtonText">Github Ticket</EuiText>
              <Box
                component="img"
                ml="3rem"
                className="buttonImage"
                src={'/static/github_ticket.svg'}
                alt={'github_ticket'}
                height={'14px'}
                width={'14px'}
              />
            </ButtonContainer>
          </Stack>
        )}
      </Stack>
    </BaseModal>
  );
};

export const useAfterDeleteNotification = () => {
  const openModal = useCreateModal();

  const openAfterDeleteNotification = (props: Omit<AfterDeleteNotificationProps, 'onClose'>) => {
    openModal({
      Component: AfterDeleteNotification,
      props
    });
  };
  return { openAfterDeleteNotification };
};
