## 安装

- 从编码编译

    git clone git@github.com:childe/ldap-nginx-golang.git
    cd ldap-nginx-golang
    make

- 或者下载编译好的二进制文件

从[https://github.com/childe/ldap-nginx-golang/releases/tag/201607](https://github.com/childe/ldap-nginx-golang/releases/tag/201607)下载对应的版本

## 使用

1. 配置config.json

    参考 config.example.json, 配置项参考后面的ldap验证原理

2. 配置nginx

        cp nginx.conf /etc/nginx
        nginx -s reload

    配置可参考后面的nginx原理

3. 运行

        ./nginx-ldap-auth-daemon --config config.json

    所有参数都可以在运行时指定, 会覆盖config.json里面的值, 如下:

        ./nginx-ldap-auth-daemon --host 0.0.0.0 --port 9000 --insecureSkipVerify true

    **useSSL和insecureSkipVerify两个参数只能在运行时指定, 不能写在config.json里面, 因为我不知道golang里面怎么处理bool类型的options参数的默认值, 没办法和config.json里面的值做合并. 用interface好像也可以做, 但太麻烦了.**

### 查看所有参数
    ./nginx-ldap-auth-daemon --help


## 原理

### nginx

依赖auth_request这个模块, 但这个模块默认是不安装的, 需要编译nginx的时候加上--with-http_auth_request_module这个参数.

官方文档在[http://nginx.org/en/docs/http/ngx_http_auth_request_module.html](http://nginx.org/en/docs/http/ngx_http_auth_request_module.html)

简单解释一下:

    location / {
        auth_request /auth-proxy;

        proxy_pass http://backend/;
    }

这个意思是说, 所有访问先转到/auth-proxy这里, /auth-proxy如果返回401或者403, 则访问被拒绝; 如果返回2xx, 访问允许,继续被nginx转到http://backend/; 返回其他值, 会被认为是个错误.

### ldap

ldap的验证步骤为:

1. 连接ldapserver
2. 根据 binddn, bindpw Bind到ldapserver (相当于登陆吧)
3. 把用户填写的用户名代入到filter模板中, 去ldap搜索
4. 用搜索到的DN去bind, 成功即验证成功


### 处理流程

参见[http://ohmycat.me/nginx/2016/06/28/nginx-ldap.html](http://ohmycat.me/nginx/2016/06/28/nginx-ldap.html)

## 感谢

nginx的[一篇官方博客](https://www.nginx.com/blog/nginx-plus-authenticate-users/)已经给出了非常详细的ldap认证办法, 并给出了[示例代码](https://github.com/nginxinc/nginx-ldap-auth)

我只是用golang实现了一下, 并做了精简.
