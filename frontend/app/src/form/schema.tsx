import * as Yup from 'yup';
import { FormField } from '../form';
// import { uiStore } from '../store/ui';

const strValidator = Yup.string().trim().required('Required');
const strValidatorNotRequired = Yup.string().trim();
const repoStrValidator = Yup.string()
  .trim()
  .matches(/^[^\/]+\/[^\/]+$/, 'Incorrect format')
  .required('Required');
const repoArrayStrValidator = Yup.array().of(
  Yup.object().shape({
    value: repoStrValidator
  })
);
const badgeObjectStrValidator = Yup.object().shape({
  value: strValidator
});

const nomValidator = Yup.number().required('Required');

const GetValue = (arr: any) => {
  return arr.map((val) => {
    return {
      label: val,
      value: val
    };
  });
};

const languages = [
  'Lightning',
  'Javascript',
  'Typescript',
  'Node',
  'Golang',
  'Swift',
  'Kotlin',
  'MySQL',
  'PHP',
  'R',
  'C#',
  'C++',
  'Java',
  'Rust'
];

const estimation = ['Less than 1 hour', 'Less than 3 hours', 'More than 3 hours', 'Not sure yet'];

const estimated_budget_15_min = ['USD $10', 'USD $20', 'USD $30', 'USD $40', 'USD $50'];

const help_wanted_coding_task_schema = [
  'Web development',
  'Mobile development',
  'Desktop app',
  'Dev ops',
  'Bitcoin / Lightning',
  'other'
];
const help_wanted_other_schema = ['Troubleshooting', 'Debugging', 'Tutoring'];

const codingLanguages = GetValue(languages);
const estimated_time = GetValue(estimation);
const helpType_wanted_coding_task_schema = GetValue(help_wanted_coding_task_schema);
const helpType_wanted_other_schema = GetValue(help_wanted_other_schema);
const estimated_budget_15_min_options = GetValue(estimated_budget_15_min);

// this is source of truth for widget items!
export const meSchema: FormField[] = [
  {
    name: 'img',
    label: 'Image',
    type: 'img',
    page: 1
  },
  {
    name: 'pubkey',
    label: 'Pubkey*',
    type: 'text',
    readOnly: true,
    page: 1
  },
  {
    name: 'owner_alias',
    label: 'Name*',
    type: 'text',
    required: true,
    validator: strValidator,
    page: 1
  },
  {
    name: 'description',
    label: 'Description',
    type: 'textarea',
    page: 1
  },
  {
    name: 'price_to_meet',
    label: 'Price to Meet',
    type: 'number',
    page: 1
  },
  {
    name: 'id',
    label: 'ID',
    type: 'hidden',
    page: 1
  },
  {
    name: 'extras',
    label: 'Widgets',
    type: 'widgets',
    validator: Yup.object().shape({
      alert: Yup.boolean(),
      twitter: Yup.object({
        handle: strValidator
      }).default(undefined),
      supportme: Yup.object({
        url: strValidator
      }).default(undefined),
      wanted: Yup.array().of(
        Yup.object()
          .shape({
            title: strValidator,
            priceMin: Yup.number().when('priceMax', (pricemax) =>
              Yup.number().max(pricemax, `Must be less than max`)
            )
          })
          .nullable()
      ),
      offer: Yup.array().of(
        Yup.object().shape({
          title: strValidator
        })
      ),
      tribes: Yup.array().of(
        Yup.object().shape({
          address: strValidator
        })
      ),
      blog: Yup.array().of(
        Yup.object().shape({
          title: strValidator,
          markdown: strValidator
        })
      )
    }),
    extras: [
      {
        name: 'alert',
        label: 'Alerts',
        type: 'switch',
        single: true
      },
      {
        name: 'twitter',
        label: 'Twitter',
        type: 'widget',
        class: 'twitter',
        single: true,
        icon: 'twitter',
        fields: [
          {
            name: 'handle',
            label: 'Twitter*',
            type: 'text',
            prepend: '@'
          },
          {
            name: 'show',
            label: 'Show In Link',
            type: 'switch'
          }
        ]
      },
      {
        name: 'supportme',
        label: 'Support Me',
        type: 'widget',
        class: 'supportme',
        single: true,
        fields: [
          {
            name: 'url',
            label: 'URL*',
            type: 'text'
          },
          {
            name: 'description',
            label: 'Description',
            type: 'textarea'
          },
          {
            name: 'gallery',
            label: 'Gallery',
            type: 'gallery'
          },
          {
            name: 'show',
            label: 'Show In Link',
            type: 'switch'
          }
        ]
      },
      {
        name: 'offer',
        label: 'Offer',
        type: 'widget',
        class: 'offer',
        fields: [
          {
            name: 'title',
            label: 'Title*',
            type: 'text'
          },
          {
            name: 'price',
            label: 'Price',
            type: 'number'
          },
          {
            name: 'gallery',
            label: 'Gallery',
            type: 'gallery'
          },
          {
            name: 'show',
            label: 'Show In Link',
            type: 'switch'
          }
        ]
      },
      {
        name: 'wanted',
        label: 'Wanted',
        type: 'widget',
        class: 'wanted',
        fields: [
          {
            name: 'title',
            label: 'Title*',
            type: 'text'
          },
          {
            name: 'priceMin',
            label: 'Price Min',
            type: 'number'
          },
          {
            name: 'priceMax',
            label: 'Price Max',
            type: 'number'
          },
          {
            name: 'description',
            label: 'Description',
            type: 'textarea'
          },
          {
            name: 'show',
            label: 'Show In Link',
            type: 'switch'
          }
        ]
      },
      {
        name: 'post',
        label: 'Post',
        type: 'widget',
        fields: [
          {
            name: 'title',
            label: 'Title',
            type: 'text'
          },
          {
            name: 'content',
            label: 'Content',
            type: 'textarea'
          },
          {
            name: 'gallery',
            label: 'Gallery',
            type: 'gallery'
          }
        ]
      },
      {
        name: 'tribes',
        label: 'Tribes',
        type: 'widget',
        fields: [
          {
            name: 'address',
            label: 'Tribe address*',
            type: 'text'
          },
          {
            name: 'show',
            label: 'Show In Link',
            type: 'switch'
          }
        ]
      },
      {
        name: 'blog',
        label: 'Blog',
        type: 'widget',
        class: 'blog',
        fields: [
          {
            name: 'title',
            label: 'Title*',
            type: 'text'
          },
          {
            name: 'markdown',
            label: 'Markdown*',
            type: 'textarea'
          },
          {
            name: 'gallery',
            label: 'Gallery',
            type: 'gallery'
          },
          {
            name: 'show',
            label: 'Show In Link',
            type: 'switch'
          }
        ]
      }
    ],
    page: 2
  }
];

export const liquidSchema: FormField[] = [
  {
    name: 'address',
    label: 'Liquid Address',
    type: 'text'
  }
];

export const firstScreenSchema: FormField[] = [
  {
    name: 'img',
    label: 'Image',
    type: 'img',
    page: 1
  },
  {
    name: 'pubkey',
    label: 'Pubkey*',
    type: 'text',
    readOnly: true,
    page: 1
  },
  {
    name: 'owner_alias',
    label: 'Name*',
    type: 'text',
    required: true,
    validator: strValidator,
    page: 1
  },
  {
    name: 'description',
    label: 'Description',
    type: 'textarea',
    page: 1
  },
  {
    name: 'price_to_meet',
    label: 'Price to Meet',
    type: 'number',
    page: 1,
    extraHTML:
      '<p>*This amount applies to users trying to connect within the Sphinx app. Older versions of the app may not support this feature.</p>'
  },
  {
    name: 'twitter',
    label: 'Twitter',
    widget: true,
    type: 'text',
    prepend: '@',
    page: 1
  }
];

export const aboutSchema: FormField[] = [
  {
    name: 'img',
    label: 'Image',
    type: 'img',
    page: 1
  },
  {
    name: 'pubkey',
    label: 'Pubkey*',
    type: 'text',
    readOnly: true,
    page: 1
  },
  {
    name: 'owner_alias',
    label: 'Name*',
    type: 'text',
    required: true,
    validator: strValidator,
    page: 1
  },
  {
    name: 'price_to_meet',
    label: 'Price to Meet',
    type: 'number',
    page: 1,
    extraHTML:
      '<p>*This amount applies to users trying to connect within the Sphinx app. Older versions of the app may not support this feature.</p>'
  },
  {
    name: 'description',
    label: 'Description',
    type: 'textarea',
    page: 1
  },
  {
    name: 'alert',
    label: 'Alerts',
    type: 'switch'
  },
  {
    name: 'tribes',
    label: 'Tribes',
    type: 'multiselect',
    options: [],
    widget: true
  },
  {
    name: 'coding_languages',
    label: 'Coding Languages',
    widget: true,
    type: 'creatablemultiselect',
    options: codingLanguages,
    page: 1
  },
  {
    name: 'github',
    label: 'Github',
    widget: true,
    type: 'text',
    prepend: 'https://github.com/',
    page: 1
  },
  {
    name: 'repos',
    label: 'Github Repository Links',
    widget: true,
    type: 'creatablemultiselect',
    options: [],
    note: 'Enter in this format: ownerName/repoName, (e.g. stakwork/sphinx-tribes).',
    validator: repoArrayStrValidator, // look for 1 slash
    page: 1
  },

  {
    name: 'lightning',
    label: 'Lightning address',
    widget: true,
    type: 'text',
    page: 1
  },
  {
    name: 'liquid',
    label: 'Liquid address',
    widget: true,
    type: 'text',
    page: 1
  },
  {
    name: 'amboss',
    label: 'Amboss address',
    widget: true,
    type: 'text',
    page: 1
  },
  {
    name: 'twitter',
    label: 'Twitter',
    widget: true,
    type: 'text',
    prepend: '@',
    page: 1
  }

  // {
  //     name: "facebook",
  //     label: "Facebook",
  //     widget: true,
  //     type: "text",
  //     page: 1,
  // },
];

export const postSchema: FormField[] = [
  {
    name: 'title',
    label: 'Title',
    type: 'text',
    validator: strValidator
  },
  {
    name: 'content',
    label: 'Content',
    type: 'textarea',
    validator: strValidator
  },
  {
    name: 'gallery',
    label: 'Gallery',
    type: 'gallery'
  }
];

export const sendBadgeSchema: FormField[] = [
  {
    name: 'recipient',
    label: 'Recipient',
    type: 'searchableselect',
    options: [],
    validator: badgeObjectStrValidator
  },
  {
    name: 'badge',
    label: 'Badge',
    type: 'searchableselect',
    options: [],
    validator: badgeObjectStrValidator
  }
  // {
  //     name: "amount",
  //     label: "Amount",
  //     type: "number",
  //     validator: nomValidator,
  // },
];

//name, webhook, price_per_use, img, description, tags

export const botSchema: FormField[] = [
  {
    name: 'img',
    label: 'Logo',
    type: 'imgcanvas'
  },
  {
    name: 'name',
    label: 'Bot Name',
    type: 'text',
    validator: strValidator
  },
  {
    name: 'webhook',
    label: 'Webhook',
    type: 'text',
    validator: strValidator
  },
  {
    name: 'description',
    label: 'How to use',
    type: 'textarea',
    validator: strValidator
  },
  {
    name: 'price_per_use',
    label: 'Price Per Use',
    type: 'number',
    validator: nomValidator
  },
  {
    name: 'tags',
    label: 'Tags',
    type: 'creatablemultiselect',
    options: [
      {
        value: 'Utility',
        label: 'Utility'
      },
      {
        value: 'Social',
        label: 'Social'
      },
      {
        value: 'Fun',
        label: 'Fun'
      },
      {
        value: 'Betting',
        label: 'Betting'
      }
    ]
  }
];

export const offerSkillSchema: FormField[] = [
  {
    name: 'title',
    label: 'Title',
    validator: strValidator,
    type: 'text'
  },
  {
    name: 'description',
    label: 'Description',
    validator: strValidator,
    type: 'textarea'
  },
  {
    name: 'gallery',
    label: 'Gallery',
    type: 'gallery'
  }
];

export const offerOtherSchema: FormField[] = [
  {
    name: 'title',
    label: 'Title',
    validator: strValidator,
    type: 'text'
  },
  {
    name: 'description',
    label: 'Description',
    validator: strValidator,
    type: 'textarea'
  },
  {
    name: 'gallery',
    label: 'Gallery',
    type: 'gallery'
  }
];

export const offerSchema: FormField[] = [
  {
    name: 'dynamicSchema',
    label: 'none',
    type: 'text',
    defaultSchema: offerSkillSchema,
    defaultSchemaName: 'offer_skill',
    dropdownOptions: 'offer',
    // these are included to allow searching by fields for all possible schema types
    dynamicSchemas: [offerSkillSchema, offerOtherSchema]
  }
];

export const wantedOtherSchema: FormField[] = [
  {
    name: 'wanted_type',
    label: 'Tell us about what you need help with',
    type: 'select',
    options: helpType_wanted_other_schema
  },
  {
    name: 'title',
    label: 'Title*',
    type: 'text',
    validator: strValidator
  },

  {
    name: 'one_sentence_summary',
    label: 'One Sentence Summary',
    type: 'text',
    validator: strValidator
  },

  {
    name: 'description',
    label: 'Description',
    type: 'textarea',
    validator: strValidator
  },
  // {
  //   name: 'technology_need_help',
  //   label: 'Technologies you need help with',
  //   type: 'select',
  //   options: helpType_wanted_coding_task_schema
  // },

  {
    name: 'codingLanguage',
    label: 'Technologies you need help with',
    type: 'creatablemultiselect',
    options: codingLanguages
  },
  // {
  //   name: 'priceMin',
  //   label: 'Price Min',
  //   validator: Yup.number().when('priceMax', (pricemax) =>
  //     Yup.number().max(pricemax, `Must be less than max`)
  //   ),
  //   type: 'number'
  // },
  // {
  //   name: 'priceMax',
  //   label: 'Price Max',
  //   validator: nomValidator,
  //   type: 'number'
  // },
  {
    name: 'estimated_budget_for_every_15_minutes',
    label: 'Estimated budget for every 15 minutes',
    type: 'select',
    options: estimated_budget_15_min_options
  },
  {
    name: 'show',
    label: 'Show to public',
    type: 'switch'
  },
  {
    name: 'gallery',
    label: 'Gallery',
    type: 'gallery'
  },

  {
    name: 'type',
    label: 'Type',
    type: 'hide'
  }

  // {
  //     name: 'show',
  //     label: "Show In Link",
  //     type: "switch",
  // },
];

export const wantedCodingTaskSchema: FormField[] = [
  {
    name: 'title',
    label: 'Title',
    type: 'hide'
    // validator: strValidator,
  },
  // {
  //     name: 'repo',
  //     label: "Github Repository",
  //     type: "text",
  //     note: 'Enter in this format: ownerName/repoName, (e.g. stakwork/sphinx-tribes).',
  //     validator: repoStrValidator, // look for 1 slash
  // },
  // {
  //     name: 'issue',
  //     label: "Issue #",
  //     type: "number",
  //     note: 'Add the "stakwork" user to your github repo for issue status updates.',
  //     validator: nomValidator,
  // },
  {
    name: 'wanted_type',
    label: 'Tell us about what you need help with',
    type: 'select',
    options: helpType_wanted_coding_task_schema
  },
  {
    name: 'one_sentence_summary',
    label: 'One Sentence Summary',
    type: 'text',
    validator: strValidator
  },
  {
    name: 'ticketUrl',
    label: 'Github Issue URL',
    type: 'text'
    // validator: strValidator
  },
  {
    name: 'github_description',
    label: 'Get Description From Github ',
    type: 'switch'
  },
  {
    name: 'description',
    label: 'Description',
    type: 'textarea'
    // validator: strValidator,
  },
  {
    name: 'price',
    label: 'Price (Sats)',
    validator: nomValidator,
    type: 'number'
  },
  // {
  //     name: 'space',
  //     label: 'space',
  //     type: 'space'
  // },

  {
    name: 'assignee',
    label: 'Assignee',
    type: 'searchableselect',
    options: []
  },
  {
    name: 'codingLanguage',
    label: 'Coding Language',
    type: 'creatablemultiselect',
    options: codingLanguages
  },
  {
    name: 'tribe',
    label: 'Tribe',
    type: 'select',
    options: [],
    validator: strValidatorNotRequired
  },
  {
    name: 'estimate_session_length',
    label: 'Estimate Seesion Length',
    type: 'select',
    options: estimated_time
  },
  {
    name: 'estimated_completion_date',
    label: 'Estimate Completion Date',
    type: 'date'
  },
  {
    name: 'deliverables',
    label: 'Deliverables',
    type: 'textarea'
  },
  {
    name: 'show',
    label: 'Show to public',
    type: 'switch'
  },
  {
    name: 'loomEmbedUrl',
    label: 'Loom Video',
    type: 'loom',
    validator: strValidatorNotRequired
  },

  {
    name: 'type',
    label: 'Type',
    type: 'hide'
  }
];

export const wantedSchema: FormField[] = [
  {
    name: 'dynamicSchema',
    label: 'none',
    type: 'text',
    defaultSchema: wantedCodingTaskSchema,
    defaultSchemaName: 'freelance_job_request',
    dropdownOptions: 'wanted',
    // these are included to allow searching by fields for all possible schema types
    dynamicSchemas: [wantedCodingTaskSchema, wantedOtherSchema]
  }
];

// this object is used to switch between schemas in form when dynamic
export const dynamicSchemasByType = {
  coding_task: wantedCodingTaskSchema,
  other: wantedOtherSchema,
  //
  freelance_job_request: wantedCodingTaskSchema,
  live_help: wantedOtherSchema,
  offer_skill: offerSkillSchema,
  offer_other: offerOtherSchema
};

// this object is used to autofill form fields if info is available in local storage
export const dynamicSchemaAutofillFieldsByType = {
  freelance_job_request: {
    ticketUrl: 'lastGithubRepo'
  }
};
