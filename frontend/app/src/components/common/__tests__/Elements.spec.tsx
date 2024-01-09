import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../__test__/__mockData__/setupStore';
import { user } from '../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../__test__/__mockFn__/useHistory';
import { Title, Date, Paragraph, Link } from '../Elements';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

/**
 * @jest-environment jsdom
 */
describe('Elements', () => {
  nock(user.url).get('/person/id/1').reply(200, {});
  test('Title', async () => {
    const TitleText = 'Title Text';
    render(<Title>{TitleText}</Title>);

    expect(screen.queryByText(TitleText)).toBeInTheDocument();
  });
  test('Date', async () => {
    const DateText = 'Date Text';
    render(<Date>{DateText}</Date>);

    expect(screen.queryByText(DateText)).toBeInTheDocument();
  });
  test('Paragraph', async () => {
    const ParagraphText = 'Paragraph Text';
    render(<Paragraph>{ParagraphText}</Paragraph>);

    expect(screen.queryByText(ParagraphText)).toBeInTheDocument();
  });
  test('Link', async () => {
    const LinkText = 'Link Text';
    render(<Link>{LinkText}</Link>);

    expect(screen.queryByText(LinkText)).toBeInTheDocument();
  });
});
