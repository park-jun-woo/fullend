package authz

import rego.v1

# @ownership gig: gigs.client_id
# @ownership proposal: proposals.freelancer_id

allow if {
    input.action == "PublishGig"
    input.resource == "gig"
    input.claims.role == "client"
    data.owners.gig[input.resource_id] == input.claims.user_id
}

allow if {
    input.action == "SubmitProposal"
    input.resource == "gig"
    input.claims.role == "freelancer"
}

allow if {
    input.action == "AcceptProposal"
    input.resource == "gig"
    input.claims.role == "client"
    data.owners.gig[input.resource_id] == input.claims.user_id
}

allow if {
    input.action == "SubmitWork"
    input.resource == "gig"
    input.claims.role == "freelancer"
}

allow if {
    input.action == "ApproveWork"
    input.resource == "gig"
    input.claims.role == "client"
    data.owners.gig[input.resource_id] == input.claims.user_id
}
