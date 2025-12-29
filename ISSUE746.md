### [\#746] HTTP delivery target
**URL**: https://github.com/foxcpp/maddy/issues/746
**State**: OPEN
**Labels**: mta-in, 

**Body:**
maddy should be able to forward messages to an arbitrary HTTP endpoint.

```
deliver_to http {
  tls_client { ... }
  header Authorization "Bearer ..."
  format rfc822/headers/json/form
  field "body" "{{ header }}" # ???
}
```

## Possible operation modes

### message/rfc822 body

* Content-Type is set to `message/rfc822`.
* The entire message is copied into HTTP body.

## message/rfc822-headers

* Content-Type is set to `message/rfc822-headers`
* Message header is copied into HTTP Body.

## JSON

* Content-Type is set to `application/json`.
* Populated values are defined in configuration.

## form-encoded

* Content-Type is set to `application/x-www-form-urlencoded`
* Populated values are defined in configuration.