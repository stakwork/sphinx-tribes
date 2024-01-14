import styled from 'styled-components';
interface SubheadingProps {
  marginTop?: string;
  marginLeft?: string;
  width?: string;
}

export const Wrapper = styled.section`
  background: var(--Search-bar-background, #f2f3f5);
  padding: 28px 47px;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  flex-shrink: 0;
`;

export const Card = styled.div`
  height: 168px;
  padding: 0px 40px 0px 40px;
  flex-shrink: 0;
  box-shadow: 0px 0px 10px 0px rgba(219, 219, 219, 0.75);
  display: flex;
  flex-direction: column;
  background: var(--Text-Messages, #fff);
  margin-right: 11px;
  margin-left: 11px;
`;

export const CardGreen = styled.div`
  height: 168px;
  padding: 0px 40px 0px 40px;
  flex-shrink: 0;
  box-shadow: 0px 0px 10px 0px rgba(219, 219, 219, 0.75);
  display: flex;
  flex-direction: column;
  background: var(--Text-Messages, #fff);
  margin-right: 11px;
  margin-left: 11px;
`;

export const CardHunter = styled.div`
  height: 168px;
  padding: 0px 40px 0px 40px;
  flex-shrink: 0;
  box-shadow: 0px 0px 10px 0px rgba(219, 219, 219, 0.75);
  display: flex;
  flex-direction: column;
  background: var(--Text-Messages, #fff);
  margin-right: 11px;
  margin-left: 11px;
`;

export const VerticaGrayLine = styled.div`
  width: 2px;
  height: 32px;
  background-color: #edecec;
  content: '';
  margin: 10px auto;
`;

export const VerticaGrayLineSecondary = styled.div`
  width: 2px;
  height: 32px;
  background-color: #edecec;
  content: '';
  margin: 10px 0px 10px 0px;
`;

export const HorizontalGrayLine = styled.div`
  width: 100%;
  height: 2px;
  background-color: #edecec;
  content: '';
  margin: auto;
`;

export const UpperCardWrapper = styled.div`
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 3em;
  padding-top: 15px;
`;

export const BelowCardWrapper = styled.div`
  flex: 1;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 2em;
`;

export const DivWrapper = styled.div`
  flex: 1;
  display: flex;
  gap: 20px;
  margin-top: 1em;
  margin-bottom: 1em;
  position: relative;
  z-index: 0;
`;

export const StatusWrapper = styled.div`
  display: flex;
  gap: 10px;
  margin-top: 1em;
  margin-bottom: 10px;
  position: relative;
  margin-left: auto;
  margin-right: 2em;
  z-index: 0;
`;

export const LeadingText = styled.h2`
  color: var(--Text-2, var(--Hover-Icon-Color, #3c3f41));
  font-family: Barlow;
  font-size: 20px;
  font-style: normal;
  font-weight: 700;
  line-height: normal;
`;

export const Title = styled.div`
  color: var(--Text-2, var(--Hover-Icon-Color, #3c3f41));
  font-family: Barlow;
  font-size: 24px;
  font-style: normal;
  font-weight: 700;
  line-height: normal;
`;

export const TitleBlue = styled.div`
  color: var(--Primary-Blue-Border, #5078f2);
  font-family: Barlow;
  font-size: 24px;
  font-style: normal;
  font-weight: 700;
  line-height: normal;
`;

export const TitleGreen = styled.div`
  color: var(--Green-Border, #2fb379);
  font-family: Barlow;
  font-size: 24px;
  font-style: normal;
  font-weight: 700;
  line-height: normal;
`;

export const Subheading = styled.h3`
  color: var(--Secondary-Text-4, #6b7a8d);
  margin-top: ${(props: SubheadingProps) => (props.marginTop ? props.marginTop : '')};
  margin-left: ${(props: SubheadingProps) => (props.marginLeft ? props.marginLeft : '')};
  width: ${(props: SubheadingProps) => (props.width ? props.width : '')};
  font-family: Barlow;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: normal;
`;

export const TitleWrapper = styled.div`
  width: 109.271px;
  height: 20px;
  flex-shrink: 0;
  display: flex;
  margin-top: 1em;
  margin-bottom: 10px;
`;

export const InfoFrame = styled.div`
  display: inline-flex;
  align-items: flex-start;
  gap: 12px;
  background-color: black;
  width: 172px;
  height: 35px;
`;
export const TitleDiv = styled.div`
  display: flex;
  flex-direction: column;
`;
