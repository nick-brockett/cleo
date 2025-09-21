# Health Metric Extraction API

### Hi, I would prefer to discuss my design choices in person, however I will write at a high-level my approach and hope it goes some way to explain what I have delivered!

### I felt that the stated requirements to extract weight and height metrics from free-text clinical notes lacked so much clarity that I could not evaluate the best coding solution or approach to take.

### The idea that I could not programmatically work out which weight (if there were more than one reference in the text) would be the correct one to report bothered me and I felt stumped to proceed without consulting and receiving more clarity.
### I asked myself the following question, what if the clinician had written in the notes the following example :-
`the patient stated their weight was 75Kg but when admitted their weight was found to be 95kg `

### Conscious of time I switched to examining the following instruction given in the test:-
`Items to consider:While implementing the above you should take into consideration key items such as validation, scalability, security and testing `

### Security & Scalability are easier and more comfortable for me to think through, so I created a new go project based on a previous one I had already completed found here
### https://github.com/nick-brockett/eagle-bank this project now sits next to it at  https://github.com/nick-brockett/cleo

### Performing some quick edits this gave me a starter for a docker application. I added a role based jwt access middleware section, picking a random role name example of CLINICAL-EDITOR. This meant that the POST/parse endpoint was secure, and I made sure my Router was protected to not receive more than 64KB in a single payload.
### I then conducted Postman tests to prove that I needed the Bearer Token set for appropriate responses. I used jwt.io website to concoct a sample token.

### I spent some time writing unit tests (not extensively) but I wanted to show enough test coverage of the adapter with a mocked out service for the parser, however that might be implemented.

### At this point I felt confident that I had enough example code to achieve the following items Scalability, Security and Testing

### The final step was basically to assume a quick regex expression (poor as it may seem) would have to suffice. So two files within the service layer (parser.go and parser_test.go) were written.
### The regex pattern and some of the utility methods were lifted straight from an internet regex generator, and then the test file was written by hand, conscious of time I omitted many unit tests which would have been needed if in any way this would be considered a viable production solution.

### Finally, I switched to Postman testing of various example hand typed dynamic clinical notes, whilst also running the api in my docker environment.
### I spotted some edge cases, which meant tweaking to the regex, and I was not happy with the error handling, as nothing
### was coming through in the docker logs, so went back to inject a logger into the service and made sure that when edge cases and illogical weights and heights were found that at least they were logged.

### approx 5.5 hours spent.











