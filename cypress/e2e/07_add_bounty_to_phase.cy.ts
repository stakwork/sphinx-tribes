import { User, HostName, Workspaces, Bounties } from '../support/objects/objects';


describe('Create Bounty', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 2; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/gobounties/`,
                    headers: { 'x-jwt': `${value}` },
                    body: Bounties[i],
                    failOnStatusCode: false
                })//.its('body').should('have.property', 'id', Workspaces[i].name.trim())
                .then( value => {
                    console.log(value);
                });
            }
        })
    })
})