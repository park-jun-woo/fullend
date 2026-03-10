package service

// @get Proposal proposal = Proposal.FindByID(request.ID)
// @empty proposal "Proposal not found"
// @get Gig gig = Gig.FindByID(proposal.GigID)
// @empty gig "Gig not found"
// @auth "reject" "gig" {id: gig.ID} "Not authorized"
// @state proposal {status: proposal.Status} "RejectProposal" "Cannot reject proposal"
// @put Proposal.UpdateStatus(proposal.ID, "rejected")
// @get Proposal proposal = Proposal.FindByID(proposal.ID)
// @response {
//   proposal: proposal
// }
func RejectProposal() {}
