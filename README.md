# env-injector

A simple tool to inject credentials into environment variables from AWS Systems Manager Parameter Store.


## Install

```bash
$ go get github.com/okzk/env-injector

```

## How to use


``` bash
# When your parameter store is configured as below,
$ aws ssm get-parameters --with-decryption --names prod.wap.DB_USER prod.wap.DB_PASSWORD
{
    "InvalidParameters": [],
    "Parameters": [
        {
            "Type": "String",
            "Name": "prod.wap.DB_USER",
            "Value": "scott"
        },
        {
            "Type": "SecureString",
            "Name": "prod.wap.DB_PASSWORD",
            "Value": "tiger"
        }
    ]
}


# Set empty environment valiables.
$ export DB_USER=
$ export DB_PASSWORD=

# And specify parameter name prefix.
$ export ENV_INJECTOR_PREFIX=prod.wap

# Then exec your command via env-injector.
$ env-injector env | grep DB_
DB_USER=scott
DB_PASSWORD=tiger
```

## DEBUG

Set `ENV_INJECTOR_VERBOSE=1`
