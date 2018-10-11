package redisc

type Config struct {
    Host        string
    MaxActive   int
    MaxIdle     int
    IdleTimeout int
    DialConnectTimeout  int
    DialWriteTimeout    int
    DialReadTimeout     int
    DialDatabase        int
    DialKeepAlive       int
    DialPassword        string
}
