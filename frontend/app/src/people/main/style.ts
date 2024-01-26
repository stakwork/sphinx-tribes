import styled, { css } from 'styled-components';

// this is where we see others posts (etc) and edit our own
export const BWrap = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding: 10px;
  min-height: 42px;
  position: absolute;
  left: 0px;
  border-bottom: 1px solid rgb(221, 225, 229);
  background: #ffffff;
  box-shadow: 0px 1px 6px rgba(0, 0, 0, 0.07);
  z-index: 100;
`;

export const EnvWithScrollBar = ({ thumbColor, trackBackgroundColor }: any) => css`
                scrollbar-color: ${thumbColor} ${trackBackgroundColor}; // Firefox support
                scrollbar-width: thin;

                &::-webkit-scrollbar {
                    width: 6px;
                height: 100%;
  }

                &::-webkit-scrollbar-thumb {
                    background - color: ${thumbColor};
                background-clip: content-box;
                border-radius: 5px;
                border: 1px solid ${trackBackgroundColor};
  }

                &::-webkit-scrollbar-corner,
                &::-webkit-scrollbar-track {
                    background - color: ${trackBackgroundColor};
  }
}

                `;
interface BProps {
  hide: boolean;
}

export const B = styled.div<BProps>`
  display: ${(p: any) => (p.hide ? 'none' : 'flex')};
  justify-content: ${(p: any) => (p.hide ? 'none' : 'center')};
  height: 100%;
  width: 100%;
  overflow-y: auto;
  box-sizing: border-box;
  ${EnvWithScrollBar({
    thumbColor: '#5a606c',
    trackBackgroundColor: 'rgba(0,0,0,0)'
  })}
`;
