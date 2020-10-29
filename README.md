# dxdtmeter2

An HTTP server that counts the requests.

## Usage

1. Run the server, for instance in docker:

   ```shell
   docker build . -t dxdtmeter2
   docker run --rm -it --net=host -e ADDR=:8080 -e CONTROL_ADDR=:8081 dxdtmeter2
   ```

2. Create the load:

   We'll use [hey](https://github.com/rakyll/hey).

   ```shell
   go get github.com/rakyll/hey
   hey http://127.0.0.1:8080
   ```

3. Read the counter value:

   ```shell
   curl http://127.0.0.1:8081/get
   ```

4. Reset the counter value:

   ```shell
   curl http://127.0.0.1:8081/reset
   ```
