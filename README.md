# Made by: Graeme Ferguson
# Golang CRUD Failure-detecting Frontend & Multithreaded, Replicated Backend
Fourth, and final, part of semester long project for Parallel and Distributed Systems. The objective of this part was to implement a replication strategy on the backend such that nodes of the backend could fail without bringing the whole system down.
## Getting Started
### Prerequisites
* A Linux operating system.

* Golang interpreter and source code.

### Building
Unzip project into a directory of your choosing. All files and folders must remain in the original hierarchic order. Navigate inside the /frontend directory in console and type:

```go build frontend.go dbcommands.go```

This will create an executable file within the directory titled "frontend".

Navigate inside the /backend directory in console and type:

```go build backend.go paxos.go```

This will create an executable file within the directory titled "backend".
### Running
To run the backend, navigate to the /backend directory in console and type with appropriate flags:

```./backend```

To run the frontend, navigate to the /frontend directory in console and type with appropriate flags:

```./frontend```

The backends must be started before the frontend. An example of a typical running scenario is as follows:

```bash
./backend --id 1 --listen 8090 --backend :8091,:8092 &
./backend --id 2 --listen 8091 --backend :8090,:8092 &
./backend --id 3 --listen 8092 --backend :8090,:8091 &
./frontend --listen 8081 --backend :8090,:8091,:8092
```
## Testing
### Initial
In regards to the backend:

* Without a listen flag, the backend will automatically bind to the port 8090 on your localhost domain.

* Mandatory flags:

    * Without a backend flag, the backend will not have peers and will not dial them

    * Without an id flag, the backend will not have an id to identify itself by

In regards to the frontend:

* Without a listen flag, the frontend will automatically bind to port 8080 on your localhost domain for listening to HTTP requests. As such, within a browser of your choosing, go to http://localhost:8080 to view the root page.

* Mandatory flags:

    * Without a backend flag, the frontend will only be able to dial one backend.

### Setting a Listen Flag
In regards to the backend:
* Using the command flag `--listen`, the backend will bind to a custom port of your choosing. For example, the backend will bind to port 8091 if you run `./backend --listen 8091` in your /backend directory. As such, the backend will then listen for requests on port 8091.

In regards to the frontend:
* Using the command flag `--listen`, the frontend will bind to a custom port of your choosing for listening to HTTP requests. For example, the frontend will bind to port 8092 if you run `./frontend --listen 8092` in your /frontend directory. As such, you will be able to find the webserver at http://localhost:8092.
### Setting a Backend Flag
In regards to the backend:
* Using the command flag `--backend`, the backend will be aware of its peer nodes on the network. For example, the backend will contact addressed `localhost:8091` and `localhost:8092` if you run `./backend --backend :8091,:8092` in your /backend directory. Addresses must be given in the form of comma separated list.


In regards to the frontend:
Using the command flag `--backend`, the frontend will contact one the custom addresses(ip, port) of your choosing for the sending of requests to the backend. For example, the frontend will send requests to the address localhost:8090 or localhost:8091 if you run `./frontend --backend :8090,:8091` in your /frontend directory. Addresses must be given in the form of comma separated list.
### Setting an ID Flag
All backends must be run with an appropriate ID flag. Flags must be unique 32-bit integers, beginning at one. For example, the backend will identify as node 1 if you run `./backend --id 1`
### Creation
On the root page of the webserver is a "Create!" link that, if followed, will take you to a /create/id page of the webserver. At this page, you can fill out the form accordingly. All form values are required and the rating form value only accepts integers from 1 to 10. Click the "Save" button to save item or the "Cancel" link to discard the item, both return automatically to the root page.
### Editing
On the root page of the webserver is a "Edit" link under each review item that, if followed, will take you to an /edit/id page of the webserver. At this page, you will be able to edit the previously entered form values accordingly. All form values have the same specifications as a /create/id page and the "Save" button and "Cancel" link work in the same fashion. Additionally, if one attempts to edit an item that is not recognized, they will be automatically redirected to a /create/id page.
### Deletion
On the root page of the webserver is a "Delete" link under each review item that, if followed, will delete the aforementioned review item. If one attempts to delete an item that is not recognized, they will be automatically redirected to the root page.
### Failure
The frontend will randomly ping one of the backends every 5 seconds; printing a message if failure is detected. Additonally, one can test the backend replication by stopping one after a normal start. As long as the number of nodes does not fall below quorum, one should still be able to write. However, if the number of nodes falls below quorum, any date operations (such as create, update, and delete) will fail silently, except for a log printed on the frontend.
## Assignment Information
### State of Assignment
Currently, the system is incapable of correctly recovering from the failure of nodes. Nodes can be restarted but they will forever have inconsistent data. A possible solution would be to save if a node fails due to ping and then send it a special start message to contact a functioning node and receive the correct data from it.
### Resources Used
In particular: 

* https://golang.org/pkg/net/

In general: 

* https://golang.org/

### Design Decisions
I chose the Paxos replication strategy. I disliked the constant pinging of Raft and liked the idea of a less leader dependent replication strategy. Viewstamped replication seemed too niche compared to the resources for Raft and Paxos. 

The cons of Paxos was the need for implementation of a large set of functions to handle the correct types of messages. I also did not implement the optimized version of MultiPaxos with leaders so I complete the very costly version of Paxos every time data is updated.

The pros of Paxos were most likely found in the difficulty of implementation. I had to be intimately familiar with how my messages could fail or be denied to implement a working version while Raft may have left some abstraction. Additionally, the leaderless system allows for the frontend to contact any node in particular without the need for rerouting data.

My system can handle up to n/2 failures of nodes but can not recover correctly from said failure.
### Additional Thoughts
Instructions were clear. Time allotment was ample. Maybe one more day would have been nice but that is a personal issue. Support was not used.