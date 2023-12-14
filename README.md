# alert-system

```
docker run -p 9001:9001 quay.io/chronosphereiotest/interview-alerts-engine:latest
```

```
// Build an alert execution engine that:
// - 1. Queries for alerts using the query_alerts API and execute them at the specified interval
// - 2. Alerts will not change over time, so only need to be loaded once at start 
// - 3. The basic alert engine will send notifications whenever it sees a value that exceeds the critical threshold. 
// - 4. Add support for repeat intervals, so that if an alert is continuously firing it will only re-notify after the repeat interval.
// - 5. Incorporate using the warn threshold in the alerting engine - now an alert can go between states PASS <-> WARN <-> CRITICAL <-> PASS.
```