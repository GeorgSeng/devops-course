package org.demo;

import io.quarkus.runtime.annotations.QuarkusMain;

import java.time.Instant;

import io.quarkus.logging.Log;
import io.quarkus.runtime.Quarkus;

@QuarkusMain
public class Main {

    public final class StartUpTime {
        private static Instant StartUpTimeInstance;
        private StartUpTime(){}
        public static Instant getStartUpTime(){
            if (StartUpTimeInstance == null){
                StartUpTimeInstance = java.time.Instant.now();
            }
            return StartUpTimeInstance;
        }
    }
    public static void main(String[] args) {
        Log.info(String.format("Started at %s", StartUpTime.getStartUpTime()));
        Quarkus.run(args);
    }
}
