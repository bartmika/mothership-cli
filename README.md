# mothership-cli
The command line interface for communicating with your [mothership-server](https://github.com/bartmika/mothership-server).

## Installation

Get our latest code.

```bash
go install github.com/bartmika/mothership-cli@latest
```

## Usage

```
The purpose of this application is to provide sub-commands for users to perform actions on their account in remote mothership server.

Usage:
  mothership-cli [flags]
  mothership-cli [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  insert      Insert a single time-series datum
  insert_bulk Insert a multiple time-series datum
  login       Log into your account
  register    Create an account
  select      List time-series data
  version     Print the version number

Flags:
  -h, --help   help for mothership-cli

Use "mothership-cli [command] --help" for more information about a command.
```

## Sub-Commands Reference

### ``register``

**Details:**

```text
Connects to the mothership server and creates your account

Usage:
  mothership-cli register [flags]

Flags:
  -c, --company string      Your companies name. If personal use then leave blank
  -e, --email string        The email you want to associate with your account
  -f, --first_name string   Your first name to use in your account
  -h, --help                help for register
  -l, --last_name string    Your last name to use in your account
  -x, --password string     The password you want to use to protect your account
  -p, --port int            The port of to connect to the server (default 50051)
  -t, --timezone string     Your accounts timezone.` (default "America/Toronto")
```

**Example Usage:**

```bash
$GOBIN/mothership-cli register -c="United Earth Fleet" -e="aki@unitedearthfleet.gov" -f="Aki" -l="Kirasagi" -x="blue-earth-good-radam-bad" -p=50051 -t="America/Toronto"
```

**Example Results:**

```text
2021/08/15 23:09:33 You have been successfully registered. Please login to begin using the system.
```

### ``login``

**Details:**

```text
Connects to the mothership server and log into your account

Usage:
  mothership-cli login [flags]

Flags:
  -e, --email string      The email you want to associate with your account
  -h, --help              help for login
  -x, --password string   The password you want to use to protect your account
  -p, --port int          The port of to connect to the server (default 50051)
```

**Example Usage:**

```bash
$GOBIN/mothership-cli login -e="aki@unitedearthfleet.gov" -x="blue-earth-good-radam-bad" -p=50051
```

**Example Failure Results:**

```text
Email or password are incorrect
```

**Example Success Results:**

```text
2021/08/15 23:22:43 Successfully logged in. Your credentials are as follows, please run this in your console.

Access Token:
export MOTHERSHIP_CLI_ACCESS_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mjk2ODg5NjMsImlzcyI6ImRiNjM5ZjEwLWNhY2ItNDE2ZC05MWIzLWNjYzBiODFlNjlkMSJ9.tXyXPpJzWUUIMZEOc0MJMaXof3knaipA9IgTDQjDrFM

Refresh Token:
export MOTHERSHIP_CLI_REFRESH_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzAyOTM3NjMsImlzcyI6ImRiNjM5ZjEwLWNhY2ItNDE2ZC05MWIzLWNjYzBiODFlNjlkMSJ9.BrsamYa7E7H98Ys1Fw9mR701043srLphgtBU6I7fPvQ
```

Please be sure you run the two exports in your console.

### ``insert``

**Details:**

```text
Connect to the gRPC server and send a single time-series datum.

Usage:
  mothership-cli insert [flags]

Flags:
  -a, --access_token string    The JWT access token for your account. Leave blank to access environment variable.
  -h, --help                   help for insert
  -m, --metric string          The metric to attach to the TSD.
  -p, --port int               The port of our server. (default 50051)
  -b, --refresh_token string   The JWT refresh token for your account. Leave blank to access environment variable.
  -t, --timestamp int          The timestamp to attach to the TSD.
  -v, --value float            The value to attach to the TSD.  
```

**Example Usage:**

```bash
$GOBIN/mothership-cli insert -p=50051 -m="solar_biodigester_temperature_in_degrees" -v=50 -t=1600000000
```

**Example Results:**

```text
2021/08/15 23:41:58 Successfully inserted
```

### ``select``

**Details:**

```text
Connect to the gRPC server and return list of time-series data results based on a selection filter.

Usage:
  mothership-cli select [flags]

Flags:
  -a, --access_token string    The JWT access token for your account. Leave blank to access environment variable.
  -e, --end int                The end timestamp to finish our range
  -h, --help                   help for select
  -m, --metric string          The metric to filter by
  -p, --port int               The port of our server. (default 50051)
  -b, --refresh_token string   The JWT refresh token for your account. Leave blank to access environment variable.
  -s, --start int              The start timestamp to begin our range
```

**Example Usage:**

```bash
$GOBIN/mothership-server select --port=50051 --metric="solar_biodigester_temperature_in_degrees" --start=1600000000 --end=1725946120
```

**Example Results:**

```text
2021/08/15 23:50:57 Server Response:
dataPoints:{value:50 timestamp:{seconds:1600000000}}
```

## Contributing
### Development
If you'd like to setup the project for development. Here are the installation steps:

1. Go to your development folder.

    ```bash
    cd ~/go/src/github.com/bartmika
    ```

2. Clone the repository.

    ```bash
    git clone https://github.com/bartmika/mothership-cli.git
    cd mothership-cli
    ```

3. Install the package dependencies

    ```bash
    go mod tidy
    ```

4. In your **terminal**, make sure we export our path (if you haven’t done this before) by writing the following:

    ```bash
    export PATH="$PATH:$(go env GOPATH)/bin"
    ```

5. You are now ready to start the command line interface and begin contributing! (Don't forget to apply the environment variables as well)

    ```bash
    go run main.go serve
    ```

### Quality Assurance

Found a bug? Need Help? Please create an [issue](https://github.com/bartmika/mothership-cli/issues).


## License

[**ISC License**](LICENSE) © Bartlomiej Mika
