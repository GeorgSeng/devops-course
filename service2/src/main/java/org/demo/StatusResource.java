package org.demo;

import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.client.ClientBuilder;
import jakarta.ws.rs.core.MediaType;

import io.quarkus.logging.Log;
import io.vertx.core.file.OpenOptions;

import java.io.File;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.net.http.HttpRequest.BodyPublisher;
import java.net.http.HttpRequest.BodyPublishers;
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
       var upTimeInHours = duration.toSeconds()/60.0;
       var usableSpace = new File("/").getUsableSpace() / (1024 * 1024); // convert bytes to mb
       var logMsg = String.format("%s: uptime %.8f hours, free disk in root: %d MBytes", currentTime.truncatedTo(ChronoUnit.SECONDS), upTimeInHours, usableSpace);
        
       // send the status to the storage service
       var logPost = HttpRequest
        .newBuilder()
        .uri(URI.create("http://storage:8080/log"))
        .POST(BodyPublishers.ofString(logMsg))
        .build();
        try {
            var httpClient = HttpClient.newHttpClient();
            httpClient.send(logPost, HttpResponse.BodyHandlers.discarding());
        } catch (Exception e) {
            Log.error("Could not wirte the log msg to the sotrage service");
            Log.error(e);
        }

       // wirte the status to the vstorage log file
        try {
            var path = java.nio.file.Path.of("./externalData/vstorage");
            Files.write(path, (logMsg+"\n").getBytes(), StandardOpenOption.CREATE,StandardOpenOption.APPEND);
        } catch (Exception e) {
            Log.error("Could not write to the vStorage File");
            Log.error(e);
        }

       return logMsg;
    }
}
