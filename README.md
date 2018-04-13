# env-injector

A simple tool to inject credentials into environment variables from AWS Secrets Manager
and/or Systems Manager Parameter Store.

## Install

Download from [releases](https://github.com/okzk/env-injector/releases).

## How to use

### Using meta config


``` bash
# When your secrets manager and parameter store are configured as below,
$ aws secretsmanager get-secret-value --secret-id prd/db1 --query SecretString --output text
{"user":"alice","password":"foo"}
$ aws secretsmanager get-secret-value --secret-id prd/db2 --query SecretString --output text
{"user":"bob","password":"bar"}
$ aws ssm get-parameters-by-path --with-decryption --path /prod/wap
{
    "Parameters": [
        {
            "Type": "SecureString",
            "Name": "/prod/wap/SOME_OTHER_CONFIG",
            "Value": "hoge"
        }
    ]
}

# And meta config yaml is stored as below, 
$ aws ssm get-parameter --name /meta/prd/wap --query Parameter.Value --output text
- secret_name: prd/db1
  env_prefix: db1
  capitalize: true
- secret_name: prd/db2
  env_prefix: db2
  capitalize: true
- parameter_store_path: /prod/wap


# Then specify meta config,
$ export ENV_INJECTOR_META_CONFIG=/meta/prd/wap

# and exec your command via env-injector.
$ env-injector env 
DB1_USER=alice
DB1_PASSWORD=foo
DB2_USER=bob
DB2_PASSWORD=var
SOME_OTHER_CONFIG=hoge
```

### Injecting form Secrets Manages

``` bash
# When your secrets manager is configured as below,
$ aws secretsmanager get-secret-value --secret-id prd/db --query SecretString --output text
{"DB_USER":"scott","DB_PASSWORD":"tiger"}

# And specify your secret name
$ export ENV_INJECTOR_SECRET_NAME=prd/db

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
                "secretsmanager:GetSecretValue"
            ],
            "Resource": [
                "arn:aws:secretsmanager:ap-northeast-1:123456789012:secret:prd/db-*"
            ]
        }
    ]
}
``` 


### Injecting form Parameter Store
You can use hierarchical parameters and/or grouped parameters.

#### Injecting hierarchical parameters

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


#### Injecting grouped parameters

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

### 

## DEBUG

Set `ENV_INJECTOR_VERBOSE=1`
