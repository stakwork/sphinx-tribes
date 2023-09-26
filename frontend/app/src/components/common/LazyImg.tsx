import { noop } from 'lodash';
import React, { HTMLAttributes, ImgHTMLAttributes, useEffect, useRef, useState } from 'react';
import styled from 'styled-components';

const Placeholder = styled.img`
  filter: blur(1rem);
`;

const Img = styled.img<{ loaded: boolean }>`
  filter: ${({ loaded }: any) => (loaded ? '"blur(0)"' : 'blur(1rem)')};
`;

export const LazyImg = ({
  placeholder = '/static/person_placeholder.png',
  ...props
}: ImgHTMLAttributes<HTMLImageElement> & { placeholder?: string }) => {
  const [inView, setInView] = useState(false);
  const [imgLoaded, setImgLoaded] = useState(false);

  const placeholderRef = useRef<HTMLImageElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries: IntersectionObserverEntry[], obs: IntersectionObserver) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            setInView(true);
            obs.disconnect();
          }
        }
      },
      {}
    );

    placeholderRef.current && observer.observe(placeholderRef.current);
    return () => {
      observer.disconnect();
    };
  }, []);

  return inView ? (
    <Img loaded={imgLoaded} {...props} alt={props.alt || ''} onLoad={() => setImgLoaded(true)} />
  ) : (
    <Placeholder {...props} ref={placeholderRef} src={placeholder} alt={props.alt || ''} />
  );
};

interface ImageProps {
  src: string;
  loaded: boolean;
}
const ImgDiv = styled.div<ImageProps>`
  background-image: url('${(p: any) => p.src}');
  filter: ${({ loaded }: any) => (loaded ? '"blur(0)"' : 'blur(1rem)')};
`;

export const LazyImgBg = (
  props: HTMLAttributes<HTMLDivElement> & { placeholder?: string; src: string }
) => {
  const [inView, setInView] = useState(false);
  const [imgLoaded, setImgLoaded] = useState(false);

  const placeholderRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries: IntersectionObserverEntry[], obs: IntersectionObserver) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            setInView(true);
            obs.disconnect();
          }
        }
      },
      {}
    );

    placeholderRef.current && observer.observe(placeholderRef.current);
    return () => {
      observer.disconnect();
    };
  }, []);

  useEffect(() => {
    if (inView) {
      const image = new Image();
      image.onload = () => {
        setImgLoaded(true);
      };
      image.src = props.src;
      return () => (image.onload = null);
    }
    return noop;
  }, [inView, props.src]);

  return inView ? (
    <ImgDiv loaded={imgLoaded} {...props} />
  ) : (
    <div {...props} ref={placeholderRef} />
  );
};
