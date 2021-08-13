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
}

type PalletType = "dark" | "light"

const colors: { [key in PalletType]: any } = {
    dark: {
        ...palette,
    },
    light: {
        ...palette,
        background: '#ffffff',
    }
}

export { colors }