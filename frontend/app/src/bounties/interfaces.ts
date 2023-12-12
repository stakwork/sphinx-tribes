import { CodingLanguageLabel } from 'people/interfaces';
import { Person } from 'store/main';

export interface BountiesDescriptionProps {
  description?: any;
  codingLanguage: Array<CodingLanguageLabel>;
  isPaid?: boolean;
  title: string;
  style?: React.CSSProperties;
  owner_alias: string;
  widget?: any;
  id: number;
  owner_pubkey: string;
  img: string;
  created?: number;
  name?: string;
  uuid?: string;
  org_uuid?: string;
  org_img?: string;
}

export interface BountiesPriceProps {
  sessionLength?: boolean | string;
  priceMin?: number;
  price: number;
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
