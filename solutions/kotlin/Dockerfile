ARG GRADLE_IMAGE=gradle:8.7.0-jdk17
ARG JAVA_IMAGE=eclipse-temurin:17-jre

# Stage 1: create gradle cache
FROM ${GRADLE_IMAGE} as gradle-cache
ENV BUILD_SRC /usr/src/app/
ENV GRADLE_USER_HOME /home/gradle/cache/
RUN mkdir -p ${BUILD_SRC}
RUN mkdir -p ${GRADLE_USER_HOME}
COPY gradle ${BUILD_SRC}gradle
COPY gradlew build.gradle.kts settings.gradle.kts ${BUILD_SRC}
WORKDIR ${BUILD_SRC}
RUN gradle build -i --build-cache || return 0 && echo "DO NOT PANIC! Error is ok at 1st stage"

# Stage 2: build
FROM gradle-cache as builder
ENV BUILD_SRC /usr/src/app/
COPY . ${BUILD_SRC}
WORKDIR ${BUILD_SRC}
RUN gradle bootJar -i --stacktrace --build-cache

# Stage 2.5: Unpack JAR
ENV EXTRACTED /opt/app
RUN mkdir -p ${EXTRACTED} && (cd ${EXTRACTED}; java -Djarmode=layertools -jar ${BUILD_SRC}build/libs/*.jar extract)
RUN du -hs ${EXTRACTED}/*

# Stage 3: run
FROM ${JAVA_IMAGE}
ENV EXTRACTED /opt/app
USER root
WORKDIR /opt/app
COPY --from=builder ${EXTRACTED}/dependencies/ ./
COPY --from=builder ${EXTRACTED}/spring-boot-loader/ ./
COPY --from=builder ${EXTRACTED}/snapshot-dependencies/ ./
COPY --from=builder ${EXTRACTED}/application/ ./
ENTRYPOINT ["java","org.springframework.boot.loader.launch.JarLauncher"]
