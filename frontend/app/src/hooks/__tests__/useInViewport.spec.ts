import { renderHook } from '@testing-library/react-hooks';
import { useInViewPort } from 'hooks/useInViewport';
import { IntersectionOptions } from 'react-intersection-observer';

const defaultOptions: IntersectionOptions = {};

describe('useInViewport hook', () => {
  interface Window {
    IntersectionObserver: unknown;
  }

  const prepareForTesting = (isIntersecting: boolean): void => {
    const mockedIntersectionObserver = jest.fn((callback: any) => {
      callback([{ isIntersecting }]);

      return {
        observe: jest.fn(),
        unobserve: jest.fn()
      };
    });

    (window as Window).IntersectionObserver = mockedIntersectionObserver;
  };

  it('should initialize the inView state to false', () => {
    prepareForTesting(false);

    const { result } = renderHook(() => useInViewPort(defaultOptions));
    expect(result.current[0]).toBe(false);
  });

  it('should not detect element when reference is not connected', () => {
    prepareForTesting(false);

    const { result } = renderHook(() => useInViewPort(defaultOptions));
    const [isVisible, ref] = result.current;

    expect(ref.current).toBe(null);
    expect(isVisible).toBe(false);
  });

  it('should detect element after scroll', () => {
    prepareForTesting(true);

    const element = document.createElement('div');

    const { result } = renderHook(() => useInViewPort(defaultOptions));

    const [isVisible, ref] = result.current;

    ref.current = element;

    renderHook(() => useInViewPort(defaultOptions));

    expect(isVisible).toBe(true);
  });
});
