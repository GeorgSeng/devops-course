package org.demo;

import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;

import io.quarkus.logging.Log;
import io.vertx.core.file.OpenOptions;

import java.io.File;
import java.nio.charset.Charset;
import java.nio.file.Files;
import java.nio.file.OpenOption;
import java.nio.file.StandardOpenOption;
import java.time.temporal.ChronoUnit;

import org.demo.Main.StartUpTime;

@Path("/status")
public class StatusResource {

    @GET
    @Produces(MediaType.TEXT_PLAIN)
    public String status() {
        io.quarkus.logging.Log.info(String.format("jvmStartTime: %s", StartUpTime.getStartUpTime()));
       var currentTime = java.time.Instant.now();
       Log.info(currentTime);
       java.time.Duration duration = java.time.Duration.between(StartUpTime.getStartUpTime(), currentTime);
       //var upTimeInHours = duration.toHours();
       var upTimeInHours = duration.toSeconds();
       var usableSpace = new File("/").getUsableSpace() / (1024 * 1024); // convert bytes to mb
       var logMsg = String.format("%s: uptime %s s hours, free disk in root: %d MBytes\n", currentTime.truncatedTo(ChronoUnit.SECONDS), upTimeInHours, usableSpace);
        
        try {
            var path = java.nio.file.Path.of("./externalData/vstorage");
            Files.write(path, logMsg.getBytes(), StandardOpenOption.CREATE,StandardOpenOption.APPEND);
        } catch (Exception e) {
            Log.error("Could not write to the vStorage File");
            Log.error(e);
        }

       return logMsg;
    }
}
