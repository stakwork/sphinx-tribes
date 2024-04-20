declare namespace Cypress {
    interface Chainable {
        upsertlogin(person: Person): Promise<string>;
    }

    type RandomObject = { [key: string]: string };

    type LoginResponse = {
        jwt: string;
        //user: Person;
    }

    type Person = {
        id?: number;
        uuid?: string;
        owner_pubkey: string;
        owner_alias: string;
        unique_name: string;
        description: string;
        tags: String[]
        img: string;
        created?: number;
        updated?: string;
        unlisted: boolean;
        deleted: boolean;
        last_login?: number;
        owner_route_hint: string;
        owner_contact_key: string;
        price_to_meet: number;
        new_ticket_time?: number;
        twitter_confirmed: boolean;
        referred_by: number;
        extras: RandomObject;
        github_issues: RandomObject;
    }
}
