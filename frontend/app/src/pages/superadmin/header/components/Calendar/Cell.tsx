import { strict } from "assert";
import React from "react";
import styled, { css } from "styled-components";
import { string } from "yup";

type Props={
  className?: string;
  isActive?: boolean;
  onClick?: () => void;
  children?:number | string;
}

const CellWrapper = styled.div<Props>`
  height: 3.2em;
  width: 3.2em;
  display: flex;
  align-items: center;
  justify-content: center;
  user-select: none;
  border-radius: 50%;
  color: rgba(0, 0, 0, 0.7);
  transition: color 150ms, background-color 150ms, border-color 150ms,
    text-decoration-color 150ms, fill 150ms, stroke 150ms;
  cursor: pointer;
  &:hover {
    background-color: ${(props:any) =>
      props.isActive ? "" : "rgb(243, 244, 246)"};
  }
  ${(props:any) =>
      props.isActive &&
    css`
      font-weight: 700;
      color: white;
      background-color: #618AFF;
    `}
`;

const Cell = ({
  onClick,
  children,
  isActive = false,
}:Props) => {
  return (
    <CellWrapper onClick={!isActive ? onClick : undefined} isActive={isActive}>
      {children}
    </CellWrapper>
  );
};

export default Cell;
