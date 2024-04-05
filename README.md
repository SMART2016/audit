## Run and Install Instructions
- This has been tested on mac and can run on linux too.
- Install Docker and docker compose for this too work.
- just navigate to audit folder where the git repo is cloned and run below command
  - docker compose up
- That is it to run all components locally.
- Then Hit localhost:8080/health to hit the rest API , where we can generate events 
  - It is a simple user service with which you can create/read/update/delet user and these will be the events that will be logged for the current user who is performing the action.
  - The idea is to answere Who, When, Where and What thats all.