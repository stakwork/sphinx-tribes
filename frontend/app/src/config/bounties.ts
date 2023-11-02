export interface ColorOption {
  readonly value: string;
  readonly label: string;
  readonly color: string;
  readonly border: string;
  readonly background: string;
  readonly isFixed?: boolean;
  readonly isDisabled?: boolean;
}

// all sphinx language options
export const coding_languages: ColorOption[] = [
  {
    value: 'lightning',
    label: 'Lightnint',
    border: '1px solid rgba(184, 37, 95, 0.1)',
    background: 'rgba(184, 37, 95, 0.1)',
    color: '#B8255F'
  },
  {
    value: 'javascript',
    label: 'Javascript',
    border: '1px solid rgba(219, 64, 53, 0.1)',
    background: 'rgba(219, 64, 53, 0.1)',
    color: '#DB4035'
  },
  {
    value: 'typescript',
    label: 'Typescript',
    border: '1px solid rgba(255, 153, 51, 0.1)',
    background: ' rgba(255, 153, 51, 0.1)',
    color: '#FF9933'
  },
  {
    value: 'node',
    label: 'Node',
    border: '1px solid rgba(255, 191, 59, 0.1)',
    background: 'rgba(255, 191, 59, 0.1)',
    color: '#FFBF3B'
  },
  {
    value: 'golang',
    label: 'Golang',
    border: '1px solid rgba(175, 184, 59, 0.1)',
    background: 'rgba(175, 184, 59, 0.1)',
    color: '#AFB83B'
  },
  {
    value: 'swift',
    label: 'Swift',
    border: '1px solid rgba(126, 204, 73, 0.1)',
    background: 'rgba(126, 204, 73, 0.1)',
    color: '#7ECC49'
  },
  {
    value: 'kotlin',
    label: 'Kotlin',
    border: '1px solid rgba(41, 148, 56, 0.1)',
    background: 'rgba(41, 148, 56, 0.1)',
    color: '#299438'
  },
  {
    value: 'mysql',
    label: 'MySQL',
    border: '1px solid rgba(106, 204, 188, 0.1)',
    background: 'rgba(106, 204, 188, 0.1)',
    color: '#6ACCBC'
  },
  {
    value: 'php',
    label: 'PHP',
    border: '1px solid rgba(21, 143, 173, 0.1)',
    background: 'rgba(21, 143, 173, 0.1)',
    color: '#158FAD'
  },
  {
    value: 'r',
    label: 'R',
    border: '1px solid rgba(64, 115, 255, 0.1)',
    background: 'rgba(64, 115, 255, 0.1)',
    color: '#4073FF'
  },
  {
    value: 'cs',
    label: 'C#',
    border: '1px solid rgba(136, 77, 255, 0.1)',
    background: 'rgba(136, 77, 255, 0.1)',
    color: '#884DFF'
  },
  {
    value: 'cpp',
    label: 'C++',
    border: '1px solid rgba(175, 56, 235, 0.1)',
    background: 'rgba(175, 56, 235, 0.1)',
    color: '#AF38EB'
  },
  {
    value: 'java',
    label: 'Java',
    border: '1px solid rgba(235, 150, 235, 0.1)',
    background: 'rgba(235, 150, 235, 0.1)',
    color: '#EB96EB'
  },
  {
    value: 'rust',
    label: 'Rust',
    border: '1px solid rgba(224, 81, 148, 0.1)',
    background: 'rgba(224, 81, 148, 0.1)',
    color: '#E05194'
  },
  {
    value: 'ruby',
    label: 'ruby',
    border: '1px solid rgba(255, 32, 110, 0.1)',
    background: 'rgba(255, 32, 110, 0.1)',
    color: '#FF206E'
  },
  {
    value: 'postgres',
    label: 'Postgres',
    border: '1px solid rgba(65, 234, 212, 0.1)',
    background: 'rgba(65, 234, 212, 0.1)',
    color: '#41EAD4'
  },
  {
    value: 'elasticSearch',
    label: 'Elastic search',
    border: '1px solid rgba(251, 255, 18, 0.1)',
    background: 'rgba(251, 255, 18, 0.1)',
    color: '#FBFF12'
  },
  {
    value: 'python',
    label: 'Python',
    border: '1px solid rgba(75, 100, 74, 0.1)',
    background: 'rgba(75, 100, 74, 0.1)',
    color: '#4B644A'
  },
  {
    value: 'other',
    label: 'Other',
    border: '1px solid rgba(21, 143, 173, 1)',
    background: 'rgba(21, 143, 173, 0.1)',
    color: '#158FAD'
  }
];

// all sphinx timeframe options
export const time_estimation = [
  'Less than 1 hour',
  'Less than 3 hours',
  'More than 3 hours',
  'Not sure yet'
];

export const estimated_budget_15_min = ['USD $10', 'USD $20', 'USD $30', 'USD $40', 'USD $50'];

// all sphinx task options
export const help_wanted_coding_task_schema = [
  'Web development',
  'Mobile development',
  'Design',
  'Desktop app',
  'Dev ops',
  'Bitcoin / Lightning',
  'Other'
];

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

// all sphinx task type options
export const help_wanted_other_schema = ['Troubleshooting', 'Debugging', 'Tutoring'];

// to whom the bounty is assigned
export const status = ['Open', 'Assigned', 'Paid'];

export const GetValue = (arr: string[]) =>
  arr.map((val: any) => ({
    id: val,
    label: val,
    value: val
  }));
