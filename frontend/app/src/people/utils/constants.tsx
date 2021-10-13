import { aboutSchema, postSchema, wantedSchema, meSchema, offerSchema, offerSkillSchema, offerOtherSchema, wantedCodingTaskSchema, wantedOtherSchema } from "../../form/schema";

const MAX_UPLOAD_SIZE = 10194304 //10MB

const widgetConfigs = {
    about: {
        label: 'About',
        name: 'about',
        single: true,
        skipEditLayer: true,
        submitText: 'Save',
        schema: aboutSchema,
        action: {
            text: 'Edit Profile',
            icon: 'edit'
        },
    },
    post: {
        label: 'Blog',
        name: 'post',
        submitText: 'Post',
        schema: postSchema,
        action: {
            text: 'Create a Post',
            icon: 'add',
            info: "What's on your mind?",
            infoIcon: 'chat_bubble_outline'
        },
        noneSpace: {
            me: {
                img: 'no_posts.png',
                text: 'What’s on your mind?',
                buttonText: 'Create a post',
                buttonIcon: 'add'
            },
            otherUser: {
                img: 'no_posts2.png',
                text: 'No Posts Yet',
                sub: 'Looks like this person hasn’t posted anything yet.'
            }
        }
    },
    offer: {
        label: 'Portfolio',
        name: 'offer',
        submitText: 'Save',
        modalStyle: {
            width: 'auto',
            maxWidth: 'auto',
            minWidth: '400px',
            height: 'auto'
        },
        schema: offerSchema,
        action: {
            text: 'Add to Portfolio',
            // text: 'Sell Something',
            icon: 'local_offer'
        },
        noneSpace: {
            me: {
                img: 'no_offers.png',
                text: 'Use lightning network to market your skills!',
                // text: 'Use lightning network to sell your digital goods!',
                buttonText: 'Add to Portfolio',
                // buttonText: 'Sell something',
                buttonIcon: 'local_offer'
            },
            otherUser: {
                img: 'no_offers2.png',
                text: 'Portfolio is Empty',
                // text: 'No Offers Yet',
                sub: "Looks like this person hasn't shared their porfolio."
                // sub: 'Looks like this person is not selling anything yet.'
            }
        }
    },
    wanted: {
        label: 'Wanted',
        name: 'wanted',
        submitText: 'Save',
        modalStyle: {
            width: 'auto',
            maxWidth: 'auto',
            minWidth: '400px',
            height: 'auto'
        },
        schema: wantedSchema,
        action: {
            text: 'Add to Wanted',
            icon: 'favorite_outline'
        },
        noneSpace: {
            me: {
                img: 'no_wanted.png',
                text: 'Make a list of items and services you need.',
                buttonText: 'Add to wanted',
                buttonIcon: 'favorite_outline'
            },
            otherUser: {
                img: 'no_wanted2.png',
                text: 'No Wanteds Yet',
                sub: 'Looks like this person doesn’t need anything yet.'
            }
        }
    },
}

const formDropdownOptions = {
    'wanted': [
        {
            value: 'wanted_coding_task',
            label: 'Coding Task',
            schema: wantedCodingTaskSchema,
            description: 'Post a coding task referencing your github repo.'
        },
        {
            value: 'wanted_other',
            label: 'Other',
            schema: wantedOtherSchema,
            description: 'Could be anything.'
        },
    ],
    'offer': [
        {
            value: 'offer_skill',
            label: 'Skill',
            schema: offerSkillSchema,
            description: 'Build your portfolio.'
        },
        {
            value: 'offer_other',
            label: 'Other',
            schema: offerOtherSchema,
            description: 'Could be anything.'
        },
    ],
}



export { MAX_UPLOAD_SIZE, widgetConfigs, formDropdownOptions }