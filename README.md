## Design
### High Level Design
<p align="center">
  <img src="docs/System Designs-Audit Logs.jpg">
  <br/>
</p>

#### Design Flow
- Internal service raise audit events in a specific form and publish to message broker (Kafka here)
  - Current pattern for the Logs are
  >RequestId: 0, CurrentUser: admin,Role: admin, System: auth-service, Action: POST:/auth-service/v1/login, IP: 172.22.0.1:56148, Agent: PostmanRuntime/7.37.3, Time: 2024-04-08T09:03:03+05:30, Status: 200
  - **Design Notes** It is better to allow services write there audit logs in a WAL and asynchronously sync the WAL's to elastic search using some agent like fluentbit or filebeat.
  - Log rotation is enabled using logrun in the current service implementation and also , filebeat handles file rotation with itself.
- Current the Audit service itself will generate events as and when user performs some operation on the service via REST API.
- Audit service subscribing to the kafka topic named: `log_events_topic`, picks audit log events from the topic and does below operations `[event-log-reader.go]`
  - Normalises the event log messages using a configured pattern as per the source system that generated the logs. The source system can register ther log patterns, which will be used while normalizing the respective logs.
    - Using a `grok` library to handle normalization right now. `[log-normalizer.go]`
  - After normalising the log messages they are pushed to the event store `(elastic search)`, where the audit log events are indexed and stored for better querying. `[event-store-client.go]`
    - Elastic Search Index Schema:
        > ``
    - **Design Note:** Elastic search is scalable document store , where we can query documents efficiently based on their fields and also on different time ranges.
  - The Audit service has exposed api's to query the events logs by registered users base on their permissions.
    - [API definitions](#api_definitions)
    - Code flow is as below for API
  > main.go --> router --> middlewares.go --> audit-request-handler.go and auth-request-handler.go --> response-log-handler.go

#### Design Considerations
- **Cross Platform Deployement**
  - I have used containers for the system and deploying the system as a whole , Using container technologies like Docker, which encapsulate applications and their dependencies into a single container that can run on any system that supports the containerization platform, thereby abstracting away the underlying OS differences.
- **Safeguarding Against Audit Tampering**
  - Immutable Audit Logs , all event are new events written to the event store and only can be red via the expose api, no API to modify the event is exposed and by default event store used are inherently immutable.
  - Implemented access control right now on a crude level which can be made better with policy based ABAC and implementing access control on the event store level as well.
  - Secure Architecture and Isolation should isolated and deployed on seperate security boundaries and subnets where access to the system can be managed based on least privilage or zero trust principle.
  - backup and archival of the log events in the event store should be implemented.
- **Scalable Deployment**
  - The implementation is containerised and stateless , the system carrying state the event store is highly scalable and distributed datastore system.

### TODO's
- Refactor the code to be more readable.
- Use WAL and async event publish instead of using message publish to kafka
- Providing a user interface for accessing the even logs data.
- Archiving data properly to keep the event store light and also to be compliant with individual state laws.
- Enabling download capability for event logs.
- Implementing Authz and auth centrally using well know Oauth2 implementations.
  - Implement Policy based Authz for handling ABAC.
- Audit log format changes , we will still have to support older logs and transform them to the new format , we cannot leave them or remove them.How much old log we need to keep depends on different compliance factors per domain or line of business.SO need to think on an effective way to handle that , right now i am normalising data to a specific format and this may work as well , as long as we can map incoming logs to this format.
- Right now the password is sent as plain text , we should apply encryption and store the password, instead of sending and storing in plain text at all for more secure data.
- No unit tests added while unit tests should be present for any code that we write and that helps mkae the code more formatter and better and more moduler.
- We should expose the API via LB and API gateway where we can handle the Authz centrally and also scale the audit service efficiently.

### Running the service

- It is assumme that docker and docker compose is installed locally if not follow the link to install docker https://download.docker.com/linux/ubuntu or run the script in `./docs/install-docker.sh`
- Navigate to the cloned repo for audit
  > cd audit
- Run below command to start the AUdit service and all dependent services.
  > docker compose up --build
  - --build flag will force to build the image of the audit-service every time, so if it is already done once , from next time just run below command to start all services locally.
    > docker compose up
#### Sending request to Audit Service
- To send request for loggin in to the Audit service and Fetch data, below are the detail information:
  - The service already has an admin user pre-configured , so that we can create other users with the admin.
    - Default username : admin and password: admin.
  - Right now could not generate Swagger from the API's, but this link will provide definition of all API's
    - https://documenter.getpostman.com/view/5673453/2sA35MyeVe, open the link and there will be a button at the top right corner to open them with postman as below:
      <p align="center">
        <img src="docs/Screenshot 2024-04-08 at 11.22.25 AM.png">
        <br/>
      </p>
    - **Design NOTE:**  In ideal case the swagger definition needs to be generated before implementation , agreed upon between legit clients and then start the implementation.

#### <a name="api_definitions"></a>API definitions
- Link : https://documenter.getpostman.com/view/5673453/2sA35MyeVe
- The service should also be available on https://audit.smartlabs.site/
- API's at high level: **API's ate protected and Need user token**
  - GET /audit-service/v1/health
  - Submit log events query
    - POST /audit-service/v1/logevents
  - /auth-service/v1/register
  - /auth-service/v1/users
  - /auth-service/v1/health
  - **Unprotected API doesn't need user token**:
    - - /auth-service/v1/login

#### Testing The system
1. Login to the system using Admin user as below:
    > `curl --location 'localhost:8080/auth-service/v1/login' \
  --header 'Content-Type: application/json' \
  --data '{
           "username": "admin",
           "password": "admin"
  }'`
   - This will respond with a token, use that token to execute rest of the requests
2. Create a new user with `user` role as below:
    > curl --location --request POST 'localhost:8080/auth-service/v1/login' \
     --header 'Authorization: Basic YWRtaW46YWRtaW4='

3. Now use the token for the respective user `admin(Role: admin) or rohan(Role: user)` created above to fetch log events.
    > `curl --location 'http://localhost:8080/audit-service/v1/logevents' \
  --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzEyNTYxNjc4fQ.ZHRmRjVt4neKYw73AV-MiVkysv1rhfflanpd1aS-7gA' \
  --header 'Content-Type: application/json' \
  --data '{
            "type":"auth-service",  
            "es_query":{"query": {
            "range": {
            "Time": {
                  "gte": "now-48h",
                  "lte": "now"
            }
           }
         }
     }
  }
  '`
   - Note the `type` attribute is something that is mandatory in this case and is used to identify permissions for the current user.
   - Also except admin users can only see events that are generated because of their own action.