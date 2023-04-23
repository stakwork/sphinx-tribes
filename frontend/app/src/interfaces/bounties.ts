import { Person } from "store/main";
import { CSSProperties } from "styled-components";

export interface BountiesDescriptionProps {
    description: any;
    codingLanguage: [{[key: string]: string}];
    isPaid: boolean;
    title: string;
    style?: React.CSSProperties
} 

export interface BountiesPriceProps {
    sessionLength: boolean;
    priceMin: number;
    price: number;
    style?: React.CSSProperties
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
        background: string
    }
}