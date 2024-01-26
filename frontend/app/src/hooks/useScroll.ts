/* eslint-disable func-style */
import { useState } from 'react';

export function useScroll() {
  const [loadingMore, setLoadingMore] = useState(false);
  const [n, setN] = useState(100);
  function handleScroll(e: any) {
    const bottom = e.target.scrollHeight - e.target.scrollTop <= e.target.clientHeight;
    if (bottom) {
      setLoadingMore(true);
      setTimeout(() => {
        setN(n + 100);
      }, 500);
      setTimeout(() => {
        setLoadingMore(false);
      }, 3000);
    }
  }
  return { handleScroll, n, loadingMore };
}

export function usePageScroll(goForward: any, goBackwards?: any) {
  const [loadingBottom, setLoadingBottom] = useState(false);

  async function handleScroll(e: any) {
    const bottom = e.target.scrollHeight - e.target.scrollTop === e.target.clientHeight;

    try {
      if (bottom) {
        setLoadingBottom(true);
        await goForward();
        setLoadingBottom(false);
      }
    } catch (e) {
      console.log('oops!', e);
    }
  }
  return { handleScroll, loadingBottom, loadingTop: false };
}
