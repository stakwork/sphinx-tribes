import styled from 'styled-components';

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

export const Head = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
`;
export const RowWrap = styled.div`
  display: flex;
  justify-content: center;

  width: 100%;
`;

export const Name = styled.div`
  font-style: normal;
  font-weight: 500;
  font-size: 24px;
  line-height: 28px;
  text-align: center;
  color: #3c3f41;
  width: 100%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
`;

interface ImageProps {
  readonly src: string;
}

export const Img = styled.div<ImageProps>`
  background-image: url('${(p: any) => p.src}');
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
