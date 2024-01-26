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

export interface Offer {
  title: string;
  price: number;
  description: string;
  gallery?: [string];
  url?: string;
  created: number;
  show?: boolean;
}

export interface BlogPost {
  title: string;
  markdown: string;
  gallery?: [string];
  created: number;
  show?: boolean;
}

export interface SingleValueExtra {
  value: string;
}

export interface Twitter {
  handle: string;
}

export interface Extras {
  twitter?: Twitter;
  email?: SingleValueExtra[];
  blog?: BlogPost[];
  offers?: Offer[];
  wanted?: Wanted[];
  supportme?: SupportMe;
  liquid?: SingleValueExtra[];
  github?: [{ [key: string]: string }];
  coding_languages?: [{ [key: string]: string }];
  tribes?: [{ [key: string]: string }];
  lightning?: [{ [key: string]: string }];
  amboss?: [{ [key: string]: string }];
}

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

export interface Post {
  title: string;
  content: string;
  created: number;
  gallery?: [string];
  show?: boolean;
}

export interface FocusedWidgetProps {
  name: string;
  values: any[];
  errors: any[];
  initialValues: any;
  setFieldTouched: (string, boolean) => void;
  setFieldValue: (string, any) => void;
  item: any;
  setShowFocused: (boolean) => void;
  setDisableFormButtons: (boolean) => void;
}

export interface InvitePeopleSearchProps {
  peopleList: any;
  handleChange?: (any) => void;
  setAssigneefunction?: (string) => void;
  newDesign?: boolean;
  isProvidingHandler: boolean;
  handleAssigneeDetails: (any) => void;
}

export interface WidgetProps {
  values: any;
  name: string;
  parentName: string;
  setFieldValue: (string, any) => void;
  setSelected: (any) => void;
  label: string;
  single: string;
  icon?: string;
}

export interface WidgetListProps {
  setSelected: (any, i: number) => void;
  deleteItem: (any, i: number) => void;
  schema: any;
  values: { [key: string]: any };
}
