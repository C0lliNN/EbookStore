# Todo

* Create endpoint to download a book
* Commit
* Deploy code
* Create a Nice README.md (Include the whole thing there)
* Commit

## Nice to have
* Use `interface{}` in the Error details field. Then, when `*validation.ValidationErrors` occurs return a slice with nicely formatted data

ADMIN TOKEN: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZW1haWwiOiJyYXBoYWVsMkB0ZXN0LmNvbSIsImlkIjoiY2JiOGE3N2YtN2Y3Yy00ZGJlLWI4ZGQtMTNlYjhlN2Q2YTQ1IiwibmFtZSI6IlJhcGhhZWwgQ29sbGluIn0.rR0uCLihDtK6m7Ck4WgpmeguVfPI_N4x9_t2Tp7xBT0
CUSTOMER TOKEN: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImVtYWlsIjoicmFwaGFlbDNAdGVzdC5jb20iLCJpZCI6IjAxYmYzYTcxLWRmNWYtNGQ4Yi05YTA0LWEzYWExOGU1YmQyZiIsIm5hbWUiOiJSYXBoYWVsIENvbGxpbiJ9.Uq-rXsNox5UVrlRqd20iUYdVpMxWIiqgbkKGP38brpA

stripe trigger --add payment_intent:metadata.orderID=some-id payment_intent.succeeded