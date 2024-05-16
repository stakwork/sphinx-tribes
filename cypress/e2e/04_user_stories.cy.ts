import { User, HostName, UserStories } from '../support/objects/objects';

describe('Create user stories for Feature', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 5; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/features/story`,
                    headers: { 'x-jwt': `${value}` },
                    body: UserStories[i]
                }).its('body').then(body => {
                    expect(body).to.have.property('uuid').and.equal(UserStories[i].uuid.trim());
                    expect(body).to.have.property('feature_uuid').and.equal(UserStories[i].feature_uuid.trim());
                    expect(body).to.have.property('description').and.equal(UserStories[i].description.trim());
                    expect(body).to.have.property('priority').and.equal(UserStories[i].priority);
                });
            }
        })
    })
})

describe('Modify user story description', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 5; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/features/story`,
                    headers: { 'x-jwt': `${value}` },
                    body: {
                        uuid: UserStories[i].uuid,
                        description: UserStories[i].description + "_addtext"
                    }
                }).its('body').then(body => {
                    expect(body).to.have.property('uuid').and.equal(UserStories[i].uuid.trim());
                    expect(body).to.have.property('feature_uuid').and.equal(UserStories[i].feature_uuid.trim());
                    expect(body).to.have.property('description').and.equal(UserStories[i].description.trim() + " _addtext");
                    expect(body).to.have.property('priority').and.equal(UserStories[i].priority);
                });
            }
        })
    })
})

describe('Get user stories for feature', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            cy.request({
                method: 'GET',
                url: `${HostName}/features/${UserStories[0].feature_uuid}/story`,
                headers: { 'x-jwt': `${ value }` },
                body: {} 
            }).then((resp) => {
                expect(resp.status).to.eq(200)
                for(let i = 0; i <= 5; i++) {
                    expect(resp.body[i]).to.have.property('uuid').and.equal(UserStories[i].uuid.trim());
                    expect(resp.body[i]).to.have.property('feature_uuid').and.equal(UserStories[i].feature_uuid.trim());
                    expect(resp.body[i]).to.have.property('description').and.equal(UserStories[i].description.trim() + " _addtext");
                    expect(resp.body[i]).to.have.property('priority').and.equal(UserStories[i].priority);
                }
            })
        })
    })
})

describe('Get story by uuid', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for(let i = 0; i <= 5; i++) {
                cy.request({
                    method: 'GET',
                    url: `${HostName}/features/${UserStories[0].feature_uuid}/story/${UserStories[i].uuid}`,
                    headers: { 'x-jwt': `${value}` },
                    body: {}
                }).then((resp) => {
                    expect(resp.status).to.eq(200);
                    expect(resp.body).to.have.property('uuid').and.equal(UserStories[i].uuid.trim());
                    expect(resp.body).to.have.property('feature_uuid').and.equal(UserStories[i].feature_uuid.trim());
                    expect(resp.body).to.have.property('description').and.equal(UserStories[i].description.trim() + " _addtext");
                    expect(resp.body).to.have.property('priority').and.equal(UserStories[i].priority);
                });
            }
        });
    });
});

describe('Delete story by uuid', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            cy.request({
                method: 'DELETE',
                url: `${HostName}/features/${UserStories[0].feature_uuid}/story/${UserStories[0].uuid}`,
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
                url: `${HostName}/features/${UserStories[0].feature_uuid}/story/${UserStories[0].uuid}`,
                headers: { 'x-jwt': `${ value }` },
                body: {},
                failOnStatusCode: false
            }).then((resp) => {
                expect(resp.status).to.eq(404);
            })
        })
    })
})
