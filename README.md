# dummy_service

This service produce cpu load on every cpu based on gotten value(%), but not less and not more than defined in config<br>
Also it has readiness and healthz methods and methods which can control this two. Just for experiments : >

### Methods:

| Name         | Http method | Params              | Description
|--------------|-------------|---------------------|-----------------------------------------------------------------------------|
| /set-cpu     | POST        | int target_cpu_load | Set target CPU utilization                                                  | 
| /show-cpu    | GET         | None                | Show current target CPU utilization, if not specified use workload.cpu.min  | 
| /healthz     | GET         | None                | Health check endpoint, receive 200 and "Ok" when app is fine                | 
| /readiness   | GET         | None                | Readiness endpoint, receive 200 and "Ok" if app is ready to handle requests | 
| /healthz-on  | GET         | None                | Enable /healthz                                                             | 
| /healthz-off | GET         | None                | Disable /healthz                                                            | 
| /ready       | GET         | None                | Enable /readiness                                                           | 
| /not-ready   | GET         | None                | Disable /readiness                                                          | 

Code which artificially produce cpu workload was borrowed from [here](|https://github.com/vikyd/go-cpu-load/blob/master/cpu_load.go)

###Example
curl -v -X POST -H "Content-Type:application/json" -d '{"target_cpu_load": 15}' 'http://localhost:8080/set-cpu'