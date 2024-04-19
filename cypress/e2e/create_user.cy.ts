describe('Create User', () => {
    it('it creates a user', () => {
        const person: Cypress.Person = {
            owner_pubkey: "test_pubkey",
            owner_alias: "test_research",
            unique_name: "test_unique_name",
            description: "this is a test",
            tags: [],
            img: "",
            unlisted: false,
            deleted: false,
            owner_route_hint: "",
            owner_contact_key: "",
            price_to_meet: 0,
            twitter_confirmed: false,
            referred_by: 0,
            extras: {},
            github_issues: {}
        }

        const user = cy.create_person(person)
    })
});