import { MeInfo } from '../../../store/ui'

export interface FormState {
    img?: string,
    pubkey: string,
    owner_alias: string,
    description?: string,
    price_to_meet: number,
    id?: number,
    extras?: Extras
}

export interface BlogPost {
    title: string,
    markdown: string,
    createdAt: number
}

export interface Offer {
    title: string,
    price: number,
    description: string,
    img?: string,
    url?: string,
    createdAt: number
}

export interface Wanted {
    title: string,
    priceMin: number,
    priceMax: number,
    description: string,
    url?: string,
    createdAt: number
}

export interface Wanted {
    title: string,
    priceMin: number,
    priceMax: number,
    description: string,
    url?: string,
    createdAt: number
}

export interface Donation {
    title: string,
    description: string,
    createdAt: number,
    url?: string,
    img?: string
}

export interface Extras {
    twitter?: string,
    blog?: BlogPost[],
    offers?: Offer[],
    wanted?: Wanted[],
    donation?: Donation,
}

function doJSONToFormState(json: MeInfo): FormState {
    let formState: FormState = {
        id: json.id || 0,
        pubkey: json.pubkey,
        owner_alias: json.alias || "",
        img: json.photo_url || "",
        price_to_meet: json.price_to_meet || 0,
        description: json.description || "",
        extras: json.extras || {}
    };

    // extras.blog

    {
        name:
        posts: [

        ]
    }

    return formState
}

function doFormStateToJSON() {

}