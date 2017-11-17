# env-injector

A simple tool to inject credentials into environment variables from AWS Systems Manager Parameter Store.


## Install

Download from [releases](https://github.com/okzk/env-injector/releases).

## How to use

You can use hierarchical parameters and/or grouped parameters.

### Injecting hierarchical parameters

``` bash
# When your parameter store is configured as below,
$ aws ssm get-parameters-by-path --with-decryption --path /prod/wap
{
    "Parameters": [
        {
            "Type": "String",
            "Name": "/prod/wap/DB_USER",
            "Value": "scott"
        },
        {
            "Type": "SecureString",
            "Name": "/prod/wap/DB_PASSWORD",
            "Value": "tiger"
        }
    ]
}

# And specify parameter name path
$ export ENV_INJECTOR_PATH=/prod/wap

# Then exec your command via env-injector.
$ env-injector env | grep DB_
DB_USER=scott
DB_PASSWORD=tiger
```

Required IAM role policy is as follows:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ssm:GetParametersByPath"
            ],
            "Resource": [
                "arn:aws:ssm:ap-northeast-1:123456789012:parameter/prod/wap"
            ]
        }
    ]
}
``` 


### Injecting grouped parameters

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

Required IAM role policy is as follows:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ssm:GetParameters"
            ],
            "Resource": [
                "arn:aws:ssm:ap-northeast-1:123456789012:parameter/prod.wap.*"
            ]
        }
    ]
}
``` 


## DEBUG

Set `ENV_INJECTOR_VERBOSE=1`
