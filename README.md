# go-maps
Geocode Project with Bloom Filter

This project is a Go application that queries geocoding data (addresses) using the Google Maps API. To improve performance, the application uses a Bloom Filter to check if a ZIP code has been queried before. If the ZIP code is not in the Bloom Filter, the application queries the Google Maps API, saves the result in MongoDB, and adds the ZIP code to the Bloom Filter. The project also includes a mock of the Google Maps API and MongoDB, all orchestrated using Docker Compose.

Project Structure
Services
1. Go Application
The Go application is the main service that exposes a /geocode endpoint. When a request is made to this endpoint with a ZIP code, the application:

Checks if the ZIP code is in the Bloom Filter.
If the ZIP code is in the Bloom Filter, it tries to fetch the address from MongoDB.
If the ZIP code is not in the Bloom Filter or not found in MongoDB, it queries the Google Maps API.
Saves the address in MongoDB and adds the ZIP code to the Bloom Filter.
Returns the address as a response.
2. MongoDB
MongoDB is used to store the queried addresses. When a new address is obtained from the Google Maps API, it is saved in MongoDB for future queries.