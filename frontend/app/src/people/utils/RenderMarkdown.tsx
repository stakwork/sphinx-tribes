import React from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeRaw from 'rehype-raw';
import { colors } from '../../config/colors';

export function renderMarkdown(markdown: any) {
  const color = colors['light'];
  return (
    <ReactMarkdown
      children={markdown}
      remarkPlugins={[remarkGfm]}
      rehypePlugins={[rehypeRaw]}
      components={{
        code({ className, children, ...props }: any) {
          return (
            <code className={className} {...props}>
              {children}
            </code>
          );
        },
        img({ className, ...props }: any) {
          return (
            <img
              alt={'Markodown'}
              className={className}
              style={{
                width: '100%',
                maxHeight: '100%',
                border: `1px solid ${color.black80}`,
                borderRadius: '4px'
              }}
              {...props}
            />
          );
        }
      }}
    />
  );
}
