// This file holds definitions for the style properties of
// bounty constants. This includes status, award, and
// language options.

import { languages } from 'config';

export const status = ['Open', 'Assigned', 'Paid'];

export const GetValue = (arr: string[]) =>
  arr.map((val: any) => ({
    id: val,
    label: val,
    value: val
  }));

export const awards = [
  {
    id: 'Admin',
    label: 'Admin',
    label_icon: '/static/awards/Admin_award.svg'
  },
  {
    id: 'Moderator',
    label: 'Moderator',
    label_icon: '/static/awards/Moderator_award.svg'
  },
  {
    id: 'Developer',
    label: 'Developer',
    label_icon: '/static/awards/Developer_award.svg'
  },
  {
    id: 'First 1000 members',
    label: 'First 1000 members',
    label_icon: '/static/awards/1st_1000_member_award.svg'
  },
  {
    id: 'Contributing 1M sats ',
    label: 'Contributing 1M sats ',
    label_icon: '/static/awards/1M_sat_award.svg'
  },
  {
    id: 'New Member',
    label: 'New Member',
    label_icon: '/static/awards/new_member_award.svg'
  },
  {
    id: 'Early/Founding Member',
    label: 'Early/Founding Member',
    label_icon: '/static/awards/early_founding_member_award.svg'
  },
  {
    id: 'Conversation Starter',
    label: 'Conversation Starter',
    label_icon: '/static/awards/conversation_award.svg'
  },
  {
    id: 'VIP member',
    label: 'VIP member',
    label_icon: '/static/awards/vip_member_award.svg'
  },
  {
    id: 'Group Expert',
    label: 'Group Expert',
    label_icon: '/static/awards/group_expert_award.svg'
  }
];

export const coding_languages = languages;

export interface ColorOption {
  readonly value: string;
  readonly label: string;
  readonly color: string;
  readonly border: string;
  readonly background: string;
  readonly isFixed?: boolean;
  readonly isDisabled?: boolean;
}

// define the colors for all language options
export const colorLanguageOptions = {
  Lightning: {
    border: '1px solid rgba(184, 37, 95, 0.1)',
    background: 'rgba(184, 37, 95, 0.1)',
    color: '#B8255F'
  },
  Javascript: {
    border: '1px solid rgba(219, 64, 53, 0.1)',
    background: 'rgba(219, 64, 53, 0.1)',
    color: '#DB4035'
  },
  Typescript: {
    border: '1px solid rgba(255, 153, 51, 0.1)',
    background: ' rgba(255, 153, 51, 0.1)',
    color: '#FF9933'
  },
  Node: {
    border: '1px solid rgba(255, 191, 59, 0.1)',
    background: 'rgba(255, 191, 59, 0.1)',
    color: '#FFBF3B'
  },
  Golang: {
    border: '1px solid rgba(175, 184, 59, 0.1)',
    background: 'rgba(175, 184, 59, 0.1)',
    color: '#AFB83B'
  },
  Swift: {
    border: '1px solid rgba(126, 204, 73, 0.1)',
    background: 'rgba(126, 204, 73, 0.1)',
    color: '#7ECC49'
  },
  Kotlin: {
    border: '1px solid rgba(41, 148, 56, 0.1)',
    background: 'rgba(41, 148, 56, 0.1)',
    color: '#299438'
  },
  MySQL: {
    border: '1px solid rgba(106, 204, 188, 0.1)',
    background: 'rgba(106, 204, 188, 0.1)',
    color: '#6ACCBC'
  },
  PHP: {
    border: '1px solid rgba(21, 143, 173, 0.1)',
    background: 'rgba(21, 143, 173, 0.1)',
    color: '#158FAD'
  },
  R: {
    border: '1px solid rgba(64, 115, 255, 0.1)',
    background: 'rgba(64, 115, 255, 0.1)',
    color: '#4073FF'
  },
  'C#': {
    border: '1px solid rgba(136, 77, 255, 0.1)',
    background: 'rgba(136, 77, 255, 0.1)',
    color: '#884DFF'
  },
  'C++': {
    border: '1px solid rgba(175, 56, 235, 0.1)',
    background: 'rgba(175, 56, 235, 0.1)',
    color: '#AF38EB'
  },
  Java: {
    border: '1px solid rgba(235, 150, 235, 0.1)',
    background: 'rgba(235, 150, 235, 0.1)',
    color: '#EB96EB'
  },
  Rust: {
    border: '1px solid rgba(224, 81, 148, 0.1)',
    background: 'rgba(224, 81, 148, 0.1)',
    color: '#E05194'
  },
  Ruby: {
    border: '1px solid rgba(255, 32, 110, 0.1)',
    background: 'rgba(255, 32, 110, 0.1)',
    color: '#FF206E'
  },
  Postgres: {
    border: '1px solid rgba(65, 234, 212, 0.1)',
    background: 'rgba(65, 234, 212, 0.1)',
    color: '#41EAD4'
  },
  'Elastic search': {
    border: '1px solid rgba(251, 255, 18, 0.1)',
    background: 'rgba(251, 255, 18, 0.1)',
    color: '#FBFF12'
  },
  Python: {
    border: '1px solid rgba(75, 100, 74, 0.1)',
    background: 'rgba(75, 100, 74, 0.1)',
    color: '#4B644A'
  },
  Other: {
    border: '1px solid rgba(21, 143, 173, 1)',
    background: 'rgba(21, 143, 173, 0.1)',
    color: '#158FAD'
  }
};

export function GetColorLanguageOptions(): readonly ColorOption[] {
  return languages.map((lan: string) => {
    if (colorLanguageOptions[lan]) {
      return {
        value: lan,
        label: lan,
        color: colorLanguageOptions[lan].color,
        border: colorLanguageOptions[lan].border,
        background: colorLanguageOptions[lan].background,
        isFixed: true,
        isDisabled: false
      };
    }

    return {
      label: 'no-lan-option',
      value: 'no-lan-option',
      border: '1px solid rgba(255, 141, 133, 0.1)',
      background: 'rgba(255, 141, 133, 0.1)',
      color: '#FF8D85',
      isFixed: true,
      isDisabled: false
    };
  });
}
