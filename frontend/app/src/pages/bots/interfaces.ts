export interface BotProps {
  name: string;
  hideActions: React.CSSProperties;
  small: boolean;
  id: number;
  img: string;
  description: any;
  selected: boolean;
  select: (number, string) => void;
  unique_name: string;
}

export interface BotViewProps {
  botUniqueName: string;
  selectBot: (any) => void;
  loading: boolean;
  goBack: () => void;
  botView?: boolean;
}

export interface BotSecretProps {
  id: string;
  secret: string;
  name: string;
  full: string;
}
