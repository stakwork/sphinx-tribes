import { observable, action } from 'mobx'

export class MainStore {
  @observable
  tribes: Tribe[] = tribs
}

export const mainStore = new MainStore()

export interface Tribe {
  uuid: string;
  name: string;
  owner: string;
  pubkey: string; // group encryption key
  price: number;
  img: string;
  tags: string[];
  description: string;
  matchCount?: number; // for tag search
}
const tribs: Tribe[] = [
  {
    uuid:'XqzlKB-h8IWHQ1fx2x0yCkcWW2zbmynWfNREJz7nZiWeLUhWrUfYjMKHlRxGCEa6p7VZdmyHw6UWEoYiTai_nt11kZWr',
    name:'Sphinx Chat ',
    owner:'asdf',
    tags:['Sphinx','Lightning','BTC'],
    pubkey:'asdf',
    price:10,
    description:'Join this chat to get help or talk about Sphinx Chat apps and services!',
    img:'https://sphinx.chat/img/Sphinx_icon_1024.png'
  },
  {
    uuid:'aqzlKB-7nZiWeLUhWrUfYjMKHlRxGCEa6p7VZdmyHw6UWEoYiTai_nt11kZWr',
    name:`Somebody's cool crypto chat`,
    owner:'asdf',
    pubkey:'asdf',
    price:10,
    img:'https://bitcoinexchangeguide.com/wp-content/uploads/2018/07/2018-Top-Artificial-Intelligence-Based-Crypto-Platforms-696x449.jpg',
    tags:['BTC','Crypto'],
    description:`Cryptocurrencies are revolutionizing the way we think about value in the moden world. Chat about bitcoin and the whole cryptosphere here`,
  },
  {
    uuid:'asdf-7nZiWeLUhadsfYjMKHlRxGCEaasdfZdmyHw6UWEoYiTai_nt11kZWr',
    name:`Lightning Network`,
    owner:'asdf',
    pubkey:'asdf',
    price:10,
    img:'https://news.bitcoin.com/wp-content/uploads/2020/02/shutterstock_1567394578.png',
    tags:['Lightning','BTC','Crypto'],
    description:`The lightning network is here!!! This chat room focuses on Lightning technology, services, and developments. Join this chat room to talk about new Lightning features and implementations.`,
  },
  {
    uuid:'asdf-7nZiWeLUasdfYjMKHlRxGCEaasdfZdmyHw6UWEoYiTai_nt11kZWr',
    name:`New Digital Tech Chat`,
    owner:'asdf',
    pubkey:'asdf',
    price:10,
    img:'https://images.idgesg.net/images/article/2019/05/cso_best_security_software_best_ideas_best_technology_lightbulb_on_horizon_of_circuit_board_landscape_with_abstract_digital_connective_technology_atmosphere_ideas_innovation_creativity_by_peshkov_gettyimages-965785212_3x2_2400x1600-100797318-large.jpg',
    tags:['Tech','Crypto'],
    description:`Chat about new decentralized technologies in cloud computing, IoT, cryptography, internet, etc.`,
  },
  {
    uuid:'asdf-7nZiWeUasdfYjMKHlRxGCEaasdfZdmyHw6UWEoYiTai_nt11kZWr',
    name:`Altcoin corner`,
    owner:'asdf',
    pubkey:'asdf',
    price:10,
    img:'https://investorplace.com/wp-content/uploads/2020/01/cryptocurrencies1600-768x432.jpg',
    tags:['Altcoins','Crypto'],
    description:`Chat about coins besides BTC! The cryptosphere is full of amazing projects, but also full of scams! Lets pick them apart...`,
  },
  {
    uuid:'asdf-7nZiWasdfeUasdfYjMKHlRxGCEaasdfZdmyHw6UWEoYiTai_nt11kZWr',
    name:`New Music`,
    owner:'asdf',
    pubkey:'asdf',
    price:10,
    img:'https://cdn2.vectorstock.com/i/1000x1000/10/86/music-equaliser-wave-vector-171086.jpg',
    tags:['Music'],
    description:`What are you listening to these days?`,
  }
]

const Tags = {
  element: {
    type: {
      
    }
  }
}
