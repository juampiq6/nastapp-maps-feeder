## Simple Go Program for Gas Station Discovery and Information Upload using Google Places API

**Functionality:**

1. **Region Subdivision:** The program first divides the region into smaller squares based on a user-defined "differential" parameter (measured in radians). Each square is represented by two points: its starting point and the diagonal endpoint.
2. **Gas Station Discovery:** The program queries for gas stations within each subregion and stores their IDs. The "differential" parameter significantly affects the number of queries required, making its selection crucial.
3. **ID Persistence:** All discovered gas station IDs are saved to the "discoveredPlaces.csv" file. This allows for resuming the process later, even if interrupted by the user.
4. **Detailed Information Acquisition:** Once all IDs are identified, the program retrieves detailed information about each gas station using the Places API details query.
5. **Data Upload:** Finally, the program uploads the gathered gas station information to a MongoDB database, making it readily accessible for map services.

**Running the Program:**

```bash
go run *.go [--flags]
```

For a list of available flags and their default values, use:

```bash
go run *.go --help
```

**Available Flags:**

* `--discover`: Controls entry into the gas station discovery phase (default: true)
* `--details`: Controls entry into the detailed information acquisition phase (default: true)
* `--lastIndex`: Defines the index to continue processing from (useful for resuming after interruption) (default: 0)
* `--differential`: Sets the differential value used for region subdivision (default: defined in constants.go)
