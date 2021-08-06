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
        type: "textarea",
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
                handle: strValidator,
            }).default(undefined),
            supportme: Yup.object({
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
                    },
                    {
                        name: 'show',
                        label: "Show In Link",
                        type: "switch",
                    },
                ]
            },
            {
                name: "supportme",
                label: "Support Me",
                type: "widget",
                class: "supportme",
                single: true,
                fields: [
                    {
                        name: 'gallery',
                        label: "Gallery",
                        type: "gallery",
                    },
                    {
                        name: 'description',
                        label: "Description",
                        type: "textarea",
                    },
                    {
                        name: 'url',
                        label: "URL",
                        type: "text",
                    },
                    {
                        name: 'show',
                        label: "Show In Link",
                        type: "switch",
                    },
                ]
            },
            {
                name: "offer",
                label: "Offer",
                type: "widget",
                class: "offer",
                fields: [
                    {
                        name: 'title',
                        label: "Title",
                        type: "text",
                    },
                    {
                        name: 'price',
                        label: "Price",
                        type: "number",
                    },
                    {
                        name: 'gallery',
                        label: "Gallery",
                        type: "gallery",
                    },
                    {
                        name: 'show',
                        label: "Show In Link",
                        type: "switch",
                    },
                ]
            },
            {
                name: "wanted",
                label: "Wanted",
                type: "widget",
                class: "wanted",
                fields: [
                    {
                        name: 'title',
                        label: "Title",
                        type: "text",
                    },
                    {
                        name: 'priceMin',
                        label: "Price Min",
                        type: "number",
                    },
                    {
                        name: 'priceMax',
                        label: "Price Max",
                        type: "number",
                    },
                    {
                        name: 'description',
                        label: "Description",
                        type: "textarea",
                    },
                    {
                        name: 'show',
                        label: "Show In Link",
                        type: "switch",
                    },
                ]
            },
            {
                name: "blog",
                label: "Blog",
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
                        type: "textarea",
                    },
                    {
                        name: 'gallery',
                        label: "Gallery",
                        type: "gallery",
                    },
                    {
                        name: 'show',
                        label: "Show In Link",
                        type: "switch",
                    },
                ],
            },
        ],
        page: 2,
    }
];


// extras.blog.existing