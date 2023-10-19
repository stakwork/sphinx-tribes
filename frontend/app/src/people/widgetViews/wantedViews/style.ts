import styled from 'styled-components';
import { Title } from '../../../components/common';

interface WrapProps {
  isClosed?: boolean;
  color?: any;
}

interface styledProps {
  color?: any;
}

export const BountyBox = styled.div<styledProps>`
  min-height: 160px;
  max-height: 160px;
  width: 1100px;
  box-shadow: 0px 1px 6px ${(p: any) => p?.color && p?.color.black100};
  border: none;
`;

export const DWrap = styled.div<WrapProps>`
  display: flex;
  flex: 1;
  height: 100%;
  min-height: 510px;
  flex-direction: column;
  width: 100%;
  min-width: 100%;
  max-height: 510px;
  font-style: normal;
  font-weight: 500;
  font-size: 17px;
  line-height: 23px;
  color: ${(p: any) => p?.color && p?.color.grayish.G10} !important;
  letter-spacing: 0px;
  justify-content: space-between;
  opacity: ${(p: any) => (p.isClosed ? '0.5' : '1')};
  filter: ${(p: any) => (p.isClosed ? 'grayscale(1)' : 'grayscale(0)')};
`;

export const Wrap = styled.div<WrapProps>`
  display: flex;
  justify-content: flex-start;
  opacity: ${(p: any) => (p.isClosed ? '0.5' : '1')};
  filter: ${(p: any) => (p.isClosed ? 'grayscale(1)' : 'grayscale(0)')};
`;

export const B = styled.span<styledProps>`
  font-size: 14px;
  font-weight: bold;
  color: ${(p: any) => p?.color && p?.color.grayish.G10};
`;

export const P = styled.div<styledProps>`
  font-weight: regular;
  font-size: 14px;
  color: ${(p: any) => p?.color && p?.color.grayish.G100};
`;

export const Body = styled.div<styledProps>`
  font-size: 15px;
  line-height: 20px;
  padding: 10px;
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  color: ${(p: any) => p?.color && p?.color.grayish.G05};
  overflow: hidden;
  min-height: 132px;
`;

export const Pad = styled.div`
  display: flex;
  flex-direction: column;
`;

export const DescriptionCodeTask = styled.div<styledProps>`
  margin-bottom: 10px;

  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 13px;
  line-height: 20px;
  color: ${(p: any) => p?.color && p?.color.grayish.G50};
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 6;
  -webkit-box-orient: vertical;
  height: 120px;
  max-height: 120px;
`;

export const DT = styled(Title)`
  margin-bottom: 9px;
  max-height: 52px;
  min-height: 43.5px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  /* Primary Text 1 */

  font-family: 'Roboto';
  font-style: normal;
  font-weight: 500;
  font-size: 17px;
  line-height: 23px;
`;

interface ImageProps {
  readonly src?: string;
}

export const Img = styled.div<ImageProps>`
  background-image: url('${(p: any) => p.src}');
  background-position: center;
  background-size: cover;
  position: relative;
  width: 22px;
  height: 22px;
`;

export const EyeDeleteTextContainerMobile = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
`;

export const EyeDeleteContainerMobile = styled.div`
  margin-top: 10px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
`;
