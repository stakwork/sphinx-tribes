import {
  aboutSchema,
  postSchema,
  wantedSchema,
  meSchema,
  offerSchema,
  offerSkillSchema,
  offerOtherSchema,
  wantedCodingTaskSchema,
  wantedOtherSchema
} from '../../form/schema';

const MAX_UPLOAD_SIZE = 10194304; //10MB

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
    }
  },
  badges: {
    label: 'Badges',
    name: 'badges',
    // hideOnMainPage: true,
    single: true,
    skipEditLayer: true,
    action: {
      text: 'Edit Profile',
      icon: 'edit'
    },
    noneSpace: {
      me: {
        img: '',
        text: 'No Badges',
        sub: 'Click here to learn about badges',
        buttonText: 'Add to Portfolio',
        // buttonText: 'Sell something',
        buttonIcon: 'work'
      },
      otherUser: {
        img: '',
        text: 'No Badges',
        // text: 'No Offers Yet',
        sub: "Looks like this person doesn't have any Badges yet."
        // sub: 'Looks like this person is not selling anything yet.'
      }
    }
  },
  // post: {
  //     label: 'Blog',
  //     name: 'post',
  //     submitText: 'Post',
  //     schema: postSchema,
  //     action: {
  //         text: 'Create a Post',
  //         icon: 'add',
  //         info: "What's on your mind?",
  //         infoIcon: 'chat_bubble_outline'
  //     },
  //     noneSpace: {
  //         me: {
  //             img: 'no_posts.png',
  //             text: 'What’s on your mind?',
  //             buttonText: 'Create a post',
  //             buttonIcon: 'add'
  //         },
  //         otherUser: {
  //             img: 'no_posts2.png',
  //             text: 'No Posts Yet',
  //             sub: 'Looks like this person hasn’t posted anything yet.'
  //         }
  //     }
  // },
  offer: {
    label: 'Portfolio',
    name: 'offer',
    submitText: 'Save',
    modalStyle: {
      width: 'auto',
      maxWidth: 'auto',
      minWidth: '400px',
      minHeight: '40%',
      maxHeight: '70%'
    },
    schema: offerSchema,
    action: {
      text: 'Add to Portfolio',
      // text: 'Sell Something',
      icon: 'work'
    },
    noneSpace: {
      me: {
        img: 'no_offers.png',
        text: 'Use lightning network to market your skills!',
        // text: 'Use lightning network to sell your digital goods!',
        buttonText: 'Add to Portfolio',
        // buttonText: 'Sell something',
        buttonIcon: 'work'
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
    label: 'Tickets',
    name: 'wanted',
    submitText: 'Save',
    modalStyle: {
      width: 'auto',
      maxWidth: 'auto',
      minWidth: '400px',
      minHeight: '40%',
      maxHeight: '70%'
    },
    schema: wantedSchema,
    action: {
      text: 'Add New Ticket',
      icon: 'local_offer'
    },
    noneSpace: {
      me: {
        img: 'no_wanted.png',
        text: 'Make a list of github tickets you want help on.',
        buttonText: 'Add New Ticket',
        buttonIcon: 'local_offer'
      },
      otherUser: {
        img: 'no_wanted2.png',
        text: 'No Tickets Yet',
        sub: 'Looks like this person doesn’t need anything yet.'
      }
    }
  }
};

const formDropdownOptions = {
  wanted: [
    {
      value: 'freelance_job_request',
      label: 'Freelance Job Request',
      schema: wantedCodingTaskSchema
      // description: 'Post a coding task referencing your github repo.',
    },
    {
      value: 'live_help',
      label: 'Live Help',
      schema: wantedOtherSchema
      // description: 'Could be anything.',
    }
  ],
  offer: [
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
    }
  ]
};

const badges = {
  // '1k': {
  //     title: '1k',
  //     src: '1K.svg',
  // },
  // 'admin': {
  //     title: 'Admin',
  //     src: 'Admin.svg',
  // },
  // 'conversationStarter': {
  //     title: 'Conversation Starter',
  //     src: 'ConversationStarter.svg',
  // },
  // 'developer': {
  //     title: 'Developer',
  //     src: 'Developer.svg',
  // },
  earlyMember: {
    title: 'Early Adopter',
    src: 'EarlyMember.svg'
  }
  // 'groupExpert': {
  //     title: 'Group Expert',
  //     src: 'GroupExpert.svg',
  // },
  // 'moderator': {
  //     title: 'Moderator',
  //     src: 'Moderator.svg',
  // },
  // 'newMember': {
  //     title: 'New Member',
  //     src: 'NewMember.svg',
  // },
  // 'vip': {
  //     title: 'VIP',
  //     src: 'VIP.svg',
  // },
  // 'whale': {
  //     title: 'Whale',
  //     src: 'Whale.svg',
  // }
};

export { MAX_UPLOAD_SIZE, widgetConfigs, formDropdownOptions, badges };
