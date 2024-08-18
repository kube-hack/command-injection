# What is command injection?

Command injection is a vulnerability that allows the execution of unauthorized commands on the host of a web-facing application. This means that someone who exploits a command injection vulnerability has a lot of options:
- Retrieving Kubernetes service account token, which they could use to run `kubectl` commands.
- Viewing environment variables that could contain API keys, passwords, etc.
- Altering the web server's static and executable files to redirect users to a malicious website.
- Accessing other services running on the same network as the web server host.

# How to identify a command injection vulnerability

If application spawns a shell and passes user input into the shell command, this makes it vulnerable to command injection:

```go
cmdString := fmt.Sprintf("ping -c 2 %s", userInput)
cmd := exec.Command("sh", "-c", cmdString)
```

In the above example, there is nothing stopping a user from passing in any command they would like to run, including `google.com; rm --no-preserve-root -rf /`, which ends the `ping` command and deletes every file on the host.

# How to guard against command injection

You should always use a built-in Golang library to run commands on the host. Rather than running `ping` via `exec.Command`, we might consider using this library instead: [pro-bing](https://github.com/prometheus-community/pro-bing).

If there is a software running on the host that doesn't have a built-in Golang library, you should never spawn a shell to run commmands. Instead, run commands directly with the user input as arguments:

```go
cmd := exec.Command("ping", "-c", "2", userInput)
```

This will prevent the user from executing any other commands other than `ping`; however, keep in mimd that this does not prevent them from adding unauthorized arguments or flags to a command.

# How to hack into the web server

Our objective is to retrieve the file contents of the `flag_to_capture` file, so we need to run a `find` command at the root. In the URL field of the web application, we'll send this string to the web server to get the contents of the file:

```sh
google.com; cat $(find / -type f -name 'flag_to_capture')
```

*****Note***: Each time the Kubernetes pod hosting the web server is restarted, the flag file's contents will be different, and will be randomly placed in a different directory.