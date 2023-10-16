export const BountyDetailsCreationData = {
  step_1: {
    step: 1,
    schemaName: '',
    heading: 'Choose Bounty type',
    sub_heading: '',
    schema: [''],
    schema2: [''],
    required: [''],
    outerContainerStyle: {
      minWidth: '712px',
      maxWidth: '712px',
      height: '560px'
    },
    headingStyle: {},
    extraText: ''
  },
  step_2: {
    step: 2,
    schemaName: 'Freelance Job Request',
    heading: 'Basic info',
    sub_heading: ' ',
    schema: ['org_uuid', 'one_sentence_summary', 'ticket_url'],
    schema2: ['wanted_type', 'coding_languages'],
    required: ['one_sentence_summary', 'wanted_type'],
    outerContainerStyle: {
      minWidth: '712px',
      maxWidth: '712px',
      height: 'auto'
    },
    headingStyle: {},
    extraText: '* Required fields'
  },
  step_3: {
    step: 3,
    schemaName: 'Freelance Job Request',
    heading: 'Description',
    sub_heading: ' ',
    schema: ['github_description', 'description'],
    schema2: [' ', 'loomEmbedUrl'],
    required: [''],
    outerContainerStyle: {
      minWidth: '712px',
      maxWidth: '712px',
      height: 'auto'
    },
    headingStyle: {},
    extraText: '* Required fields'
  },
  step_4: {
    step: 4,
    schemaName: 'Freelance Job Request',
    heading: 'Price and Estimate',
    sub_heading: ' ',
    schema: ['price', 'estimated_session_length', 'estimated_completion_date'],
    schema2: ['tribe', 'deliverables', 'show'],
    required: ['price'],
    outerContainerStyle: {
      minWidth: '712px',
      maxWidth: '712px',
      height: 'auto'
    },
    headingStyle: {},
    extraText: '* Required fields'
  },
  step_5: {
    step: 5,
    schemaName: '',
    heading: 'Assign Developer',
    sub_heading: '',
    schema: ['assignee'],
    schema2: [''],
    required: [''],
    outerContainerStyle: {
      minWidth: '388px',
      maxWidth: '388px',
      height: 'auto'
    },
    headingStyle: {},
    extraText: ''
  }
};
