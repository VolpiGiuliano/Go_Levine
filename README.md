# Go Levine

The project aims to replicate the main function of an exchange: organizing the incoming orders and processing matching orders.

![1707598960074](image/README/1707598960074.png)

## Bugs

* [X] Fix again the duplicate problem

  ```
  Error decoding order book type message: gob: duplicate type received
  ```

  Fix: the encoder of the server had to have a pointer as an argument
  `err_mm := encoder.Encode(*OB)`

  But only after the third try the Client can read it
