import sinon from 'sinon';
import request from 'request';
import { expect } from 'chai';
import { people } from '../../../__test__/__mockData__/persons';
import { Person } from '../../../store/main';
import mockBounties, {
  mockPaginatedBountiesPage
} from '../../../bounties/__mock__/mockBounties.data';
import { mockOrganizations } from '../__mock__/mockOrganizations.data';
import { Organization } from '../../../store/main';

describe('HomePage - Signed Out', () => {
  let getStub: sinon.SinonStub;

  beforeEach(() => {
    getStub = sinon.stub(request, 'get');
  });

  afterEach(() => {
    getStub.restore();
  });

  it('should be able to view the listed bounties when signed out', (done: jest.DoneCallback) => {
    getStub.yields(null, { statusCode: 200 }, JSON.stringify(mockBounties));

    request.get('/api/bounties', (err: any, res: any, body: string) => {
      if (err) {
        done(err);
        return;
      }
      expect(res.statusCode).to.equal(200);
      const responseBody = JSON.parse(body);
      expect(responseBody).to.be.an('array');
      expect(responseBody[0]).to.have.keys(['assignee', 'bounty', 'organization', 'owner']);
      done();
    });
  });

  it('should be able to view all bounties created by a specific organization', (done: jest.DoneCallback) => {
    getStub.yields(null, { statusCode: 200 }, JSON.stringify(mockOrganizations));

    request.get('/api/bounties/organization/{orgId}', (err: any, res: any, body: string) => {
      if (err) {
        done(err);
        return;
      }
      expect(res.statusCode).to.equal(200);
      const responseBody = JSON.parse(body);
      expect(responseBody).to.be.an('array');
      responseBody.forEach((bounty: Organization[]) => {
        expect(bounty).to.include.keys(
          'bounty_count',
          'budget',
          'created',
          'deleted',
          'id',
          'img',
          'name',
          'owner_pubkey',
          'show',
          'updated',
          'uuid'
        );
      });

      done();
    });
  });

  it('should be able to view the listed people when signed out', (done: jest.DoneCallback) => {
    getStub.yields(null, { statusCode: 200 }, JSON.stringify(people));

    request.get('/api/people', (err: any, res: any, body: string) => {
      if (err) {
        done(err);
        return;
      }

      expect(res.statusCode).to.equal(200);
      const responseBody = JSON.parse(body);
      expect(responseBody).to.be.an('array');
      responseBody.forEach((person: Person[]) => {
        expect(person).to.include.keys('id', 'pubkey', 'contact_key', 'alias');
      });

      done();
    });
  });

  it('should be able to load more bounties with pagination', (done: jest.DoneCallback) => {
    getStub.yields(null, { statusCode: 200 }, JSON.stringify(mockPaginatedBountiesPage));
    request.get('/api/bounties?page=2', (err: any, res: any, body: string) => {
      if (err) {
        done(err);
        return;
      }
      expect(res.statusCode).to.equal(200);
      const responseBody = JSON.parse(body);
      expect(responseBody.bounties).to.be.an('array');
      expect(responseBody.currentPage).to.equal(2);
      done();
    });
  });
});
