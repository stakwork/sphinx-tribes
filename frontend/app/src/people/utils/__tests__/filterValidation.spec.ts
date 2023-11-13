import bounty from 'components/form/bounty';
import { 
    bountyHeaderFilter, 
    bountyHeaderLanguageFilter 
} from '../filterValidation';

describe('testing helpers', () => {
    describe('bountyHeaderFilter', () => {
        test('o/t/t', () => {
            expect(bountyHeaderFilter({ Open: true }, true, true)).toEqual(false);
        });
        test('a/t/t', () => {
            expect(bountyHeaderFilter({ Assigned: true }, true, true)).toEqual(true);
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
            expect(bountyHeaderLanguageFilter(['Javascript', 'Python'], {'Javascript': true})).toEqual(true);
        });
        test('no-match', () => {
            expect(bountyHeaderLanguageFilter(['Javascript'], {'Python': true, 'Javascript': false})).toEqual(false);
        });
        test('no filters', () => {
            expect(bountyHeaderLanguageFilter(['Javascript'], {})).toEqual(false);
        });
        test('no languages', () => {
            expect(bountyHeaderLanguageFilter([], {'Javascript': true})).toEqual(false);
        });
        test('false filters', () => {
            expect(bountyHeaderLanguageFilter(['Javascript'], {'Javascript': false, 'Python': false})).toEqual(true);
        });
    });
});
