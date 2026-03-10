package service

// @get Gig gig = Gig.FindByID(request.ID)
// @empty gig "Gig not found"
// @auth "submit_work" "gig_assignee" {id: gig.ID} "Not authorized"
// @state gig {status: gig.Status} "SubmitWork" "Cannot submit work"
// @put Gig.UpdateStatus(gig.ID, "under_review")
// @get Gig gig = Gig.FindByID(gig.ID)
// @response {
//   gig: gig
// }
func SubmitWork() {}
