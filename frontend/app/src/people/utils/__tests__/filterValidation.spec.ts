import bounty from 'components/form/bounty';
import { 
    bountyHeaderFilter, 
    bountyHeaderLanguageFilter 
} from '../filterValidation';

beforeAll(() => {
    
});

afterAll(() => {
});

describe('testing helpers', () => {
    describe('bountyHeaderFilter', () => {

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
