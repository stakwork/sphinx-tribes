import styled from 'styled-components';
import './main.css';

export const FlexContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  flex-direction: column;
`;
export const TextWrapper = styled.div`
  display: flex;
  gap: 6px;
`;

export const LeadingText = styled.h1`
  color: var(--Text-2, var(--Hover-Icon-Color, #3c3f41));
  font-size: 36px;
  font-style: normal;
  font-weight: 900;
  line-height: 14px; /* 38.889% */
  font-family: 'Barlow', sans-serif;
`;
export const SecondaryText = styled.h2`
  color: var(--Text-2, var(--Hover-Icon-Color, #3c3f41));
  font-family: 'Barlow', sans-serif;
  font-size: 36px;
  font-style: normal;
  font-weight: 400;
  line-height: 14px;
`;
export const InfoText = styled.h4`
  color: var(--Main-bottom-icons, #5f6368);
  margin-top: 43.5px;
  margin-bottom: 27px;
  text-align: center;
  font-family: 'Barlow', sans-serif;
  font-size: 20px;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
`;
export const ImageContainer = styled.div`
  position: relative;
  margin-top: 27px;
`;

export const LogoSpan = styled.span`
  position: absolute;
  top: 52.65px;
  left: 53.06px;
`;

export const ButtonDiv = styled.button`
  margin: 30px;
  width: 197px;
  height: 48px;
  border-radius: 30px;
  background: var(--Primary-blue, #618aff);
  box-shadow: 0px 2px 10px 0px rgba(97, 138, 255, 0.5);
  outline: none;
  border: none;
  text-align: center;
  color: white;
  font-size: 16px;
`;
