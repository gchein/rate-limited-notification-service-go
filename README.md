# Rate Limited Notification Service API - Golang Implementation

This repo is a Golang implementation of [the first Rails implementation](https://github.com/gchein/rate-limited-notification-service)

## Installation and Configuration

1. **Clone the Repository (below using SSH) and enter the local directory:**
   ```sh
   git clone git@github.com:gchein/rate-limited-notification-service-go.git
   cd rate-limited-notification-service-go
   ```

2. **Install dependencies:**
    ```sh
    go mod tidy
    ```

3. **Create a '.env' file with the necessary variables on the root of the repo**

    The application makes use of some variables that exist on a '.env' file on the root of the repo. The only variable that is absolutely necessary is the 'DB_PASSWORD' with the correct password to access your local MySQL installation, for creating the DSN to access the database. So make sure to set that at least. The other variables will default to the below:

    ```
      ### .env file
      DB_PASSWORD=SET_THIS # Required

      # Defaults
      PUBLIC_HOST=http://localhost
      PORT=8080
      DB_USER=root
      DB_HOST=127.0.0.1
      DB_PORT=3306
      DB_NAME=rlnotif
      TEST_DB_NAME=rlnotif_test
    ```

4. **Setup the Database:**

    Each command that is related to the database call for a password by default. As the API was built using MySQL, the command calls the mysql CLI commands, with the '-p' option for the password, for secure access. If your MySQL configuration does not require a password to be typed for access usign the CLI, you should not be prompted for a password.

   ```sh
    make db_create
    make db_migrate
    make db_seed
    ```

    The `db_seed` command will populate the database with 3 test users (initial IDs 1 to 3), and an example set of rate limits.

    If you need to reset the database, run the command below, then the commands in this section again

    ```sh
    make db_drop
    ```

5. **Start the server:**
    ```sh
    make run
    ```

## Interacting with the API

  The API has the below endpoints:

  ```sh
    POST    /notifications    # Sends notification to user, if no rate limits are disrespected
    GET     /rate-limits      # Fetches all current rate limits
    POST    /rate-limits      # Creates a new rate limit
    DELETE  /rate-limits/{ID} # Deletes a rate limit
  ```

  You can interact with the application by sending HTTP JSON requests to the endpoints in any way you prefer (cURL, Postman etc.), while passing the necessary attributes.

  How to interact with each endpoint:

  **1. POST /notifications**

  This endpoint is the main entry point into the application. It will receive the notification parameters in the request body, and try to create the given notification, for the given user. If no rate limits are exceeded, then the notification is sent to the user (simulated as a JSON return with HTTP Status OK and the message in the body). If the notification would disrespect the rate limits, or if the request contains an error, you get back a response with the corresponding error.

  Below are the fields that are required on the request body:

  - notificationType (string)
  - userId (integer)
  - message (string)

  Here's an example using cURL (assuming the server is running on the default port 8080):

  ```sh
    curl -X POST http://localhost:8080/notifications \
    -H "Content-Type: application/json" \
    -d '{"notificationType":"Status Update", "userId":1, "message":"Print message if successful"}'

    {"message":"Print message if successful"}
  ```

  **2. GET /rate-limits**

  This endpoint basically fetches all the rate limits that are currently on the database.

  **3. POST /rate-limits**

  This endpoint is used if you want to create any specific rules you want. It returns with the instance of rule created, if successful.

  Each rate limit consists of a determined max limit of notifications to be sent in a given time window, for a given notification type. For instance, if you wanted to create a limit of maximum 2 notifications per hour on the "Project Update" notification type, you would use the command below.

  Below are the fields that are required on the request body:

  - notificationType (string)
  - timeWindow (string)
  - maxLimit (integer)

  Example using cURL:

  ```sh
    curl -X POST http://localhost:8080/rate-limits \
    -H "Content-Type: application/json" \
    -d '{"notificationType":"Project Update", "timeWindow":"Hour", "maxLimit":2}'

    {"id":XX,"notificationType":"Project Update","timeWindow":"Hour","maxLimit":2}
  ```

  **4. DELETE /rate-limits/{ID}**

  Finally, the DELETE endpoint should be used only if you want to delete a specific rate limit from the DB. The usage is pretty simple, you only need to send a DELETE request with a valid id (integer) as the end of the request path (without the brackets).

## Running the test suite

To run the test suite, first you need to run at least once the command below to create the test database. This command wil reset and migrate the test database. After you run it once, you do not need to run it again, as the tests are made as transactions to the DB that are rolled back, so this database should always be empty.

  ```sh
  make test_db_prepare
  ```

After you ran at least once the command above for preparing the database, run the command below to run the actual tests:

  ```sh
  make test
  ```
