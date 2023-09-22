import styled from 'styled-components';

interface styledColor {
  color?: any;
}

interface ButtonContainerProps extends styledColor {
  topMargin?: string;
}
export const ButtonContainer = styled.div<ButtonContainerProps>`
  width: 220px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  margin-top: ${(p: any) => p?.topMargin};
  background: ${(p: any) => p?.color && p?.color.pureWhite};
  border: 1px solid ${(p: any) => p?.color && p?.color.grayish.G600};
  border-radius: 30px;
  user-select: none;
  .LeadingImageContainer {
    margin-left: 14px;
    margin-right: 16px;
  }
  .ImageContainer {
    position: absolute;
    min-height: 48px;
    min-width: 48px;
    right: 37px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .buttonImage {
    filter: brightness(0) saturate(100%) invert(85%) sepia(10%) saturate(180%) hue-rotate(162deg)
      brightness(87%) contrast(83%);
  }
  :hover {
    border: 1px solid ${(p: any) => p?.color && p?.color.grayish.G300};
  }
  :active {
    border: 1px solid ${(p: any) => p?.color && p?.color.grayish.G100};
    .buttonImage {
      filter: brightness(0) saturate(100%) invert(22%) sepia(5%) saturate(563%) hue-rotate(161deg)
        brightness(91%) contrast(86%);
    }
  }
  .ButtonText {
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 500;
    font-size: 14px;
    line-height: 17px;
    color: ${(p: any) => p?.color && p?.color.grayish.G50};
  }
`;
