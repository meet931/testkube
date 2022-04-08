## kubectl-testkube watch testsuiteexecution

Watch test

### Synopsis

Watch test by test execution ID, returns results to console

```
kubectl-testkube watch testsuiteexecution <executionID> [flags]
```

### Options

```
  -h, --help   help for testsuiteexecution
```

### Options inherited from parent commands

```
      --analytics-enabled   enable analytics (default true)
  -c, --client string       client used for connecting to Testkube API one of proxy|direct (default "proxy")
  -s, --namespace string    Kubernetes namespace, default value read from config if set (default "testkube")
  -v, --verbose             show additional debug messages
```

### SEE ALSO

* [kubectl-testkube watch](kubectl-testkube_watch.md)	 - Watch tests or test suites
