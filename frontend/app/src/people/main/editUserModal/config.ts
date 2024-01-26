import { aboutSchema } from 'components/form/schema';

export const formConfig = {
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
};
