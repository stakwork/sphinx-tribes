import React, { ReactNode } from 'react';
import { Extras } from '../components/form/inputs/widgets/interfaces';
import { Person, PersonBounty } from '../store/main';
import { MeData } from '../store/ui';
import { Widget } from './main/types';

export interface AuthProps {
  style?: React.CSSProperties;
  onSuccess?: () => void;
}

export interface BountyModalProps {
  basePath: string;
  bountyOwner?: Person;
  fromPage?: 'usertickets' | 'bounties' | 'organization';
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
  bounty?: PersonBounty[];
  setRemoveNextAndPrev?: (boolean) => void;
  setAfterEdit?: (boolean) => void;
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

export interface CodingLanguageLabel {
  label: string;
  value: string;
}
export interface BountiesProps {
  price: number;
  sessionLength: string;
  priceMin: number;
  priceMax: number;
  codingLanguage: Array<CodingLanguageLabel>;
  title: string;
  person: Person;
  onPanelClick: () => void;
  created?: number;
  ticketUrl?: string;
  loomEmbedUrl?: string;
  org_uuid?: string;
  description?: any;
  isPaid: boolean;
  widget?: any;
  assignee?: Person;
  name?: string;
  uuid?: string;
  img?: string;
  id?: number;
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
  ticket_url?: string;
  loomEmbedUrl?: string;
  org_img?: string;
  org_name?: string;
  org_uuid?: string;
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
  codingLanguage: Array<CodingLanguageLabel>;
  priceMax: number;
  priceMin: number;
  price: number;
  sessionLength: string;
  assignee: Person;
  description: string;
  owner_alias: string;
  owner_pubkey: string;
  org_uuid?: string;
  img: string;
  id: number;
  widget: any;
  created: number;
  name?: string;
  org_img?: string;
}

export interface QRProps {
  type?: string;
  size: number;
  value: string;
  style?: React.CSSProperties;
}

export interface QRBarProps {
  simple?: boolean;
  value: string;
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

export interface StatusPillProps {
  assignee?: Person;
  style?: React.CSSProperties;
  paid?: boolean;
}

export interface WantedSummaryProps {
  description: any;
  priceMin?: number;
  ticket_url: string;
  person: any;
  created?: number | undefined;
  repo: string;
  issue: string;
  price?: number;
  type: string;
  tribe: string;
  paid: boolean;
  badgeRecipient: string;
  loomEmbedUrl?: string;
  coding_languages: string[];
  estimated_session_length?: string;
  assignee: Person;
  fromBountyPage: string;
  wanted_type: string;
  one_sentence_summary: string;
  github_description: string;
  show: boolean;
  setIsModalSideButton: (any) => void;
  setIsExtraStyle: (any) => void;
  formSubmit: (any, notEdit?: boolean) => void;
  title: string;
  org_uuid?: string;
  id?: number;
  owner_id?: string;
  markPaidOrUnpaid?: ReactNode;
}

export type LocalPaymeentState = 'UNKNOWN' | 'PAID' | 'UNPAID';

export interface CodingBountiesProps extends WantedSummaryProps {
  deliverables?: string;
  description: any;
  ticket_url: string;
  assignee: Person;
  created?: number;
  titleString: string;
  nametag: JSX.Element;
  labels?: Array<CodingLanguageLabel>;
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
  estimated_session_length?: string;
  loomEmbedUrl?: string;
  extraModalFunction?: () => void;
  commitment_fee?: number;
  bounty_expires?: string;
  org_uuid?: string;
  id?: number;
  localPaid: LocalPaymeentState;
  setLocalPaid: (state: LocalPaymeentState) => void;
  isMobile?: boolean;
  assigneeLabel?: { [key: string]: any };
  actionButtons?: boolean | JSX.Element;
}

export interface CodingViewProps extends WantedSummaryProps {
  // paid?: boolean;
  // price?: number;
  // description?: string;
  // estimated_session_length?: string;
  // loomEmbedUrl?: string;
  // tribe?: string;
  titleString: string;
  labels?: Array<CodingLanguageLabel>;
  envHeight?: string;
  ticketUrl?: string;
  assignee: Person;
  assigneeLabel?: { [key: string]: any };
  nametag?: JSX.Element;
  actionButtons?: boolean | JSX.Element;
  status?: string;
  handleCopyUrl?: () => void;
  isCopied?: boolean;
  payBounty?: ReactNode;
  markUnpaid?: ReactNode;
  showPayBounty?: boolean;
  paid: boolean;
  // owner_id: string;
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
  estimated_session_length?: string;
  loomEmbedUrl?: string;
  showModal?: () => void;
  setDeletePayload?: (boolean) => void;
  key?: string;
  setExtrasPropertyAndSave?: (any) => void;
  saving?: boolean;
  labels?: Array<CodingLanguageLabel>;
  isClosed?: boolean;
  onPanelClick: () => void;
  status?: string;
  isCodingTask?: boolean;
  show?: string | boolean;
  paid?: boolean;
  isMine?: boolean;
  titleString: string | JSX.Element | JSX.Element[];
  bounty_expires?: string;
  commitment_fee?: number;
  name?: string;
  img?: string;
  org_uuid?: string;
  id?: number;
}

export interface WantedViews2Props extends WantedViewsProps {
  one_sentence_summary?: string;
  title?: string;
  issue?: string;
  repo?: string;
  type?: string;
  coding_languages?: any;
  fromBountyPage?: boolean;
}

export interface AboutViewProps {
  price_to_meet?: number;
  extras?: Extras;
  twitter_confirmed?: boolean;
  owner_pubkey?: string;
  description?: string;
  canEdit?: boolean;
  owner_route_hint?: string;
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

export interface OrgBountyHeaderProps {
  onChangeStatus: (number) => void;
  checkboxIdToSelectedMap: any;
  languageString: string;
  org_uuid: string;
}
export interface PeopleHeaderProps {
  onChangeLanguage: (number) => void;
  checkboxIdToSelectedMapLanguage: any;
}

export interface DeleteTicketModalProps {
  closeModal: () => void;
  confirmDelete: () => void;
  text?: string;
  imgUrl?: string;
  userDelete?: boolean;
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
