


import '../support/objects/objects.js'

describe('Create Workspace', () => {
  it('passes', () => {
    cy.request({
      method: 'POST',
      url: `${hostname}/workspaces/`,
      headers: {
        'x-user-token': `${user.authToken}`
      },
      body: workspaces[0]
    }).then((response) => {
      id = response.body.response.id;
    })
  })
});