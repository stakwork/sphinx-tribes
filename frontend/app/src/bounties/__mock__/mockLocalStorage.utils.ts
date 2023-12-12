const mockLocalStorage = (() => {
  let store = {};

  return {
    getItem(key: string) {
      return store[key] || null;
    },

    setItem(key: string, value: string) {
      store[key] = typeof value === 'string' ? value : JSON.stringify(value);
    },

    removeItem(key: string) {
      delete store[key];
    },

    clear() {
      store = {};
    },

    getAll() {
      return store;
    }
  };
})();

Object.defineProperty(window, 'localStorage', { value: mockLocalStorage });

export default mockLocalStorage;
