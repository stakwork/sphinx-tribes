const palette = {
  blue1: '#618AFF',
  blue2: '#82b4ff',
  blue3: '#45b9f6',
  blue4: '#3C3D3F',
  borderBlue1: '#5078F2',
  light_blue100: '#a3c1ff',
  light_blue200: 'rgba(130, 180, 255, 0.25)',
  textBlue1: '#A3C1FF',
  text1: '#292C33',
  text2: '#3C3F41',
  text2_4: '#8E969C',
  border_image: '#6b7a8d',
  green1: '#49C998',
  borderGreen1: '#2FB379',
  borderGreen2: 'rgba(73, 201, 152, 0.2)',
  red1: '#ED7474',
  red2: '#FF8F80',
  red3: '#b75858',
  header: '#1A242E',
  background: '#151E27',
  divider1: '#151E27',
  divider2: '#101317',
  pureWhite: '#fff',
  white100: '#F5F6F8',
  pureBlack: '#000',
  black80: 'rgba(0, 0, 0, 0.1)',
  black85: 'rgba(0, 0, 0, 0.2)',
  black90: 'rgba(0, 0, 0, 0.25)',
  black100: 'rgba(0, 0, 0, 0.07)',
  black150: 'rgba(0, 0, 0, 0.75)',
  black200: '#272727',
  black300: '#202020',
  black400: '#222E3A',
  black500: '#3c3f41',
  background100: '#f0f1f3',
  statusAssigned: '#49C998',
  statusCompleted: '#8256D0',
  button_primary: {
    main: '#49C998',
    hover: '#3CBE88',
    active: '#2FB379',
    shadow: 'rgba(73, 201, 152, 0.5)'
  },
  button_secondary: {
    main: '#618AFF',
    hover: '#5881F8',
    active: '#5078F2',
    shadow: 'rgba(97, 138, 255, 0.5)'
  },
  grayish: {
    G05: '#292c33',
    G06: 'rgba(176, 183, 188, 0.1)',
    G07: 'rgb(60, 63, 65)',
    G10: '#3C3F41',
    G20: '#555555',
    G25: '#444851',
    G50: '#5F6368',
    G60: '#ffffff44',
    G65: '#f2f3f580',
    G60A: '#888888',
    G70: '#909090',
    G71: '#999999',
    G71A: '#f2f3f580',
    G100: '#8E969C',
    G200: '#909BAA',
    G250: '#9AAEC6',
    G300: '#B0B7BC',
    G400: '#bac1c6',
    G500: '#d0d5d8',
    G600: '#DDE1E5',
    G700: '#EBEDEF',
    G750: '#e0e0e0',
    G760: '#ddd',
    G800: '#f7f8f8',
    G900: '#F0F1F2',
    G950: '#F2F3F5',
    G1000: '#cfcfcf',
    G1100: 'rgb(104, 104, 79)'
  },
  primaryColor: {
    P100: 'rgba(73, 201, 152, 0.15)',
    P200: 'rgba(73,201, 152, 0.2)',
    P300: '#2F7460',
    P400: '#86d9b9',
    P700: '#e0f7f0'
  },
  tribesBackground: '#212539'
};

type PalletType = 'dark' | 'light';

const colors: { [key in PalletType]: typeof palette } = {
  dark: {
    ...palette
  },
  light: {
    ...palette,
    background: '#ffffff'
  }
};

export { colors };
