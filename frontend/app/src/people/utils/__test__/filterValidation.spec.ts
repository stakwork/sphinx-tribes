import { bountyHeaderFilter, bountyHeaderLanguageFilter } from '../filterValidation';
import filterByCodingLanguage from '../filterPeople';
import { users } from '../__test__/__mockData__/users';

describe('testing filters', () => {
  describe('bountyHeaderFilter', () => {
    test('o/t/t', () => {
      expect(bountyHeaderFilter({ Open: true }, true, true)).toEqual(false);
    });
    test('a/t/t', () => {
      expect(bountyHeaderFilter({ Assigned: true }, true, true)).toEqual(false);
    });
    test('p/t/t', () => {
      expect(bountyHeaderFilter({ Paid: true }, true, true)).toEqual(true);
    });
    test('/t/t', () => {
      expect(bountyHeaderFilter({}, true, true)).toEqual(true);
    });
    test('o/f/t', () => {
      expect(bountyHeaderFilter({ Open: true }, false, true)).toEqual(false);
    });
    test('a/f/t', () => {
      expect(bountyHeaderFilter({ Assigned: true }, false, true)).toEqual(true);
    });
    test('p/f/t', () => {
      expect(bountyHeaderFilter({ Paid: true }, false, true)).toEqual(false);
    });
  });
  describe('bountyHeaderLanguageFilter', () => {
    test('match', () => {
      expect(bountyHeaderLanguageFilter(['Javascript', 'Python'], { Javascript: true })).toEqual(
        true
      );
    });
    test('no-match', () => {
      expect(
        bountyHeaderLanguageFilter(['Javascript'], { Python: true, Javascript: false })
      ).toEqual(false);
    });
    test('no filters', () => {
      expect(bountyHeaderLanguageFilter(['Javascript'], {})).toEqual(true);
    });
    test('no languages', () => {
      expect(bountyHeaderLanguageFilter([], { Javascript: true })).toEqual(false);
    });
    test('false filters', () => {
      expect(
        bountyHeaderLanguageFilter(['Javascript'], { Javascript: false, Python: false })
      ).toEqual(true);
    });
  });
  describe('peopleHeaderCodingLanguageFilters', () => {
    test('match', () => {
      expect(filterByCodingLanguage(users, { Typescript: true })).toStrictEqual([users[0]]);
    });
    test('no_match', () => {
      expect(filterByCodingLanguage(users, { Rust: true })).toStrictEqual([]);
    });
    test('no filters', () => {
      expect(filterByCodingLanguage(users, {})).toEqual(users);
    });
    test('false filters', () => {
      expect(filterByCodingLanguage(users, { PHP: false, MySQL: false })).toStrictEqual(users);
    });
    test('no users', () => {
      expect(filterByCodingLanguage([], { Typescript: true })).toStrictEqual([]);
    });
  });
});
