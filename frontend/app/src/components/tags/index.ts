import {ReactComponent as btc} from './svg/btc.svg'
import {ReactComponent as lightning} from './svg/lightning.svg'
import {ReactComponent as sphinx} from './svg/sphinx.svg'
import {ReactComponent as crypto} from './svg/crypto.svg'
import {ReactComponent as music} from './svg/music.svg'
import {ReactComponent as tech} from './svg/tech.svg'
import {ReactComponent as altcoins} from './svg/altcoins.svg'

const tags:{[k:string]:any} = {
    BTC: {
        icon: btc,
        color: '#FAC917'
    },
    Lightning: {
        icon: lightning,
        color: '#9f5bca'
    },
    Sphinx: {
        icon: sphinx,
        color: '#6189ff',
    },
    Crypto: {
        icon: crypto,
        color: '#51ae95',
    },
    Music: {
        icon: music,
        color: '#81c12c',
    },
    Tech: {
        icon: tech,
        color: '#c1501f',
    },
    Altcoins: {
        icon: altcoins,
        color: '#cccccc'
    }
}

export default tags