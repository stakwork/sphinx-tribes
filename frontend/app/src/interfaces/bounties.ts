import { Person } from "store/main";
import { CSSProperties } from "styled-components";

export interface BountiesDescriptionProps {
    description: any;
    codingLanguage: [{[key: string]: string}];
    isPaid: boolean;
    title: string;
    style?: CSSProperties
} 

export interface BountiesPriceProps {
    sessionLength: boolean;
    priceMin: number;
    price: number;
    style?: CSSProperties
    priceMax?: number;
}

export interface BountiesProfileProps {
    UserProfileContainerStyle?: CSSProperties;
    UserImageStyle?: CSSProperties;
    userInfoStyle?: CSSProperties;
    statusStyle?: CSSProperties;
    assignee: Person;
    isNameClickable?: boolean;
    status: string;
    NameContainerStyle?: CSSProperties;
    canViewProfile?: boolean;
    statusStyles?: {
        width: string; 
        height: string; 
        background: string
    }
}