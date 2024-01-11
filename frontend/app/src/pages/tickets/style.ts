import styled from 'styled-components';

export const Body = styled.div`
  flex: 1;
  height: calc(100% - 105px);
  width: 100%;
  overflow: auto;
  display: flex;
  flex-direction: column;
`;

export const OrgBody = styled.div`
  display:flex;
  flex-direction:column;
  background: var(--Search-bar-background, #F2F3F5);
  height: 100vh;
`

export const Backdrop = styled.div`
  position: fixed;
  z-index: 1;
  background: rgba(0, 0, 0, 70%);
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
`;

export const Spacer = styled.div`
  display: flex;
  min-height: 10px;
  min-width: 100%;
  height: 10px;
  width: 100%;
`;
