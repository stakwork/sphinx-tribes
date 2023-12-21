import React from 'react';
import { Button } from 'components/common';
import AccessDeniedImage from '../../../public/static/access_denied.png';
import { AccessImg, Container, DeniedSmall, DeniedText, Wrap } from './style';

const AdminAccessDenied = () => {
  return (
    <Container>
      <Wrap>
        <AccessImg src={AccessDeniedImage} />
        <DeniedText>Access Denied</DeniedText>
        <DeniedSmall>You don't have access to this page</DeniedSmall>

        <Button
          style={{ borderRadius: '6px', height: '45px' }}
          leadingIcon={'arrow_back'}
          text="Go Back"
        />
      </Wrap>
    </Container>
  );
};

export default AdminAccessDenied;
