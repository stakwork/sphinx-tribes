import { User, HostName, Workspaces } from '../support/objects/objects';

describe('Create Workspaces', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i < Workspaces.length; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/workspaces`,
                    headers: { 'x-jwt': `${value}` },
                    body: Workspaces[i],
                    failOnStatusCode: false
                }).then((response) => {
                    expect(response.status).to.eq(200);
                    expect(response.body).to.have.property('name', Workspaces[i].name.trim());
                });
            }
        });
    });
});


describe('Edit Mission', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 2; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/workspaces/mission`,
                    headers: { 'x-jwt': `${value}` },
                    body: {
                        uuid: Workspaces[i].uuid,
                        owner_pubkey: Workspaces[i].owner_pubkey,
                        mission: Workspaces[i].mission + '_addedtext'
                    }
                }).then((resp) => {
                    expect(resp.status).to.eq(200)
                })
            }
        })
    })
})

describe('Edit Tactics', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 2; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/workspaces/tactics`,
                    headers: { 'x-jwt': `${value}` },
                    body: {
                        uuid: Workspaces[i].uuid,
                        owner_pubkey: Workspaces[i].owner_pubkey,
                        tactics: Workspaces[i].tactics + '_addedtext'
                    }
                }).then((resp) => {
                    expect(resp.status).to.eq(200)
                })
            }
        })
    })
})

describe('Edit Schematics Url', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 2; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/workspaces/schematicurl`,
                    headers: { 'x-jwt': `${value}` },
                    body: {
                        uuid: Workspaces[i].uuid,
                        owner_pubkey: Workspaces[i].owner_pubkey,
                        schematic_url: Workspaces[i].schematic_url + '_addedtext'
                    }
                }).then((resp) => {
                    expect(resp.status).to.eq(200)
                })
            }
        })
    })
})


describe('Check Workspace Values', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 2; i++) {
                cy.request({
                    method: 'GET',
                    url: `${HostName}/workspaces/` + Workspaces[i].uuid,
                    headers: { 'x-jwt': `${ value }` },
                    body: {} 
                }).then((resp) => {
                    expect(resp.status).to.eq(200)
                    expect(resp.body).to.have.property('mission', Workspaces[i].mission.trim() + '_addedtext')
                    expect(resp.body).to.have.property('tactics', Workspaces[i].tactics.trim() + '_addedtext')
                    expect(resp.body).to.have.property('schematic_url', Workspaces[i].schematic_url.trim() + '_addedtext')
                })
            }
        })
    })
})
