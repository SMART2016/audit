## Design
### High Level Design
![Alt text](github.com/SMART2016/audit/docs/System Designs-Audit Logs.jpg?raw=true "High Level Design")

#### Design Flow
- Internal service raise audit events in a specific form and publish to message broker (Kafka here)
  - Current pattern for the Logs are 
   >RequestId: 0, CurrentUser: admin,Role: admin, System: auth-service, Action: POST:/auth-service/v1/login, IP: [::1]:55533, Agent: PostmanRuntime/7.37.3, Time: 2024-04-08T09:03:03+05:30, Status: 200
  - **Design Notes** It is better to allow services write there audit logs in a WAL and asynchronously sync the WAL's to elastic search using some agent like fluentbit or filebeat.
  - Log rotation is enabled using logrun in the current service implementation and also , filebeat handles file rotation with itself.
- Audit service subscribing to the kafka topic named: `log_events_topic`, picks audit log events from the topic and does below operations `[event-log-reader.go]`
  - Normalises the event log messages using a configured pattern as per the source system that generated the logs. The source system can register ther log patterns, which will be used while normalizing the respective logs.
    - Using a `grok` library to handle normalization right now. `[log-normalizer.go]`
  - After normalising the log messages they are pushed to the event store `(elastic search)`, where the audit log events are indexed and stored for better querying. `[event-store-client.go]`
    - **Design Note:** Elastic search is scalable document store , where we can query documents efficiently based on there fields and also on different time ranges.
  - The Audit service has exposed api's to query the events logs by registered users base on their permissions.
    - [API definitions](#api_definitions)
    - Code flow is as below for API
   > main.go --> router --> middlewares.go --> audit-request-handler.go and auth-request-handler.go --> response-log-handler.go

### TODO's
- Refactor the code to be more readable.
- Use WAL and async event publish instead of using message publish to kafka
- Providing a user interface for accessing the even logs data.
- Archiving data properly to keep the event store light and also to be compliant with individual state laws.
- Enabling download capability for event logs.
- Implementing Authz and auth centrally using well know Oauth2 implementations.
  - Implement Policy based Authz for handling ABAC.
- Audit log format changes , we will still have to support older logs and transform them to the new format , we cannot leave them or remove them.How much old log we need to keep depends on different compliance factors per domain or line of business.SO need to think on an effective way to handle that , right now i am normalising data to a specific format and this may work as well , as long as we can map incoming logs to this format.



### Running the service

- It is assumme that docker and docker compose is installed locally if not follow the link to install docker https://download.docker.com/linux/ubuntu or run the script in ./docs/install-docker.sh
- Navigate to the cloned repo for audit
    > cd audit
- Run below command to start the AUdit service and all dependent services.
    > docker compose up --build
  - --build flag will force to build the image of the audit-service every time, so if it is already done once , from next time just run below command to start all services locally.
      > docker compose up
#### Sending request to Audit Service
- To send request for loggin in to the Audit service and Fetch data, below are the detail information:
  - The service already has an admin user pre configured , so that we can create other users with the admin.
    - Default username : admin and password: admin.
  - Right now could not generate Swagger from the API's, but this link will provide definition of all API's
    - https://documenter.getpostman.com/view/5673453/2sA35MyJQV, open the link and there will be a button at the top right corner to open them with postman as below:
    ![Alt text](./docs/Screenshot 2024-04-08 at 11.22.25 AM.png?raw=true "High Level Design")
    - **Design NOTE:**  In ideal case the swagger definition needs to be generated before implementation , agreed upon between legit clients and then start the implementation.
 
#### <a name="api_definitions"></a>API definitions
- Link : https://documenter.getpostman.com/view/5673453/2sA35MyJQV
- API's at high level: **API's ate protected and Need user token**
  - GET /audit-service/v1/health
  - Submit log events query
    - POST /audit-service/v1/logevents  
  - /auth-service/v1/register
  - /auth-service/v1/users
  - /auth-service/v1/health
  - **Unprotected API doesn't need user token**:
    - - /auth-service/v1/login
