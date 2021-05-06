# MariaDB Encryption
## Key Files
### Creating the Key file
```
openssl rand -hex 32 > encryption/keyfile
```
Add key identifier to the encryption key:
```
1;<hex_encoded_encryption_key>
```

### Encrypting the Key file
Generate a random encryption password:
```
openssl rand -hex 128 > encryption/keyfile.key
```
Encrypt the key file with the encryption password:
```
openssl enc -aes-256-cbc -md sha1 \
   -pass file:encryption/keyfile.key \
   -in encryption/keyfile \
   -out encryption/keyfile.enc
```

## Verifying MariaDB Encryption
### 1. By querying Information Schema
Open a mysql command-line client as root in the docker container (password is in docker-compose file):
```
docker-compose exec db mysql -u root -p
``` 
Use this query to see which tables are encrypted:
```
SELECT * FROM information_schema.innodb_tablespaces_encryption\G
```
### 2. By searching an existing string in a table
You can try to find an existing string with the grep command.
For example, if you know there is the string *Fabien* in the Users table, you can use this command:
```
docker-compose exec db find /var/lib/mysql/fp/users.ibd -type f -exec grep -i 'Fabien' {} +
```
If Users table is encrypted, there will be no output. Data are well encrypted, so it is impossible to read the file directly.

If Users table is not encrypted, there will be this message: 
*Binary file /var/lib/mysql/fp/users.ibd matches*.