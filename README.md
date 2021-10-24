# go-restapi-cache
A REST-API service that works as an in memory key-value store with [go-minimal-cache](https://github.com/HarunBuyuktepe/go-minimal-cache) library.                                                

The application is hosted on http://localhost:8080/ by the default. All requests method is GET. Example request is as follows,                                                                            

* /Get?key=XXX                                                                                                                              
* /Set?key=XXX&value=XXX                                                                                                                        
* /Flush                                                                                                                            
* /Delete?key=XXX                                                                                                                   
* /GetFrequency                                                                 
* /SetFrequency?Frequency=XXX                                                                               
* /GetPath                                                                              
* /SetPath?Path=XXX                                                                               
* /GetImageOfMemory                                                                                           
