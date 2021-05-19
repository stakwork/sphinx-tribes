import { useState } from "react";

export function useScroll() {
  const [loadingMore, setLoadingMore] = useState(false);
  const [n, setN] = useState(100);
  function handleScroll(e: any) {
    const bottom =
      e.target.scrollHeight - e.target.scrollTop <= e.target.clientHeight;
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
  return {handleScroll, n, loadingMore}
}
