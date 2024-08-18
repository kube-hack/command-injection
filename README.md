# Overview

This is an educational resource demonstrating a web server running in a Kubernetes cluster that has a command injection vulnerability. Here are a few suggestions for how you might use this repository:

1. Practice your hacking skills by getting the contents of the `flag_to_capture` file, which is stored somewhere on the web server.

2. Read the source code found in `main.go` and the solution guide found in `solution/README.md` to better-understand what command injection vulnerabilities look like, how to exploit them, and how to prevent them.

3. Use this as a guide/inspiration for building your own applications with vulnerabilities.

The instructions and solutions were written assuming you are using some kind of Linux distribution (sorry Windows :grimacing:), whether Ubuntu, MacOS, or another Debian-based OS.

## \*\*\*\*\**DISCLAIMER*\*\*\*\*\*

This is an application with a built-in security vulnerability. Please don't deploy the Helm chart into a production environment. There are also instructions showing how to exploit command injection vulnerabilities, so please don't use this to break any laws :grin:.

# Usage

## Requirements

- Latest version of [Helm](https://helm.sh/docs/intro/install/)
- Latest version of [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
- A fully-compliant Kubernetes distribution (i.e. microk8s, k3s, k3d) that is running on Linux/amd64, and is using containerd or Docker as the runtime.

## Deploying to Kubernetes

Add the Helm chart repository:

```sh
helm repo add kube-hack https://kube-hack.github.io/charts
```

Update the charts in your Helm repository:

```sh
helm repo update
```

Deploy the chart to your Kubernetes cluster:

```sh
helm install command-injection kube-hack/command-injection
```

## Interacting with the application

### Port-forward the application

```sh
kubectl port-forward svc/web-server-command-injection 3000:3000
```

After the application is port-forwarded (accessible via localhost), you can open your browser and enter `localhost:3000` to access the web application.

### Using the application

This is a web application that sends a `ping` request to a url on the internet, and prints the response on the page. Enter a website (i.e. google.com) into the field next to the `Enter a URL or Domain:` label and click the `Ping` button to see the response.

### Validating the value of the flag

If you have found the `flag_to_capture` file, send a `curl` request to `localhost:3000/validate` with the contents of the file as a string in the request payload:

```sh
curl --data-binary "your-string-here" localhost:3000/validate
```
