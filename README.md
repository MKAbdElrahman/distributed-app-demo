
  # Demo for A Distrubuted Application in Go

  This project aims to provide a hands-on exploration of core concepts in distributed applications for educational purposes.


  ## Current Implemented Features:

  **Service Registration:**

  - When a service starts up, it registers itself with the service registration service. This registration typically includes information such as the service's name, IP address, port and names of the dependant service names.
  - The service registration service keeps track of all registered services and their details.
  - When a service is about to shut down, it sends a deregistration request to the service registration service to remove its entry.
  - The service registration service updates its records to reflect the service's shutdown, ensuring accurate information is available to other services.



  **Service Discovery:**

  Service Discovery enables a web service to dynamically discover and connect to other dependent services without relying on static configurations. Here's how it typically works:

  - A web service, upon startup, queries the registry service for information about dependent services. This query could include parameters like the service name.
  - The registry responds to the query with information about the available services that match the specified criteria. This information may include the IP address, port, and other details needed to connect to the service.
  - To stay informed about changes in the service landscape, the registry may implement event notification mechanisms. When a new service registers or an existing service deregisters, the registry can broadcast these events to interested parties.




## Acknowledgments

Special thanks to [Mike Van Sickle](https://www.linkedin.com/in/vansimke/) for his exceptional teaching on Pluralsight.
