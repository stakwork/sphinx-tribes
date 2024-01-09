import styled from 'styled-components';

export const Container = styled.div``;

export const NavWrapper = styled.div`
  display: flex;
  padding: 23px 47px 0px 47px;
  justify-content: space-between;
  align-items: flex-start;
  border-bottom: 1px solid var(--Divider-2, #dde1e5);
  background: var(--Body, #fff);
`;
export const AlternateWrapper = styled.div`
  background: #fff;
  box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.15);
  display: flex;
  height: 72px;
  padding: 0 47px;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
  position: fixed;
  top: 62px;
  left: 0;
  width: 100%;
  z-index: 9999999;
`;
export const LeftWrapper = styled.div`
  display: flex;
  align-items: flex-start;
  gap: 20px;
`;

export const ButtonWrapper = styled.div`
  display: flex;
  align-items: flex-start;
  gap: 8px;
`;
export const RightWrapper = styled.div`
  display: flex;
  justify-content: center;
  align-items: flex-start;
  gap: 16px;
`;

export const Title = styled.h4`
  color: var(--Text-2, var(--Hover-Icon-Color, #3c3f41));
  font-family: 'Barlow', sans-serif;
  font-size: 24px;
  font-style: normal;
  font-weight: 900;
  line-height: 14px; /* 58.333% */
  display: flex;
  gap: 6px;
`;
export const Button = styled.h5`
  color: var(--Text-2, var(--Hover-Icon-Color, #3c3f41));
  text-align: center;
  font-family: 'Barlow', sans-serif;
  font-size: 14px;
  font-style: normal;
  font-weight: 600;
  line-height: normal;
  cursor: pointer;
`;

export const AlternateTitle = styled.h4`
  color: var(--Text-2, var(--Hover-Icon-Color, #3c3f41));
  font-family: 'Barlow', sans-serif;
  font-size: 24px;
  font-style: normal;
  font-weight: 400;
  line-height: 14px;
`;
export const ExportButton = styled.button`
  width: 112px;
  padding: 8px 16px;
  height: 40px;
  justify-content: center;
  align-items: center;
  gap: 6px;
  border-radius: 6px;
  border: 1px solid var(--Input-Outline-1, #d0d5d8);
  background: var(--White, #fff);
  box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.06);
  margin-right: 10px;
`;
export const ExportText = styled.p`
  color: var(--Main-bottom-icons, #5f6368);
  text-align: center;
  font-family: 'Barlow', sans-serif;
  font-size: 14px;
  font-style: normal;
  font-weight: 500;
  line-height: 0px; /* 0% */
  letter-spacing: 0.14px;
  margin-top: 10px;
`;

export const Month = styled.h4`
  padding-top: 18px;
  font-size: 18px;
  font-weight: 400;
  color: var(--Text-2, var(--Hover-Icon-Color, #3c3f41));
  font-family: 'Barlow', sans-serif;
  font-size: 20px;
  font-style: normal;
  font-weight: 500;
  line-height: 0px; /* 0% */
  letter-spacing: 0.2px;
`;

export const ArrowButton = styled.button`
  border-radius: 6px;
  border: 1px solid var(--Input-Outline-1, #d0d5d8);
  background: var(--White, #fff);
  box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.06);
  width: 40px;
  height: 40px;
`;
export const DropDown = styled.div`
  display: flex;
  width: 137px;
  height: 40px;
  padding: 8px 8px 8px 16px;
  justify-content: space-between;
  align-items: center;
  border-radius: 6px;
  background: var(--Primary-blue, #618aff);
  box-shadow: 0px 2px 10px 0px rgba(97, 138, 255, 0.5);
  outline: none;
  border: none;
  color: white;
  font-size: 14px;
`;
export const Select = styled.select`
  background: var(--Primary-blue, #618aff);
  box-shadow: 0px 2px 10px 0px rgba(97, 138, 255, 0.5);
  outline: none;
  border: none;
  width: 113px;
  color: var(--White, #fff);
`;

export const Option = styled.div`
  position: absolute;
  z-index: 1;
  top: 130px;
  right: 48px;
  width: 169px;
  height: 157px;
  display: inline-flex;
  padding: 12px 28px 12px 28px;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  border-radius: 6px;
  background: #fff;
  box-shadow: 0px 4px 20px 0px rgba(0, 0, 0, 0.25);

  ul {
    list-style: none;
    padding: 0;
    margin: 0;
    color: grey;
  }

  li {
    padding: 4px;
    cursor: pointer;
    color: grey;
    font-family: 'Barlow', sans-serif;
    font-size: 15px;
    font-style: normal;
    font-weight: 500;
    line-height: 18px;

    &:hover {
      color: #3c3f41;
    }
  }
`;

export const CustomButton = styled.button`
  display: flex;
  width: 113px;
  height: 40px;
  padding: 8px 8px 8px 16px;
  justify-content: center;
  align-items: center;
  gap: 6px;
  border: none;
  outline: none;
  border-radius: 6px;
  background: var(--Primary-blue, #618aff);
  box-shadow: 0px 2px 10px 0px rgba(97, 138, 255, 0.5);
  color: white;
`;

export const Flex = styled.div`
  display: flex;
`;
