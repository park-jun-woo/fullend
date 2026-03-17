package authz

import rego.v1

# @ownership gig: gigs.client_id
# @ownership gig_assignee: gigs.freelancer_id
# @ownership proposal: proposals.freelancer_id

default allow := false

# PublishGig — client who owns the gig
allow if {
    input.action == "PublishGig"
    input.resource == "gig"
    input.claims.role == "client"
    data.owners.gig[input.resource_id] == input.claims.user_id
}

# SubmitProposal — any freelancer
allow if {
    input.action == "SubmitProposal"
    input.resource == "proposal"
    input.claims.role == "freelancer"
}

# AcceptProposal — client who owns the gig
allow if {
    input.action == "AcceptProposal"
    input.resource == "gig"
    input.claims.role == "client"
    data.owners.gig[input.resource_id] == input.claims.user_id
}

# SubmitWork — freelancer assigned to the gig
allow if {
    input.action == "SubmitWork"
    input.resource == "gig_assignee"
    input.claims.role == "freelancer"
    data.owners.gig_assignee[input.resource_id] == input.claims.user_id
}

# ApproveWork — client who owns the gig
allow if {
    input.action == "ApproveWork"
    input.resource == "gig"
    input.claims.role == "client"
    data.owners.gig[input.resource_id] == input.claims.user_id
}
