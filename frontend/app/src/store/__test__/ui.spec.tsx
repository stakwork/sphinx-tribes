import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import { UiStore } from '../ui.ts';

describe('UiStore', () => {
  let uiStore: UiStore;
  let container: HTMLElement;

  beforeEach(() => {
    uiStore = new UiStore();
    container = document.createElement('div');
    document.body.appendChild(container);
  });

  afterEach(() => {
    document.body.removeChild(container);
    container = null!;
  });

  it('clears search text when clicking outside of search input', () => {
    const { getByTestId } = render(
      <div>
        <input
          data-testid="search-input"
          value={uiStore.searchText}
          onChange={(e) => uiStore.setSearchText(e.target.value)}
        />
        <button data-testid="outside-button">Outside</button>
      </div>
    );

    const searchInput = getByTestId('search-input') as HTMLInputElement;
    const outsideButton = getByTestId('outside-button') as HTMLButtonElement;

    fireEvent.change(searchInput, { target: { value: 'Test' } });

    expect(uiStore.searchText).toBe('test');

    fireEvent.click(outsideButton);

    expect(uiStore.searchText).toBe('');
  });
});
