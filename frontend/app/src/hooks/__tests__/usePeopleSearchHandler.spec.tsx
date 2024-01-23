import { mainStore } from 'store/main';
import { uiStore } from 'store/ui';
import { usePeopleSearchHandler } from 'hooks';
import { act, renderHook } from '@testing-library/react-hooks';

jest.mock('store/ui');
jest.mock('store/main');

describe('People Search Handler Hook', () => {
  test('Handler should pass his argument to `ui.setSearchText` method and call `main.getPeople` method', async () => {
    await act(async () => {
      const { result } = renderHook(() => usePeopleSearchHandler());
      const handleSearchChange = result.current;

      handleSearchChange('Y');

      expect(uiStore.setSearchText).toBeCalledWith('Y');

      expect(mainStore.getPeople).toBeCalledTimes(1);
    });
  });
});
