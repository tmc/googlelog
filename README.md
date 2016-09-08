googlelog
=========

cli tool to publish entries to google cloud logging

Example use
-----------

Look up project ID via:
```sh
$ gcloud projects list
```

```sh
$ googlelog -h
Usage of googlelog:
  -log string
       	log name
  -projectID string
       	project id
$ echo "test log line" | googlelog -projectID (id from above) -log test-log
```
