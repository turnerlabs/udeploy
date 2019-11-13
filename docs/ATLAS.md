### Setup Atlas via Terraform ###

1. Register an API key in Atlas `Atlas -> Project -> Access management -> API Keys`
    
    1. Click `Manage -> Create API Key`
    2. Save the public key for later use 
    3. Give the key `Project Owner` permission
    4. Click `Next`
    5. Save the private key for later use

    Optional

    6. Click `Add Whitelist Entry`
    7. Click `Use Current IP Address`
    8. Click `Done`

1. Apply changes
    ```bash
    $ terraform init -var-file=infrastructure/portals/prod/atlas/terraform.tfvars infrastructure/portals/prod/atlas
    $ terraform apply -var-file=infrastructure/portals/prod/atlas/terraform.tfvars infrastructure/portals/prod/atlas
    ```
1. Terraform will prompt for DB configuration
1. Save the output variables `DB_URI` and `DB_NAME` for later use