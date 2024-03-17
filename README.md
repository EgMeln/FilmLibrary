# Film Library API

Film Library API is a RESTful API for managing films in the film library. It allows users to perform CRUD operations on actors and movies, as well as user registration and login.


## Usage

Once the server is up and running, you can access the API endpoints using any HTTP client, such as Postman or cURL. The base URL for the API is http://localhost:10011.

## Endpoints

The API provides the following endpoints:

- **POST /register:** Register a new user with a username and password.
- **POST /login:** Log in an existing user with a username and password.
- **POST /actors/create:** Create a new actor in the film library.
- **PUT /actors/update:** Update an existing actor in the film library.
- **DELETE /actors/delete:** Delete an actor from the film library by ID.
- **GET /actors/getAllWithMovies:** Retrieve all actors from the film library along with their associated movies.
- **POST /movies/create:** Create a new movie with the provided details.
- **PUT /movies/update:** Update an existing movie with the provided details.
- **DELETE /movies/delete:** Delete an existing movie by its ID.
- **GET /movies/getAllWithSorting:** Retrieve all movies with sorting based on the provided flag.
- **GET /movies/getByTitleFragment:** Retrieve movies that match the provided title fragment.
- **GET /movies/getByActorNameFragment:** Retrieve movies associated with actors whose name matches the provided fragment.

For detailed information about the request and response formats, please refer to the Swagger documentation.