# go-maps
Geocode Project

This project is a Go application that queries geocoding data (addresses) using the Google Maps API. The application queries the Google Maps API, saves the result in MongoDB, and returns the address as a response. The project also includes a mock of the Google Maps API and MongoDB, all orchestrated using Docker Compose.

## Project Structure
## Services

### 1. Go Application

The Go application is the main service that exposes a `/address` endpoint. When a request is made to this endpoint with a ZIP code, the application:

1. Tries to fetch the address from MongoDB.
2. If the address is not found in MongoDB, it queries the Google Maps API.
3. Saves the address in MongoDB.
4. Returns the address as a response.

### 2. MongoDB

MongoDB is used to store the queried addresses. When a new address is obtained from the Google Maps API, it is saved in MongoDB for future queries.

### 3. Mockoon

Mockoon is used to mock the Google Maps API for testing purposes. The mock server is configured using the `maps.json` file.

## Running the Project

### Prerequisites

- Docker and Docker Compose installed on your machine.

### Steps

1. **Clone the repository**:
   ```sh
   git clone https://github.com/amandamata/go-maps.git
   cd go-maps
   ```

2. **Update the .env**
```
GOOGLE_MAPS_API_KEY=dummy_key
GOOGLE_MAPS_API_URL=http://localhost:3001
```

3. **Run Docker Compose**:
```sh
docker-compose up
```

4. **Access the application**:
The Go application will be available at http://localhost:8080.
The mock Google Maps API will be available at http://localhost:3001.

5. **Example Request**
To query an address, make a GET request to the /address endpoint with a zipcode parameter:
```sh
curl --location 'http://localhost:8080/address?zipcode=10001'
```

Contributing
Feel free to contribute to this project by opening issues or submitting pull requests.

