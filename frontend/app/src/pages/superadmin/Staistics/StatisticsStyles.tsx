
import styled from "styled-components";

 export const Wrapper = styled.section`
  padding: 2em;
  background: var(--Search-bar-background, #F2F3F5);
`;

 export const Card = styled.div`
  background-color:white;
  margin-top: 2em;
  box-shadow: 0px 0px 10px 0px rgba(219, 219, 219, 0.75);
  padding: 1em;
  position: relative;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1em;
`;

export const VerticaBluelLine = styled.div`
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 6px;
  background-color: #007bff; 
  content: ""; 
`;


export const VerticaGreenlLine = styled.div`
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 6px;
  background-color: #49c998; 
  content: "";
`;

export const VerticaGrayLine = styled.div`
  position: absolute;
  left: -25px;
  top: -10px;
  bottom: 0;
  width: 2px;
  height:73px;
  background-color: #edecec; 
  content: ""; 
`;


export const DivWrapper = styled.div`
 position:relative;
 display:flex;

  gap:10px;
`

export const LeadingText = styled.h2`
  color: var(--Text-2, var(--Hover-Icon-Color, #3C3F41));
  font-family: Barlow;
  font-size: 20px;
  font-style: normal;
  font-weight: 700;
  line-height: normal;
`;

export const Title = styled.h2`
  color: var(--Text-2, var(--Hover-Icon-Color, #3C3F41));
  font-family: Barlow;
  font-size: 24px;
  font-style: normal;
  font-weight: 700;
  line-height: normal;
`;

export const TitleBlue = styled.h2`
  color: var(--Primary-Blue-Border, #5078F2);
  font-family: Barlow;
  font-size: 24px;
  font-style: normal;
  font-weight: 700;
  line-height: normal;
`;

export const TitleGreen = styled.h2`
  color: var(--Green-Border, #2FB379);
  font-family: Barlow;
  font-size: 24px;
  font-style: normal;
  font-weight: 700;
  line-height: normal;
`;

export const Subheading = styled.h3`
  font-size: 0.75em;
  text-align: left;
  color:#bfc5ce;
  font-weight:bold;
`;

export const TitleWrapper = styled.div`
  padding-top:0.75em;
  padding-left:5.5em;
   display:flex;
`;