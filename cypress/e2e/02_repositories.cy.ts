import { User, HostName, Workspaces, Repositories } from '../support/objects/objects';


describe('Create Repositories for Workspace', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 1; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/workspaces/repositories`,
                    headers: { 'x-jwt': `${value}` },
                    body: Repositories[i]
                }).its('body').should('have.property', 'name', Repositories[i].name.trim())
                .its('body').should('have.property', 'url', Repositories[i].url.trim());
            }
        })
    })
})


describe('Check Repositories Values', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            cy.request({
                method: 'GET',
                url: `${HostName}/workspaces/repositories/` + Repositories[0].workspace_uuid,
                headers: { 'x-jwt': `${ value }` },
                body: {} 
            }).then((resp) => {
                expect(resp.status).to.eq(200)
                expect(resp.body[0]).to.have.property('name', Repositories[0].name.trim())
                expect(resp.body[0]).to.have.property('url', Repositories[0].url.trim())
                expect(resp.body[1]).to.have.property('name', Repositories[1].name.trim())
                expect(resp.body[1]).to.have.property('url', Repositories[1].url.trim())
            })
        })
    })
})
