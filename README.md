
# Demo for A Distrubuted Application in Go

This project aims to provide a hands-on exploration of core concepts in distributed applications for educational purposes.


## Current Implemented Features:

 **Service Registration:**

   - When a service starts up, it registers itself with the service registration service. This registration typically includes information such as the service's name, IP address, port.
   - The service registration service keeps track of all registered services and their details.
   
   
   - When a service is about to shut down, it sends a deregistration request to the service registration service to remove its entry.
   - The service registration service updates its records to reflect the service's shutdown, ensuring accurate information is available to other services.



## Acknowledgments

Special thanks to [Mike Van Sickle](https://www.linkedin.com/in/vansimke/) for his exceptional teaching on Pluralsight.
