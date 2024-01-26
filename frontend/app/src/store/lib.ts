import { random } from 'lodash';

const getUserImagePlaceholderPath = (number: number) =>
  `/static/avatarPlaceholders/placeholder_${number}.jpg`;

type UserPlaceholdersCache = Record<string, number>;
export const getUserAvatarPlaceholder = (userId: string) => {
  // caching to show the same avatars if the user has already seen them
  const storageCacheKey = 'userPlaceholdersCache';
  const cache: UserPlaceholdersCache = JSON.parse(localStorage.getItem(storageCacheKey) ?? '{}');
  if (cache[userId]) {
    return getUserImagePlaceholderPath(cache[userId]);
  }
  const randomIndex = random(1, 39, false);
  cache[userId] = randomIndex;
  localStorage.setItem(storageCacheKey, JSON.stringify(cache));

  return getUserImagePlaceholderPath(cache[userId]);
};
