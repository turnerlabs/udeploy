### Setup Azure Active Directory ###

Follow these steps to setup Azure to be used as the Authentication API for the portal.

1. Register portal in Azure.  `Azure Active Directory -> App Registrations -> New registration`
    
    1. Enter `udeploy` as name.
    2. Click `Register`
    
2. Save `client` and `tenant` ids to use in the container configuration.

3. Under the `Authentication` menu option, add three Redirect URI's per domain. (dev,prod,etc...) 

    Replace `udeploy.com` with the Route53 portal domain.

    * https://udeploy.com/oauth2/response
    * https://udeploy.com/idpresponse
    * https://udeploy.com/apps

4. Under the `Certificates and Secrets` menu option, create a `New client secret` and save the secret to use in the container configuration.