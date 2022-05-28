# generic-csp-go
An experiment with new generics available in Go 1.18+ inspired by CSP chapter in "Classic Computer Science Problems" by David Kopec. Run `make` to solve the example problems.

# Learnings
My original implementation surfaced an unexpected behavior around structs with generic parameters and embedding. You can view that version of the implementation [on this branch](https://github.com/elireisman/generic-csp-go/tree/embedded_generic_struct_implementation). I'm hoping I was just holding it wrong, but some research and [simple repro attempts](https://gotipplay.golang.org/p/M8vnSG9KYC0) outside the context of this project appear to indicate it's a legit problem :(

If you see what I did wrong there, let me know!
