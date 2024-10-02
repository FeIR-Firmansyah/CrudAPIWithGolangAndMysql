# Backend API Documentation

This backend API is written in Go and uses MySQL as the database.

## Setup Instructions

### Before Running the Application

1. **Create MySQL Connection**

   Create a connection with the following details:
   - **User**: `root`
   - **Password**: `123`
   - **HostName**: `127.0.0.7`
   - **Port**: `3306`

2. **Ensure Go Development Environment is Ready**

   Make sure you have Go installed and set up on your machine.

---

## API Endpoint Documentation

### `/animal`

#### GET Method
- **Description**: Retrieve all animal data.
- **Additional Body Data**: None.

#### POST Method
- **Description**: Insert new animal data into the database.
- **Additional Body Data**: 
  ```json
  {
      "name": "<string>",
      "class": "<string>",
      "legs": <integer>
  }

#### PUT Method
- **Description**: Update the data on the specified id, if the data dont exist create a new one.
- **Additional Body Data**: 
  ```json
  {
    "id": <integer>,
    "name": "<string>",
    "class": "<string>",
    "legs": <integer>
  }

#### DELETE Method
- **Description**: Used to delete data with specified id.
- **Additional Body Data**: 
  ```json
  {
    "id": <integer>
  }

### `/animal/<id>`

#### GET Method
- **Description**: retrieves the details of a specific animal identified by its provided ID
- **Additional Body Data**: None.