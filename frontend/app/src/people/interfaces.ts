import React from 'react';
import { Extras } from '../components/form/inputs/widgets/interfaces';
import { Person } from '../store/main';
import { MeData } from '../store/ui';
import { Widget } from './main/types';

export interface AuthProps {
  style?: React.CSSProperties;
  onSuccess?: () => void;
}

export interface BountyModalProps {
  basePath: string;
}

export interface FocusViewProps {
  goBack?: () => void;
  config: { [key: string]: any };
  selectedIndex: number;
  canEdit?: boolean;
  person: any;
  personBody?: any;
  buttonsOnBottom?: boolean;
  formHeader?: JSX.Element;
  manualGoBackOnly?: boolean;
  isFirstTimeScreen?: boolean;
  fromBountyPage?: boolean;
  newDesign?: boolean;
  setIsModalSideButton?: boolean;
  ReCallBounties?: () => Promise<void>;
  onSuccess?: () => void;
  extraModalFunction?: () => void;
  deleteExtraFunction?: () => void;
  style?: React.CSSProperties;
  setIsExtraStyle?: any;
}

export interface PeopleMobileeHeaderProps {
  goBack: () => void;
  canEdit: boolean;
  logout: () => void;
  onEdit: () => void;
}

export interface UserInfoProps {
  setShowSupport: (boolean) => void;
}

export interface BountiesProps {
  price: number;
  sessionLength: string;
  priceMin: number;
  priceMax: number;
  codingLanguage: [{ [key: string]: string }];
  title: string;
  person: Person;
  onPanelClick: () => void;
  created?: number;
  ticketUrl?: string;
  loomEmbedUrl?: string;
  description?: any;
  isPaid: boolean;
  widget?: any;
  assignee?: Person;
}

export interface BadgesProps {
  person?: Person;
  txid?: string;
  color?: string;
}

export interface ConnectCardProps {
  person: Person | MeData | undefined;
  dismiss: () => void;
  modalStyle?: React.CSSProperties;
  visible: boolean;
  created?: number;
  dismissConnectModal?: () => void;
}

export interface LoomViewProps {
  loomEmbedUrl?: string;
  onChange?: (string) => void;
  readOnly?: boolean;
  style: React.CSSProperties;
  setIsVideo?: (any) => void;
  name?: string;
  onBlur?: () => void;
  onFocus?: () => void;
}

export interface NameTagProps {
  owner_alias: string;
  owner_pubkey: string;
  img: string;
  created?: number;
  id: number;
  style?: React.CSSProperties;
  widget: any;
  iconSize?: number;
  textSize?: number;
  isPaid?: boolean;
  ticketUrl?: string;
  loomEmbedUrl?: string;
}

export interface NoneSpaceProps {
  banner?: boolean;
  style: React.CSSProperties;
  img: string;
  text: string;
  sub?: string;
  buttonText1?: string;
  buttonText?: string;
  buttonText2?: string;
  Button?: JSX.Element | boolean;
  buttonIcon?: string;
  small?: boolean;
  action?: () => void;
  action1?: () => void;
  action2?: () => void;
}

export interface PageLoadProps {
  show: boolean;
  style?: React.CSSProperties;
  noAnimate?: boolean;
}

export interface NoResultProps {
  loading: boolean;
}

export interface PaidBountiesProps {
  onPanelClick: () => void;
  title: string;
  codingLanguage: [{ [key: string]: string }];
  priceMax: number;
  priceMin: number;
  price: number;
  sessionLength: string;
  assignee: Person;
  description: string;
  owner_alias: string;
  owner_pubkey: string;
  img: string;
  id: number;
  widget: any;
  created: number;
}

export interface QRProps {
  type?: string;
  size: number;
  value: string;
  style?: React.CSSProperties;
}

export interface QRBarProps {
  simple?: boolean;
  value: string | undefined;
  style?: React.CSSProperties;
}

export interface StartUpModalProps {
  closeModal: () => void;
  dataObject: string;
  buttonColor: string;
}

export interface SvgMaskProps {
  svgStyle: React.CSSProperties;
  width: string;
  height: string;
  src: string;
  size: string;
  bgcolor: string;
}

export interface PersonProps extends Person {
  hideActions: boolean;
  small: boolean;
  id: number;
  img: string;
  selected: boolean;
  select: (id: number, unique_name: string, owner_pubkey: string) => void;
  owner_alias: string;
  owner_pubkey: string;
  unique_name: string;
  squeeze: boolean;
  description: string;
}

export interface GithubStatusPillProps {
  status?: string;
  assignee?: Person;
  style?: React.CSSProperties;
}

export interface WantedSummaryProps {
  description: any;
  priceMin: number;
  ticketUrl: string;
  person: any;
  created?: number | undefined;
  repo: string;
  issue: string;
  price: number;
  type: string;
  tribe: string;
  paid: boolean;
  badgeRecipient: string;
  loomEmbedUrl: string;
  codingLanguage: { [key: string]: any };
  estimate_session_length: string;
  assignee: Person;
  fromBountyPage: string;
  wanted_type: string;
  one_sentence_summary: string;
  github_description: string;
  show: boolean;
  setIsModalSideButton: (any) => void;
  setIsExtraStyle: (any) => void;
  formSubmit: (any) => void;
  title: string;
}

export interface CodingBountiesProps {
  deliverables?: string;
  description: any;
  ticketUrl: string;
  assignee: Person;
  created?: number;
  titleString: string;
  nametag: JSX.Element;
  labels?: { [key: string]: any };
  person: Person;
  setIsPaidStatusPopOver?: (boolean) => void;
  creatorStep: number;
  paid: boolean;
  tribe: string;
  saving?: string;
  isPaidStatusPopOver: boolean;
  isPaidStatusBadgeInfo: boolean;
  awardDetails: any;
  isAssigned: boolean;
  dataValue: { [key: string]: any };
  assigneeValue: boolean;
  assignedPerson: Person;
  changeAssignedPerson: () => void;
  sendToRedirect: (string) => void;
  handleCopyUrl: () => void;
  isCopied: boolean;
  setExtrasPropertyAndSave: (string, boolean) => void;
  setIsModalSideButton: (boolean) => void;
  replitLink: string;
  assigneeHandlerOpen: () => void;
  setCreatorStep: (number) => void;
  setIsExtraStyle: (any) => void;
  awards: { [key: string]: any };
  setExtrasPropertyAndSaveMultiple: (string, any) => void;
  handleAssigneeDetails: (any) => void;
  peopleList: Person[];
  setIsPaidStatusBadgeInfo: (any) => void;
  bountyPrice: number;
  selectedAward: string;
  handleAwards: (any) => void;
  repo: string;
  issue: string;
  isMarkPaidSaved: boolean;
  setAwardDetails: (any) => void;
  setBountyPrice: (any) => void;
  owner_idURL: string;
  createdURL: string;
  editAction?: (any) => void;
  deletingState?: boolean;
  deleteAction?: (any) => void;
  priceMin?: number;
  priceMax?: number;
  price?: number;
  estimate_session_length?: string;
  loomEmbedUrl?: string;
  extraModalFunction?: () => void;
}

export interface CodingViewProps {
  paid?: boolean;
  titleString: string;
  labels?: { [key: string]: any };
  price?: number;
  description?: string;
  envHeight?: string;
  estimate_session_length?: string;
  loomEmbedUrl?: string;
  ticketUrl?: string;
  assignee: Person;
  assigneeLabel?: { [key: string]: any };
  nametag?: JSX.Element;
  actionButtons?: boolean | JSX.Element;
  status?: string;
  handleCopyUrl?: () => void;
  isCopied?: boolean;
  tribe?: string;
}

export interface AddToFavoritesProps {
  tribe: string | undefined;
}

export interface WantedViewsProps {
  description: string;
  priceMin: number;
  priceMax: number;
  price?: number;
  person: any;
  created?: number;
  ticketUrl?: string;
  gallery?: any;
  assignee?: Person;
  estimate_session_length?: string;
  loomEmbedUrl?: string;
  showModal?: () => void;
  setDeletePayload?: (boolean) => void;
  key?: string;
  setExtrasPropertyAndSave?: (any) => void;
  saving?: boolean;
  labels?: [{ [key: string]: string }] | never[];
  isClosed?: boolean;
  onPanelClick: () => void;
  status?: string;
  isCodingTask?: boolean;
  show?: string | boolean;
  paid?: boolean;
  isMine?: boolean;
  titleString: string | JSX.Element | JSX.Element[];
}

export interface WantedViews2Props extends WantedViewsProps {
  one_sentence_summary?: string;
  title?: string;
  issue?: string;
  repo?: string;
  type?: string;
  codingLanguage?: any;
  fromBountyPage?: boolean;
}

export interface AboutViewProps {
  price_to_meet?: number;
  extras?: Extras;
  twitter_confirmed?: boolean;
  owner_pubkey?: string;
  description?: string;
  canEdit?: boolean;
}

export interface BlogViewProps {
  title: string;
  markdown: string;
  gallery: string;
  created: number;
}

export interface BountyHeaderProps {
  selectedWidget: Widget;
  scrollValue: boolean;
  onChangeStatus: (number) => void;
  onChangeLanguage: (number) => void;
  checkboxIdToSelectedMap: any;
  checkboxIdToSelectedMapLanguage: any;
}

export interface DeleteTicketModalProps {
  closeModal: () => void;
  confirmDelete: () => void;
}

export interface OfferViewProps {
  gallery: [{ [key: string]: string }];
  title: string;
  description: any;
  price: number;
  person: Person;
  created: number;
  type: string;
  content: string;
}

export interface RenderWidgetsProps {
  widget: any;
}
