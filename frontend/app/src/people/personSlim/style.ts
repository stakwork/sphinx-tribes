import styled from 'styled-components';

interface PanelProps {
  isMobile: boolean;
}

export const PeopleList = styled.div`
  position: relative;
  display: flex;
  flex-direction: column;
  background: #ffffff;
  width: 265px;
  overflow-y: overlay !important;

  * {
    scrollbar-width: 6px;
    scrollbar-color: rgba(176, 183, 188, 0.25);
  }

  /* Works on Chrome, Edge, and Safari */
  *::-webkit-scrollbar {
    width: 6px;
  }

  *::-webkit-scrollbar-thumb {
    background-color: rgba(176, 183, 188, 0.25);
    background: rgba(176, 183, 188, 0.25);
    width: 6px;
    border-radius: 10px;
    background-clip: padding-box;
  }

  ::-webkit-scrollbar-track-piece:start {
    background: transparent url('images/backgrounds/scrollbar.png') repeat-y !important;
  }

  ::-webkit-scrollbar-track-piece:end {
    background: transparent url('images/backgrounds/scrollbar.png') repeat-y !important;
  }
`;

export const PeopleScroller = styled.div`
  overflow-y: overlay !important;
  width: 100%;
  height: 100%;
`;

export const AboutWrap = styled.div`
  overflow-y: auto !important;
  ::-webkit-scrollbar-thumb {
    background-color: rgba(176, 183, 188, 0);
    background: rgba(176, 183, 188, 0);
  }

  &:hover {
    ::-webkit-scrollbar-thumb {
      background-color: rgba(176, 183, 188, 0.45);
      background: rgba(176, 183, 188, 0.45);
    }
  }
`;

export const DBack = styled.div`
  min-height: 64px;
  height: 64px;
  display: flex;
  padding-right: 10px;
  align-items: center;
  justify-content: space-between;
  background: #ffffff;
  box-shadow: 0px 1px 6px rgba(0, 0, 0, 0.07);
  z-index: 0;
`;

export const Panel = styled.div<PanelProps>`
  position: relative;
  background: #ffffff;
  color: #000000;
  padding: 20px;
  box-shadow: ${(p) => (p.isMobile ? 'none' : '0px 0px 6px rgb(0 0 0 / 7%)')};
  border-bottom: ${(p) => (p.isMobile ? '2px solid #EBEDEF' : 'none')};
`;

export const Content = styled.div`
  display: flex;
  flex-direction: column;

  width: 100%;
  height: 100%;
  align-items: center;
  color: #000000;
  background: #f0f1f3;
`;

export const Counter = styled.div`
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 11px;
  line-height: 19px;
  margin-bottom: -3px;
  /* or 173% */
  margin-left: 8px;

  display: flex;
  align-items: center;

  /* Placeholder Text */

  color: #b0b7bc;
`;

export const Tabs = styled.div`
  display: flex;
  width: 100%;
  align-items: center;
  // justify-content:center;
  overflow-x: auto;
  ::-webkit-scrollbar {
    display: none;
  }
`;

interface TagProps {
  selected: boolean;
}

export const Tab = styled.div<TagProps>`
  display: flex;
  padding: 10px;
  margin-right: 25px;
  color: ${(p) => (p.selected ? '#292C33' : '#8E969C')};
  border-bottom: ${(p) => p.selected && '4px solid #618AFF'};
  cursor: hover;
  font-weight: 500;
  font-size: 15px;
  line-height: 19px;
  cursor: pointer;
`;

export const Head = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
`;

export const Name = styled.div`
  font-style: normal;
  font-weight: 500;
  font-size: 24px;
  line-height: 28px;
  width: 100%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  /* or 73% */

  text-align: center;

  /* Text 2 */

  color: #3c3f41;
`;

export const Sleeve = styled.div``;

export const RowWrap = styled.div`
  display: flex;
  justify-content: center;

  width: 100%;
`;

interface ImageProps {
  readonly src: string;
}

export const Img = styled.div<ImageProps>`
  background-image: url('${(p) => p.src}');
  background-position: center;
  background-size: cover;
  margin-bottom: 20px;
  width: 150px;
  height: 150px;
  border-radius: 50%;
  position: relative;
  display: flex;
  align-items: flex-end;
  justify-content: flex-end;
`;
