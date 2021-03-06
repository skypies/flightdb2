package fr24

import "fmt"
import "time"

type Identifier struct {
	Fr24, Reg, Callsign, IATAFlightNumber  string
	Orig, Dest             string
	DepartureEpoch         int64
	DepartureTimeLocation  string
}

func (id Identifier)String() string {
	timeStr := ""
	if id.DepartureTimeLocation != "" {
		loc,err := time.LoadLocation(id.DepartureTimeLocation)
		if err != nil {
			timeStr = fmt.Sprintf("err with '%s': %v", id.DepartureTimeLocation, err)
		} else {
			t := time.Unix(id.DepartureEpoch, 0).In(loc)
			timeStr = fmt.Sprintf("%s", t)
		}
	}

	return fmt.Sprintf("%7s/%6s/%7s/%7s %3s-%3s %s",
		id.Fr24, id.IATAFlightNumber, id.Callsign, id.Reg, id.Orig, id.Dest, timeStr)
}

// Mechanical type definitions, to match the json responses.

// {{{ FlightPlayback

type FlightPlaybackResponse struct {
	Result struct {
		Response struct {
			Data struct {

				Flight struct {
					Airline struct {
						Name string `json:"name"`
						Code struct {
							Icao string `json:"icao"`
							Iata string `json:"iata"`
						} `json:"code"`
					} `json:"airline"`

					Identification struct {
						Hex string `json:"hex"`
						Id int `json:"id"`
						Number struct {
							Default string `json:"default"`
						} `json:"number"`
						Callsign string `json:"callsign"`
					} `json:"identification"`

					Aircraft struct {
						Identification struct {
							ModeS string `json:"modes"`
							Registration string `json:"registration"`
						} `json:"identification"`
						Model struct {
							Text string `json:"text"`
							Code string `json:"code"`
						} `json:"model"`
					} `json:"aircraft"`
					
					Track  []struct {
						Heading int `json:"heading"`
						Latitude float64 `json:"latitude"`
						Speed struct {
							Kts int `json:"kts"`
							Mph float64 `json:"mph"`
							Kph float64 `json:"kph"`
						} `json:"speed"`
						Squawk string `json:"squawk"`
						Altitude struct {
							Feet int `json:"feet"`
							Meters int `json:"meters"`
						} `json:"altitude"`
						Longitude float64 `json:"longitude"`
						Timestamp int `json:"timestamp"`
					} `json:"track"`

					Airport struct {
						Destination AirportData `json:"destination"`
						Origin AirportData `json:"origin"`
					} `json:"airport"`

				} `json:"flight"`
			} `json:"data"`

			Timestamp int `json:"timestamp"`
		} `json:"response"`

		Request struct {
			Format string `json:"format"`
			FlightId string `json:"flightId"`
			Callback string `json:"callback"`
		} `json:"request"`
		
	} `json:"result"`
}

type AirportData struct {
	Name string `json:"name"`
	Code struct {
		Icao string `json:"icao"`
		Iata string `json:"iata"`								
	} `json:"code"`
	Timezone struct {
		Offset int `json:"offset"`
		Abbr string `json:"abbr"`
		AbbrName string `json:"abbrName"`
		Name string `json:"name"`
		IsDst bool `json:"isDst"`
	} `json:"Timezone"`
	Position struct {
		Latitude float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Country struct {
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"country"`
		Region struct {
			City string `json:"city"`
		} `json:"region"`
	} `json:"position"`
}

/* http://mobile.api.fr24.com/common/v1/flight-playback.json?flightId=729a70e

{
    "result": {
        "response": {
            "data": {
                "flight": {
                    "airline": {
                        "name": "US Airways",
                        "code": {
                            "icao": "AWE",
                            "iata": "US"
                        }
                    },

                    "identification": {
                        "hex": "74c3ec3",
                        "id": 122437315,
                        "number": {
                            "default": "AA433"
                        },
                        "callsign": "AAL433"
                    },

                    "aircraft": {
                        "identification": {
                            "modes": "A71B43",
                            "registration": "N557UW"
                        },
                        "model": {
                            "text": "Airbus A321-231",
                            "code": "A321"
                        }
                    },

                    "track": [
                        {
                            "heading": 180,
                            "latitude": 33.4356,
                            "speed": {
                                "kts": 28,
                                "mph": 32.2,
                                "kmh": 51.9
                            },
                            "squawk": "0",
                            "altitude": {
                                "feet": 0,
                                "meters": 0
                            },
                            "longitude": -112.004,
                            "timestamp": 1441128357
                        },
                        {
                            "heading": 180,
                            "latitude": 33.4344,
                            "speed": {
                                "kts": 33,
                                "mph": 38,
                                "kmh": 61.1
                            },
                            "squawk": "0",
                            "altitude": {
                                "feet": 0,
                                "meters": 0
                            },
                            "longitude": -112.004,
                            "timestamp": 1441128378
                        },
                    ],

                    "airport": {
                        "destination": {
                            "name": "San Francisco International Airport",
                            "code": {
                                "icao": "KSFO",
                                "iata": "SFO"
                            },
                            "timezone": {
                                "offset": -25200,
                                "abbr": "PDT",
                                "abbrName": "Pacific Daylight Time",
                                "name": "America/Los_Angeles",
                                "isDst": true
                            },
                            "position": {
                                "country": {
                                    "name": "United States",
                                    "code": "US"
                                },
                                "latitude": 37.618969,
                                "region": {
                                    "city": "San Francisco"
                                },
                                "longitude": -122.374001
                            }
                        },
                        "origin": {
                            "name": "Phoenix Sky Harbor International Airport",
                            "code": {
                                "icao": "KPHX",
                                "iata": "PHX"
                            },
                            "timezone": {
                                "offset": -25200,
                                "abbr": "MST",
                                "abbrName": "Mountain Standard Time",
                                "name": "America/Phoenix",
                                "isDst": false
                            },
                            "position": {
                                "country": {
                                    "name": "United States",
                                    "code": "US"
                                },
                                "latitude": 33.434269,
                                "region": {
                                    "city": "Phoenix"
                                },
                                "longitude": -112.011002
                            }
                        }
                    }                    // airport
                }                        // flight
            },                           // data

            "timestamp": 1441134678
w        },                               // response

        "request": {
            "format": "json",
            "flightId": "74C3EC3",
            "callback": null
        }
    },                                   // result

    "_api": {
        "copyright": "Copyright (c) 2012-2015 Flightradar24 AB. All rights reserved.",
        "legalNotice": "The contents of this file and all derived data are the property of Flightradar24 AB for use exclusively by its products and applications. Using, modifying or redistributing the data without the prior written permission of Flightradar24 AB is not allowed and may result in prosecutions.",
        "version": "1.0.12"
    }
}

*/

// }}}
// {{{ CurrentDetails

// We want very little of this; hand-extract it via a jsonMap
type CurrentDetailsResponse struct {
	FlightNumber string
	Status string
	ScheduledDepartureUTC time.Time
	ScheduledArrivalUTC time.Time
	OriginTZOffset string
	DestinationTZOffset string
	ETAUTC time.Time
}

/*

{
  "flight":"BA287",
  "snapshot_id":null,
  "status":"landed",
  "dep_schd":1441976400,
  "arr_schd":1442015700,
  "departure":1441978943,
  "arrival":1442015551,
  "eta":1442015551,
  "from_iata":"LHR",
  "from_city":"London, London Heathrow Airport",
  "from_pos":[51.477501,-0.46138],
  "from_tz_code":"BST",
  "from_tz_offset":"1.00",
  "from_tz_name":"British Summer Time",
  "to_iata":"SFO",
  "to_city":"San Francisco, San Francisco International Airport",
  "to_pos":[37.618969,-122.374001],
  "to_tz_code":"PDT",
  "to_tz_offset":"-7.00",
  "to_tz_name":"Pacific Daylight Time",
  "airline":"British Airways",
  "aircraft":"Airbus A380-841",
  "airline_url":"http:\/\/www.flightradar24.com\/data\/airplanes\/baw-baw",
  "image":"http:\/\/doygv6stby4r3.cloudfront.net\/img\/6\/8\/7\/3\/47268_1439540378_tb.jpg?v=0",
  "imagelink":"http:\/\/external.flightradar24.com\/redir.php?url=http%3A%2F%2Fwww.jetphotos.net%2Fviewphoto_new.php%3Fid%3D8086765",
  "copyright":"Mike Barker",
  "imagesource":"Jetphotos.net",
  "image_large":"http:\/\/doygv6stby4r3.cloudfront.net\/image.php?params=QXNBS2duc01QcVkzWVVyV0tIWlJLYkRWbkJXNWx0WllpVTJsMjJMWWR0MkIvc2FtMm1rUTltZnN2c1dsMVZQRDVwS0NjRWJOSVpZSUwrWlNrakJ0aHd1Qi82Y0lBSWx2THJMcnRMRW44Y21WNnZDTFE1RlFDWWZDd0VycEpQSTc=",
  "imagelink_large":"http:\/\/external.flightradar24.com\/redir.php?url=http%3A%2F%2Fwww.jetphotos.net%2Fviewphoto_new.php%3Fid%3D8086765",
  "copyright_large":"Mike Barker",
  "q":"7648405",
  "trail":[37.6116,-122.389,0,37.6114,-122.389,0,37.6112,-122.389,0,37.6111,-122.388,0,37.6111,-122.388,0,37.611,-122.388,0,37.6109,-122.387,0,37.6106,-122.387,0,37.6105,-122.386,0,37.6103,-122.386,0,37.6101,-122.386,0,37.61,-122.385,0,37.6098,-122.385,0,37.6097,-122.385,0,37.6097,-122.384,0,37.61,-122.384,0,37.6103,-122.384,0,37.6105,-122.384,0,37.6108,-122.384,0,37.611,-122.383,0,37.611,-122.383,0,37.6109,-122.383,0,37.6109,-122.383,0,37.6111,-122.382,0,37.6114,-122.382,0,37.6116,-122.382,0,37.612,-122.382,0,37.6123,-122.382,0,37.6126,-122.381,0,37.6129,-122.381,0,37.6133,-122.381,0,37.6135,-122.381,0,37.6139,-122.381,0,37.6145,-122.38,0,37.6151,-122.38,0,37.6156,-122.379,0,37.616,-122.379,0,37.6162,-122.379,0,37.6165,-122.379,0,37.6168,-122.379,0,37.6172,-122.379,0,37.6175,-122.379,0,37.6179,-122.379,0,37.6183,-122.379,0,37.6186,-122.379,0,37.6189,-122.379,0,37.6192,-122.379,0,37.6194,-122.38,0,37.6196,-122.38,0,37.6199,-122.381,0,37.62,-122.381,0,37.6202,-122.381,0,37.6203,-122.382,0,37.6204,-122.382,0,37.6206,-122.382,0,37.6206,-122.383,0,37.6207,-122.383,0,37.6209,-122.383,0,37.6212,-122.383,0,37.6216,-122.383,0,37.6219,-122.383,0,37.6222,-122.383,0,37.6226,-122.383,0,37.6228,-122.383,0,37.6231,-122.382,0,37.6234,-122.382,0,37.6235,-122.382,0,37.6235,-122.381,0,37.6232,-122.38,0,37.6227,-122.379,0,37.6217,-122.377,0,37.6202,-122.373,0,37.619,-122.37,7.5,37.6179,-122.368,7.5,37.6157,-122.362,7.5,37.6133,-122.357,15,37.6117,-122.353,22.5,37.6092,-122.347,32.5,37.6077,-122.343,37.5,37.606,-122.339,45,37.6045,-122.336,50,37.6029,-122.332,55,37.6013,-122.328,62.5,37.5998,-122.325,67.5,37.5975,-122.319,77.5,37.5959,-122.315,82.5,37.5942,-122.311,90,37.5918,-122.305,97.5,37.5903,-122.302,105,37.5878,-122.296,112.5,37.5854,-122.29,122.5,37.5838,-122.286,127.5,37.5821,-122.282,135,37.5804,-122.278,140,37.5786,-122.274,147.5,37.5768,-122.27,155,37.575,-122.266,162.5,37.573,-122.261,170,37.5702,-122.254,180,37.5671,-122.247,190,37.5642,-122.24,200,37.561,-122.232,212.5,37.5592,-122.228,220,37.5561,-122.221,230,37.554,-122.216,237.5,37.552,-122.211,245,37.5499,-122.206,252.5,37.5477,-122.201,260,37.5455,-122.196,270,37.5437,-122.191,277.5,37.5415,-122.186,287.5,37.5395,-122.181,295,37.5375,-122.176,302.5,37.5354,-122.171,310,37.5334,-122.166,317.5,37.5311,-122.161,320,37.5286,-122.157,327.5,37.5253,-122.152,332.5,37.5221,-122.149,337.5,37.5185,-122.145,342.5,37.5148,-122.142,350,37.509,-122.136,362.5,37.5056,-122.133,370,37.5,-122.128,382.5,37.4942,-122.122,397.5,37.4906,-122.119,405,37.4827,-122.111,422.5,37.4786,-122.108,432.5,37.4745,-122.104,440,37.4707,-122.101,445,37.4662,-122.097,447.5,37.4616,-122.095,452.5,37.4566,-122.094,457.5,37.4486,-122.093,462.5,37.4437,-122.092,467.5,37.4389,-122.092,472.5,37.4192,-122.091,492.5,37.4139,-122.092,500,37.406,-122.096,510,37.401,-122.103,517.5,37.3954,-122.115,530,37.3881,-122.13,550,37.385,-122.136,557.5,37.3815,-122.143,567.5,37.3791,-122.148,575,37.3766,-122.154,582.5,37.3746,-122.161,590,37.3732,-122.172,605,37.3724,-122.186,622.5,37.3716,-122.201,632.5,37.3707,-122.217,640,37.3695,-122.236,652.5,37.3681,-122.266,670,37.3686,-122.273,675,37.3702,-122.282,682.5,37.3735,-122.29,695,37.3788,-122.298,705,37.3835,-122.303,705,37.4099,-122.311,715,37.4273,-122.306,727.5,37.441,-122.302,737.5,37.4586,-122.297,757.5,37.4652,-122.295,765,37.4888,-122.29,792.5,37.5026,-122.294,805,37.518,-122.303,820,37.5287,-122.31,850,37.5407,-122.318,895,37.5539,-122.326,937.5,37.5657,-122.333,972.5,37.5811,-122.343,1017.5,37.5872,-122.347,1035,37.5963,-122.353,1055,37.6027,-122.357,1070,37.6137,-122.368,1097.5,37.6247,-122.38,1105,37.6356,-122.391,1105,37.6464,-122.403,1105,37.657,-122.414,1105,37.6692,-122.427,1105,37.6812,-122.44,1105,37.692,-122.451,1105,37.703,-122.463,1105,37.714,-122.475,1105,37.7253,-122.487,1110,37.7404,-122.503,1122.5,37.7492,-122.513,1130,37.7622,-122.527,1140,37.7779,-122.544,1155,37.788,-122.554,1165,37.825,-122.594,1205,37.864,-122.636,1297.5,37.871,-122.643,1317.5,37.878,-122.651,1340,37.885,-122.659,1357.5,37.8927,-122.666,1372.5,37.9012,-122.672,1387.5,37.9092,-122.676,1402.5,37.9602,-122.691,1495,38.0117,-122.706,1585,38.0641,-122.721,1677.5,38.1717,-122.753,1867.5,38.1816,-122.756,1885,38.1915,-122.759,1902.5,38.2015,-122.762,1922.5,38.2119,-122.765,1937.5,38.2218,-122.765,1955,38.2679,-122.756,2052.5,38.3723,-122.726,2215,38.4853,-122.695,2412.5,38.6005,-122.662,2575,38.7163,-122.63,2775,38.8496,-122.592,3017.5,38.8595,-122.59,3035,38.8772,-122.585,3065,39.0019,-122.519,3295,39.1245,-122.448,3555,39.2497,-122.376,3830,39.3736,-122.304,3945,39.4956,-122.233,4000,39.6167,-122.163,4000,39.7337,-122.094,4000,39.8557,-122.022,4000,39.9787,-121.95,4000,40.0957,-121.88,4000,40.2192,-121.807,4000,40.3349,-121.738,4000,40.4592,-121.663,4000,40.5781,-121.592,4000,40.7023,-121.517,4000,40.7143,-121.51,4000,40.7272,-121.503,4002.5,40.7462,-121.49,4000,40.7992,-121.434,4000,40.8946,-121.323,4000,40.9956,-121.204,4000,41.0901,-121.093,4000,41.1866,-120.979,4000,41.2877,-120.86,4000,41.3829,-120.746,4000,41.4836,-120.626,4000,41.5842,-120.505,4000,41.685,-120.384,4000,41.7804,-120.268,4000,41.8798,-120.147,4000,41.9751,-120.031,4000,42.0702,-119.902,4000,42.1638,-119.771,4000,42.2531,-119.645,4000,42.347,-119.512,4000,42.436,-119.386,4000,42.5297,-119.252,4000,42.6216,-119.12,4000,42.7128,-118.989,4000,42.8021,-118.86,4000,42.896,-118.723,4000,42.9842,-118.594,4000,43.083,-118.449,4000,43.1729,-118.317,4000,43.2649,-118.181,4000,43.3516,-118.052,4000,43.4435,-117.914,4000,43.5318,-117.781,4000,43.6245,-117.641,4000,43.7116,-117.509,4000,43.8036,-117.369,4000,43.8918,-117.234,4000,43.9831,-117.093,4000,44.075,-116.951,4000,44.1665,-116.809,4000,44.2581,-116.665,4000,44.3491,-116.522,4000,44.4405,-116.378,4000,44.5338,-116.23,4000,44.6246,-116.085,4000,44.7096,-115.949,4000,44.7968,-115.808,4000,44.8886,-115.66,4000,44.9786,-115.513,4000,45.0611,-115.378,4000,45.1506,-115.231,4000,45.2361,-115.09,4000,45.3257,-114.941,4000,45.4133,-114.794,4000,45.4956,-114.656,4000,45.5868,-114.502,4000,45.6717,-114.358,4000,45.7556,-114.215,4000,45.8413,-114.068,4000,45.9289,-113.917,4000,46.013,-113.771,4000,46.1009,-113.618,4000,46.1868,-113.468,4000,46.2691,-113.323,4000,46.3567,-113.167,4000,46.4413,-113.017,4000,46.53,-112.858,4000,46.6196,-112.697,4000,46.6999,-112.551,4000,46.7843,-112.397,4000,46.8687,-112.243,4000,46.953,-112.087,4000,47.0537,-111.936,4000,47.158,-111.811,4000,47.2639,-111.683,4000,47.3775,-111.545,4000,47.4851,-111.413,4000,47.5898,-111.284,4000,47.6997,-111.149,4000,47.8039,-111.019,4000,47.9132,-110.883,4000,48.0123,-110.758,4000,48.1206,-110.621,4000,48.2256,-110.488,4000,48.3308,-110.353,4000,48.4513,-110.198,4000,48.5572,-110.061,4000,48.6652,-109.921,4000,48.7694,-109.785,4000,48.8791,-109.64,4000,48.9938,-109.488,4000,49.1243,-109.314,4000,49.2341,-109.167,4000,49.3395,-109.025,4000,49.4489,-108.876,4000,49.5533,-108.733,4000,49.6628,-108.583,4000,49.7726,-108.431,4000,49.8772,-108.285,4000,49.9799,-108.141,4000,50.0904,-107.985,4000,50.1975,-107.834,4000,50.3106,-107.675,4000,50.4227,-107.539,4000,50.5375,-107.4,4000,50.6512,-107.26,4000,50.7682,-107.116,4000,50.8853,-106.971,4000,50.9969,-106.832,4000,51.1143,-106.684,4000,51.2392,-106.526,4000,51.3579,-106.375,4000,51.474,-106.226,4000,51.5843,-106.084,4000,51.6982,-105.936,4000,51.8145,-105.783,4000,51.9315,-105.629,4000,52.0415,-105.483,4000,52.1567,-105.329,4000,52.272,-105.174,4000,52.3813,-105.026,4000,52.4971,-104.869,4000,52.6127,-104.71,4000,52.7302,-104.547,4000,52.8448,-104.388,4000,52.9545,-104.234,4000,53.0711,-104.069,4000,53.1785,-103.916,4000,53.2864,-103.762,4000,53.3999,-103.598,4000,53.5055,-103.444,4000,53.6173,-103.281,4000,53.728,-103.118,4000,53.8382,-102.954,4000,53.9524,-102.783,4000,54.062,-102.618,4000,54.1752,-102.446,4000,54.2825,-102.282,4000,54.3897,-102.117,4000,54.5032,-101.94,4000,54.6106,-101.772,4000,54.7135,-101.61,4000,54.8213,-101.439,4000,54.9282,-101.268,4000,55.0362,-101.094,4000,55.1426,-100.921,4000,55.2537,-100.739,4000,55.3596,-100.564,4000,55.823,-99.7835,4000,56.4179,-98.7494,4000,56.5174,-98.5641,4000,56.6149,-98.3801,4000,56.8956,-97.8624,4000,56.9984,-97.6672,4000,57.0877,-97.4931,4000,57.2857,-97.1275,4000,57.4788,-96.7542,4000,57.5709,-96.5626,4000,57.6777,-96.36,4000,57.8685,-95.9948,4000,57.9615,-95.8205,4000,58.1516,-95.4175,3920,58.4221,-94.7473,3800,58.5563,-94.3202,3800,58.8277,-93.7339,3800,58.9121,-93.5317,3800,59.0055,-93.3232,3800,59.1856,-92.9257,3800,59.2742,-92.7153,3800,59.4495,-92.3151,3800,59.529,-92.1321,3800,59.6213,-91.8974,3800,59.7965,-91.4808,3800,59.8841,-91.2882,3800,60.0486,-90.8765,3800,60.1362,-90.6654,3800,60.2201,-90.4487,3800,60.3054,-90.2409,3800,60.4623,-89.8296,3800,60.5515,-89.5781,3800,60.6255,-89.3685,3800,60.714,-89.1382,3800,60.868,-88.6886,3800,60.9509,-88.4844,3800,61.0331,-88.2351,3800,61.1811,-87.781,3800,61.2488,-87.5664,3800,61.33,-87.3278,3800,61.4834,-86.8796,3800,61.5551,-86.643,3800,61.6985,-86.1838,3800,61.7617,-85.9824,3800,61.8333,-85.7428,3800,61.9046,-85.5118,3800,62.0482,-85.0469,3800,62.1146,-84.8083,3800,62.1834,-84.5693,3800,62.254,-84.3252,3800,62.3946,-83.8503,3800,62.4708,-83.6051,3800,62.6044,-83.1177,3800,62.6657,-82.8813,3800,62.7342,-82.6239,3800,62.8624,-82.1297,3800,62.9186,-81.894,3800,62.9892,-81.6318,3800,63.1148,-81.1479,3800,63.1802,-80.8845,3800,63.2346,-80.6518,3800,63.3659,-80.1332,3800,63.4266,-79.8758,3800,63.4969,-79.6163,3800,63.6095,-79.1104,3800,63.6774,-78.8542,3800,63.7955,-78.3029,3800,63.8587,-78.0287,3800,63.9831,-77.4905,3800,64.049,-77.2024,3800,64.1668,-76.6541,3800,64.2227,-76.3691,3800,64.3191,-75.943,3800,64.3776,-75.7544,3800,65.1259,-72.2888,3800,65.1423,-72.1939,3800,65.1917,-71.9057,3800,65.2413,-71.6122,3800,65.2899,-71.3198,3800,65.3385,-71.0231,3800,65.3857,-70.73,3800,65.4329,-70.4327,3800,65.5247,-69.8723,3800,65.5247,-69.8723,3800,65.5915,-69.5215,3800,65.921,-66.9881,3800,65.9665,-66.7086,3800,66.0159,-66.3958,3800,66.1079,-65.807,3800,66.1596,-65.5113,3800,66.2519,-64.9346,3800,66.2981,-64.6087,3800,66.3331,-64.3355,3800,66.4153,-63.7067,3800,66.4583,-63.3866,3800,66.5027,-63.0796,3800,66.5868,-62.4194,3800,66.6246,-62.0978,3800,66.7062,-61.4256,3800,66.7447,-61.0851,3770,66.8198,-60.4783,3600,66.8525,-60.1815,3600,66.9158,-59.6877,3600,67.3949,-56.8238,3600,67.4731,-56.1116,3600,67.4894,-55.9586,3600,67.57,-55.18,3600,67.57,-55.18,3600,67.5788,-55.0924,3600,67.6145,-54.7309,3600,67.6474,-54.3898,3600,67.8072,-52.5949,3600,67.8072,-52.5949,3600,67.8072,-52.5949,3600,67.8356,-52.2475,3600,67.8356,-52.2475,3600,67.8566,-51.9843,3600,67.8966,-51.4668,3600,67.8972,-51.4581,3600,67.9243,-51.0932,3600,67.9501,-50.7342,3600,67.9749,-50.3748,3600,67.9985,-50.0262,3600,68.0103,-49.6565,3600,68.0243,-49.131,3600,68.0282,-48.9736,3600,68.0338,-48.7325,3600,68.0736,-45.8981,3600,66.6798,-28.689,3600,66.59,-28.3351,3600,66.5073,-28.0143,3600,66.4286,-27.7136,3600,66.3348,-27.3598,3600,66.2514,-27.0505,3600,66.1686,-26.747,3600,66.086,-26.4489,3600,66.0035,-26.1549,3600,65.9143,-25.8411,3600,65.8271,-25.5387,3600,65.7416,-25.2462,3600,65.656,-24.9573,3600,65.5694,-24.6688,3600,65.4825,-24.3828,3600,65.3947,-24.0974,3600,65.306,-23.8134,3600,65.2175,-23.5331,3600,65.1226,-23.2365,3600,65.0335,-22.9618,3600,64.942,-22.683,3600,64.8471,-22.3973,3600,64.75,-22.1088,3600,64.652,-21.8216,3600,64.5574,-21.5477,3600,64.4587,-21.2655,3600,64.3637,-20.9975,3600,64.2652,-20.7225,3600,64.1639,-20.4436,3600,64.067,-20.1771,3600,63.9674,-19.9165,3600,63.8649,-19.6577,3600,63.7574,-19.3899,3600,63.6504,-19.1266,3600,63.5472,-18.8753,3600,63.4401,-18.6177,3600,63.3353,-18.3687,3600,63.2217,-18.102,3600,63.1195,-17.8647,3600,63.0094,-17.6118,3600,62.8986,-17.3604,3600,62.7895,-17.1157,3600,62.6785,-16.8695,3600,62.5667,-16.6243,3600,62.4499,-16.3713,3600,62.3319,-16.1187,3600,62.2212,-15.8841,3600,62.1062,-15.6435,3600,61.9941,-15.4112,3600,61.8767,-15.1705,3600,61.7647,-14.9436,3600,61.1867,-13.8093,3600,61.1243,-13.6903,3600,60.9241,-13.3134,3600,60.703,-12.9045,3600,60.6315,-12.7742,3600,60.4159,-12.3857,3600,60.414,-12.3822,3600,60.3,-12.1796,3600,60.09,-11.8117,3600,59.9694,-11.6032,3600,59.952,-11.5732,3600,59.839,-11.3801,3600,59.7199,-11.1785,3600,59.6018,-10.9804,3600,59.4905,-10.7955,3600,59.3776,-10.6095,3600,59.2144,-10.3436,3600,59.1441,-10.23,3600,59.033,-10.0531,3600,58.9269,-9.86392,3600,58.8224,-9.67098,3600,58.7202,-9.48404,3600,58.6156,-9.29425,3600,58.4982,-9.08357,3600,58.3883,-8.88785,3600,58.2847,-8.70516,3600,58.1833,-8.52769,3600,58.0765,-8.34234,3600,57.9751,-8.16779,3600,57.8738,-7.99475,3600,57.7709,-7.82043,3600,57.6608,-7.63559,3600,57.556,-7.46087,3600,57.4538,-7.29226,3600,57.3527,-7.12641,3600,57.247,-6.95426,3600,57.1406,-6.78251,3600,57.0347,-6.61288,3600,56.9322,-6.45004,3600,56.8304,-6.2895,3600,56.7289,-6.13063,3600,56.6262,-5.97107,3600,56.5261,-5.81643,3600,56.4203,-5.65463,3600,56.3136,-5.49248,3600,56.2169,-5.34643,3600,56.1124,-5.18996,3600,56.0106,-5.03865,3600,55.9055,-4.88342,3600,55.8008,-4.72987,3600,55.6979,-4.58013,3600,55.5929,-4.42852,3600,55.4885,-4.27866,3600,55.3838,-4.12941,3600,55.2804,-3.98311,3600,55.1759,-3.83625,3600,55.0732,-3.69278,3600,54.9674,-3.54625,3600,54.8633,-3.40285,3600,54.7607,-3.26269,3600,54.6574,-3.12238,3600,54.5545,-2.98352,3600,54.4529,-2.84737,3600,54.3533,-2.71451,3600,54.2457,-2.57221,3600,54.1485,-2.44527,3600,54.1355,-2.42693,3600,54.125,-2.41221,3600,54.1123,-2.39575,3600,54.1023,-2.38492,3600,54.0122,-2.3228,3600,53.8891,-2.24451,3565,53.7602,-2.16301,3452.5,53.6332,-2.08344,3322.5,53.5153,-2.01017,3185,53.39,-1.93226,3077.5,53.265,-1.8417,2990,53.1476,-1.7322,2900,53.0317,-1.62582,2807.5,52.9214,-1.52663,2710,52.815,-1.43128,2602.5,52.8047,-1.42204,2592.5,52.7932,-1.41151,2580,52.7835,-1.40213,2572.5,52.7512,-1.36113,2540,52.6685,-1.22292,2392.5,52.5864,-1.08452,2277.5,52.5735,-1.06277,2260,52.5624,-1.04527,2242.5,52.5566,-1.03719,2232.5,52.5097,-0.988235,2170,52.4083,-0.893707,2025,52.3106,-0.803807,1875,52.2602,-0.75688,1770,52.2102,-0.709763,1662.5,52.16,-0.662918,1565,52.1102,-0.615616,1445,52.063,-0.570297,1340,52.0153,-0.524597,1230,51.9981,-0.508118,1185,51.9897,-0.499799,1165,51.9786,-0.489044,1142.5,51.9669,-0.478455,1125,51.9552,-0.470058,1115,51.9451,-0.464957,1107.5,51.925,-0.458601,1092.5,51.9059,-0.453339,1075,51.89,-0.449029,1060,51.8711,-0.443878,1040,51.8572,-0.44014,1012.5,51.8438,-0.436554,975,51.828,-0.432129,930,51.8095,-0.426908,877.5,51.8007,-0.4245,855,51.7937,-0.422592,835,51.7837,-0.419485,797.5,51.7761,-0.416219,767.5,51.7684,-0.4108,735,51.7599,-0.400927,702.5,51.7488,-0.383301,657.5,51.7373,-0.364685,620,51.7269,-0.348053,600,51.7143,-0.327835,600,51.7071,-0.31609,600,51.6934,-0.29396,600,51.6796,-0.27199,600,51.6649,-0.24868,600,51.6506,-0.22614,600,51.6315,-0.19597,600,51.6232,-0.18288,600,51.6138,-0.16891,602.5,51.6087,-0.16281,602.5,51.5853,-0.151433,597.5,51.5739,-0.153957,587.5,51.5629,-0.161603,565,51.5565,-0.168991,542.5,51.5517,-0.177043,522.5,51.547,-0.18959,505,51.5446,-0.199461,492.5,51.5413,-0.215644,457.5,51.5393,-0.225601,427.5,51.5371,-0.237097,395,51.5339,-0.253205,355,51.5318,-0.263895,327.5,51.5278,-0.279483,292.5,51.524,-0.288925,272.5,51.5193,-0.298233,245,51.5125,-0.311127,210,51.5057,-0.324097,195,51.5013,-0.332633,185,51.495,-0.344543,175,51.491,-0.352156,167.5,51.4873,-0.359268,162.5,51.4837,-0.366409,152.5,51.4784,-0.376816,145,51.4753,-0.383148,137.5,51.4721,-0.389569,125,51.4692,-0.396118,110,51.4669,-0.403061,95,51.4657,-0.410385,82.5,51.4653,-0.417852,70,51.4651,-0.429506,47.5,51.4651,-0.436401,32.5,51.4648,-0.478648,0,51.4659,-0.487728,0,51.4648,-0.486012,0,51.4648,-0.48645,0,51.4648,-0.487156,0,51.4652,-0.487834,0,51.4655,-0.487881,0,51.466,-0.487614,0],
  "first_timestamp":0}

 */

// }}}
// {{{ QueryResponse

// This is a small subset, only parsing out some interesting fields
type QueryResponse struct {
	Results []struct {
		Id string `json:"id"`

		Detail struct {
			Reg string `json:"reg"`
			Callsign string `json:"callsign"`
			Flight string `json:"flight"`
		} `json:"detail"`

		Type string `json:"type"`
	} `json:"results"`
}

/* http://www.flightradar24.com/v1/search/web/find?query=aal353&limit=8

{"results":
  [
    {"id":"87e6114",
     "label":"AA353 / AAL353 (N923NN)",
     "detail":{
       "reg":"N923NN",
       "callsign":"AAL353",
       "flight":"AA353"
     },
     "type":"live",
     "match":"begins"
   },
   {"id":"AA353",
    "label":"AA353",
    "detail":{
      "flight":"AA353"
    },
    "type":"schedule",
    "match":"begins"}
  ],
  "stats":{
    "total":{
      "all":2,
      "airport":0,
      "operator":0,
      "live":1,
      "schedule":1,
      "aircraft":0
    },
    "count":{
      "airport":0,
      "operator":0,
      "live":1,
      "schedule":1,
      "aircraft":0
    }
  }
}

*/

// }}}
// {{{ LookupHistoryResponse

// This is a subset, only parsing out some interesting fields
type LookupHistoryResponse struct {
	Result struct {
		Response struct {
			Data []struct {

				Identification struct {
					Id string `json:"id"`
					Number struct {
						Default string `json:"default"`
					} `json:"number"`
					Callsign string `json:"callsign"`
				} `json:"identification"`

				Time struct {
					Other struct {
						ETA int `json:"eta"`
						Updated int `json:"updated"`
					} `json:"other"`
					Scheduled struct {
						Departure int `json:"departure"`
						Arrival int `json:"arrival"`
					} `json:"scheduled"`
					Real struct {
						Departure int `json:"departure"`
						Arrival int `json:"arrival"`
					} `json:"real"`
					Estimated struct {
						Departure int `json:"departure"`
						Arrival int `json:"arrival"`
					} `json:"estimated"`
				} `json:"time"`
				
				Aircraft struct {
					Model struct {
						Text string `json:"text"`
						Code string `json:"code"`
					} `json:"model"`
					Registration string `json:"registration"`
				} `json:"aircraft"`

				Airport struct {
					Destination AirportData `json:"destination"`
					Origin AirportData `json:"origin"`
				} `json:"airport"`

			} `json:"data"`
		} `json:"response"`		
	} `json:"result"`
}

/* 

http://www.flightradar24.com/reg/n980uy
http://www.flightradar24.com/flight/aa1799

http://api.flightradar24.com/common/v1/flight/list.json?query=n980uy&fetchBy=reg
http://api.flightradar24.com/common/v1/flight/list.json?query=aa1799&fetchBy=flight

{
    "result": {
        "response": {
            "data": [

                    "airline": {
                        "name": "American Airlines",
                        "code": {
                            "icao": "AAL",
                            "iata": "AA"
                        }
                    },

                    "status": {
                        "ambiguous": false,
                        "text": "Landed 21:15",
                        "icon": "yellow",
                        "estimated": null,
                        "live": false,
                        "generic": {
                            "status": {
                                "color": "yellow",
                                "text": "landed",
                                "diverted": null,
                                "type": "arrival"
                            },
                            "eventTime": {
                                "local": 1452201356,
                                "utc": 1452230156
                            }
                        }
                    },

                    "identification": {
                        "codeshare": null,
                        "row": 2428550798,
                        "id": "87a0aa4",
                        "number": {
                            "alternative": "AA1799",
                            "default": "AA1799"
                        },
                        "callsign": "AAL1799"
                    },

                    "time": {
                        "other": {
                            "eta": 1452230156,
                            "updated": 1452230703
                        },
                        "scheduled": {
                            "departure": 1452207600,
                            "arrival": 1452228600
                        },
                        "real": {
                            "departure": 1452213397,
                            "arrival": 1452230160
                        },
                        "estimated": {
                            "departure": null,
                            "arrival": null
                        }
                    },

                    "aircraft": {
                        "model": {
                            "text": "Airbus A321-231",
                            "code": "A321"
                        },
                        "registration": "N552UW"
                    },

                    "airport": {
                        "destination": {
                            "visible": true,
                            "name": "San Francisco International Airport",
                            "code": {
                                "icao": "KSFO",
                                "iata": "SFO"
                            },
                            "timezone": {
                                "offset": -28800,
                                "abbr": "PST",
                                "abbrName": "Pacific Standard Time",
                                "name": "America/Los_Angeles",
                                "isDst": false
                            },
                            "position": {
                                "country": {
                                    "name": "United States",
                                    "code": "US"
                                },
                                "latitude": 37.618969,
                                "region": {
                                    "city": "San Francisco"
                                },
                                "longitude": -122.374001
                            }
                        },

                        "origin": {
                            "visible": true,
                            "name": "Charlotte Douglas International Airport",
                            "code": {
                                "icao": "KCLT",
                                "iata": "CLT"
                            },
                            "timezone": {
                                "offset": -18000,
                                "abbr": "EST",
                                "abbrName": "Eastern Standard Time",
                                "name": "America/New_York",
                                "isDst": false
                            },
                            "position": {
                                "country": {
                                    "name": "United States",
                                    "code": "US"
                                },
                                "latitude": 35.214001,
                                "region": {
                                    "city": "Charlotte"
                                },
                                "longitude": -80.9431
                            }
                        },
                        "real": null
                    }
                },
*/

// }}}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
