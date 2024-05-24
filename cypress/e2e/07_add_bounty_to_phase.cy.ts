import { User, HostName, Phases, Bounties } from '../support/objects/objects';


//This test passes! It only asserts that response contains workspace_uuid 
describe('Create Bounties - don\'t check phase_uuid or phase_priority', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for (let i = 0; i < Bounties.length; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/gobounties/`,
                    headers: { 'x-jwt': `${value}` },
                    body: Bounties[i],
                    failOnStatusCode: false
                }).then((resp) => {
                    expect(resp.status).to.eq(200);
                    expect(resp.body).to.have.property('workspace_uuid').and.equal(Bounties[i].workspace_uuid);
                    console.log(resp);
                })
            }
        })
    })
});

//This test initially does not pass! It asserts that the response should contain phase_uuid
//You need to add phase_uuid to bounties 
describe('Create Bounties - with check phase_uuid and phase_priority', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            for (let i = 0; i < Bounties.length; i++) {
                cy.request({
                    method: 'POST',
                    url: `${HostName}/gobounties/`,
                    headers: { 'x-jwt': `${value}` },
                    body: Bounties[i],
                    failOnStatusCode: false
                }).then((resp) => {
                    expect(resp.status).to.eq(200);
                    expect(resp.body).to.have.property('phase_uuid').and.equal(Bounties[i].phase_uuid);
                    expect(resp.body).to.have.property('phase_priority').and.equal(Bounties[i].phase_priority);
                    console.log(resp);
                })
            }
        })
    })
});


//This test passes! It only asserts that response contains workspace_uuid 
describe('Get All Bounties - don\'t check phase_uuid', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
                cy.request({
                    method: 'GET',
                    url: `${HostName}/gobounties/all?limit=10&sortBy=created&search=&page=1&resetPage=true&Open=true&Assigned=false&Completed=false&Paid=false&languages=`,
                    headers: { 'x-jwt': `${value}` },
                    failOnStatusCode: false
                }).then((resp) => {
                    expect(resp.status).to.eq(200);
                    JSON.parse(resp.body).forEach((bounty) => {
                        expect(bounty).to.have.property('bounty').to.have.property('workspace_uuid');
                    })
                })
        })
    })
});

//This test initially does not pass! It asserts that the response should contain phase_uuid
//You need to add phase_uuid to bounties 
describe('Get All Bounties - check phase_uuid', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
                cy.request({
                    method: 'GET',
                    url: `${HostName}/gobounties/all?limit=10&sortBy=created&search=&page=1&resetPage=true&Open=true&Assigned=false&Completed=false&Paid=false&languages=`,
                    headers: { 'x-jwt': `${value}` },
                    //body: Bounties[i],
                    failOnStatusCode: false
                }).then((resp) => {
                    expect(resp.status).to.eq(200);
                    JSON.parse(resp.body).forEach((bounty) => {
                        expect(bounty).to.have.property('bounty').to.have.property('phase_uuid');
                        expect(bounty).to.have.property('bounty').to.have.property('phase_priority');
                    })
                })
        })
    })
});

//This test initially does not pass! It asserts that the response should contain phase_uuid
//You need to create the endpoint in /handlers/features.go and the route in /routes/features.go 
describe('Get Bounties for Phase', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            cy.request({
                method: 'GET',
                url: `${HostName}/features/${Phases[1].feature_uuid}/phase/${Phases[1].uuid}/bounty`,
                headers: { 'x-jwt': `${value}` },
                failOnStatusCode: false
            }).then((resp) => {
                expect(resp.status).to.eq(200);
                console.log(resp.body);
                const responseBody = typeof resp.body === 'string' ? JSON.parse(resp.body) : resp.body;
                console.log(responseBody);
                if (Array.isArray(responseBody)) {
                    responseBody.forEach((bounty) => {
                        expect(bounty).to.have.property('bounty').to.have.property('phase_uuid');
                        expect(bounty).to.have.property('bounty').to.have.property('phase_priority');
                    });
                } else {
                    console.error('Expected an array but got:', responseBody);
                }
            });
        });
    });
});


//This test passes! It asserts that the endpoint should receive a valid phase_uuid
describe('Create Bounties with valid phase_uuid', () => {
    it('passes', () => {
        cy.upsertlogin(User).then(value => {
            cy.request({
                method: 'POST',
                url: `${HostName}/gobounties/`,
                headers: { 'x-jwt': `${value}` },
                body: {...Bounties[0], phase_uuid: Phases[0].uuid},
                failOnStatusCode: false
            }).then((resp) => {
                expect(resp.status).to.eq(200);
                expect(resp.body).to.have.property('phase_uuid').and.equal(Phases[0].uuid);
                expect(resp.body).to.have.property('phase_priority').and.equal(Phases[0].priority);
                console.log(resp);
            })
        })
    })
});