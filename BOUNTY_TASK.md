# Bounty Objective: Trigger Notifications During Ticket Assignment
Target: stakwork/sphinx-tribes Issue #2392
URL: https://github.com/stakwork/sphinx-tribes/issues/2392

## Issue Description
  **Description:**  
Modify the function responsible for assigning tickets to users so that it triggers a notification by calling the `sendNotification` function. This ensures users are notified immediately when a ticket is assigned to them.  

**Acceptance Criteria:**  

### 1. **Update the Ticket Assignment Function**  
- Integrate the `sendNotification` function call into the ticket assignment logic.  
- Parameters to be passed to `sendNotification`:  
  - **`pubkey`**: The public key of the user to whom the ticket is assigned.  
  - **`event`**: The event type, e.g., `ticket_assigned`.  
  - **`content`**: A message indicating the ticket assignment, e.g., `"You have been assigned a new ticket: {ticket_title}."`.  

This ticket depends on [this](https://github.com/stakwork/sphinx-tribes/issues/2391)

## Instructions for Kuro (Autonomous Solver):
1. Use skills: `systematic-debugging` and `code-review-and-quality`.
2. Inspect the repository structure, locate relevant files, and understand the bug/feature.
3. Apply the minimal, robust fix or implementation.
4. Run any existing test suites or linters in the workspace.
5. Ensure the fix adheres to the project's existing code style.
6. When done, output a summary of files modified and verification results.
