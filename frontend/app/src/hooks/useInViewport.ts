import { useEffect, useRef, useState } from 'react';
import { IntersectionOptions } from 'react-intersection-observer';

export const useInViewPort = (options: IntersectionOptions) => {
  const [inView, setInView] = useState<boolean>(false);
  const elementRef = useRef<HTMLElement | null>(null);

  const callback = (entries: IntersectionObserverEntry[]) => {
    const [entry] = entries;
    setInView(entry.isIntersecting);
  };

  useEffect(() => {
    // checking support and blocking execution on server-side
    if (!(typeof window !== 'undefined') && !('IntersectionObserver' in window)) return;

    const observer = new IntersectionObserver(callback, options);
    if (elementRef.current) observer.observe(elementRef.current);

    return () => observer.disconnect();
  }, [options]);

  return [inView, elementRef] as const;
};
