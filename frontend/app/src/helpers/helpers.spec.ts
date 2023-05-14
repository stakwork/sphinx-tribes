import {
  extractGithubIssueFromUrl,
  extractRepoAndIssueFromIssueUrl,
  randomString,
  satToUsd,
  calculateTimeLeft
} from './helpers';
import { uiStore } from '../store/ui';
import crypto from 'crypto';
import moment from 'moment';

beforeAll(() => {
  uiStore.setUsdToSatsExchangeRate(10);
  // for test randomString
  Object.defineProperty(window, 'crypto', {
    value: {
      getRandomValues: (arr) => crypto.randomBytes(arr.length)
    }
  });
});

afterAll(() => {
  uiStore.setUsdToSatsExchangeRate(0);
});

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
  describe('satToUsd', () => {
    test('validData', () => {
      expect(satToUsd(100)).toEqual('10.00');
      expect(satToUsd(1000000)).toEqual('100000.00');
      expect(satToUsd(1)).toEqual('0.10');
      expect(satToUsd(0)).toEqual('0.00');
    });
  });
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
      const { minutes, seconds } = calculateTimeLeft(timeLimit);
      expect(minutes).toBe(1);
      expect(seconds).toBe(59);
    });
  });
});
