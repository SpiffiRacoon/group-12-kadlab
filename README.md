# D7024E Kademlia Laboration

---

This repository contains the code for the Kademlia lab in the D7024E course at the Luleå University of Technology.
Contributors: Filip Renberg, Oskar Lundqvist and Tovah Parnes.

## Prerequisites

Before you begin, ensure you have the following installed on your machine:

- Docker
- Make (most Unix-like systems have this pre-installed)

## Getting started

### Cloning the repository

To get started, you need to clone the repository to your local machine.

```bash
git clone git@github.com:SpiffiRacoon/group-12-kadlab.git
```

### Setup

1. Ensure Docker is installed and running.
2. Navigate to the project directory.

### Make Commands

There exists 5 make commands to manage the project:

| Command      | Description                           |
| ------------ | ------------------------------------- |
| `make build` | Build the Docker image `kadlab`       |
| `make up`    | Build and run Kademlia network        |
| `make clean` | Stop and remove the Docker containers |
| `make test`  | Run all tests and show test coverage  |
| `make help`  | Print available make commands         |

## Usage

After running `make up`, the Kademlia network will be up and running. You can interact with the network using the available command-line interfaces (CLIs).
To do this you need to attach to the running container. You can do this by running dockers `attach` command in a terminal.

```bash
docker attach <container-name>
```

This will attach said terminal to the specified terminal. The existing cli commands are.
| Command | Description |
|----------------|-------------------------------------------------------|
| `put <value>` | Store the given value |
| `get <value>` | Retrieve the stored value |
| `print` | Print the routing table and the number of nodes |
| `exit` | Shut down the node |
| `help` | Print available CLI commands |

To run any of the above commands, simply type `<command>` in the terminal.

## Testing

Running the command `make test` will run all tests and show the test coverage. The code coverage is calculated using the tool `go test -cover`. This repository has a code testing coverage of 85.5% <!-- TODO: update test coverage -->

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
