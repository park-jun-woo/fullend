package service

// @get Gig gig = Gig.FindByID(request.ID)
// @empty gig "Gig not found"
// @auth "dispute" "gig" {id: gig.ID} "Not authorized"
// @state gig {status: gig.Status} "RaiseDispute" "Cannot raise dispute"
// @put Gig.UpdateStatus(gig.ID, "disputed")
// @get Gig gig = Gig.FindByID(gig.ID)
// @response {
//   gig: gig
// }
func RaiseDispute() {}
