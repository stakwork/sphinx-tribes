import { useState } from "react";
import { uiStore } from "../store/ui";

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
  return { handleScroll, n, loadingMore }
}


export function usePageScroll(goForward, goBackwards) {
  const [loadingBottom, setLoadingBottom] = useState(false);
  const [loadingTop, setLoadingTop] = useState(false);

  function handleScroll(e: any) {
    const bottom = e.target.scrollHeight - e.target.scrollTop == e.target.clientHeight;
    const top = e.target.scrollTop == 0;

    try {
      if (bottom || top) {
        setTimeout(async () => {
          if (bottom) {
            setLoadingBottom(true)
            await goForward()
            setLoadingBottom(false)
          }
          else {
            setLoadingTop(true)
            goBackwards()
            setLoadingTop(false)

          }
        }, 500);
      }
    } catch (e) {
      console.log('oops!', e)
    }
  }
  return { handleScroll, loadingBottom, loadingTop }
}
