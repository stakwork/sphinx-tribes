import * as Yup from 'yup'
import { FormField } from "../form";

export const meSchema: FormField[] = [
    {
        name: "img",
        label: "Image",
        type: "img",
        page: 1
    },
    {
        name: "pubkey",
        label: "Pubkey",
        type: "text",
        readOnly: true,
        page: 1
    },
    {
        name: "owner_alias",
        label: "Name",
        type: "text",
        required: true,
        validator: Yup.string().required('Required'),
        page: 1,
    },
    {
        name: "description",
        label: "Description",
        type: "text",
        page: 1,
    },
    {
        name: "price_to_meet",
        label: "Price to Meet",
        type: "number",
        page: 1,
    },
    {
        name: "id",
        label: "ID",
        type: "hidden",
        page: 1,
    },
    {
        name: 'extras',
        label: 'Widgets',
        type: 'widgets',
        extras: [
            {
                name: "twitter",
                label: "Twitter",
                type: "widget",
                class: "twitter",
                single: true,
                icon: '',
                fields: [
                    {
                        name: 'handle',
                        label: "Twitter Handle",
                        type: "text",
                        prepend: '@',
                    }
                ]
            },
            {
                name: "donations",
                label: "Donations",
                type: "widget",
                class: "donations",
                single: true,
                icon: '',
                fields: [
                    {
                        name: 'img',
                        label: "Image",
                        type: "img",
                    },
                    {
                        name: 'bio',
                        label: "Bio",
                        type: "text",
                    },
                    {
                        name: 'url',
                        label: "URL",
                        type: "text",
                    }
                ]
            },
            {
                name: "offer",
                label: "Offer",
                itemLabel: "Offer",
                type: "widget",
                class: "offer",
                icon: '',
                fields: [
                    {
                        name: 'img',
                        label: "image",
                        type: "img",
                    },
                    {
                        name: 'header',
                        label: "Header",
                        type: "text",
                    },
                    {
                        name: 'price',
                        label: "Price",
                        type: "number",
                    }
                ]
            },
            {
                name: "wanted",
                label: "Wanted",
                itemLabel: "Listing",
                type: "widget",
                class: "wanted",
                icon: '',
                fields: [
                    {
                        name: 'img',
                        label: "image",
                        type: "img",
                    },
                    {
                        name: 'header',
                        label: "Header",
                        type: "text",
                    },
                    {
                        name: 'price',
                        label: "Price",
                        type: "number",
                    }
                ]
            },
            {
                name: "blog",
                label: "Blog",
                itemLabel: "Post",
                type: "widget",
                class: "blog",
                icon: '',
                fields: [
                    {
                        name: 'title',
                        label: "URL",
                        type: "text",
                    }
                ],
            },
        ],
        page: 2,
    }
];


// extras.blog.existing