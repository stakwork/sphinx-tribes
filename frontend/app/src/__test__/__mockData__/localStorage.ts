export const localStorageMock = (function () {
  let store = {};

  return {
    getItem(key: string | number) {
      return store[key];
    },

    setItem(key: string | number, value: any) {
      store[key] = value;
    },

    clear() {
      store = {};
    },

    removeItem(key: string | number) {
      delete store[key];
    },

    getAll() {
      return store;
    }
  };
})();

Object.defineProperty(window, 'localStorage', { value: localStorageMock });
