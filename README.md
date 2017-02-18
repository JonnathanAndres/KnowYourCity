## KnowYourCity Campaign [Code4SierraLeone]
=============================================


**This branch mainly contains the main logic of the app that handles the Administrative Management part and can be used to both map the data as well manipulate and manage the data.**

Know Your City Campaign aims at mapping all the amenities within Informal Settlement areas in Freetown, Sierra Leone, so as to help the occupants get help in the case of a disaster in the regions. The mapping will give details in terms of directions and distance from their current locations as well as the services offered in those areas.


For the User Interface Design and Logic Implementation, I have used the following resources(please refer to the documentation(`to come`) for further details on each)
* `jQuery`
* `HTML` & `LessCSS`
* `Google Maps Javascript API v3`
* `JavaScript`
* `Golang` (Local server taste runs
* `Firebase`(Authentication and Realtime Database)

## Main Features
* Allow users to find mapped resource centers in the Informal Settlement Areas
* Allow users to have access to contacts to the mapped resource centers
* Allow users to  Identify their current location on the map, thus be able to trace the mapped locations

## Check it out at a glance

#### Home Page

![alt tag](https://raw.githubusercontent.com/Code4SierraLeone/KnowYourCity/base/assets/img/photos/13.png)

#### Web-app Page

![alt tag](https://raw.githubusercontent.com/Code4SierraLeone/KnowYourCity/base/assets/img/photos/12.png)

## Getting your own instance of Know Your City Campaign app

1. Clone this repository

   `$ https://github.com/Code4SierraLeone/KnowYourCity.git`

2. Install project dependencies via `go get`.

    `$ go get`

3. Run a development server from the root folder.

    `$ go run app/serve.go ` (This is for current tests)
    `$ go run app/main.go ` (This is for production stage)

## Data Mapping
* Collection and Implementation of data mapping.
* Included in the package is the Google Places Api which I am using to search for all the hospitals within a radius of 11000, this shows and maps any health center within that radius. Which is an amazing way of identifying amenities around the user. Thank you Google Places.

```javascript
	map = new google.maps.Map(document.getElementById('map'), {
	          center: freetown,
	          zoom: 14
	        });

	        infowindow = new google.maps.InfoWindow();
	        var service = new google.maps.places.PlacesService(map);
	        service.nearbySearch({
	          location: freetown,
	          radius: 11000,
	          type: ['store']
	        }, callback);
	      }

	      function callback(results, status) {
	        if (status === google.maps.places.PlacesServiceStatus.OK) {
	          for (var i = 0; i < results.length; i++) {
	            createMarker(results[i]);
	          }
	        }
	      }

	      function createMarker(place) {
	        var placeLoc = place.geometry.location;
	        var marker = new google.maps.Marker({
	          map: map,
	          position: place.geometry.location
	        });

	        google.maps.event.addListener(marker, 'click', function() {
	          infowindow.setContent(place.name);
	          infowindow.open(map, this);
	        });
```

* To avoid data crowding, I have decided to use the marker clustering tool, which allows me to group the markers and as the user zoom in they get to reveal each marker and thus be able to view the information on that marker.

```javascript
        // Create an array of alphabetical characters used to label the markers.
        
        var labels = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';

        // Add some markers to the map.
        // Note: The code uses the JavaScript Array.prototype.map() method to
        // create an array of markers based on a given "locations" array.
        // The map() method here has nothing to do with the Google Maps API.
        
        var markers = locations.map(function(location, i) {
          return new google.maps.Marker({
            animation: google.maps.Animation.DROP,
            position: location,
            label: labels[i % labels.length]
          });
        });
        

        // Add a marker clusterer to manage the markers.
        var markerCluster = new MarkerClusterer(map, markers,
            {imagePath: 'https://developers.google.com/maps/documentation/javascript/examples/markerclusterer/m'});

```

## Issues

* Need to figure out to use the cluster tool on the newly identified `google search places`
* Broken side-menu on the web-app page. Working to fix this ASAP **[FIXED]**

## Milestone/Backlog
* Fully integrate the SMS feature into the app to allow for quick notifications


* To view the current user interface designs just run the files on the `templates directory` on your local server. You can as well set up the `$GOPATH` and run the `main.go` server found here[This is under construction and will be documented by next week].
* The data being rendered on the User interface at the moment are dummy data, and not real especially on the Map.

## Did You Know...

#### Its mobile device ready? Yes it is.

![alt tag](https://raw.githubusercontent.com/Code4SierraLeone/KnowYourCity/base/assets/img/photos/11.png)

