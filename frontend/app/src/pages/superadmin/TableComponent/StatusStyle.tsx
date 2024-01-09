import styled from 'styled-components';

interface styledProps { color?: any }

export const StatusPopOverCheckbox = styled.div<styledProps>`
  padding: 15px 18px;
  border-right: 1px solid ${(p: any) => p.color && p.color.grayish.G700};
  user-select: none;
  .leftBoxHeading {
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 700;
    font-size: 12px;
    line-height: 32px;
    text-transform: uppercase;
    color: ${(p: any) => p.color && p.color.grayish.G100};
    margin-bottom: 10px;
  }

  &.CheckboxOuter > div {
    display: flex;
    flex-direction: column;
    flex-wrap: wrap;

    .euiCheckboxGroup__item {
      .euiCheckbox__square {
        top: 5px;
        border: 1px solid ${(p: any) => p?.color && p?.color?.grayish.G500};
        border-radius: 2px;
      }
      .euiCheckbox__input + .euiCheckbox__square {
        background: ${(p: any) => p?.color && p?.color?.pureWhite} no-repeat center;
      }
      .euiCheckbox__input:checked + .euiCheckbox__square {
        border: 1px solid ${(p: any) => p?.color && p?.color?.blue1};
        background: ${(p: any) => p?.color && p?.color?.blue1} no-repeat center;
        background-image: url('static/checkboxImage.svg');
      }
      .euiCheckbox__label {
        font-family: 'Barlow';
        font-style: normal;
        font-weight: 500;
        font-size: 13px;
        line-height: 16px;
        color: ${(p: any) => p?.color && p?.color?.grayish.G50};
        &:hover {
          color: ${(p: any) => p?.color && p?.color?.grayish.G05};
        }
      }
      input.euiCheckbox__input:checked ~ label {
        color: ${(p: any) => p?.color && p?.color?.blue1};
      }
    }
  }
`;

export const FlexDivStatus = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 4px;
  position: relative;
  cursor: pointer;
  user-select: none;
  top: 2.5px;
`;

export const FilterContainer = styled.div<styledProps>`
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  user-select: none;
    .statusText {
      color: #5f6368;
      font-family: Barlow;
      font-size: 15px;
      font-style: normal;
      font-weight: 500;
      line-height: 18px;
  }
  .subStatusText {
    color: #3c3f41;
    font-family: Barlow;
    font-size: 15px;
    font-style: normal;
    font-weight: 500;
    line-height: 18px;
    cursor: pointer;
  }

`;

export const StatusCheckboxItem = styled.div<styledProps>`
padding: 20px 28px;
.euiCheckboxGroup__item {
  .euiCheckbox__square {
    border: 1px solid ${(p: any) => p?.color && p?.color?.grayish.G500};
    border-radius: 2px;
  }
  .euiCheckbox__input + .euiCheckbox__square {
    background: ${(p: any) => p?.color && p?.color?.pureWhite} no-repeat center;
  }
  .euiCheckbox__input:checked + .euiCheckbox__square {
    border: 1px solid ${(p: any) => p?.color && p?.color?.blue1};
    background: ${(p: any) => p?.color && p?.color?.blue1} no-repeat center;
    background-image: url('static/checkboxImage.svg');
  }
  .euiCheckbox__label {
    top: 1px;
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 500;
    font-size: 15px;
    line-height: 18px;
    color: ${(p: any) => p?.color && p?.color?.grayish.G50};
    &:hover {
      color: ${(p: any) => p?.color && p?.color?.grayish.G05};
    }
  }
  input.euiCheckbox__input:checked ~ label {
    color: ${(p: any) => p?.color && p?.color?.blue1};
  }
}
`;