FROM golang:alpine AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /AuthCore ./cmd/AuthCore
RUN go build -o /Benchmark ./cmd/Benchmark
RUN go build -o /CourseEnrollmentServer ./cmd/CourseEnrollmentServer
RUN go build -o /DatabaseBatcher ./cmd/DatabaseBatcher

FROM scratch AS auth_core

COPY --from=build /AuthCore /AuthCore

ENTRYPOINT ["/AuthCore"]

FROM scratch AS benchmark

COPY --from=build /Benchmark /Benchmark

ENTRYPOINT ["/Benchmark"]

FROM scratch AS course_enrollment_server

COPY --from=build /CourseEnrollmentServer /CourseEnrollmentServer

ENTRYPOINT ["/CourseEnrollmentServer"]

FROM scratch AS database_batcher

COPY --from=build /DatabaseBatcher /DatabaseBatcher

ENTRYPOINT ["/DatabaseBatcher"]