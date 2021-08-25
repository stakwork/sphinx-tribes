import { useState, useEffect } from 'react';

function getIsMobile() {
    return window.innerWidth < 500
}

const screenWidthOffset = 36

function getScreenWidth() {
    return window.innerWidth - screenWidthOffset
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

export { useIsMobile, useScreenWidth, getScreenWidth, getIsMobile, screenWidthOffset }