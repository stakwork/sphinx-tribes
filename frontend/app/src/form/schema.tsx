import * as Yup from 'yup'
import { FormField } from "../form";

const strValidator = Yup.string().required('Required')
const nomValidator = Yup.number().required('Required')

// this is source of truth for widget items!
export const meSchema: FormField[] = [
    {
        name: "img",
        label: "Image",
        type: "img",
        page: 1
    },
    {
        name: "pubkey",
        label: "Pubkey*",
        type: "text",
        readOnly: true,
        page: 1
    },
    {
        name: "owner_alias",
        label: "Name*",
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
                    priceMin: Yup.number().when('priceMax', (pricemax) => Yup.number().max(pricemax, `Must be less than max`)
                    )
                }).nullable()
            ),
            offer: Yup.array().of(
                Yup.object().shape({
                    title: strValidator,
                })
            ),
            tribes: Yup.array().of(
                Yup.object().shape({
                    address: strValidator,
                })
            ),
            blog: Yup.array().of(
                Yup.object().shape({
                    title: strValidator,
                    markdown: strValidator,
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
                        label: "Twitter*",
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
                        name: 'url',
                        label: "URL*",
                        type: "text",
                    },
                    {
                        name: 'description',
                        label: "Description",
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
                        label: "Title*",
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
                        label: "Title*",
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
                name: "post",
                label: "Post",
                type: "widget",
                fields: [
                    {
                        name: "title",
                        label: "Title",
                        type: "text"
                    },
                    {
                        name: "content",
                        label: "Content",
                        type: "textarea",
                    },
                    {
                        name: 'gallery',
                        label: "Gallery",
                        type: "gallery",
                    },
                ]
            },
            {
                name: "tribes",
                label: "Tribes",
                type: "widget",
                fields: [
                    {
                        name: 'address',
                        label: "Tribe address*",
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
                name: "blog",
                label: "Blog",
                type: "widget",
                class: "blog",
                fields: [
                    {
                        name: 'title',
                        label: "Title*",
                        type: "text",
                    },
                    {
                        name: 'markdown',
                        label: "Markdown*",
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

export const aboutSchema: FormField[] = [
    {
        name: "img",
        label: "Image",
        type: "img",
        page: 1
    },
    {
        name: "pubkey",
        label: "Pubkey*",
        type: "text",
        readOnly: true,
        page: 1
    },
    {
        name: "owner_alias",
        label: "Name*",
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
];

export const postSchema: FormField[] = [
    {
        name: "title",
        label: "Title",
        type: "text",
        validator: strValidator,
    },
    {
        name: "content",
        label: "Content",
        type: "textarea",
        validator: strValidator,
    },
    {
        name: 'gallery',
        label: "Gallery",
        type: "gallery",
    },
];

export const offerSchema: FormField[] = [
    {
        name: "title",
        label: "Title",
        validator: strValidator,
        type: "text"
    },
    {
        name: "description",
        label: "Description",
        validator: strValidator,
        type: "textarea",
    },
    {
        name: "price",
        label: "Price",
        validator: nomValidator,
        type: "number",
    },
    {
        name: 'gallery',
        label: "Gallery",
        type: "gallery",
    },
];

export const wantedSchema: FormField[] = [
    {
        name: 'title',
        label: "Title*",
        type: "text",
        validator: strValidator,
    },
    {
        name: 'description',
        label: "Description",
        type: "textarea",
        validator: strValidator,
    },
    {
        name: 'priceMin',
        label: "Price Min",
        validator: Yup.number().when('priceMax', (pricemax) => Yup.number().max(pricemax, `Must be less than max`)),
        type: "number",
    },
    {
        name: 'priceMax',
        label: "Price Max",
        validator: nomValidator,
        type: "number",
    },
    {
        name: 'gallery',
        label: "Gallery",
        type: "gallery",
    },

    // {
    //     name: 'show',
    //     label: "Show In Link",
    //     type: "switch",
    // },
];



// extras.blog.existing