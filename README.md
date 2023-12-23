# Unattended Programming Test: The Test Signer

| Method | URL Pattern       | Handler                | Action                                                                         |
|--------|-------------------|------------------------|--------------------------------------------------------------------------------|
| GET    | /ping             | pingHandler            |                                                                                |
| POST   | /signature        | createSignatureHandler | Accepts a user JWT, questions and answers, and creates and returns a signature |
| POST   | /signature/verify | verifySignatureHandler | Accepts a user JWT and signature and returns                                   |