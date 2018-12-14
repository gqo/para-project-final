# Made by: Graeme Ferguson
# Golang CRUD Failure-detecting Frontend & Multithreaded Backend
Third part of semester long project for Parallel and Distributed Systems. The objective of this part was to extend the frontend to detect failure of the backend and report to stdout and to extend the backend to handle requests in a concurrent fashion.
## Getting Started
### Prerequisites
* A Linux operating system.

* Golang interpreter and source code.

### Building
Unzip project into a directory of your choosing. All files and folders must remain in the original hierarchic order. Navigate inside the /frontend directory in console and type:

```go build frontend.go dbcommands.go```

This will create an executable file within the directory titled "frontend".
### Running
To run the backend, navigate to the /backend directory in console and type:

```go run backend.go```

To run the frontend, navigate to the /frontend directory in console and type:

```./frontend```
## Testing
### Initial
In regards to the backend:

* Without a listen flag, the backend will automatically bind to the port 8090 on your localhost domain.

In regards to the frontend:

* Without a listen flag, the frontend will automatically bind to port 8080 on your localhost domain for listening to HTTP requests. As such, within a browser of your choosing, go to http://localhost:8080 to view the root page.

* Without a backend flag, the frontend will automatically send all requests for the backend to the address "localhost:8090".

### Setting a Listen Flag
In regards to the backend:
* Using the command flag `--listen`, the backend will bind to a custom port of your choosing. For example, the backend will bind to port 8091 if you run `go run backend.go --listen 8091` in your /backend directory. As such, the backend will then listen for requests on port 8091.

In regards to the frontend:
* Using the command flag `--listen`, the frontend will bind to a custom port of your choosing for listening to HTTP requests. For example, the frontend will bind to port 8092 if you run `./frontend --listen 8092` in your /frontend directory. As such, you will be able to find the webserver at http://localhost:8092.
### Setting a Backend Flag
Using the command flag `--backend`, the frontend will bind to a custom address(ip, port) of your choosing for the sending of requests to the backend. For example, the frontend will send requests to the address 10.18.207.22:8100 if you run `./frontend --backend 10.18.207.22:8100` in your /frontend directory. As such, the frontend will send all requests to the 10.18.207.22:8100 address.
### Creation
On the root page of the webserver is a "Create!" link that, if followed, will take you to a /create/id page of the webserver. At this page, you can fill out the form accordingly. All form values are required and the rating form value only accepts integers from 1 to 10. Click the "Save" button to save item or the "Cancel" link to discard the item, both return automatically to the root page.
### Editing
On the root page of the webserver is a "Edit" link under each review item that, if followed, will take you to an /edit/id page of the webserver. At this page, you will be able to edit the previously entered form values accordingly. All form values have the same specifications as a /create/id page and the "Save" button and "Cancel" link work in the same fashion. Additionally, if one attempts to edit an item that is not recognized, they will be automatically redirected to a /create/id page.
### Deletion
On the root page of the webserver is a "Delete" link under each review item that, if followed, will delete the aforementioned review item. If one attempts to delete an item that is not recognized, they will be automatically redirected to the root page.
### Failure
If you run the frontend while the backend is not currently running, a failure detection message will be printed every ~5 seconds. 
## Assignment Information
### State of Assignment
Assignment was completed to specifications as listed in project3.pdf and, to the author's knowledge, meets all specifications.
### Resources Used
In particular: 

* https://golang.org/pkg/net/

In general: 

* https://golang.org/

### Design Decisions
I used a single RW Mutex to wrap my map in which my database stores its data. I choose a RW Mutex due to its relative simplicity while still allowing concurrent read calls that I felt would be important for a CRUD server given that multiple frontends may simple want to read the data at the same time. I attempted to ensure some finer granularity by wrapping the minimal portions of my backend code said mutex allowing the majority of the logic to be performed outside of the locked data structure. 

The major con of a RW Mutex is that I'm still wrapping the entire map in a mutex at points and thus it is unnessarily big. I would have liked to use a thread that manages the state of the server by owning a map inside of itself and accepting stateful changes through channels but, due to the complexity, avoided this method for this assignment. I plan on at least attempting to implement this in my completion of part 4.

In terms of performance metrics, mean latency seemed the best (as stated within the project3.pdf itself). I wanted a balanced amount of POST and GET requests so I wrote a list of targets that balanced that and ran the the following vegeta commands:

Combined attack:

```
$ vegeta attack -rate=10 -workers=50 -duration=30s -targets=targets.txt | vegeta report
Requests      [total, rate]            300, 10.03
Duration      [total, attack, wait]    29.906003203s, 29.900132347s, 5.870856ms
Latencies     [mean, 50, 95, 99, max]  4.623476ms, 4.202675ms, 8.80859ms, 10.119648ms, 11.415049ms
Bytes In      [total, mean]            5896050, 19653.50
Bytes Out     [total, mean]            6900, 23.00
Success       [ratio]                  100.00%
Status Codes  [code:count]             200:300  
Error Set:
```
POST-based attack:

```
$ vegeta attack -rate=10 -workers=50 -duration=30s -targets=targets.txt | vegeta report
Requests      [total, rate]            300, 10.03
Duration      [total, attack, wait]    29.914335834s, 29.900133732s, 14.202102ms
Latencies     [mean, 50, 95, 99, max]  11.259338ms, 11.517471ms, 15.804722ms, 17.163882ms, 17.679817ms
Bytes In      [total, mean]            22547786, 75159.29
Bytes Out     [total, mean]            13800, 46.00
Success       [ratio]                  100.00%
Status Codes  [code:count]             200:300  
Error Set:
```
GET-based attack:

```
$ vegeta attack -rate=10 -workers=50 -duration=30s -targets=targets.txt | vegeta report
Requests      [total, rate]            300, 10.03
Duration      [total, attack, wait]    29.900839034s, 29.900077469s, 761.565µs
Latencies     [mean, 50, 95, 99, max]  1.333649ms, 1.288811ms, 1.983167ms, 2.43869ms, 3.38638ms
Bytes In      [total, mean]            344400, 1148.00
Bytes Out     [total, mean]            0, 0.00
Success       [ratio]                  100.00%
Status Codes  [code:count]             200:300  
Error Set:
```
As you can see in these reports, my mean latencies vary wildly depending on the composition of the attack. GET requests are processed relatively quickly at 1.33ms while POST requests are processed at a languid 11.26ms. It is clear from this that concurrent writes result in a much larger delay than concurrent reads, most likely due to the RW Mutex only exclusively locking the map for writes. The combined attack has a mean latency of 4.62ms which is respectable but I cannot shake the assumption that with a concurrency approach that lacked mutexes I would be able to bring the mean latency closer to that of the GET-based attack.

I implemented my failure detector using the ping-ack approach with a wait time of 5 seconds between each ping. My approach additionally waits up to 5 seconds for a response upon writing to the backend. If my approach waits for on read, it will skip the wait time before initiating a ping again. I considered heartbeat but felt like ping-ack would be a less brittle approach for this part of the assignment as it only featured multiple frontends. However, it may be useful to consider heartbeat when adding multiple backends.

While not related to the purpose of the assignment, I discovered an interesting bug in my program using vegeta that, due to the fact that I send the entire data structure on a ReadAll call, I was quickly overloading the buffer I created to read responses into on my frontend with a map of sufficient size. As such, I increased the buffer size form 1024 to 8096 bytes such that it would be handle larger map sizes but, obviously, there is still a limit on that size. As such, I plan on implementing a more clever method of dealing with my ReadAll calls in the future.
### Additional Thoughts
Instructions were clear. Time allotment was ample. Support was not used.