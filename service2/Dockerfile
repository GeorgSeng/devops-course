FROM maven:3.9 as builder

WORKDIR /usr/src/app
COPY . .
WORKDIR /usr/src/app/service2
RUN ls
#RUN mvn clean quarkus:add-extension -Dextensions='container-image-docker'
#RUN mvn clean install -Dquarkus.container-image.build=true
RUN mvn clean install
RUN ls
RUN ls ./target

# this is the code of the dockerfile.jvm provided by quarkus
FROM registry.access.redhat.com/ubi9/openjdk-21:1.23
ENV LANGUAGE='en_US:en'


# We make four distinct layers so if there are application changes the library layers can be re-used
RUN ls
COPY --from=builder --chown=185 /usr/src/app/service2/target/quarkus-app/lib/ /deployments/lib/
COPY --from=builder --chown=185 /usr/src/app/service2/target/quarkus-app/*.jar /deployments/
COPY --from=builder --chown=185 /usr/src/app/service2/target/quarkus-app/app/ /deployments/app/
COPY --from=builder --chown=185 /usr/src/app/service2/target/quarkus-app/quarkus/ /deployments/quarkus/

EXPOSE 8080
USER 185
ENV JAVA_OPTS_APPEND="-Dquarkus.http.host=0.0.0.0 -Djava.util.logging.manager=org.jboss.logmanager.LogManager"
ENV JAVA_APP_JAR="/deployments/quarkus-run.jar"

ENTRYPOINT [ "/opt/jboss/container/java/run/run-java.sh" ]