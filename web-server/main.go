package main

import (
    "git.akyoto.dev/go/web"
)

func main() {
    // Create a new web server
    s := web.NewServer()

    // Define a simple GET route
    s.Get("/", func(ctx web.Context) error {
        return ctx.String("Hello, World!")
    })

    // Start the server on port 8080
    s.Run(":3000")
}
