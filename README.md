# What is rate Limiting?
#### Rate Limiting simply means limiting the the number of times a user can request a particular resource / service in a given time-frame. To ensure that all the clients of the service use the service optimally. The rate limiting is both a prevention and quality control measure.

### example
```

client := redis.NewClient(&redis.Options{
		DB:       0,
		Addr:     "redis:6379",
		Username: "",
		Password: "",
	})
	if client == nil {
		fmt.Println("invalid connection")
		os.Exit(1)
	}
	ratelimit := ratelimit.New(context.Background(), ratelimit.Config{
		Redis:     client,
		MaxTokens: 100,
		Rate:      10,
		Duration:  time.Hour,
	})
	bucket := ratelimit.GetBucket(context.Background(), "Emad")
	fmt.Printf("\n%+v\n", bucket)

	if !bucket.IsRequestAllowed(5) {
		fmt.Println("Request is not Allowed")
		os.Exit(1)
	}
	
```