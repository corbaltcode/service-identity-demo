# Service Identity Demo

## Setup

1. Set AWS credentials that have admin access to an account where the demo will be deployed.

2. In `terraform/terraform.tfvars`, set the unset variables.

3. Apply Terraform to create infrastructure:

        $ cd terraform
        $ terraform apply
        $ cd ..

4. In `ansible/hosts`, set the unset variables and enter the IP addresses of the instances created by Terraform.

5. Run Ansible playbooks to configure instances

        $ cd ansible
        $ ansible-playbook site.yml -i hosts
        $ cd ..
        
6. Add entries to `/etc/hosts` for `client`, `server`, and `spire-server` (public IPs).

## Demo

### 1. Intro: No Authentication

The echo client and server are programs that will be used to demonstrate SPIRE.

#### Run server

```
$ ssh ec2-user@server
$ ./demo/echo/server/server 0.0.0.0:8443
```

#### Access server as ec2-user

```
$ ssh ec2-user@client
$ ./demo/echo/client/client server:8443
```

#### Access server as evil-user

```
$ ssh ec2-user@client
$ su evil-user  # password is "password"
$ ./demo/echo/client/client server:8443
```

### 2. Node Attestation

The SPIRE agent is authenticated so that it can grant certificates to processes.

#### Create node registration entries

```
$ ssh ec2-user@spire-server
$ sudo systemctl start spire-server
$ /opt/spire/bin/spire-server entry create -node -spiffeID spiffe://corbalt.com/echo/client-node  -selector aws_iid:iamrole:arn:aws:iam::795973919855:role/service-id-demo-client
$ /opt/spire/bin/spire-server entry create -node -spiffeID spiffe://corbalt.com/echo/server-node  -selector aws_iid:iamrole:arn:aws:iam::795973919855:role/service-id-demo-server
```

#### Client node attestation

```
$ ssh ec2-user@client
$ cd /opt/spire
$ bin/spire-agent run
```

#### Server node attestation

```
$ ssh ec2-user@server
$ cd /opt/spire
$ bin/spire-agent run
```

### 3. Workload Attestation

A process can retrieve certificates signed by the SPIRE server if and only if it is authorized.

#### Failing workload attestation

```
$ ssh ec2-user@client
$ cd /opt/spire
$ bin/spire-agent api fetch
```

#### Create workload registration entries

```
$ ssh ec2-user@spire-server
$ /opt/spire/bin/spire-server entry create -parentID spiffe://corbalt.com/echo/client-node -spiffeID spiffe://corbalt.com/echo/client -selector unix:user:ec2-user
$ /opt/spire/bin/spire-server entry create -parentID spiffe://corbalt.com/echo/server-node -spiffeID spiffe://corbalt.com/echo/server -selector unix:user:ec2-user
```

#### Successful workload attestation

```
$ ssh ec2-user@client
$ bin/spire-agent api fetch
$ bin/spire-agent api fetch x509 -write .
$ openssl x509 -text -nocert -in svid.0.pem
```

#### Failing workload attestation with evil-user

```
$ ssh ec2-user@client
$ su evil-user  # password is "password"
$ bin/spire-agent api fetch
```

### 4. mTLS with SPIRE Client Library

Modified echo and client servers present and verify SPIRE certificates.

#### Run server

```
$ ssh ec2-user@server
$ ./demo/echo-spire/server/server 0.0.0.0:8443
```

#### Run client

```
$ ssh ec2-user@client
$ ./demo/echo-spire/client/client server:8443
```

#### Run client as evil-user

```
$ ssh ec2-user@client
$ su evil-user  # password is "password"
$ ./demo/echo-spire/client/client server:8443
```

### 5. mTLS with Ghostunnel Proxy

Unmodified programs can use SPIRE certificates by routing traffic through a proxy.

#### Show modified source

```
$ ssh ec2-user@client
$ cat demo/echo-spire/client/main.go
```

#### Run server

```
$ ssh ec2-user@server
$ ghostunnel server --use-workload-api-addr unix:///tmp/spire-agent/public/api.sock --listen 0.0.0.0:8443 --target localhost:8000 --allow-uri spiffe://corbalt.com/echo/client &
$ ./demo/echo/server/server localhost:8000
```

#### Run client

```
$ ssh ec2-user@<client-public-ip>
$ ghostunnel client --use-workload-api-addr unix:///tmp/spire-agent/public/api.sock --listen localhost:8000 --target server:8443 --verify-uri spiffe://corbalt.com/echo/server &
$ ./demo/echo/client/client localhost:8000
```