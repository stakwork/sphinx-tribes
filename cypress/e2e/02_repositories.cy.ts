import { User, HostName, Repositories } from '../support/objects/objects';


describe('Create Repositories for Workspace', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for (let i = 0; i <= 1; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/workspaces/repositories`,
                    headers: { 'x-jwt': `${value}` },
                    body: Repositories[i]
                }).its('body').then(body => {
                    expect(body).to.have.property('name').and.equal(Repositories[i].name.trim());
                    expect(body).to.have.property('url').and.equal(Repositories[i].url.trim());
                });
            }
        })
    })
})

describe('Modify Repository name for Workspace', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for (let i = 0; i <= 1; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/workspaces/repositories`,
                    headers: { 'x-jwt': `${value}` },
                    body: {
                        uuid: Repositories[i].uuid,
                        name: Repositories[i].name.trim() + "_addText"
                    }
                }).its('body').then(body => {
                    expect(body).to.have.property('name').and.equal(Repositories[i].name.trim() + "_addText");
                    expect(body).to.have.property('url').and.equal(Repositories[i].url.trim());
                });
            }
        })
    })
})

describe('Modify Repository url for Workspace', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for (let i = 0; i <= 1; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/workspaces/repositories`,
                    headers: { 'x-jwt': `${value}` },
                    body: {
                        uuid: Repositories[i].uuid,
                        url: Repositories[i].url.trim() + "_addText"
                    }
                }).its('body').then(body => {
                    expect(body).to.have.property('name').and.equal(Repositories[i].name.trim() + "_addText");
                    expect(body).to.have.property('url').and.equal(Repositories[i].url.trim() + "_addText");
                });
            }
        })
    })
})


describe('Check Repositories Values', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for (let i = 0; i <= 1; i++) {
                cy.request({
                    method: 'GET',
                    url: `${HostName}/workspaces/repositories/` + Repositories[i].workspace_uuid,
                    headers: { 'x-jwt': `${value}` },
                    body: {}
                }).then((resp) => {
                    expect(resp.status).to.eq(200)
                    expect(resp.body[i]).to.have.property('name', Repositories[i].name.trim() + "_addText")
                    expect(resp.body[i]).to.have.property('url', Repositories[i].url.trim() + "_addText")
                })
            }
        })
    })
})

describe('Get repository by uuid', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for (let i = 0; i <= 1; i++) {
                cy.request({
                    method: 'GET',
                    url: `${HostName}/workspaces/${Repositories[i].workspace_uuid}/repository/${Repositories[i].uuid}`,
                    headers: { 'x-jwt': `${value}` }
                }).then((resp) => {
                    expect(resp.status).to.eq(200)
                    expect(resp.body).to.have.property('name', Repositories[i].name.trim() + "_addText")
                    expect(resp.body).to.have.property('url', Repositories[i].url.trim() + "_addText")
                })
            }
        })
    })
})

describe('Delete repository by uuid', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            cy.request({
                method: 'DELETE',
                url: `${HostName}/workspaces/${Repositories[0].workspace_uuid}/repository/${Repositories[0].uuid}`,
                headers: { 'x-jwt': `${value}` },
                body: {}
            }).then((resp) => {
                expect(resp.status).to.eq(200)
            })
        })
    })
})

describe('Check delete by uuid', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            cy.request({
                method: 'GET',
                url: `${HostName}/workspaces/${Repositories[0].workspace_uuid}/repository/${Repositories[0].uuid}`,
                headers: { 'x-jwt': `${value}` },
                body: {},
                failOnStatusCode: false
            }).then((resp) => {
                expect(resp.status).to.eq(404);
            })
        })
    })
})