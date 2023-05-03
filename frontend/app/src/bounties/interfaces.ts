import { Person } from 'store/main';
import { CSSProperties } from 'styled-components';

export interface BountiesDescriptionProps {
  description: any;
  codingLanguage: [{ [key: string]: string }];
  isPaid: boolean;
  title: string;
  style?: React.CSSProperties;
  owner_alias: string;
  widget?: any;
  id: number;
  owner_pubkey: string;
  img: string;
  created: number;
}

export interface BountiesPriceProps {
  sessionLength?: string;
  priceMin: number;
  price: Number;
  style?: React.CSSProperties;
  priceMax?: number;
}

export interface BountiesProfileProps {
  UserProfileContainerStyle?: React.CSSProperties;
  UserImageStyle?: React.CSSProperties;
  userInfoStyle?: React.CSSProperties;
  statusStyle?: React.CSSProperties;
  assignee: Person;
  isNameClickable?: boolean;
  status: string;
  NameContainerStyle?: React.CSSProperties;
  canViewProfile?: boolean;
  statusStyles?: {
    width: string;
    height: string;
    background: string;
  };
}
