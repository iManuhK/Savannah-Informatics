# Savannah-Informatics

Savannah-Informatics is a web application that manages customer and order data, and integrates with the Africa's Talking SMS API to send SMS notifications. It uses the Gin framework for routing, PostgreSQL for data storage, and integrates with Africa's Talking for SMS functionality.

## Installations

Before you run the application, make sure the following dependencies are installed:

Go: Go programming language (v1.18 or above).

PostgreSQL: Database to store customers and orders.

Gin Framework: HTTP web framework for Go.

To install Go and PostgreSQL, follow the installation guides available on their official websites.

## Dependencies

Install required dependencies using go get:

    go get -u github.com/gin-gonic/gin

For Africa's Talking SMS API, use:

    go get github.com/AfricasTalkingLtd/africastalking-go

OR

    go get . //to install all dependencies

## API Endpoints

### /customers

POST: Create a new customer.

Request Body Example:

    {
        "code": "rsfrshdd",
        "full_name": "ekreykrt",
        "phone": "72734"
    }

Response Example:

    {
        "cust_id": 1,
        "code": "rsfrshdd",
        "full_name": "ekreykrt",
        "phone": "72734"
    }

### /orders

POST: Create a new order.

Request Body Example:

    {
        "item": "snack",
        "amount": 70.1,
        "cust_id": 1
    }

Response Example:

    {
        "order_id": 2,
        "item": "snack",
        "time": "0000-01-01T00:00:00Z",
        "amount": 70.1,
        "cust_id": 1
    }

## Testing API

The application’s API endpoints can be tested using Postman. Make sure to provide the correct request body in JSON format as shown in the examples above.

## Unit Tests

The application includes tests to ensure that authentication, SMS functionality, and the main application behave as expected. These tests are located in the following files:

auth_test.go – Tests for the authentication middleware and OIDC (OpenID Connect) initialization.

sms_test.go – Tests for the SMS functionality using Africa’s Talking API.

main_test.go – Tests the overall application behavior, ensuring the routes and database interactions work as expected.

### Dependencies for Testing

Before running the tests, make sure you have the following dependencies installed:

Gin: For routing and handling HTTP requests during tests.

Testify: For assertions in tests.

GoMock: To mock external dependencies where required (e.g., SMS API).

To install the required dependencies for testing, use:

    go get -u github.com/stretchr/testify
    go get -u github.com/golang/mock/gomock

### Test Files

#### auth_test.go

This test file contains tests for the authentication logic, including the initialization of OIDC and the middleware that protects routes.

Key Tests:

TestInitOIDC_MissingEnv: Tests the case where required environment variables (such as CLIENT_ID, CLIENT_SECRET, and REDIRECT_URI) are missing, expecting the application to panic.

TestOIDCAuthMiddleware_InvalidToken: Tests the middleware that validates tokens, ensuring it handles missing or invalid tokens correctly.

#### sms_test.go

This test file ensures that SMS functionality works as expected with Africa’s Talking API. It includes tests for sending SMS and handling errors from the API.

Key Tests:

TestSendSMS_Success: Tests successful SMS sending when the correct parameters are passed.

TestSendSMS_Failure: Tests failure scenarios, such as invalid parameters or API errors.

#### main_test.go

The main_test.go file contains tests for the application as a whole. It ensures that the routes work correctly, the database interaction is successful, and that the API responds with expected values.

Key Tests:

TestCustomerCreation: Tests creating a new customer via the /customers endpoint.

TestOrderCreation: Tests creating a new order via the /orders endpoint.

TestProtectedRoute: Ensures that the protected routes require valid authentication and reject unauthorized access.

### Running Tests

To run all the tests for the project, use the following command:

    go test -v

This will run all the tests in the auth_test.go, sms_test.go, and main_test.go files and display detailed output.

If you wish to run tests for a specific file, use:

    go test -v auth_test.go
    go test -v sms_test.go
    go test -v main_test.go

Example Test Output

A successful test run will show output similar to:

    emmanuel@Emmanuels-Computer Savannah-Informatics % go test -v
    2024/11/19 14:27:22 Environment variables loaded from Render
    === RUN   TestRoutes
    2024/11/19 14:27:22 Loaded config: ClientID: mock-client-id, ClientSecret: mock-client-secret, RedirectURI: http://localhost:8080/oauth/callback
    --- PASS: TestRoutes (0.43s)
    === RUN   TestLoginRoute
    [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

    [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
    - using env:   export GIN_MODE=release
    - using code:  gin.SetMode(gin.ReleaseMode)

    [GIN-debug] GET    /login                    --> savannah%2ego.TestLoginRoute.func1 (3 handlers)
    [GIN] 2024/11/19 - 14:27:22 | 307 |     198.958µs |                 | GET      "/login"
    --- PASS: TestLoginRoute (0.00s)
    === RUN   TestGetCustomersRoute
    [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

    [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
    - using env:   export GIN_MODE=release
    - using code:  gin.SetMode(gin.ReleaseMode)

    [GIN-debug] GET    /customers                --> savannah%2ego.TestGetCustomersRoute.func1 (3 handlers)
    [GIN] 2024/11/19 - 14:27:22 | 200 |     250.542µs |                 | GET      "/customers"
    --- PASS: TestGetCustomersRoute (0.00s)
    === RUN   TestPostCustomersRoute
    [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

    [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
    - using env:   export GIN_MODE=release
    - using code:  gin.SetMode(gin.ReleaseMode)

    [GIN-debug] POST   /customers                --> savannah%2ego.TestPostCustomersRoute.func1 (3 handlers)
    [GIN] 2024/11/19 - 14:27:22 | 201 |     384.375µs |                 | POST     "/customers"
    --- PASS: TestPostCustomersRoute (0.00s)
    === RUN   TestPostOrdersRoute
    [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

    [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
    - using env:   export GIN_MODE=release
    - using code:  gin.SetMode(gin.ReleaseMode)

    [GIN-debug] POST   /orders                   --> savannah%2ego.TestPostOrdersRoute.func1 (3 handlers)
    [GIN] 2024/11/19 - 14:27:22 | 201 |        78.5µs |                 | POST     "/orders"
    --- PASS: TestPostOrdersRoute (0.00s)
    PASS
    ok      savannah.go     0.667s

## Running the Application

To run the application localy:

    go run main.go

This will:

Connect to the PostgreSQL database hosted remotely and accessed through an env variable (DATABASE_URI)
Serve HTTP requests on 0.0.0.0:8080.
Make sure you have the required PostgreSQL database and tables created, and set up the proper credentials for your database connection.

## Continous Integration/Continous Deployment (CI/CD)

This project uses GitHub Actions and Render for automated CI/CD workflows.

### GitHub Actions for CI

CI Workflow: The project includes a GitHub Actions workflow (.github/workflows/go.yml) that run automatically:

How to use:

Push changes to the main branch.
Monitor the workflow status in the GitHub Actions tab of the repository.

### Render for CD

CD Workflow: Render is used to host and deploy the application. Changes pushed to the main branch are automatically deployed to the Render service.

Configuration:

Deployment configuration is defined in the render.yaml file.
Environment variables for the deployment are managed via the Render dashboard or the render.yaml file.

How to deploy:

Push changes to the GitHub repository.
Render will build, upload, and deploy the latest version of the application.

## Screenshots

[Blank](image.png).
<img width="507" alt="image" src="https://github.com/user-attachments/assets/dbd2c51b-22af-4727-8e10-4042a5d66344">

## Sources

Tabnine

GitHub Copilot

https://chatgpt.com/

https://google.com/

https://www.postgresql.org/docs/

https://pkg.go.dev/go.chromium.org/luci/server/auth#pkg-overview

https://github.com/AfricasTalkingLtd/africastalking-go

https://stackoverflow.com/

https://openid.net/specs/openid-connect-core-1_0.html

## License

This project is licensed under the [MIT License](LICENSE).
