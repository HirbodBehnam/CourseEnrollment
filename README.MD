# Course Enrollment

A fast course enrollment service written with Go and backed by PostgreSQL and RabbitMQ.

## Features

* Very fast. Handles 700 req/s in a 1 Core (i7 4790K) / 1 GB VM
* Student and Admins (staff) endpoints
* Reserve Queues
* REST API
* JWT Authentication

## Benchmark in VirtualBox

### Specifications

* One core of Intel i7-4790K
* 1 GB RAM
* HDD
* VirtualBox 7.0.10
* Windows 10 Host
* Ubuntu 22.04 Kernel 5.15 Guest
* CourseEnrollment compiled with GO 1.21.1

### Setup

I used docker to setup Postgres and RabbitMQ. To do so, I used the following commands:

```bash
docker run -itd -e POSTGRES_USER=hirbod -e POSTGRES_PASSWORD=1234 -p 5432:5432 -v ~/data:/var/lib/postgresql/data --name postgresql postgres
docker cp database.sql postgres:/ # Move all database schema build files to data folder
docker exec -it postgresql bash
```

This will open a terminal inside the docker container. In the docker container, use `psql -U hirbod` to open a postgres
command line shell. In the shell use the following command to create a new database for course enrollment and then quit.

```postgresql
CREATE DATABASE course_enrollment;
\q
```

Now to create the schema use the following command in bash:

```bash
psql -U hirbod course_enrollment < database.sql
# Or add additional mock data
```

Use exit to exit the container terminal.

Next, we have to setup RabbitMQ. Use the following command:

```bash
docker run -itd --hostname course-broker --name course-broker -e RABBITMQ_DEFAULT_USER=hirbod -e RABBITMQ_DEFAULT_PASS=1234 -p 5672:5672 -v ~/rabbitmq:/var/lib/rabbitmq rabbitmq:3
```

Finally, we have to run each service.

```bash
export DATABASE_URL=postgres://hirbod:1234@127.0.0.1:5432/course_enrollment
export RABBITMQ_ADDRESS=amqp://hirbod:1234@127.0.0.1:5672/
./DatabaseBatcher
```

```bash
export DATABASE_URL=postgres://hirbod:1234@127.0.0.1:5432/course_enrollment
export RABBITMQ_ADDRESS=amqp://hirbod:1234@127.0.0.1:5672/
export LISTEN_PROTOCOL=unix
export LISTEN_ADDRESS=/dev/shm/CES.sock
./CourseEnrollmentServer
```

```bash
export DATABASE_URL=postgres://hirbod:1234@127.0.0.1:5432/course_enrollment
export CORE_ADDRESS=unix:///dev/shm/CES.sock
export LISTEN_ADDRESS=127.0.0.1:12345
export GIN_MODE=release
./AuthCore
```

### Benchmark

```bash
export STD_COUNT=12000

```