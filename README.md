# generic-csp-go
An experiment with new generics available in Go 1.18+ inspired by CSP chapter in "Classic Computer Science Problems" by David Kopec.

# Learnings
I hit a problem with struct embedding and generics that I'm surprised about, and makes me think generics in Go are still kind of half-baked! You can view that version of the implementation at [on this branch](). I'm hoping I was just holding it wrong, but some research and [simple repro attempts](https://gotipplay.golang.org/p/M8vnSG9KYC0) outside the context of this project appear to indicate it's a legit problem :(

If you see what I'm doing wrong there, let me know! In the meantime, I've added a workaround in the form of a `csp.Satisfied[V, D comparable]` function that's not bound to the `Constraint[V, D comparable]` type that can be passed into a `Problem[V, D comparable]` instance of a particular type.
