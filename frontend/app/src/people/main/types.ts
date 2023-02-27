export type Widget =
  | 'people'
  | 'wanted'
  | 'post'
  | 'offer'
  | 'badges'
  | 'about'
  | 'twitter'
  | 'supportme'
  | 'blog';
export type PeopleBodyProps = { selectedWidget: Widget };
