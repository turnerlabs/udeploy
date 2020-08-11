## Upgrade (AWS) ##

If running version >= **v0.30.0-rc**, follow these steps to upgrade uDeploy.

### LATEST: v0.33.2-rc

1. Update configuration.

Update `infrastructure/portals/prod/.env` based on auth provider.

```bash
$ cd infrastructure/portals/prod
$ cstore pull -t prod
```

<details>
  <summary>Portal (OIDC) -> AzureAD</summary> 

```
OAUTH_SCOPES=openid,offline_access,email
```

</details>

<details>
  <summary>Portal (OIDC) -> Okta</summary> 

```
OAUTH_SCOPES=openid,email
```

</details>

```bash
$ cstore push -t prod
```

2. Update infrastructure.

`infrastructure/portals/prod/main.tf`

```
terraform {
  required_version = ">= 0.12.0"
}

provider "aws" {
  version = ">= 2.46.0"
}

module “prod” {
    source = "github.com/turnerlabs/udeploy//infrastructure/modules/portal?ref=v0.33.2-rc"
}
```

`infrastructure/portals/prod/terraform.tfvars`

```
image = "quay.io/turner/udeploy:v0.33.2-rc.18"
```

4. Apply changes.

```bash
$ cd infrastructure/portals/prod
$ terraform init
$ terraform apply
```