A basic check to see that the cameras are still running

# RUNNING
*With Binary*
camera-check \
	--smtp-server="$SMTP_SERVER" \
	--smtp-port="$SMTP_PORT" \
	--username="$SMTP_USER" \
	--password="$SMTP_PSWD" \
        --from="$FROM" \
        --to="$TO" \
	--last-alert-file="last_alert" \
	--minute-threshold=10

*With Golang*
1. Set the env variables:
```
export SMTP_SERVER=""
export SMTP_PORT=""
export SMTP_USER=""
export SMTP_PSWD=""
export FROM=""
export TO=""
```

2. run:
`./bin/run`


## TODO
add tests
log to file/log rotation
add in the actual camera check
