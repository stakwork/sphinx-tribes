type LocType = 'English' | 'Spanish' | 'Russian';

const loc: { [key in LocType]: any } = {
  English: {
    tribes: 'Tribes',
    people: 'People',
    bots: 'Bots',
    portfolios: 'Portfolios',
    tickets: 'Tickets',

    startYourOwnProfile: 'Start your own profile',
    getStarted: 'Get Started',
    getSphinx: 'Get Sphinx',
    signIn: 'Sign in'
  },
  Spanish: {},
  Russian: {}
};

const languageOptions = ['English', 'Spanish', 'Russian'];

export { loc, languageOptions };
