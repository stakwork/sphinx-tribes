import { strict } from 'assert';
import React from 'react';
import styled, { css } from 'styled-components';

type Props = {
  className?: string;
  isActive?: boolean;
  onClick?: () => void;
  children?: number | string;
  isEmpty: boolean;
};

const CellWrapper = styled.div<Props>`
  height: 48px;
  width: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  user-select: none;
  border-radius: 50%;
  color: rgba(0, 0, 0, 0.7);
  transition:
    color 150ms,
    background-color 150ms,
    border-color 150ms,
    text-decoration-color 150ms,
    fill 150ms,
    stroke 150ms;
  cursor: pointer;
  &:hover {
    background-color: ${(props: any) =>
      props.isEmpty || props.isActive ? '' : 'rgb(243, 244, 246)'};
  }
  ${(props: any) =>
    props.isActive &&
    css`
      font-weight: 700;
      color: white;
      background-color: #618aff;
    `}
`;
const Cell = ({ onClick, children, isActive = false, isEmpty }: Props) => (
    <CellWrapper onClick={!isActive ? onClick : undefined} isActive={isActive} isEmpty={isEmpty}>
      {children}
    </CellWrapper>
  );

export default Cell;
