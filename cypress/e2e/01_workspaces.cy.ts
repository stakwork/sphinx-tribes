import {User, HostName, Workspaces, Features} from '../support/objects/objects';


describe('Create Workspaces', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 2; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/workspaces`,
                    headers: { 'x-jwt': `${value}` },
                    body: Workspaces[i]
                }).its('body').should('have.property', 'name', Workspaces[i].name.trim());
            }
        })
    })
})

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


describe('Get Features for Workspace', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            cy.request({
                method: 'GET',
                url: `${HostName}/workspaces/` + Features[0].workspace_uuid + `/features`,
                headers: { 'x-jwt': `${value}` },
                body: {}
            }).then((resp) => {
                expect(resp.status).to.eq(200);
                if (resp.status === 200) {
                    resp.body.forEach((feature) => {
                        const expectedFeature = Features.find(f => f.uuid === feature.uuid);
                        expect(feature).to.have.property('name', expectedFeature.name.trim() + " _addtext");
                        expect(feature).to.have.property('brief', expectedFeature.brief.trim() + " _addtext");
                        expect(feature).to.have.property('requirements', expectedFeature.requirements.trim() + " _addtext");
                        expect(feature).to.have.property('architecture', expectedFeature.architecture.trim() + " _addtext");
                    });
                }
            });
        })
    })
})