package service

import "github.com/gigbridge/api/billing"

// @get Proposal proposal = Proposal.FindByID(request.ID)
// @empty proposal "Proposal not found"
// @get Gig gig = Gig.FindByID(proposal.GigID)
// @empty gig "Gig not found"
// @auth "accept" "gig" {id: gig.ID} "Not authorized"
// @state proposal {status: proposal.Status} "AcceptProposal" "Cannot accept proposal"
// @state gig {status: gig.Status} "AcceptProposal" "Cannot accept proposal"
// @put Proposal.UpdateStatus(proposal.ID, "accepted")
// @put Gig.AssignFreelancer(gig.ID, proposal.FreelancerID)
// @put Gig.UpdateStatus(gig.ID, "in_progress")
// @call int64 transactionID = billing.HoldEscrow(gig.ID, gig.Budget, gig.ClientID)
// @get Gig gig = Gig.FindByID(gig.ID)
// @response {
//   gig: gig,
//   transactionID: transactionID
// }
func AcceptProposal() {}
