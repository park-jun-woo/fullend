package service

import "github.com/gigbridge/api/internal/billing"

// @get Gig gig = Gig.FindByID(request.ID)
// @empty gig "Gig not found"
// @auth "approve" "gig" {id: gig.ID} "Not authorized"
// @state gig {status: gig.Status} "ApproveWork" "Cannot approve work"
// @put Gig.UpdateStatus(gig.ID, "completed")
// @call int64 transactionID = billing.ReleaseFunds({GigID: gig.ID, Amount: gig.Budget, FreelancerID: gig.FreelancerID})
// @get Gig gig = Gig.FindByID(gig.ID)
// @response {
//   gig: gig,
//   transactionID: transactionID
// }
func ApproveWork() {}
