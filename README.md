# HTTP hook gateway

This project creates a gateway between third-party services that can reformat an incoming HTTP request to the format accepted by one or many target HTTP services.
HHGW allows to create a template of request to be sent to target service and interpolate it with data of incoming request including HTTP method, path, headers and body.
Service is intended to run in a trusted environment as it doesn't provide rate limiting and load balancing capabilities.

## Template variables

|Variable|Description|Example|
|-|-|-|
|Method|HTTP method|`{{ .Method }}`|
|RemoteAddr|Incoming request's IP address with port|`{{ .RemoteAddr }}`|
|Headers|Collection of HTTP headers. Please note that each header contains a list of values|`{{ index .Headers "Content-Type" 0 }}` or the same with helper function `{{ header . "Content-Type" }}`|
|Body|Parsed request body. It parses `application/json`, `application/x-www-form-urlencoded` and `multipart/form-data` request|`{{ .Body.event.message }}`|
|Scheme|http or https||
|User|User||
|UserPassword|User password||
|Path|Request path||
|Fragment|Request fragment||
|Query|Collection of URL query string key-value parameters|`{{ index .Query "thing_id" 0 }}`|

## Functions

|Function|Description|Example|
|-|-|-|
|json|Encodes it's argument to JSON|`{{ json .Body }}`, `{"foo": {{ .Body.some.key | json }} }`|
|query|||
|header|Takes first header from a passed template data||
|headervalues|||
|env|Fetches env variable|`{{ env "SHELL" }}`|

## Usage

### Configuration

```
hhgw.icinga.accept_url_regex=.*
hhgw.icinga.proxy_host=https://some-icinga.com
hhgw.icinga.proxy_path=/hooks/12345
hhgw.icinga.request_file=/tmp/icinga.request
```

### REQUEST file format

```
POST /path/on/target/host{{ .Path }}
Content-Type: application/json

{
    "proxied foo": {{ .Body.foo }},
    "full_body": {{ json .Body }}
}
```

### Input

```
POST /icinga

{"foo": 1, "bar": 2}
```

### Output

```
POST https://some-icinga.com/path/on/target/host/hooks/12345

{"foo": 1, "bar": 2}
```
