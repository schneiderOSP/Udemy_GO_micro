Question: Ch14 results in error on TEst microservices page (in Chromium-based browsers like Chromium or Vivaldi):

  Output shows here...
  
  Error: TypeError: Failed to fetch

When I save the log data from browsers dev tools I get:

       POST http://localhost:8080/ net::ERR_CONNECTION_RESET
(anonymous) @ (index):64

Is it sth about the docker/compose stuff? Or CORS?

In Firefox however this got logged all the time:
Cross-Origin Request Blocked: The Same Origin Policy disallows reading the remote resource at http://localhost:8080/. (Reason: CORS request did not succeed). Status code: (null).

But as the 	

mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},

allows all http at once, I'm not sure why this CORS-error comes up. 

I'm running the microservice on 80 in Docker and mapped to 8080 on the PC. And the web application works on 8081. Thus I used 8080 in the fetch command in the template as this would be called from the webapplication outside Docker which will reach inside Docker the port 80.
Do I need to configure sth more due to not using 80 but 8080 (microservice) and 8081 (websrv)? 
