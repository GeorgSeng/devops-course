package org.demo.api

import jakarta.ws.rs.GET
import jakarta.ws.rs.Path
import jakarta.ws.rs.Produces
import jakarta.ws.rs.core.MediaType
import jakarta.inject.Inject
import org.demo.service.DockerService

@Path("/docker")
class DockerMetrics {
    @Inject
    lateinit var dockerService: DockerService

    @GET()
    @Path("container")
    @Produces(MediaType.TEXT_PLAIN)
    fun getAllContainers() = dockerService.getAllRunningContainers()
}