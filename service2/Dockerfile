FROM maven:3.9 as builder

WORKDIR /usr/src/app
COPY . .
WORKDIR /usr/src/app/service2
RUN mvn clean install

# this is the a modifyed version of the code provided by Quarkus in the dockerfile.jvm
# Nomrlay a quarkus contianer is created by mavne that compiles the application and docker build is than just loading the files as shown belwo.
# Since the goal of the exercise is that the system can be set up on any machine without maven or java installed i had to used a builder stage to create the nececerry code.

FROM registry.access.redhat.com/ubi9/openjdk-21:1.23
ENV LANGUAGE='en_US:en'


# We make four distinct layers so if there are application changes the library layers can be re-used
COPY --from=builder --chown=185 /usr/src/app/service2/target/quarkus-app/lib/ /deployments/lib/
COPY --from=builder --chown=185 /usr/src/app/service2/target/quarkus-app/*.jar /deployments/
COPY --from=builder --chown=185 /usr/src/app/service2/target/quarkus-app/app/ /deployments/app/
COPY --from=builder --chown=185 /usr/src/app/service2/target/quarkus-app/quarkus/ /deployments/quarkus/

EXPOSE 8080
USER 185
ENV JAVA_OPTS_APPEND="-Dquarkus.http.host=0.0.0.0 -Djava.util.logging.manager=org.jboss.logmanager.LogManager"
ENV JAVA_APP_JAR="/deployments/quarkus-run.jar"

ENTRYPOINT [ "/opt/jboss/container/java/run/run-java.sh" ]