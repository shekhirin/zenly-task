# Problematic

At Zenly, we receive lots of geolocations that need to be corrected, enriched and published for the rest of the platform.

A GeoLocation is a Protobuf message defined as follows:
```Protobuf
message GeoLocation {
    double lat = 1;
    double lng = 2;
    google.protobug.Timestamp created_at = 3;
}
```

GeoLocations arrive in an incremental order, through a GRPC Stream.

## Your job

- A) Implement a GPRC service
- B) Enrich the received geolocations
- C) Monitoring

### The GPRC service

Your service should have 2 streams.

- *publish* to send N geolocations with an attached userID.
- *subscribe* to receive N geolocations from a list of userID given as parameters.

### Enrich

With the following code as example/pseudocode

```go
type GeoLocationEnriched struct {
    p *GeoLocation

    // sub messages for each enricher
    weather Weather
    personalPlace PersonalPlace
    transport Transport
}

// ---------------------------------------------
type Weather struct {
    temperature int
}

// ---------------------------------------------
type PersonalPlaceType = int
const (
    PERSONAL_PLACE_HOME PersonalPlaceType = iota
    PERSONAL_PLACE_WORK PersonalPlaceType = iota
    PERSONAL_PLACE_SCHOOL PersonalPlaceType = iota
)

type PersonalPlace struct {
    hws PersonalPlaceType
}

// ---------------------------------------------
type TransportType = int
const (
    TRANSPORT_CAR TransportType = iota
    TRANSPORT_TRAIN TransportType = iota
    TRANSPORT_PLANE TransportType = iota
)

type Transport struct {
    ttype TransportType
}

// ---------------------------------------------
type Enricher interface {
    Enrich(ctx context.Context, pe *GeoLocationEnriched)
}
```

*Publish* **should not be aware** of any internal details of each enricher
- it should just have an array of Enricher `[]Enricher`.
- it has only 100ms to call the maximum of enricher with the given geolocations.
- each enricher uses its own sub-message to store the data.
- publish the message GeoLocationEnriched on realtime bus like NATS.
- publish the message on a kafka topic.

*Each enricher*
 - should fill the sub-message with random data
 - sleeps between 1 and 100ms to simulate some I/Os

*Subscribe*
- it should listen on the realtime bus with the given userID.
- it should return those geolocations to the client.

### Monitoring

For this part you have to add a system of your choice to monitor what happens on the service.
Define which parts are critical and nested to be monitored.

## Contraints

Your service is to become a core piece of infrastructure: performance, scalability & reliability shouldn't be an afterthought.

There are is deadline, take all the time you need.
Feel free to ask *any* questions.


