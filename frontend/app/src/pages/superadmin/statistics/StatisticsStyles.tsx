import styled from 'styled-components';

export const Wrapper = styled.section`
  background: var(--Search-bar-background, #f2f3f5);
  padding: 28px 0px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  flex-shrink: 0;
`;

export const Card = styled.div`
  height: 105px;
  flex-shrink: 0;
  box-shadow: 0px 0px 10px 0px rgba(219, 219, 219, 0.75);
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1em;
  border-left: 6px solid var(--Primary-blue, #618aff);
  background: var(--Text-Messages, #fff);
  margin-right: 47px;
  margin-left: 47px;
`;

export const CardGreen = styled.div`
  height: 105px;
  flex-shrink: 0;
  box-shadow: 0px 0px 10px 0px rgba(219, 219, 219, 0.75);
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1em;
  border-left: 6px solid var(--Primary-Green, #49c998);
  background: var(--Text-Messages, #fff);
  margin-right: 47px;
  margin-left: 47px;
`;

export const VerticaGrayLine = styled.div`
  position: absolute;
  left: -45px;
  top: -16.42px;
  bottom: 17.42px;
  width: 2px;
  height: 73px;
  background-color: #edecec;
  content: '';
  margin-bottom: 17.42px;
`;

export const VerticaGrayLineSecondary = styled.div`
  position: absolute;
  left: -75px;
  top: -16.42px;
  bottom: 17.42px;
  width: 2px;
  height: 73px;
  background-color: #edecec;
  content: '';
  margin-bottom: 17.42px;
`;

export const VerticaGrayLineAleternative = styled.div`
  position: absolute;
  left: -25px;
  top: -16.42px;
  bottom: 17.42px;
  width: 2px;
  height: 73px;
  background-color: #edecec;
  content: '';
  margin-bottom: 17.42px;
`;

export const DivWrapper = styled.div`
  display: flex;
  gap: 10px;
  margin-top: 30px;
  margin-bottom: 45px;
  position: relative;
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
  margin-top: 40.5px;
  margin-bottom: 42.5px;
  margin-left: 46.99px;
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
