package org.demo;

import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;

@Path("/log")
public class LogResource {

    @GET
    @Produces(MediaType.TEXT_PLAIN)
    public String log() {
        return "Hello from Quarkus REST";
    }
}
