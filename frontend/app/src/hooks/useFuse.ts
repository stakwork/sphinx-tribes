import Fuse from 'fuse.js';
import { useStores } from '../store';

const fuseOptions = {
  keys: ['name', 'description'],
  shouldSort: true,
  // matchAllTokens: true,
  includeMatches: true,
  threshold: 0.35,
  location: 0,
  distance: 100,
  maxPatternLength: 32,
  minMatchCharLength: 1
};

export function useFuse(array, keys: string[] = []) {
  const { ui } = useStores();
  let theArray = array;

  if (ui.searchText) {
    const fuse = new Fuse(array, { ...fuseOptions, keys });
    const res = fuse.search(ui.searchText);
    theArray = res.map((r) => r.item);
  }

  return theArray;
}

export function useLocalFuse(searchText, array, keys: string[] = []) {
  let theArray = array;
  if (searchText) {
    const fuse = new Fuse(array, { ...fuseOptions, keys });
    const res = fuse.search(searchText);
    theArray = res.map((r) => r.item);
  }
  return theArray;
}
