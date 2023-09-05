import styled from 'styled-components';

interface styledProps {
  color?: any;
}

export const B = styled.small`
  font-weight: bold;
  display: block;
  margin-bottom: 10px;
`;
export const N = styled.div<styledProps>`
  font-family: Barlow;
  font-style: normal;
  font-weight: 500;
  font-size: 17px;
  line-height: 26px;
  text-align: center;
  margin-bottom: 10px;
  color: ${(p: any) => p?.color && p?.color.grayish.G100};
`;
export const ModalBottomText = styled.div<styledProps>`
  position: absolute;
  bottom: -36px;
  width: 310;
  background-color: transparent;
  display: flex;
  justify-content: center;
  .bottomText {
    margin-left: 12px;
    color: ${(p: any) => p?.color && p?.color.pureWhite};
  }
`;
export const InvoiceForm = styled.div`
  margin: 10px 0px;
  text-align: left;
`;
export const InvoiceLabel = styled.label`
  font-size: 0.9rem;
  font-weight: bold;
`;
export const InvoiceInput = styled.input`
  padding: 10px 20px;
  border-radius: 10px;
  border: 0.5px solid black;
`;