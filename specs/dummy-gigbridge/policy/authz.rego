package authz

# @ownership gig: gigs.client_id
# @ownership gig_assignee: gigs.freelancer_id
# @ownership proposal: proposals.freelancer_id

default allow = false

# PublishGig: client who owns the gig
allow if {
    input.action == "publish"
    input.resource == "gig"
    input.user.role == "client"
    input.user.id == input.resource_owner
}

# SubmitProposal: freelancer, cannot submit to own gig
allow if {
    input.action == "submit_proposal"
    input.resource == "gig"
    input.user.role == "freelancer"
    input.user.id != input.resource_owner
}

# AcceptProposal: client who owns the gig
allow if {
    input.action == "accept"
    input.resource == "gig"
    input.user.role == "client"
    input.user.id == input.resource_owner
}

# RejectProposal: client who owns the gig
allow if {
    input.action == "reject"
    input.resource == "gig"
    input.user.role == "client"
    input.user.id == input.resource_owner
}

# SubmitWork: freelancer who is assigned to the gig
allow if {
    input.action == "submit_work"
    input.resource == "gig_assignee"
    input.user.role == "freelancer"
    input.user.id == input.resource_owner
}

# ApproveWork: client who owns the gig
allow if {
    input.action == "approve"
    input.resource == "gig"
    input.user.role == "client"
    input.user.id == input.resource_owner
}

# RaiseDispute: client who owns the gig
allow if {
    input.action == "dispute"
    input.resource == "gig"
    input.user.role == "client"
    input.user.id == input.resource_owner
}
