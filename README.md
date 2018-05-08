# big_mammal_from_the_north_assignment

Solution for candidate test BE v 1.2 May 2018
## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

go language pack
ApacheBench

### Installing

-run go install in the directory
-run the produced executable

## Running the tests

the given tests where the following:
ab -c 3 -n 10 'http://127.0.0.1:8080/player/1'
ab -c 3 -n 10 'http://127.0.0.1:8080/player/2'

curl ‘http://127.0.0.1:8080/statistic

### To test manually

to test manually:
```
open http://127.0.0.1:8080/player/N
where N is the number of a player```

then ```open http://127.0.0.1:8080/statistic to see a summary of the results```
