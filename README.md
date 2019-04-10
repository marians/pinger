[![Docker Repository on Quay](https://quay.io/repository/marian/pinger/status "Docker Repository on Quay")](https://quay.io/repository/marian/pinger)

# pinger

pinger is a very primitive heartbeat checker. It simply tries periodically to
reach an HTTP(S) URL. You can decide to take action if the last successful
response happened too long ago.

- Sends periodic HTTP request to configured URLs
- Keeps the timestamp of the last successful request
- Provides an HTTP endpoint to list successful results

## Deployment

A Docker image automated build is available as `quay.io/marian/pinger` ([details](https://quay.io/repository/marian/pinger?tab=info)).

In addition to the pinger image, Redis is required.

The `kubernetes` folder contains example manifests to create a deployment,
configmap (contianing the URLs to check), service, and ingress.

Also the `docker-compose.yaml` file can be used as a local starting point
if you don't want to do Kubernetes.

## Configuration of URL checks

The check YAML file expects an array of objects containing the following keys:

- `url` (mandatory): The URL to check
- `method` (optional): The HTTP method to use. Defaults to `GET`.

**Hint**: Consider using method `HEAD` instead of `GET` to prevent avoidable
data transfer.

See the provided `config.yaml` for an example.

### HTTP endpoint output

The HTTP endpoint for listing results only lists URLs that have been checked
successfully within the last 24 hours. The response looks like this:

```json
[
  {
    "url": "https://example.com/",
    "method": "GET",
    "last_success": "2019-04-10T17:56:46+0000"
  }
]
```

Here, `url` and `method` are the original values from the check configuration
(see above). The `last_success` attribute shows the timestamp when a URL has
been last retrieved successfully.

**Hint:** If a configured check does not appear in the results list, it has
not been retrieved successfully in the last 24 hours. When in doubt, check the
pinger logs.
