# What is rate Limiting?
#### Rate Limiting simply means limiting the the number of times a user can request a particular resource / service in a given time-frame. To ensure that all the clients of the service use the service optimally. The rate limiting is both a prevention and quality control measure.

### example
```

client := redis.NewClient(&redis.Options{
		DB:       0,
		Addr:     "redis:6379",
		Password: "",
	})
	if client == nil {
		fmt.Println("invalid connection")
		os.Exit(1)
	}
	ratelimit := ratelimit.New(context.Background(), ratelimit.Config{
		Redis:     client,
		MaxTokens: 200,
		Rate:      10,
		Duration:  time.Minute,
	})

	identifier := "Emad"
	bucket, err := ratelimit.GetBucket(context.Background(), identifier)
	if err != nil {
		fmt.Println("error: ", err.Error())
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 250)
		if !bucket.IsRequestAllowed(context.Background(), 10) {
			fmt.Printf("\n\n Request is not Allowed %d\n\n", i)
			// os.Exit(1)
		}
		if i%2 == 0 {
			bucket.DecreaseToken(context.Background(), 10)
		}
	}
	
```