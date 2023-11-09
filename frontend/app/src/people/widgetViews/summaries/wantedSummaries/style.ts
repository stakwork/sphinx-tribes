import styled from 'styled-components';

interface colorProps {
  color?: any;
  isPaidStatusPopOver?: any;
  isPaidStatusBadgeInfo?: any;
}

interface styleProps extends colorProps {
  paid?: boolean;
}

export const Wrap = styled.div<colorProps>`
  display: flex;
  width: 100%;
  height: 100%;
  min-width: 800px;
  font-style: normal;
  font-weight: 500;
  font-size: 24px;
  color: ${(p: any) => p?.color && p.color.grayish.G10};
  justify-content: space-between;
`;

export const SectionPad = styled.div`
  padding: 38px;
  word-break: break-word;
`;

export const Pad = styled.div`
  padding: 0 20px;
  word-break: break-word;
`;

export const Y = styled.div`
  display: flex;
  justify-content: space-between;
  width: 100%;
  padding: 20px 0;
  align-items: center;
`;

export const T = styled.div`
  font-weight: 500;
  font-size: 20px;
  margin: 10px 0;
`;

export const B = styled.span<colorProps>`
  font-size: 15px;
  font-weight: bold;
  color: ${(p: any) => p?.color && p.color.grayish.G10};
`;

export const P = styled.div<colorProps>`
  font-weight: regular;
  font-size: 15px;
  color: ${(p: any) => p?.color && p.color.grayish.G100};
`;

export const D = styled.div<colorProps>`
  color: ${(p: any) => p?.color && p.color.grayish.G50};
  margin: 10px 0 30px;
`;

export const Assignee = styled.div<colorProps>`
  margin-left: 3px;
  font-weight: 500;
  cursor: pointer;

  &:hover {
    color: ${(p: any) => p?.color && p.color.pureBlack};
  }
`;

export const ButtonRow = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
`;

export const GithubIcon = styled.div`
  height: 20px;
  width: 20px;
  position: relative;
  top: -6px;
  margin-left: 20px;
  cursor: pointer;
`;

export const LoomIcon = styled.div`
  height: 20px;
  width: 20px;
  position: relative;
  top: -6px;
  margin-left: 20px;
  cursor: pointer;
`;

export const GithubIconMobile = styled.div`
  height: 20px;
  width: 20px;
  position: relative;
  margin-left: 20px;
  cursor: pointer;
`;

export const LoomIconMobile = styled.div`
  height: 20px;
  width: 20px;
  position: relative;
  margin-left: 20px;
  cursor: pointer;
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

export const Creator = styled.div`
  min-width: 892px;
  max-width: 892px;
  height: 100vh;
  display: flex;
  justify-content: space-between;
`;

export const NormalUser = styled.div`
  min-width: 892px;
  max-width: 892px;
  height: 100%;
  overflow-y: auto;
  overflow-x: hidden;
  display: flex;
  justify-content: space-between;
`;

export const CreatorDescription = styled.div<styleProps>`
  min-width: 600px;
  max-width: 600px;
  overflow: auto;
  height: 100vh;
  border-right: ${(p: any) =>
    p?.paid ? `3px solid ${p?.color?.primaryColor.P400}` : `1px solid ${p?.color.grayish.G700}`};
  background: ${(p: any) => p?.color && p.color.pureWhite};
  padding: 48px 0px 0px 48px;
  .DescriptionUpperContainerNormalView {
    padding-right: 28px;
  }
  .CreatorDescriptionOuterContainerCreatorView {
    padding-right: 28px;
  }
  .CreatorDescriptionInnerContainerCreatorView {
    display: flex;
    justify-content: space-between;
    .CreatorDescriptionExtraButton {
      min-width: 250px;
      max-width: 250px;
      min-height: 40px;
      max-height: 40px;
      display: flex;
      justify-content: space-between;
    }
  }
`;

export const TitleBox = styled.div<colorProps>`
  margin-top: 24px;
  font-family: 'Barlow';
  font-style: normal;
  font-weight: 600;
  font-size: 22px;
  line-height: 26px;
  display: flex;
  align-items: center;
  color: ${(p: any) => p?.color && p.color.text1};
`;

export const DescriptionBox = styled.div<colorProps>`
  padding-right: 44px;
  margin-right: 5px;
  max-height: calc(100% - 160px);
  overflow-y: scroll;
  overflow-wrap: anywhere;
  font-family: 'Barlow';
  font-weight: 400;
  font-size: 15px;
  line-height: 25px;
  color: ${(p: any) => p?.color && p.color.black500};
  .loomContainer {
    margin-top: 23px;
  }
  .loomHeading {
    font-family: 'Barlow';
    font-style: normal;
    font-size: 13px;
    font-weight: 700;
    line-height: 25px;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: ${(p: any) => p?.color && p.color.black500};
  }
  .deliverablesContainer {
    margin-top: 23px;
    .deliverablesHeading {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 700;
      font-size: 13px;
      line-height: 25px;
      letter-spacing: 0.06em;
      text-transform: uppercase;
      color: ${(p: any) => p?.color && p.color.black500};
    }
    .deliverablesDesc {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 400;
      font-size: 15px;
      line-height: 20px;
      color: ${(p: any) => p?.color && p.color.black500};
      white-space: pre-wrap;
    }
  }
  ::-webkit-scrollbar {
    width: 6px;
    height: 6px;
  }
  ::-webkit-scrollbar-thumb {
    background: ${(p: any) => p?.color && p.color.grayish.G300} !important;
    height: 80px;
  }
  ::-webkit-scrollbar-track-piece {
    height: 80px;
  }
`;

export const AssigneeProfile = styled.div<colorProps>`
  min-width: 292px;
  max-width: 292px;
  max-height: 100vh;
  background: ${(p: any) => p?.color && p.color.pureWhite};
  display: flex;
  flex-direction: column;
  .buttonSet {
    display: flex;
    flex-direction: column;
    flex: 1;
  }
`;

interface BountyPriceContainerProps {
  margin_top?: string;
}

export const BountyPriceContainer = styled.div<BountyPriceContainerProps>`
  padding-left: 37px;
  margin-top: ${(p: any) => p.margin_top};
`;

interface codingLangProps {
  background?: string;
  border?: string;
  color?: string;
  styledColors?: any;
}

export const LanguageContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  width: 80%;
  margin-top: 16px;
  margin-bottom: 23.25px;
`;

export const CodingLabels = styled.div<codingLangProps>`
  padding: 0px 8px;
  border: ${(p: any) => (p.border ? p?.border : `1px solid ${p?.styledColors.pureBlack}`)};
  color: ${(p: any) => (p.color ? p?.color : `${p?.styledColors.pureBlack}`)};
  background: ${(p: any) => (p.background ? p?.background : `${p?.styledColors.pureWhite}`)};
  border-radius: 4px;
  overflow: hidden;
  max-height: 22.75px;
  min-height: 22.75px;
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-right: 4px;
  .LanguageText {
    font-size: 13px;
    font-weight: 500;
    text-align: center;
    font-family: 'Barlow';
    line-height: 16px;
  }
`;

export const DividerContainer = styled.div`
  padding: 32px 36.5px;
`;

interface containerProps {
  color?: any;
  unAssignedBackgroundImage?: string;
  assignedBackgroundImage?: string;
  unassigned_border?: string;
  grayish_G200?: string;
}

export const UnassignedPersonProfile = styled.div<containerProps>`
  min-width: 228px;
  min-height: 57.6px;
  display: flex;
  padding-top: 0px;
  padding-left: 28px;
  margin-top: 43px;
  .UnassignedPersonContainer {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 57.6px;
    width: 57.6px;
    border-radius: 50%;
  }
  .UnassignedPersonalDetailContainer {
    margin-left: 25px;
    display: flex;
    align-items: center;
  }
  .BountyProfileOuterContainerCreatorView {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
  }
  .AssigneeCloseButtonContainer {
    margin-left: 6px;
    margin-top: 5px;
    align-self: center;
    height: 22px;
    width: 22px;
    cursor: pointer;
  }
`;

export const AutoCompleteContainer = styled.div<colorProps>`
  overflow: hidden;
  z-index: 10;
  padding: 25px 53px 6px 53px;
  .autoCompleteHeaderText {
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 800;
    font-size: 26px;
    line-height: 36px;
    color: ${(p: any) => p.color && p.color.text2};
    height: 44px;
    margin-bottom: 11px;
  }
`;

export const BottomButtonContainer = styled.div`
  margin-bottom: 20px;
`;

export const AdjustAmountContainer = styled.div<colorProps>`
  min-height: 460px;
  max-height: 460px;
  min-width: 440px;
  max-width: 440px;
  border-radius: 10px;
  background: ${(p: any) => p.color && p.color.pureWhite};
  padding-top: 32px;
  .TopHeader {
    max-height: 48px;
    height: 100%;
    display: flex;
    align-items: center;
    margin-left: 25px;
    cursor: pointer;
    .imageContainer {
      height: 48px;
      width: 48px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    .TopHeaderText {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 500;
      font-size: 15px;
      line-height: 18px;
      letter-spacing: 0.06em;
      text-transform: uppercase;
      color: ${(p: any) => p.color && p.color.black500};
    }
  }
  .Header {
    height: 32px;
    margin-left: 70px;
    .HeaderText {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 800;
      font-size: 36px;
      line-height: 43px;
      display: flex;
      align-items: center;
      text-align: center;
      color: ${(p: any) => p.color && p.color.black500};
    }
  }
  .AssignedProfile {
    height: 184px;
    margin-top: 30px;
    padding: 0px 31px 0px 38px;
    .InputContainer {
      display: flex;
      flex-direction: row;
      align-items: center;
      .InputContainerLeadingText {
        font-family: 'Barlow';
        font-style: normal;
        font-weight: 400;
        font-size: 14px;
        line-height: 17px;
        display: flex;
        align-items: center;
        color: ${(p: any) => (p.color ? p.color.grayish.G100 : '')};
        margin-right: 7px;
      }
      .InputContainerTextField {
        width: 296px;
        background: ${(p: any) => p?.color && p?.color?.pureWhite};
        border: 1px solid ${(p: any) => p.color && p.color.grayish.G600};
        color: ${(p: any) => p.color && p.color.pureBlack};
      }
      .InputContainerEndingText {
        font-family: 'Barlow';
        font-style: normal;
        font-weight: 400;
        font-size: 14px;
        line-height: 17px;
        display: flex;
        align-items: center;
        color: ${(p: any) => p.color && p.color.grayish.G100};
        margin-left: 14px;
      }
    }
    .USDText {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 500;
      font-size: 13px;
      line-height: 16px !important;
      display: flex;
      align-items: center;
      color: ${(p: any) => p.color && p.color.grayish.G100};
      margin-left: 42px;
      height: 32px;
    }
  }
  .BottomButton {
    margin-top: 20px;
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    padding-right: 36px;
  }
`;

export const AwardsContainer = styled.div<colorProps>`
  width: 622px;
  height: 100vh;
  max-height: 100vh;
  background: ${(p: any) => p.color && p.color.pureWhite};
  display: flex;
  flex-direction: column;
  .header {
    min-height: 159px;
    max-height: 159px;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    border-bottom: 1px solid ${(p: any) => p.color && p.color.grayish.G600};
    box-shadow: 0px 1px 4px ${(p: any) => p.color && p.color.black80};
    .headerTop {
      height: 48px;
      display: flex;
      align-items: center;
      margin: 32px 0px 0px 25px;
      cursor: pointer;
      .imageContainer {
        height: 48px;
        width: 48px;
        display: flex;
        justify-content: center;
        align-items: center;
      }
      .TopHeaderText {
        font-family: 'Barlow';
        font-style: normal;
        font-weight: 500;
        font-size: 15px;
        line-height: 18px;
        letter-spacing: 0.06em;
        text-transform: uppercase;
        color: ${(p: any) => p.color && p.color.black500};
      }
    }
    .headerText {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 800;
      font-size: 36px;
      line-height: 43px;
      display: flex;
      align-items: center;
      color: ${(p: any) => p.color && p.color.black500};
      margin-left: 73px;
      margin-bottom: 48px;
    }
  }
  .AwardContainer {
    min-height: 481px;
    height: 100%;
    display: grid;
    grid-template-columns: 0.5fr 0.5fr;
    gap: 10px;
    place-content: flex-start;
    overflow-y: scroll;
    margin-left: 63px;
    user-select: none;
    cursor: pointer;
    .RadioImageContainer {
      display: flex;
      flex-direction: row;
      height: 65px;
      width: 248px;
      align-items: center;
      padding-left: 9px;
      margin-top: 14px;
      border-radius: 6px;
      input[type='radio'] {
        border: 1px solid ${(p: any) => p.color && p.color.grayish.G500};
        border-radius: 2px;
        -webkit-appearance: none;
      }
      input[type='radio']:checked {
        background: url('/static/Checked.svg');
        background-repeat: no-repeat;
        border-radius: 2px;
        border: none;
      }
    }
    .awardImageContainer {
      height: 40px;
      width: 40px;
      margin-left: 13px;
    }
    .awardLabelText {
      margin-left: 15px;
      font-family: 'Barlow';
      font-weight: 500;
      font-size: 13px;
      line-height: 15px;
      color: ${(p: any) => p.color && p.color.grayish.G05};
    }
  }
`;

export const AwardBottomContainer = styled.div<colorProps>`
  height: 129px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-top: 1px solid ${(p: any) => p.color && p.color.grayish.G600};
  box-shadow: 0px -1px 4px ${(p: any) => p.color && p.color.black80};
`;

export const PaidStatusPopover = styled.div<colorProps>`
  position: absolute;
  background: transparent;
  height: 70px;
  width: 222px;
  right: 54px;
  top: 120px;
  background-image: url('/static/paid_popover_triangle.svg');
  background-size: 16px 16px;
  background-repeat: no-repeat;
  background-position: 16% 0%;
  filter: drop-shadow(0px 1px 20px rgba(0, 0, 0, 0.15));

  .PaidStatusContainer {
    height: 65px;
    width: 222px;
    background: ${(p: any) => p.color && p.color.green1};
    margin-top: 5px;
    padding: 18px 0px 0px 21px;
    display: flex;
    flex-direction: row;
    cursor: pointer;
    .imageContainer {
      width: 31px;
      height: 31px;
    }
    .PaidStatus {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 700;
      font-size: 17px;
      line-height: 15px;
      color: ${(p: any) => p.color && p.color.pureWhite};
      margin-top: 6px;
      margin-left: 18px;
      user-select: none;
    }
  }
  .ExtraBadgeInfo {
    display: flex;
    flex-direction: row;
    align-items: center;
    background: ${(p: any) => p.color && p.color.black400};
    height: 75px;
    width: 222px;
    padding: 14px 0px 0px 19px;
    object-fit: cover;
    border-radius: 0px 0px 6px 6px;
    opacity: ${(p: any) => (p?.isPaidStatusBadgeInfo ? 1 : 0)};
    transition: all ease 4s;
    .imageContainer {
      position: absolute;
      top: 96px;
      left: 14px;
      height: 15px;
      width: 15px;
      background: ${(p: any) => p.color && p.color.pureWhite};
      display: flex;
      justify-content: center;
      align-items: center;
      border-radius: 50%;
      border: none;
    }
    .badgeText {
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 700;
      font-size: 17px;
      line-height: 15px;
      display: flex;
      align-items: center;
      color: ${(p: any) => p.color && p.color.pureWhite};
      margin-left: 11px;
    }
  }
`;

export const CountDownTimerWrap = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin-bottom: 5px;
  width: 220px;
`;

export const CountDownText = styled.p`
  font-size: 1rem;
  margin: 0;
  text-align: center;
`;

export const CountDownTimer = styled.p`
  font-size: 2rem;
  padding: 0px;
  font-weight: bolder;
`;

export const InvoiceWrap = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 220px;
`;

export const QrWrap = styled.div`
  overflow: hidden;
  text-align: center;
  width: 100%;
`;

export const BountyTime = styled.p`
  font-size: 0.9rem;
  padding: 5px 0px;
  text-align: center;
  width: 220px;
  margin: 15px 0px;
`;
