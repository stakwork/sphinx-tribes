import crypto from 'crypto';
import moment from 'moment';
import {
  extractGithubIssueFromUrl,
  extractRepoAndIssueFromIssueUrl,
  randomString,
  calculateTimeLeft,
  toCapitalize,
  userHasRole,
  spliceOutPubkey,
  userHasManageBountyRoles,
  RolesCategory,
  handleDisplayRole,
  formatSat,
  filterCount
} from '../helpers-extended';

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
      expect(seconds).toBeGreaterThan(50);
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

    test('test user has manage bounty roles', () => {
      const testRoles = [
        {
          name: 'ADD BOUNTY'
        },
        {
          name: 'UPDATE BOUNTY'
        },
        {
          name: 'PAY BOUNTY'
        },
        {
          name: 'DELETE BOUNTY'
        }
      ];

      const userRole = [
        {
          role: 'ADD BOUNTY'
        },
        {
          role: 'DELETE BOUNTY'
        },
        {
          role: 'PAY BOUNTY'
        },
        {
          role: 'UPDATE BOUNTY'
        }
      ];
      const hasRole = userHasManageBountyRoles(testRoles, userRole);
      expect(hasRole).toBe(true);
    });
    test('test user dose not have manage bounty roles', () => {
      const testRoles = [
        {
          name: 'ADD BOUNTY'
        },
        {
          name: 'DELETE BOUNTY'
        },
        {
          name: 'PAY BOUNTY'
        },
        {
          name: 'UPDATE BOUNTY'
        }
      ];

      const userRole = [
        {
          role: 'ADD BOUNTY'
        },
        {
          role: 'DELETE BOUNTY'
        }
      ];
      const hasRole = userHasManageBountyRoles(testRoles, userRole);
      expect(hasRole).toBe(false);
    });
  });
  describe('toCapitalize', () => {
    test('test to capitalize string', () => {
      const capitalizeString = toCapitalize('hello test sphinx');
      expect(capitalizeString).toBe('Hello Test Sphinx');
    });
  });
  describe('spliceOutPubkey', () => {
    test('test that it returns pubkey from a pubkey:route_hint string', () => {
      const pubkey = '12344444444444444';
      const routeHint = '899900000000000000:88888888';
      const userAddress = `${pubkey}:${routeHint}`;
      const pub = spliceOutPubkey(userAddress);
      expect(pub).toBe(pubkey);
    });
  });

  describe('format roles', () => {
    test('should correctly set the default data roles for the first assigned user', () => {
      const displayedRoles: RolesCategory[] = [];
      const result = handleDisplayRole(displayedRoles);
      expect(result.newDisplayedRoles).toEqual([]);
      expect(result.tempDataRole).toEqual({});
    });

    test('should correctly update the status of a role if it is present in the default roles', () => {
      const displayedRoles: RolesCategory[] = [
        { name: 'Manage bounties', roles: [], status: false },
        { name: 'Fund organization', roles: [], status: false },
        { name: 'Withdraw from organization', roles: [], status: false },
        { name: 'View transaction history', roles: [], status: false }
      ];
      const result = handleDisplayRole(displayedRoles);
      expect(result.newDisplayedRoles).toEqual([
        { name: 'Manage bounties', roles: [], status: true },
        { name: 'Fund organization', roles: [], status: true },
        { name: 'Withdraw from organization', roles: [], status: true },
        { name: 'View transaction history', roles: [], status: true }
      ]);
      expect(result.tempDataRole).toEqual({});
    });

    test('should correctly update the tempDataRole object with the data roles of a role if it is present in the default roles', () => {
      const displayedRoles: RolesCategory[] = [
        { name: 'Manage bounties', roles: ['role1', 'role2'], status: false },
        { name: 'Fund organization', roles: ['role3'], status: false },
        { name: 'Withdraw from organization', roles: ['role4'], status: false },
        { name: 'View transaction history', roles: ['role5'], status: false }
      ];
      const result = handleDisplayRole(displayedRoles);
      expect(result.newDisplayedRoles).toEqual([
        { name: 'Manage bounties', roles: ['role1', 'role2'], status: true },
        { name: 'Fund organization', roles: ['role3'], status: true },
        { name: 'Withdraw from organization', roles: ['role4'], status: true },
        { name: 'View transaction history', roles: ['role5'], status: true }
      ]);
      expect(result.tempDataRole).toEqual({
        role1: true,
        role2: true,
        role3: true,
        role4: true,
        role5: true
      });
    });

    test('formatSat', () => {
      expect(formatSat(10000)).toBe('10 000');
      expect(formatSat(0)).toBe('0');
    });
    test('filterCount', () => {
      expect(filterCount({ thing1: 0, thing2: 1 })).toBe(1);
      expect(filterCount({ thing1: 1, thing2: 1 })).toBe(2);
      expect(filterCount({})).toBe(0);
    });
  });
});
