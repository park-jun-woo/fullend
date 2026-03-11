package authz

default allow = false

# Allow all actions when DISABLE_AUTHZ is set (handled in Go code).
# This is a minimal default policy — users should replace with their own.
allow {
    input.user.id != 0
}
