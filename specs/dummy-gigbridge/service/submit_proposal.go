package service

// @get Gig gig = Gig.FindByID(request.GigID)
// @empty gig "Gig not found"
// @auth "submit_proposal" "gig" {id: gig.ID} "Not authorized"
// @post Proposal proposal = Proposal.Create(request.GigID, currentUser.ID, request.BidAmount)
// @response {
//   proposal: proposal
// }
func SubmitProposal() {}
