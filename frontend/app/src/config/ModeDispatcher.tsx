import React from 'react';

export enum AppMode {
  TRIBES = 'tribes',
  PEOPLE = 'people',
  COMMUNITY = 'community'
}

const hosts: { [k: string]: AppMode } = {
  'localhost:3000': AppMode.TRIBES,
  'localhost:13000': AppMode.TRIBES,
  'localhost:23000': AppMode.TRIBES,
  'tribes.sphinx.chat': AppMode.TRIBES,
  'tribes-test.sphinx.chat': AppMode.TRIBES,
  'localhost:13007': AppMode.COMMUNITY,
  'localhost:23007': AppMode.COMMUNITY,
  'localhost:3007': AppMode.COMMUNITY,
  'people.sphinx.chat': AppMode.COMMUNITY,
  'people-test.sphinx.chat': AppMode.COMMUNITY,
  'community-test.sphinx.chat': AppMode.COMMUNITY,
  'community.sphinx.chat': AppMode.COMMUNITY,
  'bounties.sphinx.chat': AppMode.COMMUNITY
};

function getMode(): AppMode {
  const { host } = window.location;
  return hosts[host] || AppMode.TRIBES;
}

export const ModeDispatcher = ({
  children
}: {
  children: (mode: AppMode) => React.ReactElement;
}) => {
  const mode = getMode();
  return children(mode);
};
