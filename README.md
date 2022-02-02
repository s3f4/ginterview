## Case study
--------------

### Usage

Clone the project
```sh
git clone https://github.com/s3f4/ginterview
```

In order to use locally, first go to **api** folder and create **.env** file with the following content.
```sh
MONGO_DSN=Your mongo dsn
```

 Project can be built  with the following command:
```sh
make build
```
or If docker is installed on the computer, the project can be build  with the following command. 
```sh
make up
```


## Build infrastructure
In order to build infrastructure on aws [Terraform](https://www.terraform.io/downloads.html) must be installed.

- Create aws iam with administrator permission
- Create a **terraform.tfvars** file with the following content in the infra folder 
  ```
  access_key = "your access key"
  secret_key = "your secret key"
  ```

and apply 
```sh
make apply
```

## Clean Up
To destroy created EC2 and the other resources run the following command.
```sh
make destroy
```

## Prod request test
```sh
make mongo-post-prod
make inmemory-post-prod
make inmemory-get-prod
```