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
      DB_PASSWORD=<SET_THIS> # Required

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
    POST    /rate-limits      # Creates a new rate limits
    DELETE  /rate-limits/{ID} # Deletes a rate limit
  ```

  You can interact with the application by sending HTTP JSON requests to the endpoints in any way you prefer (cURL, Postman etc.), while passing the necessary attributes:
  - notification_type
  - user_id
  - message

  Here's an example using cURL (assuming the server is running on port 3000):

  ```sh
    curl -X POST http://localhost:8080/notifications \
    -H "Content-Type: application/json" \
    -d '{"notification_type":"Status Update", "user_id":"1", "message":"Print message if successful"}'
  ```
