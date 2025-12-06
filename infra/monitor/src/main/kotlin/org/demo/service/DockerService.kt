package org.demo.service

import com.github.dockerjava.api.DockerClient
import com.github.dockerjava.api.model.Container
import jakarta.enterprise.context.ApplicationScoped
import jakarta.inject.Inject

@ApplicationScoped
class DockerService {
    @Inject
    lateinit var dockerClient: DockerClient

    fun getAllRunningContainers(): List<Container> {
        return dockerClient.listContainersCmd().exec()
    }
}