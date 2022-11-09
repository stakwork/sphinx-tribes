const palette = {
  blue1: '#618AFF',
  borderBlue1: '#5078F2',
  textBlue1: '#A3C1FF',
  text1: '#292C33',
  text2: '#3C3F41',
  text2_4: '#8E969C',
  green1: '#49C998',
  borderGreen1: '#2FB379',
  red1: '#ED7474',
  red2: '#FF8F80',
  header: '#1A242E',
  background: '#151E27',
  divider1: '#151E27',
  divider2: '#101317',
  pureWhite: '#fff',
  pureBlack: '#000',
  black100: 'rgba(0, 0, 0, 0.07)',
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
    G10: '#3C3F41',
    G50: '#5F6368',
    G100: '#8E969C',
    G200: '#909BAA',
    G250: '#9AAEC6',
    G300: '#B0B7BC',
    G400: '#bac1c6',
    G500: '#d0d5d8',
    G600: '#DDE1E5',
    G700: '#EBEDEF',
    G800: '#f7f8f8',
    G900: '#F0F1F2',
    G950: '#F2F3F5'
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

const colors: { [key in PalletType]: any } = {
  dark: {
    ...palette
  },
  light: {
    ...palette,
    background: '#ffffff'
  }
};

export { colors };
