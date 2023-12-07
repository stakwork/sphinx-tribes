import styled from "styled-components"
import { createGlobalStyle } from "styled-components";

const GlobalStyle = createGlobalStyle`
body {
 font-family: 'Your Google Font', sans-serif;
}
`;

export const NavWrapper = styled.div`
display: flex;
padding: 13px 37px 0px 37px;
justify-content: space-between;
align-items: flex-start;
border-bottom: 1px solid var(--Divider-2, #DDE1E5);
background: var(--Body, #FFF);

`
export const AlternateWrapper = styled.div`
background: #FFF;
box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.15);
display: flex;
height: 72px;
padding-left: 1em;
padding-right: 2em;
justify-content: space-between;
align-items: center;
flex-shrink: 0;

`

export const ButtonWrapper = styled.div`
padding-left:1em;
padding-right:1em;
display:flex;
gap:4px;

`

export const Title = styled.h4`
color: var(--Text-2, var(--Hover-Icon-Color, #3C3F41));
font-family: Barlow;
font-size: 24px;
font-style: normal;
font-weight: 900;
line-height: 14px; /* 58.333% */
display:flex;
gap:6px;
`
export const Button = styled.h5`
color: var(--Text-2, var(--Hover-Icon-Color, #3C3F41));
text-align: center;
font-family: Barlow;
font-size: 14px;
font-style: normal;
font-weight: 600;
line-height: normal;
cursor:pointer;
`

export const AlternateTitle = styled.h4`
  color: var(--Text-2, var(--Hover-Icon-Color, #3C3F41));
  leading-trim: both;
  text-edge: cap;
  font-family: Barlow;
  font-size: 24px;
  font-style: normal;
  font-weight: 400;
  line-height: 14px;
`
export const ExportButton =styled.button`
width: 112px;
padding: 8px 16px;
height:40px;
justify-content: center;
align-items: center;
gap: 6px;
border-radius: 6px;
border: 1px solid var(--Input-Outline-1, #D0D5D8);
background: var(--White, #FFF);
box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.06);
margin-right:10px; 
`
export const ExportText =styled.p`
color: var(--Main-bottom-icons, #5F6368);
text-align: center;
leading-trim: both;
text-edge: cap;
font-family: Barlow;
font-size: 14px;
font-style: normal;
font-weight: 500;
line-height: 0px; /* 0% */
letter-spacing: 0.14px;
margin-top:10px;
`

export const Month = styled.h4`
padding-top:18px ;
font-size: 18px;
font-weight:400;
color: var(--Text-2, var(--Hover-Icon-Color, #3C3F41));
leading-trim: both;
text-edge: cap;
font-family: Barlow;
font-size: 20px;
font-style: normal;
font-weight: 500;
line-height: 0px; /* 0% */
letter-spacing: 0.2px;
` 


export const ArrowButton =styled.button`
border-radius: 6px;
border: 1px solid var(--Input-Outline-1, #D0D5D8);
background: var(--White, #FFF);
box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.06);
width:40px;
height:40px;
`
export const DropDown = styled.select`
 display: flex;
  width: 137px;
  height: 40px;
  padding: 8px 8px 8px 16px;
  justify-content: space-between;
  align-items: center;
  border-radius: 6px;
  background: var(--Primary-blue, #618AFF);
  box-shadow: 0px 2px 10px 0px rgba(97, 138, 255, 0.50);
  outline:none;
  border:none;
  color:white;
  font-size:14px;
  
`


export const DropDownOption = styled.option`
color: var(--White, #FFF);
leading-trim: both;
text-edge: cap;
font-family: Barlow;
font-size: 14px;
font-style: normal;
font-weight: 500;
line-height: 0px; /* 0% */
letter-spacing: 0.14px;
margin-right:2px;
`  


export const Flex = styled.div`
display:flex;
`