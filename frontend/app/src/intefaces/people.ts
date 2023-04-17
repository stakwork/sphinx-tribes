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
}

export interface BadgesProps {
    person?: Person;
    txid?: string;
    color?: string;
}

export interface ConnectCardProps {
    person: Person | MeData | undefined;
    dismiss: () => void;
    modalStyle: React.CSSProperties;
    visible: boolean;
}

export interface LoomViewProps {
    loomEmbedUrl: string;
    onChange: (string) => void;
    readOnly: boolean;
    style: React.CSSProperties;
    setIsVideo: (any) => void;
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
    sub: string;
    buttonText1?: string;
    buttonText?: string;
    buttonText2?: string;
    Button?: JSX.Element;
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