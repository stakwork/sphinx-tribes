import { User, HostName, Workspaces } from '../support/objects/objects';


describe('Create Workspaces', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            //Create 3 Workspaces
            cy.request({
                method: 'POST',
                url: `${HostName}/workspaces/`,
                headers: { 'x-jwt': `${value}` },
                body: Workspaces[0]
            }).its('body').should('have.property', 'name', Workspaces[0].name.trim());

            cy.request({
                method: 'POST',
                url: `${HostName}/workspaces/`,
                headers: { 'x-jwt': `${value}` },
                body: Workspaces[1]
            });

            cy.request({
                method: 'POST',
                url: `${HostName}/workspaces/`,
                headers: { 'x-jwt': `${value}` },
                body: Workspaces[2]
            });
        })
    })
})

describe('Edit Mission', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            //Create 3 Workspaces
            cy.request({
                method: 'POST',
                url: `${HostName}/workspaces/mission`,
                headers: { 'x-jwt': `${value}` },
                body: {
                    uuid: Workspaces[0].uuid,
                    mission: 'This is a sample mission for workspace'
                }
            }).then((resp) => {
                expect(resp.status).to.eq(200)
            })

        })
    })
})

describe('Edit Tactics', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            //Create 3 Workspaces
            cy.request({
                method: 'POST',
                url: `${HostName}/workspaces/tactics`,
                headers: { 'x-jwt': `${value}` },
                body: {
                    uuid: Workspaces[0].uuid,
                    mission: 'This is a sample tactics and objectives for workspace'
                }
            }).then((resp) => {
                expect(resp.status).to.eq(200)
            })

        })
    })
})

describe('Edit Schematics Url', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            //Create 3 Workspaces
            cy.request({
                method: 'POST',
                url: `${HostName}/workspaces/schematicurl`,
                headers: { 'x-jwt': `${value}` },
                body: {
                    uuid: Workspaces[0].uuid,
                    mission: 'This is a sample schematic url for workspaces'
                }
            }).then((resp) => {
                expect(resp.status).to.eq(200)
            })

        })
    })
})
