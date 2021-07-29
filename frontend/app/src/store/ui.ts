import { observable, action } from 'mobx'
import tags from '../tribes/tags'
import { Extras } from '../form/inputs/widgets/interfaces'

const tagLabels = Object.keys(tags)
const initialTags = tagLabels.map(label => {
  return <EuiSelectableOption>{ label }
})

const david = {
  "id": 0,
  "pubkey": "0350587f325dcd6eb50b1c86874961c134be3ab2b9297d88e61443bb0531d7798e",
  "contact_key": "MIIBCgKCAQEAy89ezhWOZyeAyWRepJjIFb5mZo/ggx2iN5FXDWgYPw4hPik/87Al8R7TsbU73JlxfwJkxBKs4gBF0YrydE5sZ9Tmjlp2Jcm9MPMBXVSqyIHTBYRX+5whv6TilYjz3PVoDDzVPfSxJcd4nYLCuW6l0w7fzDnDUHqI4uV8DrQUNO/qpaR73LqylNNNn9xKk9xHYDTs+gqYQ1Oe0Cvxlr5RxQTZxySxcD7HNVoJNnNzqSfq6y1V2oJUu2Lc4hgAhwfMO6foZaYrItsnnvdS0+lEeDqjQjFoMF/4Wwvp7pSYI3N9SdljCM0TH1t1Y7P3y5KRbyJ44MxRRBgf5D1GSR3IOwIDAQAB",
  "alias": "David",
  "photo_url": "",
  "route_hint": "",
  "price_to_meet": 1,
  "jwt": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJyZWxheSIsInB1YmtleSI6IjAzNTA1ODdmMzI1ZGNkNmViNTBiMWM4Njg3NDk2MWMxMzRiZTNhYjJiOTI5N2Q4OGU2MTQ0M2JiMDUzMWQ3Nzk4ZSIsInNjb3BlIjoicGVyc29uYWwiLCJqdGkiOiIyMmE5MWI0NC03NzUyLTQ4MDAtODFhZC00ZDE3MTNlNzIzNTYiLCJpYXQiOjE2Mjc1MTY4OTAsImV4cCI6MTYyNzUxNzE5MH0.OWyMRIqm-JLu1pV_WrZc_cQoo4tW34-7R3pAZ9cg0Qk",
  "url": "https://ecs-relay-arm-3-relay-b6f5b88cbc8bc3a30200.022.sphinxnodes.chat",
  "description": "",
  "verification_signature": "IOdouJWEPi60Hsb2O_EfdBoA7UohlDNFDgLcPl8lpIodeIVvLKKbpZ2vmHE98Q1vfyjHh3SsGhw86TBeG9n2-jg=",
  "extras": {}
}

export type EuiSelectableOptionCheckedType = 'on' | 'off' | undefined;

export interface EuiSelectableOption {
  label: string;
  checked?: EuiSelectableOptionCheckedType;
}

class UiStore {
  @observable ready: boolean = false
  @action setReady(ready: boolean) {
    this.ready = ready
  }

  @observable tags: EuiSelectableOption[] = initialTags
  @action setTags(t: EuiSelectableOption[]) {
    this.tags = t
  }

  @observable searchText: string = ''
  @action setSearchText(s: string) {
    this.searchText = s
  }

  @observable editMe: boolean = false
  @action setEditMe(b: boolean) {
    this.editMe = b
  }

  @observable meInfo: MeData = david
  @action setMeInfo(t: MeData) {
    this.meInfo = t
  }
}

export type MeData = MeInfo | null

export interface MeInfo {
  id?: number
  pubkey: string
  photo_url: string
  alias: string
  route_hint: string
  contact_key: string
  price_to_meet: number
  jwt: string
  url: string
  description: string
  verification_signature: string
  extras: Extras
}
const emptyMeData: MeData = { pubkey: 'asdf', alias: 'evan', route_hint: '', contact_key: '', price_to_meet: 0, photo_url: '', url: '', jwt: '', description: '', verification_signature: '', extras: {} }

export const uiStore = new UiStore()