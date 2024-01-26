import styled from 'styled-components';
import { PeopleHeaderProps } from 'people/interfaces';
import { observer } from 'mobx-react-lite';
import React, { useState, useEffect } from 'react';
import { EuiCheckboxGroup, EuiPopover, EuiText } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';
import { colors } from 'config';
import { filterCount } from '../../helpers';
import { GetValue, coding_languages } from '../utils/languageLabelStyle';

interface styledProps {
  color?: any;
}

const FilterWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const FilterTrigger = styled.div<styledProps>`
  width: 78px;
  height: 48px;
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  margin-left: 19px;
  cursor: pointer;
  user-select: none;
  .filterImageContainer {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 48px;
    width: 36px;
    .materialIconImage {
      color: ${(p: any) => p.color && p.color.grayish.G200};
      cursor: pointer;
      font-size: 18px;
      margin-top: 4px;
    }
  }
  .filterText {
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 500;
    font-size: 16px;
    line-height: 19px;
    display: flex;
    align-items: center;
    color: ${(p: any) => p.color && p.color.grayish.G200};
  }
  &:hover {
    .filterImageContainer {
      .materialIconImage {
        color: ${(p: any) => p.color && p.color.grayish.G50} !important;
        cursor: pointer;
        font-size: 18px;
        margin-top: 4px;
      }
    }
    .filterText {
      color: ${(p: any) => p.color && p.color.grayish.G50};
    }
  }
  &:active {
    .filterImageContainer {
      .materialIconImage {
        color: ${(p: any) => p.color && p.color.grayish.G10} !important;
        cursor: pointer;
        font-size: 18px;
        margin-top: 4px;
      }
    }
    .filterText {
      color: ${(p: any) => p.color && p.color.grayish.G10};
    }
  }
`;

const FilterCount = styled.div<styledProps>`
  height: 20px;
  width: 20px;
  border-radius: 50%;
  margin-left: 4px;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: -5px;
  background: ${(p: any) => p?.color && p.color.blue1};
  .filterCountText {
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 500;
    font-size: 13px;
    display: flex;
    align-items: center;
    text-align: center;
    color: ${(p: any) => p.color && p.color.pureWhite};
  }
`;

const PopOverBox = styled.div<styledProps>`
  display: flex;
  flex-direction: column;
  max-height: 304px;
  padding: 15px 0px 20px 21px;
  .rightBoxHeading {
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 700;
    font-size: 12px;
    line-height: 32px;
    text-transform: uppercase;
    color: ${(p: any) => p.color && p.color.grayish.G100};
  }
`;

const EuiPopOverCheckboxWrapper = styled.div<styledProps>`
  min-width: 285px;
  max-width: 285px;
  height: 240px;
  user-select: none;

  &.CheckboxOuter > div {
    height: 100%;
    display: grid;
    grid-template-columns: 1fr 1fr;
    justify-content: center;
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

const Coding_Languages = GetValue(coding_languages);

const PeopleHeader = ({ onChangeLanguage, checkboxIdToSelectedMapLanguage }: PeopleHeaderProps) => {
  const [isPopoverOpen, setIsPopoverOpen] = useState<boolean>(false);
  const [filterCountNumber, setFilterCountNumber] = useState<number>(0);
  const onToggleButton = () => setIsPopoverOpen((prev: boolean) => !prev);
  const closePopover = () => setIsPopoverOpen(false);
  const color = colors['light'];

  useEffect(() => {
    setFilterCountNumber(filterCount(checkboxIdToSelectedMapLanguage));
  }, [checkboxIdToSelectedMapLanguage]);

  const panelStyles = {
    border: 'none',
    boxShadow: `0px 1px 20px ${color.black90}`,
    background: `${color.pureWhite}`,
    borderRadius: '6px',
    minWidth: '300px',
    minHeight: '304px',
    marginTop: '0px',
    marginLeft: '20px'
  };

  return (
    <FilterWrapper>
      <EuiPopover
        button={
          <FilterTrigger onClick={onToggleButton} color={color}>
            <div className="filterImageContainer">
              <MaterialIcon
                className="materialIconImage"
                icon="tune"
                style={{
                  color: isPopoverOpen ? color.grayish.G10 : ''
                }}
              />
            </div>
            <EuiText
              className="filterText"
              style={{
                color: isPopoverOpen ? color.grayish.G10 : ''
              }}
            >
              Filter
            </EuiText>
          </FilterTrigger>
        }
        panelStyle={panelStyles}
        isOpen={isPopoverOpen}
        closePopover={closePopover}
        panelClassName="yourClassNameHere"
        panelPaddingSize="none"
        anchorPosition="downLeft"
      >
        <div
          style={{
            display: 'flex',
            flexDirection: 'row'
          }}
        >
          <PopOverBox color={color}>
            <EuiText className="rightBoxHeading">Skills</EuiText>
            <EuiPopOverCheckboxWrapper className="CheckboxOuter" color={color}>
              <EuiCheckboxGroup
                options={Coding_Languages}
                idToSelectedMap={checkboxIdToSelectedMapLanguage}
                onChange={(id: any) => {
                  onChangeLanguage(id);
                }}
              />
            </EuiPopOverCheckboxWrapper>
          </PopOverBox>
        </div>
      </EuiPopover>
      {filterCountNumber > 0 && (
        <FilterCount color={color}>
          <EuiText className="filterCountText">{filterCountNumber}</EuiText>
        </FilterCount>
      )}
    </FilterWrapper>
  );
};

export default observer(PeopleHeader);
