package service

// @get Gig gig = Gig.FindByID(request.ID)
// @empty gig "Gig not found"
// @auth "publish" "gig" {id: gig.ID} "Not authorized"
// @state gig {status: gig.Status} "PublishGig" "Cannot publish gig"
// @put Gig.UpdateStatus(gig.ID, "open")
// @get Gig gig = Gig.FindByID(gig.ID)
// @response {
//   gig: gig
// }
func PublishGig() {}
