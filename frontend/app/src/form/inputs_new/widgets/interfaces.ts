import { MeInfo } from '../../../store/ui';

export interface FormState {
  img?: string;
  pubkey: string;
  owner_alias?: string;
  alias?: string;
  description?: string;
  price_to_meet: number;
  id?: number;
  extras?: Extras;
}

export interface BlogPost {
  title: string;
  markdown: string;
  gallery?: [string];
  created: number;
  show?: boolean;
}

export interface Post {
  title: string;
  content: string;
  created: number;
  gallery?: [string];
  show?: boolean;
}

export interface Offer {
  title: string;
  price: number;
  description: string;
  gallery?: [string];
  url?: string;
  created: number;
  show?: boolean;
}

export interface Wanted {
  title: string;
  priceMin: number;
  priceMax: number;
  gallery?: [string];
  description: string;
  url?: string;
  created: number;
  show?: boolean;
}

export interface SupportMe {
  title: string;
  description: string;
  created: number;
  url?: string;
  gallery?: [string];
  show?: boolean;
}

export interface Twitter {
  handle: string;
}

export interface SingleValueExtra {
  value: string;
}

export interface Extras {
  twitter?: Twitter;
  blog?: BlogPost[];
  offers?: Offer[];
  wanted?: Wanted[];
  supportme?: SupportMe;
  liquid?: SingleValueExtra;
}
