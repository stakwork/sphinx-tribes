# Bounty Objective: Create route /features/uuid/phase/uuid/tickets
Target: stakwork/sphinx-tribes Issue #1982
URL: https://github.com/stakwork/sphinx-tribes/issues/1982

## Issue Description
This endpoint will be used to return the full list of ticket UUIDS for a given feature and phase uuid. 

Add a route to features.go at /features/uuid/phase/uuid/tickets

		r.Get("/{uuid}/phase/{uuid}/tickets", featureHandlers.GetTickets)

## Context
We need to create a new endpoint that returns all ticket UUIDs associated with a specific feature and phase. This endpoint will be used to fetch the complete list of tickets for a given phase within a feature.

## Task
Add a new route to features.go that handles GET requests for retrieving all tickets associated with a feature and phase UUID combination.

## Outcome
A functioning route that accepts feature UUID and phase UUID parameters and returns the associated ticket UUIDs.

## Design
### Route Definition
In features.go:
```go
func RegisterFeatureRoutes(r chi.Router) {
    // Existing routes...
    
    // New route for getting phase tickets
    r.Get("/{featureUUID}/phase/{phaseUUID}/tickets", featureHandlers.GetTickets)
}
```

## Acceptance Criteria
- [ ] Route `/features/{featureUUID}/phase/{phaseUUID}/tickets` is registered
- [ ] Route is accessible via GET method
- [ ] Route parameters are properly defined:
  - `featureUUID`: Feature identifier
  - `phaseUUID`: Phase identifier
- [ ] Route is mapped to the correct handler function
- [ ] Route follows existing project patterns for URL structure

## Instructions for Kuro (Autonomous Solver):
1. Use skills: `systematic-debugging` and `code-review-and-quality`.
2. Inspect the repository structure, locate relevant files, and understand the bug/feature.
3. Apply the minimal, robust fix or implementation.
4. Run any existing test suites or linters in the workspace.
5. Ensure the fix adheres to the project's existing code style.
6. When done, output a summary of files modified and verification results.
