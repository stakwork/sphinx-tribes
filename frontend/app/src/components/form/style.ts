import styled from 'styled-components';

interface styledProps {
  color?: any;
  show?: boolean;
}
interface WrapProps {
  newDesign?: string | boolean;
}

export const Wrap = styled.div<WrapProps>`
  padding: ${(p: any) => (p?.newDesign ? '28px 0px' : '80px 0px 0px 0px')};
  margin-bottom: ${(p: any) => !p?.newDesign && '100px'};
  display: flex;
  height: inherit;
  flex-direction: column;
  align-content: center;
  min-width: 600px;
`;

export const OrgWrap = styled.div<WrapProps>`
  padding: ${(p: any) => (p?.newDesign ? '28px 0px' : '30px 20px')};
  margin-bottom: ${(p: any) => !p?.newDesign && '100px'};
  display: flex;
  height: inherit;
  flex-direction: column;
  align-content: center;
  min-width: 550px;
  max-width: auto;
`;

interface bottomButtonProps {
  assigneeName?: string;
  color?: any;
  valid?: any;
}

export const BWrap = styled.div<styledProps>`
  display: flex;
  justify-content: space-between !important;
  align-items: center;
  width: 100%;
  padding: 10px;
  min-height: 42px;
  position: absolute;
  left: 0px;
  background: ${(p: any) => p?.color && p.color.pureWhite};
  z-index: 10;
  box-shadow: 0px 1px 6px ${(p: any) => p?.color && p.color.black100};
`;

export const CreateBountyHeaderContainer = styled.div<styledProps>`
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 0px 48px;
  margin-bottom: 30px;
  .TopContainer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    .stepText {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 500;
      font-size: 15px;
      line-height: 18px;
      letter-spacing: 0.06em;
      text-transform: uppercase;
      color: ${(p: any) => p.color && p.color.black500};
      .stepTextSpan {
        font-family: 'Barlow';
        font-style: normal;
        font-weight: 500;
        font-size: 15px;
        line-height: 18px;
        letter-spacing: 0.06em;
        text-transform: uppercase;
        color: ${(p: any) => p.color && p.color.grayish.G300};
      }
    }
    .schemaName {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 500;
      font-size: 13px;
      line-height: 23px;
      display: flex;
      align-items: center;
      text-align: right;
      letter-spacing: 0.06em;
      text-transform: uppercase;
      color: ${(p: any) => p.color && p.color.grayish.G300};
    }
  }

  .HeadingText {
    font-family: 'Barlow';
    font-size: 36px;
    font-weight: 800;
    line-height: 43px;
    color: ${(p: any) => p?.color && p.color.grayish.G10};
    margin-bottom: 11px;
    margin-top: 16px;
  }
`;

export const SchemaTagsContainer = styled.div`
  display: flex;
  justify-content: space-between;
  height: 100%;
  padding: 0px 48px;
  .LeftSchema {
    width: auto;
  }
  .RightSchema {
    width: auto;
  }
`;

export const BottomContainer = styled.div<bottomButtonProps>`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0px 48px;
  .RequiredText {
    font-size: 13px;
    font-family: 'Barlow';
    font-weight: 400;
    line-height: 35px;
    color: ${(p: any) => p?.color && p.color.grayish.G300};
    user-select: none;
  }
  .ButtonContainer {
    display: flex;
    flex-direction: row-reverse;
    justify-content: space-between;
    align-items: center;
  }
  .nextButtonDisable {
    width: 120px;
    height: 42px;
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: pointer;
    background: ${(p: any) => p?.color && p.color.grayish.G950};
    border-radius: 32px;
    user-select: none;
    .disableText {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 600;
      font-size: 14px;
      line-height: 19px;
      display: flex;
      align-items: center;
      text-align: center;
      color: ${(p: any) => p?.color && p.color.grayish.G300};
    }
  }
  .nextButton {
    height: 42px;
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: pointer;
    background: ${(p: any) =>
      p?.assigneeName === '' ? `${p?.color.button_secondary.main}` : `${p?.color.statusAssigned}`};
    box-shadow: 0px 2px 10px
      ${(p: any) =>
        p?.assigneeName === ''
          ? `${p.color.button_secondary.shadow}`
          : `${p.color.button_primary.shadow}`};
    border-radius: 32px;
    color: ${(p: any) => p?.color && p.color.pureWhite};
    :hover {
      background: ${(p: any) =>
        p?.assigneeName === ''
          ? `${p.color.button_secondary.hover}`
          : `${p.color.button_primary.hover}`};
    }
    :active {
      background: ${(p: any) =>
        p?.assigneeName === ''
          ? `${p.color.button_secondary.active}`
          : `${p.color.button_primary.active}`};
    }
    .nextText {
      font-family: 'Barlow';
      font-size: 16px;
      font-weight: 600;
      line-height: 19px;
      user-select: none;
    }
  }
`;

export const SchemaOuterContainer = styled.div`
  display: flex;
  justify-content: center;
  width: 100%;
  .SchemaInnerContainer {
    width: 70%;
    @media only screen and (max-width: 700px) {
      width: auto;
    }
  }
`;

export const AboutSchemaInner = styled.div`
  min-width: 100%;
  @media only screen and (max-width: 700px) {
    padding: 0px 25px;
  }
`;

export const ChooseBountyContainer = styled.div<styledProps>`
  display: flex;
  flex-direction: row;
  height: 100%;
  align-items: center;
  justify-content: center;
  gap: 34px;
  margin-bottom: 24px;
`;

export const BountyContainer = styled.div<styledProps>`
  min-height: 352px;
  max-height: 352px;
  min-width: 290px;
  max-width: 290px;
  background: ${(p: any) => p.color && p.color.pureWhite};
  border: 1px solid ${(p: any) => p.color && p.color.grayish.G600};
  outline: 1px solid ${(p: any) => p.color && p.color.pureWhite};
  box-shadow: 0px 1px 4px ${(p: any) => p.color && p.color.black100};
  border-radius: 20px;
  overflow: hidden;
  transition: all 0.2s;
  .freelancerContainer {
    min-height: 352px;
    max-height: 352px;
    width: 100%;
  }
  :hover {
    border: ${(p: any) =>
      p.show
        ? `1px solid ${p.color.button_primary.shadow}`
        : `1px solid ${(p: any) => p.color && p.color.grayish.G600}`};
    outline: ${(p: any) =>
      p.show
        ? `1px solid ${p.color.button_primary.shadow}`
        : `1px solid ${(p: any) => p.color && p.color.grayish.G600}`};
    box-shadow: ${(p: any) => (p.show ? `1px 1px 6px ${p.color.black85}` : ``)};
  }
  :active {
    border: ${(p: any) =>
      p.show
        ? `1px solid ${p.color.button_primary.shadow}`
        : `1px solid ${(p: any) => p.color && p.color.grayish.G600}`} !important;
  }
  .TextButtonContainer {
    height: 218px;
    width: 290px;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding-top: 60px;
    .textTop {
      height: 40px;
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 700;
      font-size: 20px;
      line-height: 23px;
      display: flex;
      align-items: center;
      text-align: center;
      color: ${(p: any) => p.color && p.color.grayish.G25};
    }
    .textBottom {
      height: 31px;
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 400;
      font-size: 14px;
      line-height: 17px;
      text-align: center;
      color: ${(p: any) => p.color && p.color.grayish.G100};
    }
    .StartButton {
      height: 42px;
      width: 120px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: ${(p: any) => p.color && p.color.button_secondary.main};
      box-shadow: 0px 2px 10px ${(p: any) => p.color && p.color.button_secondary.shadow};
      color: ${(p: any) => p.color && p.color.pureWhite};
      border-radius: 32px;
      margin-top: 10px;
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 600;
      font-size: 14px;
      line-height: 19px;
      cursor: pointer;
      :hover {
        background: ${(p: any) => p.color && p.color.button_secondary.hover};
        box-shadow: 0px 1px 5px ${(p: any) => p.color && p.color.button_secondary.shadow};
      }
      :active {
        background: ${(p: any) => p.color && p.color.button_secondary.active};
      }
      :focus-visible {
        outline: 2px solid ${(p: any) => p.color && p.color.button_primary.shadow} !important;
      }
    }
    .ComingSoonContainer {
      height: 42px;
      margin-top: 10px;
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: center;
      .ComingSoonText {
        font-family: 'Barlow';
        font-style: normal;
        font-weight: 500;
        font-size: 14px;
        line-height: 17px;
        display: flex;
        align-items: center;
        text-align: center;
        letter-spacing: 0.06em;
        text-transform: uppercase;
        color: ${(p: any) => p.color && p.color.grayish.G300};
        margin-right: 18px;
        margin-left: 18px;
      }
    }
  }
`;

export const EditBountyText = styled.h4`
  color: #3c3d3f;
  font-weight: 700;
  @media only screen and (max-width: 500px) {
    font-size: 1.1rem;
  }
`;
