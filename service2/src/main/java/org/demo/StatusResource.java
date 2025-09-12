package org.demo;

import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;

import java.io.File;
import java.time.temporal.ChronoUnit;

@Path("/status")
public class StatusResource {

    @GET
    @Produces(MediaType.TEXT_PLAIN)
    public String status() {
       long jvmStartTime = java.lang.management.ManagementFactory.getRuntimeMXBean().getStartTime();
       var startTime = java.time.Instant.ofEpochSecond(jvmStartTime);
       var currentTime = java.time.Instant.now();
       java.time.Duration duration = java.time.Duration.between(startTime, currentTime);
       var usableSpace = new File("/").getUsableSpace() / (1024 * 1024); // convert bytes to mb
       return String.format("%s: uptime %s hours, free disk in root: %d MBytes", currentTime.truncatedTo(ChronoUnit.SECONDS), duration,usableSpace);
    }
}
