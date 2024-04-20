import { User } from '../support/objects/objects'

describe('Create User', () => {
    it('it creates a user', () => {
        const response = cy.upsertlogin(User);
        return response;
    })
});