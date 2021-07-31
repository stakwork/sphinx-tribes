import * as Yup from 'yup'
import { FormField } from "../form";

const strValidator = Yup.string().required('Required')

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
        validator: strValidator,
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
        validator: Yup.object().shape({
            twitter: Yup.object({
                handle: strValidator
            }).default(undefined),
            donations: Yup.object({
                url: strValidator
            }).default(undefined),
            wanted: Yup.array().of(
                Yup.object().shape({
                    title: strValidator,
                }).nullable()
            ),
            offer: Yup.array().of(
                Yup.object().shape({
                    title: strValidator,
                })
            ),
            blog: Yup.array().of(
                Yup.object().shape({
                    title: strValidator,
                })
            ),
        }),
        extras: [
            {
                name: "twitter",
                label: "Twitter",
                type: "widget",
                class: "twitter",
                single: true,
                icon: 'twitter',
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
                fields: [
                    {
                        name: 'img',
                        label: "image",
                        type: "img",
                    },
                    {
                        name: 'title',
                        label: "Title",
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
                fields: [
                    {
                        name: 'img',
                        label: "image",
                        type: "img",
                    },
                    {
                        name: 'title',
                        label: "Title",
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
                fields: [
                    {
                        name: 'title',
                        label: "Title",
                        type: "text",
                    },
                    {
                        name: 'markdown',
                        label: "Markdown",
                        type: "text",
                    }
                ],
            },
        ],
        page: 2,
    }
];


// extras.blog.existing