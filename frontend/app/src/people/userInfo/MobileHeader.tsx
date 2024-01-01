import { Button, IconButton } from 'components/common';
import { PeopleMobileeHeaderProps } from 'people/interfaces';
import React from 'react';

export const HeaderMobile = ({ goBack, canEdit, onEdit }: PeopleMobileeHeaderProps) => (
  <div
    style={{
      position: 'absolute',
      top: 20,
      left: 0,
      display: 'flex',
      justifyContent: 'space-between',
      width: '100%',
      padding: '0 20px'
    }}
  >
    <IconButton onClick={goBack} icon="arrow_back" />
    {canEdit ? (
      <>
        <Button
          text="Edit Profile"
          onClick={onEdit}
          color="white"
          height={42}
          style={{
            fontSize: 13,
            color: '#3c3f41',
            border: 'none',
            marginLeft: 'auto'
          }}
          leadingIcon={'edit'}
          iconSize={15}
        />
      </>
    ) : (
      <div />
    )}
  </div>
);
