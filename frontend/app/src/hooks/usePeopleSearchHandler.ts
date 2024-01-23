import { useStores } from 'store';

export function usePeopleSearchHandler() {
  const { ui, main } = useStores();

  return (newSearchText: string) => {
    ui.setSearchText(newSearchText);

    main.getPeople({ page: 1, resetPage: true });
  };
}
