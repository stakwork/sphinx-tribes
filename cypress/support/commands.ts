import dotenv from 'dotenv';
dotenv.config();

import { HostName } from './objects/objects';

// @ts-check
// @ts-check
/// <reference types="cypress" />

Cypress.Commands.add('upsertlogin', (person: Cypress.Person) => {
    cy.request({
        method: 'POST',
        url: `http://${HostName}/person/upsertlogin`,
        headers: {},
        body: person
    }).then((response) => response.body);
});