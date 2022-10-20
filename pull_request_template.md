# 🌻 Description 🌻

This is my Quotes RESTful API that is deployed to Google Cloud Run. It uses the Golang library gqlgen to create GraphQL servers to query results from the REST API and add new additions to the GCP Postgres database. 

🎟️ GitHub Ticket #255 🎟️

## 🌸 How To Test 🌸

- in your terminal cd into intended directory you want to run the project
- git clone (https://github.com/KaleseCarpenter/quotes-starter.git) {New file name)
- checkout mutations branch 
- cd into quotes-starter folder
- code .
- in your terminal run -> go run server.go
- in your terminal -> command click http://localhost:8080/ for GraphQL playground to test queries

- You can copy and paste this below to run the queries and mutations with the provided ID in GraphiQL <br />
query GetQuote { 
  getRandomQuote { 
    id 
    quote 
    author 
  } <br />
  getQuoteByID(id: "2560b862-6941-4737-8a21-1e96805f786a") { <br />
    id <br />
    quote <br />
    author <br />
  } <br />
} <br />

mutation postQuote { <br />
  postQuote(input: {quote: "Happy Coding", author: "Code Reviewer"}) { <br />
    id <br />
    quote <br />
    author <br />
  } <br />
} <br />

mutation DeleteQuote { <br />
  deleteQuote(id: "2560b862-6941-4737-8a21-1e96805f786a")<br />
} <br />


## 🧪 How This Has Been Tested 🧪

These are the tests that I ran to verify my changes:

- [ ] I queried my results in GraphiQL and successfully returned a Random Quote.
- [ ] I queried my results in GraphiQL and successfully returned a Quote By ID.
- [ ] I checked my error handling my trying to query a quote with an unknown ID and returned "error: 404 Quote ID Not Found."
- [ ] I added a mutation in GraphiQL to Post and I successfully Posted a Quote.
- [ ] I tried posting a quote with less than 3 characters for the Quote or Author and returned "400 Bad Request."
- [ ] I added a mutation in GraphiQL to Delete a Quote By ID and it returns "Bruh, Delete was successful!"
- [ ] After deleting a quote, I tried to search that quote by ID again and it returns "error: 404 Quote ID Not Found." 

# Checklist ✅:

Requirements❗️❗

- [ ] My code enables a mutation for creating a quote
- [ ] My code enables a mutation for deleting a quote
- [ ] Schema has a type mutation with two entry points for Post and Delete
- [ ] Resolvers update the db through your existing REST API
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] Work is PR'd
