import { MeInfo } from '../../../store/ui'

export interface FormState {
    img?: string,
    pubkey: string,
    owner_alias?: string,
    alias?: string,
    description?: string,
    price_to_meet: number,
    id?: number,
    extras?: Extras
}

export interface BlogPost {
    title: string,
    markdown: string,
    gallery?: [string],
    createdAt: number,
    show?: boolean
}

export interface Offer {
    title: string,
    price: number,
    description: string,
    gallery?: [string],
    url?: string,
    createdAt: number,
    show?: boolean
}

export interface Wanted {
    title: string,
    priceMin: number,
    priceMax: number,
    description: string,
    url?: string,
    createdAt: number,
    show?: boolean
}

export interface SupportMe {
    title: string,
    description: string,
    createdAt: number,
    url?: string,
    gallery?: [string],
    show?: boolean
}

export interface Twitter {
    handle: string
}

export interface Extras {
    twitter?: Twitter,
    blog?: BlogPost[],
    offers?: Offer[],
    wanted?: Wanted[],
    supportme?: SupportMe
}