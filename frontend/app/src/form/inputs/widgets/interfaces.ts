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

}

export type Blog = BlogPost[]

export interface Extras {
    twitter?: string,

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