# Course Enrollment

A fast course enrollment service written with Go and backed by PostgreSQL and RabbitMQ.

## Features

* Very fast. Handles 3500 req/s in a 1 Core (i7 4790K) / 1 GB VM
* Small memory footprint: About 13 MB of memory usage for 1500 courses and 1000 students.
* No disk bottleneck due to RabbitMQ (excluding RabbitMQ disk usage itself)
* Student and Admins (staff) endpoints
* Reserve Queues
* Sex lock on courses
* Partially horizontally scalable
* REST API
* JWT Authentication

## Running

At first, you need to spin up PostgreSQL and RabbitMQ. Assume that PostgreSQL is running on `127.0.0.1:5432` and
RabbitMQ is running on `127.0.0.1:5672`. At first, you need to create the required tables for the microservices to
function. To do so, you can simply redirect the `database.sql` file to the stdin of `psql` command like so:

```bash
psql -U hirbod course_enrollment < database.sql
```

One thing you have to note is that the core, caches the students in memory. So you cannot add students while this
program is running. Like who registers students to a university on an active course enrollment? So after you have
created your schema, add your students to `students` table.

Now, you have to start the services. It is advised to start them with this order:

### Batcher

These environment values must be supplied for the batcher:

* `DATABASE_URL`: PostgreSQL connection URL
* `RABBITMQ_ADDRESS`: RabbitMQ connection URL

For example:

```bash
export DATABASE_URL=postgres://hirbod:1234@127.0.0.1:5432/course_enrollment
export RABBITMQ_ADDRESS=amqp://hirbod:1234@127.0.0.1:5672/
```

### Enrollment Server

These environment values must be supplied for the enrollment server:

* `DATABASE_URL`: PostgreSQL connection URL
* `RABBITMQ_ADDRESS`: RabbitMQ connection URL
* `LISTEN_ADDRESS`: The address which the core expects the auth core to make requests to it.
* `LISTEN_PROTOCOL` (Optional): The protocol which the enrollment server excepts the auth core to make its requests in.
  The default is `tcp`. If two services are on a single operating system, `unix` is recommended.

Example of TCP listening:

```bash
export DATABASE_URL=postgres://hirbod:1234@127.0.0.1:5432/course_enrollment
export RABBITMQ_ADDRESS=amqp://hirbod:1234@127.0.0.1:5672/
export LISTEN_ADDRESS=0.0.0.0:23621
```

Example of unix sockets:

```bash
export DATABASE_URL=postgres://hirbod:1234@127.0.0.1:5432/course_enrollment
export RABBITMQ_ADDRESS=amqp://hirbod:1234@127.0.0.1:5672/
export LISTEN_PROTOCOL=unix
export LISTEN_ADDRESS=/dev/shm/CES.sock
```

### Authorization Core

These environment values must be supplied for the authorization core:

* `DATABASE_URL`: PostgreSQL connection URL
* `CORE_ADDRESS`: The address of the enrollment server. This can be a `unix://` link if the enrollment server is
  listening on a unix socket. Otherwise, plain gRPC over TCP is used.
* `LISTEN_ADDRESS`: The address which the authorization core expects the students and the staff to access it.
* `LISTEN_PROTOCOL` (Optional): The protocol which the authorization core excepts the users their requests in.
  The default is `tcp`. If you are using a reverse proxy (for example nginx), `unix` is recommended.
* `GIN_MODE` (Optional): Set this value to `release` if you are running this core in production. This will suppress the
  logs.

Example of everything over TCP:

```bash
export DATABASE_URL=postgres://hirbod:1234@127.0.0.1:5432/course_enrollment
export CORE_ADDRESS=127.0.0.1:23621
export LISTEN_ADDRESS=0.0.0.0:80
export GIN_MODE=release
```

Example of everything over unix sockets:

```bash
export DATABASE_URL=postgres://hirbod:1234@127.0.0.1:5432/course_enrollment
export CORE_ADDRESS=unix:///dev/shm/CES.sock
export LISTEN_ADDRESS=/dev/shm/auth.sock
export LISTEN_PROTOCOL=unix
export GIN_MODE=release
```

## Details

The overview of system is shown below:

![System Design](docs/design.svg)

### Authorization Core

The authorization core takes acts as a service which is exposed to the internet and users. Users and admins can login,
enroll, get courses etc. It simply uses a REST API. Users and admins have to login with their username and password.
JWT key is generated on each launch of application and nothing is cached. So users and staff can be added to database
while the authorization cores are running. It's worth nothing that the user data is stored inside the JWT and there is
no need for cache.

One cool aspect of authorization core is that it can be horizontally distributed. You can simply just spawn multiple
clients and connect each one to database and enrollment server. Just be aware that after authorization, each user must
request the corresponding server in which they logged in. This is because the JWT key of each instance is different from
others.

It's worth noting that database access for auth core is read only and the transport between auth core and enrollment
server is based on gRPC. The postman documentation is available in the docs folder.

### Enrollment Server

The enrollment server is the heart of the system. It handles all the requests related to courses and students.
Enrolling, getting the list of enrollments, force enrolling etc.

The key to speed in this service is the caching. Once the application is loaded, ALL the course data is cached. This
enables the service to do the calculations in memory. However, the downside is that if any student is added to system or
updated without using this service, the service must be restarted. But this scenario is incredibly rare; how often do
you have to add a student in mid-enrollment time?

This service exposes its API as gRPC. It needs read only access to database only once It's starting up. I also relys on
the broker to update the database (read the Batcher section for more info).

The enrollment server _can_ be horizontally distributed in some capacity. Each service needs to have distinct
departments from other running services. Each request from the authorization core should specifically go to the
corresponding enrollment service. The authorization core should be also changed a little.

### Batcher

Batcher is a service which its solely is to read the changes from the RabbitMQ broker and apply them into database. It
could be improved a little bit in terms of reliability with not enabling auto ack while subscribing to the queue and ack
each message when we are sure that the data is written to disk.

One important aspect of the batcher is that you cannot restart the Enrollment Service due to batcher. Batcher needs to
empty the queue in RabbitMQ because Enrollment Server does not detect changes in database while it is up. Also, only and
only one batcher needs to be up. The read order is important from the message queue.

## Benchmark in VirtualBox

### Specifications

* One core of Intel i7-4790K
* 1 GB RAM
* HDD
* VirtualBox 7.0.10
* Ubutnu 22.04 Kernel 6.2 Host
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

#### Database Batcher

```bash
export DATABASE_URL=postgres://hirbod:1234@127.0.0.1:5432/course_enrollment
export RABBITMQ_ADDRESS=amqp://hirbod:1234@127.0.0.1:5672/
./DatabaseBatcher
```

#### Course Enrollment Server

```bash
export DATABASE_URL=postgres://hirbod:1234@127.0.0.1:5432/course_enrollment
export RABBITMQ_ADDRESS=amqp://hirbod:1234@127.0.0.1:5672/
export LISTEN_PROTOCOL=unix
export LISTEN_ADDRESS=/dev/shm/CES.sock
./CourseEnrollmentServer
```

#### Auth Core

```bash
export DATABASE_URL=postgres://hirbod:1234@127.0.0.1:5432/course_enrollment
export CORE_ADDRESS=unix:///dev/shm/CES.sock
export LISTEN_ADDRESS=127.0.0.1:12345
export GIN_MODE=release
./AuthCore
```

### Benchmark

Set the VM network adaptor as host only. This will allow to connect to guest VM. Let the guest IP
be `192.168.56.101`. Use the following environment variables to setup the benchmark and run it. You need to restart the
core before pressing enter when the program asks you.

```bash
export STD_COUNT=1000
export DATABASE_URL=postgres://hirbod:1234@192.168.56.101:5432/course_enrollment
export AUTH_CORE_URL=http://192.168.56.101:12345
./Benchmark
```

#### Results

On my computer the result is like this:

```
INFO[0026] Done in 1.880980622s                         
INFO[0026] Total requests: 7000                         
INFO[0026] Failed requests: 1504
```

**Note**: Windows tend to have some problems with VirtualBox network driver which makes it slow. I could only get about
700 req/s on Windows.