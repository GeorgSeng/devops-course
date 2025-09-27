# DevOps

This project consists out of 3 services, where only service1 is exposed to the host. Service2 is only offers a basic api to get it's status which gets logged to storage. The storage service recives a log message via a post rquest and anc can send the logs back via a get request. 

## Service1
- Langauge: Go
- Framework: Gin
- Folder: `./service1`
- How to start: `go run .`

## Service2
- Langauge: Java
- Framework: Quarkus
- Folder: ./service2
- How to start: mvn quarkus:dev 

## Storage
- Langauge: Go
- Framework: Gin
- Folder: `./storage`
- How to start: `go run .`

## Dockerfiles
The docker files to build the services are in there folders.

## Docker Compose
The docker compose file is in the root folder and has the following tasks. Create two virtual networks, one connected to the bridge and one that is isolated for the backend services (service2, storage). There is also a volume defined for the storage service such that the logs are persistent over multiple startups. In order to set up the log files used by service1 and service2, a setup step was added to the compose file. It's basically an Alpine Linux image that runs the shell command touch to create the vstorage file for logs. In theory, the setup stage is not needed since the task now specifies that the folder should be mounted to the container instead of binding a file. But I think it's interesting to see how an image can be used to script some setup steps. Storage and service2 offer a health API, where the API of service2 was provided by Quarkus; the other one was written by me. The health endpoints are then used by Docker to verify that the services are up and running in order to avoid issues that a service would try to access a service that is still not ready.

### Building
Building the system takes a long time since I use multi stage builds. This means that first an image with the compiler and build dependencies is fetched and used to create the executable. Then a second image is pulled that will act as runtime. The major benefit of this strategy is that the final image is much smaller and it has fewer dependencies, therefore the risk of a Software Supply Chain (SSC) related security incident gets reduced. Unfortunately, the build of service2 takes a long time since the Java build system has to be pulled and then the Maven dependencies have to be pulled as well. Quarkus itself has a command that builds a docker image but it assumes that it's run on the host machine where the dependencies are already available. Since the task was to make everything generic and don't assume anything about the host, it was necessary to remake the dockerfile in such a way that it first sets up a build environment and then copies the files.

`docker-compose up --build`

### Clean up
`docker-compose down --volumes`
