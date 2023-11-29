import { localStorageMock } from '../__test__/__mockData__/localStorage';
describe('Local Storage', () => {
  beforeEach(() => {
    Object.defineProperty(window, 'localStorage', { value: localStorageMock });
  });

  afterEach(() => {
    localStorageMock.clear();
  });

  it('should set and get an item from local storage', () => {
    const key = 'testKey';
    const value = 'testValue';

    localStorageMock.setItem(key, value);

    expect(localStorageMock.getItem(key)).toEqual(value);
  });

  it('should remove an item from local storage', () => {
    const key = 'testKey';
    const value = 'testValue';

    localStorageMock.setItem(key, value);
    localStorageMock.removeItem(key);

    expect(localStorageMock.getItem(key)).toBeUndefined();
  });

  it('should clear all items from local storage', () => {
    const key1 = 'testKey1';
    const value1 = 'testValue1';
    const key2 = 'testKey2';
    const value2 = 'testValue2';

    localStorageMock.setItem(key1, value1);
    localStorageMock.setItem(key2, value2);
    localStorageMock.clear();

    expect(localStorageMock.getItem(key1)).toBeUndefined();
    expect(localStorageMock.getItem(key2)).toBeUndefined();
  });
});
