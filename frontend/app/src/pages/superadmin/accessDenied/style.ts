import styled from 'styled-components';

export const Container = styled.div`
  width: 100vw;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 0px;
`;

export const Wrap = styled.div`
  padding: 0px;
  margin: 0px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin-top: -10%;
  @media only screen and (max-width: 500px) {
    margin-top: -30%;
  }
`;

export const AccessImg = styled.img`
  width: 150px;
  height: 150px;
  padding: 0px;
  @media only screen and (max-width: 800px) {
    width: 120px;
    height: 120px;
  }
  @media only screen and (max-width: 500px) {
    width: 100px;
    height: 100px;
  }
`;

export const DeniedText = styled.h1`
  font-weight: bolder;
  margin-top: 10px;
  margin-bottom: 0px;
  @media only screen and (max-width: 800px) {
    font-size: 1.4rem;
  }
  @media only screen and (max-width: 500px) {
    font-size: 1.3rem;
  }
`;

export const DeniedSmall = styled.p`
  font-weight: 500;
  font-size: 0.9rem;
  color: #5f6368;
  margin-top: 10px;
  margin-bottom: 20px;
  @media only screen and (max-width: 800px) {
    font-size: 0.88rem;
  }
  @media only screen and (max-width: 500px) {
    font-size: 0.83rem;
  }
`;
