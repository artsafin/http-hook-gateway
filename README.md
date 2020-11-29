### Configuration

```
hhg.icinga.accept_url_regex=.*
hhg.icinga.proxy_host=https://some-icinga.com
hhg.icinga.proxy_path=/hooks/12345
hhg.icinga.request_file=/tmp/icinga.request
```

### REQUEST file format

```
POST /path/on/target/host
Content-Type: application/json

{
    "foo": 1,
    "bar": 2
}
```

### Input

```
POST /icinga

{"foo": 1, "bar": 2}
```

### Output

```
POST https://some-icinga.com/hooks/12345

{"foo": 1, "bar": 2}
```
