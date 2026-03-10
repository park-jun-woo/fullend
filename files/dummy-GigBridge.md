# Project: GigBridge (Freelance Escrow Matching Platform)
**Target AI Agent:** Generate 10 SSOT files for the `fullend` framework based on the following domain constraints and business rules. Ensure zero mismatches in `fullend validate`.

## 1. Domain Overview
GigBridge is a platform where `clients` post projects (gigs), `freelancers` submit proposals, and payments are held in escrow until the client approves the submitted work. 10% platform fee is deducted upon completion.

## 2. Entity & DDL Requirements (`db/schema.sql`, `db/queries.sql`)
Must implement the following core tables with correct Foreign Keys.
* **users**: `id`, `email`, `password_hash`, `role` ('client', 'freelancer', 'admin'), `name`
* **gigs**: `id`, `client_id` (FK), `title`, `description`, `budget` (int), `status` (string), `freelancer_id` (FK, nullable), `created_at`
* **proposals**: `id`, `gig_id` (FK), `freelancer_id` (FK), `bid_amount` (int), `status` (string)
* **transactions**: `id`, `gig_id` (FK), `type` ('hold', 'release', 'refund'), `amount` (int), `created_at`

## 3. State Machine Rules (`states/gig.md`, `states/proposal.md`)
Implement Mermaid `stateDiagram-v2`.
* **Gig States (`gig`)**:
    * `[*] --> draft` (Default)
    * `draft --> open`: **PublishGig**
    * `open --> in_progress`: **AcceptProposal**
    * `in_progress --> under_review`: **SubmitWork**
    * `under_review --> completed`: **ApproveWork**
    * `under_review --> disputed`: **RaiseDispute**
* **Proposal States (`proposal`)**:
    * `[*] --> pending` (Default)
    * `pending --> accepted`: **AcceptProposal**
    * `pending --> rejected`: **RejectProposal**

## 4. Authorization Rules (`policy/authz.rego`)
Use OPA Rego with `@ownership` annotations.
* **Ownerships**: 
    * `gig`: `gigs.client_id`
    * `gig_assignee`: `gigs.freelancer_id`
    * `proposal`: `proposals.freelancer_id`
* **Rules (allow if)**:
    * `PublishGig`: Role 'client' AND owns `gig`
    * `SubmitProposal`: Role 'freelancer' (cannot submit to own gig)
    * `AcceptProposal`: Role 'client' AND owns `gig`
    * `SubmitWork`: Role 'freelancer' AND is `gig_assignee`
    * `ApproveWork`: Role 'client' AND owns `gig`

## 5. API & Business Logic (OpenAPI ↔ SSaC)
Map `operationId` exactly to SSaC function names. Apply JWT `bearerAuth`.

* **POST /gigs (CreateGig)**: Client only. Creates gig in `draft` state.
* **PUT /gigs/{id}/publish (PublishGig)**: Auth check + State transition `draft` -> `open`.
* **GET /gigs (ListGigs)**: Public. 
    * OpenAPI Extensions required: `x-pagination` (limit/offset), `x-sort` (budget, created_at), `x-filter` (status, budget), `x-include` (`client_id:users.id`).
    * SSaC must use `query` arg to support these extensions.
* **POST /gigs/{id}/proposals (SubmitProposal)**: Freelancer only. Creates proposal.
* **POST /proposals/{id}/accept (AcceptProposal)**: 
    * Verifies Client owns the gig related to the proposal.
    * Transitions Proposal `pending` -> `accepted`.
    * Transitions Gig `open` -> `in_progress`.
    * Assigns `freelancer_id` to the Gig.
    * **@call `billing.holdEscrow`**: Holds the budget from the client.
* **POST /gigs/{id}/submit-work (SubmitWork)**: Freelancer only (Assignee). `in_progress` -> `under_review`.
* **POST /gigs/{id}/approve (ApproveWork)**: Client only. `under_review` -> `completed`.
    * **@call `billing.releaseFunds`**: Deducts 10% fee and sends funds to freelancer.
    * **@call `mail.sendTemplateEmail`**: Notify freelancer.

## 6. Custom Functions (`func/billing/*.go`)
Declare these using `// @func` so the generator can create stubs.
* `holdEscrow(gigID, amount, clientID)`: Simulates locking funds. Returns transaction ID.
* `releaseFunds(gigID, amount, freelancerID)`: Calculates 10% platform fee, releases 90% to freelancer. Returns transaction ID.

## 7. Frontend UI (`frontend/gigs.html`, `frontend/gig-detail.html`)
Use STML `data-*` attributes.
* `gigs.html`: Must implement `data-fetch="ListGigs"`, `data-paginate`, `data-sort="created_at:desc"`, `data-filter="status"`. Render gig list using `data-each`. Include client name via `x-include` binding (e.g., `data-bind="client.name"`).
* `gig-detail.html`: Show actions based on state (e.g., `<button data-action="SubmitWork" data-state="gig.status === 'in_progress'">`).

## 8. E2E Testing Scenarios (`scenario/gig_lifecycle.feature`)
Write fixed-pattern Gherkin for Hurl generation.
* **@scenario**: Happy Path: Client creates gig -> Freelancer A submits proposal -> Client accepts -> Escrow held -> Freelancer A submits work -> Client approves -> Funds released.
* **@invariant**: Unauthorized Access: Freelancer B tries to `SubmitWork` on Freelancer A's assigned gig -> Expect `status == 403`.
* **@invariant**: Invalid State: Client tries to `ApproveWork` when gig is in `open` state -> Expect `status == 409`.

---
**Agent Instruction:**
1. Generate `fullend.yaml` matching this module config (`github.com/gigbridge/api`).
2. Generate all 10 SSOTs strictly adhering to the `fullend` validation rules.
3. Ensure no mismatches between OpenAPI, DDL, Mermaid, Rego, and SSaC `@call`/`@auth`/`@state`/`@empty` guards.
4. Output the files in their respective directory paths.