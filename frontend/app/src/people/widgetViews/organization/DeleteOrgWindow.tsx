import React from 'react';
import styled from 'styled-components';
import { useIsMobile } from 'hooks';
import { Button } from '../../../components/common';

const DeleteWrap = styled.div`
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 20;
  width: 551px;
  height: 435px;
  flex-shrink: 0;
  border-radius: 10px;
  background: rgba(235, 237, 239, 0.85);
`;

const DeleteConfirmation = styled.div`
  width: 353px;
  height: 240px;
  flex-shrink: 0;
  position: relative;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;

  padding: 32px;
  border-radius: 8px;
  border: 1px solid var(--Divider-2, #dde1e5);
  background: var(--White, #fff);
  box-shadow: 0px 4px 20px 0px rgba(0, 0, 0, 0.15);
`;

const DeleteIcon = styled.img`
  width: 27.857px;
  height: 30px;
  flex-shrink: 0;
  fill: var(--Input-Outline-1, #d0d5d8);
`;

const DeleteText = styled.p`
  margin: 5px;

  color: var(--Text-2, var(--Hover-Icon-Color, #3c3f41));
  font-family: 'Barlow';
  font-size: 20px;
  font-style: normal;
  font-weight: 700;
  line-height: 24px;

  color: var(--Text-2, var(--Hover-Icon-Color, #3c3f41));
  text-align: center;
  font-family: 'Barlow';
  font-size: 20px;
  font-style: normal;
  font-weight: 400;
  line-height: 24px; /* 120% */
`;

const DeleteOrgWindow = (props: { onDeleteOrg: () => void; close: () => void }) => {
  const isMobile = useIsMobile();

  return (
    <DeleteWrap
      style={{
        width: isMobile ? '100%' : '551px',
        height: isMobile ? '100%' : '435px'
      }}
    >
      <DeleteConfirmation>
        <DeleteIcon src="/static/Delete.svg" alt="delete icon" />
        <DeleteText style={{ marginTop: '26px' }}>Are you sure you want to</DeleteText>
        <DeleteText style={{ fontWeight: 'bold', marginBottom: '36px' }}>
          Delete this Organization?
        </DeleteText>
        <div style={{ display: 'flex', flexDirection: 'row', width: '100%' }}>
          <Button
            disabled={false}
            onClick={() => props.close()}
            loading={false}
            style={{
              width: '120',
              height: '40px',
              borderRadius: '6px',
              alignSelf: 'flex-start',
              border: '1px solid var(--Input-Outline-1, #D0D5D8)',
              background: 'var(--White, #FFF)',
              boxShadow: '0px 1px 2px 0px rgba(0, 0, 0, 0.06)'
            }}
            color={'#5F6368'}
            text={'Cancel'}
          />
          <Button
            disabled={false}
            onClick={() => props.onDeleteOrg()}
            loading={false}
            style={{
              width: '120',
              height: '40px',
              marginLeft: 'auto',
              borderRadius: '6px',
              background: 'var(--Primary-Red, #ED7474)',
              boxShadow: '0px 2px 10px 0px rgba(237, 116, 116, 0.50)'
            }}
            color={'primary'}
            text={'Delete'}
          />
        </div>
      </DeleteConfirmation>
    </DeleteWrap>
  );
};

export default DeleteOrgWindow;
