import {
  aboutSchema,
  wantedSchema,
  offerSkillSchema,
  offerOtherSchema,
  wantedCodingTaskSchema,
  wantedOtherSchema,
  organizationSchema,
  organizationUserSchema
} from '../../components/form/schema';

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
  organizations: {
    label: 'Organizations',
    name: 'organizations',
    submitText: 'Save',
    modalStyle: {
      width: 'auto',
      maxWidth: 'auto',
      minWidth: '400px',
      minHeight: '40%',
      maxHeight: '70%'
    },
    schema: organizationSchema,
    action: {
      text: 'Add Organization',
      icon: 'add'
    },
    noneSpace: {
      noUserResult: {
        img: 'no_org.png',
        text: 'Manage and organize your tickets',
        sub: 'Fund and pay bounties directly through the website, add members, organize tickets, and more!'
      },
      noResult: {
        img: 'no_org.png',
        text: 'No Organization Yet',
        sub: 'Looks like this person has not created or added to any organizations yet.'
      }
    }
  },
  badges: {
    label: 'Badges',
    name: 'badges',
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
        buttonIcon: 'work'
      },
      otherUser: {
        img: '',
        text: 'No Badges',
        sub: "Looks like this person doesn't have any Badges yet."
      }
    }
  },
  // TODO: REMOVE
  wanted: {
    label: 'Bounties',
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
  },
  usertickets: {
    label: 'Assigned Bounties',
    name: 'userwanted',
    submitText: 'Save',
    modalStyle: {
      width: 'auto',
      maxWidth: 'auto',
      minWidth: '400px',
      minHeight: '40%',
      maxHeight: '70%'
    },
    schema: [],
    action: {
      text: 'Add New Ticket',
      icon: 'local_offer'
    },
    noneSpace: {
      noResult: {
        img: 'no_wanted2.png',
        text: 'No Assigned Tickets Yet',
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
  earlyMember: {
    title: 'Early Adopter',
    src: 'EarlyMember.svg'
  }
};

const nonWidgetConfigs = {
  organizationusers: {
    label: 'Organization Users',
    name: 'organizationusers',
    submitText: 'Save',
    modalStyle: {
      width: 'auto',
      maxWidth: 'auto',
      minWidth: '400px',
      minHeight: '40%',
      maxHeight: '70%'
    },
    schema: organizationUserSchema
  }
};

export { MAX_UPLOAD_SIZE, widgetConfigs, formDropdownOptions, badges, nonWidgetConfigs };
