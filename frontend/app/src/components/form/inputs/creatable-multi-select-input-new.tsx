import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { EuiCheckboxGroup, EuiPopover, EuiText } from '@elastic/eui';
import type { Props } from './propsType';
import { colors } from '../../../config/colors';
import ImageButton from '../../common/Image_button';
import {
  coding_languages,
  GetValue,
  LanguageObject
} from '../../../people/utils/language_label_style';
import { SvgMask } from '../../../people/utils/svgMask';

const codingLanguages = GetValue(coding_languages);

export default function CreatableMultiSelectInputNew({ error, label, handleChange }: Props) {
  let labeltext = label;
  if (error) labeltext = `${labeltext} (INCORRECT FORMAT)`;
  const color = colors['light'];

  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const onButtonClick = () => setIsPopoverOpen((isPopoverOpen) => !isPopoverOpen);
  const closePopover = () => setIsPopoverOpen(false);
  const [checkboxIdToSelectedMap, setCheckboxIdToSelectedMap] = useState({});
  const [labels, setLabels] = useState<any>([]);
  const [data, setData] = useState<any>([]);

  const onChange = (optionId) => {
    let trueCount = 0;
    for (const [key, value] of Object.entries(checkboxIdToSelectedMap)) {
      if (value) {
        trueCount += 1;
      }
    }
    if (!(!checkboxIdToSelectedMap[optionId] && trueCount >= 6)) {
      const newCheckboxIdToSelectedMap = {
        ...checkboxIdToSelectedMap,
        ...{
          [optionId]: !checkboxIdToSelectedMap[optionId]
        }
      };

      setCheckboxIdToSelectedMap(newCheckboxIdToSelectedMap);
    }
  };

  useEffect(() => {
    setLabels(LanguageObject.filter((x) => checkboxIdToSelectedMap[x.label]));
  }, [checkboxIdToSelectedMap]);

  return (
    <div
      style={{
        paddingTop: '11px'
      }}
    >
      <EuiPopover
        className="EuiPopOver"
        anchorPosition="downCenter"
        panelStyle={{
          height: '180px',
          marginTop: '-8px',
          boxShadow: 'none !important',
          borderRadius: '0px 0px 6px 6px',
          backgroundRepeat: 'no-repeat',
          border: `1px solid ${color.blue2}`,
          borderTop: `1px solid ${color.grayish.G700}`
        }}
        button={
          <ImageButton
            buttonText={labeltext}
            ButtonContainerStyle={{
              width: '292px',
              height: '40px',
              border: !isPopoverOpen ? '' : `1px solid ${color?.blue1}`,
              borderBottom: !isPopoverOpen ? '' : `1px solid ${color?.grayish.G700}`,
              borderRadius: !isPopoverOpen ? '4px' : '4px 4px 0px 0px',
              display: 'flex',
              justifyContent: 'flex-start',
              paddingLeft: '18px',
              marginRight: '2px',
              marginTop: isPopoverOpen ? '0.9px' : '0px',
              boxShadow: `0px 1px 2px ${color.black100} `,
              fontWeight: '500'
            }}
            endImageSrc={'/static/Skill_drop_down.svg'}
            endingImageContainerStyle={{
              right: 0,
              top: -2
            }}
            buttonTextStyle={{
              color: `${color.grayish.G300}`,
              textAlign: 'center',
              fontSize: '14px',
              fontWeight: '500',
              fontFamily: 'Barlow'
            }}
            buttonAction={() => {
              onButtonClick();
            }}
          />
        }
        isOpen={isPopoverOpen}
        closePopover={closePopover}
      >
        <EuiPopOverCheckbox className="CheckboxOuter" color={color}>
          <EuiCheckboxGroup
            options={codingLanguages}
            idToSelectedMap={checkboxIdToSelectedMap}
            onChange={(id) => {
              onChange(id);
              setData([...data, { value: id, label: id }]);
              handleChange(data);
            }}
          />
        </EuiPopOverCheckbox>
      </EuiPopover>
      <LabelsContainer
        style={{
          padding: !isPopoverOpen && labels.length > 0 ? '16px 0px 24px 0px' : ''
        }}
      >
        {!isPopoverOpen &&
          labels &&
          labels?.map((x, index) => (
            <Label
              key={x.label}
              value={x}
              onClick={() => {
                onChange(x.label);
              }}
              style={{
                margin: 4
              }}
            >
              <EuiText className="labelText">{x.label}</EuiText>
              <SvgMask
                src={'/static/label_cross.svg'}
                bgcolor={x.color}
                height={'23px'}
                width={'16px'}
                size={'8px'}
                svgStyle={{
                  marginLeft: '2px',
                  marginTop: '1px'
                }}
              />
            </Label>
          ))}
      </LabelsContainer>
    </div>
  );
}

interface styledProps {
  color?: any;
}
interface labelProps {
  value?: any;
}

const EuiPopOverCheckbox = styled.div<styledProps>`
  height: 180px;
  padding: 10px 0px 0px 20px;
  margin-right: 3px;
  overflow-y: scroll;
  &.CheckboxOuter > div {
    height: 100%;
    display: grid;
    grid-template-columns: 1fr 1fr;
    .euiCheckboxGroup__item {
      .euiCheckbox__square {
        top: 5px;
        border: 1px solid ${(p) => p?.color && p?.color?.grayish.G500};
        border-radius: 2px;
      }
      .euiCheckbox__input + .euiCheckbox__square {
        background: ${(p) => p?.color && p?.color?.pureWhite} no-repeat center;
      }
      .euiCheckbox__input:checked + .euiCheckbox__square {
        border: 1px solid ${(p) => p?.color && p?.color?.blue1};
        background: ${(p) => p?.color && p?.color?.blue1} no-repeat center;
        background-image: url('static/checkboxImage.svg');
      }
      .euiCheckbox__label {
        font-family: 'Barlow';
        font-style: normal;
        font-weight: 500;
        font-size: 13px;
        line-height: 16px;
        color: ${(p) => p?.color && p?.color?.grayish.G50};
      }
      input.euiCheckbox__input:checked ~ label {
        color: ${(p) => p?.color && p?.color?.blue1};
      }
    }
  }
`;

const LabelsContainer = styled.div<labelProps>`
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  flex-wrap: wrap;
  min-height: 24px;
  width: 100%;
`;

const Label = styled.div<labelProps>`
  height: 23px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  text-align: center;
  border: ${(p) => p?.value && p?.value.border};
  background: ${(p) => p?.value && p?.value.background};
  margin-right: 4px;
  border-radius: 4px;
  padding: 2px 6px;
  cursor: pointer;
  .labelText {
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 500;
    font-size: 13px;
    line-height: 16px;
    color: ${(p) => p?.value && p?.value.color};
  }
`;
