/* eslint-disable jsx-a11y/img-redundant-alt */
import { EuiCheckboxGroup, EuiLoadingSpinner, EuiPopover, EuiText } from '@elastic/eui';
import React, { useCallback, useEffect, useState } from 'react';
import { useInView } from 'react-intersection-observer';
import styled from 'styled-components';
import { useStores } from '../../../../store';
import { colors } from '../../../../config/colors';
import {
  coding_languages,
  GetValue,
  LanguageObject
} from '../../../../people/utils/languageLabelStyle';
import { SvgMask } from '../../../../people/utils/SvgMask';
import ImageButton from '../../../common/ImageButton';
import { InvitePeopleSearchProps } from './interfaces';

interface styledProps {
  color?: any;
}

interface labelProps {
  value?: any;
}

const SearchOuterContainer = styled.div<styledProps>`
  min-height: 256x;
  max-height: 256x;
  max-width: 292px;
  background: ${(p: any) => p?.color && p?.color?.pureWhite};
  display: flex;
  flex-direction: column;
  align-items: center;

  .SearchSkillContainer {
    display: flex;
    flex-direction: row;
    justify-content: center;
    margin-bottom: 8px;
    height: fit-content;
    .SearchContainer {
      position: relative;
      .SearchInput {
        background: ${(p: any) => p?.color && p?.color?.pureWhite};
        border: 1px solid ${(p: any) => p?.color && p?.color?.grayish.G600};
        border-radius: 4px;
        width: 177px;
        height: 40px;
        outline: none;
        overflow: hidden;
        caret-color: ${(p: any) => p?.color && p?.color?.textBlue1};
        padding: 0px 32px 0px 18px;
        margin-right: 11px;
        font-family: Roboto !important;
        font-weight: 400;
        font-size: 13px;
        line-height: 35px;

        :focus-visible {
          background: ${(p: any) => p?.color && p?.color?.pureWhite};
          border: 1px solid ${(p: any) => p?.color && p?.color?.blue2};
          outline: none;
          .SearchText {
            outline: none;
            background: ${(p: any) => p?.color && p?.color?.pureWhite};
            border: 1px solid ${(p: any) => p?.color && p?.color?.grayish.G600};
            outline: none;
          }
        }
        ::placeholder {
          color: ${(p: any) => p?.color && p?.color?.grayish.G300};
          font-family: 'Roboto';
          font-style: normal;
          font-weight: 400;
          font-size: 13px;
          line-height: 35px;
          display: flex;
          align-items: center;
        }
      }
      .ImageContainer {
        height: 40px;
        width: 43px;
        position: absolute;
        top: 0px;
        right: 10px;
        display: flex;
        justify-content: center;
        align-items: center;
        :active {
          .crossImage {
            filter: brightness(0) saturate(100%) invert(22%) sepia(5%) saturate(563%)
              hue-rotate(161deg) brightness(91%) contrast(86%) !important;
          }
        }
      }
    }

    .EuiPopOver {
      margin-top: 0px;
      .SkillSetContainer {
        height: 40px;
        width: 103px;
        border: 1px solid ${(p: any) => p?.color && p?.color?.grayish.G600};
        border-radius: 4px;
        display: flex;
        justify-content: center;
        align-items: center;
        cursor: pointer;
        user-select: none;
      }
    }
  }

  .OuterContainer {
    width: auto;
    background: ${(p: any) => p?.color && p.color.grayish.G950};
    box-shadow: inset 0px 2px 8px ${(p: any) => p?.color && p.color.black100};
    .PeopleList {
      background: ${(p: any) => p?.color && p?.color?.grayish.G950};
      box-shadow: inset 0px 2px 8px ${(p: any) => p?.color && p.color.black100};
      width: 400px;
      padding: 0 49px 16px;
      min-height: 256px;
      max-height: 256px;
      overflow-y: scroll;
      .People {
        height: 32px;
        min-width: 291.5813903808594px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-top: 16px;
        padding: 0px 0px 0px 6px;
        .PeopleDetailContainer {
          display: flex;
          justify-content: center;
          align-items: center;
          .ImageContainer {
            height: 32px;
            width: 32px;
            border-radius: 50%;
            overflow: hidden;
            display: flex;
            justify-content: center;
            align-items: center;
            object-fit: cover;
          }
          .PeopleName {
            font-family: 'Barlow';
            font-style: normal;
            font-weight: 500;
            font-size: 13px;
            line-height: 16px;
            color: ${(p: any) => p?.color && p?.color?.grayish.G10};
            margin-left: 10px;
          }
        }
      }
      .no_result_container {
        display: flex;
        height: 210px;
        justify-content: center;
        align-items: center;
        .no_result_text {
          font-family: 'Barlow';
          font-size: 16px;
          font-weight: 600;
          color: ${(p: any) => p?.color && p?.color?.grayish.G50};
          word-spacing: 0.08em;
        }
      }
    }
  }
`;

const EuiPopOverCheckbox = styled.div<styledProps>`
  width: 292px;
  height: 293px;
  padding: 20px 10px 10px 18px;

  &.CheckboxOuter > div {
    height: 100%;
    display: grid;
    grid-template-columns: 1fr 1fr;
    overflow-y: scroll;
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
      }
      input.euiCheckbox__input:checked ~ label {
        color: ${(p: any) => p?.color && p?.color?.blue1};
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
  border: ${(p: any) => p?.value && p?.value.border};
  background: ${(p: any) => p?.value && p?.value.background};
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
    color: ${(p: any) => p?.value && p?.value.color};
  }
`;

const InvitedButton = styled.div<styledProps>`
  width: 86px;
  height: 30px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  background: ${(p: any) => p.color && p.color.button_secondary.main};
  box-shadow: 0px 2px 10px ${(p: any) => p.color && p.color.button_secondary.shadow};
  border-radius: 32px;
  color: ${(p: any) => p.color && p.color.pureWhite};
  :hover {
    background: ${(p: any) => p.color && p.color.button_secondary.hover};
  }
  :active {
    background: ${(p: any) => p.color && p.color.button_secondary.active};
  }
  .nextText {
    font-family: 'Barlow';
    font-size: 13px;
    font-weight: 500;
    line-height: 15px;
    user-select: none;
    text-align: center;
    letter-spacing: 0.01em;
  }
`;

const LoaderContainer = styled.div`
  height: 100%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
`;
const codingLanguages = GetValue(coding_languages);

const InvitePeopleSearch = (props: InvitePeopleSearchProps) => {
  const color = colors['light'];
  const [searchValue, setSearchValue] = useState<string>('');
  const [peopleData, setPeopleData] = useState<any>(props?.peopleList);
  const [inviteNameId, setInviteNameId] = useState<number>(0);
  const [checkboxIdToSelectedMap, setCheckboxIdToSelectedMap] = useState({});
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const [labels, setLabels] = useState<any>([]);
  const [initialPeopleCount, setInitialPeopleCount] = useState<number>(20);
  const onButtonClick = () => setIsPopoverOpen((isPopoverOpen: boolean) => !isPopoverOpen);
  const closePopover = () => setIsPopoverOpen(false);
  const { main } = useStores();

  const { ref, inView } = useInView({
    triggerOnce: false,
    threshold: 0
  });

  useEffect(() => {
    if (inView) {
      setTimeout(() => {
        setInitialPeopleCount(initialPeopleCount + 10);
      }, 2000);
    }
  }, [inView, initialPeopleCount]);

  useEffect(() => {
    async function updatePeopleData() {
      setLabels(LanguageObject.filter((x: any) => checkboxIdToSelectedMap[x.label]));
      let peopleList = props?.peopleList;
      if (searchValue) {
        peopleList = await main.getPeopleByNameAliasPubkey(searchValue);
      }
      setPeopleData(
        Object.keys(checkboxIdToSelectedMap).every((key: any) => !checkboxIdToSelectedMap[key])
          ? peopleList
          : peopleList?.filter(
              ({ extras }: any) =>
                extras?.coding_languages?.some(
                  ({ value }: any) => checkboxIdToSelectedMap[value] ?? false
                )
            )
      );
    }
    updatePeopleData();
  }, [checkboxIdToSelectedMap, searchValue]);

  useEffect(() => {
    if (
      searchValue === '' &&
      Object.keys(checkboxIdToSelectedMap).every((key: any) => !checkboxIdToSelectedMap[key])
    ) {
      setPeopleData(props?.peopleList);
    }
  }, [searchValue, props, checkboxIdToSelectedMap]);

  const handler = useCallback((e: any, value: any) => {
    if (value === '') {
      setSearchValue(e.target.value);
    } else {
      setSearchValue(value);
    }
  }, []);

  const onChange = (optionId: any) => {
    let trueCount = 0;
    for (const [, value] of Object.entries(checkboxIdToSelectedMap)) {
      if (value) {
        trueCount += 1;
      }
    }
    if (!(!checkboxIdToSelectedMap[optionId] && trueCount >= 4)) {
      const newCheckboxIdToSelectedMap = {
        ...checkboxIdToSelectedMap,
        ...{
          [optionId]: !checkboxIdToSelectedMap[optionId]
        }
      };

      setCheckboxIdToSelectedMap(newCheckboxIdToSelectedMap);
    }
  };

  return (
    <SearchOuterContainer color={color}>
      <div className="SearchSkillContainer">
        <div className="SearchContainer">
          <input
            value={searchValue}
            className="SearchInput"
            onChange={(e: any) => {
              handler(e, '');
            }}
            placeholder={'Type to search ...'}
            style={{
              background: color.pureWhite,
              color: color.text1,
              fontFamily: 'Barlow'
            }}
          />
          {searchValue !== '' && (
            <div
              className="ImageContainer"
              onClick={() => {
                setSearchValue('');
              }}
            >
              <img
                className="crossImage"
                src="/static/search_cross.svg"
                alt="cross_icon"
                height={'12px'}
                width={'12px'}
              />
            </div>
          )}
        </div>

        <EuiPopover
          className="EuiPopOver"
          anchorPosition="downRight"
          panelStyle={{
            marginTop: '-9px',
            boxShadow: 'none !important',
            borderRadius: '6px 0px 6px 6px',
            backgroundImage: "url('/static/panel_bg.svg')",
            backgroundRepeat: 'no-repeat',
            outline: 'none',
            border: 'none'
          }}
          button={
            <ImageButton
              buttonText={'Skills'}
              ButtonContainerStyle={{
                width: '102px',
                height: '40px',
                border: !isPopoverOpen ? '' : `1px solid ${color?.blue1}`,
                borderBottom: !isPopoverOpen ? '' : `1px solid ${color?.grayish.G700}`,
                borderRadius: !isPopoverOpen ? '4px' : '4px 4px 0px 0px',
                display: 'flex',
                justifyContent: 'flex-start',
                paddingLeft: '18px',
                marginRight: '1px',
                marginTop: isPopoverOpen ? '0.9px' : '0px'
              }}
              endImageSrc={'/static/Skill_drop_down.svg'}
              endingImageContainerStyle={{
                left: 60,
                top: -2
              }}
              buttonTextStyle={{
                color: !isPopoverOpen ? `${color.grayish.G300}` : `${color.black500}`,
                textAlign: 'center',
                fontSize: '13px',
                fontWeight: '400',
                fontFamily: 'Roboto'
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
              onChange={(id: any) => {
                onChange(id);
              }}
            />
          </EuiPopOverCheckbox>
        </EuiPopover>
      </div>
      <LabelsContainer
        style={{
          padding: !isPopoverOpen && labels.length > 0 ? '16px 0px 24px 0px' : ''
        }}
      >
        {!isPopoverOpen &&
          labels.length > 0 &&
          labels?.map((x: any) => (
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

      <div className="OuterContainer">
        <div className="PeopleList">
          {peopleData?.slice(0, initialPeopleCount)?.map((value: any) => (
            <div className="People" key={value.id}>
              <div className="PeopleDetailContainer">
                <div className="ImageContainer">
                  <img
                    src={value.img || '/static/person_placeholder.png'}
                    alt={'user-image'}
                    height={'100%'}
                    width={'100%'}
                    style={{
                      opacity: inviteNameId && inviteNameId !== value?.id ? '0.5' : ''
                    }}
                  />
                </div>
                <EuiText
                  className="PeopleName"
                  style={{
                    opacity: inviteNameId && inviteNameId !== value?.id ? '0.5' : ''
                  }}
                >
                  {value.owner_alias}
                </EuiText>
              </div>
              {inviteNameId === value?.id ? (
                <InvitedButton
                  color={color}
                  onClick={() => {
                    handler('', value.owner_alias);
                    setInviteNameId(0);
                    if (props?.handleChange)
                      props?.handleChange({
                        owner_alias: '',
                        owner_pubkey: '',
                        img: '',
                        value: '',
                        label: ''
                      });
                    if (searchValue === '') {
                      setSearchValue('');
                    }
                    if (props.setAssigneefunction) props.setAssigneefunction('');
                  }}
                >
                  <EuiText className="nextText">{props.newDesign ? 'Unassign' : 'Invited'}</EuiText>
                </InvitedButton>
              ) : (
                <ImageButton
                  buttonText={props.newDesign ? 'Assign' : 'Assign'}
                  ButtonContainerStyle={{
                    width: '86px',
                    height: '30px',
                    background: `${color.grayish.G600}`
                  }}
                  buttonAction={() => {
                    if (props.isProvidingHandler) {
                      props.handleAssigneeDetails(value);
                    } else {
                      handler('', value.owner_alias);
                      setInviteNameId(value.id);
                      if (props.handleChange)
                        props?.handleChange({
                          owner_alias: value.owner_alias,
                          owner_pubkey: value.owner_pubkey,
                          img: value.img,
                          value: value.owner_pubkey,
                          label: `${value.owner_alias} (${value.owner_alias
                            .toLowerCase()
                            .replace(' ', '')})`
                        });
                      if (searchValue === '') {
                        setSearchValue('');
                      }
                      if (props.setAssigneefunction) props.setAssigneefunction(value.owner_alias);
                    }
                  }}
                />
              )}
            </div>
          ))}
          {peopleData && peopleData.length > initialPeopleCount && (
            <LoaderContainer ref={ref}>
              <EuiLoadingSpinner size="l" />
            </LoaderContainer>
          )}

          {peopleData?.length === 0 && (
            <div className="no_result_container">
              <EuiText className="no_result_text">No Result Found</EuiText>
            </div>
          )}
        </div>
      </div>
    </SearchOuterContainer>
  );
};

export default InvitePeopleSearch;
