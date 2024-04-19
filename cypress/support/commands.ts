import dotenv from 'dotenv';
dotenv.config();

import { HostName } from './objects/objects';

// @ts-check
// @ts-check
/// <reference types="cypress" />

Cypress.Commands.add('create_person', (person: Cypress.Person) => {
    cy.request({
        method: 'POST',
        url: `http://${HostName}/person/test`,
        headers: {},
        body: person
    }).then(res => {
        return res
    });
});