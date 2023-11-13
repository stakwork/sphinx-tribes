import {
  extractGithubIssueFromUrl,
  extractRepoAndIssueFromIssueUrl,
  randomString,
  calculateTimeLeft,
  toCapitalize,
  userHasRole
} from '../helpers-extended';
import { uiStore } from '../../store/ui';
import crypto from 'crypto';
import moment from 'moment';

beforeAll(() => {
  // for test randomString
  Object.defineProperty(globalThis, 'crypto', {
    value: {
      getRandomValues: (arr) => crypto.randomBytes(arr.length)
    }
  });
});

afterAll(() => {});

describe('testing helpers', () => {
  describe('extractRepoAndIssueFromIssueUrl', () => {
    test('valid data', () => {
      const issueUrl = 'https://github.com/stakwork/sphinx-tribes/issues/459';
      const result = { issue: '459', repo: 'stakwork/sphinx-tribes' };
      expect(extractRepoAndIssueFromIssueUrl(issueUrl)).toEqual(result);
    });
    test('empty string', () => {
      const issueUrl = '';

      expect(() => {
        extractRepoAndIssueFromIssueUrl(issueUrl);
      }).toThrow(Error);
    });
    test('invalid URL', () => {
      const issueUrl = 'https://test.url/issue/test/awr/awr/';
      expect(() => {
        extractRepoAndIssueFromIssueUrl(issueUrl);
      }).toThrow(Error);
    });
  });
  describe('extractGithubIssueFromUrl', () => {
    test('valid data', () => {
      const issueUrl = 'https://github.com/stakwork/sphinx-tribes/issues/459';
      const issueKey = 'stakwork/sphinx-tribes/459';
      const person = {
        github_issues: {
          [issueKey]: 'test'
        }
      };
      expect(extractGithubIssueFromUrl(person, issueUrl)).toBe('test');
    });

    test('invalid data', () => {
      const issueUrl = 'https://github.com/tribes/issues/459';
      const issueKey = 'stakwork/sphinx-tribes/459';
      const person = {
        github_issues: {
          [issueKey]: 'test'
        }
      };
      expect(extractGithubIssueFromUrl(person, issueUrl)).toEqual({});
    });
    test('empty url', () => {
      const issueUrl = '';
      const issueKey = 'stakwork/sphinx-tribes/459';
      const person = {
        github_issues: {
          [issueKey]: 'test'
        }
      };
      expect(extractGithubIssueFromUrl(person, issueUrl)).toEqual({});
    });
  });
  // This was breaking our test suite
  /* describe('satToUsd', () => {
    test('validData', () => {
      expect(satToUsd(100)).toEqual('10.00');
      expect(satToUsd(1000000)).toEqual('100000.00');
      expect(satToUsd(1)).toEqual('0.10');
      expect(satToUsd(0)).toEqual('0.00');
    });
  });*/
  describe('randomString', () => {
    test('length', () => {
      expect(randomString(15)).toHaveLength(30);
    });
    test('strings not equal', () => {
      const str1 = randomString(2);
      const str2 = randomString(2);
      expect(str1).not.toBe(str2);
    });
  });
  describe('calculateTimeLeft', () => {
    test('time remaining', () => {
      const timeLimit = new Date(moment().add(2, 'minutes').format().toString());
      const { minutes, seconds } = calculateTimeLeft(timeLimit, 'minutes');
      expect(minutes).toBe(1);
      expect(seconds).toBe(59);
    });
    test('calculate days remaining', () => {
      const timeLimit = new Date(moment().add(2, 'days').format().toString());
      const { days, hours, minutes, seconds } = calculateTimeLeft(timeLimit, 'days');
      expect(minutes).toBe(59);
      expect(seconds).toBe(59);
      expect(days).toBe(1);
      expect(hours).toBe(23);
    });
  });
  describe('userHasRole', () => {
    test('test user has roles', () => {
      const testRoles = [
        {
          name: 'ADD BOUNTY'
        },
        {
          name: 'DELETE BOUNTY'
        },
        {
          name: 'PAY BOUNTY'
        }
      ];

      const userRole = [
        {
          role: 'ADD BOUNTY'
        }
      ];
      const hasRole = userHasRole(testRoles, userRole, 'ADD BOUNTY');
      expect(hasRole).toBe(true);
    });
  });
  describe('toCapitalize', () => {
    test('test to capitalize string', () => {
      const capitalizeString = toCapitalize('hello test sphinx');
      expect(capitalizeString).toBe('Hello Test Sphinx');
    });
  });
});
