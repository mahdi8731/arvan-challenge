# Coupon Management and Wallet Service

This GitHub project is a coupon management and wallet service designed to facilitate efficient handling of coupons and wallet transactions. It utilizes a combination of technologies including Postgres, Redis, and NATS for database management and messaging services.

## Technologies Used

- **Postgres**: Utilized for storing coupon and wallet data, providing a reliable and scalable database solution.
- **Redis**: Employed for caching and enhancing the performance of certain operations within the service.
- **NATS**: Used as a messaging system to facilitate communication between different components of the service.

## Third-Party Dependencies

The project integrates with the following third-party libraries:

- **PGX**: Chosen over Gorm for its superior performance and direct interaction with Postgres, enhancing efficiency in handling database operations.
- **go-redis**: Utilized for interacting with Redis, offering efficient caching mechanisms to optimize performance.
- **nats.go**: Selected over alternatives like Kafka or RabbitMQ for its lower latency, ensuring swift communication between service components.

## Additional Notes

- **Transactional Outbox Pattern**: The project employs the transactional outbox pattern for communication between services, ensuring reliability and consistency in distributed systems.
- **Monorepo Microservice**: This project follows a monorepo architecture, where multiple microservices are contained within a single repository, promoting code sharing and easier management of dependencies.

## Why PGX and NATS?

- **PGX over Gorm**: PGX is preferred over Gorm due to its ability to provide better performance when working with Postgres, offering more control and efficiency in database interactions.
- **NATS over Kafka or RabbitMQ**: NATS is chosen for its lower latency, enabling faster communication between components of the service, which is crucial for real-time processing of coupon and wallet transactions.

## Getting Started

To set up and run the project locally, just run this command first:

```
docker compose build
```

and then

```
docker compose up
```

## Contributing

Contributions to the project are welcome. To contribute, follow these guidelines:

- Fork the repository and create a new branch for your feature or fix.
- Make changes and ensure that they adhere to the project's coding conventions.
- Write tests for your changes to maintain code quality.
- Submit a pull request with a clear description of your changes and their purpose.

## License

This project is licensed under the [MIT License](LICENSE), allowing for open collaboration and usage with proper attribution.

## Contact

For any inquiries or support regarding the project, feel free to contact the project maintainers via [email](m.a1378.1387@gmail.com) or by opening an issue on the GitHub repository.
