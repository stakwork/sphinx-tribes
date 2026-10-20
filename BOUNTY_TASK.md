# Bounty Objective: Notification Lifecycle Functionality and Key Exchange Handling  
Target: stakwork/sphinx-tribes Issue #2391
URL: https://github.com/stakwork/sphinx-tribes/issues/2391

## Issue Description
**Description:**  
This task involves implementing a function to process notifications and setting up a cron job to handle cases where the `contact_key` is unavailable. The function will manage the end-to-end notification lifecycle, including user verification, notification sending, and status logging. The cron job will ensure notifications with the `WAITING_KEY_EXCHANGE` status are periodically rechecked and processed when the `contact_key` becomes available.  

**Acceptance Criteria:**  

### 1. **Notification Processing Function**  
Create a function named `sendNotification` to handle the following:  
- **Input Parameters:**  
  - `pubkey`: User's public key.  
  - `event`: The event type (e.g., `bounty_assigned`).  
  - `content`: Notification content.  
  - `alias`: Alias of the user

- **Steps to Implement:**  
  1. **Verify User on v2 Bot:**  
     - Call the `/contact/{pubkey}` endpoint to check if the user exists.  
     - If `contact_key` is empty:  
       - Add the user to the bot's contact list using the `/add_contact` endpoint.  
       - Re-call `/contact/{pubkey}` to verify the `contact_key`.  

     - If the `contact_key` is still empty, log the notification in the database with the status `WAITING_KEY_EXCHANGE` and exit the function

  2. **Send the Notification:**  
     - If the `contact_key` is available, call the `/send` endpoint with the appropriate payload:  
       ```typescript
       {
         dest: string;       // User's pubkey
         amt_msat: 0;        // For text messages
         content: string;    // Notification content
         is_tribe: false;    // For direct messages
         wait: true;         // Synchronous behavior
       }
       ```  

  3. **Log the Result:**  
     - Update the database with the response from the `/send` endpoint.  
     - Store the final status (`COMPLETE`, `FAILED`, or `PENDING`) along with relevant metadata.  

### 2. **Cron Job for Key Exchange Handling**  
Implement a cron job to periodically process notifications with the status `WAITING_KEY_EXCHANGE`:  

- **Behavior:**  
  1. Query the database for all notifications with the status `WAITING_KEY_EXCHANGE`.  
  2. For each record:  
     - Re-call the `/contact/{pubkey}` endpoint to check if the `contact_key` is now available.  
     - If the `contact_key` is available:  
       - Call the function in charge of sending notifications to send the notification.  
       - Update the database with the response from the `/send` endpoint.  
     - If the `contact_key` is still unavailable increment retry by 1 and leave the record for the next cron job iteration.  

**Database Schema:**  
Ensure the database table `notifications` supports the following:  
- `id` (Primary Key)  
- `event` (Event type, e.g., `bounty_assigned`)  
- `pubkey`  
- `content`  
- `retries`  
- `status` (`PENDING`, `COMPLETE`, `FAILED`, `WAITING_KEY_EXCHANGE`)  
- `uuid` (Unique identifier for the notification)  
- `created_at` (Timestamp for record creation)  

**Additional Notes:**  
- The `sendNotification` function and cron job should be modular and reusable for future enhancements.  
- Ensure robust error handling for API calls and database interactions.  

This Ticket depends on [this](https://github.com/stakwork/sphinx-tribes/issues/2390)

## Instructions for Kuro (Autonomous Solver):
1. Use skills: `systematic-debugging` and `code-review-and-quality`.
2. Inspect the repository structure, locate relevant files, and understand the bug/feature.
3. Apply the minimal, robust fix or implementation.
4. Run any existing test suites or linters in the workspace.
5. Ensure the fix adheres to the project's existing code style.
6. When done, output a summary of files modified and verification results.
