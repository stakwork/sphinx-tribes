import { User, HostName, Workspaces, Repositories, Features } from '../support/objects/objects';



describe('Create Features for Workspace', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 2; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/features`,
                    headers: { 'x-jwt': `${value}` },
                    body: Features[i]
                }).its('body').then(body => {
                    expect(body).to.have.property('name').and.equal(Features[i].name.trim());
                    expect(body).to.have.property('brief').and.equal(Features[i].brief.trim());
                    expect(body).to.have.property('requirements').and.equal(Features[i].requirements.trim());
                    expect(body).to.have.property('architecture').and.equal(Features[i].architecture.trim());
                });
            }
        })
    })
})

describe('Modify name for Feature', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 2; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/features`,
                    headers: { 'x-jwt': `${value}` },
                    body: {
                        uuid: Features[i].uuid,
                        name: Features[i].name + "_addtext"
                    }
                }).its('body').then(body => {
                    expect(body).to.have.property('name').and.equal(Features[i].name.trim() + " _addtext");
                    expect(body).to.have.property('brief').and.equal(Features[i].brief.trim());
                    expect(body).to.have.property('requirements').and.equal(Features[i].requirements.trim());
                    expect(body).to.have.property('architecture').and.equal(Features[i].architecture.trim());
                });
            }
        })
    })
})

describe('Modify brief for Feature', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 2; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/features`,
                    headers: { 'x-jwt': `${value}` },
                    body: {
                        uuid: Features[i].uuid,
                        brief: Features[i].brief + "_addtext"
                    }
                }).its('body').then(body => {
                    expect(body).to.have.property('name').and.equal(Features[i].name.trim() + " _addtext");
                    expect(body).to.have.property('brief').and.equal(Features[i].brief.trim() + " _addtext");
                    expect(body).to.have.property('requirements').and.equal(Features[i].requirements.trim());
                    expect(body).to.have.property('architecture').and.equal(Features[i].architecture.trim());
                });
            }
        })
    })
})

describe('Modify requirements for Feature', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 2; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/features`,
                    headers: { 'x-jwt': `${value}` },
                    body: {
                        uuid: Features[i].uuid,
                        requirements: Features[i].requirements + "_addtext"
                    }
                }).its('body').then(body => {
                    expect(body).to.have.property('name').and.equal(Features[i].name.trim() + " _addtext");
                    expect(body).to.have.property('brief').and.equal(Features[i].brief.trim() + " _addtext");
                    expect(body).to.have.property('requirements').and.equal(Features[i].requirements.trim() + " _addtext");
                    expect(body).to.have.property('architecture').and.equal(Features[i].architecture.trim());
                });
            }
        })
    })
})

describe('Modify architecture for Feature', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 2; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/features`,
                    headers: { 'x-jwt': `${value}` },
                    body: {
                        uuid: Features[i].uuid,
                        architecture: Features[i].architecture + "_addtext"
                    }
                }).its('body').then(body => {
                    expect(body).to.have.property('name').and.equal(Features[i].name.trim() + " _addtext");
                    expect(body).to.have.property('brief').and.equal(Features[i].brief.trim() + " _addtext");
                    expect(body).to.have.property('requirements').and.equal(Features[i].requirements.trim() + " _addtext");
                    expect(body).to.have.property('architecture').and.equal(Features[i].architecture.trim() + " _addtext");
                });
            }
        })
    })
})


describe('Get Features for Workspace', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            cy.request({
                method: 'GET',
                url: `${HostName}/workspaces/${Features[0].workspace_uuid}/features`, //changed from url: `${HostName}/features/forworkspace/` + Features[0].workspace_uuid, please update the routes file and any other change needed.
                headers: { 'x-jwt': `${ value }` },
                body: {}
            }).then((resp) => {
                expect(resp.status).to.eq(200)
                resp.body.forEach((feature) => {
                    const expectedFeature = Features.find(f => f.uuid === feature.uuid);
                    expect(feature).to.have.property('name', expectedFeature.name.trim() + " _addtext");
                    expect(feature).to.have.property('brief', expectedFeature.brief.trim() + " _addtext");
                    expect(feature).to.have.property('requirements', expectedFeature.requirements.trim() + " _addtext");
                    expect(feature).to.have.property('architecture', expectedFeature.architecture.trim() + " _addtext");
                });
            })
        })
    })
})

describe('Get Feature by uuid', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 2; i++) {
                cy.request({
                    method: 'GET',
                    url: `${HostName}/features/` + Features[i].uuid,
                    headers: { 'x-jwt': `${ value }` },
                    body: {}
                }).then((resp) => {
                    expect(resp.status).to.eq(200)
                    expect(resp.body).to.have.property('name', Features[i].name.trim() + " _addtext")
                    expect(resp.body).to.have.property('brief', Features[i].brief.trim() + " _addtext")
                    expect(resp.body).to.have.property('requirements', Features[i].requirements.trim() + " _addtext")
                    expect(resp.body).to.have.property('architecture', Features[i].architecture.trim() + " _addtext")
                })
            }
        })
    })
})

describe('Delete Feature by uuid', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            cy.request({
                method: 'DELETE',
                url: `${HostName}/features/${Features[2].uuid}`,
                headers: { 'x-jwt': `${ value }` },
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
                method: 'DELETE',
                url: `${HostName}/features/${Features[2].uuid}`,
                headers: { 'x-jwt': `${ value }` },
                body: {},
                failOnStatusCode: false
            }).then((resp) => {
                expect(resp.status).to.eq(404);
            })
        })
    })
})
