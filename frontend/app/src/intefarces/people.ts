import { Person } from "store/main";
import { MeData } from "store/ui";

export interface AuthProps  {
    style?: React.CSSProperties;
    onSuccess?: () => void;
}

export interface BountyModalProps {
    basePath: string;
}

export interface FocusViewProps {
    goBack?: () => void;
    config: {[key: string]: any}
    selectedIndex: number;
    canEdit?: boolean;
    person: any,
    personBody?: any;
    buttonsOnBottom?: boolean,
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
    setShowSupport: (boolean) => void
}

export interface BountiesProps {
    assignee: Person;
    price: number;
    sessionLength: string;
    priceMin: number;
    priceMax: number;
    codingLanguage: string | string[];
    title: string;
    person: Person;
    onPanelClick: () => void;
    created?: number;
    ticketUrl?: string;
    loomEmbedUrl?: string;
    description?: any;
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
}

export interface LoomViewProps {
    loomEmbedUrl: string;
    onChange: (string) => void;
    readOnly?: boolean;
    style: React.CSSProperties;
    setIsVideo: (any) => void;
    name?: string;
    onBlur?: () => void;
    onFocus?: () => void;
}

export interface NameTagProps {
    owner_alias: string;
    owner_pubkey: string;
    img: string;
    created: number;
    id: number;
    style: React.CSSProperties;
    widget: any;
    iconSize: number;
    textSize: number;
    isPaid: boolean;
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
    codingLanguage: string | string[];
    priceMax: number;
    priceMin: number;
    price: number;
    sessionLength: number;
    assignee: Person;
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
    svgStyle:  React.CSSProperties;
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
    assignee: Person; 
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
    price: Number;
    type: string;
    tribe: string;
    paid: string;
    badgeRecipient: string;
    loomEmbedUrl: string;
    codingLanguage: {[key: string]: any};
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
    titleString: string;
    nametag: JSX.Element;
    labels?: {[key: string]: any};
    person: Person;
    setIsPaidStatusPopOver?:  (boolean) => void;
    creatorStep: number;
    paid: string;
    tribe: string;
    saving?: string;
    isPaidStatusPopOver: boolean;
    isPaidStatusBadgeInfo: boolean;
    awardDetails: any;
    isAssigned: boolean;
    dataValue: {[key: string]: any};
    assigneeValue: boolean;
    assignedPerson: Person
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
    awards: {[key: string]: any};
    setExtrasPropertyAndSaveMultiple: (string, any) => void;
    handleAssigneeDetails: (any) =>  void;
    peopleList: Person[];
    setIsPaidStatusBadgeInfo: (any) => void
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
    editAction?: boolean;
    deletingState?: boolean;
    deleteAction?: boolean;
    priceMin?: number;
    priceMax?: number;
    price?: Number;
    estimate_session_length?: string;
    extraModalFunction?: () => void;
}