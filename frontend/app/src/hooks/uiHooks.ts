import { useState, useEffect } from 'react';
import { mobileWidht } from '../config';

function getIsMobile() {
  return window.innerWidth < mobileWidht;
}

const screenWidthOffset = 36;
const screenHeightOffset = 36;

function getScreenWidth() {
  return window.innerWidth - screenWidthOffset;
}

function getScreenHeight() {
  return window.innerHeight - screenHeightOffset;
}

function useIsMobile() {
  const [isMobile, setIsMobile] = useState(getIsMobile());

  useEffect(() => {
    function handleResize() {
      setIsMobile(getIsMobile());
    }

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  return isMobile;
}

function useScreenWidth() {
  const [width, setWidth] = useState(getScreenWidth());

  useEffect(() => {
    function handleResize2() {
      setWidth(getScreenWidth());
    }

    window.addEventListener('resize', handleResize2);
    return () => window.removeEventListener('resize', handleResize2);
  }, []);

  return width;
}

function useScreenHeight() {
  const [h, setH] = useState(getScreenHeight());

  useEffect(() => {
    function handleResize3() {
      setH(getScreenHeight());
    }

    window.addEventListener('resize', handleResize3);
    return () => window.removeEventListener('resize', handleResize3);
  }, []);

  return h;
}

export {
  useIsMobile,
  useScreenWidth,
  getScreenWidth,
  useScreenHeight,
  getIsMobile,
  screenWidthOffset
};
