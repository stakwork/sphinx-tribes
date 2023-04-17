import { Person } from "store/main";

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