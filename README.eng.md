## IMPORTANT

**the project is just a demo which is for learning. I never run it in production.**

## install

-  build from source code

    git clone git@github.com:childe/ldap-nginx-golang.git
    cd ldap-nginx-golang
    make

- download runnable bin file

download it here [https://github.com/childe/ldap-nginx-golang/releases/tag/201607](https://github.com/childe/ldap-nginx-golang/releases/tag/201607)


## usage

1. make your config.json

    refer to config.example.json.  *ladp chapter below explains what the parameters mean.*
    
2.  nginx config

        cp nginx.conf /etc/nginx
        nginx -s reload

    *nginx chapter explains what they mean.*
    

3. run it

        ./nginx-ldap-auth-daemon --config config.json


## mechanism

### nginx
The project depends on auth_request nginx module , but the module is not install by default. You need compile nginx with `--with-http_auth_request_module`

refer to [http://nginx.org/en/docs/http/ngx_http_auth_request_module.html](http://nginx.org/en/docs/http/ngx_http_auth_request_module.html)

explanation:

    location / {
        auth_request /auth-proxy;

        proxy_pass http://backend/;
    }
the config means all requests would firstly be redirected to /auth-proxy. The request will be denied if /auth-proxy return 401 or 403;  request will go on to http://backend/ if /auth-proxy return 2xx; Any other response code returned by the subrequest is considered an error.

### ldap

ldap auth check steps:

1. connect to ldapserver
2. bind(means login) ldapserver according to binddn & bindpw
3. replace username input by user to filter template , and use it to search in ldap
5. return true if username is found in ldap AND could bind the user with the password input by user.


## Thanks

One blog in nginx.com (https://www.nginx.com/blog/nginx-plus-authenticate-users/) has already givin detail method and also [example code](https://github.com/nginxinc/nginx-ldap-auth)

I just implement it with golang in a much simpler way (removed many direct steps)
